package postgres

import (
	"blog-service/internal/db/mongo"
	mongomodels "blog-service/internal/db/mongo/models"
	postgresmodels "blog-service/internal/db/postgres/models"
	"blog-service/internal/db/testutil"

	"log"
	"os"

	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
)

var postgresClient *Client
var postgresTestArticle *mongomodels.ArticleDB

func setupMongoContainer(ctx context.Context) (testcontainers.Container, *mongo.Client, error) {
	// Create mongo container and defer its termination
	mongoContainer, port, err := testutil.CreateMongoContainer(ctx)
	if err != nil {
		return nil, nil, err
	}

	os.Setenv("MONGO_PORT", port) // set port env variable with obtained port

	// Instantiate new mongo client and defer disconnection
	mongoClient, err := mongo.NewMongoClient()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	return mongoContainer, mongoClient, nil
}

func setupPostgresContainer(ctx context.Context) (testcontainers.Container, *Client, error) {
	postgresContainer, port, err := testutil.CreatePostgreSQLContainer(ctx)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to start container: %s", err)
	}

	os.Setenv("PG_PORT", port) // set port env variable with obtained port

	// Instantiate new postgres client
	postgresClient, err := NewPostgresClient()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create postgres client: %w", err)
	}

	return postgresContainer, postgresClient, nil
}

func generateRandomComment() *postgresmodels.Comment {
	comment := &postgresmodels.Comment{
		Content:       testutil.GenerateRandomString(),
		UserID:        1,
		CommenterName: testutil.GenerateRandomString(),
		ArticleID:     postgresTestArticle.ID.String(),
		CreatedAt:     time.Now(),
	}

	return comment
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error

	var postgresContainer testcontainers.Container

	postgresContainer, postgresClient, err = setupPostgresContainer(ctx)
	if err != nil {
		log.Fatalf("failed to setup postgres container")
	}

	mongoContainer, mongoClient, err := setupMongoContainer(ctx)
	if err != nil {
		log.Fatalf("Failed to setup mongo container: %s", err)
	}

	postgresTestArticle, err = testutil.CreatePostgresTestArticle(ctx, mongoClient)
	if err != nil {
		log.Fatalf("failed to generate postgres test article: %s", err)
	}

	exitCode := m.Run()

	postgresClient.ConnPool.Close()

	if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}

	if err = mongoClient.Client.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	}

	if err := testcontainers.TerminateContainer(mongoContainer); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}

	os.Exit(exitCode)
}

func TestAddComment(t *testing.T) {
	ctx := context.Background()
	comment := generateRandomComment()

	id, err := postgresClient.CreateComment(ctx, comment)
	if err != nil {
		t.Errorf("Failed to add comment: %s", err)
	}

	addedComment, err := postgresClient.GetComment(ctx, id)
	if err != nil {
		t.Errorf("Error when getting comment: %s", err)
	}

	comment.ID = id

	// Set time values for comments to zero as a workaround
	// for differing between db and golang.
	comment.CreatedAt = time.Time{}
	addedComment.CreatedAt = time.Time{}

	require.Equal(t, comment, addedComment)
}

func TestDeleteComment(t *testing.T) {
	ctx := context.Background()

	comment := postgresmodels.Comment{
		Content:   testutil.GenerateRandomString(),
		UserID:    1,
		ArticleID: postgresTestArticle.ID.String(),
		CreatedAt: time.Now(),
	}

	id, err := postgresClient.CreateComment(ctx, &comment)
	if err != nil {
		t.Errorf("Failed to add comment: %s", err)
	}

	_, err = postgresClient.DeleteComment(ctx, id)
	require.NoError(t, err, "failed to delete comment: %s", err)

	retrievedComment, err := postgresClient.GetComment(ctx, id)
	require.Error(t, err, "Delete succesful, but comment is still present: %v", retrievedComment)
}

func TestGetComments(t *testing.T) {
	ctx := context.Background()
	nrOfComments := 10
	comments := make([]postgresmodels.Comment, 0, nrOfComments)

	for range nrOfComments {
		comment := generateRandomComment()
		id, err := postgresClient.CreateComment(ctx, comment)
		if err != nil {
			t.Errorf("Failed to add comment: %s", err)
		}
		comment.ID = id
		comments = append(comments, *comment)
	}

	addedComments, err := postgresClient.GetComments(ctx, postgresTestArticle.ID.String(), 10, 1)
	if err != nil {
		t.Errorf("Failed to get first 10 comments: %s", err)
	}

	// Iterate over retrieved comments and original comments
	// in different orders, and set CreatedAt time to zero
	// due to creation times being expectedly different
	for i := range addedComments {
		comments[nrOfComments-i-1].CreatedAt = time.Time{}
		addedComments[i].CreatedAt = time.Time{}
		require.Equal(t, comments[nrOfComments-i-1], addedComments[i])
	}
}

func TestAddLike(t *testing.T) {
	ctx := context.Background()

	comment := generateRandomComment()

	id, err := postgresClient.CreateComment(ctx, comment)
	if err != nil {
		t.Errorf("Failed to add comment: %s", err)
	}

	like := postgresmodels.Like{
		CommentID: id,
		UserID:    1,
	}

	id, err = postgresClient.AddLike(ctx, like)
	require.NoError(t, err, "failed to add like: %s", err)

	like.ID = id

	addedLike, err := postgresClient.FindLike(ctx, like.CommentID, like.UserID)
	require.NoError(t, err, "failed to retrieve like: %s", err)

	require.Equal(t, &like, addedLike)
}

func TestRemoveLike(t *testing.T) {
	ctx := context.Background()

	comment := generateRandomComment()

	id, err := postgresClient.CreateComment(ctx, comment)
	if err != nil {
		t.Errorf("Failed to add comment: %s", err)
	}

	like := postgresmodels.Like{
		CommentID: id,
		UserID:    1,
	}

	id, err = postgresClient.AddLike(ctx, like)
	require.NoError(t, err, "failed to add like: %s", err)

	like.ID = id

	_, err = postgresClient.RemoveLike(ctx, like.CommentID, like.UserID)
	require.NoError(t, err, "Failed to remove like: %s", err)

	addedLike, err := postgresClient.FindLike(ctx, like.CommentID, like.UserID)
	require.Error(t, err, "found like, but it should have been deleted: %v", addedLike)
}

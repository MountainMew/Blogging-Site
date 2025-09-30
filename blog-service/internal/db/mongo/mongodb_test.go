package mongo

import (
	"blog-service/internal/db/mongo/models"
	"blog-service/internal/db/testutil"
	"context"
	"log"
	"math/rand"
	"os"
	"testing"

	"time"

	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mongoClient *Client

func setupMongoContainer(ctx context.Context) (testcontainers.Container, *Client, error) {
	// Create mongo container and defer its termination
	mongoContainer, port, err := testutil.CreateMongoContainer(ctx)
	if err != nil {
		return nil, nil, err
	}

	os.Setenv("MONGO_PORT", port) // set port env variable with obtained port

	// Instantiate new mongo client and defer disconnection
	mongoClient, err := NewMongoClient()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	return mongoContainer, mongoClient, nil
}

func generateTestArticle() *models.ArticleDB {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec // Secure rng not needed for generating usernames
	categories := [...]string{"sports", "tv", "travel", "gaming"}
	randomContent := testutil.GenerateRandomString()

	content := randomContent
	title := randomContent + "_title"
	category := categories[seededRand.Intn(len(categories))]
	publisherID := seededRand.Intn(100)

	return &models.ArticleDB{
		Title:       title,
		Content:     content,
		Category:    category,
		PublisherID: publisherID,
	}
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	var mongoContainer testcontainers.Container

	mongoContainer, mongoClient, err = setupMongoContainer(ctx)
	if err != nil {
		log.Fatalf("Failed to setup mongo container: %s", err)
	}

	exitCode := m.Run()

	if err = mongoClient.Client.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	}

	if err := testcontainers.TerminateContainer(mongoContainer); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
	os.Exit(exitCode)
}

func TestInsertAndFind(t *testing.T) {
	ctx := context.Background()

	article := generateTestArticle()

	res, err := mongoClient.InsertArticle(ctx, article)
	if err != nil {
		t.Errorf("Failed to insert article: %s\n", err)
	}

	insertedID := res.InsertedID.(primitive.ObjectID)
	insertedArticle, err := mongoClient.FindArticleByID(ctx, &insertedID)
	if err != nil {
		t.Errorf("Failed to retrieve inserted article: %s\n", err)
	}

	if insertedArticle.Category != article.Category ||
		insertedArticle.Content != article.Content ||
		insertedArticle.PublisherID != article.PublisherID ||
		insertedArticle.Title != article.Title {
		t.Errorf("Article was inserted, but is corrupted\n")
		t.Errorf("Want: %v\nGot:%v\n", article, insertedArticle)
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()

	article := generateTestArticle()

	res, err := mongoClient.InsertArticle(ctx, article)
	if err != nil {
		t.Errorf("Failed to insert article: %s\n", err)
	}

	insertedID := res.InsertedID.(primitive.ObjectID)
	err = mongoClient.DeleteArticle(ctx, &insertedID)

	if err != nil {
		t.Errorf("Failed to delete article: %s\n", err)
	}

	insertedArticle, err := mongoClient.FindArticleByID(ctx, &insertedID)
	if err == nil {
		t.Errorf("Found article, but it should have been deleted: %v\n", insertedArticle)
	}
}

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	article := generateTestArticle()
	updated := generateTestArticle()

	insertRes, err := mongoClient.InsertArticle(ctx, article)
	if err != nil {
		t.Errorf("Failed to insert article: %s\n", err)
	}

	insertedID := insertRes.InsertedID.(primitive.ObjectID)

	updateRes, err := mongoClient.UpdateArticle(ctx, &insertedID, updated)
	if err != nil {
		t.Errorf("Failed to update article: %s\n", err)
	}

	if updateRes.ModifiedCount == 0 {
		t.Errorf("Update succesful, but no documents modified\n")
	} else if updateRes.ModifiedCount > 1 {
		t.Errorf("Multiple documents modified, but only one ID supplied\n")
	}
}

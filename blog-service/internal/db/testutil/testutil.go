package testutil

import (
	mongomodels "blog-service/internal/db/mongo/models"

	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	postgresStartupTimeout = 30 * time.Second
	postgresWaitOccurence  = 2
)

type MongoClientInterface interface {
	InsertArticle(context.Context, *mongomodels.ArticleDB) (*mongo.InsertOneResult, error)
}

type PostgresClientInterface interface {
	ExecuteQuery(string, context.Context) (pgconn.CommandTag, error)
}

func GenerateRandomString() string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	const stringLen = 10

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec // Secure rng not needed for generating usernames
	randomBytes := make([]byte, stringLen)
	for i := range randomBytes {
		randomBytes[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(randomBytes)
}

func CreateMongoContainer(ctx context.Context) (testcontainers.Container, string, error) {
	// Set env variables for mongodb
	os.Setenv("MONGO_DB", "blog_db")
	const rootUser = "rootuser"
	const rootPass = "rootpass"

	env := map[string]string{
		"MONGO_INITDB_ROOT_USERNAME": rootUser,
		"MONGO_INITDB_ROOT_PASSWORD": rootPass,
		"MONGO_INITDB_DATABASE":      "blog_db",
	}

	port := "27017/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo",
			ExposedPorts: []string{port},
			Env:          env,
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, "", fmt.Errorf("failed to start mongo container")
	}

	p, err := container.MappedPort(ctx, "27017")
	if err != nil {
		return container, "", fmt.Errorf("failed to get mongo container external port: %w", err)
	}
	os.Setenv("MONGO_URI", fmt.Sprintf("mongodb://%s:%s@localhost:%s", rootUser, rootPass, p.Port()))

	log.Printf("Mongo container up and running on port: %s\n", p.Port())

	return container, p.Port(), nil
}

func CreatePostgreSQLContainer(ctx context.Context) (testcontainers.Container, string, error) {
	dbName := "my_test_db"
	dbUser := "my_test_user"
	dbPassword := "my_secret_test_password"
	dbPort := "5432"

	os.Setenv("PG_DB", dbName)
	os.Setenv("PG_USER", dbUser)
	os.Setenv("PG_PASSWORD", dbPassword)
	os.Setenv("PG_PORT", dbPort)
	os.Setenv("PG_HOST", "localhost")

	postgresContainer, err := postgres.Run(ctx,
		"postgres:13-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(postgresWaitOccurence).
				WithStartupTimeout(postgresStartupTimeout),
		),
		postgres.WithSQLDriver("pgx"),
	)

	if err != nil {
		return postgresContainer, "", fmt.Errorf("failed to start mongo container")
	}

	p, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return postgresContainer, "", fmt.Errorf("failed to get postgres container external port: %w", err)
	}

	log.Printf("Postgres container up and running on port: %s\n", p.Port())

	return postgresContainer, p.Port(), nil
}

func GenerateTestArticle() *mongomodels.ArticleDB {
	const publisherIDRange = 100

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec // Secure rng not needed for generating usernames
	categories := [...]string{"sports", "tv", "travel", "gaming"}
	randomContent := GenerateRandomString()

	content := randomContent
	title := randomContent + "_title"
	category := categories[seededRand.Intn(len(categories))]
	publisherID := seededRand.Intn(publisherIDRange)

	return &mongomodels.ArticleDB{
		Title:       title,
		Content:     content,
		Category:    category,
		PublisherID: publisherID,
	}
}

func CreatePostgresTestArticle(ctx context.Context, mongoClient MongoClientInterface) (*mongomodels.ArticleDB, error) {
	article := GenerateTestArticle()

	res, err := mongoClient.InsertArticle(ctx, article)
	if err != nil {
		return nil, err
	}

	var ok bool
	article.ID, ok = res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to typecast InsertedID to primitive.ObjectID")
	}

	return article, nil
}

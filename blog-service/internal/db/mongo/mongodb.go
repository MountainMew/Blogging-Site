package mongo

import (
	"blog-service/internal/db/mongo/models"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	connectTimeout = 10 * time.Second
	insertTimeout  = 5 * time.Second
)

const (
	articleCollection = "articles"
)

type Client struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongoClient() (*Client, error) {
	mongoURI := os.Getenv("MONGO_URI")
	mongoDB := os.Getenv("MONGO_DB")

	if mongoURI == "" || mongoDB == "" {
		return nil, fmt.Errorf("lipsește una sau mai multe variabile de mediu necesare: MONGO_URI, MONGO_DB")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB!")

	return &Client{
		Client: client,
		DB:     client.Database(mongoDB),
	}, nil
}

func (c *Client) FindArticleByID(ctx context.Context, articleID *primitive.ObjectID) (*models.ArticleDB, error) {
	var article models.ArticleDB
	collection := c.DB.Collection(articleCollection)

	err := collection.FindOne(ctx,
		bson.D{{Key: "_id", Value: articleID}}).Decode(&article)
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (c *Client) InsertArticle(ctx context.Context, article *models.ArticleDB) (*mongo.InsertOneResult, error) {
	collection := c.DB.Collection(articleCollection)

	result, err := collection.InsertOne(ctx, article)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) UpdateArticle(ctx context.Context, articleID *primitive.ObjectID, modifiedArticle *models.ArticleDB) (*mongo.UpdateResult, error) {
	if modifiedArticle.Content == "" || modifiedArticle.Category == "" ||
		modifiedArticle.Title == "" || modifiedArticle.PublisherID == 0 {
		return nil, fmt.Errorf("missing one or more fields in the modified article")
	}

	if articleID == nil {
		return nil, fmt.Errorf("missing article id")
	}

	filter := bson.D{{Key: "_id", Value: articleID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "content", Value: modifiedArticle.Content},
		{Key: "category", Value: modifiedArticle.Category},
		{Key: "title", Value: modifiedArticle.Title},
		{Key: "publisher_id", Value: modifiedArticle.PublisherID},
		{Key: "publisher_name", Value: modifiedArticle.PublisherName},
	}}}

	collection := c.DB.Collection(articleCollection)

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) DeleteArticle(ctx context.Context, articleID *primitive.ObjectID) error {
	filter := bson.D{{Key: "_id", Value: articleID}}

	collection := c.DB.Collection(articleCollection)

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Close(ctx context.Context) error {
	if c.Client == nil {
		return nil
	}
	log.Println("Closing MongoDB connection.")
	return c.Client.Disconnect(ctx)
}

// FindArticlesByPublisherID retrieves all articles published by a specific user.
func (c *Client) FindArticlesByPublisherID(ctx context.Context, publisherID int) ([]*models.ArticleDB, error) {
	collection := c.DB.Collection(articleCollection)
	filter := bson.D{{Key: "publisherId", Value: publisherID}}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to execute find command: %w", err)
	}
	defer cursor.Close(ctx)

	var articles []*models.ArticleDB
	if err = cursor.All(ctx, &articles); err != nil {
		return nil, fmt.Errorf("failed to decode articles from cursor: %w", err)
	}

	return articles, nil
}

func (c *Client) GetArticlesByEngagement(page, limit int64) ([]models.ArticleDB, error) {
	collection := c.DB.Collection("articles")
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	skip := (page - 1) * limit
	if skip < 0 {
		skip = 0
	}

	log.Printf("Page: %d, Limit: %d, Skip: %d", page, limit, skip)

	sort := bson.D{{Key: "engagement", Value: -1}}
	findOptions := options.Find().SetSort(sort).SetLimit(limit).SetSkip(skip)

	cursor, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find articles: %w", err)
	}
	defer cursor.Close(ctx)

	var articles []models.ArticleDB
	if err = cursor.All(ctx, &articles); err != nil {
		return nil, fmt.Errorf("failed to decode articles: %w", err)
	}

	if err = cursor.Err(); err != nil {
		return nil, fmt.Errorf("failed during cursor iteration: %w", err)
	}

	return articles, nil
}

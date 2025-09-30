package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArticleDB struct {
	CreatedAt     time.Time          `bson:"createdAt"`
	Title         string             `bson:"title"`
	Content       string             `bson:"content"`
	Category      string             `bson:"category"`
	PublisherName string             `bson:"publisherName"`
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	PublisherID   int                `bson:"publisherId"`
}

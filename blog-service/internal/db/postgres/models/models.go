package models

import "time"

type Comment struct {
	CreatedAt     time.Time `json:"createdAt"`
	Content       string    `json:"content"`
	ArticleID     string    `json:"articleId"`
	CommenterName string    `json:"commenterName"`
	ID            int       `json:"id"`
	UserID        int       `json:"userId"`
}

type Like struct {
	ID        int
	CommentID int
	UserID    int
}

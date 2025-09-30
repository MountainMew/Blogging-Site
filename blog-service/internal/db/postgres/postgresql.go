package postgres

import (
	"blog-service/internal/db/postgres/models"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultTimeout = 5 * time.Second

type Client struct {
	ConnPool *pgxpool.Pool
	Host     string
	User     string
	Port     string
	Password string
	Dbname   string
}

func NewPostgresClient() (*Client, error) {
	client := Client{
		Host:     os.Getenv("PG_HOST"),
		Port:     os.Getenv("PG_PORT"),
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		Dbname:   os.Getenv("PG_DB"),
	}

	databaseURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?pool_max_conns=10",
		client.User, client.Password, client.Host, client.Port, client.Dbname)

	var err error

	// Create db connection
	client.ConnPool, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Create a timedout context for a qiuck databse connection check
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Validate db connection
	err = client.ConnPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	fmt.Println("Successfully connected to PostGres DB!")

	result, err := client.ExecuteQuery(ctx, SchemaComments)
	if err != nil {
		return nil, fmt.Errorf("cannot CREATE comments and likes Table: %w", err)
	}
	fmt.Println("CREATE comments and likes TABLE Result:", result.String())

	return &client, nil
}

/*
Function used for simple query execution with no parameters
Ex : CREATE TABLE, DROP TABLE
@params
schema - string with the table schema.
@returns
error - for checking the execution of the query.
pgconn.CommandTag - to check the query result.
*/
func (db *Client) ExecuteQuery(ctx context.Context, query string) (pgconn.CommandTag, error) {
	// Check db connection
	if db.ConnPool == nil {
		return pgconn.CommandTag{}, fmt.Errorf("unable to connect to database")
	}

	// Execute schema query
	commandTag, err := db.ConnPool.Exec(ctx, query)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("failed to execute query: %w", err)
	}

	return commandTag, nil
}

func (db *Client) GetCommentLikeCount(ctx context.Context, commentID int) (int, error) {
	const query = "SELECT COUNT(*) FROM likes WHERE CommentID = $1"

	if db.ConnPool == nil {
		return 0, fmt.Errorf("unable to connect to database")
	}

	var likeCount int
	err := db.ConnPool.QueryRow(ctx, query, commentID).Scan(&likeCount)
	if err != nil {
		return 0, err
	}

	return likeCount, nil
}

func (db *Client) HasLikedComment(ctx context.Context, commentID, userID int) (bool, error) {
	const query = `SELECT id FROM likes WHERE commentid = $1 AND userid = $2`

	if db.ConnPool == nil {
		return false, fmt.Errorf("unable to connect to database")
	}

	var id int

	err := db.ConnPool.QueryRow(ctx, query, commentID, userID).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (db *Client) GetComment(ctx context.Context, commentID int) (*models.Comment, error) {
	const query = `SELECT * FROM comments WHERE ID = $1`

	if db.ConnPool == nil {
		return nil, fmt.Errorf("unable to connect to database")
	}

	comment := &models.Comment{}

	err := db.ConnPool.QueryRow(ctx, query, commentID).Scan(
		&comment.ID,
		&comment.ArticleID,
		&comment.UserID,
		&comment.CommenterName,
		&comment.Content,
		&comment.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (db *Client) GetCommentsCount(ctx context.Context, articleID string) (int, error) {
	const query = "SELECT COUNT(*) FROM comments WHERE ArticleID = $1"

	if db.ConnPool == nil {
		return 0, fmt.Errorf("unable to connect to database")
	}

	var cnt int
	err := db.ConnPool.QueryRow(ctx, query, articleID).Scan(&cnt)
	if err != nil {
		return 0, err
	}

	return cnt, nil
}

func (db *Client) GetComments(ctx context.Context, articleID string, limit, page int) ([]models.Comment, error) {
	const query = `SELECT * FROM comments WHERE ArticleID = $1 ORDER BY id DESC LIMIT $2 OFFSET $3`

	if db.ConnPool == nil {
		return nil, fmt.Errorf("unable to connect to database")
	}

	offset := limit * (page - 1)
	rows, err := db.ConnPool.Query(ctx, query, articleID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]models.Comment, 0, limit)
	comment := models.Comment{}

	for rows.Next() {
		err := rows.Scan(
			&comment.ID,
			&comment.ArticleID,
			&comment.UserID,
			&comment.CommenterName,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve row: %w", err)
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (db *Client) CreateComment(ctx context.Context, comment *models.Comment) (int, error) {
	var commentID int

	const query = `INSERT INTO comments (ArticleID, UserID, CommenterName, Content, CreatedAt) 
	          VALUES($1, $2, $3, $4, $5) RETURNING ID`

	// Check db connection
	if db.ConnPool == nil {
		return 0, fmt.Errorf("unable to connect to database")
	}

	// Execute query
	err := db.ConnPool.QueryRow(ctx, query, comment.ArticleID, comment.UserID, comment.CommenterName, comment.Content, comment.CreatedAt).Scan(&commentID)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return commentID, nil
}

func (db *Client) DeleteComment(ctx context.Context, commentID int) (pgconn.CommandTag, error) {
	const query = `DELETE FROM comments WHERE id = $1`

	// Check db connection
	if db.ConnPool == nil {
		return pgconn.CommandTag{}, fmt.Errorf("unable to connect to database")
	}

	// Execute query
	commandTag, err := db.ConnPool.Exec(ctx, query, commentID)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("failed to execute query: %w", err)
	}

	return commandTag, nil
}

func (db *Client) FindLike(ctx context.Context, commentID, userID int) (*models.Like, error) {
	const query = `SELECT id FROM likes WHERE CommentID = $1 AND UserID = $2`

	if db.ConnPool == nil {
		return nil, fmt.Errorf("unable to connect to database")
	}

	like := &models.Like{
		CommentID: commentID,
		UserID:    userID,
	}
	err := db.ConnPool.QueryRow(ctx, query, commentID, userID).Scan(
		&like.ID,
	)

	if err != nil {
		return nil, err
	}

	return like, nil
}

func (db *Client) AddLike(ctx context.Context, like models.Like) (int, error) {
	const query = `INSERT INTO likes (CommentID, UserID) 
	          VALUES($1, $2) RETURNING id`

	// Check db connection
	if db.ConnPool == nil {
		return 0, fmt.Errorf("unable to connect to database")
	}

	var id int

	// Execute query
	err := db.ConnPool.QueryRow(ctx, query, like.CommentID, like.UserID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	return id, nil
}

func (db *Client) RemoveLike(ctx context.Context, commentID, userID int) (pgconn.CommandTag, error) {
	const query = `DELETE FROM likes WHERE CommentID = $1 AND UserID = $2`

	// Check db connection
	if db.ConnPool == nil {
		return pgconn.CommandTag{}, fmt.Errorf("unable to connect to database")
	}

	// Execute query
	commandTag, err := db.ConnPool.Exec(ctx, query, commentID, userID)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("failed to execute query: %w", err)
	}

	return commandTag, nil
}

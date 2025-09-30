package db

import (
	"auth-service/internal/auth/models"
	"auth-service/internal/crypto"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const defaultTimeout = 5 * time.Second

//nolint:govet //it is already in the correct order (biggest size to smallest)
type Database struct {
	Host     string
	User     string
	Port     int
	Password string
	Dbname   string
	ConnPool *pgxpool.Pool
}

/*
Etablish database configuration.
Creates a new Databse instance.
@returns
*Databse - pointer to a database structure.
*/
func Config() (*Database, error) {
	// Convert .env PORT to int
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("error converting DB_PORT: %w", err)
	}

	db := Database{
		Host:     os.Getenv("DB_HOST"),
		Port:     dbPort,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
	}

	return &db, nil
}

/*
Creates a dabase connection.
Sets db structure filed Coon to a pointer of db connection type.
@returns
error - for checing the db coonection.
*/
func (db *Database) OpenDbConnection() error {
	// Create database URL
	databaseURL := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?pool_max_conns=10",
		db.User, db.Password, db.Host, db.Port, db.Dbname)

	var err error

	// Create db connection
	db.ConnPool, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Create a timedout context for a qiuck databse connection check
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Validate db connection
	err = db.ConnPool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	fmt.Println("Successfully connected to PostGres DB!")

	return nil
}

func (db *Database) UserExists(username string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)"
	var exists bool
	err := db.ConnPool.QueryRow(context.Background(), query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to execute user exists query: %w", err)
	}

	return exists, nil
}

func (db *Database) Init() error {
	result, err := db.ExecuteQuery(Schema)
	if err != nil {
		return fmt.Errorf("cannot CREATE users Table: %w", err)
	}
	fmt.Println("CREATE users TABLE Result:", result.String())

	defaultUserUsername := os.Getenv("DEFAULT_USER_USERNAME")

	// Check if the default user already exists
	exists, err := db.UserExists(defaultUserUsername)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %w", err)
	}
	if exists {
		fmt.Printf("Default user '%s' already exists. Skipping creation.\n", defaultUserUsername)
		return nil
	}

	// Hash Default User Password
	params := crypto.GetDefaultParams()
	hashedPassword, err := crypto.HashPassword(os.Getenv("DEFAULT_USER_PASS"), &params)
	if err != nil {
		return fmt.Errorf("error in hashing password %w", err)
	}

	defaultUser := models.User{
		Username: os.Getenv("DEFAULT_USER_USERNAME"),
		Password: hashedPassword,
		Email:    os.Getenv("DEFAULT_USER_EMAIL"),
		Role:     models.ADMIN,
	}

	result, err = db.CreateUser(&defaultUser)
	if err != nil {
		return fmt.Errorf("cannot CREATE users Table: %w", err)
	}
	fmt.Println("CREATE default user Result:", result.String())

	return nil
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
func (db *Database) ExecuteQuery(query string) (pgconn.CommandTag, error) {
	// Check db connection
	if db.ConnPool == nil {
		return pgconn.CommandTag{}, fmt.Errorf("unable to connect to database")
	}

	// Create a backgrond context ( no timeout or deadline )
	ctx := context.Background()

	// Execute schema query
	commandTag, err := db.ConnPool.Exec(ctx, query)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("failed to execute query: %w", err)
	}

	return commandTag, nil
}

/*
Function used to insert an User
@params
query - string INSERT query
user - user structure with the new db entry
@returns
error - for checking the execution of the query.
pgconn.CommandTag - to check the query result.
*/
func (db *Database) CreateUser(user *models.User) (pgconn.CommandTag, error) {
	// InsertQuery
	const query = `INSERT INTO users (Username, Password, Email, Role) 
	          VALUES($1, $2, $3, $4)`

	// Check db connection
	if db.ConnPool == nil {
		return pgconn.CommandTag{}, fmt.Errorf("unable to connect to database")
	}

	// Create a backgrond context ( no timeout or deadline )
	ctx := context.Background()

	// Execute query
	commandTag, err := db.ConnPool.Exec(ctx, query, user.Username, user.Password, user.Email, user.Role.RoleString())
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("failed to execute query: %w", err)
	}

	return commandTag, nil
}

/*
Function used to update an User
@params
query - string UPDATE query
user - user structure for the update table entry
@returns
error - for checking the execution of the query.
pgconn.CommandTag - to check the query result.
*/
func (db *Database) UpdateUser(userToUpdate *models.User) (pgconn.CommandTag, error) {
	// Update query
	const query = `UPDATE users 
              SET Username = $2, Password = $3, Email = $4, Role = $5
              WHERE ID = $1`

	// Check db connection
	if db.ConnPool == nil {
		return pgconn.CommandTag{}, fmt.Errorf("unable to connect to database")
	}

	// Create a backgrond context ( no timeout or deadline )
	ctx := context.Background()

	// Execute query
	commandTag, err := db.ConnPool.Exec(ctx, query,
		userToUpdate.ID,
		userToUpdate.Username,
		userToUpdate.Password,
		userToUpdate.Email,
		userToUpdate.Role.RoleString(),
	)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("failed to execute query: %w", err)
	}

	return commandTag, nil
}

/*
Function used to delete an User
@params
query - string DELETE query
user - user structure for the deleted table entry
@returns
error - for checking the execution of the query.
pgconn.CommandTag - to check the query result.
*/
func (db *Database) DeleteUser(userToDelete *models.User) (pgconn.CommandTag, error) {
	// Delete query
	const query = `DELETE FROM users WHERE id = $1`

	// Check db connection
	if db.ConnPool == nil {
		return pgconn.CommandTag{}, fmt.Errorf("unable to connect to database")
	}

	// Create a backgrond context ( no timeout or deadline )
	ctx := context.Background()

	// Execute query
	commandTag, err := db.ConnPool.Exec(ctx, query, userToDelete.ID)
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("failed to execute query: %w", err)
	}

	return commandTag, nil
}

/*
Function used to SELECT an user by ID
@params
query - string DELETE query
id - user Id
@returns
user - struc with the user from the table
error - for checking the execution of the query.
*/
func (db *Database) SelectUserByID(id int) (*models.User, error) {
	// Select User by ID
	const query = `SELECT id, username, password, email, role FROM users WHERE id = $1`

	// Check db connection
	if db.ConnPool == nil {
		return nil, fmt.Errorf("unable to connect to database")
	}

	user := models.User{}
	// Temporary role variable
	roleStr := ""

	// Create a backgrond context
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Execute SELECT query
	err := db.ConnPool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&roleStr,
	)

	// Set user Role
	user.Role = models.StringToRole(roleStr)

	// Error check
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

/*
Function used to SELECT an user by Username
@params
query - string DELETE query
username - user username
@returns
user - struc with the user from the table
error - for checking the execution of the query.
*/
func (db *Database) SelectUserByUsername(username string) (*models.User, error) {
	// Select ID query
	const query = `SELECT id, username, password, email, role FROM users WHERE username = $1`

	// Check db connection
	if db.ConnPool == nil {
		return nil, fmt.Errorf("unable to connect to database")
	}

	user := models.User{}
	// Temporary role variable
	roleStr := ""

	// Create a backgrond context
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Execute SELECT query
	err := db.ConnPool.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&roleStr,
	)

	// Set user Role
	user.Role = models.StringToRole(roleStr)

	// Error check
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

/*
Function used to SELECT an user by Username and Password
@params
query - string DELETE query
username - user username
pass - user hashed password
@returns
user - struc with the user from the table
error - for checking the execution of the query.
*/
func (db *Database) SelectUserByUsernameAndPass(username, pass string) (*models.User, error) {
	// Select User Name and Password
	const query = `SELECT id, username, password, email, role FROM users WHERE username = $1 AND password =$2`

	// Check db connection
	if db.ConnPool == nil {
		return nil, fmt.Errorf("unable to connect to database")
	}

	user := models.User{}
	// Temporary role variable
	roleStr := ""

	// Create a backgrond context
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Execute SELECT query
	err := db.ConnPool.QueryRow(ctx, query, username, pass).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&roleStr,
	)

	// Set user Role
	user.Role = models.StringToRole(roleStr)

	// Error check
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

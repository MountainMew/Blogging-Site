package main_test

import (
	"auth-service/internal/crypto"
	"auth-service/internal/db"
	pb "auth-service/internal/protobuf"
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const defaultTimeOut = 10 * time.Second
const TOKEN_PARTS = 3

var (
	grpcAddr string
)

func TestMain(t *testing.T) {
	ctx := context.Background()

	// define db variables
	dbName := "users"
	dbUser := "user"
	dbPassword := "password"
	dbHost := "db"

	// Hardcoded env variables
	// TODO: Change in the future
	JWT_SECRET := "test_secret"
	DEFAULT_USER_USERNAME := "test_admin"
	DEFAULT_USER_PASS := "test_pass"
	DEFAULT_USER_EMAIL := "test@email.com"

	// =================================================================
	// Creating network to share for all container services
	// =================================================================
	net, err := network.New(ctx)
	require.NoError(t, err, "Failed to create Docker network")

	// defer network termination
	defer func() {
		err := net.Remove(ctx)
		require.NoError(t, err, fmt.Sprintf("failed to remove network: %s", err))
	}()

	// =================================================================
	// Create postgres container
	// =================================================================
	postgresContainer, err := postgres.Run(ctx, "postgres",
		network.WithNetwork([]string{net.Name}, net),
		postgres.WithDatabase(dbHost),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(10*time.Second),
		),
		postgres.WithSQLDriver("pgx"),
	)

	// defer container termination
	testcontainers.CleanupContainer(t, postgresContainer)
	require.NoError(t, err)
	defer func() {
		err := postgresContainer.Terminate(ctx)
		require.NoError(t, err, "Failed to terminate postgres container")
	}()

	// =================================================================
	// Create database connection
	// =================================================================

	// Get postgres service container address from within network
	pgIP, err := postgresContainer.ContainerIP(ctx)
	require.NoError(t, err, fmt.Sprintf("failed to get pgcontainer IP: %s", err))

	// Get postgres service container host
	pgHost, err := postgresContainer.Host(ctx)
	require.NoError(t, err, fmt.Sprintf("failed to get pgcontainer host: %s", err))

	// Get postgres service container port
	pgPort, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err, fmt.Sprintf("failed to get container mapped port: %s", err))

	// Set environment variables for the test process
	t.Setenv("DB_HOST", pgHost)
	t.Setenv("DB_PORT", pgPort.Port())
	t.Setenv("DB_USER", dbUser)
	t.Setenv("DB_PASSWORD", dbPassword)
	t.Setenv("DB_NAME", dbName)

	// Config database structure
	database, err := db.Config()
	require.NoError(t, err, fmt.Sprintf("failed to config database: %s", err))
	// Establish database connection
	err = database.OpenDbConnection()
	require.NoError(t, err, fmt.Sprintf("failed to open test db connection: %s", err))

	defer database.ConnPool.Close()

	// =================================================================
	// Set up auth service container
	// =================================================================
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context: "..",
		},
		ExposedPorts: []string{"9001/tcp"},
		WaitingFor:   wait.ForLog("Listening on port 9001..."),
		Networks:     []string{net.Name},
		Env: map[string]string{
			// Inject pg container address as env variable
			"DB_HOST":               pgIP,
			"DB_PORT":               "5432",
			"DB_USER":               dbUser,
			"DB_PASSWORD":           dbPassword,
			"DB_NAME":               dbName,
			"JWT_SECRET":            JWT_SECRET,
			"DEFAULT_USER_USERNAME": DEFAULT_USER_USERNAME,
			"DEFAULT_USER_PASS":     DEFAULT_USER_PASS,
			"DEFAULT_USER_EMAIL":    DEFAULT_USER_EMAIL,
		},
	}

	// Create auth service container
	authContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	// defer container termination
	testcontainers.CleanupContainer(t, authContainer)
	require.NoError(t, err)
	defer func() {
		err := authContainer.Terminate(ctx)
		require.NoError(t, err, "Failed to terminate auth container")
	}()

	// =================================================================
	// Set up gRPC connection
	// =================================================================

	// Get auth service host
	authHost, err := authContainer.Host(ctx)
	require.NoError(t, err, "Failed to get auth container host")

	// Get auth service public mapped port
	authPort, err := authContainer.MappedPort(ctx, "9001")
	require.NoError(t, err, "Failed to get auth container port")

	// Set auth service address to enable clients to send requests
	grpcAddr = fmt.Sprintf("%s:%s", authHost, authPort.Port())
	log.Printf("gRPC service available at %s", grpcAddr)

	testCtx, cancel := context.WithTimeout(ctx, defaultTimeOut)
	defer cancel()

	// Establish gRPC connection
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "Failed to connect to grpc.")

	defer conn.Close()

	// Create new gRPC client
	client := pb.NewAuthServiceClient(conn)

	// =================================================================
	// Run subtests
	// =================================================================
	t.Run("1. Testing Login As Admin", func(t *testing.T) {

		// Set request payload with valid credentials
		LoginUserReq := pb.LoginRequest{
			Username: "test_admin",
			Password: "test_pass",
		}

		// Validate response
		LoginUserResponse, err := client.Login(testCtx, &LoginUserReq)
		require.NoError(t, err, fmt.Sprintf("Could not log in user: %v", err))

		// Check if the token is empty
		require.NotEqual(t, LoginUserResponse.Token, "", "Expected a token, but got an empty string")

		// Check for JWT structure (header.payload.signature)
		parts := strings.Split(LoginUserResponse.Token, ".")
		require.Equal(t, len(parts), TOKEN_PARTS, fmt.Sprintf("Expected a valid JWT structure, but got %d parts", len(parts)))
	})

	t.Run("2. Testing login with invalid data", func(t *testing.T) {
		// Set request payload with invalid credentials
		LoginUserReq := pb.LoginRequest{
			Username: "",
			Password: "",
		}

		// Validate response
		_, err = client.Login(testCtx, &LoginUserReq)
		require.Error(t, err, "Successful login with invalid data")

	})

	t.Run("3. Testing successful user creation", func(t *testing.T) {
		// Set request payload
		createUserReq := pb.CreateUserRequest{
			Email:    "somerandommail@example.com",
			Username: "johndoe1",
			Password: "johnspassword",
		}

		// Validate response
		_, err = client.CreateUser(testCtx, &createUserReq)
		require.NoError(t, err, fmt.Sprintf("Could not register user: %v", err))

		// Get inserted user from the db
		resultUser, err := database.SelectUserByUsername(createUserReq.Username)
		require.NoError(t, err, fmt.Sprintf("Failed to select user: %v", err))

		// Compare username
		require.Equal(t, createUserReq.Username, resultUser.Username, "The req username and result username should be the same.")

		// Compare email
		require.Equal(t, resultUser.Email, createUserReq.Email, "The req email and result email should be the same.")

		// Compare hashed passwords
		ok, err := crypto.VerifyPassword(createUserReq.Password, resultUser.Password)
		require.NoError(t, err, fmt.Sprintf("Internal crypto package error %v", err))
		require.Equal(t, ok, true, "Passwords do not match")

	})

	t.Run("4. Testing create duplicated user", func(t *testing.T) {
		// Set request payload
		createUserReq := pb.CreateUserRequest{
			Email:    "somerandommail@example.com",
			Username: "johndoe1",
			Password: "johnspassword",
		}

		// Validate response
		_, err = client.CreateUser(testCtx, &createUserReq)
		require.Error(t, err, fmt.Sprintf("Duplicated user created %v", err))
	})

	t.Run("5. Testing Login as User", func(t *testing.T) {
		//Set request payload with valid credentials
		LoginUserReq := pb.LoginRequest{
			Username: "johndoe1",
			Password: "johnspassword",
		}

		// Validate response
		LoginUserResponse, err := client.Login(testCtx, &LoginUserReq)
		require.NoError(t, err, fmt.Sprintf("Could not log in user: %v", err))

		// Check if the token is empty
		require.NotEqual(t, LoginUserResponse.Token, "", "Expected a token, but got an empty string")

		// Check for JWT structure (header.payload.signature)
		parts := strings.Split(LoginUserResponse.Token, ".")
		require.Equal(t, len(parts), TOKEN_PARTS, fmt.Sprintf("Expected a valid JWT structure, but got %d parts", len(parts)))
	})

	t.Run("6._Testing_unauthorized_valid_data_update", func(t *testing.T) {
		// Login with user credentials
		LoginUserReq := pb.LoginRequest{
			Username: "johndoe1",
			Password: "johnspassword",
		}

		// Validate Login response
		LoginUserResponse, err := client.Login(testCtx, &LoginUserReq)
		require.NoError(t, err, fmt.Sprintf("Could not log in user: %v", err))

		// Get JWT token
		token := LoginUserResponse.GetToken()

		// Create metadata with JWT token
		md := metadata.New(map[string]string{
			"authorization": "Bearer " + token,
		})

		// Create new context with metadata
		authedCtx := metadata.NewOutgoingContext(testCtx, md)

		// Set request payload
		UpdateReq := pb.UpdateUserRequest{
			Id:       2,
			Email:    "johnpork@example.com",
			Username: "johnpork",
			Password: "johnporkpassword",
		}

		// Validate response
		_, err = client.UpdateUser(authedCtx, &UpdateReq)
		require.Error(t, err, fmt.Sprintf("Successful user update with invalid credentials: %v", err))
	})

	t.Run("7._Testing_authorized_invalid_data_update", func(t *testing.T) {
		// Login with user credentials
		LoginUserReq := pb.LoginRequest{
			Username: "test_admin",
			Password: "test_pass",
		}

		// Validate Login response
		LoginUserResponse, err := client.Login(testCtx, &LoginUserReq)
		require.NoError(t, err, fmt.Sprintf("Could not log in user: %v", err))

		// Get JWT token
		token := LoginUserResponse.GetToken()

		// Create metadata with JWT token
		md := metadata.New(map[string]string{
			"authorization": "Bearer " + token,
		})

		// Create new context with metadata
		authedCtx := metadata.NewOutgoingContext(testCtx, md)

		// Set request payload
		UpdateReq := pb.UpdateUserRequest{
			Id:       99999999,
			Email:    "",
			Username: "",
			Password: "",
		}

		// Validate response
		_, err = client.UpdateUser(authedCtx, &UpdateReq)
		require.Error(t, err, fmt.Sprintf("Successful user update with invalid data: %v", err))
	})

	t.Run("8._Testing_authorized_valid_data_update", func(t *testing.T) {
		// Login with user credentials
		LoginUserReq := pb.LoginRequest{
			Username: "test_admin",
			Password: "test_pass",
		}

		// Validate Login response
		LoginUserResponse, err := client.Login(testCtx, &LoginUserReq)
		require.NoError(t, err, fmt.Sprintf("Could not log in user: %v", err))

		// Get JWT token
		token := LoginUserResponse.GetToken()

		// Create metadata with JWT token
		md := metadata.New(map[string]string{
			"authorization": "Bearer " + token,
		})

		// Create new context with metadata
		authedCtx := metadata.NewOutgoingContext(testCtx, md)

		// Set request payload
		UpdateReq := pb.UpdateUserRequest{
			Id:       2,
			Email:    "johnpork@example.com",
			Username: "johnpork",
			Password: "johnporkpassword",
		}

		// Validate response
		_, err = client.UpdateUser(authedCtx, &UpdateReq)
		require.NoError(t, err, fmt.Sprintf("Could not update the user: %v", err))

		// Get inserted user from the db
		resultUser, err := database.SelectUserByUsername(UpdateReq.Username)
		require.NoError(t, err, fmt.Sprintf("Failed to select user: %v", err))

		// Compare username
		require.Equal(t, UpdateReq.Username, resultUser.Username, "The req username and result username should be the same.")

		// Compare email
		require.Equal(t, resultUser.Email, UpdateReq.Email, "The req email and result email should be the same.")

		// Compare hashed passwords
		ok, err := crypto.VerifyPassword(UpdateReq.Password, resultUser.Password)
		require.NoError(t, err, fmt.Sprintf("Internal crypto package error %v", err))
		require.Equal(t, ok, true, "Passwords do not match")
	})

	t.Run("9._Testing_unauthorized_valid_data_delete", func(t *testing.T) {
		// Login with user credentials
		LoginUserReq := pb.LoginRequest{
			Username: "johnpork",
			Password: "johnporkpassword",
		}

		// Validate Login response
		LoginUserResponse, err := client.Login(testCtx, &LoginUserReq)
		require.NoError(t, err, fmt.Sprintf("Could not log in user: %v", err))

		// Get JWT token
		token := LoginUserResponse.GetToken()

		// Create metadata with JWT token
		md := metadata.New(map[string]string{
			"authorization": "Bearer " + token,
		})

		// Create new context with metadata
		authedCtx := metadata.NewOutgoingContext(ctx, md)

		// Set request payload
		deleteUserReq := pb.DeleteUserRequest{
			Id: 1,
		}

		// Validate response
		_, err = client.DeleteUser(authedCtx, &deleteUserReq)
		require.Error(t, err, fmt.Sprintf("Successful user deletion with invalid credentials: %v", err))

	})

	t.Run("10._Testing_authorized_invalid_data_delete", func(t *testing.T) {
		// Login with admin credentials
		LoginUserReq := pb.LoginRequest{
			Username: "test_admin",
			Password: "test_pass",
		}

		// Validate Login response
		LoginUserResponse, err := client.Login(ctx, &LoginUserReq)
		require.NoError(t, err, fmt.Sprintf("Could not log in user: %v", err))

		// Get JWT token
		token := LoginUserResponse.GetToken()

		// Create metadata with JWT token
		md := metadata.New(map[string]string{
			"authorization": "Bearer " + token,
		})

		// Create new context with metadata
		authedCtx := metadata.NewOutgoingContext(ctx, md)

		// Set request payload
		deleteUserReq := pb.DeleteUserRequest{
			Id: 999999999,
		}

		// Validate response
		_, err = client.DeleteUser(authedCtx, &deleteUserReq)
		require.Error(t, err, fmt.Sprintf("Successful user deletion with invalid data: %v", err))

	})

	t.Run("11._Testing_authorized_valid_data_delete", func(t *testing.T) {
		// Login with admin credentials
		LoginUserReq := pb.LoginRequest{
			Username: "test_admin",
			Password: "test_pass",
		}

		// Validate Login response
		LoginUserResponse, err := client.Login(ctx, &LoginUserReq)
		require.NoError(t, err, fmt.Sprintf("Could not log in user: %v", err))

		// Get JWT token
		token := LoginUserResponse.GetToken()

		// Create metadata with JWT token
		md := metadata.New(map[string]string{
			"authorization": "Bearer " + token,
		})

		// Create new context with metadata
		authedCtx := metadata.NewOutgoingContext(ctx, md)

		// Set request payload
		deleteUserReq := pb.DeleteUserRequest{
			Id: 2,
		}

		// Validate response
		_, err = client.DeleteUser(authedCtx, &deleteUserReq)
		require.NoError(t, err, fmt.Sprintf("Could not delete user: %v", err))

		// Get inserted user from the db
		_, err = database.SelectUserByID(int(deleteUserReq.Id))
		require.Error(t, err, "User still exists after delete operation")

	})

}

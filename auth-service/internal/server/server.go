package server

import (
	"auth-service/internal/auth/jwt"
	"auth-service/internal/auth/models"
	"auth-service/internal/crypto"
	"auth-service/internal/db"
	pb "auth-service/internal/protobuf"
	"context"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const PORT string = "9001"

const expirationTime = 5 * time.Minute

type userClaimsKey struct{}

type Server struct {
	pb.UnimplementedAuthServiceServer
	db *db.Database
}

func NewGRPCServer(database *db.Database) *Server {
	return &Server{
		db: database,
	}
}

// Check valid email.
func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

// Midleware to intercept API calls and validate them before reaching theire handlers.
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Ignore Login and CreatUser requsts ( they don't have a token )
	if strings.HasSuffix(info.FullMethod, "Login") || strings.HasSuffix(info.FullMethod, "CreateUser") ||
		strings.HasSuffix(info.FullMethod, "VerifyToken") {
		return handler(ctx, req)
	}

	// Extract request metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	// Extract JWT token from metadata
	tokenStrings := md.Get("authorization")
	if len(tokenStrings) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
	}

	// Extract JWT token string
	tokenStr := strings.TrimPrefix(tokenStrings[0], "Bearer ")

	// Token Validation
	_, claims, err := jwt.ValidateJWT(tokenStr)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token: "+err.Error())
	}

	// Extract JWT Claim and add it to the request context
	newCtx := context.WithValue(ctx, userClaimsKey{}, claims)

	return handler(newCtx, req)
}

// Login handler.
func (s *Server) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	//	Validate requst
	if request.Password == "" || request.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "Invalid request, missing some request parameter(s).")
	}

	// Extract user from db with request credentials
	user, err := s.db.SelectUserByUsername(request.Username)
	if err != nil {
		return nil, status.Error(codes.NotFound, "invalid username or password")
	}

	// Hash Password check
	ok, err := crypto.VerifyPassword(request.Password, user.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "invalid username or password")
	}
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "invalid username or password")
	}

	// JWT Token creation
	token, err := jwt.GenerateJWT(*user, expirationTime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &pb.LoginResponse{Token: token}, nil
}

// Auxiliar function for different request handlers.
// Checks user role to ensure each user performs operation that are authorized.
func checkRole(ctx context.Context, requiredRole models.Role) error {
	// Extract JWT claim from request context
	rawClaims := ctx.Value(userClaimsKey{})
	claims, ok := rawClaims.(*jwt.CustomClaims)
	if !ok || claims == nil {
		return status.Error(codes.Internal, "user claims not found in context")
	}

	// Check user role
	if claims.Role != requiredRole {
		return status.Error(codes.PermissionDenied, "insufficient permissions")
	}

	return nil
}

func (s *Server) VerifyToken(ctx context.Context, request *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	if request.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request, missing token")
	}

	_, claims, err := jwt.ValidateJWT(request.Token)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "invalid token supplied")
	}

	return &pb.VerifyTokenResponse{
		Username: claims.Username,
		Id:       claims.ID,
		Role:     claims.Role.RoleString(),
	}, nil
}

// CreateUser handler.
func (s *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*emptypb.Empty, error) {
	// Validate credentials
	if request.Email == "" || request.Password == "" || request.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "Invalid request, missing some request parameter(s).")
	}

	emailFormatted := strings.ToLower(request.Email)

	if !isEmailValid(emailFormatted) {
		return nil, status.Error(codes.InvalidArgument, "invalid email")
	}

	// Hash password
	params := crypto.GetDefaultParams()
	hashedPassword, err := crypto.HashPassword(request.Password, &params)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error in hashing password.")
	}

	// Inject request body into a user structure
	user := &models.User{
		Username: request.Username,
		Password: hashedPassword,
		Email:    emailFormatted,
		Role:     models.USER,
	}

	// Query db for a new user entry
	_, err = s.db.CreateUser(user)
	if err != nil {
		fmt.Printf("Failed creating user with err: %s", err)
		return nil, status.Error(codes.AlreadyExists, "email or username already taken "+err.Error())
	}
	fmt.Printf("Created user %s", user.Username)

	return nil, nil
}

// UpdateUser handler.
func (s *Server) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*emptypb.Empty, error) {
	changed := false

	// Validate incoming request
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "cannot process the request")
	}

	// Check for user permission
	if err := checkRole(ctx, models.ADMIN); err != nil {
		return nil, err
	}

	// Validate user Id
	if request.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid ID")
	}

	// Query db for user entry to be update
	user, err := s.db.SelectUserByID(int(request.Id))
	if err != nil {
		if err.Error() == "user not found" {
			return nil, status.Error(codes.NotFound, "invalid id supplied")
		}
		return nil, status.Error(codes.Internal, "internal server error"+err.Error())
	}

	emailFormatted := strings.ToLower(request.Email)
	if emailFormatted != "" && emailFormatted != user.Email {
		if !isEmailValid(emailFormatted) {
			return nil, status.Error(codes.InvalidArgument, "invalid email")
		}
		user.Email = emailFormatted
		changed = true
	}
	// Check for password change
	if request.Password != "" {
		params := crypto.GetDefaultParams()
		user.Password, err = crypto.HashPassword(request.Password, &params)
		if err != nil {
			return nil, status.Error(codes.Internal, "internal server error")
		}
		changed = true
	}
	if request.Username != "" && request.Username != user.Username {
		user.Username = request.Username
		changed = true
	}
	if request.Role != "" && request.Role != user.Role.RoleString() {
		user.Role = models.StringToRole(request.Role)
		changed = true
	}

	// Verify update
	if !changed {
		return nil, status.Error(codes.InvalidArgument, "no changes to be made")
	}

	// Query db for user entry to update
	_, err = s.db.UpdateUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			return nil, status.Error(codes.AlreadyExists, "username or email already taken")
		}
		return nil, status.Error(codes.InvalidArgument, "no changes to be made")
	}

	return nil, nil
}

// DeleteUser handler.
func (s *Server) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	// check for corect request
	if request == nil {
		return nil, status.Errorf(codes.InvalidArgument, "cannot process the request")
	}

	if err := checkRole(ctx, models.ADMIN); err != nil {
		return nil, err
	}

	// check for a valid ID
	if request.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "a valid User ID must be provided")
	}

	// convert from int32 to int
	intID := int(request.Id)

	// Select user to delte
	userToDelete, err := s.db.SelectUserByID(intID)
	if err != nil {
		// log DB level error
		log.Printf("ERROR: could not retrieve user %d from database: %v", intID, err)
		return nil, status.Errorf(codes.NotFound, "user with id %d not found", intID)
	}

	_, err = s.db.DeleteUser(userToDelete)
	if err != nil {
		// log DB level error
		log.Printf("ERROR: could not retrieve user %d from database: %v", intID, err)
		return nil, status.Errorf(codes.NotFound, "user with id %d not found", intID)
	}

	return nil, nil
}

func (s *Server) ListenAndServe(ctx context.Context) {
	lc := net.ListenConfig{}
	lis, err := lc.Listen(ctx, "tcp", fmt.Sprintf(":%s", PORT))
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
	fmt.Printf("Listening on port %s...\n", PORT)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(AuthInterceptor),
	)

	reflection.Register(grpcServer) // Register the reflection service for easier debugging

	pb.RegisterAuthServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

package authservice_test

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"auth-service/internal/auth/jwt"
	"auth-service/internal/auth/models"
	pb "auth-service/internal/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type loginType string

type dummyUser struct {
	email    string
	username string
	password string
	id       int32
}

//nolint:gochecknoglobals // Required for testing
var (
	createDummy dummyUser
	deleteDummy dummyUser
	updateDummy dummyUser
)

const (
	USER  loginType = "user"
	ADMIN loginType = "admin"
)

const usernameLen = 8

// grpc client for testing.
var client pb.AuthServiceClient //nolint:gochecknoglobals // Required for testing

const (
	// admin user credentials for testing.
	adminUsername = "admin"
	adminPassword = "admin"

	// regular user credentials for testing.
	userUsername       = "testuser"
	userPassword       = "testuser"
	userEmail          = "test@test.com"
	userID       int32 = 2
)

func initDummyUser(email, username, password string) dummyUser {
	return dummyUser{
		email:    email,
		username: username,
		password: password,
	}
}

func generateRandomUser() (email, username, password string) {
	const charset = "abcdefghijklmnopqrstuvwxyz"

	seededRand := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec // Secure rng not needed for generating usernames

	usernameBytes := make([]byte, usernameLen)
	for i := range usernameBytes {
		usernameBytes[i] = charset[seededRand.Intn(len(charset))]
	}

	username = string(usernameBytes)
	email = username + "@gmail.com"
	password = username

	return email, username, password
}

func generateDummyUsers() {
	createDummy = initDummyUser(generateRandomUser())
	updateDummy = initDummyUser(generateRandomUser())
	deleteDummy = initDummyUser(generateRandomUser())
}

func createDummyUser(user *dummyUser) error {
	createReq := &pb.CreateUserRequest{
		Email:    user.email,
		Username: user.username,
		Password: user.password,
	}

	_, err := client.CreateUser(context.Background(), createReq)
	if err != nil {
		log.Fatalf("Failed to create dummy user: %s\n", err)
	}

	loginResp, err := client.Login(context.Background(), &pb.LoginRequest{
		Username: user.username,
		Password: user.password,
	})
	if err != nil {
		log.Fatalf("Failed to fetch dummy user ID: %s\n", err)
	}

	_, userClaims, err := jwt.ValidateJWT(loginResp.Token)
	if err != nil {
		log.Fatalf("Failed to fetch dummy user ID: %s\n", err)
	}

	user.id = userClaims.ID

	return err
}

func loginUser(userType loginType) *pb.LoginResponse {
	loginReq := &pb.LoginRequest{}

	switch userType {
	case ADMIN:
		loginReq.Username = adminUsername
		loginReq.Password = adminPassword
	case USER:
		loginReq.Username = userUsername
		loginReq.Password = userPassword
	}

	loginResp, err := client.Login(context.Background(), loginReq)
	if err != nil {
		log.Fatalf("Failed to login as %s: %s\n", userType, err)
	}

	return loginResp
}

func loggedInContext(userType loginType) context.Context {
	loginResp := loginUser(userType)

	md := metadata.Pairs("authorization", loginResp.Token)
	authCtx := metadata.NewOutgoingContext(context.Background(), md)

	return authCtx
}

func malformedTokenContext() context.Context {
	loginResp := loginUser(USER)

	md := metadata.Pairs("authorization", loginResp.Token[3:])
	authCtx := metadata.NewOutgoingContext(context.Background(), md)

	return authCtx
}

func cleanupDummies() {
	toCleanup := [...]dummyUser{createDummy, updateDummy}

	for _, user := range toCleanup {
		log.Printf("deleting dummy user %d\n", user.id)
		deleteReq := &pb.DeleteUserRequest{
			Id: user.id,
		}

		_, err := client.DeleteUser(loggedInContext(ADMIN), deleteReq)
		if err != nil {
			log.Printf("Failed to delete dummy user %d: %s\n", user.id, err)
		}
	}
}

func TestMain(m *testing.M) {
	const grpcAddr = "127.0.0.1:8080"
	generateDummyUsers()

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC service: %v", err)
	}
	defer conn.Close()

	client = pb.NewAuthServiceClient(conn)

	exitCode := m.Run()
	defer os.Exit(exitCode)
	cleanupDummies()
}

//nolint:funlen // Function length is justified
func TestCreateUser(t *testing.T) {
	err := createDummyUser(&createDummy)
	if err != nil {
		t.Fatalf("Failed to create dummy user for Update test cases: %s\n", err)
	}

	ctx := context.Background()
	// Test correct user create request
	type expectation struct {
		out    *emptypb.Empty
		status codes.Code
	}

	tests := map[string]struct {
		in       *pb.CreateUserRequest
		expected expectation
	}{
		// Test email with missing TLD
		"Invalid_Email_No_TLD": {
			in: &pb.CreateUserRequest{
				Email:    strings.ReplaceAll(createDummy.email, ".com", ""),
				Username: createDummy.username,
				Password: createDummy.password,
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.InvalidArgument,
			},
		},

		// Test email with missing domain
		"Invalid_Email_No_Domain": {
			in: &pb.CreateUserRequest{
				Email:    strings.ReplaceAll(createDummy.email, "@gmail.com", ""),
				Username: createDummy.username,
				Password: createDummy.password,
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.InvalidArgument,
			},
		},

		// Test register duplicate user with duplicate email
		"Already_Registered_Email": {
			in: &pb.CreateUserRequest{
				Email:    userEmail,
				Username: createDummy.username + "_new",
				Password: createDummy.password,
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.AlreadyExists,
			},
		},

		"Already_Registered_Username": {
			in: &pb.CreateUserRequest{
				Email:    strings.ReplaceAll(createDummy.email, "@gmail.com", "@hotmail.com"),
				Username: createDummy.username,
				Password: createDummy.password,
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.AlreadyExists,
			},
		},

		// Test register with missing arguments
		"Missing_Email": {
			in: &pb.CreateUserRequest{
				Email:    "",
				Username: createDummy.username,
				Password: createDummy.password,
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.InvalidArgument,
			},
		},

		"Missing_Username": {
			in: &pb.CreateUserRequest{
				Email:    createDummy.email,
				Username: "",
				Password: createDummy.password,
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.InvalidArgument,
			},
		},

		"Missing_Password": {
			in: &pb.CreateUserRequest{
				Email:    createDummy.email,
				Username: createDummy.password,
				Password: "",
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.InvalidArgument,
			},
		},

		"Missing_All_Arguments": {
			in: &pb.CreateUserRequest{
				Email:    "",
				Username: "",
				Password: "",
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.InvalidArgument,
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.CreateUser(ctx, tt.in)
			if err != nil {
				st, _ := status.FromError(err)
				outCode := st.Code()

				if tt.expected.status == codes.OK {
					t.Fatalf("Err -> \nWant: nil\nGot: %q\nWith:%s\n%s\n%s", err, tt.in.Username, tt.in.Email, tt.in.Password)
				}
				if tt.expected.status != outCode {
					t.Errorf("Err -> \nWant: %q\nGot: %q\nWith:%s\n%s\n%s", tt.expected.status, err, tt.in.Username, tt.in.Email, tt.in.Password)
				}
			} else if tt.expected.out.String() != out.String() {
				t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	// Test correct user create request
	type expectation struct {
		out *pb.LoginResponse
		err error
	}

	tests := map[string]struct {
		in       *pb.LoginRequest
		expected expectation
	}{
		"Must_Success": {
			in: &pb.LoginRequest{
				Username: createDummy.username,
				Password: createDummy.password,
			},
			expected: expectation{
				out: &pb.LoginResponse{},
				err: nil,
			},
		},

		"Incorrect_Username": {
			in: &pb.LoginRequest{
				Username: userUsername + "_wrong",
				Password: userPassword,
			},
			expected: expectation{
				out: &pb.LoginResponse{},
				err: status.Error(codes.NotFound, "invalid username or password"),
			},
		},

		"Incorrect_Password": {
			in: &pb.LoginRequest{
				Username: userUsername,
				Password: userPassword + "_wrong",
			},
			expected: expectation{
				out: &pb.LoginResponse{},
				err: status.Error(codes.PermissionDenied, "invalid username or password"),
			},
		},

		"Missing_Username": {
			in: &pb.LoginRequest{
				Username: "",
				Password: userPassword,
			},
			expected: expectation{
				out: &pb.LoginResponse{},
				err: status.Error(codes.InvalidArgument, "Invalid request, missing some request parameter(s)."),
			},
		},

		"Missing_Password": {
			in: &pb.LoginRequest{
				Username: userUsername,
				Password: "",
			},
			expected: expectation{
				out: &pb.LoginResponse{},
				err: status.Error(codes.InvalidArgument, "Invalid request, missing some request parameter(s)."),
			},
		},

		"Missing_All_Arguments": {
			in: &pb.LoginRequest{
				Username: "",
				Password: "",
			},
			expected: expectation{
				out: &pb.LoginResponse{},
				err: status.Error(codes.InvalidArgument, "Invalid request, missing some request parameter(s)."),
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.Login(ctx, tt.in)
			if err != nil {
				if tt.expected.err == nil {
					t.Fatalf("Err -> \nWant: nil\nGot: %q\n", err)
				}
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if _, _, tokenErr := jwt.ValidateJWT(out.Token); tokenErr != nil {
					t.Errorf("Out -> Received invalid token : '%s'\n%s", out.Token, tokenErr)
				}
			}
		})
	}
}

//nolint:funlen // Function length is justified
func TestUpdateUser(t *testing.T) {
	err := createDummyUser(&updateDummy)
	if err != nil {
		t.Fatalf("Failed to create dummy user for Update test cases: %s\nUser: %v\n", err, updateDummy)
	}

	// Test correct user create request
	type expectation struct {
		out    *emptypb.Empty
		status codes.Code
	}

	tests := map[string]struct {
		ctx      context.Context
		in       *pb.UpdateUserRequest
		expected expectation
	}{
		"No_Changes": {
			in: &pb.UpdateUserRequest{
				Id:       userID,
				Email:    userEmail,
				Username: userUsername,
				Password: "",
				Role:     models.USER.RoleString(),
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.InvalidArgument,
			},
			ctx: loggedInContext(ADMIN),
		},

		"Must_Success": {
			in: &pb.UpdateUserRequest{
				Id:       updateDummy.id,
				Email:    strings.ReplaceAll(updateDummy.email, "@gmail.com", "@outlook.com"),
				Username: updateDummy.username,
				Password: updateDummy.password,
				Role:     models.USER.RoleString(),
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.OK,
			},
			ctx: loggedInContext(ADMIN),
		},

		"Field_Already_Exists": {
			in: &pb.UpdateUserRequest{
				Id:       updateDummy.id,
				Email:    userEmail,
				Username: userUsername,
				Password: updateDummy.password,
				Role:     models.USER.RoleString(),
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.AlreadyExists,
			},
			ctx: loggedInContext(ADMIN),
		},

		"Missing_All_Arguments": {
			in: &pb.UpdateUserRequest{},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.InvalidArgument,
			},
			ctx: loggedInContext(ADMIN),
		},

		"Invalid_Email": {
			in: &pb.UpdateUserRequest{
				Id:       updateDummy.id,
				Email:    strings.ReplaceAll(updateDummy.email, "test.com", ""),
				Username: updateDummy.username,
				Password: updateDummy.password,
				Role:     models.USER.RoleString(),
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.InvalidArgument,
			},
			ctx: loggedInContext(ADMIN),
		},

		"User_Not_Found": {
			in: &pb.UpdateUserRequest{
				Id:       99999, //nolint:mnd // Not magic
				Email:    "michaelmyers234@gmail.com",
				Username: "michaelmyers",
				Password: "fridaythe13th",
				Role:     models.USER.RoleString(),
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.NotFound,
			},
			ctx: loggedInContext(ADMIN),
		},

		"Missing_Auth_Token": {
			in: &pb.UpdateUserRequest{
				Id:       updateDummy.id,
				Email:    "michaelmyers234@gmail.com",
				Username: "michaelmyers",
				Password: "fridaythe13th",
				Role:     models.USER.RoleString(),
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.Unauthenticated,
			},
			ctx: context.Background(),
		},

		"Malformed_Token": {
			in: &pb.UpdateUserRequest{
				Id:       updateDummy.id,
				Email:    "michaelmyers234@gmail.com",
				Username: "michaelmyers",
				Password: "fridaythe13th",
				Role:     models.USER.RoleString(),
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.Unauthenticated,
			},
			ctx: malformedTokenContext(),
		},

		"Unauthorized_User": {
			in: &pb.UpdateUserRequest{
				Id:       updateDummy.id,
				Email:    "michaelmyers234@gmail.com",
				Username: "michaelmyers",
				Password: "fridaythe13th",
				Role:     models.USER.RoleString(),
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.PermissionDenied,
			},
			ctx: loggedInContext(USER),
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.UpdateUser(tt.ctx, tt.in)
			if err != nil {
				st, _ := status.FromError(err)
				outCode := st.Code()

				if tt.expected.status == codes.OK {
					t.Fatalf("Err -> \nWant: OK\nGot: %q\nErr msg: %s\n", outCode, err)
				}
				if tt.expected.status != outCode {
					t.Errorf("Err -> \nWant: %q\nGot: %q\nErr msg: %s\n", tt.expected.status, outCode, err)
				}
			} else if tt.expected.out.String() != out.String() {
				t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	err := createDummyUser(&deleteDummy)
	if err != nil {
		t.Fatalf("Failed to create dummy user for Delete test cases: %s\n", err)
	}

	// Test correct user create request
	type expectation struct {
		out    *emptypb.Empty
		status codes.Code
	}

	tests := map[string]struct {
		ctx      context.Context
		in       *pb.DeleteUserRequest
		expected expectation
	}{
		"Unauthorized": {
			in: &pb.DeleteUserRequest{},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.PermissionDenied,
			},
			ctx: loggedInContext(USER),
		},

		"Malformed_Token": {
			in: &pb.DeleteUserRequest{},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.Unauthenticated,
			},
			ctx: malformedTokenContext(),
		},

		"Missing_Argument": {
			in: &pb.DeleteUserRequest{},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.InvalidArgument,
			},
			ctx: loggedInContext(ADMIN),
		},

		"Missing_Auth_token": {
			in: &pb.DeleteUserRequest{},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.Unauthenticated,
			},
			ctx: context.Background(),
		},

		"User_Does_Not_Exist": {
			in: &pb.DeleteUserRequest{
				Id: 999999999, //nolint:mnd // Not magic
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.NotFound,
			},
			ctx: loggedInContext(ADMIN),
		},

		"Must_Success": {
			in: &pb.DeleteUserRequest{
				Id: deleteDummy.id,
			},
			expected: expectation{
				out:    &emptypb.Empty{},
				status: codes.OK,
			},
			ctx: loggedInContext(ADMIN),
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.DeleteUser(tt.ctx, tt.in)
			if err != nil {
				st, _ := status.FromError(err)
				outCode := st.Code()

				if tt.expected.status == codes.OK {
					t.Fatalf("Err -> \nWant: OK\nGot: %q\nErr msg: %s\n", outCode, err)
				}
				if tt.expected.status != outCode {
					t.Errorf("Err -> \nWant: %q\nGot: %q\nErr msg: %s\n", tt.expected.status, outCode, err)
				}
			} else if tt.expected.out.String() != out.String() {
				t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
			}
		})
	}
}

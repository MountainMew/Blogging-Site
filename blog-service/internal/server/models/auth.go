package models

import (
	pb "blog-service/internal/grpc/protobuf"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ContextKey string

type UserRegisterDTO struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponseDTO struct {
	Token string `json:"token"`
}

type UserClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	ID       int    `json:"id"`
}

const ClaimsKey ContextKey = "jwtClaims"

func GetClaimsFromContext(ctx context.Context) *UserClaims {
	claims, ok := ctx.Value(ClaimsKey).(*UserClaims)
	if !ok {
		return nil
	}
	return claims
}

func VerifyAuthToken(ctx context.Context, client pb.AuthServiceClient, token string) *UserClaims {
	verifyReq := &pb.VerifyTokenRequest{
		Token: token,
	}

	verifyRes, err := client.VerifyToken(ctx, verifyReq)

	if err != nil {
		return nil
	}

	return &UserClaims{
		Username: verifyRes.Username,
		ID:       int(verifyRes.Id),
		Role:     verifyRes.Role,
	}
}

func LoginUser(ctx context.Context, client pb.AuthServiceClient, user *UserLoginDTO) (string, error) {
	if user.Username == "" || user.Password == "" {
		return "", &ParamError{}
	}

	loginReq := pb.LoginRequest{
		Username: user.Username,
		Password: user.Password,
	}

	res, err := client.Login(ctx, &loginReq)

	if err != nil {
		st, _ := status.FromError(err)
		code := st.Code()

		//nolint:exhaustive // Other status codes not necessary
		switch code {
		case codes.PermissionDenied:
			return "", &InvalidLoginError{}
		case codes.NotFound:
			return "", &InvalidLoginError{}
		case codes.Internal:
			return "", &InvalidLoginError{}
		case codes.InvalidArgument:
			return "", &ParamError{}
		}

		return "", err
	}

	return res.Token, nil
}

func RegisterUser(ctx context.Context, client pb.AuthServiceClient, user *UserRegisterDTO) error {
	if user.Email == "" || user.Username == "" || user.Password == "" {
		return &ParamError{}
	}

	registerReq := pb.CreateUserRequest{
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
	}

	_, err := client.CreateUser(ctx, &registerReq)

	if err != nil {
		st, _ := status.FromError(err)
		code := st.Code()

		//nolint:exhaustive // Other status codes not necessary
		switch code {
		case codes.Internal:
			return err
		case codes.InvalidArgument:
			return &ParamError{}
		case codes.AlreadyExists:
			return &EmailOrUserTakenError{}
		}
	}

	return nil
}

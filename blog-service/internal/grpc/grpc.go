package grpc

import (
	pb "blog-service/internal/grpc/protobuf"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	Client pb.AuthServiceClient
}

func NewGRPCClient() (pb.AuthServiceClient, error) {
	addr := os.Getenv("AUTH_URI")
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewAuthServiceClient(conn)

	return client, nil
}

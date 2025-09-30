package main

import (
	"blog-service/internal/db/mongo"
	pg "blog-service/internal/db/postgres"
	"blog-service/internal/grpc"
	"blog-service/internal/server"
	"context"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Starting blog-service")

	// Create new mongodb client
	mongoClient, err := mongo.NewMongoClient()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err = mongoClient.Close(context.Background()); err != nil {
			log.Printf("Error closing MongoDB connection: %v", err)
		}
	}()

	// Create new postgres client
	postgresClient, err := pg.NewPostgresClient()
	if err != nil {
		log.Printf("Failed to create postgres client: %s", err)
		return
	}
	defer postgresClient.ConnPool.Close()

	// Create new gRPC client
	grpcClient, err := grpc.NewGRPCClient()
	if err != nil {
		log.Printf("Failed to create gRPC client: %s", err)
		return
	}

	// Create an instance of server
	srv := server.NewServer(mongoClient, postgresClient, grpcClient)

	// Start the server on port 8080
	if err := srv.Start(":8080"); err != nil {
		log.Printf("Server failed to start: %v", err)
	}

	select {}
}

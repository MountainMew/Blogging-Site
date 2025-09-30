package main

import (
	"auth-service/internal/db"
	"auth-service/internal/server"
	"context"
	"fmt"
)

func main() {
	fmt.Println("Starting auth-service")

	ctx := context.Background()

	// Config database
	database, err := db.Config()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Open connection
	err = database.OpenDbConnection()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Defer closing connection
	defer database.ConnPool.Close()

	// Initilize database
	err = database.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	grpcServer := server.NewGRPCServer(database)
	grpcServer.ListenAndServe(ctx)
}

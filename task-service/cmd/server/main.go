package main

import (
	"log"

	"github.com/Evilcodercpp/task-service/internal/database"
	"github.com/Evilcodercpp/task-service/internal/task"
	transportgrpc "github.com/Evilcodercpp/task-service/internal/transport/grpc"
)

func main() {
	database.InitDB()

	repo := task.NewRepository(database.DB)
	svc := task.NewService(repo)

	userClient, conn, err := transportgrpc.NewUserClient("localhost:50051")
	if err != nil {
		log.Fatalf("failed to connect to users: %v", err)
	}
	defer conn.Close()

	if err := transportgrpc.RunGRPC(svc, userClient); err != nil {
		log.Fatalf("Tasks gRPC server error: %v", err)
	}
}

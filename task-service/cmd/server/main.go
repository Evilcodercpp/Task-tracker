package main

import (
	"log"

	"github.com/Evilcodercpp/task-service/internal/database"
	transportgrpc "github.com/Evilcodercpp/task-service/internal/transport/grpc"
	"github.com/Evilcodercpp/task-service/internal/task"
)

func main() {
	database.InitDB()
	repo := task.NewRepository(database.DB)
	svc := task.NewService(repo)

	if err := transportgrpc.RunGRPC(svc); err != nil {
		log.Fatalf("gRPC сервер завершился с ошибкой: %v", err)
	}
}

package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	db "Task-tracker/interal"
	"Task-tracker/interal/handlers"
	taskservice "Task-tracker/interal/taskService"
	"Task-tracker/interal/web/tasks"
)

func main() {
	db.InitDB()

	repo := taskservice.NewTaskRepo(db.DB)
	svc := taskservice.NewTaskService(repo)
	h := handlers.NewHandler(svc)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	strictHandler := tasks.NewStrictHandler(h, nil)
	tasks.RegisterHandlers(e, strictHandler)

	// TODO: нет graceful shutdown.
	// При Ctrl+C сервер убивается мгновенно, не дожидаясь завершения активных запросов.
	// Решение: использовать e.Shutdown(ctx) через signal.NotifyContext.
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("failed to start with err: %v", err)
	}
}
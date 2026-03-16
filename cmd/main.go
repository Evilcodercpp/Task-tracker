package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	db "Task-tracker/interal"
	"Task-tracker/interal/handlers"
	taskservice "Task-tracker/interal/taskService"
	userservice "Task-tracker/interal/userService"
	"Task-tracker/interal/web/tasks"
	"Task-tracker/interal/web/users"
)

func main() {
	db.InitDB()

	tasksRepo := taskservice.NewTaskRepo(db.DB)
	tasksSvc := taskservice.NewTaskService(tasksRepo)
	tasksHandler := handlers.NewHandler(tasksSvc)

	userRepo := userservice.NewUserRepo(db.DB)
	userSvc := userservice.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userSvc)

	e := echo.New()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Printf("%s %s %d\n", v.Method, v.URI, v.Status)
			return nil
		},
	}))
	e.Use(middleware.Recover())

	strictTasksHandler := tasks.NewStrictHandler(tasksHandler, nil)
	tasks.RegisterHandlers(e, strictTasksHandler)

	strictUserHandler := users.NewStrictHandler(userHandler, nil)
	users.RegisterHandlers(e, strictUserHandler)

	// TODO: нет graceful shutdown.
	// При Ctrl+C сервер убивается мгновенно, не дожидаясь завершения активных запросов.
	// Решение: использовать e.Shutdown(ctx) через signal.NotifyContext.
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("failed to start with err: %v", err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"Task-tracker/interal"
	"Task-tracker/interal/handlers"
	taskservice "Task-tracker/interal/taskService"
)

func main() {
	// 1. Подключаемся к базе данных
	db.InitDB()

	// 1a. Автоматически создаём/обновляем таблицу task на основе модели
	if err := db.DB.AutoMigrate(&taskservice.Task{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	// 2. Создаём репозиторий — работает напрямую с БД
	repo := taskservice.NewTaskRepo(db.DB)

	// 3. Создаём сервис — содержит бизнес-логику, использует репозиторий
	svc := taskservice.NewTaskService(repo)

	// 4. Создаём хендлер — обрабатывает HTTP-запросы, использует сервис
	// Цепочка: Handler → TaskService → TaskRepository → DB
	h := handlers.NewHandler(svc)

	// 5. Регистрируем маршрут /task
	http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetTasks(w, r)
		case http.MethodPost:
			h.PostTask(w, r)
		case http.MethodPatch:
			h.PatchTask(w, r)
		case http.MethodDelete:
			h.DeleteTask(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// 6. Запускаем сервер
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

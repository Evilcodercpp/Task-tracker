package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// db — глобальное подключение к базе данных, одно на всё приложение
var db *gorm.DB

// Task — модель задачи, соответствует таблице в базе данных
type Task struct {
	ID   string `gorm:"primaryKey" json:"id"` // уникальный идентификатор, генерируется на сервере
	Task string `json:"task"`                 // текст задачи
}

// initDB — подключается к PostgreSQL и создаёт таблицу если её нет
func initDB() {
	// берём строку подключения из переменной окружения
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// если переменная не задана — используем локальные настройки
		dsn = "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	}

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	// AutoMigrate создаёт таблицу tasks если её нет, или добавляет новые колонки
	if err := db.AutoMigrate(&Task{}); err != nil {
		log.Fatalf("could not migrate database: %v", err)
	}
}

// getTask — GET /task
// Возвращает все задачи из базы данных в виде JSON-массива
func getTask(w http.ResponseWriter, r *http.Request) {
	var tasks []Task
	// SELECT * FROM tasks
	if err := db.Find(&tasks).Error; err != nil {
		http.Error(w, "could not fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// postTask — POST /task
// Создаёт новую задачу. ID генерируется автоматически.
// Тело запроса: {"task": "название задачи"}
func postTask(w http.ResponseWriter, r *http.Request) {
	var body Task
	// читаем JSON из тела запроса и записываем в body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// генерируем уникальный ID, клиент его не передаёт
	body.ID = uuid.NewString()

	// INSERT INTO tasks (id, task) VALUES (...)
	if err := db.Create(&body).Error; err != nil {
		http.Error(w, "could not create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(body)
}

// patchTask — PATCH /task
// Обновляет текст задачи по ID. Поле task обновляется только если оно передано.
// Тело запроса: {"id": "...", "task": "новое название"}
func patchTask(w http.ResponseWriter, r *http.Request) {
	var body Task
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// ID обязателен — без него не знаем какую задачу обновлять
	if body.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	var task Task
	// SELECT * FROM tasks WHERE id = ? LIMIT 1
	if err := db.First(&task, "id = ?", body.ID).Error; err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	// обновляем только если поле передано (защита от затирания пустой строкой)
	if body.Task != "" {
		task.Task = body.Task
	}

	// UPDATE tasks SET task = ? WHERE id = ?
	if err := db.Save(&task).Error; err != nil {
		http.Error(w, "could not update task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// deleteTask — DELETE /task
// Удаляет задачу по ID.
// Тело запроса: {"id": "..."}
func deleteTask(w http.ResponseWriter, r *http.Request) {
	var body Task
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// ID обязателен — без него не знаем что удалять
	if body.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	// DELETE FROM tasks WHERE id = ?
	result := db.Delete(&Task{}, "id = ?", body.ID)
	if result.Error != nil {
		http.Error(w, "could not delete task", http.StatusInternalServerError)
		return
	}
	// RowsAffected == 0 значит задачи с таким ID не существовало
	if result.RowsAffected == 0 {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	initDB()

	// один маршрут /task обрабатывает все методы
	http.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTask(w, r)
		case http.MethodPost:
			postTask(w, r)
		case http.MethodPatch:
			patchTask(w, r)
		case http.MethodDelete:
			deleteTask(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed) // 405 — метод не поддерживается
		}
	})

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

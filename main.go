package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// requestBody — структура задачи.
// Используется и для хранения, и для парсинга тела запроса.
type requestBody struct {
	ID   string `json:"id"`
	Task string `json:"task"`
}

// tasks — хранилище задач в памяти (сбрасывается при перезапуске сервера)
var tasks = []requestBody{}

// getTask — GET /task
// Возвращает все задачи в виде JSON-массива.
func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// postTask — POST /task
// Создаёт новую задачу. ID генерируется автоматически.
// Тело запроса: {"task": "название задачи"}
func postTask(w http.ResponseWriter, r *http.Request) {
	var body requestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	body.ID = uuid.NewString() // генерируем уникальный ID
	tasks = append(tasks, body)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(body)
}

// patchTask — PATCH /task
// Обновляет задачу по ID. Поле task обновляется только если оно передано.
// Тело запроса: {"id": "...", "task": "новое название"}
func patchTask(w http.ResponseWriter, r *http.Request) {
	var body requestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if body.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	for i := range tasks {
		if tasks[i].ID == body.ID {
			if body.Task != "" {
				tasks[i].Task = body.Task
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}

	http.Error(w, "task not found", http.StatusNotFound)
}

// deleteTask — DELETE /task
// Удаляет задачу по ID.
// Тело запроса: {"id": "..."}
func deleteTask(w http.ResponseWriter, r *http.Request) {
	var body requestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if body.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	for i := range tasks {
		if tasks[i].ID == body.ID {
			// вырезаем элемент из среза: склеиваем всё до i и всё после i
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "task not found", http.StatusNotFound)
}

func main() {
	// регистрируем один маршрут /task, метод определяется внутри
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
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	// запускаем сервер на порту 8080
	fmt.Println("Server started on :8080")
	// log.Fatal завершит программу, если ListenAndServe вернёт ошибку (например, если порт уже занят) или если сервер будет остановлен. Это полезно для отладки и предотвращения зависания.
	log.Fatal(http.ListenAndServe(":8080", nil))
}

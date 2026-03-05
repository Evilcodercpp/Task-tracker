package handlers

import (
	"encoding/json"
	"net/http"

	taskservice "Task-tracker/interal/taskService"

	"github.com/google/uuid"
)

// Handler хранит зависимости для HTTP-хендлеров
type Handler struct {
	repo taskservice.TaskRepository
}

// NewHandler создаёт хендлер с переданным репозиторием
func NewHandler(repo taskservice.TaskRepository) *Handler {
	return &Handler{repo: repo}
}

// GetTasks — GET /task
// Возвращает все задачи из БД в виде JSON-массива
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.repo.GetAllTasks()
	if err != nil {
		http.Error(w, "could not fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// PostTask — POST /task
// Создаёт новую задачу. ID генерируется на сервере.
// Тело запроса: {"task": "...", "is_done": false}
func (h *Handler) PostTask(w http.ResponseWriter, r *http.Request) {
	var body taskservice.Task
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	body.ID = uuid.NewString()

	if err := h.repo.CreateTask(&body); err != nil {
		http.Error(w, "could not create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(body)
}

// PatchTask — PATCH /task
// Обновляет задачу по ID.
// Тело запроса: {"id": "...", "task": "...", "is_done": true}
func (h *Handler) PatchTask(w http.ResponseWriter, r *http.Request) {
	var body taskservice.Task
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if body.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	task, err := h.repo.GetTaskByID(body.ID)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if body.Task != "" {
		task.Task = body.Task
	}
	task.IsDone = body.IsDone

	if err := h.repo.UpdateTask(task); err != nil {
		http.Error(w, "could not update task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// DeleteTask — DELETE /task
// Мягко удаляет задачу по ID (проставляет deleted_at).
// Тело запроса: {"id": "..."}
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	var body taskservice.Task
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if body.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteTask(body.ID); err != nil {
		http.Error(w, "could not delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

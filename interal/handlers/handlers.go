package handlers

import (
	"encoding/json"
	"net/http"

	taskservice "Task-tracker/interal/taskService"

	"github.com/google/uuid"
)

// Handler содержит сервисный слой.
// Хендлеры не работают с БД напрямую — только через TaskService.
type Handler struct {
	svc *taskservice.TaskService
}

// NewHandler создаёт Handler и принимает сервис через аргумент.
func NewHandler(svc *taskservice.TaskService) *Handler {
	return &Handler{svc: svc}
}

// GetTasks — GET /task
// Возвращает все задачи из БД в виде JSON-массива.
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.svc.GetAllTasks()
	if err != nil {
		http.Error(w, "could not fetch tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// PostTask — POST /task
// Создаёт новую задачу. ID генерируется на сервере.
// Тело запроса: {"task": "текст задачи", "is_done": false}
func (h *Handler) PostTask(w http.ResponseWriter, r *http.Request) {
	var body taskservice.Task

	// читаем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// генерируем уникальный ID на сервере
	body.ID = uuid.NewString()

	if err := h.svc.CreateTask(&body); err != nil {
		http.Error(w, "could not create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(body)
}

// PatchTask — PATCH /task
// Обновляет существующую задачу по ID.
// Тело запроса: {"id": "...", "task": "новый текст", "is_done": true}
func (h *Handler) PatchTask(w http.ResponseWriter, r *http.Request) {
	var body taskservice.Task
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// без ID не знаем какую задачу обновлять
	if body.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	// достаём задачу из БД через сервис
	task, err := h.svc.GetTaskByID(body.ID)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	// обновляем только переданные поля
	if body.Task != "" {
		task.Task = body.Task
	}
	task.IsDone = body.IsDone

	if err := h.svc.UpdateTask(task); err != nil {
		http.Error(w, "could not update task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// DeleteTask — DELETE /task
// Мягко удаляет задачу: заполняет deleted_at, строка остаётся в БД.
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

	if err := h.svc.DeleteTask(body.ID); err != nil {
		http.Error(w, "could not delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 — успех, тело не нужно
}

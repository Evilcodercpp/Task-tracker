package handlers

import (
	"context"
	"errors"

	"Task-tracker/interal/web/tasks"
	taskservice "Task-tracker/interal/taskService"

	"github.com/google/uuid"
)

// TODO: svc должен быть интерфейсом, а не конкретным типом.
// Сейчас Handler зависит напрямую от *taskservice.TaskService,
// что делает невозможным подмену в тестах (мокирование).
// Решение: объявить интерфейс TaskService и принимать его здесь.
type Handler struct {
	svc *taskservice.TaskService
}

func NewHandler(svc *taskservice.TaskService) *Handler {
	return &Handler{svc: svc}
}

// GetTasks возвращает список всех задач.
func (h *Handler) GetTasks(_ context.Context, _ tasks.GetTasksRequestObject) (tasks.GetTasksResponseObject, error) {
	allTasks, err := h.svc.GetAllTasks()
	if err != nil {
		return nil, err
	}

	response := tasks.GetTasks200JSONResponse{}

	for _, tsk := range allTasks {
		id := tsk.ID
		text := tsk.Task
		done := tsk.IsDone
		response = append(response, tasks.Task{
			Id:     &id,
			Task:   &text,
			IsDone: &done,
		})
	}

	return response, nil
}

// GetTasksUserUserId возвращает все задачи конкретного пользователя.
func (h *Handler) GetTasksUserUserId(_ context.Context, request tasks.GetTasksUserUserIdRequestObject) (tasks.GetTasksUserUserIdResponseObject, error) {
	tasksByUser, err := h.svc.GetTasksByUserID(request.UserId)
	if err != nil {
		return nil, err
	}

	response := tasks.GetTasksUserUserId200JSONResponse{}

	for _, tsk := range tasksByUser {
		id := tsk.ID
		text := tsk.Task
		done := tsk.IsDone
		userID := tsk.UserID
		response = append(response, tasks.Task{
			Id:     &id,
			Task:   &text,
			IsDone: &done,
			UserId: &userID,
		})
	}

	return response, nil
}

// PostTasks создаёт новую задачу.
func (h *Handler) PostTasks(_ context.Context, request tasks.PostTasksRequestObject) (tasks.PostTasksResponseObject, error) {
	// Проверка на nil: если клиент не передал поле task в теле запроса,
	// разыменование указателя вызовет panic.
	if request.Body.Task == nil {
		return nil, errors.New("field 'task' is required")
	}

	if request.Body.UserId == nil {
		return nil, errors.New("field 'user_id' is required")
	}

	taskToCreate := &taskservice.Task{
		ID:     uuid.NewString(),
		Task:   *request.Body.Task,
		IsDone: false,
		UserID: *request.Body.UserId,
	}

	if err := h.svc.CreateTask(taskToCreate); err != nil {
		return nil, err
	}

	id := taskToCreate.ID
	text := taskToCreate.Task
	done := taskToCreate.IsDone

	userID := taskToCreate.UserID

	return tasks.PostTasks201JSONResponse{
		Id:     &id,
		Task:   &text,
		IsDone: &done,
		UserId: &userID,
	}, nil
}

func (h *Handler) PatchTasksId(_ context.Context, request tasks.PatchTasksIdRequestObject) (tasks.PatchTasksIdResponseObject, error) {
	existing, err := h.svc.GetTaskByID(request.Id)
	if err != nil {
		return nil, err
	}

	if request.Body.Task != nil {
		existing.Task = *request.Body.Task
	}
	if request.Body.IsDone != nil {
		existing.IsDone = *request.Body.IsDone
	}

	if err := h.svc.UpdateTask(existing); err != nil {
		return nil, err
	}

	id := existing.ID
	text := existing.Task
	done := existing.IsDone

	return tasks.PatchTasksId200JSONResponse{
		Id:     &id,
		Task:   &text,
		IsDone: &done,
	}, nil
}

func (h *Handler) DeleteTasksId(_ context.Context, request tasks.DeleteTasksIdRequestObject) (tasks.DeleteTasksIdResponseObject, error) {
	if err := h.svc.DeleteTask(request.Id); err != nil {
		return nil, err
	}

	return tasks.DeleteTasksId204Response{}, nil
}
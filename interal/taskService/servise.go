package taskservice

import (
	"errors"
	"log"
)

// TaskService — сервисный слой.
// Находится между хендлерами и репозиторием.
// Здесь живёт бизнес-логика: валидация, проверки, правила.
type TaskService struct {
	repo TaskRepository
}

// NewTaskService создаёт сервис и принимает репозиторий через аргумент.
func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// CreateTask создаёт новую задачу.
// Валидирует текст перед сохранением.
func (s *TaskService) CreateTask(task *Task) error {
	if task.Task == "" {
		return errors.New("task text cannot be empty")
	}
	if len(task.Task) > 255 {
		return errors.New("task text cannot exceed 255 characters")
	}
	task.IsDone = false // новая задача всегда не выполнена
	return s.repo.CreateTask(task)
}

// GetAllTasks возвращает все задачи.
// Логирует если задач нет.
func (s *TaskService) GetAllTasks() ([]Task, error) {
	tasks, err := s.repo.GetAllTasks()
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		log.Println("no tasks found")
	}
	return tasks, nil
}

// GetTaskByID возвращает одну задачу по ID.
// Проверяет что ID не пустой.
func (s *TaskService) GetTaskByID(id string) (*Task, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return s.repo.GetTaskByID(id)
}

// UpdateTask обновляет задачу.
// Проверяет существование задачи и валидирует текст.
func (s *TaskService) UpdateTask(task *Task) error {
	if task.Task == "" {
		return errors.New("task text cannot be empty")
	}
	existing, err := s.repo.GetTaskByID(task.ID)
	if err != nil {
		return errors.New("task not found")
	}
	existing.Task = task.Task
	existing.IsDone = task.IsDone
	return s.repo.UpdateTask(existing)
}

// DeleteTask мягко удаляет задачу по ID.
// Логирует результат операции.
func (s *TaskService) DeleteTask(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	err := s.repo.DeleteTask(id)
	if err != nil {
		log.Printf("failed to delete task %s: %v", id, err)
		return err
	}
	log.Printf("task %s deleted successfully", id)
	return nil
}

package task

import (
	"errors"
	"log"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTask(t *Task) error {
	if t.Title == "" {
		return errors.New("title cannot be empty")
	}
	t.IsDone = false
	return s.repo.Create(t)
}

func (s *Service) GetAllTasks(page, pageSize int) ([]Task, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	tasks, total, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	if len(tasks) == 0 {
		log.Println("no tasks found")
	}
	return tasks, total, nil
}

func (s *Service) GetTaskByID(id uint) (*Task, error) {
	if id == 0 {
		return nil, errors.New("id cannot be zero")
	}
	return s.repo.GetByID(id)
}

func (s *Service) UpdateTaskByID(id uint, updates Task) (*Task, error) {
	if updates.Title == "" {
		return nil, errors.New("title cannot be empty")
	}
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	existing.Title = updates.Title
	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteTaskByID(id uint) error {
	if id == 0 {
		return errors.New("id cannot be zero")
	}
	err := s.repo.Delete(id)
	if err != nil {
		log.Printf("failed to delete task %d: %v", id, err)
		return err
	}
	log.Printf("task %d deleted successfully", id)
	return nil
}

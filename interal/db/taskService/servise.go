package taskservice

import "errors"

// TaskService содержит бизнес-логику работы с задачами
type TaskService struct {
	repo TaskRepository
}

func NewTaskService(r TaskRepository) *TaskService {
	return &TaskService{repo: r}
}

// Complete помечает задачу как выполненную
func (s *TaskService) Complete(id string) error {
	t, err := s.repo.GetTaskByID(id)
	if err != nil {
		return err
	}
	if t.IsDone {
		return errors.New("task is already done")
	}
	t.IsDone = true
	return s.repo.UpdateTask(t)
}

// Reopen снимает отметку выполнения с задачи
func (s *TaskService) Reopen(id string) error {
	t, err := s.repo.GetTaskByID(id)
	if err != nil {
		return err
	}
	t.IsDone = false
	return s.repo.UpdateTask(t)
}

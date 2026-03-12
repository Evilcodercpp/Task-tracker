package taskservice

import "gorm.io/gorm"

// основные методы CRUD - create, read, update, delete
type TaskRepository interface {
	CreateTask(tas *Task) error
	GetAllTasks() ([]Task, error)
	GetTaskByID(ID string) (*Task, error)
	UpdateTask(tas *Task) error
	DeleteTask(ID string) error
}

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) TaskRepository {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) CreateTask(tas *Task) error {
	return r.db.Create(tas).Error
}

func (r *TaskRepo) GetAllTasks() ([]Task, error) {
	var tasks []Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepo) GetTaskByID(ID string) (*Task, error) {
	var tas Task
	if err := r.db.First(&tas, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &tas, nil
}

func (r *TaskRepo) UpdateTask(tas *Task) error {
	return r.db.Save(tas).Error
}

func (r *TaskRepo) DeleteTask(ID string) error {
	return r.db.Delete(&Task{}, "id = ?", ID).Error
}

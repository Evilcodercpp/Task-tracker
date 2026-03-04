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

type TasRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) TaskRepository {
	return &TasRepo{db: db}
}

func (r *TasRepo) CreateTask(tas *Task) error {
	return r.db.Create(tas).Error
}

func (r *TasRepo) GetAllTasks() ([]Task, error) {
	var tasks []Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *TasRepo) GetTaskByID(ID string) (*Task, error) {
	var tas Task
	if err := r.db.First(&tas, "id = ?", ID).Error; err != nil {
		return nil, err
	}
	return &tas, nil
}

func (r *TasRepo) UpdateTask(tas *Task) error {
	return r.db.Save(tas).Error
}

func (r *TasRepo) DeleteTask(ID string) error {
	return r.db.Delete(&Task{}, "id = ?", ID).Error
}

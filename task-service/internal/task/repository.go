package task

import "gorm.io/gorm"

type Repository interface {
	Create(task *Task) error
	GetAll(page, pageSize int) ([]Task, int64, error)
	GetByID(id uint) (*Task, error)
	GetByUserID(userID uint) ([]Task, int64, error)
	Update(task *Task) error
	Delete(id uint) error
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}

func (r *repo) Create(task *Task) error {
	return r.db.Create(task).Error
}

func (r *repo) GetAll(page, pageSize int) ([]Task, int64, error) {
	var tasks []Task
	var total int64

	if err := r.db.Model(&Task{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Find(&tasks).Error
	return tasks, total, err
}

func (r *repo) GetByID(id uint) (*Task, error) {
	var task Task
	if err := r.db.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *repo) GetByUserID(userID uint) ([]Task, int64, error) {
	var tasks []Task
	err := r.db.Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, int64(len(tasks)), err
}

func (r *repo) Update(task *Task) error {
	return r.db.Save(task).Error
}

func (r *repo) Delete(id uint) error {
	return r.db.Delete(&Task{}, id).Error
}

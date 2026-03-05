package taskservice

import (
	"gorm.io/gorm"
)

const (
	StatusNew        = "new"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
	StatusCanceled   = "canceled"
)

// Task — модель задачи, соответствует таблице в базе данных
type Task struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Task      string         `json:"task"`
	IsDone    bool           `json:"is_done"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // мягкое удаление
}

// Разрешённые статусы
var allowedStatuses = map[string]bool{
	StatusNew:        true,
	StatusInProgress: true,
	StatusDone:       true,
	StatusCanceled:   true,
}

// IsValidStatus проверяет допустимость статуса
func IsValidStatus(status string) bool {
	return allowedStatuses[status]
}

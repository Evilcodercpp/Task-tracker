package taskservice

import (
	"gorm.io/gorm"
)

// Task — модель задачи, соответствует таблице в базе данных
type Task struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Task      string         `json:"task"`
	IsDone    bool           `json:"is_done"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // мягкое удаление
}

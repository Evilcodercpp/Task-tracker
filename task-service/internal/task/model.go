package task

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string         `json:"title"`
	IsDone    bool           `json:"is_done"`
	UserID    uint           `json:"user_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Task) TableName() string {
	return "tasks"
}

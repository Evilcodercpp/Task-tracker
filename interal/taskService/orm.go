package taskservice

import (
	"gorm.io/gorm"
)

// Task — модель задачи, соответствует таблице в базе данных.
//
// TODO: поле Task.Task — поле называется так же, как сама структура.
// Это создаёт путаницу: task.Task — это текст задачи или вложенная задача?
// Лучше переименовать поле в Title или Description, и обновить JSON-тег.
type Task struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Task      string         `json:"task"`
	IsDone    bool           `json:"is_done"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // мягкое удаление
	UserID    string         `gorm:"column:user_id" json:"user_id"`
}

func (Task) TableName() string {
	return "tasks"
}

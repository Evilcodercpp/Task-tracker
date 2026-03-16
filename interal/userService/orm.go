package userservice

import (
	"time"

	taskservice "Task-tracker/interal/taskService"

	"gorm.io/gorm"
)

// User — модель пользователя, соответствует таблице users в базе данных.
type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Tasks []taskservice.Task `gorm:"foreignKey:UserID" json:"tasks,omitempty"`
}

// TableName явно задаёт имя таблицы.
// Нужно потому что gorm с SingularTable:true маппит User → "user",
// а "user" — зарезервированное слово в PostgreSQL.
func (User) TableName() string {
	return "users"
}

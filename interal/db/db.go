package db

import (
	"Task-tracker/interal/db/taskService"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// Инициализация базы данных
func InitDB() {
	dsn := "host=localhost user=postgres password=yourpassword dbname=postgres port=5433 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Автоматическая миграция таблицы
	if err := DB.AutoMigrate(&taskservice.Task{}); err != nil {
		log.Fatalf("Could not migrate database: %v", err)
	}
}

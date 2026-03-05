package db

import (
	"Task-tracker/interal/taskService"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

// Инициализация базы данных
func InitDB() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5434 sslmode=disable"
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

	if err := DB.AutoMigrate(&taskservice.Task{}); err != nil {
		log.Fatalf("Could not migrate database: %v", err)
	}
}

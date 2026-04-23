package database

import (
	"log"

	"github.com/Evilcodercpp/task-service/internal/task"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=localhost user=milka dbname=postgres port=5432 sslmode=disable"
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
	if err := DB.AutoMigrate(&task.Task{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
}

package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// TODO: глобальная переменная — нарушение принципа DI.
// Лучше возвращать *gorm.DB из InitDB и передавать его явно в зависимости.
var DB *gorm.DB

// InitDB инициализирует подключение к базе данных.
//
// TODO: DSN захардкожен — пароль и параметры подключения находятся прямо в коде.
// Это критическая уязвимость: секреты не должны попадать в репозиторий.
// Решение: использовать переменные окружения (os.Getenv) или .env-файл (например, godotenv).
//
// Пример:
//   dsn := os.Getenv("DATABASE_URL")
func InitDB() {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5434 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// TODO: logger.Silent скрывает все SQL-запросы и ошибки.
		// В разработке лучше использовать logger.Info,
		// в проде — настраиваемый уровень через переменную окружения.
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
}

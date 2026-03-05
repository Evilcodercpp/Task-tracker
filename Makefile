# Makefile для Task Tracker

# Переменные
DB_DSN := "postgres://postgres:postgres@localhost:5434/postgres?sslmode=disable"
MIGRATE := migrate -path ./migrations -database $(DB_DSN)

# Создание новой миграции: make migrate-new NAME=create_tasks
migrate-new:
	migrate create -ext sql -dir ./migrations ${NAME}

# Применение миграций
migrate:
	$(MIGRATE) up

# Откат миграций
migrate-down:
	$(MIGRATE) down

# Запуск сервера
run:
	go run cmd/main.go

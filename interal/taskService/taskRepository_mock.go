package taskservice

import (
	"github.com/stretchr/testify/mock"
)

// MockTaskRepository — фейковая реализация интерфейса TaskRepository.
// Используется в юнит-тестах вместо реальной базы данных.
// Внутри хранит список ожидаемых вызовов (через mock.Mock),
// что позволяет проверять: какой метод вызван, с какими аргументами,
// и что он должен вернуть.
type MockTaskRepository struct {
	mock.Mock
}

// CreateTask имитирует создание задачи.
// Возвращает только ошибку — как в реальном интерфейсе.
func (m *MockTaskRepository) CreateTask(task *Task) error {
	args := m.Called(task)
	return args.Error(0)
}

// GetAllTasks имитирует получение всех задач из базы.
// m.Called() сигнализирует mock-фреймворку, что метод был вызван,
// и возвращает заранее заданные значения из теста.
func (m *MockTaskRepository) GetAllTasks() ([]Task, error) {
	args := m.Called()
	var tasks []Task
	if res := args.Get(0); res != nil {
		tasks = res.([]Task)
	}
	return tasks, args.Error(1)
}

// GetTaskByID имитирует поиск задачи по строковому ID.
// Проверка на nil нужна, чтобы не упасть с паникой при type assertion.
func (m *MockTaskRepository) GetTaskByID(ID string) (*Task, error) {
	args := m.Called(ID)
	var t *Task
	if res := args.Get(0); res != nil {
		t = res.(*Task)
	}
	return t, args.Error(1)
}

// UpdateTask имитирует обновление задачи.
// Принимает готовую задачу и возвращает только ошибку.
func (m *MockTaskRepository) UpdateTask(task *Task) error {
	args := m.Called(task)
	return args.Error(0)
}

// DeleteTask имитирует удаление задачи по строковому ID.
// Возвращает только ошибку (nil — если удаление "прошло успешно" в тесте).
func (m *MockTaskRepository) DeleteTask(ID string) error {
	args := m.Called(ID)
	return args.Error(0)
}

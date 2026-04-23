package taskservice

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestCreateTask проверяет метод CreateTask сервисного слоя.
// Тесты не обращаются к базе данных — вместо неё используется MockTaskRepository.
//
// Проверяемые сценарии:
//   - нормальное создание задачи
//   - ошибка со стороны базы данных
//   - валидация: пустой текст задачи
func TestCreateTask(t *testing.T) {
	// table-driven тест: каждый случай описывается в виде строки таблицы
	tests := []struct {
		name      string
		input     *Task
		mockSetup func(m *MockTaskRepository, input *Task) // настройка ожиданий мока
		wantErr   bool                                     // ожидаем ошибку или нет
	}{
		{
			name:  "успешное создание задачи",
			input: &Task{Task: "Test", IsDone: false},
			mockSetup: func(m *MockTaskRepository, input *Task) {
				// говорим моку: когда вызовут CreateTask с этим аргументом — вернуть nil (без ошибки)
				m.On("CreateTask", input).Return(nil)
			},
			wantErr: false,
		},
		{
			name:  "ошибка базы данных при создании",
			input: &Task{Task: "Bad task", IsDone: false},
			mockSetup: func(m *MockTaskRepository, input *Task) {
				// имитируем сбой базы данных
				m.On("CreateTask", input).Return(errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name:      "пустой текст задачи",
			input:     &Task{Task: ""},
			mockSetup: func(m *MockTaskRepository, input *Task) {}, // repo не вызывается: сервис отклоняет ещё на валидации
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.mockSetup(mockRepo, tt.input) // регистрируем ожидания до вызова

			service := NewTaskService(mockRepo)
			err := service.CreateTask(tt.input)

			if tt.wantErr {
				assert.Error(t, err) // ожидаем, что ошибка есть
			} else {
				assert.NoError(t, err) // ожидаем, что ошибки нет
			}

			// проверяем, что все зарегистрированные ожидания были вызваны
			mockRepo.AssertExpectations(t)
		})
	}
}

// TestGetAllTasks проверяет метод GetAllTasks сервисного слоя.
//
// Проверяемые сценарии:
//   - успешное получение списка задач
//   - ошибка со стороны базы данных
func TestGetAllTasks(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(m *MockTaskRepository)
		wantErr   bool
	}{
		{
			name: "успешное получение задач",
			mockSetup: func(m *MockTaskRepository) {
				// мок вернёт срез из двух задач
				m.On("GetAllTasks").Return([]Task{{Task: "Task 1"}, {Task: "Task 2"}}, nil)
			},
			wantErr: false,
		},
		{
			name: "ошибка при получении задач",
			mockSetup: func(m *MockTaskRepository) {
				// мок вернёт nil-срез и ошибку
				m.On("GetAllTasks").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.mockSetup(mockRepo)

			service := NewTaskService(mockRepo)
			result, err := service.GetAllTasks()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, 2)               // должно вернуться ровно 2 задачи
				assert.Equal(t, "Task 1", result[0].Task) // проверяем первую задачу
				assert.Equal(t, "Task 2", result[1].Task) // проверяем вторую задачу
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestUpdateTask проверяет метод UpdateTask сервисного слоя.
//
// Логика сервиса при обновлении:
//  1. Валидирует, что текст не пустой
//  2. Загружает существующую задачу из репозитория по ID
//  3. Обновляет поля и сохраняет
//
// Проверяемые сценарии:
//   - успешное обновление
//   - задача не найдена в базе
//   - валидация: пустой текст задачи
func TestUpdateTask(t *testing.T) {
	tests := []struct {
		name      string
		inputTask *Task
		mockSetup func(m *MockTaskRepository, task *Task)
		wantErr   bool
	}{
		{
			name:      "успешное обновление задачи",
			inputTask: &Task{ID: "1", Task: "Updated Task", IsDone: true},
			mockSetup: func(m *MockTaskRepository, task *Task) {
				// сервис сначала загружает старую задачу...
				existing := &Task{ID: "1", Task: "Old Task", IsDone: false}
				m.On("GetTaskByID", task.ID).Return(existing, nil)
				// ...затем изменяет поля existing и передаёт в UpdateTask.
				// mock.Anything используется, потому что existing модифицируется
				// внутри сервиса до вызова — точное значение заранее не известно.
				m.On("UpdateTask", mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "задача не найдена",
			inputTask: &Task{ID: "2", Task: "Another Task"},
			mockSetup: func(m *MockTaskRepository, task *Task) {
				// GetTaskByID вернёт ошибку — UpdateTask вызываться не должен
				m.On("GetTaskByID", task.ID).Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name:      "пустой текст задачи",
			inputTask: &Task{ID: "3", Task: ""},
			mockSetup: func(m *MockTaskRepository, task *Task) {}, // repo не вызывается: сервис отклоняет на валидации
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.mockSetup(mockRepo, tt.inputTask)

			service := NewTaskService(mockRepo)
			err := service.UpdateTask(tt.inputTask)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// TestDeleteTask проверяет метод DeleteTask сервисного слоя.
//
// Логика сервиса: валидирует ID, затем вызывает репозиторий.
//
// Проверяемые сценарии:
//   - успешное удаление
//   - ошибка со стороны базы данных
//   - валидация: пустой ID
func TestDeleteTask(t *testing.T) {
	tests := []struct {
		name      string
		inputID   string
		mockSetup func(m *MockTaskRepository, id string)
		wantErr   bool
	}{
		{
			name:    "успешное удаление задачи",
			inputID: "1",
			mockSetup: func(m *MockTaskRepository, id string) {
				m.On("DeleteTask", id).Return(nil)
			},
			wantErr: false,
		},
		{
			name:    "ошибка при удалении задачи",
			inputID: "2",
			mockSetup: func(m *MockTaskRepository, id string) {
				// имитируем, что задача с таким ID не найдена
				m.On("DeleteTask", id).Return(errors.New("not found"))
			},
			wantErr: true,
		},
		{
			name:      "пустой ID",
			inputID:   "",
			mockSetup: func(m *MockTaskRepository, id string) {}, // repo не вызывается: сервис отклоняет на валидации
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTaskRepository)
			tt.mockSetup(mockRepo, tt.inputID)

			service := NewTaskService(mockRepo)
			err := service.DeleteTask(tt.inputID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

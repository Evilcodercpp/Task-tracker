package taskservice

// TaskService — сервисный слой.
// Находится между хендлерами и репозиторием.
// Здесь должна жить бизнес-логика: валидация, проверки, правила.
// Хендлеры не знают про БД — они работают только через сервис.
type TaskService struct {
	repo TaskRepository
}

// NewTaskService создаёт сервис и принимает репозиторий через аргумент.
func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// CreateTask создаёт новую задачу.
func (s *TaskService) CreateTask(task *Task) error {
	return s.repo.CreateTask(task)
}

// GetAllTasks возвращает все задачи.
func (s *TaskService) GetAllTasks() ([]Task, error) {
	return s.repo.GetAllTasks()
}

// GetTaskByID возвращает одну задачу по ID.
func (s *TaskService) GetTaskByID(id string) (*Task, error) {
	return s.repo.GetTaskByID(id)
}

// UpdateTask обновляет задачу.
func (s *TaskService) UpdateTask(task *Task) error {
	return s.repo.UpdateTask(task)
}

// DeleteTask мягко удаляет задачу по ID.
func (s *TaskService) DeleteTask(id string) error {
	return s.repo.DeleteTask(id)
}

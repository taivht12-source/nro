package ports

import "nro/src/internal/core/domain"

// TaskRepository defines operations for accessing Task data.
type TaskRepository interface {
	// GetTasks loads all task templates.
	GetTasks() ([]*domain.Task, error)
}

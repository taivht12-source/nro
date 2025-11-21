package ports

import "nro/src/internal/core/domain"

// MapRepository định nghĩa các phương thức tương tác với dữ liệu Map.
type MapRepository interface {
	GetTemplate(id int) (*domain.MapTemplate, error)
	GetAllTemplates() ([]*domain.MapTemplate, error)
}

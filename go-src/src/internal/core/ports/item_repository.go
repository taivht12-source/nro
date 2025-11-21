package ports

import "nro/src/internal/core/domain"

type ItemRepository interface {
	GetTemplate(id int) (*domain.ItemTemplate, error)
	GetAllTemplates() ([]*domain.ItemTemplate, error)
}

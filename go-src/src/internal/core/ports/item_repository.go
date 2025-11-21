package ports

import "nro-go/internal/core/domain"

type ItemRepository interface {
	GetTemplate(id int) (*domain.ItemTemplate, error)
	GetAllTemplates() ([]*domain.ItemTemplate, error)
}

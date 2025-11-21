package ports

import "nro/src/internal/core/domain"

// SkillRepository defines methods for accessing skill data.
type SkillRepository interface {
	GetTemplate(id int) (*domain.SkillTemplate, error)
	GetAllTemplates() ([]*domain.SkillTemplate, error)
}

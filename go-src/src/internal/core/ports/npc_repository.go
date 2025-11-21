package ports

import "nro/src/internal/core/domain"

// NPCRepository defines operations for accessing NPC data.
type NPCRepository interface {
	// GetTemplates loads all NPC templates.
	GetTemplates() ([]*domain.NPCTemplate, error)
}

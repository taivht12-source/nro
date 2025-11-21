package services

import (
	"encoding/json"
	"fmt"
	"nro/src/internal/core/domain"
	"os"
)

// BossTemplateLoader loads boss templates from JSON files.
type BossTemplateLoader struct {
	templates map[int]*domain.BossTemplate
}

// NewBossTemplateLoader creates a new template loader.
func NewBossTemplateLoader() *BossTemplateLoader {
	return &BossTemplateLoader{
		templates: make(map[int]*domain.BossTemplate),
	}
}

// LoadFromFile loads boss templates from a JSON file.
func (l *BossTemplateLoader) LoadFromFile(filepath string) error {
	// Read file
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read boss templates file: %w", err)
	}

	// Parse JSON
	var templates []domain.BossTemplate
	if err := json.Unmarshal(data, &templates); err != nil {
		return fmt.Errorf("failed to parse boss templates JSON: %w", err)
	}

	// Store templates
	for i := range templates {
		tmpl := &templates[i]
		l.templates[tmpl.ID] = tmpl
		fmt.Printf("[BOSS LOADER] Loaded template: %s (ID: %d)\n", tmpl.Name, tmpl.ID)
	}

	fmt.Printf("[BOSS LOADER] Successfully loaded %d boss templates\n", len(templates))
	return nil
}

// GetTemplate returns a boss template by ID.
func (l *BossTemplateLoader) GetTemplate(id int) (*domain.BossTemplate, bool) {
	tmpl, ok := l.templates[id]
	return tmpl, ok
}

// GetAllTemplates returns all loaded templates.
func (l *BossTemplateLoader) GetAllTemplates() map[int]*domain.BossTemplate {
	return l.templates
}

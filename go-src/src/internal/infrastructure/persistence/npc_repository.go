package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nro-go/internal/core/domain"
)

// MySQLNPCRepository implements ports.NPCRepository.
type MySQLNPCRepository struct {
	db *sql.DB
}

// NewMySQLNPCRepository creates a new instance.
func NewMySQLNPCRepository(db *sql.DB) *MySQLNPCRepository {
	return &MySQLNPCRepository{db: db}
}

// GetTemplates loads all NPC templates from the database.
func (r *MySQLNPCRepository) GetTemplates() ([]*domain.NPCTemplate, error) {
	query := `SELECT id, name, map_id, x, y, avatar, menu FROM npc_template`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*domain.NPCTemplate
	for rows.Next() {
		var t domain.NPCTemplate
		var menuJSON string
		err := rows.Scan(&t.ID, &t.Name, &t.MapID, &t.X, &t.Y, &t.Avatar, &menuJSON)
		if err != nil {
			return nil, err
		}

		// Parse menu JSON to Dialogue
		// Assuming menu is stored as JSON array of strings: ["Hello", "How are you?"]
		if menuJSON != "" {
			if err := json.Unmarshal([]byte(menuJSON), &t.Dialogue); err != nil {
				fmt.Printf("Error parsing NPC menu JSON for ID %d: %v\n", t.ID, err)
			}
		}

		templates = append(templates, &t)
	}

	return templates, nil
}

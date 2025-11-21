package persistence

import (
	"database/sql"
	"encoding/json"
	"nro/src/internal/core/domain"
	"nro/src/internal/core/ports"
)

// MySQLSkillRepository implements SkillRepository using MySQL.
type MySQLSkillRepository struct {
	db *sql.DB
}

// NewMySQLSkillRepository creates a new MySQL skill repository.
func NewMySQLSkillRepository(db *sql.DB) ports.SkillRepository {
	return &MySQLSkillRepository{db: db}
}

// GetTemplate loads a skill template by ID.
func (r *MySQLSkillRepository) GetTemplate(id int) (*domain.SkillTemplate, error) {
	query := `
		SELECT id, nclass_id, NAME, max_point, mana_use_type, TYPE, icon_id, dam_info, slot, skills
		FROM skill_template WHERE id = ?
	`

	t := &domain.SkillTemplate{}
	err := r.db.QueryRow(query, id).Scan(
		&t.ID, &t.NClassID, &t.Name, &t.MaxPoint, &t.ManaUseType, &t.Type, &t.IconID, &t.DamInfo, &t.Slot, &t.Skills,
	)

	if err != nil {
		return nil, err
	}

	t.SkillData = r.parseSkillData(t.Skills)

	return t, nil
}

// GetAllTemplates loads all skill templates.
func (r *MySQLSkillRepository) GetAllTemplates() ([]*domain.SkillTemplate, error) {
	query := `
		SELECT id, nclass_id, NAME, max_point, mana_use_type, TYPE, icon_id, dam_info, slot, skills
		FROM skill_template
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*domain.SkillTemplate
	for rows.Next() {
		t := &domain.SkillTemplate{}
		err := rows.Scan(
			&t.ID, &t.NClassID, &t.Name, &t.MaxPoint, &t.ManaUseType, &t.Type, &t.IconID, &t.DamInfo, &t.Slot, &t.Skills,
		)
		if err != nil {
			return nil, err
		}
		t.SkillData = r.parseSkillData(t.Skills)
		templates = append(templates, t)
	}

	return templates, nil
}

func (r *MySQLSkillRepository) parseSkillData(data string) []*domain.SkillLevelData {
	// The data is a JSON array of strings, where each string is a JSON object.
	// Example: ["{\"power_require\":1000,...}", "{\"power_require\":10000,...}"]

	var rawList []string
	if err := json.Unmarshal([]byte(data), &rawList); err != nil {
		// Try parsing as direct array of objects if the format is different
		var directList []*domain.SkillLevelData
		if err2 := json.Unmarshal([]byte(data), &directList); err2 == nil {
			return directList
		}
		return nil
	}

	var result []*domain.SkillLevelData
	for _, s := range rawList {
		// Check if s is wrapped in quotes again?
		// The example showed: "[\"{\"power_require\":...}\"]"
		// So rawList contains strings like "{\"power_require\":...}"

		var sd domain.SkillLevelData
		if err := json.Unmarshal([]byte(s), &sd); err != nil {
			// Handle potential double escaping or different format
			// Sometimes it might be just a string that needs unquoting?
			// Let's assume standard JSON object string.
			continue
		}
		result = append(result, &sd)
	}

	return result
}

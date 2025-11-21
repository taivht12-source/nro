package persistence

import (
	"database/sql"
	"nro-go/internal/core/domain"
	"nro-go/internal/core/ports"
)

type MySQLItemRepository struct {
	db *sql.DB
}

func NewMySQLItemRepository(db *sql.DB) ports.ItemRepository {
	return &MySQLItemRepository{db: db}
}

func (r *MySQLItemRepository) GetTemplate(id int) (*domain.ItemTemplate, error) {
	query := `SELECT id, TYPE, gender, NAME, description, level, icon_id, part, is_up_to_up, power_require, gold, gem, head, body, leg 
	          FROM item_template WHERE id = ?`

	row := r.db.QueryRow(query, id)

	var it domain.ItemTemplate
	var isUpToUp int // MySQL boolean is tinyint(1)

	err := row.Scan(
		&it.ID, &it.Type, &it.Gender, &it.Name, &it.Description, &it.Level, &it.IconID, &it.Part,
		&isUpToUp, &it.PowerRequire, &it.Gold, &it.Gem, &it.Head, &it.Body, &it.Leg,
	)
	if err != nil {
		return nil, err
	}

	it.IsUpToUp = isUpToUp == 1

	return &it, nil
}

func (r *MySQLItemRepository) GetAllTemplates() ([]*domain.ItemTemplate, error) {
	query := `SELECT id, TYPE, gender, NAME, description, level, icon_id, part, is_up_to_up, power_require, gold, gem, head, body, leg 
	          FROM item_template`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*domain.ItemTemplate

	for rows.Next() {
		var it domain.ItemTemplate
		var isUpToUp int

		err := rows.Scan(
			&it.ID, &it.Type, &it.Gender, &it.Name, &it.Description, &it.Level, &it.IconID, &it.Part,
			&isUpToUp, &it.PowerRequire, &it.Gold, &it.Gem, &it.Head, &it.Body, &it.Leg,
		)
		if err != nil {
			continue
		}

		it.IsUpToUp = isUpToUp == 1
		templates = append(templates, &it)
	}

	return templates, nil
}

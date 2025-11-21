package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nro/src/internal/core/domain"
)

// MySQLTaskRepository implements ports.TaskRepository.
type MySQLTaskRepository struct {
	db *sql.DB
}

// NewMySQLTaskRepository creates a new instance.
func NewMySQLTaskRepository(db *sql.DB) *MySQLTaskRepository {
	return &MySQLTaskRepository{db: db}
}

// GetTasks loads all task templates from the database.
func (r *MySQLTaskRepository) GetTasks() ([]*domain.Task, error) {
	query := `SELECT id, name, description, npc_id, require_level, objectives, rewards FROM task_template`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*domain.Task
	for rows.Next() {
		var t domain.Task
		var objectivesJSON string
		var rewardsJSON string

		err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.NPCID, &t.RequireLevel, &objectivesJSON, &rewardsJSON)
		if err != nil {
			return nil, err
		}

		// Parse Objectives
		if objectivesJSON != "" {
			if err := json.Unmarshal([]byte(objectivesJSON), &t.Objectives); err != nil {
				fmt.Printf("Error parsing Task Objectives JSON for ID %d: %v\n", t.ID, err)
			}
		}

		// Parse Rewards
		if rewardsJSON != "" {
			if err := json.Unmarshal([]byte(rewardsJSON), &t.Rewards); err != nil {
				fmt.Printf("Error parsing Task Rewards JSON for ID %d: %v\n", t.ID, err)
			}
		}

		tasks = append(tasks, &t)
	}

	return tasks, nil
}

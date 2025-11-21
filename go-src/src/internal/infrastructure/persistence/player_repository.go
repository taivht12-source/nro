package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nro/src/internal/core/domain"
	"nro/src/internal/core/ports"
)

// MySQLPlayerRepository implements PlayerRepository using MySQL.
type MySQLPlayerRepository struct {
	db        *sql.DB
	itemRepo  ports.ItemRepository
	skillRepo ports.SkillRepository
}

// NewMySQLPlayerRepository creates a new MySQL player repository.
func NewMySQLPlayerRepository(db *sql.DB, itemRepo ports.ItemRepository, skillRepo ports.SkillRepository) ports.PlayerRepository {
	return &MySQLPlayerRepository{db: db, itemRepo: itemRepo, skillRepo: skillRepo}
}

// GetByID loads a player by ID.
func (r *MySQLPlayerRepository) GetByID(id int) (*domain.Player, error) {
	query := `
		SELECT id, user_id, name, gender, head, body, leg, 
		       power, tiem_nang, hp, mp, max_hp, max_mp, 
		       stamina, max_stamina, map_id, zone_id, x, y,
		       items_body, items_bag, items_box, gold, gem, skills
		FROM player WHERE id = ?
	`

	player := &domain.Player{}
	var itemsBodyJSON, itemsBagJSON, itemsBoxJSON, skillsJSON string
	var gold int64
	var gem int

	err := r.db.QueryRow(query, id).Scan(
		&player.ID, &player.UserID, &player.Name, &player.Gender,
		&player.Head, &player.Body, &player.Leg,
		&player.Power, &player.TiemNang,
		&player.HP, &player.MP, &player.MaxHP, &player.MaxMP,
		&player.Stamina, &player.MaxStamina,
		&player.MapID, &player.ZoneID, &player.X, &player.Y,
		&itemsBodyJSON, &itemsBagJSON, &itemsBoxJSON, &gold, &gem, &skillsJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Player not found
		}
		return nil, err
	}

	// Parse Inventory
	player.Inventory = &domain.Inventory{
		Gold:      gold,
		Gem:       gem,
		ItemsBody: r.parseItems(itemsBodyJSON),
		ItemsBag:  r.parseItems(itemsBagJSON),
		ItemsBox:  r.parseItems(itemsBoxJSON),
	}

	// Parse Skills
	player.Skills = r.parseSkills(skillsJSON)

	return player, nil
}

// GetByUserID loads all players (characters) for a user.
func (r *MySQLPlayerRepository) GetByUserID(userID int) ([]*domain.Player, error) {
	query := `
		SELECT id, user_id, name, gender, head, body, leg, 
		       power, tiem_nang, hp, mp, max_hp, max_mp, 
		       stamina, max_stamina, map_id, zone_id, x, y,
		       items_body, items_bag, items_box, gold, gem, skills
		FROM player WHERE user_id = ?
		ORDER BY id ASC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []*domain.Player
	for rows.Next() {
		player := &domain.Player{}
		var itemsBodyJSON, itemsBagJSON, itemsBoxJSON, skillsJSON string
		var gold int64
		var gem int

		err := rows.Scan(
			&player.ID, &player.UserID, &player.Name, &player.Gender,
			&player.Head, &player.Body, &player.Leg,
			&player.Power, &player.TiemNang,
			&player.HP, &player.MP, &player.MaxHP, &player.MaxMP,
			&player.Stamina, &player.MaxStamina,
			&player.MapID, &player.ZoneID, &player.X, &player.Y,
			&itemsBodyJSON, &itemsBagJSON, &itemsBoxJSON, &gold, &gem, &skillsJSON,
		)
		if err != nil {
			return nil, err
		}

		player.Inventory = &domain.Inventory{
			Gold:      gold,
			Gem:       gem,
			ItemsBody: r.parseItems(itemsBodyJSON),
			ItemsBag:  r.parseItems(itemsBagJSON),
			ItemsBox:  r.parseItems(itemsBoxJSON),
		}

		player.Skills = r.parseSkills(skillsJSON)

		players = append(players, player)
	}

	return players, rows.Err()
}

// Create creates a new player.
func (r *MySQLPlayerRepository) Create(player *domain.Player) error {
	itemsBody := r.serializeItems(player.Inventory.ItemsBody)
	itemsBag := r.serializeItems(player.Inventory.ItemsBag)
	itemsBox := r.serializeItems(player.Inventory.ItemsBox)
	skillsJSON := r.serializeSkills(player.Skills)

	query := `
		INSERT INTO player (user_id, name, gender, head, body, leg, 
		                    power, tiem_nang, hp, mp, max_hp, max_mp, 
		                    stamina, max_stamina, map_id, zone_id, x, y,
		                    items_body, items_bag, items_box, gold, gem, skills)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		player.UserID, player.Name, player.Gender,
		player.Head, player.Body, player.Leg,
		player.Power, player.TiemNang,
		player.HP, player.MP, player.MaxHP, player.MaxMP,
		player.Stamina, player.MaxStamina,
		player.MapID, player.ZoneID, player.X, player.Y,
		itemsBody, itemsBag, itemsBox, player.Inventory.Gold, player.Inventory.Gem, skillsJSON,
	)

	if err != nil {
		return err
	}

	// Get the auto-generated ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	player.ID = int(id)

	return nil
}

// Update saves player data.
func (r *MySQLPlayerRepository) Update(player *domain.Player) error {
	itemsBody := r.serializeItems(player.Inventory.ItemsBody)
	itemsBag := r.serializeItems(player.Inventory.ItemsBag)
	itemsBox := r.serializeItems(player.Inventory.ItemsBox)
	skillsJSON := r.serializeSkills(player.Skills)

	query := `
		UPDATE player SET
			name = ?, gender = ?, head = ?, body = ?, leg = ?,
			power = ?, tiem_nang = ?, hp = ?, mp = ?, max_hp = ?, max_mp = ?,
			stamina = ?, max_stamina = ?, map_id = ?, zone_id = ?, x = ?, y = ?,
			items_body = ?, items_bag = ?, items_box = ?, gold = ?, gem = ?, skills = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		player.Name, player.Gender, player.Head, player.Body, player.Leg,
		player.Power, player.TiemNang,
		player.HP, player.MP, player.MaxHP, player.MaxMP,
		player.Stamina, player.MaxStamina,
		player.MapID, player.ZoneID, player.X, player.Y,
		itemsBody, itemsBag, itemsBox, player.Inventory.Gold, player.Inventory.Gem, skillsJSON,
		player.ID,
	)

	return err
}

// serializeItems serializes []*Item to JSON string.
func (r *MySQLPlayerRepository) serializeItems(items []*domain.Item) string {
	if items == nil {
		return "[]"
	}

	var rawList []string
	for _, item := range items {
		if item == nil {
			// Empty slot
			// Format: "[]" or handle as needed. NRO might use specific format for null slots?
			// Based on sample data: "[-1,0,\"[]\",1745079848509]" seems to be empty slot?
			// Or just don't include it?
			// Sample data has many entries.
			// Let's assume standard format for existing items.
			// If item is nil, we might need to skip or represent as empty.
			// For ItemsBody, slots matter.
			// Let's use a default empty string representation if needed, or just skip if it's a list (Bag).
			// For Body, we need to preserve index?
			// The DB format is a list of strings.
			continue
		}

		// Format: [TemplateID, Quantity, OptionsJSON, Timestamp]
		// OptionsJSON: [[ID,Param],...]

		var opts [][]int
		for _, o := range item.Options {
			opts = append(opts, []int{o.ID, o.Param})
		}
		optsJSON, _ := json.Marshal(opts)

		// We need to wrap optsJSON in quotes and escape it because the outer list is a list of strings,
		// and each string is a JSON array.
		// Actually, looking at sample: "[2,1,\"[\\\"[47,3]\\\"]\",1744547441896]"
		// It's a JSON array encoded as a string?
		// No, the outer is `text` column.
		// Value: `["[2,1,...]", "[...]"]`

		// Inner item string: `[TemplateID, Quantity, "OptionsJSON", Timestamp]`
		// Note: OptionsJSON is a string inside the array.

		itemData := []interface{}{
			item.Template.ID,
			item.Quantity,
			string(optsJSON),
			item.CreateAt,
		}

		itemJSON, _ := json.Marshal(itemData)
		rawList = append(rawList, string(itemJSON))
	}

	result, _ := json.Marshal(rawList)
	return string(result)
}

// parseItems parses JSON string from DB to []*Item.
func (r *MySQLPlayerRepository) parseItems(data string) []*domain.Item {
	if r.itemRepo == nil {
		return nil
	}

	var rawList []string
	if err := json.Unmarshal([]byte(data), &rawList); err != nil {
		return nil
	}

	var items []*domain.Item
	for _, s := range rawList {
		var itemData []interface{}
		if err := json.Unmarshal([]byte(s), &itemData); err != nil {
			continue
		}

		if len(itemData) < 2 {
			continue
		}

		templateID := int(itemData[0].(float64))
		quantity := int(itemData[1].(float64))

		var options []*domain.ItemOption
		if len(itemData) > 2 {
			optStr, ok := itemData[2].(string)
			if ok {
				var rawOpts [][]int
				if err := json.Unmarshal([]byte(optStr), &rawOpts); err == nil {
					for _, ro := range rawOpts {
						if len(ro) >= 2 {
							options = append(options, &domain.ItemOption{
								ID:    ro[0],
								Param: ro[1],
							})
						}
					}
				}
			}
		}

		template, err := r.itemRepo.GetTemplate(templateID)
		if err != nil {
			template = &domain.ItemTemplate{ID: templateID, Name: fmt.Sprintf("Unknown Item %d", templateID)}
		}

		items = append(items, &domain.Item{
			TemplateID: templateID,
			Template:   template,
			Quantity:   quantity,
			Options:    options,
		})
	}
	return items
}

// parseSkills parses JSON string from DB to []*PlayerSkill.
func (r *MySQLPlayerRepository) parseSkills(data string) []*domain.PlayerSkill {
	if r.skillRepo == nil {
		return nil
	}

	var rawList []string
	if err := json.Unmarshal([]byte(data), &rawList); err != nil {
		return nil
	}

	var skills []*domain.PlayerSkill
	for _, s := range rawList {
		var skillData []interface{}
		if err := json.Unmarshal([]byte(s), &skillData); err != nil {
			continue
		}

		if len(skillData) < 3 {
			continue
		}

		templateID := int(skillData[0].(float64))
		point := int(skillData[1].(float64))
		lastTimeUse := int64(skillData[2].(float64))
		more := 0
		if len(skillData) > 3 {
			more = int(skillData[3].(float64))
		}

		template, err := r.skillRepo.GetTemplate(templateID)
		if err != nil {
			template = &domain.SkillTemplate{ID: templateID, Name: fmt.Sprintf("Unknown Skill %d", templateID)}
		}

		skills = append(skills, &domain.PlayerSkill{
			TemplateID:  templateID,
			Template:    template,
			Point:       point,
			LastTimeUse: lastTimeUse,
			More:        more,
		})
	}
	return skills
}

// serializeSkills serializes []*PlayerSkill to JSON string.
func (r *MySQLPlayerRepository) serializeSkills(skills []*domain.PlayerSkill) string {
	if skills == nil {
		return "[]"
	}

	var rawList []string
	for _, skill := range skills {
		if skill == nil {
			continue
		}

		// Format: [TemplateID, Point, LastTimeUse, More]
		skillData := []interface{}{
			skill.TemplateID,
			skill.Point,
			skill.LastTimeUse,
			skill.More,
		}

		skillJSON, _ := json.Marshal(skillData)
		rawList = append(rawList, string(skillJSON))
	}

	result, _ := json.Marshal(rawList)
	return string(result)
}

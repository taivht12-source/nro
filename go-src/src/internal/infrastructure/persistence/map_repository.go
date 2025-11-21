package persistence

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nro/src/internal/core/domain"
	"nro/src/internal/core/ports"
	"strconv"
	"strings"
)

type MySQLMapRepository struct {
	db *sql.DB
}

func NewMySQLMapRepository(db *sql.DB) ports.MapRepository {
	return &MySQLMapRepository{db: db}
}

func (r *MySQLMapRepository) GetTemplate(id int) (*domain.MapTemplate, error) {
	query := `SELECT id, NAME, type, planet_id, tile_id, bg_id, bg_type, zones, max_player, waypoints, mobs 
	          FROM map_template WHERE id = ?`

	row := r.db.QueryRow(query, id)

	var mt domain.MapTemplate
	var waypointsJSON, mobsJSON string

	err := row.Scan(
		&mt.ID, &mt.Name, &mt.Type, &mt.PlanetID, &mt.TileID, &mt.BgID, &mt.BgType,
		&mt.Zones, &mt.MaxPlayer, &waypointsJSON, &mobsJSON,
	)
	if err != nil {
		return nil, err
	}

	mt.Waypoints = r.parseWaypoints(waypointsJSON)
	mt.Mobs = r.parseMobs(mobsJSON)

	return &mt, nil
}

func (r *MySQLMapRepository) GetAllTemplates() ([]*domain.MapTemplate, error) {
	query := `SELECT id, NAME, type, planet_id, tile_id, bg_id, bg_type, zones, max_player, waypoints, mobs 
	          FROM map_template`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var maps []*domain.MapTemplate

	for rows.Next() {
		var mt domain.MapTemplate
		var waypointsJSON, mobsJSON string

		err := rows.Scan(
			&mt.ID, &mt.Name, &mt.Type, &mt.PlanetID, &mt.TileID, &mt.BgID, &mt.BgType,
			&mt.Zones, &mt.MaxPlayer, &waypointsJSON, &mobsJSON,
		)
		if err != nil {
			continue
		}

		mt.Waypoints = r.parseWaypoints(waypointsJSON)
		mt.Mobs = r.parseMobs(mobsJSON)

		maps = append(maps, &mt)
	}

	return maps, nil
}

// parseWaypoints parses the weird JSON format: ["[\"Name\",MinX,MinY,MaxX,MaxY,IsEnter,IsOffline,GoMapID,GoX,GoY]"]
func (r *MySQLMapRepository) parseWaypoints(data string) []*domain.Waypoint {
	var rawList []string
	if err := json.Unmarshal([]byte(data), &rawList); err != nil {
		fmt.Println("Error parsing waypoints JSON:", err)
		return nil
	}

	var waypoints []*domain.Waypoint
	for _, s := range rawList {
		// Remove brackets [ ]
		s = strings.Trim(s, "[]")
		// Split by comma
		parts := strings.Split(s, ",")
		if len(parts) < 10 {
			continue
		}

		// Parse fields
		// Name is parts[0] (quoted)
		minX, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		minY, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
		maxX, _ := strconv.Atoi(strings.TrimSpace(parts[3]))
		maxY, _ := strconv.Atoi(strings.TrimSpace(parts[4]))
		goMapID, _ := strconv.Atoi(strings.TrimSpace(parts[7]))
		goX, _ := strconv.Atoi(strings.TrimSpace(parts[8]))
		goY, _ := strconv.Atoi(strings.TrimSpace(parts[9]))

		wp := &domain.Waypoint{
			MinX:    int16(minX),
			MinY:    int16(minY),
			MaxX:    int16(maxX),
			MaxY:    int16(maxY),
			GoMapID: goMapID,
			GoX:     int16(goX),
			GoY:     int16(goY),
		}
		waypoints = append(waypoints, wp)
	}
	return waypoints
}

// parseMobs parses the weird JSON format: ["[TemplateID,Level,Hp,X,Y]"]
func (r *MySQLMapRepository) parseMobs(data string) []*domain.Mob {
	var rawList []string
	if err := json.Unmarshal([]byte(data), &rawList); err != nil {
		return nil
	}

	var mobs []*domain.Mob
	for _, s := range rawList {
		s = strings.Trim(s, "[]")
		parts := strings.Split(s, ",")
		if len(parts) < 5 {
			continue
		}

		templateID, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		level, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		hp, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
		x, _ := strconv.Atoi(strings.TrimSpace(parts[3]))
		y, _ := strconv.Atoi(strings.TrimSpace(parts[4]))

		mob := &domain.Mob{
			Template: &domain.MobTemplate{
				ID:    templateID,
				Level: level,
				Hp:    hp,
				// Name will be loaded later from MobTemplate repository
			},
			X:     int16(x),
			Y:     int16(y),
			Hp:    hp,
			MaxHp: hp,
		}
		mobs = append(mobs, mob)
	}
	return mobs
}

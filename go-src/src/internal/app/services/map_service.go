package services

import (
	"fmt"
	"nro/src/internal/core/domain"
	"nro/src/internal/core/ports"
	"sync"
)

// MapService quản lý map templates và waypoints.
//
// EXPLANATION:
// MapService chịu trách nhiệm:
// 1. Load và quản lý thông tin bản đồ (MapTemplate) và các khu vực (Zone).
// 2. Cung cấp thông tin về Waypoints (điểm dịch chuyển) để kết nối các map.
// 3. Xử lý logic tải dữ liệu map khi server khởi động.
type MapService struct {
	maps    map[int]*domain.MapTemplate
	mapRepo ports.MapRepository
	mu      sync.RWMutex
}

var mapServiceInstance *MapService
var mapOnce sync.Once

// GetMapService returns singleton instance.
func GetMapService() *MapService {
	mapOnce.Do(func() {
		mapServiceInstance = &MapService{
			maps: make(map[int]*domain.MapTemplate),
		}
		// mapServiceInstance.loadMockMaps() // Disable mock maps by default
	})
	return mapServiceInstance
}

// SetRepository injects the MapRepository.
func (ms *MapService) SetRepository(repo ports.MapRepository) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.mapRepo = repo
	ms.loadMapsFromRepo()
}

// loadMapsFromRepo loads all maps from the repository into memory.
func (ms *MapService) loadMapsFromRepo() {
	if ms.mapRepo == nil {
		return
	}

	maps, err := ms.mapRepo.GetAllTemplates()
	if err != nil {
		fmt.Printf("[MAP] Failed to load maps from DB: %v\n", err)
		return
	}

	for _, m := range maps {
		ms.maps[m.ID] = m
	}
	fmt.Printf("[MAP] Loaded %d maps from DB\n", len(ms.maps))
}

// loadMockMaps loads default maps into memory.
func (ms *MapService) loadMockMaps() {
	mockMaps := []*domain.MapTemplate{
		// Map 0: Làng Aru (Starting village)
		{
			ID:        0,
			Name:      "Làng Aru",
			Type:      0,
			PlanetID:  0, // Trái Đất
			TileID:    0,
			BgID:      0,
			BgType:    0,
			Zones:     3,
			MaxPlayer: 30,
			Waypoints: []*domain.Waypoint{
				{MinX: 400, MinY: 300, MaxX: 450, MaxY: 350, GoMapID: 1, GoX: 100, GoY: 100},
				{MinX: 50, MinY: 50, MaxX: 100, MaxY: 100, GoMapID: 5, GoX: 200, GoY: 200},
			},
		},
		// Map 1: Rừng Khỉ
		{
			ID:        1,
			Name:      "Rừng Khỉ",
			Type:      1,
			PlanetID:  0,
			TileID:    1,
			BgID:      1,
			BgType:    0,
			Zones:     2,
			MaxPlayer: 20,
			Waypoints: []*domain.Waypoint{
				{MinX: 50, MinY: 50, MaxX: 100, MaxY: 100, GoMapID: 0, GoX: 425, GoY: 325},
			},
		},
		// Map 5: Nhà Kame
		{
			ID:        5,
			Name:      "Nhà Kame",
			Type:      2,
			PlanetID:  0,
			TileID:    5,
			BgID:      5,
			BgType:    1,
			Zones:     1,
			MaxPlayer: 10,
			Waypoints: []*domain.Waypoint{
				{MinX: 100, MinY: 100, MaxX: 150, MaxY: 150, GoMapID: 0, GoX: 75, GoY: 75},
			},
		},
		// Map 2: Đại Hội Võ Thuật
		{
			ID:        2,
			Name:      "Đại Hội Võ Thuật",
			Type:      3,
			PlanetID:  0,
			TileID:    2,
			BgID:      2,
			BgType:    0,
			Zones:     1,
			MaxPlayer: 50,
			Waypoints: []*domain.Waypoint{
				{MinX: 200, MinY: 200, MaxX: 250, MaxY: 250, GoMapID: 0, GoX: 300, GoY: 300},
			},
		},
	}

	for _, m := range mockMaps {
		ms.maps[m.ID] = m
	}

	fmt.Printf("[MAP] Loaded %d mock maps\n", len(ms.maps))
}

// GetMap returns map template by ID.
func (ms *MapService) GetMap(mapID int) *domain.MapTemplate {
	ms.mu.RLock()
	// Try cache first
	if m, ok := ms.maps[mapID]; ok {
		ms.mu.RUnlock()
		return m
	}
	ms.mu.RUnlock()

	// If not in cache and repo exists, try to load it
	if ms.mapRepo != nil {
		ms.mu.Lock()
		defer ms.mu.Unlock()
		// Double check
		if m, ok := ms.maps[mapID]; ok {
			return m
		}

		m, err := ms.mapRepo.GetTemplate(mapID)
		if err == nil && m != nil {
			ms.maps[mapID] = m
			return m
		}
	}

	return nil
}

// GetAllMaps returns all available maps.
func (ms *MapService) GetAllMaps() []*domain.MapTemplate {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	maps := make([]*domain.MapTemplate, 0, len(ms.maps))
	for _, m := range ms.maps {
		maps = append(maps, m)
	}
	return maps
}

// CheckWaypoint kiểm tra xem vị trí có phải là waypoint không.
func (ms *MapService) CheckWaypoint(mapID int, x, y int16) *domain.Waypoint {
	mapTemplate := ms.GetMap(mapID)
	if mapTemplate == nil {
		return nil
	}

	for _, wp := range mapTemplate.Waypoints {
		if x >= wp.MinX && x <= wp.MaxX && y >= wp.MinY && y <= wp.MaxY {
			fmt.Printf("[MAP] Waypoint detected at (%d,%d) on map %d -> Go to map %d\n", x, y, mapID, wp.GoMapID)
			return wp
		}
	}

	return nil
}

// GetMapName returns map name by ID.
func (ms *MapService) GetMapName(mapID int) string {
	mapTemplate := ms.GetMap(mapID)
	if mapTemplate == nil {
		return "Unknown Map"
	}
	return mapTemplate.Name
}

// GetZoneCount returns number of zones in a map.
func (ms *MapService) GetZoneCount(mapID int) int {
	mapTemplate := ms.GetMap(mapID)
	if mapTemplate == nil {
		return 1 // Default 1 zone
	}
	return mapTemplate.Zones
}

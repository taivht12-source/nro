package services

import (
	"fmt"
	"nro/src/internal/core/domain"
	"nro/src/internal/infrastructure/session"
	"nro/src/pkg/protocol"
	"sync"
)

// ZoneService quản lý các khu vực (Zone) trong Map.
//
// EXPLANATION:
// ZoneService chịu trách nhiệm:
// 1. Quản lý trạng thái của từng khu vực (Zone): Danh sách người chơi, quái vật (Mobs), và vật phẩm rơi (Items).
// 2. Xử lý logic Broadcast: Gửi gói tin cho tất cả người chơi trong cùng một khu vực (ví dụ: Chat, Di chuyển, Skill).
// 3. Đồng bộ hóa dữ liệu giữa các người chơi trong Zone.
type ZoneService struct {
	zones map[int]map[int]*domain.Zone // Map[MapID][ZoneID] -> Zone
	mu    sync.RWMutex
}

var zoneServiceInstance *ZoneService
var zoneOnce sync.Once

func GetZoneService() *ZoneService {
	zoneOnce.Do(func() {
		zoneServiceInstance = &ZoneService{
			zones: make(map[int]map[int]*domain.Zone),
		}
	})
	return zoneServiceInstance
}

// GetZone lấy thông tin Zone. Nếu chưa có sẽ khởi tạo.
func (zs *ZoneService) GetZone(mapID, zoneID int) *domain.Zone {
	zs.mu.Lock()
	defer zs.mu.Unlock()

	if _, ok := zs.zones[mapID]; !ok {
		zs.zones[mapID] = make(map[int]*domain.Zone)
	}

	if _, ok := zs.zones[mapID][zoneID]; !ok {
		zs.zones[mapID][zoneID] = &domain.Zone{
			ID:      zoneID,
			MapID:   mapID,
			ZoneID:  zoneID,
			Players: make(map[int]*domain.Player),
		}
	}

	return zs.zones[mapID][zoneID]
}

// EnterZone xử lý người chơi vào khu vực.
func (zs *ZoneService) EnterZone(player *domain.Player, mapID, zoneID int) {
	zone := zs.GetZone(mapID, zoneID)

	// Cập nhật thông tin player
	player.MapID = mapID
	player.ZoneID = zoneID

	// Thêm vào danh sách zone
	zone.Players[player.ID] = player

	fmt.Printf("[ZONE] Player %s entered Map %d, Zone %d\n", player.Name, mapID, zoneID)
}

// LeaveZone xử lý người chơi rời khu vực.
func (zs *ZoneService) LeaveZone(player *domain.Player) {
	zone := zs.GetZone(player.MapID, player.ZoneID)
	delete(zone.Players, player.ID)

	fmt.Printf("[ZONE] Player %s left Map %d, Zone %d\n", player.Name, player.MapID, player.ZoneID)
}

// GetPlayersInZone returns all players in a specific zone.
func (zs *ZoneService) GetPlayersInZone(mapID, zoneID int) []*domain.Player {
	zone := zs.GetZone(mapID, zoneID)

	players := make([]*domain.Player, 0, len(zone.Players))
	for _, p := range zone.Players {
		players = append(players, p)
	}

	return players
}

// GetPlayerCount returns number of players in a zone.
func (zs *ZoneService) GetPlayerCount(mapID, zoneID int) int {
	zone := zs.GetZone(mapID, zoneID)
	return len(zone.Players)
}

// TeleportPlayer moves player to a different map/zone.
func (zs *ZoneService) TeleportPlayer(player *domain.Player, newMapID, newZoneID int, newX, newY int16) {
	// Leave current zone
	zs.LeaveZone(player)

	// Update position
	player.X = newX
	player.Y = newY

	// Enter new zone
	zs.EnterZone(player, newMapID, newZoneID)

	fmt.Printf("[ZONE] Player %s teleported to Map %d, Zone %d at (%d, %d)\n",
		player.Name, newMapID, newZoneID, newX, newY)
}

// CheckWaypoint checks if the player is standing on a waypoint.
func (zs *ZoneService) CheckWaypoint(player *domain.Player) *domain.Waypoint {
	// Use MapService to get map template
	mapService := GetMapService()
	mapTemplate := mapService.GetMap(player.MapID)

	if mapTemplate == nil {
		return nil
	}

	for _, wp := range mapTemplate.Waypoints {
		if player.X >= wp.MinX && player.X <= wp.MaxX &&
			player.Y >= wp.MinY && player.Y <= wp.MaxY {
			return wp
		}
	}

	return nil
}

// ChangeMap handles the logic of moving a player to a new map via waypoint.
func (zs *ZoneService) ChangeMap(player *domain.Player, wp *domain.Waypoint) {
	// 1. Leave current zone
	zs.LeaveZone(player)

	// 2. Update Player Position
	player.MapID = wp.GoMapID
	player.X = wp.GoX
	player.Y = wp.GoY
	// ZoneID usually resets to 0 or handled by logic (e.g. find less crowded zone)
	// For now, default to Zone 0
	player.ZoneID = 0

	// 3. Enter new zone
	zs.EnterZone(player, player.MapID, player.ZoneID)

	// 4. Send Map Info to Player (Service_MapTransporter equivalent)
	// TODO: Implement sending Map Info packet
	fmt.Printf("[WAYPOINT] Player %s changed map to %d at (%d, %d)\n", player.Name, player.MapID, player.X, player.Y)
}

// BroadcastMove gửi gói tin di chuyển cho người chơi khác trong cùng Zone.
func (zs *ZoneService) BroadcastMove(player *domain.Player, msg *protocol.Message) {
	zone := zs.GetZone(player.MapID, player.ZoneID)

	var receiverIDs []int
	for pid := range zone.Players {
		if pid != player.ID { // Không gửi lại cho chính mình
			receiverIDs = append(receiverIDs, pid)
		}
	}

	// Dùng SessionManager để gửi
	session.GetManager().Broadcast(receiverIDs, msg)
}

// BroadcastToZone sends a message to all players in a zone.
func (zs *ZoneService) BroadcastToZone(mapID, zoneID int, msg *protocol.Message) {
	zone := zs.GetZone(mapID, zoneID)

	var receiverIDs []int
	for pid := range zone.Players {
		receiverIDs = append(receiverIDs, pid)
	}

	session.GetManager().Broadcast(receiverIDs, msg)
}

// BroadcastToZoneExcept sends a message to all players except one.
func (zs *ZoneService) BroadcastToZoneExcept(mapID, zoneID int, exceptPlayerID int, msg *protocol.Message) {
	zone := zs.GetZone(mapID, zoneID)

	var receiverIDs []int
	for pid := range zone.Players {
		if pid != exceptPlayerID {
			receiverIDs = append(receiverIDs, pid)
		}
	}

	session.GetManager().Broadcast(receiverIDs, msg)
}

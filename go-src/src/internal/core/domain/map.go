package domain

// MapTemplate đại diện cho mẫu bản đồ (dữ liệu tĩnh).
type MapTemplate struct {
	ID           int
	Name         string
	Type         int8
	PlanetID     int8
	TileID       int8
	BgID         int8
	BgType       int8
	IsMapOffline bool
	Zones        int
	MaxPlayer    int
	Waypoints    []*Waypoint
	Mobs         []*Mob // Danh sách quái trong map (Spawn default)
}

// Zone đại diện cho một khu vực cụ thể trong Map (dữ liệu động).
type Zone struct {
	ID      int
	MapID   int
	ZoneID  int             // ID của khu vực (0 đến n)
	Players map[int]*Player // Danh sách người chơi trong khu vực
	Mobs    []*Mob          // Danh sách quái trong khu vực
	Items   []*ItemMap      // Vật phẩm rơi trên đất
}

// Waypoint điểm chuyển map.
type Waypoint struct {
	MinX, MinY int16
	MaxX, MaxY int16
	GoMapID    int
	GoX, GoY   int16
}

// MobTemplate mẫu quái.
type MobTemplate struct {
	ID    int
	Name  string
	Hp    int
	Level int
}

// Mob quái vật thực thể.
type Mob struct {
	ID       int
	Template *MobTemplate
	X, Y     int16
	Hp       int
	MaxHp    int
	Status   int // 0: Chết, 1: Sống, ...
}

// ItemMap vật phẩm rơi.
type ItemMap struct {
	ID       int
	PlayerID int // Người sở hữu
	Template *ItemTemplate
	X, Y     int16
}

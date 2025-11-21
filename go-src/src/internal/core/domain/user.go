package domain

import "time"

// User đại diện cho tài khoản người dùng (Account).
type User struct {
	ID       int
	Username string
	Password string
	Role     int // 0: User, 1: Admin...
	Ban      int // 1: Banned
	Active   bool
	CreateAt time.Time
	UpdateAt time.Time
}

// Player đại diện cho nhân vật trong game.
type Player struct {
	ID           int
	UserID       int
	Name         string
	Gender       int8 // 0: Trai Dat, 1: Namec, 2: Xayda
	Head         int16
	Body         int16
	Leg          int16
	Power        int64 // Sức mạnh
	TiemNang     int64 // Tiềm năng
	HP           int
	MP           int
	MaxHP        int
	MaxMP        int
	Stamina      int // Thể lực
	MaxStamina   int
	MapID        int
	ZoneID       int
	X            int16
	Y            int16
	LastMoveTime time.Time // Thời gian di chuyển cuối cùng (cho anti-cheat)
	Inventory    *Inventory
	Skills       []*PlayerSkill // Danh sách kỹ năng
	Effects      []*Effect      // Danh sách hiệu ứng (Buff/Debuff)
	Tasks        []*PlayerTask  // Danh sách nhiệm vụ
	Level        int            // Cấp độ
	Exp          int64          // Kinh nghiệm
	ClanID       int            // ID bang hội
	// ... còn nhiều thuộc tính khác sẽ thêm dần
}

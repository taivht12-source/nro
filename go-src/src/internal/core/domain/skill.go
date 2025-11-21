package domain

// SkillTemplate represents a skill template from the database.
type SkillTemplate struct {
	ID          int
	NClassID    int // Class ID (0: Trai Dat, 1: Namec, 2: Xayda)
	Name        string
	MaxPoint    int
	ManaUseType int
	Type        int
	IconID      int
	DamInfo     string
	Slot        int
	Skills      string            // JSON string detailing skill levels
	SkillData   []*SkillLevelData // Parsed skill data
}

// SkillLevelData represents stats for a specific skill level.
type SkillLevelData struct {
	PowerRequire int64  `json:"power_require"`
	Damage       int    `json:"damage"`
	Dx           int    `json:"dx"`
	Dy           int    `json:"dy"`
	Price        int    `json:"price"`
	MaxFight     int    `json:"max_fight"`
	ManaUse      int    `json:"mana_use"`
	CoolDown     int    `json:"cool_down"`
	ID           int    `json:"id"`
	Point        int    `json:"point"`
	Info         string `json:"info"`
}

// PlayerSkill represents a skill learned by a player.
type PlayerSkill struct {
	TemplateID  int
	Template    *SkillTemplate
	Point       int   // Current level of the skill từ này sẽ mapping SkillData để lấy thông tin skilldata
	LastTimeUse int64 // Timestamp of last use
	More        int   // Extra data (e.g., remaining active time)
}

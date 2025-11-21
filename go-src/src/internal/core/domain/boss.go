package domain

// BossStatus defines the state of a boss.
type BossStatus int

const (
	BossStatusRest     BossStatus = 0
	BossStatusRespawn  BossStatus = 1
	BossStatusJoinMap  BossStatus = 2
	BossStatusChatS    BossStatus = 3
	BossStatusActive   BossStatus = 4
	BossStatusDie      BossStatus = 5
	BossStatusChatE    BossStatus = 6
	BossStatusLeaveMap BossStatus = 7
	BossStatusAfk      BossStatus = 8
)

// BossTemplate defines the static data for a boss (matching BossData.java).
type BossTemplate struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Gender      int8     `json:"gender"`
	Outfit      []int16  `json:"outfit"` // Head, Body, Leg, Flag, Aura, Eff
	Damage      int      `json:"damage"`
	HP          []int64  `json:"hp"`       // HP scaling options
	MapJoin     []int    `json:"map_join"` // Maps where boss can appear
	SkillTemp   [][]int  `json:"skills"`   // [SkillID, Level, Cooldown]
	TextS       []string `json:"text_s"`   // Start chat
	TextM       []string `json:"text_m"`   // Mid-fight chat
	TextE       []string `json:"text_e"`   // End chat
	SecondsRest int      `json:"seconds_rest"`
	TypeAppear  int      `json:"type_appear"`          // 0: Default, 1: With Another, 2: Another Level, 3: Call by Another
	AIType      string   `json:"ai_type,omitempty"`    // AI type identifier
	AndroidID   int      `json:"android_id,omitempty"` // For Android bosses
	IsZamasu    bool     `json:"is_zamasu,omitempty"`  // For Black Goku/Zamasu
}

// Boss represents a runtime instance of a boss.
type Boss struct {
	ID         int
	TemplateID int
	Name       string
	Gender     int8
	HP         int64
	MaxHP      int64
	MP         int64
	MaxMP      int64
	Damage     int
	Def        int
	Crit       int

	MapID  int
	ZoneID int
	X      int16
	Y      int16

	Status   BossStatus
	TargetID int // Player ID being attacked

	LastTimeRest   int64
	LastTimeChatS  int64
	LastTimeChatM  int64
	LastTimeChatE  int64
	LastTimeAttack int64
	LastTimeTarget int64

	TimeChatS int
	TimeChatM int
	TimeChatE int

	IndexChatS int
	IndexChatE int

	CurrentLevel int // For multi-stage bosses

	Template    *BossTemplate
	PlayerSkill *PlayerSkill // Reuse PlayerSkill struct if possible, or define BossSkill
	AI          BossAI       // Custom AI behavior
}

// IsDead checks if the boss is dead.
func (b *Boss) IsDead() bool {
	return b.HP <= 0
}

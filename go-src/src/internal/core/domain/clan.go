package domain

import "time"

// ClanMember represents a member of a clan.
type ClanMember struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Role     int8   `json:"role"` // 0: Member, 1: Deputy, 2: Leader
	Power    int64  `json:"power"`
	JoinTime int64  `json:"join_time"`
}

// Clan represents a guild/clan.
type Clan struct {
	ID         int            `json:"id"`
	Name       string         `json:"name"`
	Slogan     string         `json:"slogan"`
	LeaderID   int            `json:"leader_id"`
	Members    []*ClanMember  `json:"members"`
	MaxMembers int            `json:"max_members"`
	Level      int            `json:"level"`
	Exp        int64          `json:"exp"`
	CreateTime int64          `json:"create_time"`
	Messages   []*ClanMessage `json:"messages"` // Chat history
}

// ClanMessage represents a chat message in the clan.
type ClanMessage struct {
	SenderID   int    `json:"sender_id"`
	SenderName string `json:"sender_name"`
	Content    string `json:"content"`
	Time       int64  `json:"time"`
}

// NewClan creates a new clan.
func NewClan(id int, name string, leader *Player) *Clan {
	clan := &Clan{
		ID:         id,
		Name:       name,
		Slogan:     "Welcome to our clan!",
		LeaderID:   leader.ID,
		Members:    make([]*ClanMember, 0),
		MaxMembers: 10,
		Level:      1,
		Exp:        0,
		CreateTime: time.Now().Unix(),
		Messages:   make([]*ClanMessage, 0),
	}

	// Add leader as member
	clan.AddMember(leader, 2) // 2 = Leader
	return clan
}

// AddMember adds a player to the clan.
func (c *Clan) AddMember(player *Player, role int8) {
	member := &ClanMember{
		ID:       player.ID,
		Name:     player.Name,
		Role:     role,
		Power:    player.Power,
		JoinTime: time.Now().Unix(),
	}
	c.Members = append(c.Members, member)
}

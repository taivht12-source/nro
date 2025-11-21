package services

import (
	"fmt"
	"nro-go/internal/core/domain"
	"sync"
	"time"
)

// ClanService manages clan operations.
type ClanService struct {
	clans  map[int]*domain.Clan
	mu     sync.RWMutex
	nextID int
}

var clanServiceInstance *ClanService
var clanOnce sync.Once

// GetClanService returns singleton instance.
func GetClanService() *ClanService {
	clanOnce.Do(func() {
		clanServiceInstance = &ClanService{
			clans:  make(map[int]*domain.Clan),
			nextID: 1,
		}
	})
	return clanServiceInstance
}

// CreateClan creates a new clan.
func (s *ClanService) CreateClan(name string, leader *domain.Player) (*domain.Clan, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if name exists
	for _, c := range s.clans {
		if c.Name == name {
			return nil, fmt.Errorf("clan name already exists")
		}
	}

	// Check if leader is already in a clan
	if leader.ClanID != 0 { // Assuming Player has ClanID field, if not we need to add it
		return nil, fmt.Errorf("you are already in a clan")
	}

	clan := domain.NewClan(s.nextID, name, leader)
	s.clans[clan.ID] = clan
	s.nextID++

	// Update leader's clan ID
	leader.ClanID = clan.ID

	fmt.Printf("[CLAN] Clan %s created by %s\n", name, leader.Name)
	return clan, nil
}

// GetClan returns a clan by ID.
func (s *ClanService) GetClan(id int) *domain.Clan {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.clans[id]
}

// AddMember adds a player to a clan.
func (s *ClanService) AddMember(clanID int, player *domain.Player) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	clan, ok := s.clans[clanID]
	if !ok {
		return fmt.Errorf("clan not found")
	}

	if len(clan.Members) >= clan.MaxMembers {
		return fmt.Errorf("clan is full")
	}

	if player.ClanID != 0 {
		return fmt.Errorf("player is already in a clan")
	}

	clan.AddMember(player, 0) // 0 = Member
	player.ClanID = clan.ID

	fmt.Printf("[CLAN] %s joined clan %s\n", player.Name, clan.Name)
	return nil
}

// SendClanMessage sends a chat message to the clan.
func (s *ClanService) SendClanMessage(clanID int, sender *domain.Player, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	clan, ok := s.clans[clanID]
	if !ok {
		return fmt.Errorf("clan not found")
	}

	msg := &domain.ClanMessage{
		SenderID:   sender.ID,
		SenderName: sender.Name,
		Content:    content,
		Time:       time.Now().Unix(),
	}

	clan.Messages = append(clan.Messages, msg)
	// Keep only last 50 messages
	if len(clan.Messages) > 50 {
		clan.Messages = clan.Messages[1:]
	}

	fmt.Printf("[CLAN] [%s] %s: %s\n", clan.Name, sender.Name, content)
	return nil
}

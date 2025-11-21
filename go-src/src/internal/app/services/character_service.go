package services

import (
	"errors"
	"fmt"
	"nro/src/internal/core/domain"
	"nro/src/internal/core/ports"
)

// CharacterService handles character management logic.
//
// EXPLANATION:
// CharacterService chịu trách nhiệm:
// 1. Quản lý danh sách nhân vật của người dùng (User).
// 2. Xử lý logic tạo nhân vật mới (bao gồm khởi tạo chỉ số, vật phẩm, kỹ năng ban đầu).
// 3. Xử lý logic chọn nhân vật để vào game (Login).
type CharacterService struct {
	playerRepo ports.PlayerRepository
}

// NewCharacterService creates a new character service.
func NewCharacterService(playerRepo ports.PlayerRepository) *CharacterService {
	return &CharacterService{playerRepo: playerRepo}
}

// GetCharacterList returns all characters for a user.
func (s *CharacterService) GetCharacterList(userID int) ([]*domain.Player, error) {
	if s.playerRepo == nil {
		return nil, errors.New("player repository not available")
	}

	players, err := s.playerRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get characters: %w", err)
	}

	return players, nil
}

// SelectCharacter validates and loads a character for a user.
func (s *CharacterService) SelectCharacter(userID int, charID int) (*domain.Player, error) {
	if s.playerRepo == nil {
		return nil, errors.New("player repository not available")
	}

	// Load the character
	player, err := s.playerRepo.GetByID(charID)
	if err != nil {
		return nil, fmt.Errorf("failed to load character: %w", err)
	}

	if player == nil {
		return nil, errors.New("character not found")
	}

	// Verify the character belongs to this user
	if player.UserID != userID {
		return nil, errors.New("character does not belong to this user")
	}

	fmt.Printf("[CHAR] User %d selected character %d (%s)\n", userID, charID, player.Name)
	return player, nil
}

// CreateCharacter creates a new character for a user.
func (s *CharacterService) CreateCharacter(userID int, name string, gender int8) (*domain.Player, error) {
	if s.playerRepo == nil {
		return nil, errors.New("player repository not available")
	}

	// Create new player with default stats
	player := &domain.Player{
		UserID:     userID,
		Name:       name,
		Gender:     gender,
		Head:       0,
		Body:       0,
		Leg:        0,
		Power:      0,
		TiemNang:   0,
		HP:         100,
		MP:         100,
		MaxHP:      100,
		MaxMP:      100,
		Stamina:    100,
		MaxStamina: 100,
		MapID:      0,   // Starting map
		ZoneID:     0,   // Starting zone
		X:          100, // Starting position
		Y:          100,
		Inventory: &domain.Inventory{
			ItemsBag:  make([]*domain.Item, 0),
			ItemsBody: make([]*domain.Item, 10), // 10 slots for equipment
			ItemsBox:  make([]*domain.Item, 0),
			Gold:      0,
			Gem:       0,
		},
		Skills: make([]*domain.PlayerSkill, 0),
		Tasks:  make([]*domain.PlayerTask, 0),
		Level:  1,
		Exp:    0,
	}

	// Give starter items
	itemService := GetItemService()
	invService := GetInventoryService()

	// Starter weapon
	starterWeapon := itemService.CreateItem(1, 1) // Gậy Như Ý
	if starterWeapon != nil {
		invService.AddItem(player, starterWeapon)
	}

	// Starter armor based on gender
	armorID := 10 + int(gender) // 10: Áo Kame, 11: Áo Namek, 12: Áo Saiyan
	starterArmor := itemService.CreateItem(armorID, 1)
	if starterArmor != nil {
		invService.AddItem(player, starterArmor)
	}

	// Starter consumables
	dauThan := itemService.CreateItem(30, 5) // 5x Đậu Thần
	if dauThan != nil {
		invService.AddItem(player, dauThan)
	}

	// Give starter skills based on race
	skillService := GetSkillService()
	starterSkills := skillService.GetSkillsByRace(gender)
	for _, skillID := range starterSkills {
		skillService.LearnSkill(player, skillID)
	}

	err := s.playerRepo.Create(player)
	if err != nil {
		return nil, fmt.Errorf("failed to create character: %w", err)
	}

	fmt.Printf("[CHAR] Created new character %d (%s) for user %d with starter items\n", player.ID, player.Name, userID)
	return player, nil
}

package services

import (
	"fmt"
	"nro/src/internal/core/domain"
	"nro/src/internal/core/ports"
	"sync"
)

// SkillService manages skill templates and player skills.
//
// EXPLANATION:
// SkillService chịu trách nhiệm:
// 1. Load và quản lý thông tin kỹ năng (SkillTemplate) từ Database.
// 2. Xử lý logic học kỹ năng mới và nâng cấp kỹ năng cho nhân vật.
// 3. Cung cấp thông tin chi tiết về kỹ năng (sát thương, mana, cooldown) cho CombatService.
type SkillService struct {
	templates map[int]*domain.SkillTemplate
	repo      ports.SkillRepository
	mu        sync.RWMutex
}

var skillServiceInstance *SkillService
var skillOnce sync.Once

// GetSkillService returns singleton instance.
func GetSkillService() *SkillService {
	skillOnce.Do(func() {
		skillServiceInstance = &SkillService{
			templates: make(map[int]*domain.SkillTemplate),
		}
	})
	return skillServiceInstance
}

func (s *SkillService) SetRepository(repo ports.SkillRepository) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.repo = repo
}

// {ID: 1, Name: "Kamehameha", Type: 0, Damage: 100, MPCost: 10, Cooldown: 5, Range: 200, Description: "Sóng năng lượng mạnh mẽ"},
// {ID: 2, Name: "Thái Dương Hạ San", Type: 0, Damage: 150, MPCost: 20, Cooldown: 10, Range: 150, Description: "Đấm mạnh từ trên cao"},
// {ID: 3, Name: "Masenko", Type: 0, Damage: 80, MPCost: 8, Cooldown: 3, Range: 180, Description: "Năng lượng vàng"},
// {ID: 4, Name: "Kaioken", Type: 1, Damage: 0, MPCost: 30, Cooldown: 60, Range: 0, Description: "Tăng sức mạnh tạm thời"},

// // Namek (Race 1)
// {ID: 10, Name: "Đấm Liên Hoàn", Type: 0, Damage: 90, MPCost: 12, Cooldown: 4, Range: 100, Description: "Đấm nhanh liên tục"},
// {ID: 11, Name: "Makankosappo", Type: 0, Damage: 200, MPCost: 30, Cooldown: 15, Range: 250, Description: "Súng xoắn ốc ma quỷ"},
// {ID: 12, Name: "Tái Tạo", Type: 2, Damage: 0, MPCost: 50, Cooldown: 30, Range: 0, Description: "Hồi phục HP"},
// {ID: 13, Name: "Đấm Từ Trên Trời", Type: 0, Damage: 120, MPCost: 15, Cooldown: 7, Range: 120, Description: "Đấm từ xa"},

// // Xayda (Race 2)
// {ID: 20, Name: "Galick Gun", Type: 0, Damage: 120, MPCost: 15, Cooldown: 6, Range: 200, Description: "Tia năng lượng tím"},
// {ID: 21, Name: "Final Flash", Type: 0, Damage: 250, MPCost: 40, Cooldown: 20, Range: 300, Description: "Năng lượng cực mạnh"},
// {ID: 22, Name: "Big Bang Attack", Type: 0, Damage: 180, MPCost: 25, Cooldown: 12, Range: 180, Description: "Quả cầu năng lượng"},
// {ID: 23, Name: "Tự Hủy", Type: 0, Damage: 500, MPCost: 100, Cooldown: 120, Range: 150, Description: "Nổ tung (nguy hiểm)"},

// GetTemplate returns skill template by ID.
func (s *SkillService) GetTemplate(id int) *domain.SkillTemplate {
	s.mu.RLock()
	template, ok := s.templates[id]
	s.mu.RUnlock()

	if ok {
		return template
	}

	// If not in cache and repo is set, try to load from repo
	if s.repo != nil {
		t, err := s.repo.GetTemplate(id)
		if err == nil && t != nil {
			s.mu.Lock()
			s.templates[id] = t
			s.mu.Unlock()
			return t
		}
	}

	return nil
}

// GetTemplatesByClass returns all skill templates for a specific class.
func (s *SkillService) GetTemplatesByClass(nClassID int) []*domain.SkillTemplate {
	// This is inefficient if we don't have all templates loaded.
	// Ideally we should load all templates on startup.
	// For now, let's assume we might need to fetch from DB or cache.

	// Let's try to load all if cache is empty
	s.mu.RLock()
	count := len(s.templates)
	s.mu.RUnlock()

	if count == 0 && s.repo != nil {
		templates, err := s.repo.GetAllTemplates()
		if err == nil {
			s.mu.Lock()
			for _, t := range templates {
				s.templates[t.ID] = t
			}
			s.mu.Unlock()
		}
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*domain.SkillTemplate
	for _, t := range s.templates {
		if t.NClassID == nClassID || t.NClassID == 0 { // 0 might be shared? Or check DB.
			// Based on DB, nclass_id: 0=Trai Dat, 1=Namec, 2=Xayda.
			// Shared skills might have specific ID or logic.
			result = append(result, t)
		}
	}
	return result
}

// LearnSkill adds a skill to player's list.
func (s *SkillService) LearnSkill(player *domain.Player, templateID int) error {
	template := s.GetTemplate(templateID)
	if template == nil {
		return fmt.Errorf("skill template %d not found", templateID)
	}

	// Check if already learned
	for _, skill := range player.Skills {
		if skill.TemplateID == templateID {
			return fmt.Errorf("skill already learned")
		}
	}

	newSkill := &domain.PlayerSkill{
		TemplateID:  templateID,
		Template:    template,
		Point:       1, // Start at level 1
		LastTimeUse: 0,
		More:        0,
	}

	player.Skills = append(player.Skills, newSkill)
	return nil
}

// GetSkillsByRace returns starter skills for a race.
// This is a helper for character creation.
func (s *SkillService) GetSkillsByRace(race int8) []int {
	// Return IDs of starter skills.
	// Based on NRO logic:
	// Trai Dat (0): 0 (Kamejoko?) - Wait, need to check DB IDs.
	// Namec (1): ?
	// Xayda (2): ?

	// Let's use some common IDs if we don't know exact ones, or query DB.
	// For now, return empty or some default IDs if we know them.
	// Assuming:
	// 0: Trai Dat -> Skill ID 0 (Kamejoko)
	// 1: Namec -> Skill ID 2 (Masenko)
	// 2: Xayda -> Skill ID 4 (Antomic)
	// These are just guesses. Real IDs are in DB.
	// Let's just return empty for now and let CharacterService handle it or update later.

	switch race {
	case 0: // Trái Đất
		return []int{0, 1} // Chiêu đấm Dragon, Chiêu Kamejoko
	case 1: // Namek
		return []int{2, 3} // Chiêu đấm Demon, Chiêu Masenko
	case 2: // Xayda
		return []int{4, 5} // Chiêu đấm Galick, Chiêu Antomic
	default:
		return []int{0, 1}
	}
}

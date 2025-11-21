package services

import (
	"fmt"
	"nro/src/internal/core/domain"
	"nro/src/internal/core/ports"
	"sync"
)

// NPCService quản lý NPC templates.
//
// EXPLANATION:
// NPCService chịu trách nhiệm:
// 1. Load và quản lý thông tin NPC (NPCTemplate) từ Database.
// 2. Cung cấp dữ liệu về vị trí, hình ảnh, và hội thoại của NPC cho các service khác (như MenuService, Controller).
// 3. Tìm kiếm NPC trong các map cụ thể.
type NPCService struct {
	repo      ports.NPCRepository
	templates map[int]*domain.NPCTemplate
	mu        sync.RWMutex
}

var npcServiceInstance *NPCService
var npcOnce sync.Once

// GetNPCService returns singleton instance.
func GetNPCService() *NPCService {
	npcOnce.Do(func() {
		npcServiceInstance = &NPCService{
			templates: make(map[int]*domain.NPCTemplate),
		}
	})
	return npcServiceInstance
}

func (s *NPCService) SetRepository(repo ports.NPCRepository) {
	s.repo = repo
}

// LoadTemplates loads all NPC templates from the repository.
func (s *NPCService) LoadTemplates() error {
	if s.repo == nil {
		return fmt.Errorf("NPC repository not set")
	}

	templates, err := s.repo.GetTemplates()
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range templates {
		s.templates[t.ID] = t
	}
	fmt.Printf("[NPC] Loaded %d NPC templates from DB\n", len(templates))
	return nil
}

// GetTemplate returns NPC template by ID.
func (s *NPCService) GetTemplate(id int) *domain.NPCTemplate {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.templates[id]
}

// GetNPCsByMap returns all NPCs in a specific map.
func (s *NPCService) GetNPCsByMap(mapID int) []*domain.NPCTemplate {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var npcs []*domain.NPCTemplate
	for _, t := range s.templates {
		if t.MapID == mapID {
			npcs = append(npcs, t)
		}
	}
	return npcs
}

// GetDialogue returns dialogue text.
func (s *NPCService) GetDialogue(npcID int, index int) string {
	npc := s.GetTemplate(npcID)
	if npc == nil || index < 0 || index >= len(npc.Dialogue) {
		return ""
	}
	return npc.Dialogue[index]
}

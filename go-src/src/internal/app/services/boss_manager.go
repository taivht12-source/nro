package services

import (
	"fmt"
	"math"
	"nro/src/internal/core/domain"
	"strings"
	"sync"
	"time"
)

// BossManager manages boss spawning and AI.
type BossManager struct {
	activeBosses map[int]*domain.Boss
	templates    map[int]*domain.BossTemplate
	mu           sync.RWMutex
}

var bossManagerInstance *BossManager
var bossOnce sync.Once

// GetBossManager returns singleton instance.
func GetBossManager() *BossManager {
	bossOnce.Do(func() {
		bossManagerInstance = &BossManager{
			activeBosses: make(map[int]*domain.Boss),
			templates:    make(map[int]*domain.BossTemplate),
		}
		// Start update loop
		go bossManagerInstance.updateLoop()
	})
	return bossManagerInstance
}

// LoadTemplates loads boss templates (mock for now).
func (m *BossManager) LoadTemplates() {
	// Mock Broly
	m.templates[1] = &domain.BossTemplate{
		ID:          1,
		Name:        "Broly",
		Gender:      2, // Xayda
		Outfit:      []int16{294, 295, 296, -1, -1, -1},
		Damage:      1000,
		HP:          []int64{100000, 200000, 500000},
		MapJoin:     []int{5},                            // Dao Kame
		SkillTemp:   [][]int{{0, 1, 1000}, {1, 1, 2000}}, // Dragon, Kamejoko
		TextS:       []string{"Ta là Broly", "Ta sẽ tiêu diệt tất cả"},
		TextM:       []string{"Đỡ đòn này", "Yaaaa!"},
		TextE:       []string{"Ta sẽ quay lại", "Không thể nào..."},
		SecondsRest: 10,
		TypeAppear:  0,
		AIType:      "Broly",
	}
	fmt.Println("[BOSS] Loaded mock boss templates")

	// Spawn test boss
	m.CreateBoss(1)
}

// CreateBoss creates a new boss instance from template.
func (m *BossManager) CreateBoss(templateID int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	tmpl, ok := m.templates[templateID]
	if !ok {
		fmt.Printf("[BOSS] Template %d not found\n", templateID)
		return
	}

	boss := &domain.Boss{
		ID:           len(m.activeBosses) + 1,
		TemplateID:   templateID,
		Name:         tmpl.Name,
		Gender:       tmpl.Gender,
		Damage:       tmpl.Damage,
		Status:       domain.BossStatusRest,
		Template:     tmpl,
		CurrentLevel: 0,
		LastTimeRest: time.Now().UnixMilli(),
	}
	// Init HP based on level 0
	if len(tmpl.HP) > 0 {
		boss.MaxHP = tmpl.HP[0]
		boss.HP = boss.MaxHP
	}

	m.activeBosses[boss.ID] = boss

	// Assign AI controller based on boss type
	m.assignAI(boss)

	fmt.Printf("[BOSS] Created boss %s (ID: %d)\n", boss.Name, boss.ID)
}

// updateLoop runs the AI loop for bosses.
func (m *BossManager) updateLoop() {
	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		m.mu.Lock()
		for _, boss := range m.activeBosses {
			m.updateBoss(boss)
		}
		m.mu.Unlock()
	}
}

func (m *BossManager) updateBoss(boss *domain.Boss) {
	now := time.Now().UnixMilli()

	// AI Hook: OnUpdate
	if boss.AI != nil {
		boss.AI.OnUpdate(boss)
	}

	switch boss.Status {
	case domain.BossStatusRest:
		if now-boss.LastTimeRest >= int64(boss.Template.SecondsRest*1000) {
			boss.Status = domain.BossStatusRespawn
		}
	case domain.BossStatusRespawn:
		boss.Status = domain.BossStatusJoinMap
		boss.HP = boss.MaxHP
		boss.MP = boss.MaxMP
		fmt.Printf("[BOSS] %s is respawning...\n", boss.Name)
	case domain.BossStatusJoinMap:
		m.joinMap(boss)
		boss.Status = domain.BossStatusChatS
		boss.LastTimeChatS = now
		boss.IndexChatS = 0
		if boss.AI != nil {
			boss.AI.OnJoinMap(boss)
		}
	case domain.BossStatusChatS:
		if m.chatS(boss) {
			boss.Status = domain.BossStatusActive
			boss.LastTimeChatM = now
			boss.TimeChatM = 5000
		}
	case domain.BossStatusActive:
		m.chatM(boss)
		m.active(boss)
	case domain.BossStatusDie:
		boss.Status = domain.BossStatusChatE
		boss.LastTimeChatE = now
		boss.IndexChatE = 0
	case domain.BossStatusChatE:
		if m.chatE(boss) {
			boss.Status = domain.BossStatusLeaveMap
		}
	case domain.BossStatusLeaveMap:
		m.leaveMap(boss)
	}
}

func (m *BossManager) joinMap(boss *domain.Boss) {
	// Simple logic: Pick first map in MapJoin
	if len(boss.Template.MapJoin) > 0 {
		mapID := boss.Template.MapJoin[0]
		// TODO: Get Zone from MapService
		boss.MapID = mapID
		boss.ZoneID = 0 // Default zone 0
		boss.X = 100
		boss.Y = 300
		fmt.Printf("[BOSS] %s joined Map %d Zone %d\n", boss.Name, boss.MapID, boss.ZoneID)
	}
}

func (m *BossManager) chatS(boss *domain.Boss) bool {
	now := time.Now().UnixMilli()
	if now-boss.LastTimeChatS >= 1000 { // 1s delay between chats
		if boss.IndexChatS >= len(boss.Template.TextS) {
			return true
		}
		text := boss.Template.TextS[boss.IndexChatS]
		fmt.Printf("[BOSS] %s: %s\n", boss.Name, text)
		// TODO: Send chat packet to zone
		boss.LastTimeChatS = now
		boss.IndexChatS++
	}
	return false
}

func (m *BossManager) chatM(boss *domain.Boss) {
	now := time.Now().UnixMilli()
	if now-boss.LastTimeChatM >= int64(boss.TimeChatM) {
		if len(boss.Template.TextM) > 0 {
			text := boss.Template.TextM[now%int64(len(boss.Template.TextM))] // Random-ish
			fmt.Printf("[BOSS] %s: %s\n", boss.Name, text)
			// TODO: Send chat packet
		}
		boss.LastTimeChatM = now
		boss.TimeChatM = 5000 + int(now%5000) // Random 5-10s
	}
}

func (m *BossManager) active(boss *domain.Boss) {
	// 1. Find target
	target := m.findTarget(boss)
	if target == nil {
		return
	}

	// 2. Check distance
	dist := m.calculateDistance(boss.X, boss.Y, target.X, target.Y)

	// 3. Move or Attack
	if dist > 50 { // Attack range 50
		m.moveTo(boss, target.X, target.Y)
	} else {
		// Attack
		now := time.Now().UnixMilli()
		if now-boss.LastTimeAttack >= 2000 { // 2s attack speed
			// AI Hook: OnAttack
			if boss.AI != nil {
				boss.AI.OnAttack(boss, target)
			} else {
				// Default attack if no AI override
				combatService := GetCombatService()
				combatService.AttackPlayer(boss, target)
			}
			boss.LastTimeAttack = now
		}
	}
}

func (m *BossManager) findTarget(boss *domain.Boss) *domain.Player {
	zoneService := GetZoneService()
	players := zoneService.GetPlayersInZone(boss.MapID, boss.ZoneID)

	var nearest *domain.Player
	minDist := 999999.0

	for _, p := range players {
		if p.HP <= 0 {
			continue
		}
		dist := m.calculateDistance(boss.X, boss.Y, p.X, p.Y)
		if dist < minDist {
			minDist = dist
			nearest = p
		}
	}
	return nearest
}

func (m *BossManager) calculateDistance(x1, y1, x2, y2 int16) float64 {
	return math.Sqrt(math.Pow(float64(x1-x2), 2) + math.Pow(float64(y1-y2), 2))
}

func (m *BossManager) moveTo(boss *domain.Boss, targetX, targetY int16) {
	// Simple move logic: move 20 units towards target
	dx := targetX - boss.X
	dy := targetY - boss.Y

	// Normalize
	dist := m.calculateDistance(boss.X, boss.Y, targetX, targetY)
	if dist > 0 {
		step := 20.0
		if step > dist {
			step = dist
		}
		boss.X += int16(float64(dx) / dist * step)
		boss.Y += int16(float64(dy) / dist * step)
		// fmt.Printf("[BOSS] %s moved to (%d, %d)\n", boss.Name, boss.X, boss.Y)
	}
}

func (m *BossManager) chatE(boss *domain.Boss) bool {
	now := time.Now().UnixMilli()
	if now-boss.LastTimeChatE >= 1000 {
		if boss.IndexChatE >= len(boss.Template.TextE) {
			return true
		}
		text := boss.Template.TextE[boss.IndexChatE]
		fmt.Printf("[BOSS] %s: %s\n", boss.Name, text)
		boss.LastTimeChatE = now
		boss.IndexChatE++
	}
	return false
}

func (m *BossManager) leaveMap(boss *domain.Boss) {
	fmt.Printf("[BOSS] %s left the map\n", boss.Name)
	boss.Status = domain.BossStatusRest
	boss.LastTimeRest = time.Now().UnixMilli()
	// TODO: Remove from zone
}

// assignAI creates and assigns the appropriate AI controller for a boss.
func (m *BossManager) assignAI(boss *domain.Boss) {
	if boss.Template.AIType == "" {
		// Try to guess from name if not specified
		if strings.Contains(strings.ToLower(boss.Name), "broly") {
			boss.Template.AIType = "Broly"
		} else {
			boss.Template.AIType = "Default"
		}
	}

	ai, err := GetBossRegistry().CreateAI(boss.Template.AIType)
	if err != nil {
		fmt.Printf("[BOSS] Failed to create AI for type %s: %v\n", boss.Template.AIType, err)
		return
	}

	boss.AI = ai
	ai.OnSpawn(boss)
	fmt.Printf("[BOSS] Assigned AI %s to %s\n", boss.Template.AIType, boss.Name)
}

package services

import (
	"fmt"
	"nro-go/internal/core/domain"
	"strings"
	"sync"
	"time"
)

// BossManager manages boss spawning and AI.
type BossManager struct {
	activeBosses map[int]*domain.Boss
	templates    map[int]*domain.BossTemplate
	// androidAI    map[int]*AndroidAI   // Android AI controllers by boss ID
	blackGokuAI map[int]*BlackGokuAI // Black Goku/Zamasu AI controllers by boss ID
	brolyAI     map[int]*BrolyAI     // Broly AI controllers by boss ID
	cellAI      map[int]*CellAI      // Cell AI controllers by boss ID
	friezaAI    map[int]*FriezaAI    // Frieza AI controllers by boss ID
	nappaAI     map[int]*NappaAI     // Nappa/Saiyan AI controllers by boss ID
	minorBossAI map[int]*MinorBossAI // Minor boss AI controllers by boss ID
	mu          sync.RWMutex
}

var bossManagerInstance *BossManager
var bossOnce sync.Once

// GetBossManager returns singleton instance.
func GetBossManager() *BossManager {
	bossOnce.Do(func() {
		bossManagerInstance = &BossManager{
			activeBosses: make(map[int]*domain.Boss),
			templates:    make(map[int]*domain.BossTemplate),
			// androidAI:    make(map[int]*AndroidAI),
			blackGokuAI: make(map[int]*BlackGokuAI),
			brolyAI:     make(map[int]*BrolyAI),
			cellAI:      make(map[int]*CellAI),
			friezaAI:    make(map[int]*FriezaAI),
			nappaAI:     make(map[int]*NappaAI),
			minorBossAI: make(map[int]*MinorBossAI),
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
	// TODO: Implement attack logic
	// 1. Find target
	// 2. Move to target
	// 3. Use skill
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
	bossName := strings.ToLower(boss.Name)
	bossID := boss.ID

	// // Android bosses
	// if strings.Contains(bossName, "android") || strings.Contains(bossName, "số") {
	// 	androidID := 13 // Default
	// 	if strings.Contains(bossName, "14") {
	// 		androidID = 14
	// 	} else if strings.Contains(bossName, "15") {
	// 		androidID = 15
	// 	} else if strings.Contains(bossName, "20") || strings.Contains(bossName, "dr gero") {
	// 		androidID = 20
	// 	}
	// 	ai := NewAndroidAI(androidID)
	// 	m.androidAI[bossID] = ai
	// 	ai.OnSpawn(boss)
	// 	boss.AI = ai
	// 	fmt.Printf("[BOSS] Assigned AndroidAI (ID: %d) to %s\n", androidID, boss.Name)
	// 	return
	// }

	// Black Goku / Zamasu
	if strings.Contains(bossName, "black") || strings.Contains(bossName, "zamasu") {
		isZamasu := strings.Contains(bossName, "zamasu")
		ai := NewBlackGokuAI(isZamasu)
		m.blackGokuAI[bossID] = ai
		ai.OnSpawn(boss)
		boss.AI = ai
		fmt.Printf("[BOSS] Assigned BlackGokuAI (Zamasu: %v) to %s\n", isZamasu, boss.Name)
		return
	}

	// Broly
	if strings.Contains(bossName, "broly") {
		ai := NewBrolyAI()
		m.brolyAI[bossID] = ai
		ai.OnSpawn(boss)
		boss.AI = ai
		fmt.Printf("[BOSS] Assigned BrolyAI to %s\n", boss.Name)
		return
	}

	// Cell
	if strings.Contains(bossName, "cell") || strings.Contains(bossName, "xên") {
		ai := NewCellAI()
		m.cellAI[bossID] = ai
		ai.OnSpawn(boss)
		boss.AI = ai
		fmt.Printf("[BOSS] Assigned CellAI to %s\n", boss.Name)
		return
	}

	// Frieza
	if strings.Contains(bossName, "frieza") || strings.Contains(bossName, "fide") {
		isGolden := strings.Contains(bossName, "golden") || strings.Contains(bossName, "vàng")
		ai := NewFriezaAI(isGolden)
		m.friezaAI[bossID] = ai
		ai.OnSpawn(boss)
		boss.AI = ai
		fmt.Printf("[BOSS] Assigned FriezaAI (Golden: %v) to %s\n", isGolden, boss.Name)
		return
	}

	// Nappa / Saiyan
	if strings.Contains(bossName, "nappa") || strings.Contains(bossName, "vegeta") ||
		strings.Contains(bossName, "raditz") || strings.Contains(bossName, "kakarot") {
		saiyanType := "nappa"
		if strings.Contains(bossName, "vegeta") {
			saiyanType = "vegeta"
		} else if strings.Contains(bossName, "raditz") {
			saiyanType = "raditz"
		}
		ai := NewNappaAI(saiyanType)
		m.nappaAI[bossID] = ai
		ai.OnSpawn(boss)
		boss.AI = ai
		fmt.Printf("[BOSS] Assigned NappaAI (Type: %s) to %s\n", saiyanType, boss.Name)
		return
	}

	// Default: Minor Boss AI
	ai := NewMinorBossAI(bossName)
	m.minorBossAI[bossID] = ai
	ai.OnSpawn(boss)
	boss.AI = ai
	fmt.Printf("[BOSS] Assigned MinorBossAI to %s\n", boss.Name)
}

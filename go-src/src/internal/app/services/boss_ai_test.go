package services

import (
	"nro/src/internal/core/domain"
	"testing"
)

// MockBossAI for testing
type MockBossAI struct {
	domain.BaseBossAI
	OnDamagedCalled bool
	OnDieCalled     bool
	OnRewardCalled  bool
}

func (m *MockBossAI) OnDamaged(boss *domain.Boss, damage int64, attacker *domain.Player) {
	m.OnDamagedCalled = true
}

func (m *MockBossAI) OnDie(boss *domain.Boss, killer *domain.Player) {
	m.OnDieCalled = true
}

func (m *MockBossAI) OnReward(boss *domain.Boss, killer *domain.Player) {
	m.OnRewardCalled = true
}

func TestCombatService_BossCallbacks(t *testing.T) {
	// Setup
	combatService := GetCombatService()
	mockAI := &MockBossAI{}
	boss := &domain.Boss{
		ID:     1,
		Name:   "Test Boss",
		HP:     1000,
		MaxHP:  1000,
		Damage: 100,
		AI:     mockAI,
	}
	player := &domain.Player{
		ID:   1,
		Name: "Test Player",
	}

	// Test OnDamaged
	combatService.TakeDamageBoss(boss, 100, player)
	if !mockAI.OnDamagedCalled {
		t.Errorf("Expected OnDamaged to be called")
	}
	if boss.HP != 900 {
		t.Errorf("Expected Boss HP to be 900, got %d", boss.HP)
	}

	// Test OnDie and OnReward
	combatService.TakeDamageBoss(boss, 900, player) // Kill boss
	if !mockAI.OnDieCalled {
		t.Errorf("Expected OnDie to be called")
	}
	if !mockAI.OnRewardCalled {
		t.Errorf("Expected OnReward to be called")
	}
	if boss.HP != 0 {
		t.Errorf("Expected Boss HP to be 0, got %d", boss.HP)
	}
}

func TestBrolyAI_RageMode(t *testing.T) {
	// Setup
	// We need to register Broly AI first, but we can just instantiate it directly for testing
	// assuming BrolyAI is exported or we can access it.
	// BrolyAI is in services package, so we can access it.

	// Re-create BrolyAI instance manually since we can't easily reset the registry
	brolyAI := &BrolyAI{rageActivated: false}

	boss := &domain.Boss{
		ID:     1,
		Name:   "Broly",
		HP:     100000,
		MaxHP:  100000,
		Damage: 1000,
		Template: &domain.BossTemplate{
			Damage: 1000,
		},
		AI: brolyAI,
	}
	player := &domain.Player{ID: 1, Name: "Goku"}

	// Damage to 60% - No Rage
	brolyAI.OnDamaged(boss, 40000, player) // 60k left
	if boss.Damage != 1000 {
		t.Errorf("Expected Damage to be 1000, got %d", boss.Damage)
	}

	// Damage to 40% - Rage Mode
	boss.HP = 40000
	brolyAI.OnDamaged(boss, 0, player) // Trigger check
	if boss.Damage != 1500 {
		t.Errorf("Expected Damage to be 1500 (1.5x), got %d", boss.Damage)
	}
}

func TestCellAI_Regeneration(t *testing.T) {
	cellAI := NewCellAI()
	boss := &domain.Boss{
		ID:    2,
		Name:  "Cell",
		HP:    100000,
		MaxHP: 100000,
		AI:    cellAI,
	}
	player := &domain.Player{ID: 1, Name: "Gohan"}

	// Damage to 40%
	boss.HP = 40000

	// Trigger regeneration
	// We need to mock time or wait, but the logic uses time.Now().
	// For testing, we can assume the condition (now - lastRegenerateTime > 60000) is met initially since lastRegenerateTime is 0.

	cellAI.OnDamaged(boss, 0, player)

	expectedHP := 40000 + 25000 // 40k + 25k (1/4 MaxHP)
	if boss.HP != int64(expectedHP) {
		t.Errorf("Expected HP to be %d, got %d", expectedHP, boss.HP)
	}

	// Try again immediately - should not regenerate
	cellAI.OnDamaged(boss, 0, player)
	if boss.HP != int64(expectedHP) {
		t.Errorf("Expected HP to remain %d, got %d", expectedHP, boss.HP)
	}
}

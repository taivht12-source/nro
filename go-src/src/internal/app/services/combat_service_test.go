package services

import (
	"nro-go/internal/core/domain"
	"testing"
)

func TestBossCombat(t *testing.T) {
	combatService := GetCombatService()

	boss := &domain.Boss{
		ID:     1,
		Name:   "Test Boss",
		Damage: 1000,
		HP:     10000,
		MaxHP:  10000,
	}

	player := &domain.Player{
		ID:    1,
		Name:  "Test Player",
		HP:    2000,
		MaxHP: 2000,
	}

	// Test CalculateBossDamage
	damage := combatService.CalculateBossDamage(boss, player)
	if damage < 900 || damage > 1100 {
		t.Errorf("Expected damage between 900 and 1100, got %d", damage)
	}

	// Test AttackPlayer
	initialHP := player.HP
	dealt, err := combatService.AttackPlayer(boss, player)
	if err != nil {
		t.Errorf("AttackPlayer returned error: %v", err)
	}

	if player.HP != initialHP-dealt {
		t.Errorf("Player HP not updated correctly. Expected %d, got %d", initialHP-dealt, player.HP)
	}

	// Test TakeDamageBoss
	combatService.TakeDamageBoss(boss, 500, player)
	if boss.HP != 9500 {
		t.Errorf("Boss HP not updated correctly. Expected 9500, got %d", boss.HP)
	}
}

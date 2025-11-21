package services

import (
	"fmt"
	"nro-go/internal/core/domain"
)

// MinorBossAI implements behavior for minor bosses.
type MinorBossAI struct {
	domain.BaseBossAI
	bossType string
}

func NewMinorBossAI(bossType string) *MinorBossAI {
	return &MinorBossAI{
		bossType: bossType,
	}
}

func (m *MinorBossAI) OnSpawn(boss *domain.Boss) {
	fmt.Printf("[BOSS] %s has appeared!\n", boss.Name)
}

func (m *MinorBossAI) OnDamaged(boss *domain.Boss, damage int64, attacker *domain.Player) {
	// Special behaviors for specific minor bosses
	switch m.bossType {
	case "tau_paypay":
		// Mercenary Tao - uses pillar for mobility
		hpPercent := float64(boss.HP) / float64(boss.MaxHP)
		if hpPercent < 0.5 {
			// Flee behavior
			fmt.Printf("[TAU PAYPAY] %s attempts to flee!\n", boss.Name)
		}
	case "king_cold":
		// King Cold - similar to Frieza but weaker
		// No special mechanics, just basic boss
	}
}

func (m *MinorBossAI) OnDie(boss *domain.Boss, killer *domain.Player) {
	fmt.Printf("[BOSS] %s has been defeated!\n", boss.Name)
}

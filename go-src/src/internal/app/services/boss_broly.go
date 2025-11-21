package services

import (
	"fmt"
	"nro-go/internal/core/domain"
)

// BrolyAI implements Broly-specific behavior.
type BrolyAI struct {
	rageActivated bool
}

func init() {
	GetBossRegistry().Register("Broly", func() domain.BossAI {
		return &BrolyAI{
			rageActivated: false,
		}
	})
}

func (b *BrolyAI) OnSpawn(boss *domain.Boss) {
	b.rageActivated = false
	fmt.Printf("[BROLY] %s has spawned!\n", boss.Name)
}

func (b *BrolyAI) OnUpdate(boss *domain.Boss) {
	// Broly specific update logic
}

func (b *BrolyAI) OnAttack(boss *domain.Boss, target *domain.Player) {
	// Broly specific attack logic
	fmt.Printf("[BROLY] %s attacks %s with overwhelming power!\n", boss.Name, target.Name)
}

func (b *BrolyAI) OnDamaged(boss *domain.Boss, damage int64, attacker *domain.Player) {
	// Rage mode when HP drops below 50%
	hpPercent := float64(boss.HP) / float64(boss.MaxHP)

	if !b.rageActivated && hpPercent < 0.5 {
		b.rageActivated = true
		// Increase damage by 50%
		boss.Damage = int(float64(boss.Damage) * 1.5)
		fmt.Printf("[BROLY] %s enters RAGE MODE! Damage increased!\n", boss.Name)
		// TODO: Send chat message to zone
	}

	// Further power increase at 25% HP
	if hpPercent < 0.25 && boss.Damage < boss.Template.Damage*2 {
		boss.Damage = boss.Template.Damage * 2
		fmt.Printf("[BROLY] %s is going BERSERK!\n", boss.Name)
	}
}

func (b *BrolyAI) OnDie(boss *domain.Boss, killer *domain.Player) {
	b.rageActivated = false
	fmt.Printf("[BROLY] %s has been defeated by %s!\n", boss.Name, killer.Name)
}

func (b *BrolyAI) OnJoinMap(boss *domain.Boss) {
	fmt.Printf("[BROLY] %s joined the map! The earth shakes!\n", boss.Name)
}

func (b *BrolyAI) OnReward(boss *domain.Boss, killer *domain.Player) {
	fmt.Printf("[BROLY] %s dropped rewards for %s\n", boss.Name, killer.Name)
}

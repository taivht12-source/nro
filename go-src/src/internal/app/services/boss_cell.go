package services

import (
	"fmt"
	"nro/src/internal/core/domain"
	"time"
)

// CellAI implements Cell-specific behavior.
type CellAI struct {
	domain.BaseBossAI
	form               int // 0: Imperfect, 1: Semi-Perfect, 2: Perfect
	canRegenerate      bool
	lastRegenerateTime int64
}

func NewCellAI() *CellAI {
	return &CellAI{
		form:          0,
		canRegenerate: true,
	}
}

func (c *CellAI) OnSpawn(boss *domain.Boss) {
	c.form = 0
	c.canRegenerate = true
	fmt.Printf("[CELL] %s has spawned in Imperfect form!\n", boss.Name)
}

func (c *CellAI) OnDamaged(boss *domain.Boss, damage int64, attacker *domain.Player) {
	// Regeneration ability (once per minute)
	now := time.Now().UnixMilli()
	if c.canRegenerate && boss.HP < boss.MaxHP/2 && now-c.lastRegenerateTime > 60000 {
		healAmount := boss.MaxHP / 4
		boss.HP += healAmount
		if boss.HP > boss.MaxHP {
			boss.HP = boss.MaxHP
		}
		c.lastRegenerateTime = now
		fmt.Printf("[CELL] %s regenerates %d HP!\n", boss.Name, healAmount)
	}
}

func (c *CellAI) OnUpdate(boss *domain.Boss) {
	// Evolution logic based on HP and kills
	// This would need integration with combat system
}

func (c *CellAI) Evolve(boss *domain.Boss) {
	if c.form >= 2 {
		return // Already perfect
	}

	c.form++
	// Increase stats
	boss.MaxHP = int64(float64(boss.MaxHP) * 1.5)
	boss.HP = boss.MaxHP
	boss.Damage = int(float64(boss.Damage) * 1.3)

	formName := []string{"Semi-Perfect", "Perfect"}
	fmt.Printf("[CELL] %s has evolved to %s form!\n", boss.Name, formName[c.form-1])
}

func (c *CellAI) OnDie(boss *domain.Boss, killer *domain.Player) {
	// Cell can self-destruct and regenerate (once)
	if c.canRegenerate && c.form == 2 {
		fmt.Printf("[CELL] %s self-destructs!\n", boss.Name)
		// TODO: Deal area damage
		c.canRegenerate = false
	}
}

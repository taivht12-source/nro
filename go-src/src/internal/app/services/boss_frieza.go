package services

import (
	"fmt"
	"nro/src/internal/core/domain"
)

// FriezaAI implements Frieza-specific behavior.
type FriezaAI struct {
	domain.BaseBossAI
	currentForm int // 1-4 (First, Second, Third, Final)
	isGolden    bool
}

func NewFriezaAI(isGolden bool) *FriezaAI {
	return &FriezaAI{
		currentForm: 1,
		isGolden:    isGolden,
	}
}

func (f *FriezaAI) OnSpawn(boss *domain.Boss) {
	f.currentForm = 1
	if f.isGolden {
		fmt.Printf("[FRIEZA] Golden %s has appeared!\n", boss.Name)
	} else {
		fmt.Printf("[FRIEZA] %s has spawned in First Form!\n", boss.Name)
	}
}

func (f *FriezaAI) OnDamaged(boss *domain.Boss, damage int64, attacker *domain.Player) {
	if f.isGolden {
		// Golden form drains stamina but has high power
		return
	}

	// Transform based on HP thresholds
	hpPercent := float64(boss.HP) / float64(boss.MaxHP)

	if f.currentForm == 1 && hpPercent < 0.75 {
		f.Transform(boss, 2)
	} else if f.currentForm == 2 && hpPercent < 0.50 {
		f.Transform(boss, 3)
	} else if f.currentForm == 3 && hpPercent < 0.25 {
		f.Transform(boss, 4)
	}
}

func (f *FriezaAI) Transform(boss *domain.Boss, newForm int) {
	if newForm <= f.currentForm || newForm > 4 {
		return
	}

	f.currentForm = newForm

	// Stat increases per form
	multiplier := 1.0 + float64(newForm-1)*0.3
	boss.Damage = int(float64(boss.Template.Damage) * multiplier)
	boss.MaxHP = int64(float64(boss.Template.HP[0]) * multiplier)
	boss.HP = boss.MaxHP // Full heal on transform

	formNames := []string{"", "First", "Second", "Third", "Final"}
	fmt.Printf("[FRIEZA] %s transforms to %s Form!\n", boss.Name, formNames[newForm])
}

func (f *FriezaAI) OnDie(boss *domain.Boss, killer *domain.Player) {
	if f.isGolden {
		fmt.Printf("[FRIEZA] Golden %s has been defeated!\n", boss.Name)
	} else {
		fmt.Printf("[FRIEZA] %s (%s Form) has been defeated!\n", boss.Name,
			[]string{"", "First", "Second", "Third", "Final"}[f.currentForm])
	}
}

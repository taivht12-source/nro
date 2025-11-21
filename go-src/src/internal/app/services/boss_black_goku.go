package services

import (
	"fmt"
	"nro/src/internal/core/domain"
	"time"
)

// BlackGokuAI implements Black Goku and Zamasu behavior.
type BlackGokuAI struct {
	domain.BaseBossAI
	isZamasu      bool
	isRoseForm    bool
	isFused       bool
	fusionPartner *domain.Boss
	lastHealTime  int64
}

func NewBlackGokuAI(isZamasu bool) *BlackGokuAI {
	return &BlackGokuAI{
		isZamasu:   isZamasu,
		isRoseForm: false,
		isFused:    false,
	}
}

func (bg *BlackGokuAI) OnSpawn(boss *domain.Boss) {
	if bg.isZamasu {
		fmt.Printf("[ZAMASU] %s has appeared!\n", boss.Name)
	} else {
		fmt.Printf("[BLACK GOKU] %s has appeared!\n", boss.Name)
	}
}

func (bg *BlackGokuAI) OnUpdate(boss *domain.Boss) {
	// Zamasu's immortality - regenerate HP over time
	if bg.isZamasu && !bg.isFused {
		now := time.Now().UnixMilli()
		if now-bg.lastHealTime > 5000 { // Heal every 5 seconds
			healAmount := boss.MaxHP / 20 // 5% HP
			boss.HP += healAmount
			if boss.HP > boss.MaxHP {
				boss.HP = boss.MaxHP
			}
			bg.lastHealTime = now
			fmt.Printf("[ZAMASU] %s regenerates %d HP (Immortality)!\n", boss.Name, healAmount)
		}
	}

	// Check fusion conditions
	bg.checkFusionConditions(boss)
}

func (bg *BlackGokuAI) OnDamaged(boss *domain.Boss, damage int64, attacker *domain.Player) {
	// Black Goku transforms to Rose form at 50% HP
	if !bg.isZamasu && !bg.isRoseForm {
		hpPercent := float64(boss.HP) / float64(boss.MaxHP)
		if hpPercent < 0.5 {
			bg.transformToRose(boss)
		}
	}

	// Fused Zamasu takes reduced damage
	if bg.isFused {
		// Damage already applied, but log it
		fmt.Printf("[FUSED ZAMASU] %s shrugs off the attack!\n", boss.Name)
	}
}

func (bg *BlackGokuAI) OnDie(boss *domain.Boss, killer *domain.Player) {
	if bg.isFused {
		fmt.Printf("[FUSED ZAMASU] The immortal fusion has been defeated!\n")
	} else if bg.isZamasu {
		fmt.Printf("[ZAMASU] %s has been vanquished!\n", boss.Name)
	} else {
		fmt.Printf("[BLACK GOKU] %s has been defeated!\n", boss.Name)
	}
}

func (bg *BlackGokuAI) transformToRose(boss *domain.Boss) {
	bg.isRoseForm = true
	boss.Name = "Black Goku (Super Saiyan Rose)"
	boss.Damage = int(float64(boss.Damage) * 1.5)

	// Visual change would update outfit here
	fmt.Printf("[BLACK GOKU] %s transforms into Super Saiyan Rose!\n", boss.Name)
}

func (bg *BlackGokuAI) checkFusionConditions(boss *domain.Boss) {
	// Check if both Black Goku and Zamasu are present and below 30% HP
	// This would require zone/boss manager integration
	if bg.fusionPartner != nil && !bg.isFused {
		// Check conditions
		hpPercent := float64(boss.HP) / float64(boss.MaxHP)
		partnerHpPercent := float64(bg.fusionPartner.HP) / float64(bg.fusionPartner.MaxHP)

		if hpPercent < 0.3 && partnerHpPercent < 0.3 {
			bg.fuse(boss)
		}
	}
}

func (bg *BlackGokuAI) fuse(boss *domain.Boss) {
	bg.isFused = true
	boss.Name = "Fused Zamasu"
	boss.MaxHP = int64(float64(boss.MaxHP) * 2.5)
	boss.HP = boss.MaxHP // Full heal on fusion
	boss.Damage = int(float64(boss.Damage) * 2.0)

	// Remove fusion partner
	if bg.fusionPartner != nil {
		bg.fusionPartner.HP = 0
	}

	fmt.Printf("[FUSION] Black Goku and Zamasu fuse into Fused Zamasu!\n")
}

// SetFusionPartner sets the other boss for fusion.
func (bg *BlackGokuAI) SetFusionPartner(partner *domain.Boss) {
	bg.fusionPartner = partner
}

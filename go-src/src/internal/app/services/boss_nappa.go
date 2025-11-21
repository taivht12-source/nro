package services

import (
	"fmt"
	"nro/src/internal/core/domain"
	"time"
)

// NappaAI implements Nappa and Saiyan-specific behavior.
type NappaAI struct {
	domain.BaseBossAI
	saiyanType      string // "nappa", "vegeta", "raditz"
	isGreatApe      bool
	canTransform    bool
	saibamenSpawned int
	maxSaibamen     int
	lastSpawnTime   int64
}

func NewNappaAI(saiyanType string) *NappaAI {
	return &NappaAI{
		saiyanType:   saiyanType,
		isGreatApe:   false,
		canTransform: true,
		maxSaibamen:  6,
	}
}

func (n *NappaAI) OnSpawn(boss *domain.Boss) {
	n.isGreatApe = false
	n.saibamenSpawned = 0
	fmt.Printf("[SAIYAN] %s has arrived on Earth!\n", boss.Name)
}

func (n *NappaAI) OnUpdate(boss *domain.Boss) {
	// Nappa spawns Saibamen at the start
	if n.saiyanType == "nappa" && n.saibamenSpawned < n.maxSaibamen {
		now := time.Now().UnixMilli()
		if now-n.lastSpawnTime > 10000 { // Spawn every 10 seconds
			n.spawnSaibaman(boss)
			n.lastSpawnTime = now
		}
	}
}

func (n *NappaAI) OnDamaged(boss *domain.Boss, damage int64, attacker *domain.Player) {
	// Transform to Great Ape at 30% HP (if moon is present)
	if n.canTransform && !n.isGreatApe {
		hpPercent := float64(boss.HP) / float64(boss.MaxHP)
		if hpPercent < 0.3 {
			n.transformToGreatApe(boss)
		}
	}
}

func (n *NappaAI) OnDie(boss *domain.Boss, killer *domain.Player) {
	if n.isGreatApe {
		fmt.Printf("[SAIYAN] The Great Ape %s has been defeated!\n", boss.Name)
	} else {
		fmt.Printf("[SAIYAN] %s has been defeated!\n", boss.Name)
	}
}

func (n *NappaAI) transformToGreatApe(boss *domain.Boss) {
	n.isGreatApe = true
	n.canTransform = false

	// Massive stat boost
	boss.MaxHP = int64(float64(boss.MaxHP) * 10.0)
	boss.HP = boss.MaxHP // Full heal
	boss.Damage = int(float64(boss.Damage) * 5.0)

	// Update appearance (would change outfit here)
	boss.Name = fmt.Sprintf("Great Ape %s", boss.Name)

	fmt.Printf("[SAIYAN] %s transforms into a Great Ape!\n", boss.Name)
}

func (n *NappaAI) spawnSaibaman(boss *domain.Boss) {
	n.saibamenSpawned++
	fmt.Printf("[NAPPA] Nappa plants a Saibaman! (%d/%d)\n", n.saibamenSpawned, n.maxSaibamen)
	// TODO: Actually spawn Saibaman mob in zone
	// This would require integration with zone/mob system
}

// RemoveTail can be called to prevent Great Ape transformation.
func (n *NappaAI) RemoveTail() {
	if n.isGreatApe {
		// Revert from Great Ape form
		n.isGreatApe = false
		fmt.Println("[SAIYAN] The tail has been cut! Reverting to normal form!")
	}
	n.canTransform = false
}

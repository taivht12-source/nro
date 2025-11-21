package services

import (
	"fmt"
	"nro-go/internal/core/domain"
)

// AndroidAI implements Android-specific behavior with group mechanics.
type AndroidAI struct {
	androidID       int
	canSelfDestruct bool
	groupMembers    []*domain.Boss // Other androids in the group
}

func init() {
	GetBossRegistry().Register("Android", func() domain.BossAI {
		return &AndroidAI{
			androidID:       13, // Default, will be updated in OnSpawn or similar if needed
			canSelfDestruct: false,
			groupMembers:    make([]*domain.Boss, 0),
		}
	})
}

func (a *AndroidAI) OnSpawn(boss *domain.Boss) {
	if boss.Template.AndroidID != 0 {
		a.androidID = boss.Template.AndroidID
	}
	a.canSelfDestruct = a.androidID == 20
	fmt.Printf("[ANDROID] Android %d has spawned!\n", a.androidID)
}

func (a *AndroidAI) OnUpdate(boss *domain.Boss) {
	// Check for fusion conditions
	// Android 13, 14, 15 can fuse into Super Android 13
	if a.androidID == 13 {
		a.checkFusionConditions(boss)
	}
}

func (a *AndroidAI) OnAttack(boss *domain.Boss, target *domain.Player) {
	// Default attack behavior - can be customized per android type
	fmt.Printf("[ANDROID] Android %d attacks %s!\n", a.androidID, target.Name)
	// Call default attack from CombatService if needed, or implement custom skill logic here
}

func (a *AndroidAI) OnDamaged(boss *domain.Boss, damage int64, attacker *domain.Player) {
	// Android 20 can absorb energy
	if a.androidID == 20 {
		// Absorb a portion of damage as HP
		absorbAmount := damage / 10
		boss.HP += absorbAmount
		if boss.HP > boss.MaxHP {
			boss.HP = boss.MaxHP
		}
		fmt.Printf("[ANDROID] Android 20 absorbs %d HP!\n", absorbAmount)
	}
}

func (a *AndroidAI) OnDie(boss *domain.Boss, killer *domain.Player) {
	if a.canSelfDestruct && a.androidID == 20 {
		fmt.Printf("[ANDROID] Android 20 self-destructs!\n")
		// TODO: Deal area damage to nearby players
		// Damage calculation: boss.Damage * 2 in radius 200
	}

	fmt.Printf("[ANDROID] Android %d has been destroyed!\n", a.androidID)
}

func (a *AndroidAI) OnJoinMap(boss *domain.Boss) {
	fmt.Printf("[ANDROID] Android %d joined the map!\n", a.androidID)
}

func (a *AndroidAI) OnReward(boss *domain.Boss, killer *domain.Player) {
	fmt.Printf("[ANDROID] Android %d dropped rewards for %s\n", a.androidID, killer.Name)
}

func (a *AndroidAI) checkFusionConditions(boss *domain.Boss) {
	// Check if Android 14 and 15 are nearby and alive
	// If so, trigger fusion into Super Android 13
	// This would require zone/boss manager integration

	// Placeholder logic
	if len(a.groupMembers) >= 2 {
		fmt.Printf("[ANDROID] Fusion conditions met! Transforming into Super Android 13!\n")
		a.fuse(boss)
	}
}

func (a *AndroidAI) fuse(boss *domain.Boss) {
	// Transform into Super Android 13
	boss.Name = "Super Android 13"
	boss.MaxHP = int64(float64(boss.MaxHP) * 2.0)
	boss.HP = boss.MaxHP
	boss.Damage = int(float64(boss.Damage) * 1.8)

	// Remove absorbed androids from zone
	for _, member := range a.groupMembers {
		member.HP = 0 // Mark as dead
	}
	a.groupMembers = nil
}

// AddGroupMember adds another android to this group.
func (a *AndroidAI) AddGroupMember(boss *domain.Boss) {
	a.groupMembers = append(a.groupMembers, boss)
}

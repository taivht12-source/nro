package services

import (
	"errors"
	"fmt"
	"math/rand"
	"nro-go/internal/core/domain"
	"time"
)

// CombatService quản lý combat mechanics.
//
// EXPLANATION:
// CombatService chịu trách nhiệm:
// 1. Tính toán sát thương (Damage Calculation) dựa trên chỉ số nhân vật và kỹ năng.
// 2. Xử lý logic sử dụng kỹ năng (Skill Usage): Kiểm tra MP, Cooldown, và áp dụng hiệu ứng.
// 3. Xử lý logic nhận sát thương (Take Damage), hồi phục (Heal), và chết/hồi sinh (Death/Respawn).
type CombatService struct{}

// GetCombatService returns combat service instance.
func GetCombatService() *CombatService {
	return &CombatService{}
}

// getSkillData trả về thông tin chỉ số cho cấp độ hiện tại của kỹ năng.
func (cs *CombatService) getSkillData(skill *domain.PlayerSkill) *domain.SkillLevelData {
	if skill.Template == nil || len(skill.Template.SkillData) == 0 {
		return nil
	}

	// Tìm dữ liệu cho cấp độ (point) hiện tại
	for _, data := range skill.Template.SkillData {
		if data.Point == skill.Point {
			return data
		}
	}

	// Nếu không tìm thấy, trả về cấp độ đầu tiên
	return skill.Template.SkillData[0]
}

// UseSkill sử dụng skill của player.
func (cs *CombatService) UseSkill(player *domain.Player, skillIndex int, targetID int) (int, error) {
	if skillIndex < 0 || skillIndex >= len(player.Skills) {
		return 0, errors.New("invalid skill index")
	}

	skill := player.Skills[skillIndex]
	skillData := cs.getSkillData(skill)
	if skillData == nil {
		return 0, errors.New("skill data not found")
	}

	// Check MP
	if player.MP < skillData.ManaUse {
		return 0, errors.New("not enough MP")
	}

	// Check cooldown
	currentTime := time.Now().UnixMilli()
	if currentTime-skill.LastTimeUse < int64(skillData.CoolDown) {
		remaining := int64(skillData.CoolDown) - (currentTime - skill.LastTimeUse)
		return 0, fmt.Errorf("skill on cooldown (%.1fs remaining)", float64(remaining)/1000.0)
	}

	// Deduct MP
	player.MP -= skillData.ManaUse

	// Update last use time
	skill.LastTimeUse = currentTime

	// Calculate damage
	damage := cs.CalculateDamage(player, skill, skillData)

	// Apply damage to target if it's an attack skill
	// Based on skill_service.go comments:
	// Type 0: Attack (Kamejoko, Masenko, Antomic...)
	// Type 1: Buff (Kaioken...)
	// Type 2: Heal/Support (Tái Tạo...)

	if skill.Template.Type == 0 { // Attack
		if targetID > 0 {
			// Find target player in zone
			zoneService := GetZoneService()
			targetPlayer := cs.findPlayerInZone(player.MapID, player.ZoneID, targetID, zoneService)

			if targetPlayer != nil {
				// Apply damage to target
				isDead := cs.TakeDamage(targetPlayer, damage)

				if isDead {
					fmt.Printf("[COMBAT] %s killed %s with %s!\n",
						player.Name, targetPlayer.Name, skill.Template.Name)
					// TODO: Handle death (respawn, drop items, etc.)
				}
			} else {
				fmt.Printf("[COMBAT] Target %d not found in zone\n", targetID)
			}
		}
	} else if skill.Template.Type == 1 { // Buff
		// Apply buff
		// SKILL vs EFFECT Explanation:
		// - Skill (Kỹ năng): Là hành động kích hoạt (Cause). Ví dụ: Dùng chiêu Kaioken.
		// - Effect (Hiệu ứng): Là trạng thái kết quả (Result). Ví dụ: Tăng sức mạnh trong 60s.
		// Khi dùng Skill Buff (Type 1), ta tạo ra một Effect và gắn vào Player.

		effectService := GetEffectService()
		// Example: Kaioken (ID 4) -> Effect Type 1 (Power Up)

		// Need mapping from Skill ID to Effect Type
		effectType := 0
		if skill.Template.ID == 4 { // Kaioken
			effectType = 1 // Power Up
		}

		effect := &domain.Effect{
			Type:     effectType,
			Value:    skillData.Damage,          // Use damage field as value (e.g. % increase)
			Duration: skillData.CoolDown * 1000, // Duration usually related to cooldown or specific field
		}
		effectService.AddEffect(player, effect)
		fmt.Printf("[COMBAT] %s used Buff %s (Created Effect Type %d)\n", player.Name, skill.Template.Name, effectType)

	} else if skill.Template.Type == 2 { // Heal
		// Logic for Heal (Tái Tạo)
		// Usually restores HP/MP over time or instantly
		// For Tái Tạo (ID 12), it might be a buff that heals over time or instant heal.
		// Let's assume instant heal for now based on previous implementation.

		healAmount := skillData.Damage // Use damage field as heal amount/percent
		healVal := int(float64(player.MaxHP) * float64(healAmount) / 100.0)
		cs.RestoreHP(player, healVal)
		cs.RestoreMP(player, healVal)
		fmt.Printf("[COMBAT] %s used Heal %s - Restored %d\n", player.Name, skill.Template.Name, healVal)
	}

	fmt.Printf("[COMBAT] %s used %s (Level %d) - Damage/Val: %d, MP: %d/%d\n",
		player.Name, skill.Template.Name, skill.Point, damage, player.MP, player.MaxMP)

	return damage, nil
}

// findPlayerInZone tìm player trong zone.
func (cs *CombatService) findPlayerInZone(mapID, zoneID, playerID int, zoneService *ZoneService) *domain.Player {
	players := zoneService.GetPlayersInZone(mapID, zoneID)
	for _, p := range players {
		if p.ID == playerID {
			return p
		}
	}
	return nil
}

// CalculateDamage tính damage của skill.
func (cs *CombatService) CalculateDamage(attacker *domain.Player, skill *domain.PlayerSkill, data *domain.SkillLevelData) int {
	// Base damage from skill
	baseDamage := data.Damage

	// Add player power bonus (1% of power)
	powerBonus := int(attacker.Power / 100)

	// Total damage
	totalDamage := baseDamage + powerBonus

	// Random variance (90-110%)
	variance := rand.Intn(21) - 10 // -10 to +10
	finalDamage := totalDamage * (100 + variance) / 100

	if finalDamage < 1 {
		finalDamage = 1
	}

	return finalDamage
}

// TakeDamage player nhận damage.
func (cs *CombatService) TakeDamage(player *domain.Player, damage int) bool {
	player.HP -= damage

	if player.HP < 0 {
		player.HP = 0
	}

	isDead := player.HP == 0
	fmt.Printf("[COMBAT] %s took %d damage, HP: %d/%d %s\n",
		player.Name, damage, player.HP, player.MaxHP,
		map[bool]string{true: "(DEAD)", false: ""}[isDead])

	return isDead
}

// RestoreHP hồi phục HP.
func (cs *CombatService) RestoreHP(player *domain.Player, amount int) {
	player.HP += amount
	if player.HP > player.MaxHP {
		player.HP = player.MaxHP
	}
	fmt.Printf("[COMBAT] %s restored %d HP, HP: %d/%d\n",
		player.Name, amount, player.HP, player.MaxHP)
}

// RestoreMP hồi phục MP.
func (cs *CombatService) RestoreMP(player *domain.Player, amount int) {
	player.MP += amount
	if player.MP > player.MaxMP {
		player.MP = player.MaxMP
	}
	fmt.Printf("[COMBAT] %s restored %d MP, MP: %d/%d\n",
		player.Name, amount, player.MP, player.MaxMP)
}

// UseHealSkill sử dụng skill hồi phục.
func (cs *CombatService) UseHealSkill(player *domain.Player, skillIndex int) error {
	if skillIndex < 0 || skillIndex >= len(player.Skills) {
		return errors.New("invalid skill index")
	}

	skill := player.Skills[skillIndex]
	skillData := cs.getSkillData(skill)
	if skillData == nil {
		return errors.New("skill data not found")
	}

	// Check if it's a heal skill (Type 2)
	if skill.Template.Type != 2 {
		return errors.New("skill is not a heal skill")
	}

	// Check MP
	if player.MP < skillData.ManaUse {
		return errors.New("not enough MP")
	}

	// Check cooldown
	currentTime := time.Now().UnixMilli()
	if currentTime-skill.LastTimeUse < int64(skillData.CoolDown) {
		return fmt.Errorf("skill on cooldown")
	}

	// Deduct MP
	player.MP -= skillData.ManaUse

	// Update last use time
	skill.LastTimeUse = currentTime

	// Heal amount based on skill level (damage field usually holds heal % or amount)
	healAmount := skillData.Damage // Assuming damage field is used for heal amount/percent

	// If it's percentage based (e.g., 10 means 10%)
	// Need to check description or logic.
	// For now assume direct amount + % of MaxHP?
	// Let's just use the value as direct heal for simplicity, or 10% of MaxHP per point.

	// Sample for Heal (7): damage: 50 (50%?)
	// "Phục hồi #% HP và KI cho đồng đội"
	// So it's percentage.

	healVal := int(float64(player.MaxHP) * float64(healAmount) / 100.0)
	cs.RestoreHP(player, healVal)
	// Also restore MP? "HP và KI"
	cs.RestoreMP(player, healVal)

	fmt.Printf("[COMBAT] %s used %s (Heal) - Restored %d HP/MP\n",
		player.Name, skill.Template.Name, healVal)

	return nil
}

// CheckDeath kiểm tra xem entity có chết không.
func (cs *CombatService) CheckDeath(hp int) bool {
	return hp <= 0
}

// Respawn hồi sinh player.
func (cs *CombatService) Respawn(player *domain.Player) {
	player.HP = player.MaxHP
	player.MP = player.MaxMP
	// Reset to spawn point
	player.X = 100
	player.Y = 100
	fmt.Printf("[COMBAT] %s respawned at spawn point\n", player.Name)
}

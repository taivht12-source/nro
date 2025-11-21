package services

import (
	"fmt"
	"nro/src/internal/core/domain"
	"sync"
	"time"
)

// EffectService quản lý các hiệu ứng (Buffs/Debuffs) trên nhân vật.
//
// EXPLANATION:
// EffectService chịu trách nhiệm:
// 1. Quản lý vòng đời của hiệu ứng: Thêm, xóa, và tự động hết hạn (expire) các hiệu ứng.
// 2. Xử lý logic cộng dồn (stacking) hoặc làm mới (refresh) khi nhận hiệu ứng trùng loại.
// 3. Cung cấp API để kiểm tra xem nhân vật có đang chịu ảnh hưởng của một hiệu ứng cụ thể không (ví dụ: Choáng, Tàng hình).
type EffectService struct {
	mu sync.RWMutex
}

var effectServiceInstance *EffectService
var effectOnce sync.Once

func GetEffectService() *EffectService {
	effectOnce.Do(func() {
		effectServiceInstance = &EffectService{}
	})
	return effectServiceInstance
}

// AddEffect adds an effect to the player.
func (s *EffectService) AddEffect(player *domain.Player, effect *domain.Effect) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if effect of same type exists
	for _, e := range player.Effects {
		if e.Type == effect.Type {
			// Update existing effect (refresh duration or stack?)
			// For now, refresh duration and update value if higher
			e.StartTime = time.Now().UnixMilli()
			e.Duration = effect.Duration
			if effect.Value > e.Value {
				e.Value = effect.Value
			}
			fmt.Printf("[EFFECT] Refreshed effect %d for %s\n", effect.Type, player.Name)
			return
		}
	}

	effect.StartTime = time.Now().UnixMilli()
	player.Effects = append(player.Effects, effect)
	fmt.Printf("[EFFECT] Added effect %d to %s (Duration: %dms)\n", effect.Type, player.Name, effect.Duration)
}

// RemoveEffect removes an effect by type.
func (s *EffectService) RemoveEffect(player *domain.Player, effectType int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, e := range player.Effects {
		if e.Type == effectType {
			player.Effects = append(player.Effects[:i], player.Effects[i+1:]...)
			fmt.Printf("[EFFECT] Removed effect %d from %s\n", effectType, player.Name)
			return
		}
	}
}

// UpdateEffects checks for expired effects and removes them.
func (s *EffectService) UpdateEffects(player *domain.Player) {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentTime := time.Now().UnixMilli()
	var activeEffects []*domain.Effect

	for _, e := range player.Effects {
		if !e.IsExpired(currentTime) {
			activeEffects = append(activeEffects, e)
		} else {
			fmt.Printf("[EFFECT] Effect %d expired for %s\n", e.Type, player.Name)
		}
	}

	player.Effects = activeEffects
}

// HasEffect checks if player has a specific effect.
func (s *EffectService) HasEffect(player *domain.Player, effectType int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, e := range player.Effects {
		if e.Type == effectType {
			return true
		}
	}
	return false
}

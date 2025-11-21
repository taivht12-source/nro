package services

import (
	"nro-go/internal/core/domain"
	"time"
)

// MovementValidator xử lý validation cho player movement.
//
// EXPLANATION:
// MovementValidator chịu trách nhiệm:
// 1. Kiểm tra tính hợp lệ của tọa độ di chuyển (tránh đi xuyên tường, ra khỏi map).
// 2. Phát hiện và ngăn chặn hack tốc độ (Speed Hack) bằng cách tính toán khoảng cách và thời gian di chuyển.
type MovementValidator struct {
	MaxSpeed      int // Pixel per second
	MapBoundaries map[int]*MapBounds
}

// MapBounds định nghĩa giới hạn của map.
type MapBounds struct {
	MinX int
	MaxX int
	MinY int
	MaxY int
}

var validatorInstance *MovementValidator

// GetMovementValidator trả về singleton instance.
func GetMovementValidator() *MovementValidator {
	if validatorInstance == nil {
		validatorInstance = &MovementValidator{
			MaxSpeed: 500, // 500 pixels/second (có thể điều chỉnh)
			MapBoundaries: map[int]*MapBounds{
				0: {MinX: 0, MaxX: 1000, MinY: 0, MaxY: 1000}, // Map 0 default
				// TODO: Load từ config hoặc DB
			},
		}
	}
	return validatorInstance
}

// ValidatePosition kiểm tra xem vị trí có hợp lệ không.
func (v *MovementValidator) ValidatePosition(x, y int16, mapID int) bool {
	bounds, ok := v.MapBoundaries[mapID]
	if !ok {
		// Nếu không có bounds cho map này, dùng default
		bounds = &MapBounds{MinX: 0, MaxX: 1000, MinY: 0, MaxY: 1000}
	}

	if int(x) < bounds.MinX || int(x) > bounds.MaxX {
		return false
	}
	if int(y) < bounds.MinY || int(y) > bounds.MaxY {
		return false
	}

	return true
}

// ValidateSpeed kiểm tra xem tốc độ di chuyển có hợp lệ không (anti-cheat).
func (v *MovementValidator) ValidateSpeed(player *domain.Player, newX, newY int16) bool {
	// Nếu chưa có thông tin di chuyển trước đó, chấp nhận
	if player.LastMoveTime.IsZero() {
		return true
	}

	// Tính khoảng cách di chuyển
	dx := float64(newX - player.X)
	dy := float64(newY - player.Y)
	distance := dx*dx + dy*dy // Không cần sqrt, so sánh bình phương

	// Tính thời gian đã trôi qua (seconds)
	elapsed := time.Since(player.LastMoveTime).Seconds()
	if elapsed <= 0 {
		elapsed = 0.001 // Tránh chia cho 0
	}

	// Tính tốc độ (pixels/second)
	// distance = speed^2 * time^2, so speed^2 = distance / time^2
	maxDistanceSquared := float64(v.MaxSpeed*v.MaxSpeed) * elapsed * elapsed

	// Nếu di chuyển quá nhanh, reject
	if distance > maxDistanceSquared*1.2 { // Thêm 20% tolerance
		return false
	}

	return true
}

// UpdateLastMove cập nhật thời gian di chuyển cuối cùng.
func (v *MovementValidator) UpdateLastMove(player *domain.Player, x, y int16) {
	player.X = x
	player.Y = y
	player.LastMoveTime = time.Now()
}

package domain

// Effect represents a buff or debuff applied to a player.
type Effect struct {
	Type      int   // Effect type (e.g., 0: Stun, 1: Shield, 2: PowerUp)
	Value     int   // Effect value (e.g., amount of shield, % power increase)
	StartTime int64 // Timestamp when effect started
	Duration  int   // Duration in milliseconds
}

// IsExpired checks if the effect has expired.
func (e *Effect) IsExpired(currentTime int64) bool {
	return currentTime > e.StartTime+int64(e.Duration)
}

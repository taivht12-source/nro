package domain

// BossAI defines the interface for boss-specific behavior.
type BossAI interface {
	OnSpawn(boss *Boss)
	OnUpdate(boss *Boss)
	OnAttack(boss *Boss, target *Player)
	OnDamaged(boss *Boss, damage int64, attacker *Player)
	OnDie(boss *Boss, killer *Player)
}

// BaseBossAI provides default implementations for BossAI.
type BaseBossAI struct{}

func (b *BaseBossAI) OnSpawn(boss *Boss) {
	// Default: do nothing
}

func (b *BaseBossAI) OnUpdate(boss *Boss) {
	// Default: do nothing
}

func (b *BaseBossAI) OnAttack(boss *Boss, target *Player) {
	// Default: do nothing
}

func (b *BaseBossAI) OnDamaged(boss *Boss, damage int64, attacker *Player) {
	// Default: do nothing
}

func (b *BaseBossAI) OnDie(boss *Boss, killer *Player) {
	// Default: do nothing
}

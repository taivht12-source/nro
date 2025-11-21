package services

import (
	"fmt"
	"nro/src/internal/core/domain"
	"sync"
)

// BossRegistry manages the registration and creation of Boss AIs.
type BossRegistry struct {
	factories map[string]func() domain.BossAI
	mu        sync.RWMutex
}

var (
	registryInstance *BossRegistry
	registryOnce     sync.Once
)

// GetBossRegistry returns the singleton instance of BossRegistry.
func GetBossRegistry() *BossRegistry {
	registryOnce.Do(func() {
		registryInstance = &BossRegistry{
			factories: make(map[string]func() domain.BossAI),
		}
	})
	return registryInstance
}

// Register registers a factory function for a specific boss type.
func (r *BossRegistry) Register(bossType string, factory func() domain.BossAI) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories[bossType] = factory
}

// CreateAI creates a new instance of BossAI for the given boss type.
func (r *BossRegistry) CreateAI(bossType string) (domain.BossAI, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, ok := r.factories[bossType]
	if !ok {
		return nil, fmt.Errorf("boss AI not found for type: %s", bossType)
	}
	return factory(), nil
}

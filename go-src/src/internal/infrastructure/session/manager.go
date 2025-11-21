package session

import (
	"nro/src/pkg/protocol"
	"sync"
)

// Manager quản lý các Session đang hoạt động.
type Manager struct {
	sessions map[int]*protocol.Session // Map PlayerID -> Session
	mu       sync.RWMutex
}

var instance *Manager
var once sync.Once

// GetManager trả về singleton instance của Manager.
func GetManager() *Manager {
	once.Do(func() {
		instance = &Manager{
			sessions: make(map[int]*protocol.Session),
		}
	})
	return instance
}

// Add thêm một session vào quản lý (khi đăng nhập thành công).
func (m *Manager) Add(playerID int, session *protocol.Session) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[playerID] = session
}

// Remove xóa session (khi disconnect).
func (m *Manager) Remove(playerID int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, playerID)
}

// Get lấy session theo PlayerID.
func (m *Manager) Get(playerID int) *protocol.Session {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessions[playerID]
}

// Broadcast gửi tin nhắn cho danh sách người chơi.
func (m *Manager) Broadcast(playerIDs []int, msg *protocol.Message) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, pid := range playerIDs {
		if session, ok := m.sessions[pid]; ok {
			session.SendMessage(msg)
		}
	}
}

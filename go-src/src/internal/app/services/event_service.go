package services

import (
	"fmt"
	"sync"
	"time"
)

// EventType defines the type of event.
type EventType int

const (
	EventShenron EventType = 1
	EventBudokai EventType = 2
)

// EventService manages global events.
type EventService struct {
	activeEvents map[EventType]bool
	mu           sync.RWMutex
}

var eventServiceInstance *EventService
var eventOnce sync.Once

// GetEventService returns singleton instance.
func GetEventService() *EventService {
	eventOnce.Do(func() {
		eventServiceInstance = &EventService{
			activeEvents: make(map[EventType]bool),
		}
		// Start event scheduler
		go eventServiceInstance.schedulerLoop()
	})
	return eventServiceInstance
}

// StartEvent starts an event manually.
func (s *EventService) StartEvent(eventType EventType) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.activeEvents[eventType] {
		fmt.Printf("[EVENT] Event %d is already active\n", eventType)
		return
	}

	s.activeEvents[eventType] = true
	fmt.Printf("[EVENT] Started event %d\n", eventType)

	// Event specific logic
	switch eventType {
	case EventShenron:
		s.startShenronEvent()
	case EventBudokai:
		s.startBudokaiEvent()
	}
}

// StopEvent stops an event.
func (s *EventService) StopEvent(eventType EventType) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.activeEvents[eventType] {
		return
	}

	delete(s.activeEvents, eventType)
	fmt.Printf("[EVENT] Stopped event %d\n", eventType)
}

func (s *EventService) startShenronEvent() {
	fmt.Println("[EVENT] Shenron has appeared! Gather the Dragon Balls!")
	// Logic to spawn Shenron NPC or enable wish functionality
}

func (s *EventService) startBudokaiEvent() {
	fmt.Println("[EVENT] The Budokai Tournament has begun! Register now!")
	// Logic to open registration
}

// schedulerLoop checks for scheduled events.
func (s *EventService) schedulerLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		now := time.Now()
		// Example: Start Budokai every day at 20:00
		if now.Hour() == 20 && now.Minute() == 0 {
			s.StartEvent(EventBudokai)
		}
	}
}

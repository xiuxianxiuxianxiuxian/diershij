package service

import (
	"fmt"
	"time"
)

// WorldEventType represents the type of world event.
type WorldEventType string

const (
	EventDemonBeastTide   WorldEventType = "demon_beast_tide"
	EventHeavenlyAnomaly  WorldEventType = "heavenly_anomaly"
	EventSpiritualTide    WorldEventType = "spiritual_tide"
	EventRealmCollision   WorldEventType = "realm_collision"
	EventAncientTombOpen  WorldEventType = "ancient_tomb_open"
	EventSectWar          WorldEventType = "sect_war"
)

// WorldEventPhase represents the current phase of a world event.
type WorldEventPhase string

const (
	PhaseGathering  WorldEventPhase = "gathering"  // event is forming
	PhaseActive     WorldEventPhase = "active"     // event is ongoing
	PhaseClimax     WorldEventPhase = "climax"     // peak intensity
	PhaseResolving  WorldEventPhase = "resolving"  // winding down
	PhaseEnded      WorldEventPhase = "ended"      // event has ended
)

// WorldEvent represents a world-scale event.
type WorldEvent struct {
	ID             string
	Type           WorldEventType
	Title          string
	Description    string
	Phase          WorldEventPhase
	Severity       int       // 1-10, higher = more dangerous
	StartTime      time.Time
	EndTime        time.Time
	AffectedRegions []string
	Participants   []string // entity IDs who participated
	RewardsDistributed bool
}

// WorldEventService manages world events lifecycle.
type WorldEventService struct {
	events      map[string]*WorldEvent
	activeEvents []string
	eventCounter int
}

// NewWorldEventService creates a new WorldEventService.
func NewWorldEventService() *WorldEventService {
	return &WorldEventService{
		events:      make(map[string]*WorldEvent),
		activeEvents: []string{},
	}
}

// GenerateEventInput holds inputs for generating a world event.
type GenerateEventInput struct {
	Type           WorldEventType
	Title          string
	Description    string
	Severity       int
	DurationHours  int
	AffectedRegions []string
}

// GenerateEvent creates a new world event.
func (s *WorldEventService) GenerateEvent(input GenerateEventInput, now time.Time) (*WorldEvent, error) {
	if input.Severity < 1 || input.Severity > 10 {
		return nil, fmt.Errorf("严重程度必须在 1-10 之间")
	}

	s.eventCounter++
	id := fmt.Sprintf("event_%d", s.eventCounter)

	event := &WorldEvent{
		ID:              id,
		Type:            input.Type,
		Title:           input.Title,
		Description:     input.Description,
		Phase:           PhaseGathering,
		Severity:        input.Severity,
		StartTime:       now,
		EndTime:         now.Add(time.Duration(input.DurationHours) * time.Hour),
		AffectedRegions: input.AffectedRegions,
	}

	s.events[id] = event
	s.activeEvents = append(s.activeEvents, id)

	return event, nil
}

// AdvanceEvent updates the phase of an event based on elapsed time.
func (s *WorldEventService) AdvanceEvent(eventID string, now time.Time) (*WorldEvent, error) {
	event, exists := s.events[eventID]
	if !exists {
		return nil, fmt.Errorf("事件 '%s' 不存在", eventID)
	}

	if event.Phase == PhaseEnded {
		return event, nil
	}

	elapsed := now.Sub(event.StartTime)
	duration := event.EndTime.Sub(event.StartTime)

	if elapsed <= 0 {
		event.Phase = PhaseGathering
	} else if elapsed < duration/4 {
		event.Phase = PhaseActive
	} else if elapsed < (duration*3)/4 {
		event.Phase = PhaseClimax
	} else if elapsed < duration {
		event.Phase = PhaseResolving
	} else {
		event.Phase = PhaseEnded
	}

	return event, nil
}

// RegisterParticipant records an entity's participation in an event.
func (s *WorldEventService) RegisterParticipant(eventID, entityID string) error {
	event, exists := s.events[eventID]
	if !exists {
		return fmt.Errorf("事件 '%s' 不存在", eventID)
	}

	if event.Phase == PhaseEnded {
		return fmt.Errorf("事件已结束")
	}

	// Check if already registered
	for _, p := range event.Participants {
		if p == entityID {
			return nil // already registered
		}
	}

	event.Participants = append(event.Participants, entityID)
	return nil
}

// CalculateRewards computes rewards for event participants.
func (s *WorldEventService) CalculateRewards(eventID string) (map[string]int64, error) {
	event, exists := s.events[eventID]
	if !exists {
		return nil, fmt.Errorf("事件 '%s' 不存在", eventID)
	}

	if event.Phase != PhaseResolving && event.Phase != PhaseEnded {
		return nil, fmt.Errorf("事件未结束，无法发放奖励")
	}

	if event.RewardsDistributed {
		return nil, fmt.Errorf("奖励已发放")
	}

	rewards := make(map[string]int64)
	baseReward := int64(event.Severity * 10)

	for _, entityID := range event.Participants {
		// Bonus for higher severity events
		rewards[entityID] = baseReward
	}

	event.RewardsDistributed = true
	return rewards, nil
}

// GetActiveEvents returns all currently active events.
func (s *WorldEventService) GetActiveEvents(now time.Time) []*WorldEvent {
	var active []*WorldEvent
	for _, id := range s.activeEvents {
		event := s.events[id]
		if event.Phase != PhaseEnded && now.Before(event.EndTime) {
			active = append(active, event)
		}
	}
	return active
}

// GetEvent retrieves an event by ID.
func (s *WorldEventService) GetEvent(id string) (*WorldEvent, bool) {
	e, ok := s.events[id]
	return e, ok
}

// ListEvents returns all events.
func (s *WorldEventService) ListEvents() []*WorldEvent {
	events := make([]*WorldEvent, 0, len(s.events))
	for _, e := range s.events {
		events = append(events, e)
	}
	return events
}

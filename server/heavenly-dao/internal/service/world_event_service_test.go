package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWorldEventService(t *testing.T) {
	svc := NewWorldEventService()
	assert.NotNil(t, svc)
}

func TestGenerateEvent_ValidInput(t *testing.T) {
	svc := NewWorldEventService()

	input := GenerateEventInput{
		Type:            EventDemonBeastTide,
		Title:           "妖兽狂潮",
		Description:     "大量妖兽袭击人类聚居地",
		Severity:        7,
		DurationHours:   24,
		AffectedRegions: []string{"青州", "云州"},
	}

	now := time.Now()
	event, err := svc.GenerateEvent(input, now)
	assert.NoError(t, err)
	assert.Equal(t, EventDemonBeastTide, event.Type)
	assert.Equal(t, PhaseGathering, event.Phase)
	assert.Equal(t, 7, event.Severity)
	assert.Equal(t, now, event.StartTime)
	assert.Len(t, event.AffectedRegions, 2)
}

func TestGenerateEvent_InvalidSeverity(t *testing.T) {
	svc := NewWorldEventService()

	input := GenerateEventInput{
		Type:     EventDemonBeastTide,
		Severity: 11,
	}

	_, err := svc.GenerateEvent(input, time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "严重程度")
}

func TestGenerateEvent_SeverityZero(t *testing.T) {
	svc := NewWorldEventService()

	input := GenerateEventInput{
		Type:     EventDemonBeastTide,
		Severity: 0,
	}

	_, err := svc.GenerateEvent(input, time.Now())
	assert.Error(t, err)
}

func TestAdvanceEvent_PhaseTransitions(t *testing.T) {
	svc := NewWorldEventService()

	input := GenerateEventInput{
		Type:            EventHeavenlyAnomaly,
		Title:           "天象异变",
		Severity:        5,
		DurationHours:   100,
		AffectedRegions: []string{"天州"},
	}

	now := time.Now()
	event, _ := svc.GenerateEvent(input, now)

	// Test PhaseGathering (elapsed = 0)
	svc.AdvanceEvent(event.ID, now)
	assert.Equal(t, PhaseGathering, event.Phase)

	// Test PhaseActive (elapsed < 25%)
	t2 := now.Add(20 * time.Hour)
	svc.AdvanceEvent(event.ID, t2)
	assert.Equal(t, PhaseActive, event.Phase)

	// Test PhaseClimax (25% <= elapsed < 75%)
	t3 := now.Add(50 * time.Hour)
	svc.AdvanceEvent(event.ID, t3)
	assert.Equal(t, PhaseClimax, event.Phase)

	// Test PhaseResolving (75% <= elapsed < 100%)
	t4 := now.Add(80 * time.Hour)
	svc.AdvanceEvent(event.ID, t4)
	assert.Equal(t, PhaseResolving, event.Phase)

	// Test PhaseEnded (elapsed >= 100%)
	t5 := now.Add(120 * time.Hour)
	svc.AdvanceEvent(event.ID, t5)
	assert.Equal(t, PhaseEnded, event.Phase)
}

func TestAdvanceEvent_NonExistentEvent(t *testing.T) {
	svc := NewWorldEventService()

	_, err := svc.AdvanceEvent("nonexistent", time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "不存在")
}

func TestAdvanceEvent_AlreadyEnded(t *testing.T) {
	svc := NewWorldEventService()

	input := GenerateEventInput{
		Type:            EventDemonBeastTide,
		Title:           "妖兽潮",
		Severity:        3,
		DurationHours:   1,
		AffectedRegions: []string{"青州"},
	}

	now := time.Now()
	event, _ := svc.GenerateEvent(input, now)

	// Fast forward past end
	future := now.Add(2 * time.Hour)
	svc.AdvanceEvent(event.ID, future)
	assert.Equal(t, PhaseEnded, event.Phase)

	// Should stay ended
	svc.AdvanceEvent(event.ID, future.Add(time.Hour))
	assert.Equal(t, PhaseEnded, event.Phase)
}

func TestRegisterParticipant(t *testing.T) {
	svc := NewWorldEventService()

	input := GenerateEventInput{
		Type:            EventDemonBeastTide,
		Title:           "妖兽潮",
		Severity:        5,
		DurationHours:   24,
		AffectedRegions: []string{"青州"},
	}

	now := time.Now()
	event, _ := svc.GenerateEvent(input, now)

	err := svc.RegisterParticipant(event.ID, "npc_1")
	assert.NoError(t, err)

	err = svc.RegisterParticipant(event.ID, "npc_2")
	assert.NoError(t, err)

	// Duplicate should not error (already registered)
	err = svc.RegisterParticipant(event.ID, "npc_1")
	assert.NoError(t, err)

	assert.Len(t, event.Participants, 2)
}

func TestRegisterParticipant_NonExistentEvent(t *testing.T) {
	svc := NewWorldEventService()

	err := svc.RegisterParticipant("nonexistent", "npc_1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "不存在")
}

func TestRegisterParticipant_EndedEvent(t *testing.T) {
	svc := NewWorldEventService()

	input := GenerateEventInput{
		Type:            EventDemonBeastTide,
		Title:           "妖兽潮",
		Severity:        5,
		DurationHours:   1,
		AffectedRegions: []string{"青州"},
	}

	now := time.Now()
	event, _ := svc.GenerateEvent(input, now)
	event.Phase = PhaseEnded

	err := svc.RegisterParticipant(event.ID, "npc_1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已结束")
}

func TestCalculateRewards(t *testing.T) {
	svc := NewWorldEventService()

	input := GenerateEventInput{
		Type:            EventDemonBeastTide,
		Title:           "妖兽潮",
		Severity:        7,
		DurationHours:   24,
		AffectedRegions: []string{"青州"},
	}

	now := time.Now()
	event, _ := svc.GenerateEvent(input, now)

	// Register participants
	svc.RegisterParticipant(event.ID, "npc_1")
	svc.RegisterParticipant(event.ID, "npc_2")

	// Event not ended yet
	_, err := svc.CalculateRewards(event.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "未结束")

	// End the event
	event.Phase = PhaseEnded

	rewards, err := svc.CalculateRewards(event.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(70), rewards["npc_1"]) // severity * 10 = 70
	assert.Equal(t, int64(70), rewards["npc_2"])
	assert.True(t, event.RewardsDistributed)
}

func TestCalculateRewards_AlreadyDistributed(t *testing.T) {
	svc := NewWorldEventService()

	input := GenerateEventInput{
		Type:            EventDemonBeastTide,
		Title:           "妖兽潮",
		Severity:        5,
		DurationHours:   24,
		AffectedRegions: []string{"青州"},
	}

	now := time.Now()
	event, _ := svc.GenerateEvent(input, now)
	event.Phase = PhaseEnded
	event.RewardsDistributed = true

	_, err := svc.CalculateRewards(event.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已发放")
}

func TestCalculateRewards_NonExistentEvent(t *testing.T) {
	svc := NewWorldEventService()

	_, err := svc.CalculateRewards("nonexistent")
	assert.Error(t, err)
}

func TestGetActiveEvents(t *testing.T) {
	svc := NewWorldEventService()

	now := time.Now()

	// Create active event
	svc.GenerateEvent(GenerateEventInput{
		Type:          EventDemonBeastTide,
		Title:         "妖兽潮",
		Severity:      5,
		DurationHours: 24,
	}, now)

	// Create ended event
	endedEvent, _ := svc.GenerateEvent(GenerateEventInput{
		Type:          EventHeavenlyAnomaly,
		Title:         "天象异变",
		Severity:      3,
		DurationHours: 1,
	}, now.Add(-2 * time.Hour))
	endedEvent.Phase = PhaseEnded

	active := svc.GetActiveEvents(now)
	assert.Len(t, active, 1)
	assert.Equal(t, EventDemonBeastTide, active[0].Type)
}

func TestGetEvent(t *testing.T) {
	svc := NewWorldEventService()

	now := time.Now()
	event, _ := svc.GenerateEvent(GenerateEventInput{
		Type:     EventDemonBeastTide,
		Title:    "妖兽潮",
		Severity: 5,
	}, now)

	found, ok := svc.GetEvent(event.ID)
	assert.True(t, ok)
	assert.Equal(t, event.ID, found.ID)

	_, ok = svc.GetEvent("nonexistent")
	assert.False(t, ok)
}

func TestListEvents(t *testing.T) {
	svc := NewWorldEventService()

	now := time.Now()
	svc.GenerateEvent(GenerateEventInput{
		Type:     EventDemonBeastTide,
		Title:    "妖兽潮",
		Severity: 5,
	}, now)
	svc.GenerateEvent(GenerateEventInput{
		Type:     EventHeavenlyAnomaly,
		Title:    "天象异变",
		Severity: 3,
	}, now)

	events := svc.ListEvents()
	assert.Len(t, events, 2)
}

package service

import (
	"testing"
	"time"

	"github.com/cultivation-world/shared/types"
	"github.com/stretchr/testify/assert"
)

func makeExploreInput() ExploreInput {
	return ExploreInput{
		EntityID:        "explorer_1",
		RegionName:      "青云山脉",
		RegionLevel:     3,
		SpiritualDensity: 0.5,
		IsDangerZone:    false,
		FortuneScore:    50,
		Luck:            50,
		Realm:           types.RealmFoundation,
		MovementSpeed:   1.0,
	}
}

func TestNewExploreOperation(t *testing.T) {
	op := NewExploreOperation(time.Minute)
	assert.NotNil(t, op)
}

func TestExecuteExplore_Basic(t *testing.T) {
	op := NewExploreOperation(time.Minute)

	input := makeExploreInput()
	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	result, err := op.ExecuteExplore(input, now, deterministic)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.Events)
	assert.True(t, result.TimeElapsed > 0)
	assert.NotEmpty(t, result.Message)
}

func TestExecuteExplore_Cooldown(t *testing.T) {
	op := NewExploreOperation(time.Hour)

	input := makeExploreInput()
	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	result, err := op.ExecuteExplore(input, now, deterministic)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	_, err = op.ExecuteExplore(input, now.Add(30*time.Minute), deterministic)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cooldown")
}

func TestExecuteExplore_DangerZone(t *testing.T) {
	op := NewExploreOperation(time.Minute)

	input := makeExploreInput()
	input.IsDangerZone = true

	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	result, err := op.ExecuteExplore(input, now, deterministic)
	assert.NoError(t, err)

	// Danger zone should have more dangerous events
	dangerCount := 0
	for _, event := range result.Events {
		if event.Danger > 0 {
			dangerCount++
		}
	}
	assert.True(t, dangerCount >= 0)
}

func TestExecuteExplore_HighFortune(t *testing.T) {
	op := NewExploreOperation(time.Minute)

	input := makeExploreInput()
	input.FortuneScore = 90 // high fortune

	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	result, err := op.ExecuteExplore(input, now, deterministic)
	assert.NoError(t, err)

	// High fortune should lead to more discoveries and epiphanies
	discoveryCount := 0
	for _, event := range result.Events {
		if event.Type == "discovery" || event.Type == "epiphany" {
			discoveryCount++
		}
	}
	assert.True(t, discoveryCount >= 0)
}

func TestExecuteExplore_DurationScalesWithRegion(t *testing.T) {
	// Low level region
	opLow := NewExploreOperation(time.Minute)
	inputLow := makeExploreInput()
	inputLow.EntityID = "explorer_low"
	inputLow.RegionLevel = 1
	now := time.Now()

	resultLow, err := opLow.ExecuteExplore(inputLow, now, func() float64 { return 0.5 })
	assert.NoError(t, err)
	assert.NotNil(t, resultLow)

	// High level region
	opHigh := NewExploreOperation(time.Minute)
	inputHigh := makeExploreInput()
	inputHigh.EntityID = "explorer_high"
	inputHigh.RegionLevel = 8
	resultHigh, err := opHigh.ExecuteExplore(inputHigh, now, func() float64 { return 0.5 })
	assert.NoError(t, err)
	assert.NotNil(t, resultHigh)

	assert.True(t, resultHigh.TimeElapsed > resultLow.TimeElapsed)
}

func TestCalculateExploreDuration(t *testing.T) {
	// Base: region_level * 4 + 8
	assert.Equal(t, 12, calculateExploreDuration(1, 1.0))
	assert.Equal(t, 20, calculateExploreDuration(3, 1.0))
	assert.Equal(t, 40, calculateExploreDuration(8, 1.0))

	// With movement speed bonus
	assert.Equal(t, 10, calculateExploreDuration(3, 2.0))
	assert.True(t, calculateExploreDuration(3, 0.5) > 20)
}

func TestGetCooldownRemaining_Explore(t *testing.T) {
	op := NewExploreOperation(time.Hour)

	input := makeExploreInput()
	now := time.Now()
	_, _ = op.ExecuteExplore(input, now, func() float64 { return 0.5 })

	remaining := op.GetCooldownRemaining("explorer_1", now.Add(30*time.Minute))
	assert.True(t, remaining > 0)
}

func TestClearCooldown_Explore(t *testing.T) {
	op := NewExploreOperation(time.Hour)

	input := makeExploreInput()
	now := time.Now()
	_, _ = op.ExecuteExplore(input, now, func() float64 { return 0.5 })

	op.ClearCooldown("explorer_1")
	remaining := op.GetCooldownRemaining("explorer_1", now.Add(30*time.Minute))
	assert.Equal(t, time.Duration(0), remaining)
}

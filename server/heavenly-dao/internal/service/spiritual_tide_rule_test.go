package service

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSpiritualTideRule(t *testing.T) {
	rule := NewSpiritualTideRule()
	assert.NotNil(t, rule)
	assert.Equal(t, 1.0, rule.baseDensity)
	assert.Equal(t, 0.30, rule.tideAmplitude)
	assert.Equal(t, 30, rule.tideCycleDays)
}

func TestCalculateCurrentTide_Base(t *testing.T) {
	rule := NewSpiritualTideRule()

	// Day 0: sin(0) = 0, so density = base = 1.0
	result := rule.CalculateCurrentTide(0, 1.0)
	assert.InDelta(t, 1.0, result.Density, 0.01)
	assert.Equal(t, TidePhaseRising, result.Phase)
}

func TestCalculateCurrentTide_HighTide(t *testing.T) {
	rule := NewSpiritualTideRule()

	// Day 7.5: sin(2π*7.5/30) = sin(π/2) = 1
	// density = 1.0 * (1 + 0.30) = 1.30
	result := rule.CalculateCurrentTide(7, 1.0)
	// Phase at π/2 is High (angle >= π/2 and < π)
	assert.True(t, result.Density > 1.2)
	assert.True(t, result.Density <= 1.30)
}

func TestCalculateCurrentTide_LowTide(t *testing.T) {
	rule := NewSpiritualTideRule()

	// Day 22.5: sin(2π*22.5/30) = sin(3π/2) = -1
	// density = 1.0 * (1 - 0.30) = 0.70
	result := rule.CalculateCurrentTide(22, 1.0)
	assert.True(t, result.Density < 0.8)
	assert.True(t, result.Density >= 0.70)
}

func TestCalculateCurrentTide_RegionBonus(t *testing.T) {
	rule := NewSpiritualTideRule()

	result := rule.CalculateCurrentTide(0, 2.0)
	// density = 1.0 * 2.0 * (1 + 0) = 2.0
	assert.InDelta(t, 2.0, result.Density, 0.01)
}

func TestCalculateCurrentTide_PhaseTransitions(t *testing.T) {
	rule := NewSpiritualTideRule()

	// Test all 4 phases
	phases := make(map[TidePhase]bool)
	for day := 0; day < 30; day++ {
		result := rule.CalculateCurrentTide(day, 1.0)
		phases[result.Phase] = true
	}

	assert.True(t, phases[TidePhaseRising])
	assert.True(t, phases[TidePhaseHigh])
	assert.True(t, phases[TidePhaseFalling])
	assert.True(t, phases[TidePhaseLow])
}

func TestCalculateCurrentTide_DaysUntilShift(t *testing.T) {
	rule := NewSpiritualTideRule()

	result := rule.CalculateCurrentTide(0, 1.0)
	assert.True(t, result.DaysUntilShift > 0)
	assert.True(t, result.DaysUntilShift <= 30)
}

func TestAdjustSpiritualDensity_Positive(t *testing.T) {
	rule := NewSpiritualTideRule()

	density := rule.AdjustSpiritualDensity(1.0, []float64{0.1, 0.2})
	// 1.0 * 1.1 * 1.2 = 1.32
	assert.InDelta(t, 1.32, density, 0.001)
}

func TestAdjustSpiritualDensity_Negative(t *testing.T) {
	rule := NewSpiritualTideRule()

	density := rule.AdjustSpiritualDensity(1.0, []float64{-0.3})
	// 1.0 * 0.7 = 0.7
	assert.InDelta(t, 0.7, density, 0.001)
}

func TestAdjustSpiritualDensity_Multiple(t *testing.T) {
	rule := NewSpiritualTideRule()

	density := rule.AdjustSpiritualDensity(2.0, []float64{0.5, -0.2, 0.1})
	// 2.0 * 1.5 * 0.8 * 1.1 = 2.64
	assert.InDelta(t, 2.64, density, 0.001)
}

func TestAdjustSpiritualDensity_MinClamp(t *testing.T) {
	rule := NewSpiritualTideRule()

	density := rule.AdjustSpiritualDensity(0.01, []float64{-0.99})
	// 0.01 * 0.01 = 0.0001 → clamped to 0.01
	assert.InDelta(t, 0.01, density, 0.001)
}

func TestDeterminePhase_AllQuadrants(t *testing.T) {
	rule := NewSpiritualTideRule()

	tests := []struct {
		angle    float64
		expected TidePhase
	}{
		{0, TidePhaseRising},
		{math.Pi / 4, TidePhaseRising},
		{math.Pi / 2, TidePhaseHigh},
		{3 * math.Pi / 4, TidePhaseHigh},
		{math.Pi, TidePhaseFalling},
		{5 * math.Pi / 4, TidePhaseFalling},
		{3 * math.Pi / 2, TidePhaseLow},
		{7 * math.Pi / 4, TidePhaseLow},
	}

	for _, tt := range tests {
		phase := rule.determinePhase(tt.angle)
		assert.Equal(t, tt.expected, phase, "angle=%v", tt.angle)
	}
}

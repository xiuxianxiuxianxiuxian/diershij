package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultResourceSpawnRule(t *testing.T) {
	rule := DefaultResourceSpawnRule()
	assert.NotNil(t, rule)
	assert.Equal(t, 10.0, rule.BaseSpawnRate)
	assert.Equal(t, 1.5, rule.SpiritualMultiplier)
	assert.Equal(t, 0.5, rule.PressurePenalty)
	assert.Equal(t, 0.7, rule.BalanceThreshold)
}

func TestCalculateSpawnRate_BaseCase(t *testing.T) {
	rule := DefaultResourceSpawnRule()

	input := ResourceSpawnInput{
		SpiritualDensity: 0.5,
		SpiritualPressure: 0.3,
		BalanceFactor:    0.2,
		BaseRate:         10.0,
	}

	rate := rule.CalculateSpawnRate(input)
	assert.Greater(t, rate, 0.0)
	// Spiritual boost: 1 + 0.5*1.5 = 1.75
	// Pressure factor: 1 - 0.3*0.5 = 0.85
	// Balance: 1.0 (below threshold)
	// Expected: 10 * 1.75 * 0.85 * 1.0 = 14.875
	assert.InDelta(t, 14.875, rate, 0.01)
}

func TestCalculateSpawnRate_HighSpiritual(t *testing.T) {
	rule := DefaultResourceSpawnRule()

	input := ResourceSpawnInput{
		SpiritualDensity: 1.0,
		SpiritualPressure: 0.0,
		BalanceFactor:    0.0,
		BaseRate:         10.0,
	}

	rate := rule.CalculateSpawnRate(input)
	// 10 * (1 + 1.0*1.5) * 1.0 * 1.0 = 10 * 2.5 = 25
	assert.InDelta(t, 25.0, rate, 0.01)
}

func TestCalculateSpawnRate_HighPressure(t *testing.T) {
	rule := DefaultResourceSpawnRule()

	input := ResourceSpawnInput{
		SpiritualDensity: 0.0,
		SpiritualPressure: 1.0,
		BalanceFactor:    0.0,
		BaseRate:         10.0,
	}

	rate := rule.CalculateSpawnRate(input)
	// 10 * 1.0 * (1 - 1.0*0.5) * 1.0 = 10 * 0.5 = 5
	assert.InDelta(t, 5.0, rate, 0.01)
}

func TestCalculateSpawnRate_MaxPressure(t *testing.T) {
	rule := DefaultResourceSpawnRule()

	input := ResourceSpawnInput{
		SpiritualDensity: 0.0,
		SpiritualPressure: 1.0,
		BalanceFactor:    0.0,
		BaseRate:         10.0,
	}

	rule.PressurePenalty = 2.0 // would go negative
	rate := rule.CalculateSpawnRate(input)
	// Should be clamped to minimum 10%
	assert.GreaterOrEqual(t, rate, 1.0)
}

func TestCalculateSpawnRate_BalancePenalty(t *testing.T) {
	rule := DefaultResourceSpawnRule()

	input := ResourceSpawnInput{
		SpiritualDensity: 0.5,
		SpiritualPressure: 0.0,
		BalanceFactor:    0.9, // severely unbalanced
		BaseRate:         10.0,
	}

	rate := rule.CalculateSpawnRate(input)
	// Balance penalty: excess = 0.9 - 0.7 = 0.2
	// penalty = 1 - 0.2^2 * 2 = 1 - 0.08 = 0.92
	// Expected: 10 * 1.75 * 1.0 * 0.92 = 16.1
	assert.InDelta(t, 16.1, rate, 0.1)
}

func TestCalculateSpawnRate_SevereImbalance(t *testing.T) {
	rule := DefaultResourceSpawnRule()

	input := ResourceSpawnInput{
		SpiritualDensity: 0.5,
		SpiritualPressure: 0.0,
		BalanceFactor:    1.0, // maximum imbalance
		BaseRate:         10.0,
	}

	rate := rule.CalculateSpawnRate(input)
	// Should be clamped to minimum 20% balance penalty
	assert.GreaterOrEqual(t, rate, 10.0*1.75*1.0*0.2)
}

func TestCalculateSpawnRate_MinimumRate(t *testing.T) {
	rule := DefaultResourceSpawnRule()

	input := ResourceSpawnInput{
		SpiritualDensity: 0.0,
		SpiritualPressure: 1.0,
		BalanceFactor:    1.0,
		BaseRate:         1.0,
	}

	rate := rule.CalculateSpawnRate(input)
	// Should never go below 0.1
	assert.GreaterOrEqual(t, rate, 0.1)
}

func TestCalculateSpawnRate_ZeroBaseRate(t *testing.T) {
	rule := DefaultResourceSpawnRule()

	input := ResourceSpawnInput{
		SpiritualDensity: 0.5,
		SpiritualPressure: 0.0,
		BalanceFactor:    0.0,
		BaseRate:         0, // should use default
	}

	rate := rule.CalculateSpawnRate(input)
	// Should use default base rate of 10
	assert.Greater(t, rate, 10.0)
}

func TestCalculateRareSpawnChance_BaseCase(t *testing.T) {
	rule := DefaultResourceSpawnRule()

	input := ResourceSpawnInput{
		SpiritualDensity: 0.5,
		SpiritualPressure: 0.5,
		BalanceFactor:    0.5,
	}

	chance := rule.CalculateRareSpawnChance(input)
	// 0.05 + 0.5*0.15 + (1-0.5)*0.05 + 0 = 0.05 + 0.075 + 0.025 = 0.15
	assert.InDelta(t, 0.15, chance, 0.01)
}

func TestCalculateRareSpawnChance_HighSpiritual(t *testing.T) {
	rule := DefaultResourceSpawnRule()

	input := ResourceSpawnInput{
		SpiritualDensity: 1.0,
		SpiritualPressure: 0.0,
		BalanceFactor:    0.0,
	}

	chance := rule.CalculateRareSpawnChance(input)
	// 0.05 + 1.0*0.15 + 1.0*0.05 + 0.05 (balanced bonus) = 0.30
	assert.InDelta(t, 0.30, chance, 0.01)
}

func TestCalculateRareSpawnChance_MaxChance(t *testing.T) {
	rule := DefaultResourceSpawnRule()

	input := ResourceSpawnInput{
		SpiritualDensity: 1.0,
		SpiritualPressure: 0.0,
		BalanceFactor:    0.0,
		BaseRate:         1000, // shouldn't matter
	}

	chance := rule.CalculateRareSpawnChance(input)
	// Should be capped at 50%
	assert.LessOrEqual(t, chance, 0.5)
}

func TestValidateSpawnInput_Valid(t *testing.T) {
	rule := DefaultResourceSpawnRule()

	input := ResourceSpawnInput{
		SpiritualDensity: 0.5,
		SpiritualPressure: 0.3,
		BalanceFactor:    0.7,
	}

	err := rule.ValidateSpawnInput(input)
	assert.NoError(t, err)
}

func TestValidateSpawnInput_InvalidDensity(t *testing.T) {
	rule := DefaultResourceSpawnRule()

	input := ResourceSpawnInput{
		SpiritualDensity: 1.5,
	}

	err := rule.ValidateSpawnInput(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "spiritual density")
}

func TestValidateSpawnInput_InvalidPressure(t *testing.T) {
	rule := DefaultResourceSpawnRule()

	input := ResourceSpawnInput{
		SpiritualDensity: 0.5,
		SpiritualPressure: -0.1,
	}

	err := rule.ValidateSpawnInput(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "spiritual pressure")
}

func TestValidateSpawnInput_InvalidBalance(t *testing.T) {
	rule := DefaultResourceSpawnRule()

	input := ResourceSpawnInput{
		SpiritualDensity: 0.5,
		SpiritualPressure: 0.5,
		BalanceFactor:    1.5,
	}

	err := rule.ValidateSpawnInput(input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "balance factor")
}

func TestDefaultRefreshCycleConfig(t *testing.T) {
	cfg := DefaultRefreshCycleConfig()
	assert.NotNil(t, cfg)
	assert.Equal(t, 24, cfg.RefreshIntervalHours)
	assert.Equal(t, 100, cfg.MaxResourcesPerRegion)
	assert.InDelta(t, 0.3, cfg.RefreshPercentage, 0.01)
}

func TestCalculateRefreshAmount(t *testing.T) {
	cfg := DefaultRefreshCycleConfig()

	// Normal case: 50 current, 100 max
	refresh := cfg.CalculateRefreshAmount(50, 100)
	// Deficit = 50, refresh = 50 * 0.3 = 15
	assert.Equal(t, 15, refresh)
}

func TestCalculateRefreshAmount_Full(t *testing.T) {
	cfg := DefaultRefreshCycleConfig()

	refresh := cfg.CalculateRefreshAmount(100, 100)
	assert.Equal(t, 0, refresh)
}

func TestCalculateRefreshAmount_ExceedsMax(t *testing.T) {
	cfg := DefaultRefreshCycleConfig()

	// 95 current, 100 max: deficit = 5, refresh = 5*0.3 = 1.5 -> 1
	refresh := cfg.CalculateRefreshAmount(95, 100)
	assert.Equal(t, 1, refresh)
	// 95 + 1 = 96, still under max
}

func TestCalculateRefreshAmount_MinimumRefresh(t *testing.T) {
	cfg := DefaultRefreshCycleConfig()

	// 99 current, 100 max: deficit = 1, refresh = 1*0.3 = 0.3 -> 0 -> clamp to 1
	refresh := cfg.CalculateRefreshAmount(99, 100)
	assert.Equal(t, 1, refresh)
}

func TestCalculateRefreshAmount_LargeDeficit(t *testing.T) {
	cfg := DefaultRefreshCycleConfig()

	// 0 current, 100 max: deficit = 100, refresh = 30
	refresh := cfg.CalculateRefreshAmount(0, 100)
	assert.Equal(t, 30, refresh)
}

func TestShouldRefresh(t *testing.T) {
	cfg := DefaultRefreshCycleConfig()

	// Not enough time passed
	assert.False(t, cfg.ShouldRefresh(20, 23))

	// Exactly enough time
	assert.True(t, cfg.ShouldRefresh(0, 24))

	// More than enough time
	assert.True(t, cfg.ShouldRefresh(0, 48))
}

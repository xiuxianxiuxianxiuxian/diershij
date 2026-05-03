package service

import (
	"fmt"
	"math"
)

// ResourceSpawnRule defines how resources spawn in the world.
type ResourceSpawnRule struct {
	// Base spawn rate per region (resources per day)
	BaseSpawnRate float64

	// Spiritual density multiplier (higher spiritual = more resources)
	SpiritualMultiplier float64

	// Pressure penalty (too much pressure reduces spawn)
	PressurePenalty float64

	// Balance factor (0 = perfectly balanced, 1 = completely unbalanced)
	BalanceThreshold float64
}

// DefaultResourceSpawnRule returns the default resource spawn rule.
func DefaultResourceSpawnRule() *ResourceSpawnRule {
	return &ResourceSpawnRule{
		BaseSpawnRate:       10.0,
		SpiritualMultiplier: 1.5,
		PressurePenalty:     0.5,
		BalanceThreshold:    0.7,
	}
}

// ResourceSpawnInput holds inputs for calculating spawn rates.
type ResourceSpawnInput struct {
	SpiritualDensity float64 // 0.0 - 1.0
	SpiritualPressure float64 // 0.0 - 1.0
	BalanceFactor    float64 // 0.0 - 1.0 (Gini-like)
	BaseRate         float64
}

// CalculateSpawnRate calculates the resource spawn rate based on world conditions.
func (r *ResourceSpawnRule) CalculateSpawnRate(input ResourceSpawnInput) float64 {
	if input.BaseRate <= 0 {
		input.BaseRate = r.BaseSpawnRate
	}

	// Spiritual density boost: more spiritual = more resources
	spiritualBoost := 1.0 + input.SpiritualDensity*r.SpiritualMultiplier

	// Pressure penalty: high pressure reduces spawn
	pressureFactor := 1.0 - input.SpiritualPressure*r.PressurePenalty
	if pressureFactor < 0.1 {
		pressureFactor = 0.1 // minimum 10%
	}

	// Balance penalty: unbalanced worlds have reduced spawn
	balancePenalty := 1.0
	if input.BalanceFactor > r.BalanceThreshold {
		// Exponential penalty for severe imbalance
		excess := input.BalanceFactor - r.BalanceThreshold
		balancePenalty = 1.0 - math.Pow(excess, 2)*2
		if balancePenalty < 0.2 {
			balancePenalty = 0.2
		}
	}

	rate := input.BaseRate * spiritualBoost * pressureFactor * balancePenalty
	return math.Max(0.1, rate) // minimum 0.1 resources per day
}

// CalculateRareSpawnChance calculates the chance of rare resource spawning.
func (r *ResourceSpawnRule) CalculateRareSpawnChance(input ResourceSpawnInput) float64 {
	// Base rare chance: 5%
	baseChance := 0.05

	// High spiritual density increases rare spawns
	spiritualBonus := input.SpiritualDensity * 0.15

	// Low pressure slightly increases rare spawns (cultivators seeking rare resources)
	pressureBonus := (1.0 - input.SpiritualPressure) * 0.05

	// Balanced world has slightly better rare spawn distribution
	balanceBonus := 0.0
	if input.BalanceFactor < 0.3 {
		balanceBonus = 0.05
	}

	chance := baseChance + spiritualBonus + pressureBonus + balanceBonus
	return math.Min(0.5, chance) // cap at 50%
}

// ValidateSpawnInput validates the spawn input parameters.
func (r *ResourceSpawnRule) ValidateSpawnInput(input ResourceSpawnInput) error {
	if input.SpiritualDensity < 0 || input.SpiritualDensity > 1 {
		return fmt.Errorf("spiritual density must be between 0 and 1, got %.2f", input.SpiritualDensity)
	}
	if input.SpiritualPressure < 0 || input.SpiritualPressure > 1 {
		return fmt.Errorf("spiritual pressure must be between 0 and 1, got %.2f", input.SpiritualPressure)
	}
	if input.BalanceFactor < 0 || input.BalanceFactor > 1 {
		return fmt.Errorf("balance factor must be between 0 and 1, got %.2f", input.BalanceFactor)
	}
	return nil
}

// RefreshCycleConfig defines the periodic resource refresh configuration.
type RefreshCycleConfig struct {
	// Interval between refreshes (in game hours)
	RefreshIntervalHours int

	// Maximum resources per region
	MaxResourcesPerRegion int

	// Refresh percentage (percentage of depleted resources to restore)
	RefreshPercentage float64
}

// DefaultRefreshCycleConfig returns the default refresh configuration.
func DefaultRefreshCycleConfig() *RefreshCycleConfig {
	return &RefreshCycleConfig{
		RefreshIntervalHours:  24, // once per day
		MaxResourcesPerRegion: 100,
		RefreshPercentage:     0.3, // 30% refresh
	}
}

// CalculateRefreshAmount calculates how many resources to refresh.
func (c *RefreshCycleConfig) CalculateRefreshAmount(currentAmount, maxAmount int) int {
	if currentAmount >= maxAmount {
		return 0
	}

	deficit := maxAmount - currentAmount
	refresh := int(float64(deficit) * c.RefreshPercentage)

	if refresh <= 0 {
		refresh = 1 // minimum 1 resource
	}

	// Don't exceed max
	if currentAmount+refresh > maxAmount {
		refresh = maxAmount - currentAmount
	}

	return refresh
}

// ShouldRefresh checks if it's time for a resource refresh.
func (c *RefreshCycleConfig) ShouldRefresh(lastRefreshHours, currentHours int) bool {
	return currentHours-lastRefreshHours >= c.RefreshIntervalHours
}

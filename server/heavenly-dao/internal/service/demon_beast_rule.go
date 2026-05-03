package service

import (
	"math"

	"github.com/cultivation-world/shared/types"
)

// DemonBeastRule handles demon beast spawn rate and level distribution.
type DemonBeastRule struct {
	baseSpawnRate    float64
	maxSpawnRate     float64
}

// NewDemonBeastRule creates a new DemonBeastRule with default configuration.
func NewDemonBeastRule() *DemonBeastRule {
	return &DemonBeastRule{
		baseSpawnRate: 0.05, // 5% base spawn chance per tick
		maxSpawnRate:  0.80,
	}
}

// RegionState describes a region's properties for beast spawning.
type RegionState struct {
	SpiritualDensity float64               // spiritual energy density (0.0-1.0)
	RegionLevel      int                   // region danger level (1-10)
	IsDangerZone     bool                  // whether this is a danger zone
	BeastPopulation  int                   // current beast count
	MaxBeastCapacity int                   // maximum beast capacity
}

// CalculateBeastSpawnRate computes the spawn rate of demon beasts in a region.
//
// Formula:
//
//	rate = base_rate * spiritual_density * region_factor * population_factor * danger_bonus
func (r *DemonBeastRule) CalculateBeastSpawnRate(state RegionState) float64 {
	// Spiritual density factor
	spiritualFactor := math.Max(0.0, math.Min(1.0, state.SpiritualDensity))
	spiritualFactor = 0.1 + spiritualFactor*1.9 // range: 0.1-2.0

	// Region level factor
	regionFactor := 1.0 + float64(state.RegionLevel-1)*0.15

	// Population pressure: fewer beasts = higher spawn rate
	popRatio := 0.0
	if state.MaxBeastCapacity > 0 {
		popRatio = float64(state.BeastPopulation) / float64(state.MaxBeastCapacity)
	}
	populationFactor := math.Max(0.0, 1.0-popRatio)

	// Danger zone bonus
	dangerBonus := 1.0
	if state.IsDangerZone {
		dangerBonus = 2.0
	}

	rate := r.baseSpawnRate * spiritualFactor * regionFactor * populationFactor * dangerBonus
	rate = math.Max(0.0, math.Min(r.maxSpawnRate, rate))

	return rate
}

// BeastLevelDistribution returns the probability distribution of beast levels for a region.
// Returns a slice where index i corresponds to level i+1 (Qi=0, Foundation=1, etc.).
func (r *DemonBeastRule) GenerateBeastLevelDistribution(regionLevel int, spiritualDensity float64) map[types.CultivationRealm]float64 {
	// Max realm based on region level and spiritual density
	maxRealmIdx := int(float64(regionLevel)/2.0 + spiritualDensity*3.0)
	if maxRealmIdx < 0 {
		maxRealmIdx = 0
	}
	if maxRealmIdx > 5 {
		maxRealmIdx = 5 // cap at SoulTransform for normal regions
	}

	// Build distribution: higher levels near maxRealmIdx have more weight
	distribution := make(map[types.CultivationRealm]float64)
	realms := []types.CultivationRealm{
		types.RealmQiCondensation,
		types.RealmFoundation,
		types.RealmGoldenCore,
		types.RealmNascentSoul,
		types.RealmSoulTransform,
		types.RealmVoidRefinement,
	}

	totalWeight := 0.0
	weights := make([]float64, len(realms))
	for i := range realms {
		if i > maxRealmIdx {
			weights[i] = 0
		} else {
			// Bell curve: peak at maxRealmIdx-1
			dist := float64(i) - float64(maxRealmIdx) + 1.0
			weights[i] = math.Exp(-dist * dist / 2.0)
		}
		totalWeight += weights[i]
	}

	if totalWeight <= 0 {
		distribution[types.RealmQiCondensation] = 1.0
		return distribution
	}

	for i, realm := range realms {
		distribution[realm] = weights[i] / totalWeight
	}

	return distribution
}

// BeastSpawnResult represents the result of a beast spawn check.
type BeastSpawnResult struct {
	Spawns     bool
	Level      types.CultivationRealm
	Count      int
}

// ResolveBeastSpawn determines whether beasts spawn and at what level.
func (r *DemonBeastRule) ResolveBeastSpawn(state RegionState, randFloat func() float64) BeastSpawnResult {
	spawnRate := r.CalculateBeastSpawnRate(state)

	if randFloat() >= spawnRate {
		return BeastSpawnResult{Spawns: false}
	}

	// Determine level from distribution
	distribution := r.GenerateBeastLevelDistribution(state.RegionLevel, state.SpiritualDensity)

	// Weighted random selection
	randVal := randFloat()
	cumulative := 0.0
	selectedRealm := types.RealmQiCondensation
	for realm, prob := range distribution {
		cumulative += prob
		if randVal <= cumulative {
			selectedRealm = realm
			break
		}
	}

	// Determine count (1-3 based on spawn rate)
	count := 1
	if spawnRate > 0.5 && randFloat() < 0.3 {
		count = 2
	}
	if spawnRate > 0.7 && randFloat() < 0.1 {
		count = 3
	}

	return BeastSpawnResult{
		Spawns: true,
		Level:  selectedRealm,
		Count:  count,
	}
}

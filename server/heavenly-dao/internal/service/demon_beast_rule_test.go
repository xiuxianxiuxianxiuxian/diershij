package service

import (
	"testing"

	"github.com/cultivation-world/shared/types"
	"github.com/stretchr/testify/assert"
)

func TestNewDemonBeastRule(t *testing.T) {
	rule := NewDemonBeastRule()
	assert.NotNil(t, rule)
	assert.Equal(t, 0.05, rule.baseSpawnRate)
	assert.Equal(t, 0.80, rule.maxSpawnRate)
}

func TestCalculateBeastSpawnRate_Baseline(t *testing.T) {
	rule := NewDemonBeastRule()

	state := RegionState{
		SpiritualDensity: 0.5,
		RegionLevel:      3,
		BeastPopulation:  10,
		MaxBeastCapacity: 100,
	}

	rate := rule.CalculateBeastSpawnRate(state)
	// spiritual = 0.1 + 0.5*1.9 = 1.05
	// region = 1.0 + 2*0.15 = 1.3
	// pop = 1 - 10/100 = 0.9
	// rate = 0.05 * 1.05 * 1.3 * 0.9 = 0.061425
	assert.InDelta(t, 0.0614, rate, 0.001)
}

func TestCalculateBeastSpawnRate_HighDensity(t *testing.T) {
	rule := NewDemonBeastRule()

	state := RegionState{
		SpiritualDensity: 1.0,
		RegionLevel:      5,
		BeastPopulation:  5,
		MaxBeastCapacity: 100,
	}

	rate := rule.CalculateBeastSpawnRate(state)
	// spiritual = 0.1 + 1.0*1.9 = 2.0
	// region = 1.0 + 4*0.15 = 1.6
	// pop = 1 - 5/100 = 0.95
	// rate = 0.05 * 2.0 * 1.6 * 0.95 = 0.152
	assert.InDelta(t, 0.152, rate, 0.001)
}

func TestCalculateBeastSpawnRate_DangerZone(t *testing.T) {
	rule := NewDemonBeastRule()

	state := RegionState{
		SpiritualDensity: 0.8,
		RegionLevel:      7,
		IsDangerZone:     true,
		BeastPopulation:  20,
		MaxBeastCapacity: 200,
	}

	rate := rule.CalculateBeastSpawnRate(state)
	// spiritual = 0.1 + 0.8*1.9 = 1.62
	// region = 1.0 + 6*0.15 = 1.9
	// pop = 1 - 20/200 = 0.9
	// danger = 2.0
	// rate = 0.05 * 1.62 * 1.9 * 0.9 * 2.0 = 0.277
	assert.InDelta(t, 0.277, rate, 0.001)
}

func TestCalculateBeastSpawnRate_FullPopulation(t *testing.T) {
	rule := NewDemonBeastRule()

	state := RegionState{
		SpiritualDensity: 1.0,
		RegionLevel:      5,
		BeastPopulation:  100,
		MaxBeastCapacity: 100,
	}

	rate := rule.CalculateBeastSpawnRate(state)
	// pop = 1 - 100/100 = 0 → rate = 0
	assert.InDelta(t, 0.0, rate, 0.001)
}

func TestCalculateBeastSpawnRate_EmptyRegion(t *testing.T) {
	rule := NewDemonBeastRule()

	state := RegionState{
		SpiritualDensity: 1.0,
		RegionLevel:      1,
		BeastPopulation:  0,
		MaxBeastCapacity: 100,
	}

	rate := rule.CalculateBeastSpawnRate(state)
	// spiritual = 2.0, region = 1.0, pop = 1.0
	// rate = 0.05 * 2.0 * 1.0 * 1.0 = 0.10
	assert.InDelta(t, 0.10, rate, 0.001)
}

func TestGenerateBeastLevelDistribution_LowRegion(t *testing.T) {
	rule := NewDemonBeastRule()

	dist := rule.GenerateBeastLevelDistribution(1, 0.3)
	// maxRealmIdx = 1/2 + 0.3*3 = 0.5 + 0.9 = 1.4 → 1
	// Most weight should be on QiCondensation
	assert.True(t, dist[types.RealmQiCondensation] > 0.5)
	assert.True(t, dist[types.RealmFoundation] > 0)
}

func TestGenerateBeastLevelDistribution_HighRegion(t *testing.T) {
	rule := NewDemonBeastRule()

	dist := rule.GenerateBeastLevelDistribution(8, 1.0)
	// maxRealmIdx = 8/2 + 1.0*3 = 4 + 3 = 7 → capped at 5
	// Should have distribution across multiple realms
	total := 0.0
	for _, prob := range dist {
		total += prob
	}
	assert.InDelta(t, 1.0, total, 0.001)
}

func TestResolveBeastSpawn_NoSpawn(t *testing.T) {
	rule := NewDemonBeastRule()
	alwaysFail := func() float64 { return 0.99 }

	state := RegionState{
		SpiritualDensity: 0.1,
		RegionLevel:      1,
		BeastPopulation:  90,
		MaxBeastCapacity: 100,
	}

	result := rule.ResolveBeastSpawn(state, alwaysFail)
	assert.False(t, result.Spawns)
}

func TestResolveBeastSpawn_SpawnSingle(t *testing.T) {
	rule := NewDemonBeastRule()
	callCount := 0
	deterministic := func() float64 {
		callCount++
		if callCount == 1 {
			return 0.01 // spawn
		}
		return 0.3 // select realm
	}

	state := RegionState{
		SpiritualDensity: 1.0,
		RegionLevel:      3,
		BeastPopulation:  0,
		MaxBeastCapacity: 100,
	}

	result := rule.ResolveBeastSpawn(state, deterministic)
	assert.True(t, result.Spawns)
	assert.True(t, result.Count >= 1)
}

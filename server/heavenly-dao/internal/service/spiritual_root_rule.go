package service

import (
	"math"
)

// SpiritualRootRule handles spiritual root awakening and mutation calculations.
type SpiritualRootRule struct {
	baseAwakeningRate float64
	mutationRate      float64
}

// NewSpiritualRootRule creates a new SpiritualRootRule with default configuration.
func NewSpiritualRootRule() *SpiritualRootRule {
	return &SpiritualRootRule{
		baseAwakeningRate: 0.01, // 1% base chance per year
		mutationRate:      0.10, // 10% chance of mutation when awakening
	}
}

// AwakeningInput holds inputs for calculating spiritual root awakening.
type AwakeningInput struct {
	Age              int     // person's age
	Luck             int     // luck attribute (0-100)
	FamilyBloodline  float64 // family bloodline quality (0.0-1.0)
	SpiritualDensity float64 // local spiritual density (0.0-1.0)
	IsMortal         bool    // whether the person is currently a mortal
}

// CalculateAwakeningRate computes the probability of spiritual root awakening.
//
// Formula:
//
//	rate = base_rate * age_factor * luck_factor * bloodline * spiritual_density
func (r *SpiritualRootRule) CalculateAwakeningRate(input AwakeningInput) float64 {
	if !input.IsMortal {
		return 0 // already awakened
	}

	// Age factor: peaks at 6-16 years old
	var ageFactor float64
	if input.Age < 6 {
		ageFactor = float64(input.Age) / 6.0 * 0.5
	} else if input.Age <= 16 {
		ageFactor = 1.0
	} else if input.Age <= 30 {
		ageFactor = 1.0 - (float64(input.Age)-16.0)/14.0*0.5
	} else {
		ageFactor = 0.5 - (float64(input.Age)-30.0)/70.0*0.5
		ageFactor = math.Max(0.01, ageFactor)
	}

	// Luck factor
	luckFactor := 0.5 + float64(input.Luck)/100.0
	luckFactor = math.Max(0.5, math.Min(1.5, luckFactor))

	// Bloodline quality
	bloodline := math.Max(0.0, math.Min(1.0, input.FamilyBloodline))
	if bloodline == 0.0 {
		bloodline = 0.1 // minimum baseline
	}

	// Spiritual density
	spiritualDensity := math.Max(0.0, math.Min(1.0, input.SpiritualDensity))
	if spiritualDensity == 0.0 {
		spiritualDensity = 0.05
	}

	rate := r.baseAwakeningRate * ageFactor * luckFactor * bloodline * spiritualDensity

	// Clamp to reasonable range
	return math.Max(0.0001, math.Min(0.50, rate))
}

// GenerateMutatedElement generates a mutated element from base elements.
// Mutations: water+wind→ice, wind+fire→lightning, wood+fire→poison
func (r *SpiritualRootRule) GenerateMutatedElement(primaryElement string, randFloat func() float64) string {
	// Check if mutation occurs
	if randFloat() >= r.mutationRate {
		return "" // no mutation
	}

	// Mutation table
	mutations := map[string][]string{
		"water": {"ice"},
		"wind":  {"ice", "lightning"},
		"fire":  {"lightning"},
		"wood":  {"poison"},
	}

	candidates, ok := mutations[primaryElement]
	if !ok || len(candidates) == 0 {
		return ""
	}

	// Randomly pick a candidate
	idx := int(randFloat() * float64(len(candidates)))
	if idx >= len(candidates) {
		idx = len(candidates) - 1
	}
	return candidates[idx]
}

// CalculateRootQuality computes the quality score of a spiritual root configuration.
func (r *SpiritualRootRule) CalculateRootQuality(element string, purity int) float64 {
	// Base quality from purity (0-100)
	baseQuality := float64(purity) / 100.0

	// Element type bonus
	rootBonus := 1.0
	switch element {
	case "fire", "water", "wood", "metal", "earth":
		rootBonus = 1.0 // single elemental
	case "ice", "lightning", "poison", "wind", "thunder":
		rootBonus = 1.2 // mutated/rare elements
	case "heavenly", "void", "time", "space", "light", "dark":
		rootBonus = 1.5 // supreme elements
	default:
		rootBonus = 0.8
	}

	return baseQuality * rootBonus
}

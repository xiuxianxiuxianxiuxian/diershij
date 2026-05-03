package service

import (
	"math"
)

// FormationRule handles formation power calculation and break formation rate.
type FormationRule struct {
	baseBreakRate    float64
	minBreakRate     float64
	maxBreakRate     float64
}

// NewFormationRule creates a new FormationRule with default configuration.
func NewFormationRule() *FormationRule {
	return &FormationRule{
		baseBreakRate: 0.10,
		minBreakRate:  0.01,
		maxBreakRate:  0.80,
	}
}

// FormationNode represents a node in a formation array.
type FormationNode struct {
	ElementType   string
	Power         float64 // node power contribution
	IsCore        bool    // whether this is the core node
}

// Formation represents a complete formation array.
type Formation struct {
	Name         string
	RequiredLevel int      // required formation skill level
	BasePower    float64   // base formation power
	Difficulty   float64   // formation complexity (1.0-3.0)
	ElementType  string    // primary element
	Nodes        []FormationNode
}

// FormationInput holds inputs for calculating formation effectiveness.
type FormationInput struct {
	FormationLevel   int     // formation skill level (1-10)
	MentalStability  int     // mental stability (0-100)
	LocationBonus    float64 // location bonus (1.0+)
	EquipmentBonus   float64 // formation flags/tools bonus (1.0+)
	ElementAffinity  float64 // element match (0.0-1.0)
	SpiritStoneCount int     // spirit stones used to power formation
	Luck             int     // luck attribute (0-100)
}

// CalculateFormationPower computes the total power of an active formation.
//
// Formula:
//
//	power = base_power * node_sum * skill_factor * mental * location * equipment * element * stone_bonus * luck
func (r *FormationRule) CalculateFormationPower(formation Formation, input FormationInput) float64 {
	// Sum node powers, with core node bonus
	nodePower := 0.0
	for _, node := range formation.Nodes {
		if node.IsCore {
			nodePower += node.Power * 1.5
		} else {
			nodePower += node.Power
		}
	}
	if len(formation.Nodes) == 0 {
		nodePower = 1.0
	}

	// Skill factor
	skillFactor := float64(input.FormationLevel) / float64(formation.RequiredLevel)
	if skillFactor <= 0 {
		skillFactor = 0.1
	}
	skillFactor = math.Min(2.0, skillFactor)

	// Mental factor
	mentalFactor := float64(input.MentalStability) / 100.0
	mentalFactor = math.Max(0.1, math.Min(1.0, mentalFactor))

	// Location bonus
	locationBonus := math.Max(1.0, input.LocationBonus)

	// Equipment bonus
	equipmentBonus := math.Max(1.0, input.EquipmentBonus)

	// Element affinity
	elementAffinity := math.Max(0.0, math.Min(1.0, input.ElementAffinity))
	if elementAffinity == 0.0 {
		elementAffinity = 1.0
	}

	// Spirit stone bonus: each stone adds 5%, max 50%
	stoneBonus := 1.0 + math.Min(0.50, float64(input.SpiritStoneCount)*0.05)

	// Luck factor
	luckFactor := 1.0 + (float64(input.Luck)-50.0)/200.0
	luckFactor = math.Max(0.75, math.Min(1.25, luckFactor))

	power := formation.BasePower * nodePower * skillFactor * mentalFactor *
		locationBonus * equipmentBonus * elementAffinity * stoneBonus * luckFactor

	return math.Max(0, power)
}

// CalculateBreakFormationRate computes the probability of breaking through an enemy formation.
func (r *FormationRule) CalculateBreakFormationRate(
	attackerPower float64,
	defenderFormation Formation,
	defenderInput FormationInput,
) float64 {
	defensePower := r.CalculateFormationPower(defenderFormation, defenderInput)

	if defensePower <= 0 {
		return r.maxBreakRate
	}

	// Ratio-based break rate
	ratio := attackerPower / defensePower

	// Map ratio to break rate: ratio=0.5 → low, ratio=1.0 → base, ratio=2.0 → high
	breakRate := r.baseBreakRate * math.Pow(ratio, 2.0)
	breakRate = math.Max(r.minBreakRate, math.Min(r.maxBreakRate, breakRate))

	return breakRate
}

// FormationResult represents the result of a formation interaction.
type FormationResult struct {
	Success       bool
	Power         float64
	BreakRate     float64
	Message       string
}

// ResolveFormation evaluates a formation's state and potential break rate.
func (r *FormationRule) ResolveFormation(
	formation Formation,
	input FormationInput,
	attackerPower float64,
) FormationResult {
	power := r.CalculateFormationPower(formation, input)
	breakRate := r.CalculateBreakFormationRate(attackerPower, formation, input)

	// Formation holds if break rate < 50%
	holds := breakRate < 0.50

	if holds {
		return FormationResult{
			Success:   true,
			Power:     power,
			BreakRate: breakRate,
			Message:   "阵法稳固",
		}
	}

	return FormationResult{
		Success:   false,
		Power:     power,
		BreakRate: breakRate,
		Message:   "阵法即将被破",
	}
}

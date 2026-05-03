package service

import (
	"math"
)

// AlchemyRule handles pill refinement success rate, quality generation, and failure handling.
type AlchemyRule struct {
	baseSuccessRate    float64
	minSuccessRate     float64
	maxSuccessRate     float64
	explosionChance    float64 // chance of furnace explosion on failure
	materialLossRate   float64 // material loss ratio on failure
}

// NewAlchemyRule creates a new AlchemyRule with default configuration.
func NewAlchemyRule() *AlchemyRule {
	return &AlchemyRule{
		baseSuccessRate:   0.50,
		minSuccessRate:    0.05,
		maxSuccessRate:    0.95,
		explosionChance:   0.10,
		materialLossRate:  0.50,
	}
}

// PillRecipe represents a pill recipe for alchemy.
type PillRecipe struct {
	Name             string
	RequiredLevel    int     // required alchemy skill level
	BaseQuality      float64 // base quality tier (1-5)
	Difficulty       float64 // recipe difficulty multiplier (1.0-3.0)
	ElementType      string  // primary element requirement
	IngredientCount  int     // number of ingredients
}

// AlchemyInput holds the inputs needed to calculate pill success rate.
type AlchemyInput struct {
	AlchemyLevel      int     // alchemist's skill level (1-10)
	MentalStability   int     // mental stability during refinement (0-100)
	LocationBonus     float64 // alchemy lab bonus multiplier (1.0+)
	EquipmentBonus    float64 // equipment/tool bonus (1.0+)
	ElementAffinity   float64 // element match factor (0.0-1.0)
	Luck              int     // luck attribute (0-100)
}

// CalculatePillSuccessRate computes the probability of successfully refining a pill.
//
// Formula:
//
//	prob = base_rate * skill_factor * mental * location * equipment * element * luck
//
// Where:
//   - skill_factor = min(2.0, alchemy_level / recipe_required_level)
//   - mental = mental_stability / 100.0
//   - location = location bonus (default 1.0)
//   - equipment = equipment bonus (default 1.0)
//   - element = element affinity (default 1.0)
//   - luck_factor = 1.0 + (luck - 50) / 200.0
func (a *AlchemyRule) CalculatePillSuccessRate(recipe PillRecipe, input AlchemyInput) float64 {
	// Skill factor: capped at 2.0
	skillFactor := float64(input.AlchemyLevel) / float64(recipe.RequiredLevel)
	if skillFactor <= 0 {
		skillFactor = 0.1
	}
	skillFactor = math.Min(2.0, skillFactor)

	// Mental factor
	mentalFactor := float64(input.MentalStability) / 100.0
	mentalFactor = math.Max(0.1, math.Min(1.0, mentalFactor))

	// Location bonus (default 1.0)
	locationBonus := math.Max(1.0, input.LocationBonus)

	// Equipment bonus (default 1.0)
	equipmentBonus := math.Max(1.0, input.EquipmentBonus)

	// Element affinity (0.0-1.0, default 1.0 if not applicable)
	elementAffinity := math.Max(0.0, math.Min(1.0, input.ElementAffinity))
	if elementAffinity == 0.0 {
		elementAffinity = 1.0 // no element requirement
	}

	// Luck factor
	luckFactor := 1.0 + (float64(input.Luck)-50.0)/200.0
	luckFactor = math.Max(0.75, math.Min(1.25, luckFactor))

	// Recipe difficulty penalty
	difficultyPenalty := 1.0 / recipe.Difficulty

	prob := a.baseSuccessRate * skillFactor * mentalFactor * locationBonus *
		equipmentBonus * elementAffinity * luckFactor * difficultyPenalty

	// Clamp to [min, max]
	prob = math.Max(a.minSuccessRate, math.Min(a.maxSuccessRate, prob))

	return prob
}

// PillGrade represents the quality tier of a refined pill.
type PillGrade string

const (
	PillGradeDefective    PillGrade = "defective"    // 劣品 (1)
	PillGradeCommon       PillGrade = "common"       // 普通 (2)
	PillGradeFine         PillGrade = "fine"         // 精品 (3)
	PillGradeExcellent    PillGrade = "excellent"    // 极品 (4)
	PillGradePerfect      PillGrade = "perfect"      // 绝品 (5)
)

// PillResult represents the outcome of an alchemy attempt.
type PillResult struct {
	Success  bool
	Grade    PillGrade
	Quality  float64
	Explosion bool
	Message  string
}

// GeneratePillQuality determines the grade of a successfully refined pill.
//
// Based on success rate and base quality:
//   - Higher success rate → higher chance of better grade
//   - Base quality sets the expected tier
func (a *AlchemyRule) GeneratePillQuality(successRate float64, baseQuality float64, randFloat func() float64) PillGrade {
	// Quality tier = base_quality + bonus from high success rate
	expectedTier := baseQuality + (successRate-0.5)*4.0 // range: base-2 to base+2
	expectedTier = math.Max(1, math.Min(5, expectedTier))

	// Use random to determine actual tier with variance
	r := randFloat()
	if r < 0.1 && expectedTier >= 5 {
		return PillGradePerfect
	}
	if r < 0.3 && expectedTier >= 4 {
		return PillGradeExcellent
	}
	if r < 0.5 && expectedTier >= 3 {
		return PillGradeFine
	}
	if r < 0.7 && expectedTier >= 2 {
		return PillGradeCommon
	}
	return PillGradeDefective
}

// ResolveAlchemy determines the full outcome of an alchemy attempt.
func (a *AlchemyRule) ResolveAlchemy(recipe PillRecipe, input AlchemyInput, randFloat func() float64) PillResult {
	successRate := a.CalculatePillSuccessRate(recipe, input)

	if randFloat() < successRate {
		grade := a.GeneratePillQuality(successRate, recipe.BaseQuality, randFloat)
		return PillResult{
			Success: true,
			Grade:   grade,
			Quality: float64(gradeToInt(grade)),
			Message: "炼丹成功",
		}
	}

	// Failure - check for explosion
	explosion := randFloat() < a.explosionChance

	if explosion {
		return PillResult{
			Success:   false,
			Explosion: true,
			Message:   "炼丹失败，丹炉爆炸",
		}
	}

	return PillResult{
		Success: false,
		Message: "炼丹失败，材料损毁",
	}
}

// CalculateMaterialLoss returns the ratio of materials lost on failure.
func (a *AlchemyRule) CalculateMaterialLoss(explosion bool) float64 {
	if explosion {
		return 1.0 // total loss
	}
	return a.materialLossRate
}

func gradeToInt(g PillGrade) int {
	switch g {
	case PillGradeDefective:
		return 1
	case PillGradeCommon:
		return 2
	case PillGradeFine:
		return 3
	case PillGradeExcellent:
		return 4
	case PillGradePerfect:
		return 5
	default:
		return 1
	}
}

package service

import (
	"math"
)

// ForgingRule handles artifact forging success rate and grade generation.
type ForgingRule struct {
	baseSuccessRate  float64
	minSuccessRate   float64
	maxSuccessRate   float64
	breakRate        float64 // chance of material breakage on failure
}

// NewForgingRule creates a new ForgingRule with default configuration.
func NewForgingRule() *ForgingRule {
	return &ForgingRule{
		baseSuccessRate: 0.45,
		minSuccessRate:  0.05,
		maxSuccessRate:  0.90,
		breakRate:       0.15,
	}
}

// ArtifactRecipe represents a recipe for forging an artifact.
type ArtifactRecipe struct {
	Name          string
	RequiredLevel int     // required artificing skill level
	BaseGrade     float64 // base grade tier (1-5)
	Difficulty    float64 // recipe difficulty (1.0-3.0)
	ElementType   string  // primary element requirement
	MaterialCount int     // number of materials needed
}

// ForgingInput holds the inputs needed to calculate forging success rate.
type ForgingInput struct {
	ArtificingLevel   int     // artificer's skill level (1-10)
	MentalStability   int     // mental stability during forging (0-100)
	LocationBonus     float64 // forging hall bonus (1.0+)
	EquipmentBonus    float64 // equipment/tool bonus (1.0+)
	ElementAffinity   float64 // element match factor (0.0-1.0)
	Luck              int     // luck attribute (0-100)
	SpiritualRoot     string  // spiritual root element (for element affinity bonus)
}

// CalculateForgingSuccessRate computes the probability of successfully forging an artifact.
//
// Formula:
//
//	prob = base_rate * skill_factor * mental * location * equipment * element * luck * difficulty_penalty
func (f *ForgingRule) CalculateForgingSuccessRate(recipe ArtifactRecipe, input ForgingInput) float64 {
	// Skill factor: capped at 2.0
	skillFactor := float64(input.ArtificingLevel) / float64(recipe.RequiredLevel)
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

	// Element affinity
	elementAffinity := math.Max(0.0, math.Min(1.0, input.ElementAffinity))
	if elementAffinity == 0.0 {
		elementAffinity = 1.0
	}

	// Luck factor
	luckFactor := 1.0 + (float64(input.Luck)-50.0)/200.0
	luckFactor = math.Max(0.75, math.Min(1.25, luckFactor))

	// Recipe difficulty penalty
	difficultyPenalty := 1.0 / recipe.Difficulty

	prob := f.baseSuccessRate * skillFactor * mentalFactor * locationBonus *
		equipmentBonus * elementAffinity * luckFactor * difficultyPenalty

	prob = math.Max(f.minSuccessRate, math.Min(f.maxSuccessRate, prob))
	return prob
}

// ArtifactGrade represents the quality tier of a forged artifact.
type ArtifactGrade string

const (
	ArtifactGradeMortal     ArtifactGrade = "mortal"     // 凡器 (1)
	ArtifactGradeSpiritual  ArtifactGrade = "spiritual"  // 灵器 (2)
	ArtifactGradeEarthly    ArtifactGrade = "earthly"    // 地器 (3)
	ArtifactGradeHeavenly   ArtifactGrade = "heavenly"   // 天器 (4)
	ArtifactGradeDivine     ArtifactGrade = "divine"     // 神器 (5)
)

// ForgingResult represents the outcome of a forging attempt.
type ForgingResult struct {
	Success      bool
	Grade        ArtifactGrade
	Quality      float64
	MaterialBroken bool
	Message      string
}

// GenerateArtifactGrade determines the grade of a successfully forged artifact.
func (f *ForgingRule) GenerateArtifactGrade(successRate float64, baseGrade float64, randFloat func() float64) ArtifactGrade {
	expectedTier := baseGrade + (successRate-0.45)*4.0
	expectedTier = math.Max(1, math.Min(5, expectedTier))

	r := randFloat()
	if r < 0.05 && expectedTier >= 5 {
		return ArtifactGradeDivine
	}
	if r < 0.2 && expectedTier >= 4 {
		return ArtifactGradeHeavenly
	}
	if r < 0.4 && expectedTier >= 3 {
		return ArtifactGradeEarthly
	}
	if r < 0.6 && expectedTier >= 2 {
		return ArtifactGradeSpiritual
	}
	return ArtifactGradeMortal
}

// ResolveForging determines the full outcome of a forging attempt.
func (f *ForgingRule) ResolveForging(recipe ArtifactRecipe, input ForgingInput, randFloat func() float64) ForgingResult {
	successRate := f.CalculateForgingSuccessRate(recipe, input)

	if randFloat() < successRate {
		grade := f.GenerateArtifactGrade(successRate, recipe.BaseGrade, randFloat)
		return ForgingResult{
			Success: true,
			Grade:   grade,
			Quality: float64(artifactGradeToInt(grade)),
			Message: "炼器成功",
		}
	}

	// Failure - check for material breakage
	materialBroken := randFloat() < f.breakRate

	if materialBroken {
		return ForgingResult{
			Success:        false,
			MaterialBroken: true,
			Message:        "炼器失败，材料损毁",
		}
	}

	return ForgingResult{
		Success: false,
		Message: "炼器失败，材料可回收",
	}
}

// CalculateMaterialLoss returns the ratio of materials lost on failure.
func (f *ForgingRule) CalculateMaterialLoss(materialBroken bool) float64 {
	if materialBroken {
		return 1.0
	}
	return 0.3 // 30% loss when materials can be partially recovered
}

func artifactGradeToInt(g ArtifactGrade) int {
	switch g {
	case ArtifactGradeMortal:
		return 1
	case ArtifactGradeSpiritual:
		return 2
	case ArtifactGradeEarthly:
		return 3
	case ArtifactGradeHeavenly:
		return 4
	case ArtifactGradeDivine:
		return 5
	default:
		return 1
	}
}

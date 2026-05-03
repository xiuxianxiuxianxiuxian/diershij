package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAlchemyRule(t *testing.T) {
	rule := NewAlchemyRule()
	assert.NotNil(t, rule)
	assert.Equal(t, 0.50, rule.baseSuccessRate)
	assert.Equal(t, 0.05, rule.minSuccessRate)
	assert.Equal(t, 0.95, rule.maxSuccessRate)
}

func TestCalculatePillSuccessRate_Baseline(t *testing.T) {
	rule := NewAlchemyRule()

	recipe := PillRecipe{
		Name:          "basic_pill",
		RequiredLevel: 3,
		BaseQuality:   2,
		Difficulty:    1.0,
	}

	input := AlchemyInput{
		AlchemyLevel:    3, // skill_factor = 1.0
		MentalStability: 100,
		LocationBonus:   1.0,
		EquipmentBonus:  1.0,
		ElementAffinity: 1.0,
		Luck:            50, // luck_factor = 1.0
	}

	prob := rule.CalculatePillSuccessRate(recipe, input)
	// 0.50 * 1.0 * 1.0 * 1.0 * 1.0 * 1.0 * 1.0 * 1.0 = 0.50
	assert.InDelta(t, 0.50, prob, 0.001)
}

func TestCalculatePillSuccessRate_HighSkill(t *testing.T) {
	rule := NewAlchemyRule()

	recipe := PillRecipe{
		Name:          "advanced_pill",
		RequiredLevel: 2,
		BaseQuality:   3,
		Difficulty:    1.0,
	}

	input := AlchemyInput{
		AlchemyLevel:    5, // skill_factor = min(2.0, 2.5) = 2.0
		MentalStability: 100,
		LocationBonus:   1.0,
		EquipmentBonus:  1.0,
		ElementAffinity: 1.0,
		Luck:            50,
	}

	prob := rule.CalculatePillSuccessRate(recipe, input)
	// 0.50 * 2.0 * 1.0 * 1.0 * 1.0 * 1.0 * 1.0 = 1.0 → clamped to 0.95
	assert.InDelta(t, 0.95, prob, 0.001)
}

func TestCalculatePillSuccessRate_LowSkill(t *testing.T) {
	rule := NewAlchemyRule()

	recipe := PillRecipe{
		Name:          "hard_pill",
		RequiredLevel: 5,
		BaseQuality:   4,
		Difficulty:    1.0,
	}

	input := AlchemyInput{
		AlchemyLevel:    2, // skill_factor = 0.4
		MentalStability: 80,
		LocationBonus:   1.0,
		EquipmentBonus:  1.0,
		ElementAffinity: 1.0,
		Luck:            50,
	}

	prob := rule.CalculatePillSuccessRate(recipe, input)
	// 0.50 * 0.4 * 0.8 = 0.16
	assert.InDelta(t, 0.16, prob, 0.001)
}

func TestCalculatePillSuccessRate_VeryLowSkill(t *testing.T) {
	rule := NewAlchemyRule()

	recipe := PillRecipe{
		Name:          "expert_pill",
		RequiredLevel: 8,
		BaseQuality:   5,
		Difficulty:    1.0,
	}

	input := AlchemyInput{
		AlchemyLevel:    1, // skill_factor = 0.125
		MentalStability: 50,
		LocationBonus:   1.0,
		EquipmentBonus:  1.0,
		ElementAffinity: 1.0,
		Luck:            50,
	}

	prob := rule.CalculatePillSuccessRate(recipe, input)
	// 0.50 * 0.125 * 0.5 = 0.03125 → clamped to 0.05
	assert.InDelta(t, 0.05, prob, 0.001)
}

func TestCalculatePillSuccessRate_LocationBonus(t *testing.T) {
	rule := NewAlchemyRule()

	recipe := PillRecipe{
		Name:          "basic_pill",
		RequiredLevel: 3,
		Difficulty:    1.0,
	}

	input := AlchemyInput{
		AlchemyLevel:    3,
		MentalStability: 100,
		LocationBonus:   1.5, // alchemy lab bonus
		EquipmentBonus:  1.0,
		ElementAffinity: 1.0,
		Luck:            50,
	}

	prob := rule.CalculatePillSuccessRate(recipe, input)
	// 0.50 * 1.0 * 1.0 * 1.5 = 0.75
	assert.InDelta(t, 0.75, prob, 0.001)
}

func TestCalculatePillSuccessRate_EquipmentBonus(t *testing.T) {
	rule := NewAlchemyRule()

	recipe := PillRecipe{
		Name:          "basic_pill",
		RequiredLevel: 3,
		Difficulty:    1.0,
	}

	input := AlchemyInput{
		AlchemyLevel:    3,
		MentalStability: 100,
		LocationBonus:   1.0,
		EquipmentBonus:  1.3, // high-quality furnace
		ElementAffinity: 1.0,
		Luck:            50,
	}

	prob := rule.CalculatePillSuccessRate(recipe, input)
	// 0.50 * 1.0 * 1.0 * 1.3 = 0.65
	assert.InDelta(t, 0.65, prob, 0.001)
}

func TestCalculatePillSuccessRate_Difficulty(t *testing.T) {
	rule := NewAlchemyRule()

	recipe := PillRecipe{
		Name:          "legendary_pill",
		RequiredLevel: 5,
		Difficulty:    2.5,
	}

	input := AlchemyInput{
		AlchemyLevel:    5,
		MentalStability: 100,
		LocationBonus:   1.0,
		EquipmentBonus:  1.0,
		ElementAffinity: 1.0,
		Luck:            50,
	}

	prob := rule.CalculatePillSuccessRate(recipe, input)
	// 0.50 * 1.0 * 1.0 / 2.5 = 0.20
	assert.InDelta(t, 0.20, prob, 0.001)
}

func TestCalculatePillSuccessRate_Luck(t *testing.T) {
	rule := NewAlchemyRule()

	recipe := PillRecipe{
		Name:          "basic_pill",
		RequiredLevel: 3,
		Difficulty:    1.0,
	}

	input := AlchemyInput{
		AlchemyLevel:    3,
		MentalStability: 100,
		LocationBonus:   1.0,
		EquipmentBonus:  1.0,
		ElementAffinity: 1.0,
		Luck:            100, // luck_factor = 1.25
	}

	prob := rule.CalculatePillSuccessRate(recipe, input)
	// 0.50 * 1.0 * 1.0 * 1.25 = 0.625
	assert.InDelta(t, 0.625, prob, 0.001)

	input.Luck = 0 // luck_factor = 0.75
	prob = rule.CalculatePillSuccessRate(recipe, input)
	// 0.50 * 1.0 * 1.0 * 0.75 = 0.375
	assert.InDelta(t, 0.375, prob, 0.001)
}

func TestCalculatePillSuccessRate_AllFactors(t *testing.T) {
	rule := NewAlchemyRule()

	recipe := PillRecipe{
		Name:          "complex_pill",
		RequiredLevel: 4,
		BaseQuality:   3,
		Difficulty:    1.5,
	}

	input := AlchemyInput{
		AlchemyLevel:    5,   // skill = min(2.0, 1.25) = 1.25
		MentalStability: 90,  // mental = 0.9
		LocationBonus:   1.2,
		EquipmentBonus:  1.1,
		ElementAffinity: 0.8,
		Luck:            70,  // luck = 1.1
	}

	prob := rule.CalculatePillSuccessRate(recipe, input)
	// 0.50 * 1.25 * 0.9 * 1.2 * 1.1 * 0.8 * 1.1 / 1.5
	expected := 0.50 * 1.25 * 0.9 * 1.2 * 1.1 * 0.8 * 1.1 / 1.5
	assert.InDelta(t, expected, prob, 0.001)
}

func TestGeneratePillQuality_HighSuccess(t *testing.T) {
	rule := NewAlchemyRule()

	// Always roll high quality with high success rate
	alwaysHigh := func() float64 { return 0.01 }
	grade := rule.GeneratePillQuality(0.95, 4.0, alwaysHigh)
	assert.Equal(t, PillGradePerfect, grade)

	// Always roll low quality with high success but unlucky rolls
	alwaysLow := func() float64 { return 0.99 }
	grade = rule.GeneratePillQuality(0.95, 4.0, alwaysLow)
	// expectedTier = 4 + (0.95-0.5)*4 = 5.8, clamped to 5
	// r=0.99 >= 0.7 → defective
	assert.Equal(t, PillGradeDefective, grade)
}

func TestGeneratePillQuality_LowSuccess(t *testing.T) {
	rule := NewAlchemyRule()

	// Low success rate: expected tier = 2 + (0.1-0.5)*4 = 0.4 → clamped to 1
	lowSuccess := func() float64 { return 0.5 }
	grade := rule.GeneratePillQuality(0.10, 2.0, lowSuccess)
	// expectedTier = 2 + (0.1-0.5)*4 = 0.4 → clamped to 1
	// r=0.5 < 0.7 but expectedTier < 2, so falls to defective
	assert.Equal(t, PillGradeDefective, grade)
}

func TestResolveAlchemy_Success(t *testing.T) {
	rule := NewAlchemyRule()
	deterministicRand := func() float64 { return 0.1 } // always succeeds

	recipe := PillRecipe{
		Name:          "healing_pill",
		RequiredLevel: 2,
		BaseQuality:   2,
		Difficulty:    1.0,
	}

	input := AlchemyInput{
		AlchemyLevel:    3,
		MentalStability: 100,
		LocationBonus:   1.0,
		EquipmentBonus:  1.0,
		ElementAffinity: 1.0,
		Luck:            50,
	}

	result := rule.ResolveAlchemy(recipe, input, deterministicRand)
	assert.True(t, result.Success)
	assert.False(t, result.Explosion)
	assert.NotEmpty(t, result.Grade)
	assert.Contains(t, result.Message, "成功")
}

func TestResolveAlchemy_Failure(t *testing.T) {
	rule := NewAlchemyRule()
	deterministicRand := func() float64 { return 0.99 } // always fails

	recipe := PillRecipe{
		Name:          "expert_pill",
		RequiredLevel: 8,
		Difficulty:    1.0,
	}

	input := AlchemyInput{
		AlchemyLevel:    1,
		MentalStability: 50,
		LocationBonus:   1.0,
		EquipmentBonus:  1.0,
		ElementAffinity: 1.0,
		Luck:            50,
	}

	result := rule.ResolveAlchemy(recipe, input, deterministicRand)
	assert.False(t, result.Success)
}

func TestResolveAlchemy_Explosion(t *testing.T) {
	rule := NewAlchemyRule()
	callCount := 0
	deterministicRand := func() float64 {
		callCount++
		if callCount == 1 {
			return 0.99 // fail (99% > ~0.05 success rate)
		}
		return 0.01 // explosion (1% < 10%)
	}

	recipe := PillRecipe{
		Name:          "expert_pill",
		RequiredLevel: 8,
		Difficulty:    1.0,
	}

	input := AlchemyInput{
		AlchemyLevel:    1,
		MentalStability: 50,
	}

	result := rule.ResolveAlchemy(recipe, input, deterministicRand)
	assert.False(t, result.Success)
	assert.True(t, result.Explosion)
	assert.Contains(t, result.Message, "爆炸")
}

func TestCalculateMaterialLoss(t *testing.T) {
	rule := NewAlchemyRule()

	// Normal failure: partial loss
	loss := rule.CalculateMaterialLoss(false)
	assert.Equal(t, 0.50, loss)

	// Explosion: total loss
	loss = rule.CalculateMaterialLoss(true)
	assert.Equal(t, 1.0, loss)
}

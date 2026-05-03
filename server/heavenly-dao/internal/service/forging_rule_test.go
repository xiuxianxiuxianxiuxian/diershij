package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewForgingRule(t *testing.T) {
	rule := NewForgingRule()
	assert.NotNil(t, rule)
	assert.Equal(t, 0.45, rule.baseSuccessRate)
	assert.Equal(t, 0.05, rule.minSuccessRate)
	assert.Equal(t, 0.90, rule.maxSuccessRate)
}

func TestCalculateForgingSuccessRate_Baseline(t *testing.T) {
	rule := NewForgingRule()

	recipe := ArtifactRecipe{
		Name:          "basic_sword",
		RequiredLevel: 3,
		BaseGrade:     2,
		Difficulty:    1.0,
	}

	input := ForgingInput{
		ArtificingLevel:  3,
		MentalStability:  100,
		LocationBonus:    1.0,
		EquipmentBonus:   1.0,
		ElementAffinity:  1.0,
		Luck:             50,
	}

	prob := rule.CalculateForgingSuccessRate(recipe, input)
	// 0.45 * 1.0 * 1.0 * 1.0 * 1.0 * 1.0 * 1.0 = 0.45
	assert.InDelta(t, 0.45, prob, 0.001)
}

func TestCalculateForgingSuccessRate_HighSkill(t *testing.T) {
	rule := NewForgingRule()

	recipe := ArtifactRecipe{
		Name:          "spirit_sword",
		RequiredLevel: 2,
		Difficulty:    1.0,
	}

	input := ForgingInput{
		ArtificingLevel:  5, // skill = min(2.0, 2.5) = 2.0
		MentalStability:  100,
		LocationBonus:    1.0,
		EquipmentBonus:   1.0,
		ElementAffinity:  1.0,
		Luck:             50,
	}

	prob := rule.CalculateForgingSuccessRate(recipe, input)
	// 0.45 * 2.0 = 0.90 → at max
	assert.InDelta(t, 0.90, prob, 0.001)
}

func TestCalculateForgingSuccessRate_LowSkill(t *testing.T) {
	rule := NewForgingRule()

	recipe := ArtifactRecipe{
		Name:          "heavenly_armor",
		RequiredLevel: 7,
		Difficulty:    1.0,
	}

	input := ForgingInput{
		ArtificingLevel:  2, // skill = ~0.286
		MentalStability:  70,
		LocationBonus:    1.0,
		EquipmentBonus:   1.0,
		ElementAffinity:  1.0,
		Luck:             50,
	}

	prob := rule.CalculateForgingSuccessRate(recipe, input)
	// 0.45 * (2/7) * 0.7 = 0.45 * 0.2857 * 0.7 ≈ 0.09
	assert.InDelta(t, 0.09, prob, 0.01)
}

func TestCalculateForgingSuccessRate_HallBonus(t *testing.T) {
	rule := NewForgingRule()

	recipe := ArtifactRecipe{
		Name:          "basic_item",
		RequiredLevel: 3,
		Difficulty:    1.0,
	}

	input := ForgingInput{
		ArtificingLevel:  3,
		MentalStability:  100,
		LocationBonus:    1.5, // forging hall
		EquipmentBonus:   1.0,
		ElementAffinity:  1.0,
		Luck:             50,
	}

	prob := rule.CalculateForgingSuccessRate(recipe, input)
	// 0.45 * 1.0 * 1.0 * 1.5 = 0.675
	assert.InDelta(t, 0.675, prob, 0.001)
}

func TestCalculateForgingSuccessRate_Difficulty(t *testing.T) {
	rule := NewForgingRule()

	recipe := ArtifactRecipe{
		Name:          "legendary_artifact",
		RequiredLevel: 5,
		Difficulty:    2.5,
	}

	input := ForgingInput{
		ArtificingLevel:  5,
		MentalStability:  100,
		LocationBonus:    1.0,
		EquipmentBonus:   1.0,
		ElementAffinity:  1.0,
		Luck:             50,
	}

	prob := rule.CalculateForgingSuccessRate(recipe, input)
	// 0.45 * 1.0 / 2.5 = 0.18
	assert.InDelta(t, 0.18, prob, 0.001)
}

func TestGenerateArtifactGrade_HighSuccess(t *testing.T) {
	rule := NewForgingRule()

	alwaysHigh := func() float64 { return 0.01 }
	grade := rule.GenerateArtifactGrade(0.90, 4.0, alwaysHigh)
	assert.Equal(t, ArtifactGradeDivine, grade)
}

func TestGenerateArtifactGrade_MidSuccess(t *testing.T) {
	rule := NewForgingRule()

	mid := func() float64 { return 0.3 }
	grade := rule.GenerateArtifactGrade(0.50, 3.0, mid)
	// expected = 3 + (0.50-0.45)*4 = 3.2
	// r=0.3 < 0.4 → earthly
	assert.Equal(t, ArtifactGradeEarthly, grade)
}

func TestResolveForging_Success(t *testing.T) {
	rule := NewForgingRule()
	deterministic := func() float64 { return 0.1 }

	recipe := ArtifactRecipe{
		Name:          "spirit_sword",
		RequiredLevel: 2,
		Difficulty:    1.0,
	}

	input := ForgingInput{
		ArtificingLevel:  4,
		MentalStability:  100,
	}

	result := rule.ResolveForging(recipe, input, deterministic)
	assert.True(t, result.Success)
	assert.False(t, result.MaterialBroken)
	assert.NotEmpty(t, result.Grade)
	assert.Contains(t, result.Message, "成功")
}

func TestResolveForging_Failure(t *testing.T) {
	rule := NewForgingRule()
	deterministic := func() float64 { return 0.99 }

	recipe := ArtifactRecipe{
		Name:          "divine_armor",
		RequiredLevel: 8,
		Difficulty:    1.0,
	}

	input := ForgingInput{
		ArtificingLevel:  1,
		MentalStability:  30,
	}

	result := rule.ResolveForging(recipe, input, deterministic)
	assert.False(t, result.Success)
}

func TestForgingMaterialLoss(t *testing.T) {
	rule := NewForgingRule()

	// Material broken: total loss
	loss := rule.CalculateMaterialLoss(true)
	assert.Equal(t, 1.0, loss)

	// Recoverable: partial loss
	loss = rule.CalculateMaterialLoss(false)
	assert.Equal(t, 0.3, loss)
}

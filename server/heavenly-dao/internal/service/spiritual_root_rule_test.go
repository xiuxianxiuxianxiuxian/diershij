package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSpiritualRootRule(t *testing.T) {
	rule := NewSpiritualRootRule()
	assert.NotNil(t, rule)
	assert.Equal(t, 0.01, rule.baseAwakeningRate)
	assert.Equal(t, 0.10, rule.mutationRate)
}

func TestCalculateAwakeningRate_NotMortal(t *testing.T) {
	rule := NewSpiritualRootRule()

	input := AwakeningInput{
		Age:              10,
		IsMortal:         false,
	}

	rate := rule.CalculateAwakeningRate(input)
	assert.Equal(t, 0.0, rate)
}

func TestCalculateAwakeningRate_PeakAge(t *testing.T) {
	rule := NewSpiritualRootRule()

	input := AwakeningInput{
		Age:              12,
		IsMortal:         true,
		Luck:             50,
		FamilyBloodline:  0.5,
		SpiritualDensity: 0.5,
	}

	rate := rule.CalculateAwakeningRate(input)
	// 0.01 * 1.0 * 1.0 * 0.5 * 0.5 = 0.0025
	assert.InDelta(t, 0.0025, rate, 0.0001)
}

func TestCalculateAwakeningRate_YoungChild(t *testing.T) {
	rule := NewSpiritualRootRule()

	input := AwakeningInput{
		Age:              3,
		IsMortal:         true,
		Luck:             50,
		FamilyBloodline:  0.5,
		SpiritualDensity: 0.5,
	}

	rate := rule.CalculateAwakeningRate(input)
	// age_factor = 3/6 * 0.5 = 0.25
	// 0.01 * 0.25 * 1.0 * 0.5 * 0.5 = 0.000625
	assert.InDelta(t, 0.000625, rate, 0.0001)
}

func TestCalculateAwakeningRate_OldPerson(t *testing.T) {
	rule := NewSpiritualRootRule()

	input := AwakeningInput{
		Age:              60,
		IsMortal:         true,
		Luck:             50,
		FamilyBloodline:  0.5,
		SpiritualDensity: 0.5,
	}

	rate := rule.CalculateAwakeningRate(input)
	// age_factor = 0.5 - (60-30)/70*0.5 = 0.5 - 0.214 = 0.286
	// 0.01 * 0.286 * 1.0 * 0.5 * 0.5 ≈ 0.000714
	assert.InDelta(t, 0.000714, rate, 0.0001)
}

func TestCalculateAwakeningRate_HighLuck(t *testing.T) {
	rule := NewSpiritualRootRule()

	input := AwakeningInput{
		Age:              12,
		IsMortal:         true,
		Luck:             100,
		FamilyBloodline:  0.5,
		SpiritualDensity: 0.5,
	}

	rate := rule.CalculateAwakeningRate(input)
	// luck_factor = 0.5 + 100/100 = 1.5
	// 0.01 * 1.0 * 1.5 * 0.5 * 0.5 = 0.00375
	assert.InDelta(t, 0.00375, rate, 0.0001)
}

func TestCalculateAwakeningRate_ZeroBloodline(t *testing.T) {
	rule := NewSpiritualRootRule()

	input := AwakeningInput{
		Age:              12,
		IsMortal:         true,
		Luck:             50,
		FamilyBloodline:  0.0,
		SpiritualDensity: 0.5,
	}

	rate := rule.CalculateAwakeningRate(input)
	// bloodline defaults to 0.1
	// 0.01 * 1.0 * 1.0 * 0.1 * 0.5 = 0.0005
	assert.InDelta(t, 0.0005, rate, 0.0001)
}

func TestCalculateAwakeningRate_MinClamp(t *testing.T) {
	rule := NewSpiritualRootRule()

	input := AwakeningInput{
		Age:              100,
		IsMortal:         true,
		Luck:             0,
		FamilyBloodline:  0.0,
		SpiritualDensity: 0.0,
	}

	rate := rule.CalculateAwakeningRate(input)
	// Should be clamped to minimum 0.0001
	assert.InDelta(t, 0.0001, rate, 0.0001)
}

func TestGenerateMutatedElement_NoMutation(t *testing.T) {
	rule := NewSpiritualRootRule()
	alwaysNoMutation := func() float64 { return 0.5 } // > 0.10

	result := rule.GenerateMutatedElement("water", alwaysNoMutation)
	assert.Empty(t, result)
}

func TestGenerateMutatedElement_Mutation(t *testing.T) {
	rule := NewSpiritualRootRule()
	alwaysMutate := func() float64 { return 0.01 } // < 0.10

	// Water → ice
	result := rule.GenerateMutatedElement("water", alwaysMutate)
	assert.Equal(t, "ice", result)

	// Wind → ice or lightning
	result = rule.GenerateMutatedElement("wind", func() float64 { return 0.01 })
	assert.Contains(t, []string{"ice", "lightning"}, result)
}

func TestGenerateMutatedElement_NoMutationPath(t *testing.T) {
	rule := NewSpiritualRootRule()
	alwaysMutate := func() float64 { return 0.01 }

	// Metal has no mutation path
	result := rule.GenerateMutatedElement("metal", alwaysMutate)
	assert.Empty(t, result)
}

func TestCalculateRootQuality_SingleElement(t *testing.T) {
	rule := NewSpiritualRootRule()

	// Pure single element
	quality := rule.CalculateRootQuality("fire", 100)
	assert.InDelta(t, 1.0, quality, 0.001)

	// Half purity
	quality = rule.CalculateRootQuality("fire", 50)
	assert.InDelta(t, 0.5, quality, 0.001)
}

func TestCalculateRootQuality_MutatedElement(t *testing.T) {
	rule := NewSpiritualRootRule()

	// Mutated element has 1.2x bonus
	quality := rule.CalculateRootQuality("ice", 100)
	assert.InDelta(t, 1.2, quality, 0.001)
}

func TestCalculateRootQuality_SupremeElement(t *testing.T) {
	rule := NewSpiritualRootRule()

	quality := rule.CalculateRootQuality("heavenly", 100)
	assert.InDelta(t, 1.5, quality, 0.001)
}

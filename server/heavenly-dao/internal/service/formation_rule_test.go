package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFormationRule(t *testing.T) {
	rule := NewFormationRule()
	assert.NotNil(t, rule)
	assert.Equal(t, 0.10, rule.baseBreakRate)
	assert.Equal(t, 0.01, rule.minBreakRate)
	assert.Equal(t, 0.80, rule.maxBreakRate)
}

func TestCalculateFormationPower_Baseline(t *testing.T) {
	rule := NewFormationRule()

	formation := Formation{
		Name:          "simple_array",
		RequiredLevel: 3,
		BasePower:     100.0,
		Difficulty:    1.0,
		Nodes: []FormationNode{
			{ElementType: "fire", Power: 10.0, IsCore: true},
			{ElementType: "fire", Power: 5.0, IsCore: false},
		},
	}

	input := FormationInput{
		FormationLevel:   3, // skill = 1.0
		MentalStability:  100,
		LocationBonus:    1.0,
		EquipmentBonus:   1.0,
		ElementAffinity:  1.0,
		SpiritStoneCount: 0,
		Luck:             50,
	}

	power := rule.CalculateFormationPower(formation, input)
	// node_power = 10*1.5 + 5 = 20
	// 100 * 20 * 1.0 * 1.0 * 1.0 * 1.0 * 1.0 * 1.0 * 1.0 = 2000
	assert.InDelta(t, 2000.0, power, 0.001)
}

func TestCalculateFormationPower_NoNodes(t *testing.T) {
	rule := NewFormationRule()

	formation := Formation{
		Name:          "empty",
		RequiredLevel: 1,
		BasePower:     50.0,
		Nodes:         []FormationNode{},
	}

	input := FormationInput{
		FormationLevel:  1,
		MentalStability: 100,
		Luck:            50,
	}

	power := rule.CalculateFormationPower(formation, input)
	// node_power defaults to 1.0
	// 50 * 1.0 * 1.0 * 1.0 * 1.0 * 1.0 * 1.0 * 1.0 = 50
	assert.InDelta(t, 50.0, power, 0.001)
}

func TestCalculateFormationPower_StoneBonus(t *testing.T) {
	rule := NewFormationRule()

	formation := Formation{
		Name:          "powered",
		RequiredLevel: 3,
		BasePower:     100.0,
		Nodes:         []FormationNode{{Power: 1.0, IsCore: false}},
	}

	input := FormationInput{
		FormationLevel:   3,
		MentalStability:  100,
		SpiritStoneCount: 10, // +50% bonus
		Luck:             50,
	}

	power := rule.CalculateFormationPower(formation, input)
	// 100 * 1.0 * 1.0 * 1.0 * 1.0 * 1.0 * 1.0 * 1.5 * 1.0 = 150
	assert.InDelta(t, 150.0, power, 0.001)
}

func TestCalculateFormationPower_MaxStoneBonus(t *testing.T) {
	rule := NewFormationRule()

	formation := Formation{
		Name:          "overpowered",
		RequiredLevel: 3,
		BasePower:     100.0,
		Nodes:         []FormationNode{{Power: 1.0}},
	}

	input := FormationInput{
		FormationLevel:   3,
		MentalStability:  100,
		SpiritStoneCount: 100, // should cap at 50%
		Luck:             50,
	}

	power := rule.CalculateFormationPower(formation, input)
	// Stone bonus capped at 1.5
	assert.InDelta(t, 150.0, power, 0.001)
}

func TestCalculateBreakFormationRate_WeakAttack(t *testing.T) {
	rule := NewFormationRule()

	formation := Formation{
		Name:          "strong_ward",
		RequiredLevel: 5,
		BasePower:     100.0,
		Nodes:         []FormationNode{{Power: 5.0, IsCore: true}},
	}

	input := FormationInput{
		FormationLevel:  5,
		MentalStability: 100,
	}

	breakRate := rule.CalculateBreakFormationRate(50.0, formation, input)
	// Formation power = 100 * (5*1.5) = 750
	// ratio = 50/750 = 0.0667
	// break_rate = 0.10 * 0.0667^2 ≈ 0.000444 → clamped to 0.01
	assert.InDelta(t, 0.01, breakRate, 0.001)
}

func TestCalculateBreakFormationRate_StrongAttack(t *testing.T) {
	rule := NewFormationRule()

	formation := Formation{
		Name:          "weak_ward",
		RequiredLevel: 2,
		BasePower:     50.0,
		Nodes:         []FormationNode{{Power: 1.0}},
	}

	input := FormationInput{
		FormationLevel:  2,
		MentalStability: 100,
	}

	breakRate := rule.CalculateBreakFormationRate(200.0, formation, input)
	// Formation power = 50 * 1.0 = 50
	// ratio = 200/50 = 4.0
	// break_rate = 0.10 * 16 = 1.6 → clamped to 0.80
	assert.InDelta(t, 0.80, breakRate, 0.001)
}

func TestResolveFormation_Holds(t *testing.T) {
	rule := NewFormationRule()

	formation := Formation{
		Name:          "stable_formation",
		RequiredLevel: 3,
		BasePower:     100.0,
		Nodes:         []FormationNode{{Power: 2.0, IsCore: true}},
	}

	input := FormationInput{
		FormationLevel:  3,
		MentalStability: 100,
	}

	result := rule.ResolveFormation(formation, input, 50.0)
	assert.True(t, result.Success)
	assert.Contains(t, result.Message, "稳固")
}

func TestResolveFormation_Breaking(t *testing.T) {
	rule := NewFormationRule()

	formation := Formation{
		Name:          "failing_formation",
		RequiredLevel: 2,
		BasePower:     10.0,
		Nodes:         []FormationNode{{Power: 1.0}},
	}

	input := FormationInput{
		FormationLevel:  1,
		MentalStability: 30,
	}

	result := rule.ResolveFormation(formation, input, 500.0)
	assert.False(t, result.Success)
	assert.Contains(t, result.Message, "被破")
}

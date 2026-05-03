package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func makeCraftInput(craftType CraftType) CraftInput {
	return CraftInput{
		EntityID:        "crafter_1",
		CraftType:       craftType,
		MentalStability: 80,
		LocationBonus:   1.0,
		EquipmentBonus:  1.0,
		ElementAffinity: 1.0,
		Luck:            50,
		SkillLevel:      3,
		SpiritStoneCount: 5,
	}
}

func TestNewCraftOperation(t *testing.T) {
	op := NewCraftOperation(time.Minute)
	assert.NotNil(t, op)
}

func TestExecuteCraft_Alchemy(t *testing.T) {
	op := NewCraftOperation(time.Minute)

	input := makeCraftInput(CraftTypeAlchemy)
	now := time.Now()
	deterministic := func() float64 { return 0.1 } // high success chance

	result, err := op.ExecuteCraft(input, now, deterministic)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, CraftTypeAlchemy, result.CraftType)
	assert.NotEmpty(t, result.Quality)
	assert.Contains(t, result.Message, "炼丹成功")
}

func TestExecuteCraft_Forging(t *testing.T) {
	op := NewCraftOperation(time.Minute)

	input := makeCraftInput(CraftTypeForging)
	now := time.Now()
	deterministic := func() float64 { return 0.1 }

	result, err := op.ExecuteCraft(input, now, deterministic)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, CraftTypeForging, result.CraftType)
	assert.NotEmpty(t, result.Quality)
	assert.Contains(t, result.Message, "炼器成功")
}

func TestExecuteCraft_Formation(t *testing.T) {
	op := NewCraftOperation(time.Minute)

	input := makeCraftInput(CraftTypeFormation)
	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	result, err := op.ExecuteCraft(input, now, deterministic)
	assert.NoError(t, err)
	// With skill level 3, threshold = 150
	// Power = 100 * (10*1.5 + 5 + 5) * 1.0 * 0.8 * 1.0 * 1.0 * 1.0 * 1.25 * 1.0
	// = 100 * 25 * 0.8 * 1.25 = 2500 > 150
	assert.True(t, result.Success)
	assert.Equal(t, CraftTypeFormation, result.CraftType)
	assert.Contains(t, result.Message, "阵法")
}

func TestExecuteCraft_UnknownType(t *testing.T) {
	op := NewCraftOperation(time.Minute)

	input := makeCraftInput(CraftType("unknown"))
	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	_, err := op.ExecuteCraft(input, now, deterministic)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown craft type")
}

func TestExecuteCraft_Cooldown(t *testing.T) {
	op := NewCraftOperation(time.Hour)

	input := makeCraftInput(CraftTypeAlchemy)
	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	result, err := op.ExecuteCraft(input, now, deterministic)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	_, err = op.ExecuteCraft(input, now.Add(30*time.Minute), deterministic)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cooldown")
}

func TestExecuteCraft_LowSkillFailure(t *testing.T) {
	op := NewCraftOperation(time.Minute)

	input := makeCraftInput(CraftTypeAlchemy)
	input.SkillLevel = 1
	input.MentalStability = 20
	input.Luck = 0

	now := time.Now()
	deterministic := func() float64 { return 0.9 } // very unlucky

	result, err := op.ExecuteCraft(input, now, deterministic)
	assert.NoError(t, err)
	// Low skill + bad luck should likely fail
	assert.False(t, result.Success)
}

func TestGetCooldownRemaining_Craft(t *testing.T) {
	op := NewCraftOperation(time.Hour)

	input := makeCraftInput(CraftTypeAlchemy)
	now := time.Now()
	_, _ = op.ExecuteCraft(input, now, func() float64 { return 0.5 })

	remaining := op.GetCooldownRemaining("crafter_1", now.Add(30*time.Minute))
	assert.True(t, remaining > 0)
}

package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func makeResourceNode() ResourceNode {
	return ResourceNode{
		ID:           "node_1",
		Name:         "千年灵芝",
		Type:         "herb",
		RegionID:     "region_1",
		Level:        3,
		Quantity:     10,
		MaxQuantity:  20,
		RegenRate:    0.5,
		RegenTimer:   0,
		Quality:      "common",
		GatherSkill:  "herbalism",
		MinSkillLevel: 3,
	}
}

func makeGatherInput() GatherInput {
	return GatherInput{
		EntityID:    "gatherer_1",
		NodeID:      "node_1",
		GatherSkill: "herbalism",
		SkillLevel:  5,
		Luck:        50,
		ToolsBonus:  0.1,
	}
}

func TestNewGatherOperation(t *testing.T) {
	op := NewGatherOperation(time.Minute)
	assert.NotNil(t, op)
}

func TestExecuteGather_Success(t *testing.T) {
	op := NewGatherOperation(time.Minute)

	input := makeGatherInput()
	node := makeResourceNode()
	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	result, err := op.ExecuteGather(input, node, now, deterministic)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.True(t, result.Amount > 0)
	assert.True(t, result.Amount <= node.Quantity)
	assert.Contains(t, result.Message, "采集了")
}

func TestExecuteGather_DepletedNode(t *testing.T) {
	op := NewGatherOperation(time.Minute)

	input := makeGatherInput()
	node := makeResourceNode()
	node.Quantity = 0
	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	result, err := op.ExecuteGather(input, node, now, deterministic)
	assert.NoError(t, err)
	assert.False(t, result.Success)
	assert.True(t, result.NodeState.Depleted)
	assert.Contains(t, result.Message, "枯竭")
}

func TestExecuteGather_WrongSkill(t *testing.T) {
	op := NewGatherOperation(time.Minute)

	input := makeGatherInput()
	input.GatherSkill = "mining" // wrong skill
	node := makeResourceNode()
	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	_, err := op.ExecuteGather(input, node, now, deterministic)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "herbalism")
}

func TestExecuteGather_InsufficientSkill(t *testing.T) {
	op := NewGatherOperation(time.Minute)

	input := makeGatherInput()
	input.SkillLevel = 1 // below min level 3
	node := makeResourceNode()
	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	_, err := op.ExecuteGather(input, node, now, deterministic)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "等级")
}

func TestExecuteGather_Cooldown(t *testing.T) {
	op := NewGatherOperation(time.Hour)

	input := makeGatherInput()
	node := makeResourceNode()
	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	result, err := op.ExecuteGather(input, node, now, deterministic)
	assert.NoError(t, err)
	assert.True(t, result.Success)

	_, err = op.ExecuteGather(input, node, now.Add(30*time.Minute), deterministic)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cooldown")
}

func TestCalculateGatherSuccessRate_HighSkill(t *testing.T) {
	op := NewGatherOperation(time.Minute)

	input := GatherInput{
		SkillLevel:  10,
		Luck:        100,
		ToolsBonus:  0.5,
	}
	node := ResourceNode{MinSkillLevel: 3}

	rate := op.calculateGatherSuccessRate(input, node)
	assert.True(t, rate > 0.7)
	assert.True(t, rate <= 0.95)
}

func TestCalculateGatherSuccessRate_LowSkill(t *testing.T) {
	op := NewGatherOperation(time.Minute)

	input := GatherInput{
		SkillLevel:  3,
		Luck:        0,
		ToolsBonus:  0,
	}
	node := ResourceNode{MinSkillLevel: 5}

	rate := op.calculateGatherSuccessRate(input, node)
	assert.True(t, rate >= 0.20)
	assert.True(t, rate <= 0.95)
}

func TestRegenerateNode(t *testing.T) {
	op := NewGatherOperation(time.Minute)

	node := &ResourceNode{
		Quantity:    5,
		MaxQuantity: 20,
		RegenRate:   2.0,
	}

	op.RegenerateNode(node, 3.0)
	// Should have regenerated 2.0 * 3 = 6 units
	assert.Equal(t, 11, node.Quantity)
}

func TestRegenerateNode_Full(t *testing.T) {
	op := NewGatherOperation(time.Minute)

	node := &ResourceNode{
		Quantity:    20,
		MaxQuantity: 20,
		RegenRate:   2.0,
	}

	op.RegenerateNode(node, 3.0)
	// Already full, should not change
	assert.Equal(t, 20, node.Quantity)
}

func TestRegenerateNode_CapAtMax(t *testing.T) {
	op := NewGatherOperation(time.Minute)

	node := &ResourceNode{
		Quantity:    18,
		MaxQuantity: 20,
		RegenRate:   2.0,
	}

	op.RegenerateNode(node, 10.0)
	// Would regenerate 20, but capped at max
	assert.Equal(t, 20, node.Quantity)
}

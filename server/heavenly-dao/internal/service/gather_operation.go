package service

import (
	"fmt"
	"time"
)

// GatherOperation handles resource gathering operations.
type GatherOperation struct {
	cooldownMap    map[string]time.Time
	cooldownPeriod time.Duration
}

// NewGatherOperation creates a new GatherOperation.
func NewGatherOperation(cooldownPeriod time.Duration) *GatherOperation {
	return &GatherOperation{
		cooldownMap:    make(map[string]time.Time),
		cooldownPeriod: cooldownPeriod,
	}
}

// ResourceNode represents a gatherable resource in the world.
type ResourceNode struct {
	ID           string
	Name         string
	Type         string // herb / ore / crystal / wood / water
	RegionID     string
	Level        int
	Quantity     int
	MaxQuantity  int
	RegenRate    float64 // units per hour
	RegenTimer   float64 // hours until next regen
	Quality      string  // common / rare / epic / legendary
	GatherSkill  string  // required gathering skill type
	MinSkillLevel int    // minimum skill level to gather
}

// GatherInput holds inputs for a gathering operation.
type GatherInput struct {
	EntityID    string
	NodeID      string
	GatherSkill string
	SkillLevel  int
	Luck        int
	ToolsBonus  float64
}

// GatherResult holds the outcome of a gathering operation.
type GatherResult struct {
	Success    bool
	Amount     int
	Quality    string
	BonusDrop  bool
	NodeState  NodeState
	Message    string
}

// NodeState represents the state of a resource node after gathering.
type NodeState struct {
	Remaining  int
	Depleted   bool
	RegenHours float64
}

// ExecuteGather attempts to gather resources from a node.
func (op *GatherOperation) ExecuteGather(input GatherInput, node ResourceNode, now time.Time, randFloat func() float64) (*GatherResult, error) {
	// Check cooldown
	if lastTime, ok := op.cooldownMap[input.EntityID]; ok {
		if now.Sub(lastTime) < op.cooldownPeriod {
			return nil, fmt.Errorf("gather cooldown: %v remaining", op.cooldownPeriod-now.Sub(lastTime))
		}
	}

	// Check if node is depleted
	if node.Quantity <= 0 {
		return &GatherResult{
			Success: false,
			Message: "资源已枯竭",
			NodeState: NodeState{
				Remaining:  0,
				Depleted:   true,
				RegenHours: node.RegenTimer,
			},
		}, nil
	}

	// Check skill requirement
	if input.GatherSkill != node.GatherSkill && node.GatherSkill != "" {
		return nil, fmt.Errorf("需要 %s 技能", node.GatherSkill)
	}
	if input.SkillLevel < node.MinSkillLevel {
		return nil, fmt.Errorf("需要 %s 等级 %d", node.GatherSkill, node.MinSkillLevel)
	}

	// Calculate success rate
	successRate := op.calculateGatherSuccessRate(input, node)
	if randFloat() >= successRate {
		op.cooldownMap[input.EntityID] = now
		return &GatherResult{
			Success: false,
			Message: "采集失败",
			NodeState: NodeState{
				Remaining:  node.Quantity,
				Depleted:   false,
				RegenHours: node.RegenTimer,
			},
		}, nil
	}

	// Calculate amount gathered
	baseAmount := 1 + int(randFloat()*3)
	skillBonus := input.SkillLevel / node.MinSkillLevel
	if skillBonus > 3 {
		skillBonus = 3
	}
	amount := baseAmount * (1 + skillBonus)

	// Cap at node quantity
	if amount > node.Quantity {
		amount = node.Quantity
	}

	// Determine quality
	quality := node.Quality
	if randFloat() < 0.1 {
		quality = "rare"
		if randFloat() < 0.1 {
			quality = "epic"
		}
	}

	// Bonus drop chance
	bonusDrop := randFloat() < 0.15*(float64(input.Luck)/50.0)*(1.0+input.ToolsBonus)

	// Update node state
	remaining := node.Quantity - amount
	if remaining < 0 {
		remaining = 0
	}
	depleted := remaining <= 0

	op.cooldownMap[input.EntityID] = now

	return &GatherResult{
		Success:   true,
		Amount:    amount,
		Quality:   quality,
		BonusDrop: bonusDrop,
		NodeState: NodeState{
			Remaining:  remaining,
			Depleted:   depleted,
			RegenHours: node.RegenTimer,
		},
		Message: fmt.Sprintf("采集了 %d 个 %s", amount, node.Name),
	}, nil
}

func (op *GatherOperation) calculateGatherSuccessRate(input GatherInput, node ResourceNode) float64 {
	// Base rate: 70%
	baseRate := 0.70

	// Skill level ratio
	skillRatio := float64(input.SkillLevel) / float64(node.MinSkillLevel)
	if skillRatio <= 0 {
		skillRatio = 0.1
	}
	skillFactor := 0.5 + skillRatio*0.5
	if skillFactor > 1.5 {
		skillFactor = 1.5
	}

	// Tools bonus
	toolsFactor := 1.0 + input.ToolsBonus
	if toolsFactor > 2.0 {
		toolsFactor = 2.0
	}

	// Luck factor
	luckFactor := 0.8 + float64(input.Luck)/200.0

	rate := baseRate * skillFactor * toolsFactor * luckFactor
	if rate > 0.95 {
		rate = 0.95
	}
	if rate < 0.20 {
		rate = 0.20
	}
	return rate
}

// RegenerateNode updates a node's quantity based on its regen rate.
func (op *GatherOperation) RegenerateNode(node *ResourceNode, elapsedHours float64) {
	if node.Quantity >= node.MaxQuantity {
		return
	}

	regenAmount := int(node.RegenRate * elapsedHours)
	if regenAmount <= 0 {
		regenAmount = 1
	}

	node.Quantity += regenAmount
	if node.Quantity > node.MaxQuantity {
		node.Quantity = node.MaxQuantity
	}

	node.RegenTimer = 0
}

// GetCooldownRemaining returns remaining cooldown time for an entity.
func (op *GatherOperation) GetCooldownRemaining(entityID string, now time.Time) time.Duration {
	if lastTime, ok := op.cooldownMap[entityID]; ok {
		remaining := op.cooldownPeriod - now.Sub(lastTime)
		if remaining > 0 {
			return remaining
		}
	}
	return 0
}

// ClearCooldown removes cooldown for an entity.
func (op *GatherOperation) ClearCooldown(entityID string) {
	delete(op.cooldownMap, entityID)
}

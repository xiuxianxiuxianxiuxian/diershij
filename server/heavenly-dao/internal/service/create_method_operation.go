package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cultivation-world/shared/types"
)

const (
	createMethodCost = 10000 // premium spirit stones required
)

// CreateMethodOperation handles the create method (功法自创) operation.
type CreateMethodOperation struct {
	cooldownMap    map[string]time.Time
	cooldownPeriod time.Duration
	methodRegistry map[string]*types.CultivationMethod
	creatorRewards map[string]int64 // creatorID -> reward amount
}

// NewCreateMethodOperation creates a new CreateMethodOperation.
func NewCreateMethodOperation(cooldownPeriod time.Duration) *CreateMethodOperation {
	return &CreateMethodOperation{
		cooldownMap:    make(map[string]time.Time),
		cooldownPeriod: cooldownPeriod,
		methodRegistry: make(map[string]*types.CultivationMethod),
		creatorRewards: make(map[string]int64),
	}
}

// CreateMethodInput holds inputs for creating a new cultivation method.
type CreateMethodInput struct {
	EntityID        string
	EntityName      string
	MethodName      string
	MethodCategory  string // 主修功法/秘术/身法/神识/辅助/生活
	ElementType     string
	RequiredRealm   string
	PremiumStones   int64
	Comprehension   int // comprehension attribute (0-100)
	MentalStability int // mental stability (0-100)
	DaoHeart        int // dao heart level (0-100)
}

// CreateMethodResult holds the outcome of a create method operation.
type CreateMethodResult struct {
	Success       bool
	Method        *types.CultivationMethod
	Cost          int64
	Quality       float64
	CreatorReward int64
	Message       string
}

// ExecuteCreateMethod attempts to create a new cultivation method.
func (op *CreateMethodOperation) ExecuteCreateMethod(input CreateMethodInput, now time.Time) (*CreateMethodResult, error) {
	// Check cooldown
	if lastTime, ok := op.cooldownMap[input.EntityID]; ok {
		if now.Sub(lastTime) < op.cooldownPeriod {
			return nil, fmt.Errorf("create method cooldown: %v remaining", op.cooldownPeriod-now.Sub(lastTime))
		}
	}

	// Validate premium spirit stones
	if input.PremiumStones < createMethodCost {
		return nil, fmt.Errorf("需要 %d 极品灵石，当前有 %d", createMethodCost, input.PremiumStones)
	}

	// Validate method name uniqueness
	if _, exists := op.methodRegistry[input.MethodName]; exists {
		return nil, fmt.Errorf("功法 '%s' 已存在", input.MethodName)
	}

	// Calculate method quality based on creator attributes
	quality := op.calculateMethodQuality(input)

	// Generate the method
	method := op.generateMethod(input, quality)

	// Register the method
	op.methodRegistry[input.MethodName] = method

	op.cooldownMap[input.EntityID] = now

	return &CreateMethodResult{
		Success:       true,
		Method:        method,
		Cost:          createMethodCost,
		Quality:       quality,
		CreatorReward: 0,
		Message:       fmt.Sprintf("成功创功 '%s'，品质：%.1f", input.MethodName, quality),
	}, nil
}

func (op *CreateMethodOperation) calculateMethodQuality(input CreateMethodInput) float64 {
	// Base quality from comprehension
	comprehensionFactor := float64(input.Comprehension) / 100.0

	// Mental stability factor
	mentalFactor := float64(input.MentalStability) / 100.0

	// Dao heart factor (most important)
	daoHeartFactor := float64(input.DaoHeart) / 100.0

	// Random variance
	variance := (rand.Float64() - 0.5) * 0.2 // +/-10%

	quality := (comprehensionFactor*0.3 + mentalFactor*0.2 + daoHeartFactor*0.5) * 10.0
	quality += variance * 10.0

	quality = maxFloat(1.0, minFloat(10.0, quality))
	return quality
}

func (op *CreateMethodOperation) generateMethod(input CreateMethodInput, quality float64) *types.CultivationMethod {
	// Map quality to rank
	var rank string
	switch {
	case quality >= 9.0:
		rank = "天极品"
	case quality >= 8.0:
		rank = "天上品"
	case quality >= 7.0:
		rank = "地极品"
	case quality >= 6.0:
		rank = "地上品"
	case quality >= 5.0:
		rank = "玄极品"
	case quality >= 4.0:
		rank = "玄上品"
	case quality >= 3.0:
		rank = "黄极品"
	default:
		rank = "黄上品"
	}

	// Calculate bonuses based on quality and category
	attacks := make(map[string]float64)
	defenses := make(map[string]float64)
	utilities := make(map[string]float64)

	switch input.MethodCategory {
	case "主修功法":
		utilities["cultivation_speed"] = quality * 0.05
	case "秘术":
		attacks["damage_bonus"] = quality * 0.05
	case "身法":
		utilities["dodge_bonus"] = quality * 0.03
		attacks["speed_bonus"] = quality * 0.03
	case "神识":
		defenses["mental_resist"] = quality * 0.05
	case "辅助":
		utilities["alchemy_success"] = quality * 0.02
		utilities["forging_success"] = quality * 0.02
	}

	return &types.CultivationMethod{
		Name:                  input.MethodName,
		CreatorID:             input.EntityID,
		Category:              input.MethodCategory,
		ElementAffinity:       input.ElementType,
		Rank:                  rank,
		RealmRequirement:      input.RequiredRealm,
		PowerScore:            int(quality * 10),
		Potential:             int(quality * 15),
		Popularity:            0,
		AttackBonuses:         attacks,
		DefenseBonuses:        defenses,
		UtilityBonuses:        utilities,
		CanModify:             true,
		Complexity:            int(quality * 20),
	}
}

// TrackLearner records a new learner and potentially rewards the creator.
func (op *CreateMethodOperation) TrackLearner(methodName, learnerID string) (creatorReward int64, err error) {
	method, exists := op.methodRegistry[methodName]
	if !exists {
		return 0, fmt.Errorf("功法 '%s' 不存在", methodName)
	}

	method.Popularity++
	learnerCount := method.Popularity

	// Creator gets reward at milestones: 5-10 premium stones at 10 learners
	if learnerCount == 10 {
		reward := int64(5 + rand.Intn(6)) // 5-10 premium stones
		op.creatorRewards[method.CreatorID] = reward
		return reward, nil
	}

	return 0, nil
}

// GetMethod retrieves a registered method by name.
func (op *CreateMethodOperation) GetMethod(name string) (*types.CultivationMethod, bool) {
	m, ok := op.methodRegistry[name]
	return m, ok
}

// ListMethods returns all registered methods.
func (op *CreateMethodOperation) ListMethods() []*types.CultivationMethod {
	methods := make([]*types.CultivationMethod, 0, len(op.methodRegistry))
	for _, m := range op.methodRegistry {
		methods = append(methods, m)
	}
	return methods
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

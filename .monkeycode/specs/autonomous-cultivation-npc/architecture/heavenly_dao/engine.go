package heavenlydao

import (
	"fmt"
	"math"
)

type StateRepository interface {
	AddKarma(entityID EntityID, delta float64) error
	SaveOperationResult(actorID EntityID, actionType ActionType, result map[string]any) error
}

type HeavenlyDaoEngine interface {
	CalculateKarma(ctx *RuleContext) (KarmaResult, error)
	ApplyKarma(ctx *RuleContext) (KarmaResult, error)
	CalculateBreakthrough(ctx *RuleContext) (BreakthroughResult, error)
	CalculateDamage(ctx *RuleContext) (DamageResult, error)
	CalculateTradeTax(amount float64) float64
	CanCreateMethod(ctx *RuleContext, requiredRealm int, requiredComprehension float64) error
}

type Engine struct {
	config *HeavenlyDaoConfig
	repo   StateRepository
	bus    EventBus
}

func NewEngine(config *HeavenlyDaoConfig, repo StateRepository, bus EventBus) *Engine {
	return &Engine{config: config, repo: repo, bus: bus}
}

func (e *Engine) CalculateTradeTax(amount float64) float64 {
	if amount > 1_000_000 {
		return amount * 0.015
	}
	if amount > 100_000 {
		return amount * 0.012
	}
	return amount * 0.01
}

func (e *Engine) CanCreateMethod(ctx *RuleContext, requiredRealm int, requiredComprehension float64) error {
	if ctx.ActorState == nil {
		return fmt.Errorf("missing actor state")
	}
	if ctx.ActorState.SpiritStoneBalances.PremiumGrade < 10000 {
		return fmt.Errorf("自创功法需要10000极品灵石")
	}
	if ctx.ActorState.RealmLevel < requiredRealm {
		return fmt.Errorf("境界不足，无法自创功法")
	}
	if ctx.ActorState.Comprehension < requiredComprehension {
		return fmt.Errorf("悟性不足，无法自创功法")
	}
	return nil
}

func clamp(v, low, high float64) float64 {
	return math.Max(low, math.Min(high, v))
}

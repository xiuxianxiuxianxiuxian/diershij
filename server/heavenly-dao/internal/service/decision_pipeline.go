package service

import (
	"context"
	"fmt"
	"time"
)

// DecisionCycleType represents the type of decision cycle.
type DecisionCycleType string

const (
	CycleFast DecisionCycleType = "fast" // behavior tree (1-5s)
	CycleSlow DecisionCycleType = "slow" // LLM (30-120s)
)

// DecisionCycleConfig holds configuration for the decision cycle.
type DecisionCycleConfig struct {
	// Fast cycle interval (behavior tree evaluation)
	FastInterval time.Duration

	// Slow cycle interval (LLM decision)
	SlowInterval time.Duration

	// Context building timeout
	ContextTimeout time.Duration

	// Action execution timeout
	ActionTimeout time.Duration
}

// DefaultDecisionCycleConfig returns the default configuration.
func DefaultDecisionCycleConfig() *DecisionCycleConfig {
	return &DecisionCycleConfig{
		FastInterval:   3 * time.Second,
		SlowInterval:   60 * time.Second,
		ContextTimeout: 5 * time.Second,
		ActionTimeout:  10 * time.Second,
	}
}

// NPCDecisionPipeline manages the NPC decision making process.
type NPCDecisionPipeline struct {
	config        *DecisionCycleConfig
	behaviorTree  *BehaviorTree
	llmClient     *LLMClient
	templateLib   *TemplateLibrary
	fallback      *DecisionFallback
	lastFastCycle  time.Time
	lastSlowCycle  time.Time
	currentAction string
}

// NewNPCDecisionPipeline creates a new decision pipeline.
func NewNPCDecisionPipeline(cfg *DecisionCycleConfig, tree *BehaviorTree, llm *LLMClient, templates *TemplateLibrary) *NPCDecisionPipeline {
	fallback := NewDecisionFallback(llm, templates, tree)

	return &NPCDecisionPipeline{
		config:       cfg,
		behaviorTree: tree,
		llmClient:    llm,
		templateLib:  templates,
		fallback:     fallback,
	}
}

// DecisionCycleResult holds the result of a decision cycle.
type DecisionCycleResult struct {
	CycleType   DecisionCycleType
	Decision    *DecisionResult
	Source      string
	Duration    time.Duration
	ContextInfo string
}

// RunCycle executes a single decision cycle.
func (p *NPCDecisionPipeline) RunCycle(ctx context.Context, npcCtx *NPCContext, now time.Time) (*DecisionCycleResult, error) {
	// Determine which cycle to run
	cycleType := p.determineCycle(now)

	switch cycleType {
	case CycleFast:
		return p.runFastCycle(ctx, npcCtx, now)
	case CycleSlow:
		return p.runSlowCycle(ctx, npcCtx, now)
	default:
		return nil, fmt.Errorf("未知决策周期类型")
	}
}

func (p *NPCDecisionPipeline) determineCycle(now time.Time) DecisionCycleType {
	timeSinceFast := now.Sub(p.lastFastCycle)
	timeSinceSlow := now.Sub(p.lastSlowCycle)

	// If slow cycle is due, run it
	if timeSinceSlow >= p.config.SlowInterval {
		return CycleSlow
	}

	// Otherwise run fast cycle if due
	if timeSinceFast >= p.config.FastInterval {
		return CycleFast
	}

	// No cycle due yet
	return ""
}

func (p *NPCDecisionPipeline) runFastCycle(ctx context.Context, npcCtx *NPCContext, now time.Time) (*DecisionCycleResult, error) {
	start := time.Now()

	// Build context
	contextInfo := p.buildContext(npcCtx)

	// Evaluate behavior tree
	status := p.behaviorTree.Evaluate(npcCtx)

	p.lastFastCycle = now

	return &DecisionCycleResult{
		CycleType: CycleFast,
		Decision: &DecisionResult{
			Action:     string(status),
			Reasoning:  "行为树快速决策",
			Confidence: 0.5,
		},
		Source:      "behavior_tree",
		Duration:    time.Since(start),
		ContextInfo: contextInfo,
	}, nil
}

func (p *NPCDecisionPipeline) runSlowCycle(ctx context.Context, npcCtx *NPCContext, now time.Time) (*DecisionCycleResult, error) {
	start := time.Now()

	// Build context
	contextInfo := p.buildContext(npcCtx)

	// Generate system prompt
	systemPrompt := SystemPromptTemplate(
		npcCtx.EntityID,
		"neutral",
		"unknown",
	)

	// Generate user prompt from context
	userPrompt := fmt.Sprintf("当前状态: 健康=%d, 气=%.0f, 战斗=%v, 目标=%v",
		npcCtx.Health, npcCtx.Qi, npcCtx.IsInCombat, npcCtx.HasTarget)

	// Use fallback chain to get decision
	result, source := p.fallback.Decide(ctx, npcCtx, systemPrompt, userPrompt)

	p.lastSlowCycle = now

	return &DecisionCycleResult{
		CycleType:   CycleSlow,
		Decision:    result,
		Source:      source,
		Duration:    time.Since(start),
		ContextInfo: contextInfo,
	}, nil
}

// buildContext creates a context string from the NPC state.
func (p *NPCDecisionPipeline) buildContext(npcCtx *NPCContext) string {
	return fmt.Sprintf("Entity=%s, Health=%d, Qi=%.0f/%.0f, Combat=%v, Danger=%v, Target=%v",
		npcCtx.EntityID,
		npcCtx.Health,
		npcCtx.Qi,
		npcCtx.MaxQi,
		npcCtx.IsInCombat,
		npcCtx.IsInDanger,
		npcCtx.HasTarget,
	)
}

// ShouldTakeAction checks if the NPC should execute the current action.
func (p *NPCDecisionPipeline) ShouldTakeAction(result *DecisionCycleResult) bool {
	if result == nil || result.Decision == nil {
		return false
	}

	// Only take action if confidence is high enough
	return result.Decision.Confidence >= 0.3
}

// ExecuteAction executes the decision action.
func (p *NPCDecisionPipeline) ExecuteAction(ctx context.Context, npcCtx *NPCContext, decision *DecisionResult) error {
	if decision == nil {
		return fmt.Errorf("决策为空")
	}

	actionCtx, cancel := context.WithTimeout(ctx, p.config.ActionTimeout)
	defer cancel()

	// Execute based on action type
	switch decision.Action {
	case "cultivate":
		return p.executeCultivate(actionCtx, npcCtx, decision)
	case "breakthrough":
		return p.executeBreakthrough(actionCtx, npcCtx, decision)
	case "explore":
		return p.executeExplore(actionCtx, npcCtx, decision)
	case "combat":
		return p.executeCombat(actionCtx, npcCtx, decision)
	case "flee":
		return p.executeFlee(actionCtx, npcCtx, decision)
	case "rest":
		return p.executeRest(actionCtx, npcCtx, decision)
	case "meditate":
		return p.executeMeditate(actionCtx, npcCtx, decision)
	default:
		npcCtx.Log(fmt.Sprintf("未知动作: %s", decision.Action))
		return nil
	}
}

func (p *NPCDecisionPipeline) executeCultivate(ctx context.Context, npcCtx *NPCContext, decision *DecisionResult) error {
	npcCtx.CurrentAction = "cultivating"
	npcCtx.Log("开始修炼")
	return nil
}

func (p *NPCDecisionPipeline) executeBreakthrough(ctx context.Context, npcCtx *NPCContext, decision *DecisionResult) error {
	npcCtx.CurrentAction = "breaking_through"
	npcCtx.Log("尝试突破")
	return nil
}

func (p *NPCDecisionPipeline) executeExplore(ctx context.Context, npcCtx *NPCContext, decision *DecisionResult) error {
	npcCtx.CurrentAction = "exploring"
	npcCtx.Log("开始探索")
	return nil
}

func (p *NPCDecisionPipeline) executeCombat(ctx context.Context, npcCtx *NPCContext, decision *DecisionResult) error {
	npcCtx.CurrentAction = "combat"
	npcCtx.Log("进入战斗")
	return nil
}

func (p *NPCDecisionPipeline) executeFlee(ctx context.Context, npcCtx *NPCContext, decision *DecisionResult) error {
	npcCtx.CurrentAction = "fleeing"
	npcCtx.Log("逃跑")
	return nil
}

func (p *NPCDecisionPipeline) executeRest(ctx context.Context, npcCtx *NPCContext, decision *DecisionResult) error {
	npcCtx.CurrentAction = "resting"
	npcCtx.Health = minInt(npcCtx.Health+20, 100)
	npcCtx.Log("休息恢复")
	return nil
}

func (p *NPCDecisionPipeline) executeMeditate(ctx context.Context, npcCtx *NPCContext, decision *DecisionResult) error {
	npcCtx.CurrentAction = "meditating"
	if npcCtx.Qi+30 < npcCtx.MaxQi {
		npcCtx.Qi += 30
	} else {
		npcCtx.Qi = npcCtx.MaxQi
	}
	npcCtx.Log("冥想恢复气")
	return nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetNextCycleTime returns when the next cycle should run.
func (p *NPCDecisionPipeline) GetNextCycleTime(now time.Time) time.Duration {
	timeToFast := p.config.FastInterval - now.Sub(p.lastFastCycle)
	timeToSlow := p.config.SlowInterval - now.Sub(p.lastSlowCycle)

	if timeToFast < timeToSlow {
		return timeToFast
	}
	return timeToSlow
}

// GetPipelineStats returns statistics about the pipeline.
func (p *NPCDecisionPipeline) GetPipelineStats(now time.Time) map[string]interface{} {
	return map[string]interface{}{
		"last_fast_cycle":  p.lastFastCycle,
		"last_slow_cycle":  p.lastSlowCycle,
		"current_action":   p.currentAction,
		"next_cycle_in":    p.GetNextCycleTime(now),
	}
}

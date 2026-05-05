package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultDecisionCycleConfig(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	assert.NotNil(t, cfg)
	assert.Equal(t, 3*time.Second, cfg.FastInterval)
	assert.Equal(t, 60*time.Second, cfg.SlowInterval)
	assert.Equal(t, 5*time.Second, cfg.ContextTimeout)
	assert.Equal(t, 10*time.Second, cfg.ActionTimeout)
}

func TestNewNPCDecisionPipeline(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	tree := NewBehaviorTree("test", &LeafNode{
		BaseNode: BaseNode{Name: "root"},
		Action:   func(ctx *NPCContext) NodeStatus { return StatusSuccess },
	})
	llm := NewLLMClient("test-key", "https://api.test.com")
	templates := NewTemplateLibrary()

	pipeline := NewNPCDecisionPipeline(cfg, tree, llm, templates)
	assert.NotNil(t, pipeline)
	assert.NotNil(t, pipeline.fallback)
}

func TestRunCycle_FastCycle(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	tree := NewBehaviorTree("test", &LeafNode{
		BaseNode: BaseNode{Name: "root"},
		Action:   func(ctx *NPCContext) NodeStatus { return StatusSuccess },
	})
	llm := NewLLMClient("test-key", "https://api.test.com")
	templates := NewTemplateLibrary()

	pipeline := NewNPCDecisionPipeline(cfg, tree, llm, templates)

	// First run should be slow cycle (no previous cycles)
	ctx := context.Background()
	npcCtx := NewNPCContext("test_npc")
	now := time.Now()

	result, err := pipeline.RunCycle(ctx, npcCtx, now)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	// First run should be slow cycle since lastSlowCycle is zero
	assert.Equal(t, CycleSlow, result.CycleType)

	// Second run within fast interval should not run
	now2 := now.Add(1 * time.Second)
	result2, err := pipeline.RunCycle(ctx, npcCtx, now2)
	assert.NoError(t, err)
	assert.Nil(t, result2) // No cycle due yet
}

func TestRunCycle_SlowCycle(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	tree := NewBehaviorTree("test", &LeafNode{
		BaseNode: BaseNode{Name: "root"},
		Action:   func(ctx *NPCContext) NodeStatus { return StatusSuccess },
	})
	llm := NewLLMClient("test-key", "https://api.test.com")
	templates := NewTemplateLibrary()

	pipeline := NewNPCDecisionPipeline(cfg, tree, llm, templates)

	ctx := context.Background()
	npcCtx := NewNPCContext("test_npc")
	now := time.Now()

	// First run - slow cycle
	result, err := pipeline.RunCycle(ctx, npcCtx, now)
	assert.NoError(t, err)
	assert.Equal(t, CycleSlow, result.CycleType)
	// Source is "behavior_tree" because fallback chain ends at behavior tree
	assert.Equal(t, "behavior_tree", result.Source)
}

func TestRunCycle_FastCycleAfterSlow(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	tree := NewBehaviorTree("test", &LeafNode{
		BaseNode: BaseNode{Name: "root"},
		Action:   func(ctx *NPCContext) NodeStatus { return StatusSuccess },
	})
	llm := NewLLMClient("test-key", "https://api.test.com")
	templates := NewTemplateLibrary()

	pipeline := NewNPCDecisionPipeline(cfg, tree, llm, templates)

	ctx := context.Background()
	npcCtx := NewNPCContext("test_npc")

	// Run slow cycle
	now := time.Now()
	result1, _ := pipeline.RunCycle(ctx, npcCtx, now)
	assert.Equal(t, CycleSlow, result1.CycleType)

	// Wait for fast cycle to be due (but not slow cycle)
	now2 := now.Add(4 * time.Second)
	result2, _ := pipeline.RunCycle(ctx, npcCtx, now2)
	assert.Equal(t, CycleFast, result2.CycleType)
	assert.Equal(t, "behavior_tree", result2.Source)
}

func TestBuildContext(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	tree := NewBehaviorTree("test", nil)
	pipeline := NewNPCDecisionPipeline(cfg, tree, nil, nil)

	npcCtx := NewNPCContext("test_npc")
	npcCtx.Health = 80
	npcCtx.Qi = 50
	npcCtx.MaxQi = 100
	npcCtx.IsInCombat = true

	contextInfo := pipeline.buildContext(npcCtx)
	assert.Contains(t, contextInfo, "test_npc")
	assert.Contains(t, contextInfo, "Health=80")
	assert.Contains(t, contextInfo, "Qi=50")
	assert.Contains(t, contextInfo, "Combat=true")
}

func TestShouldTakeAction(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	tree := NewBehaviorTree("test", nil)
	pipeline := NewNPCDecisionPipeline(cfg, tree, nil, nil)

	// Nil result
	assert.False(t, pipeline.ShouldTakeAction(nil))

	// Nil decision
	assert.False(t, pipeline.ShouldTakeAction(&DecisionCycleResult{}))

	// Low confidence
	result := &DecisionCycleResult{
		Decision: &DecisionResult{Confidence: 0.2},
	}
	assert.False(t, pipeline.ShouldTakeAction(result))

	// High confidence
	result.Decision.Confidence = 0.5
	assert.True(t, pipeline.ShouldTakeAction(result))
}

func TestExecuteAction_Cultivate(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	tree := NewBehaviorTree("test", nil)
	pipeline := NewNPCDecisionPipeline(cfg, tree, nil, nil)

	ctx := context.Background()
	npcCtx := NewNPCContext("test_npc")
	decision := &DecisionResult{Action: "cultivate"}

	err := pipeline.ExecuteAction(ctx, npcCtx, decision)
	assert.NoError(t, err)
	assert.Equal(t, "cultivating", npcCtx.CurrentAction)
	assert.Contains(t, npcCtx.EventLog, "开始修炼")
}

func TestExecuteAction_Rest(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	tree := NewBehaviorTree("test", nil)
	pipeline := NewNPCDecisionPipeline(cfg, tree, nil, nil)

	ctx := context.Background()
	npcCtx := NewNPCContext("test_npc")
	npcCtx.Health = 50
	decision := &DecisionResult{Action: "rest"}

	err := pipeline.ExecuteAction(ctx, npcCtx, decision)
	assert.NoError(t, err)
	assert.Equal(t, "resting", npcCtx.CurrentAction)
	assert.Equal(t, 70, npcCtx.Health) // 50 + 20
}

func TestExecuteAction_Meditate(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	tree := NewBehaviorTree("test", nil)
	pipeline := NewNPCDecisionPipeline(cfg, tree, nil, nil)

	ctx := context.Background()
	npcCtx := NewNPCContext("test_npc")
	npcCtx.Qi = 50
	npcCtx.MaxQi = 100
	decision := &DecisionResult{Action: "meditate"}

	err := pipeline.ExecuteAction(ctx, npcCtx, decision)
	assert.NoError(t, err)
	assert.Equal(t, "meditating", npcCtx.CurrentAction)
	assert.Equal(t, 80.0, npcCtx.Qi) // 50 + 30
}

func TestExecuteAction_Flee(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	tree := NewBehaviorTree("test", nil)
	pipeline := NewNPCDecisionPipeline(cfg, tree, nil, nil)

	ctx := context.Background()
	npcCtx := NewNPCContext("test_npc")
	decision := &DecisionResult{Action: "flee"}

	err := pipeline.ExecuteAction(ctx, npcCtx, decision)
	assert.NoError(t, err)
	assert.Equal(t, "fleeing", npcCtx.CurrentAction)
	assert.Contains(t, npcCtx.EventLog, "逃跑")
}

func TestExecuteAction_Unknown(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	tree := NewBehaviorTree("test", nil)
	pipeline := NewNPCDecisionPipeline(cfg, tree, nil, nil)

	ctx := context.Background()
	npcCtx := NewNPCContext("test_npc")
	decision := &DecisionResult{Action: "unknown_action"}

	err := pipeline.ExecuteAction(ctx, npcCtx, decision)
	assert.NoError(t, err) // Unknown actions don't error
	assert.Contains(t, npcCtx.EventLog[0], "未知动作")
}

func TestExecuteAction_NilDecision(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	tree := NewBehaviorTree("test", nil)
	pipeline := NewNPCDecisionPipeline(cfg, tree, nil, nil)

	ctx := context.Background()
	npcCtx := NewNPCContext("test_npc")

	err := pipeline.ExecuteAction(ctx, npcCtx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "决策为空")
}

func TestGetNextCycleTime(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	tree := NewBehaviorTree("test", nil)
	pipeline := NewNPCDecisionPipeline(cfg, tree, nil, nil)

	now := time.Now()
	pipeline.lastFastCycle = now
	pipeline.lastSlowCycle = now

	// Next fast cycle in 3 seconds
	next := pipeline.GetNextCycleTime(now)
	assert.Equal(t, 3*time.Second, next)

	// After 1 second
	next = pipeline.GetNextCycleTime(now.Add(1 * time.Second))
	assert.Equal(t, 2*time.Second, next)
}

func TestGetPipelineStats(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	tree := NewBehaviorTree("test", nil)
	pipeline := NewNPCDecisionPipeline(cfg, tree, nil, nil)

	now := time.Now()
	pipeline.lastFastCycle = now.Add(-2 * time.Second)
	pipeline.lastSlowCycle = now.Add(-30 * time.Second)
	pipeline.currentAction = "cultivating"

	stats := pipeline.GetPipelineStats(now)
	assert.NotNil(t, stats)
	assert.Equal(t, "cultivating", stats["current_action"])
	assert.Contains(t, stats, "next_cycle_in")
}

func TestExecuteAction_HealthClamp(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	tree := NewBehaviorTree("test", nil)
	pipeline := NewNPCDecisionPipeline(cfg, tree, nil, nil)

	ctx := context.Background()
	npcCtx := NewNPCContext("test_npc")
	npcCtx.Health = 90
	decision := &DecisionResult{Action: "rest"}

	err := pipeline.ExecuteAction(ctx, npcCtx, decision)
	assert.NoError(t, err)
	assert.Equal(t, 100, npcCtx.Health) // Should clamp to 100
}

func TestExecuteAction_QiClamp(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()
	tree := NewBehaviorTree("test", nil)
	pipeline := NewNPCDecisionPipeline(cfg, tree, nil, nil)

	ctx := context.Background()
	npcCtx := NewNPCContext("test_npc")
	npcCtx.Qi = 80
	npcCtx.MaxQi = 100
	decision := &DecisionResult{Action: "meditate"}

	err := pipeline.ExecuteAction(ctx, npcCtx, decision)
	assert.NoError(t, err)
	assert.Equal(t, 100.0, npcCtx.Qi) // Should clamp to MaxQi
}

func TestDecisionPipeline_Integration(t *testing.T) {
	cfg := DefaultDecisionCycleConfig()

	// Create behavior tree
	tree := NewBehaviorTree("npc_behavior", &SelectorNode{
		BaseNode: BaseNode{Name: "root"},
		Children: []BehaviorTreeNode{
			&SequenceNode{
				BaseNode: BaseNode{Name: "combat_branch"},
				Children: []BehaviorTreeNode{
					Condition("is_in_combat", func(ctx *NPCContext) bool { return ctx.IsInCombat }),
					Action("fight", func(ctx *NPCContext) {
						ctx.Log("fight")
					}),
				},
			},
			Action("cultivate", func(ctx *NPCContext) {
				ctx.Log("cultivate")
			}),
		},
	})

	llm := NewLLMClient("test-key", "https://api.test.com")
	templates := NewTemplateLibrary()
	for _, tmpl := range DefaultBehaviorTemplates() {
		templates.AddTemplate(tmpl)
	}

	pipeline := NewNPCDecisionPipeline(cfg, tree, llm, templates)

	ctx := context.Background()
	npcCtx := NewNPCContext("test_npc")

	// Run slow cycle first
	now := time.Now()
	result1, _ := pipeline.RunCycle(ctx, npcCtx, now)
	assert.NotNil(t, result1)
	assert.True(t, pipeline.ShouldTakeAction(result1))

	// Execute the action
	if result1.Decision != nil {
		pipeline.ExecuteAction(ctx, npcCtx, result1.Decision)
	}

	// Fast cycle after delay
	now2 := now.Add(4 * time.Second)
	result2, _ := pipeline.RunCycle(ctx, npcCtx, now2)
	assert.NotNil(t, result2)
	assert.Equal(t, CycleFast, result2.CycleType)

	// Execute based on context
	if npcCtx.IsInCombat {
		pipeline.ExecuteAction(ctx, npcCtx, &DecisionResult{Action: "combat"})
		assert.Equal(t, "combat", npcCtx.CurrentAction)
	} else {
		pipeline.ExecuteAction(ctx, npcCtx, &DecisionResult{Action: "cultivate"})
		assert.Equal(t, "cultivating", npcCtx.CurrentAction)
	}
}

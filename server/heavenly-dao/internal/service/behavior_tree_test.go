package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNPCContext(t *testing.T) {
	ctx := NewNPCContext("test_npc")
	assert.Equal(t, "test_npc", ctx.EntityID)
	assert.Equal(t, 100, ctx.Health)
	assert.Equal(t, 100.0, ctx.Qi)
	assert.Equal(t, 100.0, ctx.MaxQi)
}

func TestContextHasItem(t *testing.T) {
	ctx := NewNPCContext("test_npc")
	ctx.Inventory["灵石"] = 50

	assert.True(t, ctx.HasItem("灵石", 50))
	assert.True(t, ctx.HasItem("灵石", 10))
	assert.False(t, ctx.HasItem("灵石", 51))
	assert.False(t, ctx.HasItem("不存在的物品", 1))
}

func TestContextIsLowOnQi(t *testing.T) {
	ctx := NewNPCContext("test_npc")
	ctx.MaxQi = 100

	ctx.Qi = 10
	assert.True(t, ctx.IsLowOnQi())

	ctx.Qi = 29
	assert.True(t, ctx.IsLowOnQi())

	ctx.Qi = 30
	assert.False(t, ctx.IsLowOnQi())

	ctx.Qi = 50
	assert.False(t, ctx.IsLowOnQi())
}

func TestContextIsInjured(t *testing.T) {
	ctx := NewNPCContext("test_npc")

	ctx.Health = 49
	assert.True(t, ctx.IsInjured())

	ctx.Health = 50
	assert.False(t, ctx.IsInjured())

	ctx.Health = 80
	assert.False(t, ctx.IsInjured())
}

func TestContextIsHealthy(t *testing.T) {
	ctx := NewNPCContext("test_npc")

	ctx.Health = 81
	assert.True(t, ctx.IsHealthy())

	ctx.Health = 80
	assert.False(t, ctx.IsHealthy())

	ctx.Health = 50
	assert.False(t, ctx.IsHealthy())
}

func TestContextLog(t *testing.T) {
	ctx := NewNPCContext("test_npc")
	ctx.Log("开始修炼")
	ctx.Log("突破成功")

	assert.Len(t, ctx.EventLog, 2)
	assert.Equal(t, "开始修炼", ctx.EventLog[0])
	assert.Equal(t, "突破成功", ctx.EventLog[1])
}

func TestLeafNode_Action(t *testing.T) {
	actionCalled := false
	node := &LeafNode{
		BaseNode: BaseNode{Name: "test_action"},
		Action: func(ctx *NPCContext) NodeStatus {
			actionCalled = true
			return StatusSuccess
		},
	}

	ctx := NewNPCContext("test_npc")
	status := node.Evaluate(ctx)
	assert.Equal(t, StatusSuccess, status)
	assert.True(t, actionCalled)
}

func TestLeafNode_ConditionTrue(t *testing.T) {
	node := &LeafNode{
		BaseNode: BaseNode{Name: "is_healthy"},
		Action: func(ctx *NPCContext) NodeStatus {
			if ctx.IsHealthy() {
				return StatusSuccess
			}
			return StatusFailure
		},
	}

	ctx := NewNPCContext("test_npc")
	ctx.Health = 90
	status := node.Evaluate(ctx)
	assert.Equal(t, StatusSuccess, status)
}

func TestLeafNode_ConditionFalse(t *testing.T) {
	node := &LeafNode{
		BaseNode: BaseNode{Name: "is_healthy"},
		Action: func(ctx *NPCContext) NodeStatus {
			if ctx.IsHealthy() {
				return StatusSuccess
			}
			return StatusFailure
		},
	}

	ctx := NewNPCContext("test_npc")
	ctx.Health = 30
	status := node.Evaluate(ctx)
	assert.Equal(t, StatusFailure, status)
}

func TestLeafNode_NilAction(t *testing.T) {
	node := &LeafNode{
		BaseNode: BaseNode{Name: "nil_action"},
	}

	ctx := NewNPCContext("test_npc")
	status := node.Evaluate(ctx)
	assert.Equal(t, StatusFailure, status)
}

func TestSequenceNode_AllSuccess(t *testing.T) {
	seq := &SequenceNode{
		BaseNode: BaseNode{Name: "test_sequence"},
		Children: []BehaviorTreeNode{
			&LeafNode{BaseNode: BaseNode{Name: "step1"}, Action: func(*NPCContext) NodeStatus { return StatusSuccess }},
			&LeafNode{BaseNode: BaseNode{Name: "step2"}, Action: func(*NPCContext) NodeStatus { return StatusSuccess }},
			&LeafNode{BaseNode: BaseNode{Name: "step3"}, Action: func(*NPCContext) NodeStatus { return StatusSuccess }},
		},
	}

	ctx := NewNPCContext("test_npc")
	status := seq.Evaluate(ctx)
	assert.Equal(t, StatusSuccess, status)
}

func TestSequenceNode_FailureStops(t *testing.T) {
	executed := []string{}

	seq := &SequenceNode{
		BaseNode: BaseNode{Name: "test_sequence"},
		Children: []BehaviorTreeNode{
			&LeafNode{
				BaseNode: BaseNode{Name: "step1"},
				Action: func(*NPCContext) NodeStatus {
					executed = append(executed, "step1")
					return StatusSuccess
				},
			},
			&LeafNode{
				BaseNode: BaseNode{Name: "step2"},
				Action: func(*NPCContext) NodeStatus {
					executed = append(executed, "step2")
					return StatusFailure
				},
			},
			&LeafNode{
				BaseNode: BaseNode{Name: "step3"},
				Action: func(*NPCContext) NodeStatus {
					executed = append(executed, "step3")
					return StatusSuccess
				},
			},
		},
	}

	ctx := NewNPCContext("test_npc")
	status := seq.Evaluate(ctx)
	assert.Equal(t, StatusFailure, status)
	assert.Equal(t, []string{"step1", "step2"}, executed)
}

func TestSequenceNode_RunningStops(t *testing.T) {
	seq := &SequenceNode{
		BaseNode: BaseNode{Name: "test_sequence"},
		Children: []BehaviorTreeNode{
			&LeafNode{BaseNode: BaseNode{Name: "step1"}, Action: func(*NPCContext) NodeStatus { return StatusSuccess }},
			&LeafNode{BaseNode: BaseNode{Name: "step2"}, Action: func(*NPCContext) NodeStatus { return StatusRunning }},
			&LeafNode{BaseNode: BaseNode{Name: "step3"}, Action: func(*NPCContext) NodeStatus { return StatusSuccess }},
		},
	}

	ctx := NewNPCContext("test_npc")
	status := seq.Evaluate(ctx)
	assert.Equal(t, StatusRunning, status)
}

func TestSelectorNode_FirstSuccess(t *testing.T) {
	executed := []string{}

	sel := &SelectorNode{
		BaseNode: BaseNode{Name: "test_selector"},
		Children: []BehaviorTreeNode{
			&LeafNode{
				BaseNode: BaseNode{Name: "option1"},
				Action: func(*NPCContext) NodeStatus {
					executed = append(executed, "option1")
					return StatusSuccess
				},
			},
			&LeafNode{
				BaseNode: BaseNode{Name: "option2"},
				Action: func(*NPCContext) NodeStatus {
					executed = append(executed, "option2")
					return StatusSuccess
				},
			},
		},
	}

	ctx := NewNPCContext("test_npc")
	status := sel.Evaluate(ctx)
	assert.Equal(t, StatusSuccess, status)
	assert.Equal(t, []string{"option1"}, executed)
}

func TestSelectorNode_SkipFailures(t *testing.T) {
	executed := []string{}

	sel := &SelectorNode{
		BaseNode: BaseNode{Name: "test_selector"},
		Children: []BehaviorTreeNode{
			&LeafNode{
				BaseNode: BaseNode{Name: "option1"},
				Action: func(*NPCContext) NodeStatus {
					executed = append(executed, "option1")
					return StatusFailure
				},
			},
			&LeafNode{
				BaseNode: BaseNode{Name: "option2"},
				Action: func(*NPCContext) NodeStatus {
					executed = append(executed, "option2")
					return StatusFailure
				},
			},
			&LeafNode{
				BaseNode: BaseNode{Name: "option3"},
				Action: func(*NPCContext) NodeStatus {
					executed = append(executed, "option3")
					return StatusSuccess
				},
			},
		},
	}

	ctx := NewNPCContext("test_npc")
	status := sel.Evaluate(ctx)
	assert.Equal(t, StatusSuccess, status)
	assert.Equal(t, []string{"option1", "option2", "option3"}, executed)
}

func TestSelectorNode_AllFailure(t *testing.T) {
	sel := &SelectorNode{
		BaseNode: BaseNode{Name: "test_selector"},
		Children: []BehaviorTreeNode{
			&LeafNode{BaseNode: BaseNode{Name: "option1"}, Action: func(*NPCContext) NodeStatus { return StatusFailure }},
			&LeafNode{BaseNode: BaseNode{Name: "option2"}, Action: func(*NPCContext) NodeStatus { return StatusFailure }},
		},
	}

	ctx := NewNPCContext("test_npc")
	status := sel.Evaluate(ctx)
	assert.Equal(t, StatusFailure, status)
}

func TestSelectorNode_RunningStops(t *testing.T) {
	sel := &SelectorNode{
		BaseNode: BaseNode{Name: "test_selector"},
		Children: []BehaviorTreeNode{
			&LeafNode{BaseNode: BaseNode{Name: "option1"}, Action: func(*NPCContext) NodeStatus { return StatusFailure }},
			&LeafNode{BaseNode: BaseNode{Name: "option2"}, Action: func(*NPCContext) NodeStatus { return StatusRunning }},
			&LeafNode{BaseNode: BaseNode{Name: "option3"}, Action: func(*NPCContext) NodeStatus { return StatusSuccess }},
		},
	}

	ctx := NewNPCContext("test_npc")
	status := sel.Evaluate(ctx)
	assert.Equal(t, StatusRunning, status)
}

func TestDecoratorNode_Inverter(t *testing.T) {
	inv := Inverter(&LeafNode{
		BaseNode: BaseNode{Name: "always_success"},
		Action:   func(*NPCContext) NodeStatus { return StatusSuccess },
	})

	ctx := NewNPCContext("test_npc")
	status := inv.Evaluate(ctx)
	assert.Equal(t, StatusFailure, status)

	inv2 := Inverter(&LeafNode{
		BaseNode: BaseNode{Name: "always_failure"},
		Action:   func(*NPCContext) NodeStatus { return StatusFailure },
	})
	status = inv2.Evaluate(ctx)
	assert.Equal(t, StatusSuccess, status)

	// Running stays running
	inv3 := Inverter(&LeafNode{
		BaseNode: BaseNode{Name: "always_running"},
		Action:   func(*NPCContext) NodeStatus { return StatusRunning },
	})
	status = inv3.Evaluate(ctx)
	assert.Equal(t, StatusRunning, status)
}

func TestDecoratorNode_NilChild(t *testing.T) {
	dec := &DecoratorNode{
		BaseNode: BaseNode{Name: "nil_child"},
	}

	ctx := NewNPCContext("test_npc")
	status := dec.Evaluate(ctx)
	assert.Equal(t, StatusFailure, status)
}

func TestBehaviorTree_Evaluate(t *testing.T) {
	root := &SequenceNode{
		BaseNode: BaseNode{Name: "root"},
		Children: []BehaviorTreeNode{
			Condition("is_healthy", func(ctx *NPCContext) bool { return ctx.IsHealthy() }),
			Action("cultivate", func(ctx *NPCContext) {
				ctx.Qi += 10
			}),
		},
	}

	tree := NewBehaviorTree("cultivation_tree", root)

	ctx := NewNPCContext("test_npc")
	ctx.Health = 90
	status := tree.Evaluate(ctx)
	assert.Equal(t, StatusSuccess, status)
	assert.Equal(t, 110.0, ctx.Qi)
}

func TestBehaviorTree_NilRoot(t *testing.T) {
	tree := NewBehaviorTree("empty_tree", nil)

	ctx := NewNPCContext("test_npc")
	status := tree.Evaluate(ctx)
	assert.Equal(t, StatusFailure, status)
}

func TestBehaviorTree_GetExecutionTrace(t *testing.T) {
	root := &SequenceNode{
		BaseNode: BaseNode{Name: "root"},
		Children: []BehaviorTreeNode{
			Condition("is_healthy", func(ctx *NPCContext) bool { return ctx.IsHealthy() }),
			&SelectorNode{
				BaseNode: BaseNode{Name: "action_selector"},
				Children: []BehaviorTreeNode{
					Action("cultivate", func(ctx *NPCContext) {}),
					Action("rest", func(ctx *NPCContext) {}),
				},
			},
		},
	}

	tree := NewBehaviorTree("test_tree", root)
	trace := tree.GetExecutionTrace()

	assert.Contains(t, trace, "[sequence] root")
	assert.Contains(t, trace, "[leaf] is_healthy")
	assert.Contains(t, trace, "[selector] action_selector")
	assert.Contains(t, trace, "[leaf] cultivate")
	assert.Contains(t, trace, "[leaf] rest")
}

func TestNodeGetTypes(t *testing.T) {
	seq := &SequenceNode{}
	assert.Equal(t, "sequence", seq.GetType())

	sel := &SelectorNode{}
	assert.Equal(t, "selector", sel.GetType())

	dec := &DecoratorNode{}
	assert.Equal(t, "decorator", dec.GetType())

	leaf := &LeafNode{}
	assert.Equal(t, "leaf", leaf.GetType())
}

func TestConditionHelper(t *testing.T) {
	cond := Condition("check_health", func(ctx *NPCContext) bool {
		return ctx.Health > 50
	})

	ctx := NewNPCContext("test_npc")
	ctx.Health = 80
	assert.Equal(t, StatusSuccess, cond.Evaluate(ctx))

	ctx.Health = 30
	assert.Equal(t, StatusFailure, cond.Evaluate(ctx))
}

func TestActionHelper(t *testing.T) {
	actionCalled := false
	act := Action("heal", func(ctx *NPCContext) {
		ctx.Health = 100
		actionCalled = true
	})

	ctx := NewNPCContext("test_npc")
	ctx.Health = 30
	status := act.Evaluate(ctx)
	assert.Equal(t, StatusSuccess, status)
	assert.True(t, actionCalled)
	assert.Equal(t, 100, ctx.Health)
}

func TestComplexBehaviorTree_Deterministic(t *testing.T) {
	// Create a deterministic behavior tree
	root := &SequenceNode{
		BaseNode: BaseNode{Name: "npc_daily_routine"},
		Children: []BehaviorTreeNode{
			&SelectorNode{
				BaseNode: BaseNode{Name: "survival_check"},
				Children: []BehaviorTreeNode{
					&SequenceNode{
						BaseNode: BaseNode{Name: "flee_if_danger"},
						Children: []BehaviorTreeNode{
							Condition("is_in_danger", func(ctx *NPCContext) bool { return ctx.IsInDanger }),
							Action("flee", func(ctx *NPCContext) {
								ctx.CurrentAction = "fleeing"
							}),
						},
					},
					&SequenceNode{
						BaseNode: BaseNode{Name: "heal_if_injured"},
						Children: []BehaviorTreeNode{
							Condition("is_injured", func(ctx *NPCContext) bool { return ctx.IsInjured() }),
							Action("heal", func(ctx *NPCContext) {
								ctx.Health = 100
							}),
						},
					},
			Action("continue", func(ctx *NPCContext) {
				// always succeeds
			}),
				},
			},
			Action("cultivate", func(ctx *NPCContext) {
				ctx.Qi += 5
			}),
		},
	}

	tree := NewBehaviorTree("npc_daily_routine", root)

	// Test with danger
	ctx1 := NewNPCContext("npc1")
	ctx1.IsInDanger = true
	ctx1.Health = 80
	tree.Evaluate(ctx1)
	assert.Equal(t, "fleeing", ctx1.CurrentAction)

	// Test with injury
	ctx2 := NewNPCContext("npc2")
	ctx2.Health = 30
	tree.Evaluate(ctx2)
	assert.Equal(t, 100, ctx2.Health)

	// Test normal
	ctx3 := NewNPCContext("npc3")
	ctx3.IsInDanger = false
	ctx3.Health = 90
	tree.Evaluate(ctx3)
	assert.Equal(t, 105.0, ctx3.Qi)
}

func TestComplexBehaviorTree_DeterministicSameInput(t *testing.T) {
	root := &SelectorNode{
		BaseNode: BaseNode{Name: "combat_decision"},
		Children: []BehaviorTreeNode{
			&SequenceNode{
				BaseNode: BaseNode{Name: "attack_if_enemy_weak"},
				Children: []BehaviorTreeNode{
					Condition("has_target", func(ctx *NPCContext) bool { return ctx.HasTarget }),
					Condition("is_strong", func(ctx *NPCContext) bool { return ctx.Health > 70 }),
					Action("attack", func(ctx *NPCContext) {
						ctx.Log("attack")
					}),
				},
			},
			Action("defend", func(ctx *NPCContext) {
				ctx.Log("defend")
			}),
		},
	}

	tree := NewBehaviorTree("combat", root)

	// Same input should produce same output
	for i := 0; i < 10; i++ {
	ctx := NewNPCContext("npc")
	ctx.HasTarget = true
	ctx.Health = 80
	tree.Evaluate(ctx)
	assert.Equal(t, []string{"attack"}, ctx.EventLog)
}

// Different input
for i := 0; i < 10; i++ {
	ctx := NewNPCContext("npc")
	ctx.HasTarget = false
	ctx.Health = 80
	tree.Evaluate(ctx)
	assert.Equal(t, []string{"defend"}, ctx.EventLog)
	}
}

func TestContextMemory(t *testing.T) {
	ctx := NewNPCContext("test_npc")
	ctx.Memory["last_action"] = "cultivate"
	ctx.Memory["target_id"] = "enemy_1"

	assert.Equal(t, "cultivate", ctx.Memory["last_action"])
	assert.Equal(t, "enemy_1", ctx.Memory["target_id"])
}

func TestNodeGetName(t *testing.T) {
	node := &LeafNode{BaseNode: BaseNode{Name: "my_action"}}
	assert.Equal(t, "my_action", node.GetName())
}

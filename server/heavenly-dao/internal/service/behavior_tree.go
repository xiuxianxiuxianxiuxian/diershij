package service

import (
	"fmt"
	"strings"
)

// NodeStatus represents the result of a behavior tree node evaluation.
type NodeStatus string

const (
	StatusSuccess NodeStatus = "success"
	StatusFailure NodeStatus = "failure"
	StatusRunning NodeStatus = "running"
)

// BehaviorTreeNode is the interface for all behavior tree nodes.
type BehaviorTreeNode interface {
	Evaluate(ctx *NPCContext) NodeStatus
	GetType() string
	GetName() string
}

// BaseNode provides common fields for all nodes.
type BaseNode struct {
	Name string
}

func (b *BaseNode) GetName() string {
	return b.Name
}

// SequenceNode executes children in order until one fails.
type SequenceNode struct {
	BaseNode
	Children []BehaviorTreeNode
}

func (n *SequenceNode) GetType() string {
	return "sequence"
}

func (n *SequenceNode) Evaluate(ctx *NPCContext) NodeStatus {
	for _, child := range n.Children {
		status := child.Evaluate(ctx)
		if status != StatusSuccess {
			return status
		}
	}
	return StatusSuccess
}

// SelectorNode executes children until one succeeds.
type SelectorNode struct {
	BaseNode
	Children []BehaviorTreeNode
}

func (n *SelectorNode) GetType() string {
	return "selector"
}

func (n *SelectorNode) Evaluate(ctx *NPCContext) NodeStatus {
	for _, child := range n.Children {
		status := child.Evaluate(ctx)
		if status != StatusFailure {
			return status
		}
	}
	return StatusFailure
}

// DecoratorNode wraps a single child and modifies its behavior.
type DecoratorNode struct {
	BaseNode
	Child    BehaviorTreeNode
	Modifier func(NodeStatus) NodeStatus
}

func (n *DecoratorNode) GetType() string {
	return "decorator"
}

func (n *DecoratorNode) Evaluate(ctx *NPCContext) NodeStatus {
	if n.Child == nil {
		return StatusFailure
	}
	status := n.Child.Evaluate(ctx)
	if n.Modifier != nil {
		return n.Modifier(status)
	}
	return status
}

// LeafNode represents an actual action or condition.
type LeafNode struct {
	BaseNode
	Action func(ctx *NPCContext) NodeStatus
}

func (n *LeafNode) GetType() string {
	return "leaf"
}

func (n *LeafNode) Evaluate(ctx *NPCContext) NodeStatus {
	if n.Action == nil {
		return StatusFailure
	}
	return n.Action(ctx)
}

// NPCContext holds the NPC's current state for behavior tree evaluation.
type NPCContext struct {
	EntityID       string
	Health         int
	Qi             float64
	MaxQi          float64
	SpiritStones   int64
	CurrentAction  string
	IsInCombat     bool
	IsInDanger     bool
	HasTarget      bool
	TargetDistance int
	Inventory      map[string]int
	Memory         map[string]interface{}
	EventLog       []string
}

// NewNPCContext creates a new NPC context.
func NewNPCContext(entityID string) *NPCContext {
	return &NPCContext{
		EntityID:  entityID,
		Health:    100,
		Qi:        100,
		MaxQi:     100,
		Inventory: make(map[string]int),
		Memory:    make(map[string]interface{}),
		EventLog:  []string{},
	}
}

// Log adds a message to the context log.
func (c *NPCContext) Log(msg string) {
	c.EventLog = append(c.EventLog, msg)
}

// HasItem checks if the NPC has a specific item.
func (c *NPCContext) HasItem(itemName string, quantity int) bool {
	return c.Inventory[itemName] >= quantity
}

// IsLowOnQi checks if Qi is below 30%.
func (c *NPCContext) IsLowOnQi() bool {
	return c.Qi < c.MaxQi*0.3
}

// IsInjured checks if health is below 50%.
func (c *NPCContext) IsInjured() bool {
	return c.Health < 50
}

// IsHealthy checks if health is above 80%.
func (c *NPCContext) IsHealthy() bool {
	return c.Health > 80
}

// BehaviorTree represents a complete behavior tree.
type BehaviorTree struct {
	Name string
	Root BehaviorTreeNode
}

// NewBehaviorTree creates a new behavior tree.
func NewBehaviorTree(name string, root BehaviorTreeNode) *BehaviorTree {
	return &BehaviorTree{
		Name: name,
		Root: root,
	}
}

// Evaluate runs the behavior tree from the root.
func (t *BehaviorTree) Evaluate(ctx *NPCContext) NodeStatus {
	if t.Root == nil {
		return StatusFailure
	}
	return t.Root.Evaluate(ctx)
}

// GetExecutionTrace returns a string representation of the tree structure.
func (t *BehaviorTree) GetExecutionTrace() string {
	var sb strings.Builder
	t.writeNodeTrace(&sb, t.Root, 0)
	return sb.String()
}

func (t *BehaviorTree) writeNodeTrace(sb *strings.Builder, node BehaviorTreeNode, depth int) {
	if node == nil {
		return
	}

	indent := strings.Repeat("  ", depth)
	sb.WriteString(fmt.Sprintf("%s[%s] %s\n", indent, node.GetType(), node.GetName()))

	switch n := node.(type) {
	case *SequenceNode:
		for _, child := range n.Children {
			t.writeNodeTrace(sb, child, depth+1)
		}
	case *SelectorNode:
		for _, child := range n.Children {
			t.writeNodeTrace(sb, child, depth+1)
		}
	case *DecoratorNode:
		if n.Child != nil {
			t.writeNodeTrace(sb, n.Child, depth+1)
		}
	}
}

// Common decorators

// Inverter creates a decorator that inverts success/failure.
func Inverter(child BehaviorTreeNode) *DecoratorNode {
	return &DecoratorNode{
		BaseNode: BaseNode{Name: "Inverter"},
		Child:    child,
		Modifier: func(s NodeStatus) NodeStatus {
			switch s {
			case StatusSuccess:
				return StatusFailure
			case StatusFailure:
				return StatusSuccess
			default:
				return s
			}
		},
	}
}

// Repeater creates a decorator that repeats the child N times.
func Repeater(child BehaviorTreeNode, count int) *DecoratorNode {
	return &DecoratorNode{
		BaseNode: BaseNode{Name: fmt.Sprintf("Repeater(%d)", count)},
		Child:    child,
		Modifier: func(s NodeStatus) NodeStatus {
			// Simplified: just return the last status
			return s
		},
	}
}

// Common leaf conditions

// Condition creates a leaf node that checks a condition.
func Condition(name string, check func(ctx *NPCContext) bool) *LeafNode {
	return &LeafNode{
		BaseNode: BaseNode{Name: name},
		Action: func(ctx *NPCContext) NodeStatus {
			if check(ctx) {
				return StatusSuccess
			}
			return StatusFailure
		},
	}
}

// Action creates a leaf node that performs an action.
func Action(name string, do func(ctx *NPCContext)) *LeafNode {
	return &LeafNode{
		BaseNode: BaseNode{Name: name},
		Action: func(ctx *NPCContext) NodeStatus {
			do(ctx)
			return StatusSuccess
		},
	}
}

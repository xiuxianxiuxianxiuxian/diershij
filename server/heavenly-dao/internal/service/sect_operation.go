package service

import (
	"fmt"
	"time"
)

// SectOperation handles sect-related operations (form, join, leave).
type SectOperation struct {
	cooldownMap    map[string]time.Time
	cooldownPeriod time.Duration
	sects          map[string]*Sect
	memberContributions map[string]map[string]int // sectID -> memberID -> contribution
}

// Sect represents a cultivation sect.
type Sect struct {
	ID           string
	Name         string
	FounderID    string
	FounderName  string
	MemberCount  int
	MaxMembers   int
	Level        int
	Prestige     int
	Territories  []string
	CreatedAt    time.Time
}

// NewSectOperation creates a new SectOperation.
func NewSectOperation(cooldownPeriod time.Duration) *SectOperation {
	return &SectOperation{
		cooldownMap:       make(map[string]time.Time),
		cooldownPeriod:    cooldownPeriod,
		sects:             make(map[string]*Sect),
		memberContributions: make(map[string]map[string]int),
	}
}

// FormSectInput holds inputs for forming a new sect.
type FormSectInput struct {
	FounderID   string
	FounderName string
	SectName    string
	MaxMembers  int
}

// SectResult holds the outcome of a sect operation.
type SectResult struct {
	Success bool
	Sect    *Sect
	Message string
}

// ExecuteFormSect creates a new sect.
func (op *SectOperation) ExecuteFormSect(input FormSectInput, now time.Time) (*SectResult, error) {
	// Check cooldown
	if lastTime, ok := op.cooldownMap[input.FounderID]; ok {
		if now.Sub(lastTime) < op.cooldownPeriod {
			return nil, fmt.Errorf("sect cooldown: %v remaining", op.cooldownPeriod-now.Sub(lastTime))
		}
	}

	// Check sect name uniqueness
	for _, sect := range op.sects {
		if sect.Name == input.SectName {
			return nil, fmt.Errorf("宗门 '%s' 已存在", input.SectName)
		}
	}

	// Validate max members
	maxMembers := input.MaxMembers
	if maxMembers < 2 {
		maxMembers = 2
	}
	if maxMembers > 10000 {
		maxMembers = 10000
	}

	sect := &Sect{
		ID:          input.SectName,
		Name:        input.SectName,
		FounderID:   input.FounderID,
		FounderName: input.FounderName,
		MemberCount: 1,
		MaxMembers:  maxMembers,
		Level:       1,
		Prestige:    100,
		CreatedAt:   now,
	}

	op.sects[sect.ID] = sect
	op.memberContributions[sect.ID] = make(map[string]int)
	op.memberContributions[sect.ID][input.FounderID] = 1000 // founder gets initial contribution

	op.cooldownMap[input.FounderID] = now

	return &SectResult{
		Success: true,
		Sect:    sect,
		Message: fmt.Sprintf("成功创建宗门 '%s'", sect.Name),
	}, nil
}

// JoinSectInput holds inputs for joining a sect.
type JoinSectInput struct {
	EntityID string
	SectID   string
}

// ExecuteJoinSect adds a member to a sect.
func (op *SectOperation) ExecuteJoinSect(input JoinSectInput, now time.Time) (*SectResult, error) {
	sect, exists := op.sects[input.SectID]
	if !exists {
		return nil, fmt.Errorf("宗门 '%s' 不存在", input.SectID)
	}

	if sect.MemberCount >= sect.MaxMembers {
		return nil, fmt.Errorf("宗门已满员 (%d/%d)", sect.MemberCount, sect.MaxMembers)
	}

	sect.MemberCount++

	if op.memberContributions[sect.ID] == nil {
		op.memberContributions[sect.ID] = make(map[string]int)
	}
	op.memberContributions[sect.ID][input.EntityID] = 100

	op.cooldownMap[input.EntityID] = now

	return &SectResult{
		Success: true,
		Sect:    sect,
		Message: fmt.Sprintf("成功加入宗门 '%s'", sect.Name),
	}, nil
}

// LeaveSectInput holds inputs for leaving a sect.
type LeaveSectInput struct {
	EntityID string
	SectID   string
}

// ExecuteLeaveSect removes a member from a sect.
func (op *SectOperation) ExecuteLeaveSect(input LeaveSectInput, now time.Time) (*SectResult, error) {
	sect, exists := op.sects[input.SectID]
	if !exists {
		return nil, fmt.Errorf("宗门 '%s' 不存在", input.SectID)
	}

	// Founder cannot leave
	if sect.FounderID == input.EntityID {
		return nil, fmt.Errorf("宗主不能退出宗门")
	}

	if sect.MemberCount <= 1 {
		return nil, fmt.Errorf("宗门只剩最后一人，无法退出")
	}

	sect.MemberCount--
	delete(op.memberContributions[sect.ID], input.EntityID)

	op.cooldownMap[input.EntityID] = now

	return &SectResult{
		Success: true,
		Sect:    sect,
		Message: fmt.Sprintf("已退出宗门 '%s'", sect.Name),
	}, nil
}

// AddContribution adds contribution points to a member.
func (op *SectOperation) AddContribution(sectID, entityID string, amount int) error {
	if contribs, exists := op.memberContributions[sectID]; exists {
		contribs[entityID] += amount
		return nil
	}
	return fmt.Errorf("宗门 '%s' 不存在", sectID)
}

// GetContribution returns a member's contribution points.
func (op *SectOperation) GetContribution(sectID, entityID string) (int, error) {
	if contribs, exists := op.memberContributions[sectID]; exists {
		return contribs[entityID], nil
	}
	return 0, fmt.Errorf("宗门 '%s' 不存在", sectID)
}

// GetSect retrieves a sect by ID.
func (op *SectOperation) GetSect(id string) (*Sect, bool) {
	s, ok := op.sects[id]
	return s, ok
}

// ListSects returns all sects.
func (op *SectOperation) ListSects() []*Sect {
	sects := make([]*Sect, 0, len(op.sects))
	for _, s := range op.sects {
		sects = append(sects, s)
	}
	return sects
}

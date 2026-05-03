package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSectOperation(t *testing.T) {
	op := NewSectOperation(time.Minute)
	assert.NotNil(t, op)
}

func TestExecuteFormSect_Success(t *testing.T) {
	op := NewSectOperation(time.Minute)

	input := FormSectInput{
		FounderID:   "founder_1",
		FounderName: "张三",
		SectName:    "青云宗",
		MaxMembers:  100,
	}

	now := time.Now()
	result, err := op.ExecuteFormSect(input, now)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Sect)
	assert.Equal(t, "青云宗", result.Sect.Name)
	assert.Equal(t, 1, result.Sect.MemberCount)
	assert.Contains(t, result.Message, "成功创建")
}

func TestExecuteFormSect_DuplicateName(t *testing.T) {
	op := NewSectOperation(time.Minute)

	input := FormSectInput{
		FounderID:   "founder_1",
		SectName:    "青云宗",
	}

	now := time.Now()
	_, err := op.ExecuteFormSect(input, now)
	assert.NoError(t, err)

	_, err = op.ExecuteFormSect(input, now.Add(2*time.Minute))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已存在")
}

func TestExecuteJoinSect_Success(t *testing.T) {
	op := NewSectOperation(time.Minute)

	// Create sect first
	formInput := FormSectInput{
		FounderID:   "founder_1",
		SectName:    "青云宗",
		MaxMembers:  10,
	}
	now := time.Now()
	_, _ = op.ExecuteFormSect(formInput, now)

	// Join sect
	joinInput := JoinSectInput{
		EntityID: "member_1",
		SectID:   "青云宗",
	}

	result, err := op.ExecuteJoinSect(joinInput, now.Add(time.Minute))
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 2, result.Sect.MemberCount)
}

func TestExecuteJoinSect_Full(t *testing.T) {
	op := NewSectOperation(time.Minute)

	formInput := FormSectInput{
		FounderID:  "founder_full",
		SectName:   "SmallSect",
		MaxMembers: 2,
	}
	now := time.Now()
	_, _ = op.ExecuteFormSect(formInput, now)

	// First member joins
	joinResult, err := op.ExecuteJoinSect(JoinSectInput{EntityID: "member_full_1", SectID: "SmallSect"}, now.Add(time.Minute))
	assert.NoError(t, err)
	assert.Equal(t, 2, joinResult.Sect.MemberCount)

	// Try to join when full - should fail
	_, err = op.ExecuteJoinSect(JoinSectInput{EntityID: "member_full_2", SectID: "SmallSect"}, now.Add(2*time.Minute))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "满员")
}

func TestExecuteLeaveSect_Success(t *testing.T) {
	op := NewSectOperation(time.Minute)

	formInput := FormSectInput{
		FounderID:  "founder_1",
		SectName:   "青云宗",
		MaxMembers: 100,
	}
	now := time.Now()
	_, _ = op.ExecuteFormSect(formInput, now)

	// Member joins and leaves
	_, _ = op.ExecuteJoinSect(JoinSectInput{EntityID: "member_1", SectID: "青云宗"}, now.Add(time.Minute))

	result, err := op.ExecuteLeaveSect(LeaveSectInput{EntityID: "member_1", SectID: "青云宗"}, now.Add(2*time.Minute))
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, 1, result.Sect.MemberCount)
}

func TestExecuteLeaveSect_Founder(t *testing.T) {
	op := NewSectOperation(time.Minute)

	formInput := FormSectInput{
		FounderID:  "founder_1",
		SectName:   "青云宗",
	}
	now := time.Now()
	_, _ = op.ExecuteFormSect(formInput, now)

	_, err := op.ExecuteLeaveSect(LeaveSectInput{EntityID: "founder_1", SectID: "青云宗"}, now.Add(time.Minute))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "宗主")
}

func TestAddContribution(t *testing.T) {
	op := NewSectOperation(time.Minute)

	formInput := FormSectInput{
		FounderID:  "founder_1",
		SectName:   "青云宗",
	}
	now := time.Now()
	_, _ = op.ExecuteFormSect(formInput, now)

	err := op.AddContribution("青云宗", "founder_1", 500)
	assert.NoError(t, err)

	contrib, err := op.GetContribution("青云宗", "founder_1")
	assert.NoError(t, err)
	assert.Equal(t, 1500, contrib) // 1000 initial + 500
}

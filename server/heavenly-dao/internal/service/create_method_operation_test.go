package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func makeCreateMethodInput() CreateMethodInput {
	return CreateMethodInput{
		EntityID:        "creator_1",
		EntityName:      "张三",
		MethodName:      "烈火剑法",
		MethodCategory:  "秘术",
		ElementType:     "fire",
		RequiredRealm:   "qi_condensation",
		PremiumStones:   15000,
		Comprehension:   80,
		MentalStability: 70,
		DaoHeart:        60,
	}
}

func TestNewCreateMethodOperation(t *testing.T) {
	op := NewCreateMethodOperation(time.Minute)
	assert.NotNil(t, op)
}

func TestExecuteCreateMethod_Success(t *testing.T) {
	op := NewCreateMethodOperation(time.Minute)

	input := makeCreateMethodInput()
	now := time.Now()

	result, err := op.ExecuteCreateMethod(input, now)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Method)
	assert.Equal(t, input.MethodName, result.Method.Name)
	assert.True(t, result.Quality > 0)
	assert.Equal(t, int64(createMethodCost), result.Cost)
	assert.Contains(t, result.Message, "成功创功")
}

func TestExecuteCreateMethod_InsufficientStones(t *testing.T) {
	op := NewCreateMethodOperation(time.Minute)

	input := makeCreateMethodInput()
	input.PremiumStones = 5000 // below 10000
	now := time.Now()

	_, err := op.ExecuteCreateMethod(input, now)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "极品灵石")
}

func TestExecuteCreateMethod_DuplicateName(t *testing.T) {
	op := NewCreateMethodOperation(time.Minute)

	input := makeCreateMethodInput()
	now := time.Now()

	_, err := op.ExecuteCreateMethod(input, now)
	assert.NoError(t, err)

	// Try again with same name
	_, err = op.ExecuteCreateMethod(input, now.Add(2*time.Minute))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "已存在")
}

func TestTrackLearner_Reward(t *testing.T) {
	op := NewCreateMethodOperation(time.Minute)

	input := makeCreateMethodInput()
	now := time.Now()
	result, _ := op.ExecuteCreateMethod(input, now)

	// Add 9 learners to reach 10
	for i := 0; i < 9; i++ {
		_, _ = op.TrackLearner(result.Method.Name, "learner")
	}

	// 10th learner should trigger reward
	reward, err := op.TrackLearner(result.Method.Name, "learner_10")
	assert.NoError(t, err)
	assert.True(t, reward >= 5)
	assert.True(t, reward <= 10)
	assert.Equal(t, 10, result.Method.Popularity)
}

func TestGetMethod(t *testing.T) {
	op := NewCreateMethodOperation(time.Minute)

	input := makeCreateMethodInput()
	now := time.Now()
	_, _ = op.ExecuteCreateMethod(input, now)

	method, exists := op.GetMethod("烈火剑法")
	assert.True(t, exists)
	assert.Equal(t, "烈火剑法", method.Name)

	_, exists = op.GetMethod("nonexistent")
	assert.False(t, exists)
}

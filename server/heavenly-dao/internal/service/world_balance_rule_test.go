package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorldBalanceRule(t *testing.T) {
	rule := NewWorldBalanceRule()
	assert.NotNil(t, rule)
	assert.Equal(t, 0.40, rule.healthyGiniThreshold)
	assert.Equal(t, 0.60, rule.unhealthyGiniThreshold)
	assert.Equal(t, 3, rule.minSectDiversity)
}

func TestCalculateGiniCoefficient_PerfectEquality(t *testing.T) {
	rule := NewWorldBalanceRule()

	values := []float64{100, 100, 100, 100, 100}
	gini := rule.calculateGiniCoefficient(values)
	assert.InDelta(t, 0.0, gini, 0.001)
}

func TestCalculateGiniCoefficient_HighInequality(t *testing.T) {
	rule := NewWorldBalanceRule()

	// One person has almost everything
	values := []float64{1, 1, 1, 1, 1000}
	gini := rule.calculateGiniCoefficient(values)
	assert.True(t, gini > 0.7)
}

func TestCalculateGiniCoefficient_Empty(t *testing.T) {
	rule := NewWorldBalanceRule()

	gini := rule.calculateGiniCoefficient([]float64{})
	assert.Equal(t, 0.0, gini)
}

func TestCalculateGiniCoefficient_AllZero(t *testing.T) {
	rule := NewWorldBalanceRule()

	gini := rule.calculateGiniCoefficient([]float64{0, 0, 0})
	assert.Equal(t, 0.0, gini)
}

func TestCalculateGiniCoefficient_Moderate(t *testing.T) {
	rule := NewWorldBalanceRule()

	values := []float64{10, 20, 30, 40, 50}
	gini := rule.calculateGiniCoefficient(values)
	// Should be between 0 and 1
	assert.True(t, gini > 0 && gini < 0.5)
}

func TestCalculateResourceRatio_Balanced(t *testing.T) {
	rule := NewWorldBalanceRule()

	ratio := rule.calculateResourceRatio(100, 100)
	assert.InDelta(t, 1.0, ratio, 0.001)
}

func TestCalculateResourceRatio_Surplus(t *testing.T) {
	rule := NewWorldBalanceRule()

	ratio := rule.calculateResourceRatio(150, 100)
	assert.InDelta(t, 1.5, ratio, 0.001)
}

func TestCalculateResourceRatio_Shortage(t *testing.T) {
	rule := NewWorldBalanceRule()

	ratio := rule.calculateResourceRatio(50, 100)
	assert.InDelta(t, 0.5, ratio, 0.001)
}

func TestCalculateResourceRatio_NoConsumption(t *testing.T) {
	rule := NewWorldBalanceRule()

	ratio := rule.calculateResourceRatio(100, 0)
	assert.InDelta(t, 1.0, ratio, 0.001)
}

func TestCalculateSectDiversityScore_NoSects(t *testing.T) {
	rule := NewWorldBalanceRule()

	score := rule.calculateSectDiversityScore(0, []int{})
	assert.Equal(t, 0.0, score)
}

func TestCalculateSectDiversityScore_EqualSects(t *testing.T) {
	rule := NewWorldBalanceRule()

	// 5 equal sects
	score := rule.calculateSectDiversityScore(5, []int{100, 100, 100, 100, 100})
	// count_score = min(1.0, 5/10) = 0.5
	// size_score = 1.0 (no variance)
	// total = (0.5 + 1.0) / 2 = 0.75
	assert.InDelta(t, 0.75, score, 0.01)
}

func TestCalculateSectDiversityScore_UnequalSects(t *testing.T) {
	rule := NewWorldBalanceRule()

	// One dominant sect
	score := rule.calculateSectDiversityScore(3, []int{1000, 10, 10})
	// High CV → lower size_score
	assert.True(t, score < 0.6)
}

func TestCalculateStdDev_Zero(t *testing.T) {
	rule := NewWorldBalanceRule()

	stdDev := rule.calculateStdDev([]float64{5, 5, 5, 5})
	assert.InDelta(t, 0.0, stdDev, 0.001)
}

func TestCalculateStdDev_NonZero(t *testing.T) {
	rule := NewWorldBalanceRule()

	stdDev := rule.calculateStdDev([]float64{10, 20, 30, 40, 50})
	// mean = 30, variance = (400+100+0+100+400)/5 = 200, stdDev = sqrt(200) ≈ 14.14
	assert.InDelta(t, 14.14, stdDev, 0.1)
}

func TestEvaluateWorldHealth_Healthy(t *testing.T) {
	rule := NewWorldBalanceRule()

	metrics := WorldMetrics{
		EntitySpiritStones:      []float64{100, 120, 110, 90, 105},
		ResourceSpawnRate:       100,
		ResourceConsumptionRate: 100,
		ActiveSectCount:         5,
		SectSizes:               []int{100, 100, 100, 100, 100},
		KarmaValues:             []float64{50, 55, 45, 52, 48},
	}

	health := rule.EvaluateWorldHealth(metrics)
	assert.True(t, health.GiniCoefficient < 0.1)
	assert.InDelta(t, 1.0, health.ResourceRatio, 0.001)
	assert.True(t, health.OverallScore > 70)
	assert.Equal(t, "healthy", health.Status)
}

func TestEvaluateWorldHealth_Critical(t *testing.T) {
	rule := NewWorldBalanceRule()

	metrics := WorldMetrics{
		EntitySpiritStones:      []float64{1, 1, 1, 1, 10000},
		ResourceSpawnRate:       10,
		ResourceConsumptionRate: 100,
		ActiveSectCount:         1,
		SectSizes:               []int{1000},
		KarmaValues:             []float64{-100, -50, 0, 50, 200},
	}

	health := rule.EvaluateWorldHealth(metrics)
	assert.True(t, health.GiniCoefficient > 0.7)
	assert.InDelta(t, 0.1, health.ResourceRatio, 0.01)
	assert.True(t, health.OverallScore < 40)
	assert.Equal(t, "critical", health.Status)
}

func TestApplyBalanceAdjustment_Healthy(t *testing.T) {
	rule := NewWorldBalanceRule()

	health := WorldHealth{
		GiniCoefficient:    0.1,
		ResourceRatio:      1.0,
		SectDiversityScore: 0.8,
		KarmaStdDev:        5,
		OverallScore:       85,
		Status:             "healthy",
	}

	metrics := WorldMetrics{ActiveSectCount: 5}
	adjustments := rule.ApplyBalanceAdjustment(health, metrics)
	assert.Empty(t, adjustments)
}

func TestApplyBalanceAdjustment_MultipleIssues(t *testing.T) {
	rule := NewWorldBalanceRule()

	health := WorldHealth{
		GiniCoefficient:    0.7,
		ResourceRatio:      0.4,
		SectDiversityScore: 0.1,
		KarmaStdDev:        80,
		OverallScore:       15,
		Status:             "critical",
	}

	metrics := WorldMetrics{ActiveSectCount: 1}
	adjustments := rule.ApplyBalanceAdjustment(health, metrics)
	assert.True(t, len(adjustments) >= 3)

	// Check each adjustment type exists
	types := make(map[string]bool)
	for _, adj := range adjustments {
		types[adj.Type] = true
	}
	assert.True(t, types["wealth_tax"])
	assert.True(t, types["resource_boost"])
	assert.True(t, types["sect_promotion"])
	assert.True(t, types["karma_redistribution"])
}

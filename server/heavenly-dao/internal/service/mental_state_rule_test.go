package service

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestNewMentalStateRule(t *testing.T) {
    rule := NewMentalStateRule()
    assert.NotNil(t, rule)
    assert.Equal(t, 80, rule.highStabilityThresh)
    assert.Equal(t, 50, rule.lowStabilityThresh)
}

func TestCalculateMentalFactor_HighStability(t *testing.T) {
    rule := NewMentalStateRule()

    // >= 80: full factor
    assert.Equal(t, 1.0, rule.CalculateMentalFactor(80))
    assert.Equal(t, 1.0, rule.CalculateMentalFactor(90))
    assert.Equal(t, 1.0, rule.CalculateMentalFactor(100))
}

func TestCalculateMentalFactor_MediumStability(t *testing.T) {
    rule := NewMentalStateRule()

    // 65: midpoint between 50 and 80 → 0.5
    assert.InDelta(t, 0.5, rule.CalculateMentalFactor(65), 0.001)

    // 70: (70-50)/(80-50) = 20/30 ≈ 0.667
    assert.InDelta(t, 0.667, rule.CalculateMentalFactor(70), 0.001)

    // 55: (55-50)/30 = 5/30 ≈ 0.167
    assert.InDelta(t, 0.167, rule.CalculateMentalFactor(55), 0.001)
}

func TestCalculateMentalFactor_LowStability(t *testing.T) {
    rule := NewMentalStateRule()

    // <= 50: zero factor
    assert.Equal(t, 0.0, rule.CalculateMentalFactor(50))
    assert.Equal(t, 0.0, rule.CalculateMentalFactor(30))
    assert.Equal(t, 0.0, rule.CalculateMentalFactor(0))
}

func TestCalculateMentalRecovery_NoTime(t *testing.T) {
    rule := NewMentalStateRule()

    recovery := rule.CalculateMentalRecovery(50, 80, 0, 0)
    assert.Equal(t, 0.0, recovery)

    recovery = rule.CalculateMentalRecovery(50, 80, 0, -1)
    assert.Equal(t, 0.0, recovery)
}

func TestCalculateMentalRecovery_Basic(t *testing.T) {
    rule := NewMentalStateRule()

    // dao_heart = 100, no obsessions, 1 day
    // rate = 1.0 * 1.0 = 1.0 per day
    recovery := rule.CalculateMentalRecovery(50, 100, 0, 1)
    assert.InDelta(t, 1.0, recovery, 0.01)

    // 7 days
    recovery = rule.CalculateMentalRecovery(50, 100, 0, 7)
    assert.InDelta(t, 7.0, recovery, 0.01)
}

func TestCalculateMentalRecovery_DaoHeart(t *testing.T) {
    rule := NewMentalStateRule()

    // dao_heart = 50 → factor = 0.5
    // rate = 1.0 * 0.5 = 0.5 per day
    recovery := rule.CalculateMentalRecovery(50, 50, 0, 1)
    assert.InDelta(t, 0.5, recovery, 0.01)

    // dao_heart = 10 → factor = 0.1 (minimum)
    recovery = rule.CalculateMentalRecovery(50, 10, 0, 1)
    assert.InDelta(t, 0.1, recovery, 0.01)

    // dao_heart = 0 → factor clamped to 0.1
    recovery = rule.CalculateMentalRecovery(50, 0, 0, 1)
    assert.InDelta(t, 0.1, recovery, 0.01)
}

func TestCalculateMentalRecovery_ObsessionPenalty(t *testing.T) {
    rule := NewMentalStateRule()

    // 1 obsession: rate = 1.0 - 0.5 = 0.5 per day
    recovery := rule.CalculateMentalRecovery(50, 100, 1, 1)
    assert.InDelta(t, 0.5, recovery, 0.01)

    // 2 obsessions: rate = 1.0 - 1.0 = 0.0
    recovery = rule.CalculateMentalRecovery(50, 100, 2, 1)
    assert.InDelta(t, 0.0, recovery, 0.01)

    // 3 obsessions: rate = max(0, -0.5) = 0.0
    recovery = rule.CalculateMentalRecovery(50, 100, 3, 1)
    assert.InDelta(t, 0.0, recovery, 0.01)
}

func TestApplyMentalChange_Improvement(t *testing.T) {
    rule := NewMentalStateRule()

    result := rule.ApplyMentalChange(60, 20)
    assert.Equal(t, 80, result.NewStability)
    assert.True(t, result.Changed)
    assert.True(t, result.ThresholdCrossed) // crossed from <80 to >=80
    assert.Equal(t, "improved", result.Direction)
}

func TestApplyMentalChange_Degradation(t *testing.T) {
    rule := NewMentalStateRule()

    result := rule.ApplyMentalChange(85, -30)
    assert.Equal(t, 55, result.NewStability)
    assert.True(t, result.Changed)
    assert.True(t, result.ThresholdCrossed) // crossed from >=80 to <80
    assert.Equal(t, "degraded", result.Direction)
}

func TestApplyMentalChange_ClampBounds(t *testing.T) {
    rule := NewMentalStateRule()

    // Can't go above 100
    result := rule.ApplyMentalChange(95, 20)
    assert.Equal(t, 100, result.NewStability)

    // Can't go below 0
    result = rule.ApplyMentalChange(10, -50)
    assert.Equal(t, 0, result.NewStability)
}

func TestApplyMentalChange_NoChange(t *testing.T) {
    rule := NewMentalStateRule()

    result := rule.ApplyMentalChange(70, 0)
    assert.Equal(t, 70, result.NewStability)
    assert.False(t, result.Changed)
}

func TestApplyMentalChange_WithinZone(t *testing.T) {
    rule := NewMentalStateRule()

    // Stay within high zone
    result := rule.ApplyMentalChange(85, 5)
    assert.True(t, result.Changed)
    assert.False(t, result.ThresholdCrossed)

    // Stay within low zone
    result = rule.ApplyMentalChange(30, -10)
    assert.True(t, result.Changed)
    assert.False(t, result.ThresholdCrossed)
}

func TestNewInnerDemonRule(t *testing.T) {
    rule := NewInnerDemonRule()
    assert.NotNil(t, rule)
    assert.Equal(t, 0.001, rule.baseProbability)
    assert.Equal(t, 0.30, rule.qiDeviationProb)
    assert.Equal(t, 0.10, rule.victoryBonus)
}

func TestCalculateInnerDemonProbability_Base(t *testing.T) {
    rule := NewInnerDemonRule()

    input := InnerDemonInput{
        MentalStability:           50, // mental_factor = 0.5
        Karma:                     0,  // karma_factor = 1.0
        RecentBreakthroughAttempt: false,
        MethodConflictScore:       0.0,
    }

    prob := rule.CalculateInnerDemonProbability(input)
    // 0.001 * 0.5 * 1.0 * 1.0 * 1.0 = 0.0005
    assert.InDelta(t, 0.0005, prob, 0.00001)
}

func TestCalculateInnerDemonProbability_HighKarma(t *testing.T) {
    rule := NewInnerDemonRule()

    input := InnerDemonInput{
        MentalStability:           0,  // mental_factor = 1.0
        Karma:                     1000, // karma_factor = 2.0
        RecentBreakthroughAttempt: false,
        MethodConflictScore:       0.0,
    }

    prob := rule.CalculateInnerDemonProbability(input)
    // 0.001 * 1.0 * 2.0 * 1.0 * 1.0 = 0.002
    assert.InDelta(t, 0.002, prob, 0.00001)
}

func TestCalculateInnerDemonProbability_BreakthroughStress(t *testing.T) {
    rule := NewInnerDemonRule()

    input := InnerDemonInput{
        MentalStability:           0,
        Karma:                     0,
        RecentBreakthroughAttempt: true,
        MethodConflictScore:       0.0,
    }

    prob := rule.CalculateInnerDemonProbability(input)
    // 0.001 * 1.0 * 1.0 * 2.0 * 1.0 = 0.002
    assert.InDelta(t, 0.002, prob, 0.00001)
}

func TestCalculateInnerDemonProbability_MethodConflict(t *testing.T) {
    rule := NewInnerDemonRule()

    input := InnerDemonInput{
        MentalStability:           0,
        Karma:                     0,
        RecentBreakthroughAttempt: false,
        MethodConflictScore:       0.5, // conflict_factor = 1.0 + 0.5*2 = 2.0
    }

    prob := rule.CalculateInnerDemonProbability(input)
    // 0.001 * 1.0 * 1.0 * 1.0 * 2.0 = 0.002
    assert.InDelta(t, 0.002, prob, 0.00001)
}

func TestCalculateInnerDemonProbability_AllFactors(t *testing.T) {
    rule := NewInnerDemonRule()

    input := InnerDemonInput{
        MentalStability:           20, // mental_factor = 0.8
        Karma:                     500, // karma_factor = 1.5
        RecentBreakthroughAttempt: true, // breakthrough_stress = 2.0
        MethodConflictScore:       0.3, // conflict_factor = 1.6
    }

    prob := rule.CalculateInnerDemonProbability(input)
    // 0.001 * 0.8 * 1.5 * 2.0 * 1.6 = 0.00384
    assert.InDelta(t, 0.00384, prob, 0.00001)
}

func TestCalculateInnerDemonProbability_LowRisk(t *testing.T) {
    rule := NewInnerDemonRule()

    input := InnerDemonInput{
        MentalStability:           100, // mental_factor = 0.0
        Karma:                     0,
        RecentBreakthroughAttempt: false,
        MethodConflictScore:       0.0,
    }

    prob := rule.CalculateInnerDemonProbability(input)
    // mental_factor = 0 → prob = 0
    assert.Equal(t, 0.0, prob)
}

func TestCalculateInnerDemonStrength(t *testing.T) {
    rule := NewInnerDemonRule()

    // Realm level 0, no obsessions
    strength := rule.CalculateInnerDemonStrength(0, 0)
    assert.Equal(t, 0.0, strength)

    // Realm level 5 (Nascent Soul), 2 obsessions
    strength = rule.CalculateInnerDemonStrength(5, 2)
    // 5*100 + 2*50 = 600
    assert.Equal(t, 600.0, strength)

    // Realm level 10 (Tribulation), 5 obsessions
    strength = rule.CalculateInnerDemonStrength(10, 5)
    // 10*100 + 5*50 = 1250
    assert.Equal(t, 1250.0, strength)
}

func TestResolveInnerDemon_Victory(t *testing.T) {
    rule := NewInnerDemonRule()
    deterministicRand := func() float64 { return 0.1 } // always victory (resistance > 10%)

    // High resistance: 80 → victory_prob = 0.80
    result := rule.ResolveInnerDemon(80, deterministicRand)
    assert.True(t, result.Triggered)
    assert.Equal(t, "victory", result.Outcome)
    assert.Equal(t, "enlightenment", result.Effect)
}

func TestResolveInnerDemon_Defeat(t *testing.T) {
    rule := NewInnerDemonRule()
    deterministicRand := func() float64 { return 0.7 } // always > victory_prob for low resistance

    // Low resistance: 20 → victory_prob = 0.20
    result := rule.ResolveInnerDemon(20, deterministicRand)
    assert.True(t, result.Triggered)
    assert.Contains(t, []string{"defeat", "qi_deviation"}, result.Outcome)
}

func TestResolveInnerDemon_QiDeviation(t *testing.T) {
    rule := NewInnerDemonRule()
    callCount := 0
    deterministicRand := func() float64 {
        callCount++
        if callCount == 1 {
            return 0.9 // lose (0.9 > 0.20)
        }
        return 0.1 // qi_deviation (0.1 < 0.30)
    }

    result := rule.ResolveInnerDemon(20, deterministicRand)
    assert.True(t, result.Triggered)
    assert.Equal(t, "qi_deviation", result.Outcome)
    assert.Equal(t, "walk_fire_entrance", result.Effect)
}

func TestResolveInnerDemon_DefeatNoQiDeviation(t *testing.T) {
    rule := NewInnerDemonRule()
    callCount := 0
    deterministicRand := func() float64 {
        callCount++
        if callCount == 1 {
            return 0.9 // lose (0.9 > 0.20)
        }
        return 0.5 // no qi_deviation (0.5 > 0.30)
    }

    result := rule.ResolveInnerDemon(20, deterministicRand)
    assert.True(t, result.Triggered)
    assert.Equal(t, "defeat", result.Outcome)
    assert.Equal(t, "cultivation_regression", result.Effect)
}

func TestResolveInnerDemon_GuaranteedVictory(t *testing.T) {
    rule := NewInnerDemonRule()
    deterministicRand := func() float64 { return 0.0 }

    result := rule.ResolveInnerDemon(100, deterministicRand)
    assert.True(t, result.Triggered)
    assert.Equal(t, "victory", result.Outcome)
}

func TestResolveInnerDemon_GuaranteedDefeat(t *testing.T) {
    rule := NewInnerDemonRule()
    deterministicRand := func() float64 { return 1.0 }

    // resistance = 0 → victory_prob = 0.0, always lose
    // second rand = 1.0 > 0.30 → defeat, not qi_deviation
    result := rule.ResolveInnerDemon(0, deterministicRand)
    assert.True(t, result.Triggered)
    assert.Equal(t, "defeat", result.Outcome)
}

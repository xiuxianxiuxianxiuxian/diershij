package service

import (
    "testing"

    "github.com/cultivation-world/shared/types"
    "github.com/stretchr/testify/assert"
)

func TestNewBreakthroughRule(t *testing.T) {
    rule := NewBreakthroughRule()
    assert.NotNil(t, rule)
    assert.Equal(t, 9, len(rule.baseSuccessRateByRealm))
}

func TestCalculateBreakthroughSuccess_BaseRealms(t *testing.T) {
    rule := NewBreakthroughRule()

    // Baseline: perfect conditions except base rate
    input := BreakthroughInput{
        CultivationTime:    1.0,  // time_ratio = 1.0 (when required_time = 1.0)
        RequiredTime:       1.0,
        MethodQuality:      100,
        ResourceBonus:      0.0,
        MentalStability:    100,
        Luck:               50, // luck factor = 1.0
    }

    // For qi_condensation: base = 0.60, accumulation = min(1.5, 2.0) = 1.5
    // prob = 0.60 * 1.5 * 1.0 * 1.0 * 1.0 * 1.0 = 0.90 → clamped to 0.80
    input.TargetRealm = types.RealmQiCondensation
    prob := rule.CalculateBreakthroughSuccess(input)
    assert.InDelta(t, 0.80, prob, 0.001)

    // For foundation: base = 0.40, same factors
    // prob = 0.40 * 1.5 = 0.60
    input.TargetRealm = types.RealmFoundation
    prob = rule.CalculateBreakthroughSuccess(input)
    assert.InDelta(t, 0.60, prob, 0.001)

    // For golden_core: base = 0.25
    // prob = 0.25 * 1.5 = 0.375
    input.TargetRealm = types.RealmGoldenCore
    prob = rule.CalculateBreakthroughSuccess(input)
    assert.InDelta(t, 0.375, prob, 0.001)
}

func TestCalculateBreakthroughSuccess_HigherRealms(t *testing.T) {
    rule := NewBreakthroughRule()
    input := BreakthroughInput{
        CultivationTime:    1.0,
        RequiredTime:       1.0,
        MethodQuality:      100,
        ResourceBonus:      0.0,
        MentalStability:    100,
        Luck:               50,
    }

    realms := []struct {
        realm    types.CultivationRealm
        baseRate float64
    }{
        {types.RealmNascentSoul, 0.15},
        {types.RealmSoulTransform, 0.10},
        {types.RealmVoidRefinement, 0.08},
        {types.RealmIntegration, 0.05},
        {types.RealmMahayana, 0.03},
        {types.RealmTribulation, 0.01},
    }

    for _, tc := range realms {
        input.TargetRealm = tc.realm
        prob := rule.CalculateBreakthroughSuccess(input)
        expected := tc.baseRate * 1.5
        if expected > 0.80 {
            expected = 0.80
        }
        if expected < 0.05 {
            expected = 0.05
        }
        assert.InDelta(t, expected, prob, 0.001, "realm: %s", tc.realm)
    }
}

func TestCalculateBreakthroughSuccess_AccumulationTime(t *testing.T) {
    rule := NewBreakthroughRule()
    input := BreakthroughInput{
        TargetRealm:     types.RealmFoundation,
        MethodQuality:   100,
        ResourceBonus:   0.0,
        MentalStability: 100,
        Luck:            50,
    }

    // Half time: accumulation = 1.5
    input.CultivationTime = 0.5
    input.RequiredTime = 1.0
    prob := rule.CalculateBreakthroughSuccess(input)
    assert.InDelta(t, 0.60, prob, 0.001)

    // Full time: accumulation = min(1.5, 2.0) = 1.5
    input.CultivationTime = 1.0
    prob = rule.CalculateBreakthroughSuccess(input)
    assert.InDelta(t, 0.60, prob, 0.001)

    // Double time: accumulation = min(1.5, 3.0) = 1.5
    input.CultivationTime = 2.0
    prob = rule.CalculateBreakthroughSuccess(input)
    assert.InDelta(t, 0.60, prob, 0.001)

    // Zero time: accumulation = 1.0
    input.CultivationTime = 0.0
    prob = rule.CalculateBreakthroughSuccess(input)
    // 0.40 * 1.0 = 0.40
    assert.InDelta(t, 0.40, prob, 0.001)
}

func TestCalculateBreakthroughSuccess_MethodQuality(t *testing.T) {
    rule := NewBreakthroughRule()
    input := BreakthroughInput{
        TargetRealm:      types.RealmFoundation,
        CultivationTime:  1.0,
        RequiredTime:     1.0,
        ResourceBonus:    0.0,
        MentalStability:  100,
        Luck:             50,
    }

    // Perfect quality: method_factor = 1.0
    input.MethodQuality = 100
    prob := rule.CalculateBreakthroughSuccess(input)
    assert.InDelta(t, 0.60, prob, 0.001)

    // Half quality: method_factor = 0.5
    input.MethodQuality = 50
    prob = rule.CalculateBreakthroughSuccess(input)
    // 0.40 * 1.5 * 0.5 = 0.30
    assert.InDelta(t, 0.30, prob, 0.001)

    // Low quality: method_factor = 0.2
    input.MethodQuality = 20
    prob = rule.CalculateBreakthroughSuccess(input)
    // 0.40 * 1.5 * 0.2 = 0.12
    assert.InDelta(t, 0.12, prob, 0.001)
}

func TestCalculateBreakthroughSuccess_ResourceBonus(t *testing.T) {
    rule := NewBreakthroughRule()
    input := BreakthroughInput{
        TargetRealm:      types.RealmFoundation,
        CultivationTime:  1.0,
        RequiredTime:     1.0,
        MethodQuality:    100,
        MentalStability:  100,
        Luck:             50,
    }

    // No bonus
    input.ResourceBonus = 0.0
    prob := rule.CalculateBreakthroughSuccess(input)
    assert.InDelta(t, 0.60, prob, 0.001)

    // 20% bonus from pills
    input.ResourceBonus = 0.20
    prob = rule.CalculateBreakthroughSuccess(input)
    // 0.40 * 1.5 * 1.2 = 0.72
    assert.InDelta(t, 0.72, prob, 0.001)

    // 50% bonus
    input.ResourceBonus = 0.50
    prob = rule.CalculateBreakthroughSuccess(input)
    // 0.40 * 1.5 * 1.5 = 0.90 → clamped to 0.80
    assert.InDelta(t, 0.80, prob, 0.001)
}

func TestCalculateBreakthroughSuccess_MentalStability(t *testing.T) {
    rule := NewBreakthroughRule()
    input := BreakthroughInput{
        TargetRealm:      types.RealmFoundation,
        CultivationTime:  1.0,
        RequiredTime:     1.0,
        MethodQuality:    100,
        ResourceBonus:    0.0,
        Luck:             50,
    }

    // Perfect mental state
    input.MentalStability = 100
    prob := rule.CalculateBreakthroughSuccess(input)
    assert.InDelta(t, 0.60, prob, 0.001)

    // Half mental stability
    input.MentalStability = 50
    prob = rule.CalculateBreakthroughSuccess(input)
    // 0.40 * 1.5 * 0.5 = 0.30
    assert.InDelta(t, 0.30, prob, 0.001)

    // Zero mental stability
    input.MentalStability = 0
    prob = rule.CalculateBreakthroughSuccess(input)
    // 0.40 * 1.5 * 0 = 0 → clamped to min 0.05
    assert.InDelta(t, 0.05, prob, 0.001)
}

func TestCalculateBreakthroughSuccess_Luck(t *testing.T) {
    rule := NewBreakthroughRule()
    input := BreakthroughInput{
        TargetRealm:      types.RealmFoundation,
        CultivationTime:  1.0,
        RequiredTime:     1.0,
        MethodQuality:    100,
        ResourceBonus:    0.0,
        MentalStability:  100,
    }

    // Average luck (50): factor = 1.0
    input.Luck = 50
    prob := rule.CalculateBreakthroughSuccess(input)
    assert.InDelta(t, 0.60, prob, 0.001)

    // High luck (100): factor = 1.25
    input.Luck = 100
    prob = rule.CalculateBreakthroughSuccess(input)
    // 0.40 * 1.5 * 1.25 = 0.75
    assert.InDelta(t, 0.75, prob, 0.001)

    // Low luck (0): factor = 0.75
    input.Luck = 0
    prob = rule.CalculateBreakthroughSuccess(input)
    // 0.40 * 1.5 * 0.75 = 0.45
    assert.InDelta(t, 0.45, prob, 0.001)
}

func TestCalculateBreakthroughSuccess_MinClamp(t *testing.T) {
    rule := NewBreakthroughRule()
    input := BreakthroughInput{
        TargetRealm:      types.RealmTribulation,
        CultivationTime:  0.0,
        RequiredTime:     1.0,
        MethodQuality:    10,
        ResourceBonus:    0.0,
        MentalStability:  10,
        Luck:             0,
    }

    // base = 0.01, accumulation = 1.0, method = 0.1, mental = 0.1, luck = 0.75
    // prob = 0.01 * 1.0 * 0.1 * 1.0 * 0.1 * 0.75 = 0.000075 → clamped to 0.05
    prob := rule.CalculateBreakthroughSuccess(input)
    assert.InDelta(t, 0.05, prob, 0.001)
}

func TestCalculateBreakthroughSuccess_UnknownRealm(t *testing.T) {
    rule := NewBreakthroughRule()
    input := BreakthroughInput{
        TargetRealm:      types.CultivationRealm("unknown"),
        CultivationTime:  1.0,
        RequiredTime:     1.0,
        MethodQuality:    100,
        ResourceBonus:    0.0,
        MentalStability:  100,
        Luck:             50,
    }

    prob := rule.CalculateBreakthroughSuccess(input)
    assert.InDelta(t, 0.05, prob, 0.001) // min probability
}

func TestCalculateBreakthroughSuccess_ComplexScenario(t *testing.T) {
    rule := NewBreakthroughRule()

    // Golden Core breakthrough with:
    // - 75% cultivation time
    // - Good method (quality 80)
    // - Pill bonus (+15%)
    // - Good mental state (80)
    // - Average luck (50)
    input := BreakthroughInput{
        TargetRealm:      types.RealmGoldenCore,
        CultivationTime:  0.75,
        RequiredTime:     1.0,
        MethodQuality:    80,
        ResourceBonus:    0.15,
        MentalStability:  80,
        Luck:             50,
    }

    // base = 0.25, accumulation = 1.75 → clamped to 1.5
    // method = 0.80, resource = 1.15, mental = 0.80, luck = 1.0
    // prob = 0.25 * 1.5 * 0.80 * 1.15 * 0.80 * 1.0 = 0.276
    prob := rule.CalculateBreakthroughSuccess(input)
    assert.InDelta(t, 0.276, prob, 0.001)
}

func TestCalculateFailurePenalty_Default(t *testing.T) {
    rule := NewBreakthroughRule()

    // Foundation (level 2): cooldown = 24 * 2 = 48 hours
    penalty := rule.CalculateFailurePenalty(2)
    assert.Equal(t, 0.10, penalty.ProgressLoss)
    assert.InDelta(t, 48.0, penalty.CooldownHours, 0.1)
    assert.Equal(t, 20, penalty.MentalDamage)
    assert.InDelta(t, 0.30, penalty.InjuryProb, 0.001)
    assert.Equal(t, "half_hp", penalty.InjurySeverity)
}

func TestCalculateFailurePenalty_HigherRealm(t *testing.T) {
    rule := NewBreakthroughRule()

    // Nascent Soul (level 4): cooldown = 24 * 4 = 96 hours
    penalty := rule.CalculateFailurePenalty(4)
    assert.InDelta(t, 96.0, penalty.CooldownHours, 0.1)
}

func TestCalculateFailurePenalty_ZeroLevel(t *testing.T) {
    rule := NewBreakthroughRule()

    // Mortal (level 0): cooldown = 0 hours
    penalty := rule.CalculateFailurePenalty(0)
    assert.InDelta(t, 0.0, penalty.CooldownHours, 0.1)
}

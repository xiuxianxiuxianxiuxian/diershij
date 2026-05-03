package service

import (
    "testing"

    "github.com/cultivation-world/shared/types"
    "github.com/stretchr/testify/assert"
)

func TestNewTribulationRule(t *testing.T) {
    rule := NewTribulationRule()
    assert.NotNil(t, rule)
    assert.Equal(t, 9, len(rule.baseProbabilityByRealm))
}

func TestCalculateTribulationProbability_QiCondensation(t *testing.T) {
    rule := NewTribulationRule()
    input := DefaultTribulationInput()
    input.TargetRealm = types.RealmQiCondensation

    prob := rule.CalculateTribulationProbability(input)
    assert.Equal(t, 0.10, prob)
}

func TestCalculateTribulationProbability_HigherRealms(t *testing.T) {
    rule := NewTribulationRule()
    input := DefaultTribulationInput()

    realms := []struct {
        realm    types.CultivationRealm
        expected float64
    }{
        {types.RealmFoundation, 0.15},
        {types.RealmGoldenCore, 0.20},
        {types.RealmNascentSoul, 0.30},
        {types.RealmSoulTransform, 0.50},
        {types.RealmVoidRefinement, 0.60},
        {types.RealmIntegration, 0.70},
        {types.RealmMahayana, 0.80},
        {types.RealmTribulation, 0.90},
    }

    for _, tc := range realms {
        input.TargetRealm = tc.realm
        prob := rule.CalculateTribulationProbability(input)
        assert.Equal(t, tc.expected, prob, "realm: %s", tc.realm)
    }
}

func TestCalculateTribulationProbability_KarmaIncrease(t *testing.T) {
    rule := NewTribulationRule()
    input := DefaultTribulationInput()
    input.TargetRealm = types.RealmFoundation

    input.Karma = 500
    prob := rule.CalculateTribulationProbability(input)
    assert.InDelta(t, 0.30, prob, 0.001)

    input.Karma = 1000
    prob = rule.CalculateTribulationProbability(input)
    assert.InDelta(t, 0.45, prob, 0.001)

    input.Karma = 5000
    prob = rule.CalculateTribulationProbability(input)
    assert.InDelta(t, 1.0, prob, 0.001)
}

func TestCalculateTribulationProbability_MeritDecrease(t *testing.T) {
    rule := NewTribulationRule()
    input := DefaultTribulationInput()
    input.TargetRealm = types.RealmSoulTransform

    input.Merit = 2000
    input.Karma = 0
    prob := rule.CalculateTribulationProbability(input)
    assert.InDelta(t, 0.45, prob, 0.001)

    input.Merit = 4000
    prob = rule.CalculateTribulationProbability(input)
    assert.InDelta(t, 0.40, prob, 0.001)

    input.Merit = 20000
    prob = rule.CalculateTribulationProbability(input)
    assert.InDelta(t, 0.15, prob, 0.001)
}

func TestCalculateTribulationProbability_MinClamp(t *testing.T) {
    rule := NewTribulationRule()
    input := DefaultTribulationInput()
    input.TargetRealm = types.RealmQiCondensation
    input.Karma = 0
    input.Merit = 200000

    prob := rule.CalculateTribulationProbability(input)
    // merit floor 0.30 * base 0.10 = 0.03, which is above min 0.01
    assert.InDelta(t, 0.03, prob, 0.001)
}

func TestCalculateTribulationProbability_UnknownRealm(t *testing.T) {
    rule := NewTribulationRule()
    input := DefaultTribulationInput()
    input.TargetRealm = types.CultivationRealm("unknown_realm")

    prob := rule.CalculateTribulationProbability(input)
    assert.InDelta(t, 0.01, prob, 0.001)
}

func TestCalculateTribulationStrength_NoKarma(t *testing.T) {
    rule := NewTribulationRule()

    strength := rule.CalculateTribulationStrength(0, 0)
    assert.InDelta(t, 100.0, strength, 0.1)
}

func TestCalculateTribulationStrength_WithKarma(t *testing.T) {
    rule := NewTribulationRule()

    strength := rule.CalculateTribulationStrength(500, 0)
    assert.InDelta(t, 271.5, strength, 1.0)

    strength = rule.CalculateTribulationStrength(1000, 0)
    assert.InDelta(t, 737.0, strength, 10.0)
}

func TestCalculateTribulationStrength_RecentBreakthroughs(t *testing.T) {
    rule := NewTribulationRule()

    strength := rule.CalculateTribulationStrength(0, 3)
    assert.InDelta(t, 130.0, strength, 0.1)

    strength = rule.CalculateTribulationStrength(0, 10)
    assert.InDelta(t, 200.0, strength, 0.1)
}

func TestCalculateTribulationStrength_Combined(t *testing.T) {
    rule := NewTribulationRule()

    strength := rule.CalculateTribulationStrength(5000, 5)
    // 1.002^5000 ≈ 21915, times 100, times 1.5
    expected := 100.0 * 21915.0 * 1.5
    assert.InDelta(t, expected, strength, expected*0.05)
}

func TestAssess_Triggered(t *testing.T) {
    rule := NewTribulationRule()
    input := DefaultTribulationInput()
    input.TargetRealm = types.RealmSoulTransform

    input.Karma = 100
    result := rule.Assess(input)
    assert.True(t, result.Triggered)
    assert.True(t, result.Probability >= 0.5)
}

func TestAssess_NotTriggered(t *testing.T) {
    rule := NewTribulationRule()
    input := DefaultTribulationInput()
    input.TargetRealm = types.RealmQiCondensation

    input.Karma = 0
    result := rule.Assess(input)
    assert.False(t, result.Triggered)
    assert.True(t, result.Probability < 0.5)
}

func TestAssess_StrengthIncluded(t *testing.T) {
    rule := NewTribulationRule()
    input := DefaultTribulationInput()
    input.TargetRealm = types.RealmNascentSoul
    input.Karma = 500
    input.RecentBreakthroughs = 2

    result := rule.Assess(input)
    assert.True(t, result.Strength > 100.0)
    assert.True(t, result.Probability > 0.30)
}

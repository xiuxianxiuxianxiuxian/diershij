package service

import (
    "testing"

    "github.com/cultivation-world/shared/types"
    "github.com/stretchr/testify/assert"
)

func TestNewLifespanRule(t *testing.T) {
    rule := NewLifespanRule()
    assert.NotNil(t, rule)
    assert.Equal(t, 10, len(rule.baseLifespanByRealm))
}

func TestCalculateBaseLifespan_AllRealms(t *testing.T) {
    rule := NewLifespanRule()

    realms := []struct {
        realm    types.CultivationRealm
        expected int
    }{
        {types.RealmMortal, 80},
        {types.RealmQiCondensation, 120},
        {types.RealmFoundation, 200},
        {types.RealmGoldenCore, 500},
        {types.RealmNascentSoul, 1000},
        {types.RealmSoulTransform, 3000},
        {types.RealmVoidRefinement, 5000},
        {types.RealmIntegration, 8000},
        {types.RealmMahayana, 10000},
        {types.RealmTribulation, 15000},
    }

    for _, tc := range realms {
        lifespan := rule.CalculateBaseLifespan(tc.realm)
        assert.Equal(t, tc.expected, lifespan, "realm: %s", tc.realm)
    }
}

func TestCalculateBaseLifespan_UnknownRealm(t *testing.T) {
    rule := NewLifespanRule()
    lifespan := rule.CalculateBaseLifespan(types.CultivationRealm("unknown"))
    assert.Equal(t, 0, lifespan)
}

func TestCalculateRemainingLifespan_Young(t *testing.T) {
    rule := NewLifespanRule()

    // Foundation realm: 200 years base, age 50 → 150 remaining
    remaining := rule.CalculateRemainingLifespan(types.RealmFoundation, 50)
    assert.Equal(t, 150, remaining)
}

func TestCalculateRemainingLifespan_MiddleAge(t *testing.T) {
    rule := NewLifespanRule()

    // Golden Core: 500 years base, age 300 → 200 remaining
    remaining := rule.CalculateRemainingLifespan(types.RealmGoldenCore, 300)
    assert.Equal(t, 200, remaining)
}

func TestCalculateRemainingLifespan_Exhausted(t *testing.T) {
    rule := NewLifespanRule()

    // Mortal: 80 years base, age 80 → 0 remaining
    remaining := rule.CalculateRemainingLifespan(types.RealmMortal, 80)
    assert.Equal(t, 0, remaining)

    // Over age
    remaining = rule.CalculateRemainingLifespan(types.RealmMortal, 100)
    assert.Equal(t, 0, remaining)
}

func TestCalculateRemainingLifespanRatio(t *testing.T) {
    rule := NewLifespanRule()

    // Nascent Soul: 1000 years, age 200 → 800 remaining → 0.8 ratio
    ratio := rule.CalculateRemainingLifespanRatio(types.RealmNascentSoul, 200)
    assert.InDelta(t, 0.8, ratio, 0.001)

    // Exact half
    ratio = rule.CalculateRemainingLifespanRatio(types.RealmFoundation, 100)
    assert.InDelta(t, 0.5, ratio, 0.001)

    // Exhausted
    ratio = rule.CalculateRemainingLifespanRatio(types.RealmMortal, 80)
    assert.InDelta(t, 0.0, ratio, 0.001)
}

func TestCalculateAgingPenalty_NoPenalty(t *testing.T) {
    rule := NewLifespanRule()

    // > 50% remaining: no penalty
    assert.InDelta(t, 0.0, rule.CalculateAgingPenalty(0.51), 0.001)
    assert.InDelta(t, 0.0, rule.CalculateAgingPenalty(0.80), 0.001)
    assert.InDelta(t, 0.0, rule.CalculateAgingPenalty(1.0), 0.001)
}

func TestCalculateAgingPenalty_20To50Percent(t *testing.T) {
    rule := NewLifespanRule()

    // Exactly at 50%: no penalty
    assert.InDelta(t, 0.0, rule.CalculateAgingPenalty(0.50), 0.001)

    // At 20%: full threshold_20pct penalty
    assert.InDelta(t, 0.09, rule.CalculateAgingPenalty(0.20), 0.001)

    // At 35% (midpoint): half of threshold_20pct
    assert.InDelta(t, 0.045, rule.CalculateAgingPenalty(0.35), 0.001)
}

func TestCalculateAgingPenalty_10To20Percent(t *testing.T) {
    rule := NewLifespanRule()

    // At 20%: threshold_20pct
    assert.InDelta(t, 0.09, rule.CalculateAgingPenalty(0.20), 0.001)

    // At 10%: threshold_10pct
    assert.InDelta(t, 0.14, rule.CalculateAgingPenalty(0.10), 0.001)

    // At 15% (midpoint): between 0.09 and 0.14
    penalty := rule.CalculateAgingPenalty(0.15)
    assert.True(t, penalty > 0.09 && penalty < 0.14)
}

func TestCalculateAgingPenalty_Below10Percent(t *testing.T) {
    rule := NewLifespanRule()

    // Below 10%: final penalty
    assert.InDelta(t, 0.20, rule.CalculateAgingPenalty(0.09), 0.001)
    assert.InDelta(t, 0.20, rule.CalculateAgingPenalty(0.01), 0.001)
    assert.InDelta(t, 0.20, rule.CalculateAgingPenalty(0.0), 0.001)
}

func TestCalculateCultivationPenaltyMultiplier(t *testing.T) {
    rule := NewLifespanRule()

    // No penalty
    assert.InDelta(t, 1.0, rule.CalculateCultivationPenaltyMultiplier(0.60), 0.001)

    // Some penalty
    mult := rule.CalculateCultivationPenaltyMultiplier(0.15)
    assert.True(t, mult > 0.80 && mult < 0.95)

    // Max penalty
    assert.InDelta(t, 0.80, rule.CalculateCultivationPenaltyMultiplier(0.05), 0.001)
}

func TestHandleLifespanDepletion_ForcedAttempt(t *testing.T) {
    rule := NewLifespanRule()

    // Progress at threshold: forced attempt
    result := rule.HandleLifespanDepletion(0.80, 0.80)
    assert.True(t, result.Depleted)
    assert.True(t, result.ForcedAttempt)
    assert.False(t, result.Death)

    // Progress above threshold: forced attempt
    result = rule.HandleLifespanDepletion(0.90, 0.80)
    assert.True(t, result.ForcedAttempt)
    assert.False(t, result.Death)
}

func TestHandleLifespanDepletion_Death(t *testing.T) {
    rule := NewLifespanRule()

    // Progress below threshold: death
    result := rule.HandleLifespanDepletion(0.50, 0.80)
    assert.True(t, result.Depleted)
    assert.False(t, result.ForcedAttempt)
    assert.True(t, result.Death)

    // Zero progress: death
    result = rule.HandleLifespanDepletion(0.0, 0.80)
    assert.True(t, result.Death)
}

func TestCalculateRemainingLifespan_EdgeCases(t *testing.T) {
    rule := NewLifespanRule()

    // Age 0: full remaining
    remaining := rule.CalculateRemainingLifespan(types.RealmFoundation, 0)
    assert.Equal(t, 200, remaining)

    // Negative age: treat as 0
    remaining = rule.CalculateRemainingLifespan(types.RealmFoundation, -10)
    assert.Equal(t, 210, remaining)
}

package service

import (
    "testing"

    "github.com/cultivation-world/shared/types"
    "github.com/stretchr/testify/assert"
)

func TestNewCultivationEfficiencyRule(t *testing.T) {
    rule := NewCultivationEfficiencyRule()
    assert.NotNil(t, rule)
    assert.Equal(t, 0.1, rule.baseRateMultiplier)
    assert.Equal(t, 0.2, rule.realmPenaltyFactor)
}

func TestCalculateCultivationRate_Baseline(t *testing.T) {
    rule := NewCultivationEfficiencyRule()
    input := CultivationRateInput{
        Comprehension:   80,
        SpiritualDensity: 80,
        MethodMatch:     1.0,
        RealmLevel:      0,
        MentalState:     100,
        AgingPenalty:    0.0,
    }

    // base_rate = 80 * 0.1 = 8.0
    // spiritual = 80/100 = 0.8
    // method = 1.0
    // realm = 1/(1+0*0.2) = 1.0
    // mental = 100/100 = 1.0
    // aging = 1.0
    // result = 8.0 * 0.8 * 1.0 * 1.0 * 1.0 * 1.0 = 6.4
    rate := rule.CalculateCultivationRate(input)
    assert.InDelta(t, 6.4, rate, 0.01)
}

func TestCalculateCultivationRate_LowSpiritualDensity(t *testing.T) {
    rule := NewCultivationEfficiencyRule()
    input := CultivationRateInput{
        Comprehension:   80,
        SpiritualDensity: 20,
        MethodMatch:     1.0,
        RealmLevel:      0,
        MentalState:     100,
        AgingPenalty:    0.0,
    }

    // spiritual = 0.2, so rate = 8.0 * 0.2 * 1.0 * 1.0 * 1.0 * 1.0 = 1.6
    rate := rule.CalculateCultivationRate(input)
    assert.InDelta(t, 1.6, rate, 0.01)
}

func TestCalculateCultivationRate_MethodMismatch(t *testing.T) {
    rule := NewCultivationEfficiencyRule()
    input := CultivationRateInput{
        Comprehension:   80,
        SpiritualDensity: 80,
        MethodMatch:     0.5,
        RealmLevel:      0,
        MentalState:     100,
        AgingPenalty:    0.0,
    }

    // method = 0.5, rate = 6.4 * 0.5 = 3.2
    rate := rule.CalculateCultivationRate(input)
    assert.InDelta(t, 3.2, rate, 0.01)
}

func TestCalculateCultivationRate_HigherRealm(t *testing.T) {
    rule := NewCultivationEfficiencyRule()
    input := CultivationRateInput{
        Comprehension:   80,
        SpiritualDensity: 80,
        MethodMatch:     1.0,
        RealmLevel:      4, // Nascent Soul (level 4)
        MentalState:     100,
        AgingPenalty:    0.0,
    }

    // realm_penalty = 1/(1+4*0.2) = 1/1.8 ≈ 0.5556
    // rate = 6.4 * 0.5556 ≈ 3.556
    rate := rule.CalculateCultivationRate(input)
    assert.True(t, rate > 3.4 && rate < 3.7)
}

func TestCalculateCultivationRate_PoorMentalState(t *testing.T) {
    rule := NewCultivationEfficiencyRule()
    input := CultivationRateInput{
        Comprehension:   80,
        SpiritualDensity: 80,
        MethodMatch:     1.0,
        RealmLevel:      0,
        MentalState:     30,
        AgingPenalty:    0.0,
    }

    // mental = 30/100 = 0.3, rate = 6.4 * 0.3 = 1.92
    rate := rule.CalculateCultivationRate(input)
    assert.InDelta(t, 1.92, rate, 0.01)
}

func TestCalculateCultivationRate_AgingPenalty(t *testing.T) {
    rule := NewCultivationEfficiencyRule()
    input := CultivationRateInput{
        Comprehension:   80,
        SpiritualDensity: 80,
        MethodMatch:     1.0,
        RealmLevel:      0,
        MentalState:     100,
        AgingPenalty:    0.09,
    }

    // aging = 1 - 0.09 = 0.91, rate = 6.4 * 0.91 = 5.824
    rate := rule.CalculateCultivationRate(input)
    assert.InDelta(t, 5.824, rate, 0.01)
}

func TestCalculateCultivationRate_AllFactors(t *testing.T) {
    rule := NewCultivationEfficiencyRule()
    input := CultivationRateInput{
        Comprehension:   80,
        SpiritualDensity: 50,
        MethodMatch:     0.5,
        RealmLevel:      3, // Golden Core
        MentalState:     70,
        AgingPenalty:    0.05,
    }

    // base = 8.0, spiritual = 0.5, method = 0.5
    // realm = 1/(1+3*0.2) = 1/1.6 = 0.625
    // mental = 0.7, aging = 0.95
    // rate = 8.0 * 0.5 * 0.5 * 0.625 * 0.7 * 0.95
    expected := 8.0 * 0.5 * 0.5 * 0.625 * 0.7 * 0.95
    rate := rule.CalculateCultivationRate(input)
    assert.InDelta(t, expected, rate, 0.001)
}

func TestCalculateMethodCompatibility_PerfectMatch(t *testing.T) {
    rule := NewCultivationEfficiencyRule()
    roots := []types.SpiritualRoot{
        {Element: "fire", Purity: 90},
        {Element: "metal", Purity: 70},
    }
    required := []string{"fire", "metal"}

    score := rule.CalculateMethodCompatibility(roots, required)
    assert.Equal(t, 1.0, score)
}

func TestCalculateMethodCompatibility_PartialMatch(t *testing.T) {
    rule := NewCultivationEfficiencyRule()
    roots := []types.SpiritualRoot{
        {Element: "fire", Purity: 90},
        {Element: "water", Purity: 70},
    }
    required := []string{"fire", "metal", "wood"}

    score := rule.CalculateMethodCompatibility(roots, required)
    assert.InDelta(t, 1.0/3.0, score, 0.001)
}

func TestCalculateMethodCompatibility_NoMatch(t *testing.T) {
    rule := NewCultivationEfficiencyRule()
    roots := []types.SpiritualRoot{
        {Element: "fire", Purity: 90},
    }
    required := []string{"water", "wood"}

    score := rule.CalculateMethodCompatibility(roots, required)
    assert.Equal(t, 0.0, score)
}

func TestCalculateMethodCompatibility_NoRequirements(t *testing.T) {
    rule := NewCultivationEfficiencyRule()
    roots := []types.SpiritualRoot{
        {Element: "fire", Purity: 90},
    }
    required := []string{}

    score := rule.CalculateMethodCompatibility(roots, required)
    assert.Equal(t, 1.0, score)
}

func TestCalculateMethodCompatibility_EmptyRoots(t *testing.T) {
    rule := NewCultivationEfficiencyRule()
    roots := []types.SpiritualRoot{}
    required := []string{"fire"}

    score := rule.CalculateMethodCompatibility(roots, required)
    assert.Equal(t, 0.0, score)
}

func TestGetRealmLevel_AllRealms(t *testing.T) {
    rule := NewCultivationEfficiencyRule()

    expectedLevels := []struct {
        realm types.CultivationRealm
        level int
    }{
        {types.RealmMortal, 0},
        {types.RealmQiCondensation, 1},
        {types.RealmFoundation, 2},
        {types.RealmGoldenCore, 3},
        {types.RealmNascentSoul, 4},
        {types.RealmSoulTransform, 5},
        {types.RealmVoidRefinement, 6},
        {types.RealmIntegration, 7},
        {types.RealmMahayana, 8},
        {types.RealmTribulation, 9},
    }

    for _, tc := range expectedLevels {
        level := rule.GetRealmLevel(tc.realm)
        assert.Equal(t, tc.level, level, "realm: %s", tc.realm)
    }
}

func TestGetRealmLevel_Unknown(t *testing.T) {
    rule := NewCultivationEfficiencyRule()
    level := rule.GetRealmLevel(types.CultivationRealm("unknown"))
    assert.Equal(t, 0, level)
}

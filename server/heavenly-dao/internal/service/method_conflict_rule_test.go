package service

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestNewMethodConflictRule(t *testing.T) {
    rule := NewMethodConflictRule()
    assert.NotNil(t, rule)
    assert.Equal(t, 0.5, rule.oppositeElementConflict)
    assert.Equal(t, 0.3, rule.alignmentConflict)
    assert.Equal(t, 0.2, rule.realmGapConflict)
}

func TestCalculateMethodConflict_NoConflict(t *testing.T) {
    rule := NewMethodConflictRule()

    a := MethodInfo{ID: "fire_method", Element: "fire", Alignment: "righteous", RealmRequired: 1}
    b := MethodInfo{ID: "wood_method", Element: "wood", Alignment: "righteous", RealmRequired: 1}

    conflict := rule.CalculateMethodConflict(a, b)
    assert.Equal(t, 0.0, conflict.Score)
    assert.Equal(t, "none", conflict.ConflictType)
}

func TestCalculateMethodConflict_ElementOpposite(t *testing.T) {
    rule := NewMethodConflictRule()

    // Fire vs Water → opposite
    a := MethodInfo{ID: "fire_method", Element: "fire", Alignment: "neutral", RealmRequired: 1}
    b := MethodInfo{ID: "water_method", Element: "water", Alignment: "neutral", RealmRequired: 1}

    conflict := rule.CalculateMethodConflict(a, b)
    assert.Equal(t, 0.5, conflict.Score)
    assert.Equal(t, "element", conflict.ConflictType)

    // Wood vs Metal → opposite
    a = MethodInfo{ID: "wood_method", Element: "wood", Alignment: "neutral", RealmRequired: 1}
    b = MethodInfo{ID: "metal_method", Element: "metal", Alignment: "neutral", RealmRequired: 1}

    conflict = rule.CalculateMethodConflict(a, b)
    assert.Equal(t, 0.5, conflict.Score)
}

func TestCalculateMethodConflict_AlignmentConflict(t *testing.T) {
    rule := NewMethodConflictRule()

    // Righteous vs Demonic
    a := MethodInfo{ID: "righteous_method", Element: "fire", Alignment: "righteous", RealmRequired: 1}
    b := MethodInfo{ID: "demonic_method", Element: "fire", Alignment: "demonic", RealmRequired: 1}

    conflict := rule.CalculateMethodConflict(a, b)
    assert.Equal(t, 0.3, conflict.Score)
    assert.Equal(t, "alignment", conflict.ConflictType)

    // Same alignment → no conflict
    a = MethodInfo{ID: "method_a", Element: "fire", Alignment: "righteous", RealmRequired: 1}
    b = MethodInfo{ID: "method_b", Element: "fire", Alignment: "righteous", RealmRequired: 1}

    conflict = rule.CalculateMethodConflict(a, b)
    assert.Equal(t, 0.0, conflict.Score)
}

func TestCalculateMethodConflict_RealmGap(t *testing.T) {
    rule := NewMethodConflictRule()

    // Realm gap > 2
    a := MethodInfo{ID: "low_method", Element: "fire", Alignment: "neutral", RealmRequired: 0}
    b := MethodInfo{ID: "high_method", Element: "fire", Alignment: "neutral", RealmRequired: 3}

    conflict := rule.CalculateMethodConflict(a, b)
    assert.Equal(t, 0.2, conflict.Score)
    assert.Equal(t, "realm", conflict.ConflictType)

    // Realm gap <= 2 → no conflict
    b.RealmRequired = 2
    conflict = rule.CalculateMethodConflict(a, b)
    assert.Equal(t, 0.0, conflict.Score)
}

func TestCalculateMethodConflict_MultipleConflicts(t *testing.T) {
    rule := NewMethodConflictRule()

    // Fire/Water (element opposite) + Righteous/Demonic (alignment) + realm gap > 2
    a := MethodInfo{ID: "fire_righteous", Element: "fire", Alignment: "righteous", RealmRequired: 0}
    b := MethodInfo{ID: "water_demonic", Element: "water", Alignment: "demonic", RealmRequired: 4}

    conflict := rule.CalculateMethodConflict(a, b)
    // 0.5 (element) + 0.3 (alignment) + 0.2 (realm) = 1.0
    assert.Equal(t, 1.0, conflict.Score)
    assert.Equal(t, "multiple", conflict.ConflictType)
}

func TestCalculateMethodConflict_EmptyElement(t *testing.T) {
    rule := NewMethodConflictRule()

    a := MethodInfo{ID: "method_a", Element: "", Alignment: "neutral", RealmRequired: 1}
    b := MethodInfo{ID: "method_b", Element: "fire", Alignment: "neutral", RealmRequired: 1}

    conflict := rule.CalculateMethodConflict(a, b)
    assert.Equal(t, 0.0, conflict.Score)
}

func TestCalculateOverallConflict_SingleMethod(t *testing.T) {
    rule := NewMethodConflictRule()

    methods := []MethodInfo{
        {ID: "fire_method", Element: "fire", Alignment: "righteous", RealmRequired: 1},
    }

    result := rule.CalculateOverallConflict(methods)
    assert.Equal(t, 0.0, result.NormalizedScore)
    assert.Equal(t, 0.0, result.BacklashProb)
    assert.Equal(t, "none", result.Severity)
}

func TestCalculateOverallConflict_NoConflict(t *testing.T) {
    rule := NewMethodConflictRule()

    methods := []MethodInfo{
        {ID: "fire_a", Element: "fire", Alignment: "righteous", RealmRequired: 1},
        {ID: "fire_b", Element: "fire", Alignment: "righteous", RealmRequired: 1},
    }

    result := rule.CalculateOverallConflict(methods)
    assert.Equal(t, 0.0, result.NormalizedScore)
    assert.Equal(t, 0.0, result.BacklashProb)
    assert.Equal(t, "none", result.Severity)
}

func TestCalculateOverallConflict_ElementConflict(t *testing.T) {
    rule := NewMethodConflictRule()

    methods := []MethodInfo{
        {ID: "fire_method", Element: "fire", Alignment: "neutral", RealmRequired: 1},
        {ID: "water_method", Element: "water", Alignment: "neutral", RealmRequired: 1},
    }

    result := rule.CalculateOverallConflict(methods)
    // 1 pair, score = 0.5, normalized = 0.5/1 = 0.5
    assert.InDelta(t, 0.5, result.NormalizedScore, 0.001)
    assert.Equal(t, 1, len(result.Pairs))
    assert.Equal(t, "moderate", result.Severity)
}

func TestCalculateOverallConflict_ThreeMethods(t *testing.T) {
    rule := NewMethodConflictRule()

    methods := []MethodInfo{
        {ID: "fire_method", Element: "fire", Alignment: "neutral", RealmRequired: 0},
        {ID: "water_method", Element: "water", Alignment: "neutral", RealmRequired: 0},
        {ID: "earth_method", Element: "earth", Alignment: "neutral", RealmRequired: 0},
    }

    // Pairs: fire-water (0.5, opposite), fire-earth (0, no conflict), water-earth (0, no conflict)
    // Total = 0.5, max_pairs = 3, normalized = 0.5/3 ≈ 0.167
    result := rule.CalculateOverallConflict(methods)
    assert.InDelta(t, 0.167, result.NormalizedScore, 0.01)
    assert.Equal(t, 3, len(result.Pairs))
    assert.Equal(t, "none", result.Severity) // normalized < 0.2, backlash = 0.0
}

func TestCalculateOverallConflict_AllConflicts(t *testing.T) {
    rule := NewMethodConflictRule()

    methods := []MethodInfo{
        {ID: "fire_righteous", Element: "fire", Alignment: "righteous", RealmRequired: 0},
        {ID: "water_demonic", Element: "water", Alignment: "demonic", RealmRequired: 4},
    }

    // score = 0.5 + 0.3 + 0.2 = 1.0, normalized = 1.0
    result := rule.CalculateOverallConflict(methods)
    assert.InDelta(t, 1.0, result.NormalizedScore, 0.001)
    assert.Equal(t, 0.6, result.BacklashProb)
    assert.Equal(t, "severe", result.Severity)
}

func TestCalculateCultivationPenalty_Tiers(t *testing.T) {
    rule := NewMethodConflictRule()

    // < 0.2: no penalty
    assert.Equal(t, 0.0, rule.CalculateCultivationPenalty(0.0))
    assert.Equal(t, 0.0, rule.CalculateCultivationPenalty(0.19))

    // 0.2-0.5: -20%
    assert.Equal(t, 0.20, rule.CalculateCultivationPenalty(0.20))
    assert.Equal(t, 0.20, rule.CalculateCultivationPenalty(0.49))

    // 0.5-0.8: -50%
    assert.Equal(t, 0.50, rule.CalculateCultivationPenalty(0.50))
    assert.Equal(t, 0.50, rule.CalculateCultivationPenalty(0.79))

    // > 0.8: -80%
    assert.Equal(t, 0.80, rule.CalculateCultivationPenalty(0.80))
    assert.Equal(t, 0.80, rule.CalculateCultivationPenalty(1.0))
}

func TestCalculateBacklashProbability_Tiers(t *testing.T) {
    rule := NewMethodConflictRule()

    assert.Equal(t, 0.0, rule.calculateBacklashProbability(0.0))
    assert.Equal(t, 0.0, rule.calculateBacklashProbability(0.19))

    assert.Equal(t, 0.1, rule.calculateBacklashProbability(0.20))
    assert.Equal(t, 0.1, rule.calculateBacklashProbability(0.49))

    assert.Equal(t, 0.3, rule.calculateBacklashProbability(0.50))
    assert.Equal(t, 0.3, rule.calculateBacklashProbability(0.79))

    assert.Equal(t, 0.6, rule.calculateBacklashProbability(0.80))
    assert.Equal(t, 0.6, rule.calculateBacklashProbability(1.0))
}

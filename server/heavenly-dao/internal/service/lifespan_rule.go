package service

import (
    "math"

    "github.com/cultivation-world/shared/types"
)

// LifespanRule handles lifespan calculation, aging penalty, and depletion handling.
type LifespanRule struct {
    baseLifespanByRealm map[types.CultivationRealm]int
    agingPenalty50Pct   float64
    agingPenalty20Pct   float64
    agingPenalty10Pct   float64
    agingPenaltyFinal   float64
}

// NewLifespanRule creates a new LifespanRule with default configuration.
func NewLifespanRule() *LifespanRule {
    return &LifespanRule{
        baseLifespanByRealm: map[types.CultivationRealm]int{
            types.RealmMortal:         80,
            types.RealmQiCondensation: 120,
            types.RealmFoundation:     200,
            types.RealmGoldenCore:     500,
            types.RealmNascentSoul:    1000,
            types.RealmSoulTransform:  3000,
            types.RealmVoidRefinement: 5000,
            types.RealmIntegration:    8000,
            types.RealmMahayana:       10000,
            types.RealmTribulation:    15000,
        },
        agingPenalty50Pct: 0.0,
        agingPenalty20Pct: 0.09,
        agingPenalty10Pct: 0.14,
        agingPenaltyFinal: 0.20,
    }
}

// CalculateBaseLifespan returns the maximum lifespan for a given realm.
// Returns 0 for unknown realms.
func (l *LifespanRule) CalculateBaseLifespan(realm types.CultivationRealm) int {
    years, ok := l.baseLifespanByRealm[realm]
    if !ok {
        return 0
    }
    return years
}

// CalculateRemainingLifespan returns the remaining years of life.
// If age >= base lifespan, returns 0.
func (l *LifespanRule) CalculateRemainingLifespan(realm types.CultivationRealm, age int) int {
    base := l.CalculateBaseLifespan(realm)
    remaining := base - age
    if remaining <= 0 {
        return 0
    }
    return remaining
}

// CalculateRemainingLifespanRatio returns the ratio of remaining life to base lifespan (0.0 to 1.0).
// Returns 0.0 for unknown realms or if age >= base lifespan.
func (l *LifespanRule) CalculateRemainingLifespanRatio(realm types.CultivationRealm, age int) float64 {
    base := l.CalculateBaseLifespan(realm)
    if base <= 0 {
        return 0.0
    }
    remaining := base - age
    if remaining <= 0 {
        return 0.0
    }
    return float64(remaining) / float64(base)
}

// CalculateAgingPenalty computes the cultivation efficiency reduction based on remaining life ratio.
//
// Tiers:
//   - > 50% remaining: no penalty (returns 0.0)
//   - 20% - 50% remaining: linear interpolation from 0% to threshold_20pct penalty
//   - 10% - 20% remaining: linear interpolation from threshold_20pct to threshold_10pct
//   - < 10% remaining: threshold_final (20%)
func (l *LifespanRule) CalculateAgingPenalty(remainingRatio float64) float64 {
    if remainingRatio > 0.50 {
        return 0.0
    }
    if remainingRatio > 0.20 {
        t := (remainingRatio - 0.20) / (0.50 - 0.20)
        return l.agingPenalty20Pct * (1.0 - t)
    }
    if remainingRatio >= 0.10 {
        t := (remainingRatio - 0.10) / (0.20 - 0.10)
        return l.agingPenalty20Pct + (l.agingPenalty10Pct-l.agingPenalty20Pct)*(1.0-t)
    }
    return l.agingPenaltyFinal
}

// LifespanDepletionResult describes the outcome when an entity's lifespan is exhausted.
type LifespanDepletionResult struct {
    Depleted       bool
    ForcedAttempt  bool  // forced breakthrough attempt (if eligible)
    Death          bool  // entity dies (if not eligible)
    Message        string
}

// HandleLifespanDepletion determines the outcome when remaining lifespan reaches 0.
//
// Rules:
//   - If cultivationProgress >= realmBreakthroughThreshold, trigger a forced breakthrough attempt
//   - Otherwise, the entity dies (身死道消)
func (l *LifespanRule) HandleLifespanDepletion(cultivationProgress float64, breakthroughThreshold float64) LifespanDepletionResult {
    if cultivationProgress >= breakthroughThreshold {
        return LifespanDepletionResult{
            Depleted:      true,
            ForcedAttempt: true,
            Death:         false,
            Message:       "寿元将尽，触发九死一生突破尝试",
        }
    }
    return LifespanDepletionResult{
        Depleted:      true,
        ForcedAttempt: false,
        Death:         true,
        Message:       "寿元耗尽，身死道消",
    }
}

// CalculateCultivationPenaltyMultiplier returns the multiplier to apply to cultivation rate
// based on aging. Returns 1.0 (no penalty) when there is no aging effect.
func (l *LifespanRule) CalculateCultivationPenaltyMultiplier(remainingRatio float64) float64 {
    penalty := l.CalculateAgingPenalty(remainingRatio)
    return math.Max(0.0, 1.0-penalty)
}

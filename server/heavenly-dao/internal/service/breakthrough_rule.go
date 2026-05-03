package service

import (
    "math"

    "github.com/cultivation-world/shared/types"
)

// BreakthroughRule handles breakthrough success probability calculations.
type BreakthroughRule struct {
    baseSuccessRateByRealm map[types.CultivationRealm]float64
    maxAccumulationFactor  float64
    minProbability         float64
    maxProbability         float64
}

// NewBreakthroughRule creates a new BreakthroughRule with default configuration.
func NewBreakthroughRule() *BreakthroughRule {
    return &BreakthroughRule{
        baseSuccessRateByRealm: map[types.CultivationRealm]float64{
            types.RealmQiCondensation:    0.60,
            types.RealmFoundation:        0.40,
            types.RealmGoldenCore:        0.25,
            types.RealmNascentSoul:       0.15,
            types.RealmSoulTransform:     0.10,
            types.RealmVoidRefinement:    0.08,
            types.RealmIntegration:       0.05,
            types.RealmMahayana:          0.03,
            types.RealmTribulation:       0.01,
        },
        maxAccumulationFactor: 1.5,
        minProbability:        0.05,
        maxProbability:        0.80,
    }
}

// BreakthroughInput holds the inputs needed to calculate breakthrough success probability.
type BreakthroughInput struct {
    TargetRealm        types.CultivationRealm
    CultivationTime    float64 // time spent cultivating at current level
    RequiredTime       float64 // standard time required for breakthrough
    MethodQuality      float64 // quality of active cultivation method (0-100)
    ResourceBonus      float64 // sum of consumed items' breakthrough bonuses
    MentalStability    int     // mental stability (0-100)
    Luck               int     // luck attribute (0-100)
}

// CalculateBreakthroughSuccess computes the probability of a successful breakthrough.
//
// Formula:
//
//	prob = base_rate × accumulation × method_quality × resource × mental × luck
//
// Where:
//   - base_rate: realm-specific base success rate
//   - accumulation: min(maxAccumulationFactor, 1.0 + cultivation_time / required_time)
//   - method_quality: method_quality / 100.0
//   - resource: 1.0 + resource_bonus
//   - mental: mental_stability / 100.0
//   - luck: 1.0 + (luck - 50) / 200.0
//
// Result is clamped to [minProbability, maxProbability] (default [5%, 80%]).
func (b *BreakthroughRule) CalculateBreakthroughSuccess(input BreakthroughInput) float64 {
    baseRate, ok := b.baseSuccessRateByRealm[input.TargetRealm]
    if !ok {
        return b.minProbability
    }

    // Accumulation factor: min(max, 1.0 + time_ratio)
    var timeRatio float64
    if input.RequiredTime > 0 {
        timeRatio = input.CultivationTime / input.RequiredTime
    }
    accumulationFactor := 1.0 + timeRatio
    if accumulationFactor > b.maxAccumulationFactor {
        accumulationFactor = b.maxAccumulationFactor
    }

    // Method quality factor (0.0-1.0)
    methodQualityFactor := input.MethodQuality / 100.0
    methodQualityFactor = clamp01(methodQualityFactor)

    // Resource bonus factor
    resourceFactor := 1.0 + input.ResourceBonus

    // Mental stability factor (0.0-1.0)
    mentalFactor := float64(input.MentalStability) / 100.0
    mentalFactor = clamp01(mentalFactor)

    // Luck factor: ranges from 0.75 (luck=0) to 1.25 (luck=100)
    luckFactor := 1.0 + (float64(input.Luck)-50.0)/200.0

    prob := baseRate * accumulationFactor * methodQualityFactor * resourceFactor * mentalFactor * luckFactor

    // Clamp to [min, max]
    prob = math.Max(prob, b.minProbability)
    prob = math.Min(prob, b.maxProbability)

    return prob
}

// BreakthroughFailurePenalty describes the penalties when a breakthrough attempt fails.
type BreakthroughFailurePenalty struct {
    ProgressLoss    float64 // percentage of cultivation progress lost (default 10%)
    CooldownHours   float64 // cooldown before next attempt
    MentalDamage    int     // reduction in mental stability (default 20)
    InjuryProb      float64 // probability of injury (default 30%)
    InjurySeverity  string  // description of injury severity
}

// CalculateFailurePenalty computes the penalties for a failed breakthrough attempt.
//
// Penalties:
//   - Progress loss: 10% of current cultivation progress
//   - Cooldown: 24 × realm_level hours
//   - Mental damage: -20 mental stability
//   - Injury: 30% chance to halve current HP
func (b *BreakthroughRule) CalculateFailurePenalty(realmLevel int) BreakthroughFailurePenalty {
    cooldownHours := 24.0 * float64(realmLevel)

    return BreakthroughFailurePenalty{
        ProgressLoss:   0.10,
        CooldownHours:  cooldownHours,
        MentalDamage:   20,
        InjuryProb:     0.30,
        InjurySeverity: "half_hp",
    }
}

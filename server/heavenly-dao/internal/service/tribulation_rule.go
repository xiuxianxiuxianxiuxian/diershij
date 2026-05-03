package service

import (
    "math"

    "github.com/cultivation-world/shared/types"
)

// TribulationRule handles tribulation probability and strength calculations.
type TribulationRule struct {
    baseProbabilityByRealm map[types.CultivationRealm]float64
    strengthPerKarma       float64
    recentBreakBonus       float64
}

// NewTribulationRule creates a new TribulationRule with default configuration.
func NewTribulationRule() *TribulationRule {
    return &TribulationRule{
        baseProbabilityByRealm: map[types.CultivationRealm]float64{
            types.RealmQiCondensation:    0.10,
            types.RealmFoundation:        0.15,
            types.RealmGoldenCore:        0.20,
            types.RealmNascentSoul:       0.30,
            types.RealmSoulTransform:     0.50,
            types.RealmVoidRefinement:    0.60,
            types.RealmIntegration:       0.70,
            types.RealmMahayana:          0.80,
            types.RealmTribulation:       0.90,
        },
        strengthPerKarma: 1.002,
        recentBreakBonus: 0.1,
    }
}

// TribulationInput holds the inputs needed to calculate tribulation parameters.
type TribulationInput struct {
    TargetRealm            types.CultivationRealm
    Karma                  int
    Merit                  int
    Luck                   int
    RecentBreakthroughs    int
    MinProbability         float64
    MaxProbability         float64
    MeritReductionPer2000  float64
    MeritFloorFactor       float64
}

// DefaultTribulationInput returns a TribulationInput with reasonable defaults.
func DefaultTribulationInput() TribulationInput {
    return TribulationInput{
        MinProbability:        0.01,
        MaxProbability:        1.0,
        MeritReductionPer2000: 0.10,
        MeritFloorFactor:      0.30,
    }
}

// CalculateTribulationProbability computes the probability that a breakthrough attempt
// will trigger a heavenly tribulation.
//
// Formula:
//
//	prob = base_prob × karma_factor × merit_factor
//
// Where:
//   - base_prob: realm-specific base probability
//   - karma_factor: 1 + karma / 500
//   - merit_factor: max(floor, 1 - merit/2000 × reduction_per_2000)
//
// Result is clamped to [MinProbability, MaxProbability].
func (t *TribulationRule) CalculateTribulationProbability(input TribulationInput) float64 {
    baseProb, ok := t.baseProbabilityByRealm[input.TargetRealm]
    if !ok {
        return input.MinProbability
    }

    karmaFactor := 1.0 + float64(input.Karma)/500.0

    meritReductions := float64(input.Merit) / 2000.0
    meritFactor := 1.0 - meritReductions*input.MeritReductionPer2000
    meritFactor = math.Max(meritFactor, input.MeritFloorFactor)

    prob := baseProb * karmaFactor * meritFactor

    prob = math.Max(prob, input.MinProbability)
    prob = math.Min(prob, input.MaxProbability)

    return prob
}

// CalculateTribulationStrength computes the intensity of a heavenly tribulation.
//
// Formula:
//
//	strength = base_strength × (strength_per_karma ^ karma) × (1 + recent_breakthroughs × bonus)
func (t *TribulationRule) CalculateTribulationStrength(karma int, recentBreakthroughs int) float64 {
    baseStrength := 100.0

    karmaMultiplier := math.Pow(t.strengthPerKarma, float64(karma))

    recentBonus := 1.0 + float64(recentBreakthroughs)*t.recentBreakBonus

    return baseStrength * karmaMultiplier * recentBonus
}

// TribulationResult holds the full result of a tribulation assessment.
type TribulationResult struct {
    Triggered   bool
    Probability float64
    Strength    float64
}

// Assess evaluates whether a tribulation will trigger and its strength.
func (t *TribulationRule) Assess(input TribulationInput) TribulationResult {
    prob := t.CalculateTribulationProbability(input)
    strength := t.CalculateTribulationStrength(input.Karma, input.RecentBreakthroughs)

    return TribulationResult{
        Triggered:   prob >= 0.5,
        Probability: prob,
        Strength:    strength,
    }
}

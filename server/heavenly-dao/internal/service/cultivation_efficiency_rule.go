package service

import (
    "github.com/cultivation-world/shared/types"
)

// CultivationEfficiencyRule handles cultivation rate and method compatibility calculations.
type CultivationEfficiencyRule struct {
    baseRateMultiplier float64
    realmPenaltyFactor float64
}

// NewCultivationEfficiencyRule creates a new CultivationEfficiencyRule with default configuration.
func NewCultivationEfficiencyRule() *CultivationEfficiencyRule {
    return &CultivationEfficiencyRule{
        baseRateMultiplier: 0.1,
        realmPenaltyFactor: 0.2,
    }
}

// CultivationRateInput holds the inputs for calculating cultivation rate.
type CultivationRateInput struct {
    Comprehension   int
    SpiritualDensity int // location spiritual density (0-100 scale)
    MethodMatch     float64 // compatibility score (0.0-1.0)
    RealmLevel      int // numeric level of current realm (0-based)
    MentalState     int // mental stability (0-100)
    AgingPenalty    float64 // from LifespanRule.CalculateAgingPenalty
}

// CalculateCultivationRate computes the cultivation efficiency rate based on all factors.
//
// Formula:
//
//	rate = base_rate × spiritual_factor × method_match × realm_penalty × mental_factor × (1 - aging_penalty)
//
// Where:
//   - base_rate = comprehension × baseRateMultiplier (default 0.1)
//   - spiritual_factor = spiritual_density / 100.0
//   - method_match = pre-computed compatibility score (0.0-1.0)
//   - realm_penalty = 1.0 / (1.0 + realm_level × realmPenaltyFactor)
//   - mental_factor = mental_state / 100.0
//   - aging_penalty = from lifespan rule
func (c *CultivationEfficiencyRule) CalculateCultivationRate(input CultivationRateInput) float64 {
    baseRate := float64(input.Comprehension) * c.baseRateMultiplier

    spiritualFactor := float64(input.SpiritualDensity) / 100.0
    spiritualFactor = clamp01(spiritualFactor)

    realmPenalty := 1.0 / (1.0 + float64(input.RealmLevel)*c.realmPenaltyFactor)

    mentalFactor := float64(input.MentalState) / 100.0
    mentalFactor = clamp01(mentalFactor)

    agingFactor := 1.0 - input.AgingPenalty
    if agingFactor < 0.0 {
        agingFactor = 0.0
    }

    return baseRate * spiritualFactor * input.MethodMatch * realmPenalty * mentalFactor * agingFactor
}

// CalculateMethodCompatibility computes how well an entity's spiritual roots match a cultivation method.
//
// Returns a score from 0.0 (no match) to 1.0 (perfect match).
// Score = matching_root_count / max(required_root_count, 1)
func (c *CultivationEfficiencyRule) CalculateMethodCompatibility(roots []types.SpiritualRoot, requiredRoots []string) float64 {
    if len(requiredRoots) == 0 {
        return 1.0 // No requirements → fully compatible
    }

    rootSet := make(map[string]bool)
    for _, root := range roots {
        rootSet[root.Element] = true
    }

    matchCount := 0
    for _, required := range requiredRoots {
        if rootSet[required] {
            matchCount++
        }
    }

    return float64(matchCount) / float64(len(requiredRoots))
}

// GetRealmLevel returns the numeric level (0-based) for a realm, used in realm penalty calculation.
func (c *CultivationEfficiencyRule) GetRealmLevel(realm types.CultivationRealm) int {
    realmOrder := []types.CultivationRealm{
        types.RealmMortal,
        types.RealmQiCondensation,
        types.RealmFoundation,
        types.RealmGoldenCore,
        types.RealmNascentSoul,
        types.RealmSoulTransform,
        types.RealmVoidRefinement,
        types.RealmIntegration,
        types.RealmMahayana,
        types.RealmTribulation,
    }

    for i, r := range realmOrder {
        if r == realm {
            return i
        }
    }
    return 0
}

func clamp01(v float64) float64 {
    if v < 0.0 {
        return 0.0
    }
    if v > 1.0 {
        return 1.0
    }
    return v
}

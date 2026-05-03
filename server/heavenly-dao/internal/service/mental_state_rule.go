package service

import (
    "math"
)

// MentalStateRule handles mental stability, recovery, and related calculations.
type MentalStateRule struct {
    baseRecoveryRate   float64 // base mental stability recovery per day
    highStabilityThresh int    // threshold for full mental factor (default 80)
    lowStabilityThresh  int    // threshold for linear decay start (default 50)
}

// NewMentalStateRule creates a new MentalStateRule with default configuration.
func NewMentalStateRule() *MentalStateRule {
    return &MentalStateRule{
        baseRecoveryRate:    1.0,
        highStabilityThresh: 80,
        lowStabilityThresh:  50,
    }
}

// CalculateMentalFactor computes the mental state modifier for cultivation/breakthrough.
//
// Rules:
//   - stability >= 80: factor = 1.0 (no penalty)
//   - 50 <= stability < 80: linear from 1.0 to 0.0
//   - stability < 50: factor = 0.0 (severe penalty)
func (m *MentalStateRule) CalculateMentalFactor(mentalStability int) float64 {
    if mentalStability >= m.highStabilityThresh {
        return 1.0
    }
    if mentalStability <= m.lowStabilityThresh {
        return 0.0
    }
    // Linear interpolation from 0.0 at 50 to 1.0 at 80
    return float64(mentalStability-m.lowStabilityThresh) / float64(m.highStabilityThresh-m.lowStabilityThresh)
}

// CalculateMentalRecovery computes the mental stability recovery over a given period.
//
// Recovery rate:
//   - Base recovery per day (default 1.0)
//   - Multiplied by dao_heart factor: dao_heart / 100.0
//   - Reduced by obsession: -obsession_count * 0.5 per day
func (m *MentalStateRule) CalculateMentalRecovery(mentalStability int, daoHeart int, obsessionCount int, daysElapsed float64) float64 {
    if daysElapsed <= 0 {
        return 0.0
    }

    daoHeartFactor := float64(daoHeart) / 100.0
    daoHeartFactor = math.Max(0.1, math.Min(1.0, daoHeartFactor))

    dailyRate := m.baseRecoveryRate * daoHeartFactor
    obsessionPenalty := float64(obsessionCount) * 0.5
    dailyRate = math.Max(0.0, dailyRate-obsessionPenalty)

    return dailyRate * daysElapsed
}

// MentalStateChange represents the result of applying a mental state change.
type MentalStateChange struct {
    NewStability  int
    Changed       bool
    ThresholdCrossed bool // crossed into/out of high-stability zone
    Direction     string // "improved" or "degraded"
}

// ApplyMentalChange applies a change to mental stability with bounds checking.
func (m *MentalStateRule) ApplyMentalChange(currentStability int, change int) MentalStateChange {
    oldFactor := m.CalculateMentalFactor(currentStability)
    newStability := currentStability + change
    newStability = int(math.Max(0, math.Min(100, float64(newStability))))

    newFactor := m.CalculateMentalFactor(newStability)
    thresholdCrossed := (oldFactor >= 1.0 && newFactor < 1.0) || (oldFactor < 1.0 && newFactor >= 1.0)

    direction := "unchanged"
    if newStability > currentStability {
        direction = "improved"
    } else if newStability < currentStability {
        direction = "degraded"
    }

    return MentalStateChange{
        NewStability:     newStability,
        Changed:          change != 0,
        ThresholdCrossed: thresholdCrossed,
        Direction:        direction,
    }
}

// InnerDemonRule handles inner demon trigger probability, strength, and resolution.
type InnerDemonRule struct {
    baseProbability    float64
    qiDeviationProb    float64 // probability of qi_deviation on defeat
    victoryBonus       float64 // cultivation progress bonus on victory
    defeatProgressLoss float64 // cultivation progress loss on defeat
    defeatStabilityLoss int    // mental stability loss on defeat
}

// NewInnerDemonRule creates a new InnerDemonRule with default configuration.
func NewInnerDemonRule() *InnerDemonRule {
    return &InnerDemonRule{
        baseProbability:    0.001, // 0.1% per day
        qiDeviationProb:    0.30,
        victoryBonus:       0.10,
        defeatProgressLoss: 0.20,
        defeatStabilityLoss: 30,
    }
}

// InnerDemonInput holds inputs for inner demon probability calculation.
type InnerDemonInput struct {
    MentalStability           int
    Karma                     int
    RecentBreakthroughAttempt bool // attempted breakthrough in last 7 days
    MethodConflictScore       float64 // normalized conflict score (0-1)
}

// InnerDemonResult holds the result of an inner demon encounter.
type InnerDemonResult struct {
    Triggered     bool
    Strength      float64
    Outcome       string // "victory", "defeat", "qi_deviation", "not_triggered"
    Effect        string
}

// CalculateInnerDemonProbability computes the daily probability of inner demon invasion.
//
// Formula:
//
//	prob = base_prob * mental_factor * karma_factor * breakthrough_stress * conflict_factor
//
// Where:
//   - base_prob = 0.001 (0.1% per day)
//   - mental_factor = (100 - mental_stability) / 100.0
//   - karma_factor = 1.0 + karma / 1000.0
//   - breakthrough_stress = 2.0 if recent attempt, else 1.0
//   - conflict_factor = 1.0 + conflict_score * 2.0
func (i *InnerDemonRule) CalculateInnerDemonProbability(input InnerDemonInput) float64 {
    mentalFactor := float64(100-input.MentalStability) / 100.0
    mentalFactor = math.Max(0.0, math.Min(1.0, mentalFactor))

    karmaFactor := 1.0 + float64(input.Karma)/1000.0
    karmaFactor = math.Max(1.0, karmaFactor)

    breakthroughStress := 1.0
    if input.RecentBreakthroughAttempt {
        breakthroughStress = 2.0
    }

    conflictFactor := 1.0 + input.MethodConflictScore*2.0

    return i.baseProbability * mentalFactor * karmaFactor * breakthroughStress * conflictFactor
}

// CalculateInnerDemonStrength computes the intensity of an inner demon.
//
// Formula:
//
//	strength = realm_level * 100 + obsession_count * 50
func (i *InnerDemonRule) CalculateInnerDemonStrength(realmLevel int, obsessionCount int) float64 {
    baseStrength := float64(realmLevel) * 100.0
    obsessionBonus := float64(obsessionCount) * 50.0
    return baseStrength + obsessionBonus
}

// ResolveInnerDemon determines the outcome of an inner demon encounter.
//
// Uses a probability-based resolution:
//   - victory_prob = inner_demon_resistance / 100.0
//   - If victory: mental_stability +20, cultivation_progress +10%
//   - If defeat: cultivation_progress -20%, mental_stability -30
//   - 30% chance of qi_deviation on defeat (realm drop, permanent -10% attributes)
func (i *InnerDemonRule) ResolveInnerDemon(innerDemonResistance int, randFloat func() float64) InnerDemonResult {
    victoryProb := float64(innerDemonResistance) / 100.0
    victoryProb = math.Max(0.0, math.Min(1.0, victoryProb))

    if randFloat() < victoryProb {
        return InnerDemonResult{
            Triggered: true,
            Outcome:   "victory",
            Effect:    "enlightenment",
        }
    }

    // Defeat - check for qi_deviation
    if randFloat() < i.qiDeviationProb {
        return InnerDemonResult{
            Triggered: true,
            Outcome:   "qi_deviation",
            Effect:    "walk_fire_entrance",
        }
    }

    return InnerDemonResult{
        Triggered: true,
        Outcome:   "defeat",
        Effect:    "cultivation_regression",
    }
}

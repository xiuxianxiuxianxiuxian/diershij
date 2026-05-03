package service

import (
    "math"
)

// MethodConflictRule handles cultivation method conflict detection and backlash calculation.
type MethodConflictRule struct {
    oppositeElementConflict float64
    alignmentConflict       float64
    realmGapConflict        float64
}

// NewMethodConflictRule creates a new MethodConflictRule with default configuration.
func NewMethodConflictRule() *MethodConflictRule {
    return &MethodConflictRule{
        oppositeElementConflict: 0.5,
        alignmentConflict:       0.3,
        realmGapConflict:        0.2,
    }
}

// MethodInfo represents a cultivation method for conflict analysis.
type MethodInfo struct {
    ID            string
    Element       string // fire, water, wood, metal, earth
    Alignment     string // righteous, demonic, neutral
    RealmRequired int    // minimum realm level required
}

// elementOpposites defines which elements are mutually opposing.
var elementOpposites = map[string]string{
    "fire":  "water",
    "water": "fire",
    "wood":  "metal",
    "metal": "wood",
    "earth": "wood",
}

// isOppositeElements returns true if two elements are opposing.
func isOppositeElements(a, b string) bool {
    if a == "" || b == "" {
        return false
    }
    return elementOpposites[a] == b || elementOpposites[b] == a
}

// ConflictPair represents a conflict between two methods with its score.
type ConflictPair struct {
    MethodA     string
    MethodB     string
    Score       float64
    ConflictType string // element, alignment, realm, multiple
}

// CalculateMethodConflict computes the conflict score between two methods.
//
// Conflict types:
//   - Element opposite: +0.5
//   - Alignment conflict (righteous vs demonic): +0.3
//   - Realm requirement gap > 2 levels: +0.2
//
// Returns the conflict score for this pair.
func (m *MethodConflictRule) CalculateMethodConflict(a, b MethodInfo) ConflictPair {
    var score float64
    var types []string

    // Element opposition
    if isOppositeElements(a.Element, b.Element) {
        score += m.oppositeElementConflict
        types = append(types, "element")
    }

    // Alignment conflict: righteous vs demonic
    if (a.Alignment == "righteous" && b.Alignment == "demonic") ||
        (a.Alignment == "demonic" && b.Alignment == "righteous") {
        score += m.alignmentConflict
        types = append(types, "alignment")
    }

    // Realm gap > 2 levels
    realmDiff := int(math.Abs(float64(a.RealmRequired - b.RealmRequired)))
    if realmDiff > 2 {
        score += m.realmGapConflict
        types = append(types, "realm")
    }

    conflictType := "none"
    if len(types) == 1 {
        conflictType = types[0]
    } else if len(types) > 1 {
        conflictType = "multiple"
    }

    return ConflictPair{
        MethodA:     a.ID,
        MethodB:     b.ID,
        Score:       score,
        ConflictType: conflictType,
    }
}

// OverallConflictResult holds the overall conflict assessment for multiple methods.
type OverallConflictResult struct {
    TotalScore     float64
    NormalizedScore float64 // score / max_possible_score (0-1)
    BacklashProb   float64
    Pairs          []ConflictPair
    Severity       string // none, minor, moderate, severe
}

// CalculateOverallConflict computes the overall conflict score for a set of methods.
//
// Formula:
//
//	conflict_score = sum(conflict_scores) / (n * (n-1) / 2)
//
// Where n is the number of methods. Returns normalized score (0-1).
func (m *MethodConflictRule) CalculateOverallConflict(methods []MethodInfo) OverallConflictResult {
    n := len(methods)
    if n <= 1 {
        return OverallConflictResult{
            NormalizedScore: 0.0,
            BacklashProb:    0.0,
            Severity:        "none",
        }
    }

    maxPairs := n * (n - 1) / 2
    var pairs []ConflictPair
    totalScore := 0.0

    for i := 0; i < n; i++ {
        for j := i + 1; j < n; j++ {
            conflict := m.CalculateMethodConflict(methods[i], methods[j])
            pairs = append(pairs, conflict)
            totalScore += conflict.Score
        }
    }

    normalizedScore := totalScore / float64(maxPairs)

    backlashProb := m.calculateBacklashProbability(normalizedScore)
    severity := m.determineSeverity(normalizedScore, backlashProb)

    return OverallConflictResult{
        TotalScore:     totalScore,
        NormalizedScore: normalizedScore,
        BacklashProb:   backlashProb,
        Pairs:          pairs,
        Severity:       severity,
    }
}

func (m *MethodConflictRule) calculateBacklashProbability(normalizedScore float64) float64 {
    if normalizedScore < 0.2 {
        return 0.0
    }
    if normalizedScore < 0.5 {
        return 0.1
    }
    if normalizedScore < 0.8 {
        return 0.3
    }
    return 0.6
}

func (m *MethodConflictRule) determineSeverity(normalizedScore, backlashProb float64) string {
    if backlashProb == 0.0 {
        return "none"
    }
    if backlashProb < 0.15 {
        return "minor"
    }
    if backlashProb < 0.4 {
        return "moderate"
    }
    return "severe"
}

// CalculateCultivationPenalty returns the cultivation rate reduction based on conflict severity.
//
// Penalty tiers:
//   - < 0.2 conflict: 0% penalty
//   - 0.2-0.5 conflict: -20% rate
//   - 0.5-0.8 conflict: -50% rate
//   - > 0.8 conflict: -80% rate
func (m *MethodConflictRule) CalculateCultivationPenalty(normalizedScore float64) float64 {
    if normalizedScore < 0.2 {
        return 0.0
    }
    if normalizedScore < 0.5 {
        return 0.20
    }
    if normalizedScore < 0.8 {
        return 0.50
    }
    return 0.80
}

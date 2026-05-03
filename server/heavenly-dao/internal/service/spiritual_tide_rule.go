package service

import (
	"math"
)

// SpiritualTideRule handles spiritual tide (灵气潮汐) calculations.
type SpiritualTideRule struct {
	baseDensity       float64
	tideAmplitude     float64 // max deviation from base density
	tideCycleDays     int     // full tide cycle in days
}

// NewSpiritualTideRule creates a new SpiritualTideRule with default configuration.
func NewSpiritualTideRule() *SpiritualTideRule {
	return &SpiritualTideRule{
		baseDensity:    1.0,
		tideAmplitude:  0.30, // ±30% swing
		tideCycleDays:  30,   // 30-day cycle
	}
}

// TidePhase represents the current phase of the spiritual tide.
type TidePhase string

const (
	TidePhaseLow       TidePhase = "low"       // 灵气低谷
	TidePhaseRising    TidePhase = "rising"    // 灵气上升
	TidePhaseHigh      TidePhase = "high"      // 灵气高峰
	TidePhaseFalling   TidePhase = "falling"   // 灵气下降
)

// TideResult holds the current spiritual tide state.
type TideResult struct {
	Phase          TidePhase
	Density        float64 // current spiritual density multiplier
	DaysUntilShift int     // days until next phase transition
}

// CalculateCurrentTide computes the current spiritual tide state.
//
// Uses a sinusoidal model:
//
//	density = base + amplitude * sin(2π * day / cycle)
func (r *SpiritualTideRule) CalculateCurrentTide(currentDay int, regionBonus float64) TideResult {
	// Sinusoidal tide model
	cycle := float64(r.tideCycleDays)
	angle := 2.0 * math.Pi * float64(currentDay) / cycle
	tideOffset := r.tideAmplitude * math.Sin(angle)

	// Base density with region bonus
	baseDensity := r.baseDensity * math.Max(0.1, regionBonus)

	// Current density
	density := baseDensity * (1.0 + tideOffset)
	density = math.Max(0.1, density)

	// Determine phase
	phase := r.determinePhase(angle)

	// Days until next phase shift (every π/2 radians)
	currentPhaseAngle := math.Mod(angle, 2*math.Pi)
	if currentPhaseAngle < 0 {
		currentPhaseAngle += 2 * math.Pi
	}
	nextShiftAngle := math.Ceil(currentPhaseAngle/(math.Pi/2.0)) * (math.Pi / 2.0)
	if nextShiftAngle <= currentPhaseAngle {
		nextShiftAngle += math.Pi / 2.0
	}
	daysUntilShift := int(math.Ceil((nextShiftAngle - currentPhaseAngle) / (2.0 * math.Pi) * cycle))
	if daysUntilShift <= 0 {
		daysUntilShift = 1
	}

	return TideResult{
		Phase:          phase,
		Density:        density,
		DaysUntilShift: daysUntilShift,
	}
}

func (r *SpiritualTideRule) determinePhase(angle float64) TidePhase {
	// Normalize angle to [0, 2π)
	a := math.Mod(angle, 2*math.Pi)
	if a < 0 {
		a += 2 * math.Pi
	}

	switch {
	case a < math.Pi/2:
		return TidePhaseRising
	case a < math.Pi:
		return TidePhaseHigh
	case a < 3*math.Pi/2:
		return TidePhaseFalling
	default:
		return TidePhaseLow
	}
}

// AdjustSpiritualDensity applies temporary adjustments to spiritual density from events.
func (r *SpiritualTideRule) AdjustSpiritualDensity(baseDensity float64, adjustments []float64) float64 {
	total := baseDensity
	for _, adj := range adjustments {
		total *= (1.0 + adj)
	}
	return math.Max(0.01, total)
}

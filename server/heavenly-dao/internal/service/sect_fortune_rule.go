package service

import (
	"math"
)

// SectFortuneRule handles sect fortune (宗门气运) calculation and trajectory prediction.
type SectFortuneRule struct {
	baseFortune     float64
	fortuneDecay    float64 // daily decay rate
}

// NewSectFortuneRule creates a new SectFortuneRule with default configuration.
func NewSectFortuneRule() *SectFortuneRule {
	return &SectFortuneRule{
		baseFortune:  50.0,
		fortuneDecay: 0.001, // 0.1% daily decay
	}
}

// SectState holds the current state of a sect.
type SectState struct {
	MemberCount        int     // total members
	EliteCount         int     // members at GoldenCore+ realm
	LeaderRealm        int     // sect leader's realm level (0-9)
	Treasury           float64 // spirit stone treasury
	TerritoryCount     int     // number of controlled territories
	AllianceCount      int     // number of active alliances
	EnemyCount         int     // number of active enemies/conflicts
	RecentBreakthrough int     // breakthroughs in last 30 days
	RecentDeaths       int     // member deaths in last 30 days
}

// CalculateSectFortune computes the current fortune score of a sect.
//
// Formula combines multiple factors:
//
//	fortune = base + members + elites*3 + leader_realm*5 + territory*2
//	         + alliance*2 - enemy*3 + breakthrough*1 - death*2 + log(treasury)*0.5
func (r *SectFortuneRule) CalculateSectFortune(state SectState) float64 {
	fortune := r.baseFortune

	// Member contribution
	fortune += float64(state.MemberCount) * 0.1

	// Elite contribution (weighted more)
	fortune += float64(state.EliteCount) * 3.0

	// Leader realm contribution
	fortune += float64(state.LeaderRealm) * 5.0

	// Territory contribution
	fortune += float64(state.TerritoryCount) * 2.0

	// Alliance bonus
	fortune += float64(state.AllianceCount) * 2.0

	// Enemy penalty
	fortune -= float64(state.EnemyCount) * 3.0

	// Recent events
	fortune += float64(state.RecentBreakthrough) * 1.0
	fortune -= float64(state.RecentDeaths) * 2.0

	// Treasury (diminishing returns)
	if state.Treasury > 0 {
		fortune += math.Log10(state.Treasury) * 0.5
	}

	return math.Max(0, fortune)
}

// SectFortuneGrade represents the overall fortune grade of a sect.
type SectFortuneGrade string

const (
	SectGradeDeclining   SectFortuneGrade = "declining"   // 衰落
	SectGradeStable      SectFortuneGrade = "stable"      // 稳定
	SectGradeProspering  SectFortuneGrade = "prospering"  // 兴盛
	SectGradeDominant    SectFortuneGrade = "dominant"    // 霸主
)

// GetSectFortuneGrade maps a fortune score to a grade.
func (r *SectFortuneRule) GetSectFortuneGrade(fortune float64) SectFortuneGrade {
	switch {
	case fortune < 30:
		return SectGradeDeclining
	case fortune < 60:
		return SectGradeStable
	case fortune < 100:
		return SectGradeProspering
	default:
		return SectGradeDominant
	}
}

// SectTrajectory represents the predicted trajectory of a sect.
type SectTrajectory struct {
	CurrentFortune float64
	PredictedFortune float64
	Trend          string // "rising", "stable", "falling"
	DaysToNextGrade int   // estimated days until grade change
}

// PredictSectTrajectory predicts the sect's fortune trajectory over the next N days.
func (r *SectFortuneRule) PredictSectTrajectory(state SectState, days int) SectTrajectory {
	currentFortune := r.CalculateSectFortune(state)

	// Simple projection: apply decay and expected changes
	projectedState := state
	projectedState.MemberCount = int(float64(state.MemberCount) * math.Pow(0.999, float64(days)))
	projectedState.RecentBreakthrough = int(float64(state.RecentBreakthrough) * math.Pow(0.95, float64(days)/7.0))
	projectedState.RecentDeaths = int(float64(state.RecentDeaths) * math.Pow(0.95, float64(days)/7.0))

	predictedFortune := r.CalculateSectFortune(projectedState)

	// Apply fortune decay
	predictedFortune *= (1.0 - r.fortuneDecay*float64(days))

	// Determine trend
	change := predictedFortune - currentFortune
	var trend string
	switch {
	case change > 1.0:
		trend = "rising"
	case change < -1.0:
		trend = "falling"
	default:
		trend = "stable"
	}

	// Estimate days to next grade change
	currentGrade := r.GetSectFortuneGrade(currentFortune)
	daysToNext := r.estimateDaysToGradeChange(currentFortune, currentGrade, change, days)

	return SectTrajectory{
		CurrentFortune:   currentFortune,
		PredictedFortune: predictedFortune,
		Trend:            trend,
		DaysToNextGrade:  daysToNext,
	}
}

func (r *SectFortuneRule) estimateDaysToGradeChange(fortune float64, grade SectFortuneGrade, change float64, days int) int {
	if days <= 0 {
		return -1
	}

	dailyRate := change / float64(days)
	if math.Abs(dailyRate) < 0.01 {
		return -1 // too slow to predict
	}

	// Distance to next threshold
	var threshold float64
	switch grade {
	case SectGradeDeclining:
		threshold = 30
	case SectGradeStable:
		threshold = 60
	case SectGradeProspering:
		threshold = 100
	default:
		return -1 // already at max
	}

	if dailyRate > 0 {
		distance := threshold - fortune
		if distance > 0 {
			return int(math.Ceil(distance / dailyRate))
		}
	}

	return -1
}

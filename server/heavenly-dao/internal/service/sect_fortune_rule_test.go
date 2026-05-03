package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSectFortuneRule(t *testing.T) {
	rule := NewSectFortuneRule()
	assert.NotNil(t, rule)
	assert.Equal(t, 50.0, rule.baseFortune)
	assert.Equal(t, 0.001, rule.fortuneDecay)
}

func TestCalculateSectFortune_BaseSect(t *testing.T) {
	rule := NewSectFortuneRule()

	state := SectState{
		MemberCount:        100,
		EliteCount:         5,
		LeaderRealm:        3, // GoldenCore
		Treasury:           10000,
		TerritoryCount:     2,
		AllianceCount:      1,
		EnemyCount:         0,
		RecentBreakthrough: 2,
		RecentDeaths:       0,
	}

	fortune := rule.CalculateSectFortune(state)
	// base = 50
	// members = 100 * 0.1 = 10
	// elites = 5 * 3 = 15
	// leader = 3 * 5 = 15
	// territory = 2 * 2 = 4
	// alliance = 1 * 2 = 2
	// breakthrough = 2 * 1 = 2
	// treasury = log10(10000) * 0.5 = 4 * 0.5 = 2
	// total = 50 + 10 + 15 + 15 + 4 + 2 + 2 + 2 = 100
	assert.InDelta(t, 100.0, fortune, 0.01)
}

func TestCalculateSectFortune_EnemySect(t *testing.T) {
	rule := NewSectFortuneRule()

	state := SectState{
		MemberCount:        50,
		EliteCount:         2,
		LeaderRealm:        2,
		Treasury:           1000,
		TerritoryCount:     1,
		AllianceCount:      0,
		EnemyCount:         3,
		RecentBreakthrough: 0,
		RecentDeaths:       4,
	}

	fortune := rule.CalculateSectFortune(state)
	// base = 50
	// members = 50 * 0.1 = 5
	// elites = 2 * 3 = 6
	// leader = 2 * 5 = 10
	// territory = 1 * 2 = 2
	// enemy = 3 * 3 = -9
	// deaths = 4 * 2 = -8
	// treasury = log10(1000) * 0.5 = 3 * 0.5 = 1.5
	// total = 50 + 5 + 6 + 10 + 2 - 9 - 8 + 1.5 = 57.5
	assert.InDelta(t, 57.5, fortune, 0.01)
}

func TestCalculateSectFortune_MinClamp(t *testing.T) {
	rule := NewSectFortuneRule()

	state := SectState{
		MemberCount:        0,
		EliteCount:         0,
		LeaderRealm:        0,
		Treasury:           0,
		TerritoryCount:     0,
		AllianceCount:      0,
		EnemyCount:         100,
		RecentBreakthrough: 0,
		RecentDeaths:       50,
	}

	fortune := rule.CalculateSectFortune(state)
	// fortune will be negative → clamped to 0
	assert.InDelta(t, 0.0, fortune, 0.01)
}

func TestGetSectFortuneGrade_AllGrades(t *testing.T) {
	rule := NewSectFortuneRule()

	assert.Equal(t, SectGradeDeclining, rule.GetSectFortuneGrade(10))
	assert.Equal(t, SectGradeDeclining, rule.GetSectFortuneGrade(29.9))
	assert.Equal(t, SectGradeStable, rule.GetSectFortuneGrade(30))
	assert.Equal(t, SectGradeStable, rule.GetSectFortuneGrade(59.9))
	assert.Equal(t, SectGradeProspering, rule.GetSectFortuneGrade(60))
	assert.Equal(t, SectGradeProspering, rule.GetSectFortuneGrade(99.9))
	assert.Equal(t, SectGradeDominant, rule.GetSectFortuneGrade(100))
	assert.Equal(t, SectGradeDominant, rule.GetSectFortuneGrade(200))
}

func TestPredictSectTrajectory_Stable(t *testing.T) {
	rule := NewSectFortuneRule()

	state := SectState{
		MemberCount:        100,
		EliteCount:         5,
		LeaderRealm:        3,
		Treasury:           10000,
		TerritoryCount:     2,
		AllianceCount:      1,
		EnemyCount:         0,
		RecentBreakthrough: 2,
		RecentDeaths:       2,
	}

	trajectory := rule.PredictSectTrajectory(state, 30)
	assert.True(t, trajectory.CurrentFortune > 0)
	assert.True(t, trajectory.PredictedFortune > 0)
	assert.NotEmpty(t, trajectory.Trend)
}

func TestPredictSectTrajectory_Falling(t *testing.T) {
	rule := NewSectFortuneRule()

	// Sect with many allies and territories that decay over time
	state := SectState{
		MemberCount:        200,
		EliteCount:         10,
		LeaderRealm:        4,
		Treasury:           50000,
		TerritoryCount:     5,
		AllianceCount:      3,
		EnemyCount:         1,
		RecentBreakthrough: 0,
		RecentDeaths:       0,
	}

	trajectory := rule.PredictSectTrajectory(state, 365)
	// Over a year, members and alliances decay, causing decline
	assert.True(t, trajectory.PredictedFortune < trajectory.CurrentFortune)
	assert.Equal(t, "falling", trajectory.Trend)
}

func TestEstimateDaysToGradeChange_Rising(t *testing.T) {
	rule := NewSectFortuneRule()

	// At stable grade (55), trying to reach 60 (prospering)
	days := rule.estimateDaysToGradeChange(55, SectGradeStable, 5.0, 10)
	// Need 5 more points, rate = 0.5/day → 10 days
	assert.True(t, days > 0)
}

func TestEstimateDaysToGradeChange_TooSlow(t *testing.T) {
	rule := NewSectFortuneRule()

	days := rule.estimateDaysToGradeChange(55, SectGradeStable, 0.05, 10)
	// daily rate = 0.005 < 0.01 → returns -1
	assert.Equal(t, -1, days)
}

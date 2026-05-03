package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewEpochService(t *testing.T) {
	svc := NewEpochService()
	assert.NotNil(t, svc)
}

func TestDefaultEpochTransitionCriteria(t *testing.T) {
	criteria := DefaultEpochTransitionCriteria()
	assert.NotNil(t, criteria)
	assert.Equal(t, 365, criteria.MinimumDurationDays)
	assert.InDelta(t, 0.15, criteria.SpiritualThreshold, 0.01)
	assert.Equal(t, 50, criteria.EventCountThreshold)
	assert.InDelta(t, 0.85, criteria.BalanceThreshold, 0.01)
}

func TestStartNewEpoch_ValidInput(t *testing.T) {
	svc := NewEpochService()

	input := EpochInput{
		Phase:            EpochSpiritualAge,
		Description:      "灵气充沛的时代",
		SpiritualBase:    0.8,
		ResourceCap:      150,
		EventFrequency:   0.5,
		SecretRealmCount: 5,
	}

	now := time.Now()
	epoch, err := svc.StartNewEpoch(input, now)
	assert.NoError(t, err)
	assert.Equal(t, 1, epoch.EpochNumber)
	assert.Equal(t, EpochSpiritualAge, epoch.Phase)
	assert.Equal(t, 0.8, epoch.SpiritualBase)
	assert.Equal(t, now, epoch.StartTime)
}

func TestStartNewEpoch_InvalidSpiritualBase(t *testing.T) {
	svc := NewEpochService()

	input := EpochInput{
		Phase:         EpochSpiritualAge,
		SpiritualBase: 1.5,
	}

	_, err := svc.StartNewEpoch(input, time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "spiritual base")
}

func TestStartNewEpoch_SecondEpochEndsFirst(t *testing.T) {
	svc := NewEpochService()

	now := time.Now()
	svc.StartNewEpoch(EpochInput{Phase: EpochSpiritualAge}, now)

	// Start second epoch
	nextYear := now.Add(366 * 24 * time.Hour)
	epoch2, _ := svc.StartNewEpoch(EpochInput{Phase: EpochDharmaDecline}, nextYear)

	// First epoch should be ended and in history
	history := svc.GetEpochHistory()
	assert.Len(t, history, 1)
	assert.Equal(t, nextYear, history[0].EndTime)
	assert.Equal(t, epoch2.EpochNumber, 2)
}

func TestCheckTransition_NoEpoch(t *testing.T) {
	svc := NewEpochService()

	result := svc.CheckTransition(0.5, 0, 0.3, time.Now())
	assert.True(t, result.ShouldTransition)
	assert.Equal(t, EpochSpiritualAge, result.SuggestedPhase)
}

func TestCheckTransition_DurationNotMet(t *testing.T) {
	svc := NewEpochService()

	now := time.Now()
	svc.StartNewEpoch(EpochInput{Phase: EpochSpiritualAge}, now)

	// Only 100 days passed
	later := now.Add(100 * 24 * time.Hour)
	result := svc.CheckTransition(0.5, 0, 0.3, later)
	assert.False(t, result.ShouldTransition)
	assert.Contains(t, result.Reason, "持续时间不足")
}

func TestCheckTransition_SpiritualCollapse(t *testing.T) {
	svc := NewEpochService()

	now := time.Now()
	svc.StartNewEpoch(EpochInput{Phase: EpochSpiritualAge}, now)

	// After 1 year with very low spiritual density
	later := now.Add(366 * 24 * time.Hour)
	result := svc.CheckTransition(0.1, 0, 0.3, later) // 0.1 < 0.15 threshold
	assert.True(t, result.ShouldTransition)
	assert.Equal(t, EpochDharmaDecline, result.SuggestedPhase)
	assert.Contains(t, result.Reason, "灵气密度过低")
}

func TestCheckTransition_EventOverload(t *testing.T) {
	svc := NewEpochService()

	now := time.Now()
	svc.StartNewEpoch(EpochInput{Phase: EpochSpiritualAge}, now)

	later := now.Add(366 * 24 * time.Hour)
	result := svc.CheckTransition(0.5, 60, 0.3, later) // 60 >= 50 threshold
	assert.True(t, result.ShouldTransition)
	assert.Equal(t, EpochChaos, result.SuggestedPhase)
	assert.Contains(t, result.Reason, "事件过多")
}

func TestCheckTransition_SevereImbalance(t *testing.T) {
	svc := NewEpochService()

	now := time.Now()
	svc.StartNewEpoch(EpochInput{Phase: EpochSpiritualAge}, now)

	later := now.Add(366 * 24 * time.Hour)
	result := svc.CheckTransition(0.5, 0, 0.9, later) // 0.9 >= 0.85 threshold
	assert.True(t, result.ShouldTransition)
	assert.Equal(t, EpochDharmaDecline, result.SuggestedPhase)
	assert.Contains(t, result.Reason, "严重失衡")
}

func TestCheckTransition_RecoveryFromDecline(t *testing.T) {
	svc := NewEpochService()

	now := time.Now()
	svc.StartNewEpoch(EpochInput{Phase: EpochDharmaDecline}, now)

	later := now.Add(366 * 24 * time.Hour)
	result := svc.CheckTransition(0.6, 0, 0.3, later) // 0.6 > 0.5
	assert.True(t, result.ShouldTransition)
	assert.Equal(t, EpochRenewal, result.SuggestedPhase)
	assert.Contains(t, result.Reason, "灵气复苏")
}

func TestCheckTransition_StableEpoch(t *testing.T) {
	svc := NewEpochService()

	now := time.Now()
	svc.StartNewEpoch(EpochInput{Phase: EpochSpiritualAge}, now)

	later := now.Add(366 * 24 * time.Hour)
	result := svc.CheckTransition(0.5, 10, 0.3, later)
	assert.False(t, result.ShouldTransition)
	assert.Contains(t, result.Reason, "稳定")
}

func TestTransitionToNewEpoch_SpiritualAge(t *testing.T) {
	svc := NewEpochService()

	now := time.Now()
	epoch, err := svc.TransitionToNewEpoch(EpochSpiritualAge, now)
	assert.NoError(t, err)
	assert.Equal(t, EpochSpiritualAge, epoch.Phase)
	assert.Equal(t, 0.8, epoch.SpiritualBase)
	assert.Equal(t, 150, epoch.ResourceCap)
	assert.Equal(t, 5, epoch.SecretRealmCount)
}

func TestTransitionToNewEpoch_DharmaDecline(t *testing.T) {
	svc := NewEpochService()

	now := time.Now()
	epoch, err := svc.TransitionToNewEpoch(EpochDharmaDecline, now)
	assert.NoError(t, err)
	assert.Equal(t, EpochDharmaDecline, epoch.Phase)
	assert.Equal(t, 0.2, epoch.SpiritualBase)
	assert.Equal(t, 50, epoch.ResourceCap)
}

func TestTransitionToNewEpoch_Chaos(t *testing.T) {
	svc := NewEpochService()

	now := time.Now()
	epoch, err := svc.TransitionToNewEpoch(EpochChaos, now)
	assert.NoError(t, err)
	assert.Equal(t, EpochChaos, epoch.Phase)
	assert.Equal(t, 2.0, epoch.EventFrequency)
	assert.Equal(t, 3, epoch.SecretRealmCount)
}

func TestTransitionToNewEpoch_Renewal(t *testing.T) {
	svc := NewEpochService()

	now := time.Now()
	epoch, err := svc.TransitionToNewEpoch(EpochRenewal, now)
	assert.NoError(t, err)
	assert.Equal(t, EpochRenewal, epoch.Phase)
	assert.Equal(t, 0.6, epoch.SpiritualBase)
	assert.Equal(t, 4, epoch.SecretRealmCount)
}

func TestTransitionToNewEpoch_UnknownPhase(t *testing.T) {
	svc := NewEpochService()

	_, err := svc.TransitionToNewEpoch(EpochPhase("unknown"), time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "未知")
}

func TestGetCurrentEpoch(t *testing.T) {
	svc := NewEpochService()

	// No epoch yet
	_, ok := svc.GetCurrentEpoch()
	assert.False(t, ok)

	now := time.Now()
	svc.StartNewEpoch(EpochInput{Phase: EpochSpiritualAge}, now)

	epoch, ok := svc.GetCurrentEpoch()
	assert.True(t, ok)
	assert.Equal(t, EpochSpiritualAge, epoch.Phase)
}

func TestGetEpochHistory(t *testing.T) {
	svc := NewEpochService()

	now := time.Now()
	svc.StartNewEpoch(EpochInput{Phase: EpochSpiritualAge}, now)
	svc.StartNewEpoch(EpochInput{Phase: EpochDharmaDecline}, now.Add(366*24*time.Hour))
	svc.StartNewEpoch(EpochInput{Phase: EpochChaos}, now.Add(732*24*time.Hour))

	history := svc.GetEpochHistory()
	assert.Len(t, history, 2) // 2 past epochs
}

func TestGetEpochDuration(t *testing.T) {
	svc := NewEpochService()

	_, err := svc.GetEpochDuration(time.Now())
	assert.Error(t, err)

	now := time.Now()
	svc.StartNewEpoch(EpochInput{Phase: EpochSpiritualAge}, now)

	duration, err := svc.GetEpochDuration(now.Add(100 * 24 * time.Hour))
	assert.NoError(t, err)
	assert.InDelta(t, 100*24*float64(time.Hour), float64(duration), float64(time.Hour))
}

func TestGetCurrentPhase(t *testing.T) {
	svc := NewEpochService()

	_, err := svc.GetCurrentPhase()
	assert.Error(t, err)

	now := time.Now()
	svc.StartNewEpoch(EpochInput{Phase: EpochDharmaDecline}, now)

	phase, err := svc.GetCurrentPhase()
	assert.NoError(t, err)
	assert.Equal(t, EpochDharmaDecline, phase)
}

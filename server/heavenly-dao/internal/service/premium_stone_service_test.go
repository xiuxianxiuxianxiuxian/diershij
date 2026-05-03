package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPremiumStoneService(t *testing.T) {
	svc := NewPremiumStoneService()
	assert.NotNil(t, svc)
}

func TestGrant_ValidSource(t *testing.T) {
	svc := NewPremiumStoneService()

	input := GrantInput{
		EntityID:    "entity_1",
		Source:      SourceCreateMethod,
		Amount:      5,
		GrantID:     "method_001",
		Description: "功法自创奖励",
	}

	now := time.Now()
	result, err := svc.Grant(input, now)
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, int64(5), result.AmountGranted)
	assert.Contains(t, result.Reason, "5")
}

func TestGrant_InvalidSource(t *testing.T) {
	svc := NewPremiumStoneService()

	input := GrantInput{
		EntityID: "entity_1",
		Source:   PremiumStoneSource("invalid_source"),
		Amount:   5,
	}

	_, err := svc.Grant(input, time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无效")
}

func TestGrant_ZeroAmount(t *testing.T) {
	svc := NewPremiumStoneService()

	input := GrantInput{
		EntityID: "entity_1",
		Source:   SourceCreateMethod,
		Amount:   0,
	}

	_, err := svc.Grant(input, time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "大于0")
}

func TestGrant_DuplicatePrevention(t *testing.T) {
	svc := NewPremiumStoneService()

	input := GrantInput{
		EntityID: "entity_1",
		Source:   SourceCreateMethod,
		Amount:   5,
		GrantID:  "duplicate_test_001",
	}

	now := time.Now()
	result1, err := svc.Grant(input, now)
	assert.NoError(t, err)
	assert.True(t, result1.Success)

	// Same grant ID should fail
	result2, err := svc.Grant(input, now)
	assert.NoError(t, err)
	assert.False(t, result2.Success)
	assert.Contains(t, result2.Reason, "重复")
}

func TestGrant_DuplicatePrevention_DifferentIDs(t *testing.T) {
	svc := NewPremiumStoneService()

	input1 := GrantInput{
		EntityID: "entity_1",
		Source:   SourceCreateMethod,
		Amount:   5,
		GrantID:  "grant_001",
	}

	input2 := GrantInput{
		EntityID: "entity_1",
		Source:   SourceCreateMethod,
		Amount:   5,
		GrantID:  "grant_002",
	}

	now := time.Now()
	result1, _ := svc.Grant(input1, now)
	result2, _ := svc.Grant(input2, now)
	assert.True(t, result1.Success)
	assert.True(t, result2.Success)
}

func TestGrant_AnnualCap(t *testing.T) {
	svc := NewPremiumStoneService()

	now := time.Now()
	// Grant up to the cap
	input := GrantInput{
		EntityID: "entity_cap",
		Source:   SourceCreateMethod,
		Amount:   annualPremiumStoneCap - 10,
		GrantID:  "cap_test_1",
	}
	result, _ := svc.Grant(input, now)
	assert.True(t, result.Success)
	assert.Equal(t, int64(annualPremiumStoneCap-10), result.AmountGranted)

	// Grant remaining
	input2 := GrantInput{
		EntityID: "entity_cap",
		Source:   SourceWorldBoss,
		Amount:   20,
		GrantID:  "cap_test_2",
	}
	result2, _ := svc.Grant(input2, now)
	assert.True(t, result2.Success)
	assert.Equal(t, int64(10), result2.AmountGranted) // only 10 remaining
}

func TestGrant_AnnualCapFull(t *testing.T) {
	svc := NewPremiumStoneService()

	now := time.Now()
	// Fill the cap
	input := GrantInput{
		EntityID: "entity_full",
		Source:   SourceCreateMethod,
		Amount:   annualPremiumStoneCap,
		GrantID:  "full_test_1",
	}
	result, _ := svc.Grant(input, now)
	assert.True(t, result.Success)

	// Try to grant more
	input2 := GrantInput{
		EntityID: "entity_full",
		Source:   SourceWorldBoss,
		Amount:   5,
		GrantID:  "full_test_2",
	}
	result2, _ := svc.Grant(input2, now)
	assert.False(t, result2.Success)
	assert.Contains(t, result2.Reason, "上限已满")
}

func TestGrant_NewYearResets(t *testing.T) {
	svc := NewPremiumStoneService()

	// Grant in year 2025
	input := GrantInput{
		EntityID: "entity_year",
		Source:   SourceCreateMethod,
		Amount:   annualPremiumStoneCap,
		GrantID:  "year_test_1",
	}
	now := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)
	result1, _ := svc.Grant(input, now)
	assert.True(t, result1.Success)

	// Grant in year 2026 should work (new year)
	input2 := GrantInput{
		EntityID: "entity_year",
		Source:   SourceCreateMethod,
		Amount:   5,
		GrantID:  "year_test_2",
	}
	nextYear := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	result2, _ := svc.Grant(input2, nextYear)
	assert.True(t, result2.Success)
	assert.Equal(t, int64(5), result2.AmountGranted)
}

func TestGetAnnualSummary(t *testing.T) {
	svc := NewPremiumStoneService()

	now := time.Now()
	svc.Grant(GrantInput{EntityID: "entity_sum", Source: SourceCreateMethod, Amount: 5, GrantID: "sum_1"}, now)
	svc.Grant(GrantInput{EntityID: "entity_sum", Source: SourceWorldBoss, Amount: 10, GrantID: "sum_2"}, now)

	summary := svc.GetAnnualSummary("entity_sum", now.Year())
	assert.Equal(t, int64(5), summary[SourceCreateMethod])
	assert.Equal(t, int64(10), summary[SourceWorldBoss])
}

func TestGetRemainingAnnualCap(t *testing.T) {
	svc := NewPremiumStoneService()

	now := time.Now()
	svc.Grant(GrantInput{EntityID: "entity_rem", Source: SourceCreateMethod, Amount: 100, GrantID: "rem_1"}, now)

	remaining := svc.GetRemainingAnnualCap("entity_rem", now)
	assert.Equal(t, int64(annualPremiumStoneCap-100), remaining)
}

func TestGetValidSources(t *testing.T) {
	sources := GetValidSources()
	assert.True(t, len(sources) >= 7) // at least 7 sources
}

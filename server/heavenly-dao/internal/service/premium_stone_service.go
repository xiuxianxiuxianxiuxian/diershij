package service

import (
	"fmt"
	"time"
)

const (
	annualPremiumStoneCap = 1000 // max premium stones produced per year
)

// PremiumStoneSource represents a valid source of premium spirit stones.
type PremiumStoneSource string

const (
	SourceCreateMethod   PremiumStoneSource = "create_method"    // 功法自创奖励
	SourceWorldBoss      PremiumStoneSource = "world_boss"       // 世界Boss
	SourceSecretRealm    PremiumStoneSource = "secret_realm"     // 秘境探索
	SourceTribulation    PremiumStoneSource = "tribulation"      // 渡劫成功
	SourceSectConquest   PremiumStoneSource = "sect_conquest"    // 宗门征战
	SourceAnnualDividend PremiumStoneSource = "annual_dividend"  // 年度分红
	SourceAchievement    PremiumStoneSource = "achievement"      // 成就奖励
)

// validPremiumSources is the set of all allowed premium stone sources.
var validPremiumSources = map[PremiumStoneSource]bool{
	SourceCreateMethod:   true,
	SourceWorldBoss:      true,
	SourceSecretRealm:    true,
	SourceTribulation:    true,
	SourceSectConquest:   true,
	SourceAnnualDividend: true,
	SourceAchievement:    true,
}

// PremiumStoneService manages premium spirit stone distribution with strict controls.
type PremiumStoneService struct {
	// entityID -> source -> amount granted this year
	annualGrants map[string]map[PremiumStoneSource]int64
	// entityID -> year -> total granted
	annualTotals map[string]map[int]int64
	// entityID -> source -> set of unique grant IDs (for duplicate prevention)
	grantRecords map[string]map[PremiumStoneSource]map[string]bool
}

// NewPremiumStoneService creates a new PremiumStoneService.
func NewPremiumStoneService() *PremiumStoneService {
	return &PremiumStoneService{
		annualGrants: make(map[string]map[PremiumStoneSource]int64),
		annualTotals: make(map[string]map[int]int64),
		grantRecords: make(map[string]map[PremiumStoneSource]map[string]bool),
	}
}

// GrantInput holds inputs for granting premium spirit stones.
type GrantInput struct {
	EntityID    string
	Source      PremiumStoneSource
	Amount      int64
	GrantID     string // unique identifier for duplicate prevention
	Description string
}

// GrantResult holds the outcome of a grant operation.
type GrantResult struct {
	Success      bool
	AmountGranted int64
	Reason       string
}

// Grant attempts to grant premium spirit stones to an entity.
func (s *PremiumStoneService) Grant(input GrantInput, now time.Time) (*GrantResult, error) {
	// Validate source
	if !validPremiumSources[input.Source] {
		return nil, fmt.Errorf("无效的极品灵石来源: %s", input.Source)
	}

	// Validate amount
	if input.Amount <= 0 {
		return nil, fmt.Errorf("发放数量必须大于0")
	}

	// Check for duplicate grant
	if s.isDuplicateGrant(input.EntityID, input.Source, input.GrantID) {
		return &GrantResult{
			Success: false,
			Reason:  fmt.Sprintf("重复发放: %s 已从 %s 领取", input.GrantID, input.Source),
		}, nil
	}

	// Check annual cap
	currentYear := now.Year()
	currentTotal := s.getAnnualTotal(input.EntityID, currentYear)
	if currentTotal+input.Amount > annualPremiumStoneCap {
		remaining := annualPremiumStoneCap - currentTotal
		if remaining <= 0 {
			return &GrantResult{
				Success: false,
				Reason:  fmt.Sprintf("年度上限已满 (%d/%d)", currentTotal, annualPremiumStoneCap),
			}, nil
		}
		// Grant only remaining amount
		input.Amount = remaining
	}

	// Record the grant
	s.recordGrant(input, currentYear)

	return &GrantResult{
		Success:      true,
		AmountGranted: input.Amount,
		Reason:        fmt.Sprintf("成功发放 %d 极品灵石 (%s)", input.Amount, input.Source),
	}, nil
}

func (s *PremiumStoneService) isDuplicateGrant(entityID string, source PremiumStoneSource, grantID string) bool {
	if grantID == "" {
		return false // no ID means no duplicate check
	}

	if sources, exists := s.grantRecords[entityID]; exists {
		if grants, exists := sources[source]; exists {
			return grants[grantID]
		}
	}
	return false
}

func (s *PremiumStoneService) getAnnualTotal(entityID string, year int) int64 {
	if years, exists := s.annualTotals[entityID]; exists {
		return years[year]
	}
	return 0
}

func (s *PremiumStoneService) recordGrant(input GrantInput, year int) {
	// Initialize maps if needed
	if _, exists := s.annualGrants[input.EntityID]; !exists {
		s.annualGrants[input.EntityID] = make(map[PremiumStoneSource]int64)
	}
	if _, exists := s.annualTotals[input.EntityID]; !exists {
		s.annualTotals[input.EntityID] = make(map[int]int64)
	}
	if _, exists := s.grantRecords[input.EntityID]; !exists {
		s.grantRecords[input.EntityID] = make(map[PremiumStoneSource]map[string]bool)
	}

	// Record the grant
	s.annualGrants[input.EntityID][input.Source] += input.Amount
	s.annualTotals[input.EntityID][year] += input.Amount

	if _, exists := s.grantRecords[input.EntityID][input.Source]; !exists {
		s.grantRecords[input.EntityID][input.Source] = make(map[string]bool)
	}
	if input.GrantID != "" {
		s.grantRecords[input.EntityID][input.Source][input.GrantID] = true
	}
}

// GetAnnualSummary returns the annual premium stone summary for an entity.
func (s *PremiumStoneService) GetAnnualSummary(entityID string, year int) map[PremiumStoneSource]int64 {
	if sources, exists := s.annualGrants[entityID]; exists {
		result := make(map[PremiumStoneSource]int64)
		for src, amt := range sources {
			result[src] = amt
		}
		return result
	}
	return make(map[PremiumStoneSource]int64)
}

// GetRemainingAnnualCap returns how many more premium stones an entity can receive this year.
func (s *PremiumStoneService) GetRemainingAnnualCap(entityID string, now time.Time) int64 {
	total := s.getAnnualTotal(entityID, now.Year())
	remaining := annualPremiumStoneCap - total
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetValidSources returns all valid premium stone sources.
func GetValidSources() []PremiumStoneSource {
	sources := make([]PremiumStoneSource, 0, len(validPremiumSources))
	for src := range validPremiumSources {
		sources = append(sources, src)
	}
	return sources
}

package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cultivation-world/shared/types"
)

// OperationService is the unified service that manages all game operations.
type OperationService struct {
	combat        *CombatOperation
	explore       *ExploreOperation
	gather        *GatherOperation
	craft         *CraftOperation
	createMethod  *CreateMethodOperation
	trade         *TradeOperation
	sect          *SectOperation
	spell         *SpellOperation

	breakthrough  *BreakthroughRule
	cultivation   *CultivationEfficiencyRule

	defaultCooldown time.Duration
}

// NewOperationService creates a new OperationService with all sub-operations.
func NewOperationService(defaultCooldown time.Duration) *OperationService {
	return &OperationService{
		combat:        NewCombatOperation(defaultCooldown * 2),
		explore:       NewExploreOperation(defaultCooldown * 5),
		gather:        NewGatherOperation(defaultCooldown),
		craft:         NewCraftOperation(defaultCooldown * 3),
		createMethod:  NewCreateMethodOperation(defaultCooldown * 100),
		trade:         NewTradeOperation(defaultCooldown),
		sect:          NewSectOperation(defaultCooldown * 50),
		spell:         NewSpellOperation(),
		breakthrough:  NewBreakthroughRule(),
		cultivation:   NewCultivationEfficiencyRule(),
		defaultCooldown: defaultCooldown,
	}
}

// CultivateInput holds inputs for cultivation.
type CultivateInput struct {
	EntityID         string
	Realm            types.CultivationRealm
	SpiritualRoots   []types.SpiritualRoot
	MainMethod       *types.CultivationMethod
	SpiritualDensity float64 // 0.0-1.0
	Comprehension    int
	MentalStability  int
	BaseLifespan     int
	CurrentAge       int
}

// CultivateResult holds the outcome of cultivation.
type CultivateResult struct {
	Success           bool
	CultivationGained float64
	Rate              float64
	RealmBonus        float64
	SpiritualMult     float64
	MethodMatch       float64
	MentalFactor      float64
	AgingPenalty      float64
	Message           string
}

// ExecuteCultivate performs one cultivation cycle using the cultivation efficiency algorithm.
func (s *OperationService) ExecuteCultivate(input CultivateInput, now time.Time, randFloat func() float64) (*CultivateResult, error) {
	// Calculate method compatibility
	methodMatch := 1.0
	if input.MainMethod != nil {
		requiredRoots := []string{}
		if input.MainMethod.RequiredRoots != nil {
			requiredRoots = input.MainMethod.RequiredRoots
		}
		methodMatch = s.cultivation.CalculateMethodCompatibility(input.SpiritualRoots, requiredRoots)
	}

	agingPenalty := s.calculateAgingPenalty(input.BaseLifespan, input.CurrentAge)
	realmLevel := types.CultivationRealmLevel(input.Realm)
	if realmLevel < 0 {
		realmLevel = 0
	}

	// Convert spiritual density to 0-100 scale
	spiritualDensity100 := int(input.SpiritualDensity * 100)

	// Calculate cultivation rate
	rate := s.cultivation.CalculateCultivationRate(CultivationRateInput{
		Comprehension:    input.Comprehension,
		SpiritualDensity: spiritualDensity100,
		MethodMatch:      methodMatch,
		RealmLevel:       realmLevel,
		MentalState:      input.MentalStability,
		AgingPenalty:     agingPenalty,
	})

	// Apply randomness (+/- 10%)
	variance := 0.9 + randFloat()*0.2
	effectiveRate := rate * variance

	return &CultivateResult{
		Success:           true,
		CultivationGained: effectiveRate,
		Rate:              rate,
		RealmBonus:        getRealmMultiplier(input.Realm),
		SpiritualMult:     getSpiritualMult(input.SpiritualRoots, input.SpiritualDensity),
		MethodMatch:       methodMatch,
		MentalFactor:      getMentalFactor(input.MentalStability),
		AgingPenalty:      agingPenalty,
		Message:           fmt.Sprintf("修炼完成，获得 %.2f 修为", effectiveRate),
	}, nil
}

// OpBreakthroughInput holds inputs for breakthrough.
type OpBreakthroughInput struct {
	EntityID        string
	CurrentRealm    types.CultivationRealm
	TargetRealm     types.CultivationRealm
	CultivationTime float64
	RequiredTime    float64
	MethodQuality   float64
	ResourceBonus   float64
	MentalStability int
	Luck            int
}

// OpBreakthroughResult holds the outcome of breakthrough.
type OpBreakthroughResult struct {
	Success     bool
	NewRealm    types.CultivationRealm
	SuccessRate float64
	Tribulation *TribulationResult
	Penalty     *BreakthroughFailurePenalty
	Message     string
}

// ExecuteBreakthrough attempts a realm breakthrough using the breakthrough algorithm.
func (s *OperationService) ExecuteBreakthrough(input OpBreakthroughInput, now time.Time, randFloat func() float64) (*OpBreakthroughResult, error) {
	tribInput := DefaultTribulationInput()
	tribInput.TargetRealm = input.TargetRealm
	tribInput.Luck = input.Luck

	// Calculate breakthrough success rate
	successRate := s.breakthrough.CalculateBreakthroughSuccess(BreakthroughInput{
		TargetRealm:     input.TargetRealm,
		CultivationTime: input.CultivationTime,
		RequiredTime:    input.RequiredTime,
		MethodQuality:   input.MethodQuality,
		ResourceBonus:   input.ResourceBonus,
		MentalStability: input.MentalStability,
		Luck:            input.Luck,
	})

	// Check if breakthrough succeeds
	if randFloat() >= successRate {
		realmLevel := types.CultivationRealmLevel(input.CurrentRealm)
		if realmLevel < 0 {
			realmLevel = 0
		}
		penalty := s.breakthrough.CalculateFailurePenalty(realmLevel)

		return &OpBreakthroughResult{
			Success:     false,
			NewRealm:    input.CurrentRealm,
			SuccessRate: successRate,
			Penalty:     &penalty,
			Message:     fmt.Sprintf("突破失败！成功率 %.1f%%", successRate*100),
		}, nil
	}

	// Success - check for tribulation at high realms
	var tribResult *TribulationResult
	targetLevel := types.CultivationRealmLevel(input.TargetRealm)
	if targetLevel >= 4 { // Nascent Soul and above
		tribRule := NewTribulationRule()
		assessment := tribRule.Assess(tribInput)
		tribResult = &assessment
	}

	return &OpBreakthroughResult{
		Success:     true,
		NewRealm:    input.TargetRealm,
		SuccessRate: successRate,
		Tribulation: tribResult,
		Message:     fmt.Sprintf("突破成功！进入 %s", input.TargetRealm),
	}, nil
}

func (s *OperationService) calculateAgingPenalty(baseLifespan, currentAge int) float64 {
	if baseLifespan <= 0 {
		return 0
	}
	remaining := baseLifespan - currentAge
	if remaining <= 0 {
		return 0.20
	}
	remainingRatio := float64(remaining) / float64(baseLifespan)
	if remainingRatio > 0.50 {
		return 0
	}
	if remainingRatio > 0.20 {
		return (0.50 - remainingRatio) / 0.30 * 0.09
	}
	if remainingRatio > 0.10 {
		return 0.09 + (0.20-remainingRatio)/0.10*0.05
	}
	return 0.20
}

func getRealmMultiplier(realm types.CultivationRealm) float64 {
	multipliers := map[types.CultivationRealm]float64{
		types.RealmMortal:         0.0,
		types.RealmQiCondensation: 1.0,
		types.RealmFoundation:     1.5,
		types.RealmGoldenCore:     2.0,
		types.RealmNascentSoul:    2.5,
		types.RealmSoulTransform:  3.0,
		types.RealmVoidRefinement: 3.5,
		types.RealmIntegration:    4.0,
		types.RealmMahayana:       4.5,
		types.RealmTribulation:    5.0,
	}
	if m, ok := multipliers[realm]; ok {
		return m
	}
	return 1.0
}

func getSpiritualMult(roots []types.SpiritualRoot, density float64) float64 {
	if len(roots) == 0 {
		return 0.1
	}
	return density
}

func getMentalFactor(stability int) float64 {
	if stability >= 80 {
		return 1.0
	}
	if stability >= 50 {
		return float64(stability-50) / 30.0
	}
	return 0
}

// Accessors for sub-operations
func (s *OperationService) Combat() *CombatOperation            { return s.combat }
func (s *OperationService) Explore() *ExploreOperation          { return s.explore }
func (s *OperationService) Gather() *GatherOperation            { return s.gather }
func (s *OperationService) Craft() *CraftOperation              { return s.craft }
func (s *OperationService) CreateMethod() *CreateMethodOperation { return s.createMethod }
func (s *OperationService) Trade() *TradeOperation              { return s.trade }
func (s *OperationService) Sect() *SectOperation                { return s.sect }
func (s *OperationService) Spell() *SpellOperation              { return s.spell }

// DefaultRand returns a standard random function.
func DefaultRand() func() float64 {
	return rand.Float64
}

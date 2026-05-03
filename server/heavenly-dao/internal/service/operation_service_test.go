package service

import (
	"testing"
	"time"

	"github.com/cultivation-world/shared/types"
	"github.com/stretchr/testify/assert"
)

func TestNewOperationService(t *testing.T) {
	svc := NewOperationService(time.Minute)
	assert.NotNil(t, svc)
	assert.NotNil(t, svc.combat)
	assert.NotNil(t, svc.explore)
	assert.NotNil(t, svc.gather)
	assert.NotNil(t, svc.craft)
	assert.NotNil(t, svc.createMethod)
	assert.NotNil(t, svc.trade)
	assert.NotNil(t, svc.sect)
	assert.NotNil(t, svc.spell)
	assert.NotNil(t, svc.breakthrough)
	assert.NotNil(t, svc.cultivation)
}

func TestExecuteCultivate_Basic(t *testing.T) {
	svc := NewOperationService(time.Minute)

	input := CultivateInput{
		EntityID:         "cultivator_1",
		Realm:            types.RealmQiCondensation,
		SpiritualRoots:   []types.SpiritualRoot{{Element: "fire", Purity: 80}},
		MainMethod: &types.CultivationMethod{
			ElementAffinity: "fire",
			CultivationSpeedMult: 1.5,
		},
		SpiritualDensity: 0.6,
		Comprehension:    70,
		MentalStability:  85,
		BaseLifespan:     100,
		CurrentAge:       20,
	}

	now := time.Now()
	result, err := svc.ExecuteCultivate(input, now, func() float64 { return 0.5 })
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.True(t, result.CultivationGained > 0)
	assert.Contains(t, result.Message, "修炼完成")
}

func TestExecuteCultivate_PoorConditions(t *testing.T) {
	svc := NewOperationService(time.Minute)

	input := CultivateInput{
		EntityID:         "cultivator_2",
		Realm:            types.RealmQiCondensation,
		SpiritualRoots:   []types.SpiritualRoot{{Element: "fire", Purity: 30}},
		SpiritualDensity: 0.1,
		Comprehension:    20,
		MentalStability:  30,
		BaseLifespan:     100,
		CurrentAge:       95,
	}

	now := time.Now()
	result, err := svc.ExecuteCultivate(input, now, func() float64 { return 0.5 })
	assert.NoError(t, err)
	assert.True(t, result.Success)
	// Should still work but with very low rate
	assert.True(t, result.CultivationGained >= 0)
}

func TestExecuteBreakthrough_Success(t *testing.T) {
	svc := NewOperationService(time.Minute)

	input := OpBreakthroughInput{
		EntityID:        "breaker_1",
		CurrentRealm:    types.RealmQiCondensation,
		TargetRealm:     types.RealmFoundation,
		CultivationTime: 100,
		RequiredTime:    100,
		MethodQuality:   80,
		ResourceBonus:   0.2,
		MentalStability: 90,
		Luck:            80,
	}

	now := time.Now()
	result, err := svc.ExecuteBreakthrough(input, now, func() float64 { return 0.1 })
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, types.RealmFoundation, result.NewRealm)
	assert.Contains(t, result.Message, "突破成功")
}

func TestExecuteBreakthrough_Failure(t *testing.T) {
	svc := NewOperationService(time.Minute)

	input := OpBreakthroughInput{
		EntityID:        "breaker_3",
		CurrentRealm:    types.RealmQiCondensation,
		TargetRealm:     types.RealmFoundation,
		CultivationTime: 100,
		RequiredTime:    100,
		MethodQuality:   30,
		ResourceBonus:   0,
		MentalStability: 20,
		Luck:            0,
	}

	now := time.Now()
	result, err := svc.ExecuteBreakthrough(input, now, func() float64 { return 0.99 })
	assert.NoError(t, err)
	assert.False(t, result.Success)
	assert.NotNil(t, result.Penalty)
	assert.Contains(t, result.Message, "失败")
}

func TestExecuteBreakthrough_WithTribulation(t *testing.T) {
	svc := NewOperationService(time.Minute)

	input := OpBreakthroughInput{
		EntityID:        "breaker_4",
		CurrentRealm:    types.RealmGoldenCore,
		TargetRealm:     types.RealmNascentSoul, // triggers tribulation
		CultivationTime: 200,
		RequiredTime:    100,
		MethodQuality:   90,
		ResourceBonus:   0.5,
		MentalStability: 90,
		Luck:            80,
	}

	now := time.Now()
	result, err := svc.ExecuteBreakthrough(input, now, func() float64 { return 0.1 })
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.NotNil(t, result.Tribulation)
}

func TestOperationService_Accessors(t *testing.T) {
	svc := NewOperationService(time.Minute)

	assert.NotNil(t, svc.Combat())
	assert.NotNil(t, svc.Explore())
	assert.NotNil(t, svc.Gather())
	assert.NotNil(t, svc.Craft())
	assert.NotNil(t, svc.CreateMethod())
	assert.NotNil(t, svc.Trade())
	assert.NotNil(t, svc.Sect())
	assert.NotNil(t, svc.Spell())
}

func TestHelperFunctions(t *testing.T) {
	// Test realm multipliers
	assert.Equal(t, 0.0, getRealmMultiplier(types.RealmMortal))
	assert.Equal(t, 1.0, getRealmMultiplier(types.RealmQiCondensation))
	assert.Equal(t, 5.0, getRealmMultiplier(types.RealmTribulation))

	// Test mental factor
	assert.Equal(t, 1.0, getMentalFactor(100))
	assert.Equal(t, 1.0, getMentalFactor(80))
	assert.True(t, getMentalFactor(65) > 0 && getMentalFactor(65) < 1.0)
	assert.Equal(t, 0.0, getMentalFactor(50))
	assert.Equal(t, 0.0, getMentalFactor(20))

	// Test spiritual mult
	assert.Equal(t, 0.1, getSpiritualMult(nil, 0.5))
	assert.Equal(t, 0.6, getSpiritualMult([]types.SpiritualRoot{{}}, 0.6))
}

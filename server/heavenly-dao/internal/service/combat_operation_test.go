package service

import (
	"testing"
	"time"

	"github.com/cultivation-world/shared/types"
	"github.com/stretchr/testify/assert"
)

func makeCombatInput() CombatOpInput {
	return CombatOpInput{
		AttackerID:        "attacker_1",
		DefenderID:        "defender_1",
		AttackerName:      "张三",
		DefenderName:      "李四",
		AttackerRealm:     types.RealmFoundation,
		DefenderRealm:     types.RealmFoundation,
		AttackerHP:        100,
		DefenderHP:        100,
		AttackerMaxHP:     100,
		DefenderMaxHP:     100,
		AttackerAttack:    50,
		DefenderAttack:    40,
		AttackerDefense:   30,
		DefenderDefense:   25,
		AttackerElement:   "fire",
		DefenderElement:   "wood",
		AttackerPenetrate: 5,
		DefenderPenetrate: 3,
		AttackerCritRate:  0.1,
		DefenderCritRate:  0.05,
		AttackerDodgeRate: 0.05,
		DefenderDodgeRate: 0.1,
		IsSelfDefense:     false,
	}
}

func TestNewCombatOperation(t *testing.T) {
	op := NewCombatOperation(time.Minute)
	assert.NotNil(t, op)
}

func TestExecuteCombat_AttackerWins(t *testing.T) {
	op := NewCombatOperation(time.Minute)

	input := makeCombatInput()
	// Attacker has higher stats, should win
	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	result, err := op.ExecuteCombat(input, now, deterministic)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.WinnerID)
	assert.True(t, result.Rounds > 0)
	assert.True(t, result.TotalDamage > 0)
	assert.NotEmpty(t, result.Message)
}

func TestExecuteCombat_Cooldown(t *testing.T) {
	op := NewCombatOperation(time.Hour)

	input := makeCombatInput()
	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	// First combat succeeds
	result, err := op.ExecuteCombat(input, now, deterministic)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Second combat fails due to cooldown
	_, err = op.ExecuteCombat(input, now.Add(30*time.Minute), deterministic)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cooldown")
}

func TestExecuteCombat_CooldownExpired(t *testing.T) {
	op := NewCombatOperation(time.Minute)

	input := makeCombatInput()
	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	result1, err := op.ExecuteCombat(input, now, deterministic)
	assert.NoError(t, err)
	assert.NotNil(t, result1)

	// After cooldown expires, combat should succeed
	result2, err := op.ExecuteCombat(input, now.Add(2*time.Minute), deterministic)
	assert.NoError(t, err)
	assert.NotNil(t, result2)
}

func TestExecuteCombat_SelfDefense(t *testing.T) {
	op := NewCombatOperation(time.Minute)

	input := makeCombatInput()
	input.IsSelfDefense = true

	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	result, err := op.ExecuteCombat(input, now, deterministic)
	assert.NoError(t, err)
	assert.Equal(t, 0, result.KarmaChange)
}

func TestExecuteCombat_LootGeneration(t *testing.T) {
	op := NewCombatOperation(time.Minute)

	input := makeCombatInput()
	now := time.Now()
	// Deterministic loot generation: first call is for combat, then for loot
	callCount := 0
	deterministic := func() float64 {
		callCount++
		if callCount <= 100 {
			return 0.5 // combat
		}
		return 0.01 // always generate loot
	}

	result, err := op.ExecuteCombat(input, now, deterministic)
	assert.NoError(t, err)

	// If attacker wins, there should be loot
	if result.WinnerID == input.AttackerID {
		assert.NotEmpty(t, result.LootItems)
		// First item is always spirit stones
		assert.Equal(t, "spirit_stone", result.LootItems[0].Type)
	}
}

func TestExecuteCombat_Injuries(t *testing.T) {
	op := NewCombatOperation(time.Minute)

	input := makeCombatInput()
	// Make attacker very strong to cause significant damage
	input.AttackerAttack = 500
	now := time.Now()
	deterministic := func() float64 { return 0.5 }

	result, err := op.ExecuteCombat(input, now, deterministic)
	assert.NoError(t, err)

	// If defender loses and took significant damage, should have injuries
	if result.LoserID == input.DefenderID && len(result.Injuries) > 0 {
		assert.True(t, result.Injuries[0].Severity > 0)
		assert.True(t, result.Injuries[0].HealTime > 0)
	}
}

func TestCalculateCombatKarma_Bully(t *testing.T) {
	op := NewCombatOperation(time.Minute)

	input := CombatOpInput{
		AttackerRealm: types.RealmNascentSoul,    // level 4
		DefenderRealm: types.RealmQiCondensation,  // level 1, diff = 3 > 2
		IsSelfDefense: false,
	}

	karma := op.calculateCombatKarma(input, true)
	assert.Equal(t, -20, karma) // bullying penalty
}

func TestCalculateCombatKarma_FairFight(t *testing.T) {
	op := NewCombatOperation(time.Minute)

	input := CombatOpInput{
		AttackerRealm: types.RealmFoundation,
		DefenderRealm: types.RealmFoundation,
		IsSelfDefense: false,
	}

	karma := op.calculateCombatKarma(input, true)
	assert.Equal(t, -10, karma) // normal attack penalty
}

func TestCalculateCombatKarma_AttackerLoses(t *testing.T) {
	op := NewCombatOperation(time.Minute)

	input := CombatOpInput{
		AttackerRealm: types.RealmFoundation,
		DefenderRealm: types.RealmFoundation,
		IsSelfDefense: false,
	}

	karma := op.calculateCombatKarma(input, false)
	assert.Equal(t, -5, karma) // picked wrong fight
}

func TestGetCooldownRemaining(t *testing.T) {
	op := NewCombatOperation(time.Hour)

	// No cooldown set
	remaining := op.GetCooldownRemaining("test_id", time.Now())
	assert.Equal(t, time.Duration(0), remaining)

	// Set cooldown via combat
	input := makeCombatInput()
	now := time.Now()
	_, _ = op.ExecuteCombat(input, now, func() float64 { return 0.5 })

	remaining = op.GetCooldownRemaining("attacker_1", now.Add(30*time.Minute))
	assert.True(t, remaining > 0)

	// After cooldown expired
	remaining = op.GetCooldownRemaining("attacker_1", now.Add(2*time.Hour))
	assert.Equal(t, time.Duration(0), remaining)
}

func TestClearCooldown(t *testing.T) {
	op := NewCombatOperation(time.Hour)

	input := makeCombatInput()
	now := time.Now()
	_, _ = op.ExecuteCombat(input, now, func() float64 { return 0.5 })

	remaining := op.GetCooldownRemaining("attacker_1", now.Add(30*time.Minute))
	assert.True(t, remaining > 0)

	op.ClearCooldown("attacker_1")
	remaining = op.GetCooldownRemaining("attacker_1", now.Add(30*time.Minute))
	assert.Equal(t, time.Duration(0), remaining)
}

func TestGenerateLoot_ScalesWithRealm(t *testing.T) {
	op := NewCombatOperation(time.Minute)

	randAlways05 := func() float64 { return 0.5 }

	// Low realm loot
	lowLoot := op.generateLoot(types.RealmQiCondensation, randAlways05)
	assert.True(t, lowLoot[0].Quantity > 0)

	// High realm loot
	highLoot := op.generateLoot(types.RealmNascentSoul, randAlways05)
	assert.True(t, highLoot[0].Quantity > lowLoot[0].Quantity)
}

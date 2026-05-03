package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSpellOperation(t *testing.T) {
	op := NewSpellOperation()
	assert.NotNil(t, op)
}

func TestExecuteCastSpell_Success(t *testing.T) {
	op := NewSpellOperation()

	input := CastSpellInput{
		CasterID:      "caster_1",
		TargetID:      "target_1",
		SkillID:       "fireball",
		SkillName:     "火球术",
		SkillDamage:   50,
		SkillCooldown: 5,
		CasterAttack:  100,
		CasterElement: "fire",
		TargetElement: "wood",
		CasterRealm:   2,
		TargetRealm:   1,
	}

	now := time.Now()
	result, err := op.ExecuteCastSpell(input, now, func() float64 { return 0.5 })
	assert.NoError(t, err)
	assert.True(t, result.Success)
	assert.True(t, result.Damage > 0)
	assert.Contains(t, result.Message, "施展了")
}

func TestExecuteCastSpell_GlobalCooldown(t *testing.T) {
	op := NewSpellOperation()

	input := CastSpellInput{
		CasterID:     "caster_1",
		SkillID:      "spell1",
		SkillName:    "spell",
		SkillDamage:  50,
		CasterAttack: 100,
	}

	now := time.Now()
	_, err := op.ExecuteCastSpell(input, now, func() float64 { return 0.5 })
	assert.NoError(t, err)

	// Immediate second cast should fail
	_, err = op.ExecuteCastSpell(input, now, func() float64 { return 0.5 })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "global cooldown")

	// After 2 seconds, should succeed
	_, err = op.ExecuteCastSpell(input, now.Add(3*time.Second), func() float64 { return 0.5 })
	assert.NoError(t, err)
}

func TestExecuteCastSpell_SkillCooldown(t *testing.T) {
	op := NewSpellOperation()

	input := CastSpellInput{
		CasterID:      "caster_1",
		SkillID:       "fireball",
		SkillName:     "火球术",
		SkillDamage:   50,
		SkillCooldown: 10, // 10 second cooldown
		CasterAttack:  100,
	}

	now := time.Now()
	_, err := op.ExecuteCastSpell(input, now, func() float64 { return 0.5 })
	assert.NoError(t, err)

	// Try again after 5 seconds (should fail)
	_, err = op.ExecuteCastSpell(input, now.Add(5*time.Second), func() float64 { return 0.5 })
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "冷却")

	// After 11 seconds, should succeed
	_, err = op.ExecuteCastSpell(input, now.Add(11*time.Second), func() float64 { return 0.5 })
	assert.NoError(t, err)
}

func TestCalculateElementMult_Advantage(t *testing.T) {
	assert.Equal(t, 1.5, calculateElementMult("fire", "wood"))
	assert.Equal(t, 1.5, calculateElementMult("water", "fire"))
	assert.Equal(t, 1.5, calculateElementMult("wood", "earth"))
}

func TestCalculateElementMult_Disadvantage(t *testing.T) {
	assert.Equal(t, 0.7, calculateElementMult("wood", "fire"))
	assert.Equal(t, 0.7, calculateElementMult("fire", "water"))
}

func TestCalculateElementMult_Neutral(t *testing.T) {
	assert.Equal(t, 1.0, calculateElementMult("fire", "metal"))
	assert.Equal(t, 1.0, calculateElementMult("", "fire"))
	assert.Equal(t, 1.0, calculateElementMult("fire", ""))
}

func TestCalculateSpellDamage_RealmSuppression(t *testing.T) {
	op := NewSpellOperation()

	// Higher realm caster
	input := CastSpellInput{
		SkillDamage:   100,
		CasterAttack:  0,
		CasterRealm:   5,
		TargetRealm:   1, // 4 realm difference
	}

	damage := op.calculateSpellDamage(input, func() float64 { return 1.0 })
	// realm_mult = 1.0 + 4*0.2 = 1.8
	assert.True(t, damage > 150)
}

func TestCalculateSpellDamage_RealmDisadvantage(t *testing.T) {
	op := NewSpellOperation()

	input := CastSpellInput{
		SkillDamage:   100,
		CasterAttack:  0,
		CasterRealm:   1,
		TargetRealm:   5, // -4 realm difference
	}

	damage := op.calculateSpellDamage(input, func() float64 { return 0.0 })
	// realm_mult = max(0.5, 1.0 + (-4)*0.3) = max(0.5, -0.2) = 0.5
	// variance = 0.9 + 0*0.2 = 0.9
	// damage = 100 * 1.0 * 1.0 * 0.5 * 0.9 = 45
	assert.True(t, damage <= 50)
}

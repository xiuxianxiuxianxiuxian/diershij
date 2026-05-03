package service

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestNewCombatRule(t *testing.T) {
    rule := NewCombatRule()
    assert.NotNil(t, rule)
    assert.Equal(t, 0.15, rule.realmSuppressionPerLevel)
    assert.Equal(t, 0.05, rule.baseCritRate)
    assert.Equal(t, 1.5, rule.baseCritDamage)
}

func TestCalculateElementInteraction_Advantage(t *testing.T) {
    rule := NewCombatRule()

    // Fire attacks Metal → 1.5x
    factor := rule.CalculateElementInteraction("fire", "metal")
    assert.Equal(t, 1.5, factor)

    // Water attacks Fire → 1.5x
    factor = rule.CalculateElementInteraction("water", "fire")
    assert.Equal(t, 1.5, factor)

    // Wood attacks Earth → 1.5x
    factor = rule.CalculateElementInteraction("wood", "earth")
    assert.Equal(t, 1.5, factor)

    // Metal attacks Wood → 1.5x
    factor = rule.CalculateElementInteraction("metal", "wood")
    assert.Equal(t, 1.5, factor)

    // Earth attacks Water → 1.5x
    factor = rule.CalculateElementInteraction("earth", "water")
    assert.Equal(t, 1.5, factor)

    // Fire attacks Wood → 1.3x
    factor = rule.CalculateElementInteraction("fire", "wood")
    assert.Equal(t, 1.3, factor)
}

func TestCalculateElementInteraction_Disadvantage(t *testing.T) {
    rule := NewCombatRule()

    // Metal attacked by Fire → reciprocal 1/1.5 ≈ 0.667
    factor := rule.CalculateElementInteraction("metal", "fire")
    assert.InDelta(t, 1.0/1.5, factor, 0.001)

    // Fire attacked by Water → reciprocal 1/1.5 ≈ 0.667
    factor = rule.CalculateElementInteraction("fire", "water")
    assert.InDelta(t, 1.0/1.5, factor, 0.001)

    // Wood attacked by Metal → reciprocal 1/1.5 ≈ 0.667
    factor = rule.CalculateElementInteraction("wood", "metal")
    assert.InDelta(t, 1.0/1.5, factor, 0.001)
}

func TestCalculateElementInteraction_Neutral(t *testing.T) {
    rule := NewCombatRule()

    // Same element → 1.0
    factor := rule.CalculateElementInteraction("fire", "fire")
    assert.Equal(t, 1.0, factor)

    // Empty element → 1.0
    factor = rule.CalculateElementInteraction("", "fire")
    assert.Equal(t, 1.0, factor)

    factor = rule.CalculateElementInteraction("fire", "")
    assert.Equal(t, 1.0, factor)
}

func TestCalculateDamage_BasicAttack(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.5 } // always miss crit/dodge

    input := DamageInput{
        AttackerAttackPower: 100,
        DefenderDefense:     50,
        DefenderDodgeRate:   0.0,
        AttackerPenetration: 0.0,
        AttackerCritRate:    0.0,
        AttackerCritDamage:  0.0,
        SkillMultiplier:     1.0,
        AttackerElement:     "fire",
        DefenderElement:     "metal",
        RealmDiff:           0,
        MethodDamageBonus:   0.0,
        DefenderReduction:   0.0,
    }

    result := rule.CalculateDamage(input, deterministicRand)
    assert.False(t, result.IsDodged)

    // base = 100, skill = 1.0, realm = 1.0, element = 1.5
    // def_reduction = 50/(50+100) = 0.333, defense_factor = 0.667
    // method = 1.0
    // damage = 100 * 1.0 * 1.0 * 1.5 * 0.667 * 1.0 = 100.0
    assert.InDelta(t, 100.0, result.FinalDamage, 0.1)
    assert.Equal(t, 1.5, result.ElementFactor)
    assert.Equal(t, 1.0, result.RealmFactor)
}

func TestCalculateDamage_RealmSuppression(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.5 }

    input := DamageInput{
        AttackerAttackPower: 100,
        DefenderDefense:     0,
        SkillMultiplier:     1.0,
        RealmDiff:           2, // attacker 2 realms higher
    }

    result := rule.CalculateDamage(input, deterministicRand)
    // realm = 1 + 2*0.15 = 1.30
    // damage = 100 * 1.0 * 1.30 * 1.0 * 1.0 * 1.0 = 130.0
    assert.InDelta(t, 130.0, result.FinalDamage, 0.1)
    assert.InDelta(t, 1.30, result.RealmFactor, 0.001)
}

func TestCalculateDamage_RealmDisadvantage(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.5 }

    input := DamageInput{
        AttackerAttackPower: 100,
        DefenderDefense:     0,
        SkillMultiplier:     1.0,
        RealmDiff:           -1, // attacker 1 realm lower
    }

    result := rule.CalculateDamage(input, deterministicRand)
    // realm = 1 - 0.15 = 0.85
    // damage = 100 * 1.0 * 0.85 * 1.0 * 1.0 * 1.0 = 85.0
    assert.InDelta(t, 85.0, result.FinalDamage, 0.1)
    assert.InDelta(t, 0.85, result.RealmFactor, 0.001)
}

func TestCalculateDamage_DefenseReduction(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.5 }

    // 50 defense: reduction = 50/150 = 0.333
    input := DamageInput{
        AttackerAttackPower: 100,
        DefenderDefense:     50,
        SkillMultiplier:     1.0,
        DefenderDodgeRate:   0.0,
        AttackerPenetration: 0.0,
        AttackerCritRate:    0.0,
        AttackerCritDamage:  0.0,
    }

    result := rule.CalculateDamage(input, deterministicRand)
    // def_reduction = 0.333, defense_factor = 0.667
    // damage = 100 * 1.0 * 1.0 * 1.0 * 0.667 * 1.0 = 66.67
    assert.InDelta(t, 66.67, result.FinalDamage, 0.1)
    assert.InDelta(t, 0.667, result.DefenseFactor, 0.001)
}

func TestCalculateDamage_HighDefense(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.5 }

    // 900 defense: reduction = 900/1000 = 0.90
    input := DamageInput{
        AttackerAttackPower: 100,
        DefenderDefense:     900,
        SkillMultiplier:     1.0,
    }

    result := rule.CalculateDamage(input, deterministicRand)
    // def_reduction = 0.90, defense_factor = 0.10
    // damage = 100 * 1.0 * 1.0 * 1.0 * 0.10 * 1.0 = 10.0
    assert.InDelta(t, 10.0, result.FinalDamage, 0.1)
}

func TestCalculateDamage_Penetration(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.5 }

    input := DamageInput{
        AttackerAttackPower: 100,
        DefenderDefense:     50,
        SkillMultiplier:     1.0,
        AttackerPenetration: 0.5, // 50% penetration
    }

    result := rule.CalculateDamage(input, deterministicRand)
    // def_reduction = 0.333 * (1 - 0.5) = 0.167
    // defense_factor = 0.833
    // damage = 100 * 1.0 * 1.0 * 1.0 * 0.833 * 1.0 = 83.3
    assert.InDelta(t, 83.3, result.FinalDamage, 0.1)
}

func TestCalculateDamage_Crit(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.01 } // always crit (1%)

    input := DamageInput{
        AttackerAttackPower: 100,
        DefenderDefense:     0,
        SkillMultiplier:     1.0,
        AttackerCritRate:    0.10, // 10% crit rate
        AttackerCritDamage:  2.0,  // 200% crit damage
    }

    result := rule.CalculateDamage(input, deterministicRand)
    assert.True(t, result.IsCrit)
    // damage = 100 * 1.0 * 1.0 * 1.0 * 1.0 * 1.0 * 2.0 = 200.0
    assert.InDelta(t, 200.0, result.FinalDamage, 0.1)
}

func TestCalculateDamage_Dodge(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.10 } // 10%

    input := DamageInput{
        AttackerAttackPower: 100,
        DefenderDefense:     0,
        SkillMultiplier:     1.0,
        DefenderDodgeRate:   0.50, // 50% dodge
    }

    result := rule.CalculateDamage(input, deterministicRand)
    // 0.10 < 0.50 → dodged
    assert.True(t, result.IsDodged)
    assert.Equal(t, 0.0, result.FinalDamage)
}

func TestCalculateDamage_NoDodge(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.60 } // 60%

    input := DamageInput{
        AttackerAttackPower: 100,
        DefenderDefense:     0,
        SkillMultiplier:     1.0,
        DefenderDodgeRate:   0.50,
    }

    result := rule.CalculateDamage(input, deterministicRand)
    // 0.60 >= 0.50 → not dodged
    assert.False(t, result.IsDodged)
    assert.Equal(t, 100.0, result.FinalDamage)
}

func TestCalculateDamage_MethodBonus(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.5 }

    input := DamageInput{
        AttackerAttackPower: 100,
        DefenderDefense:     0,
        SkillMultiplier:     1.0,
        MethodDamageBonus:   0.20, // +20% from method
    }

    result := rule.CalculateDamage(input, deterministicRand)
    // damage = 100 * 1.0 * 1.0 * 1.0 * 1.0 * 1.2 = 120.0
    assert.InDelta(t, 120.0, result.FinalDamage, 0.1)
}

func TestCalculateDamage_SkillMultiplier(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.5 }

    input := DamageInput{
        AttackerAttackPower: 100,
        DefenderDefense:     0,
        SkillMultiplier:     2.5, // powerful skill
    }

    result := rule.CalculateDamage(input, deterministicRand)
    // damage = 100 * 2.5 = 250.0
    assert.InDelta(t, 250.0, result.FinalDamage, 0.1)
}

func TestCalculateDamage_CombinedScenario(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.5 }

    // Golden Core attacker (level 3) vs Foundation defender (level 2)
    // fire vs metal (element advantage)
    // High defense, some penetration, crit check fails, method bonus
    input := DamageInput{
        AttackerAttackPower: 200,
        DefenderDefense:     80,
        SkillMultiplier:     1.8,
        AttackerElement:     "fire",
        DefenderElement:     "metal",
        RealmDiff:           1, // +1 realm
        AttackerPenetration: 0.3,
        AttackerCritRate:    0.15,
        AttackerCritDamage:  1.8,
        MethodDamageBonus:   0.10,
        DefenderReduction:   0.05,
    }

    result := rule.CalculateDamage(input, deterministicRand)
    // realm = 1.15, element = 1.5
    // def_reduction = 80/180 = 0.444, effective = 0.444 * (1-0.3) = 0.311
    // total_reduction = 0.311 + 0.05*(1-0.311) = 0.311 + 0.034 = 0.345
    // defense_factor = 0.655
    // method = 1.10
    // raw = 200 * 1.8 * 1.15 * 1.5 * 0.655 * 1.10 = 200 * 1.8 * 1.15 * 1.5 * 0.655 * 1.10
    // raw = 200 * 1.8 = 360, * 1.15 = 414, * 1.5 = 621, * 0.655 = 406.7, * 1.10 = 447.4
    assert.True(t, result.FinalDamage > 400 && result.FinalDamage < 500)
    assert.False(t, result.IsDodged)
    assert.Equal(t, 1.5, result.ElementFactor)
    assert.InDelta(t, 1.15, result.RealmFactor, 0.001)
}

func TestResolveCombat_SimpleWin(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.5 }

    attacker := CombatEntity{
        ID: "player_1", HP: 500, MaxHP: 500,
        AttackPower: 100, Defense: 20,
        SkillMultiplier: 1.0,
    }
    defender := CombatEntity{
        ID: "player_2", HP: 300, MaxHP: 300,
        AttackPower: 50, Defense: 10,
        SkillMultiplier: 1.0,
    }

    result := rule.ResolveCombat(attacker, defender, deterministicRand)
    assert.NotEmpty(t, result.Turns)
    assert.True(t, result.Rounds >= 1)
    assert.True(t, result.TotalDamage > 0)
}

func TestResolveCombat_EvenMatch(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.5 }

    a := CombatEntity{
        ID: "a", HP: 1000, MaxHP: 1000,
        AttackPower: 50, Defense: 30,
        SkillMultiplier: 1.0,
    }
    b := CombatEntity{
        ID: "b", HP: 1000, MaxHP: 1000,
        AttackPower: 50, Defense: 30,
        SkillMultiplier: 1.0,
    }

    result := rule.ResolveCombat(a, b, deterministicRand)
    // Equal stats → should eventually end (someone wins by max rounds or HP depletion)
    assert.NotEmpty(t, result.Turns)
    assert.NotEmpty(t, result.Winner)
    assert.NotEmpty(t, result.Loser)
}

func TestResolveCombat_HighDodge(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.01 } // always dodge (1%)

    attacker := CombatEntity{
        ID: "attacker", HP: 100, MaxHP: 100,
        AttackPower: 100, Defense: 0,
        SkillMultiplier: 1.0,
    }
    defender := CombatEntity{
        ID: "defender", HP: 50, MaxHP: 50,
        AttackPower: 1, Defense: 0,
        DodgeRate:     0.99, // 99% dodge
        SkillMultiplier: 1.0,
    }

    result := rule.ResolveCombat(attacker, defender, deterministicRand)
    // Defender should dodge most attacks, attacker eventually wins
    assert.NotEmpty(t, result.Turns)
    dodgedCount := 0
    for _, turn := range result.Turns {
        if turn.IsDodged {
            dodgedCount++
        }
    }
    assert.True(t, dodgedCount > 0)
}

func TestResolveCombat_MaxRounds(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.5 }

    // Very high HP, low attack → should hit max rounds
    a := CombatEntity{
        ID: "a", HP: 100000, MaxHP: 100000,
        AttackPower: 1, Defense: 100,
        SkillMultiplier: 1.0,
    }
    b := CombatEntity{
        ID: "b", HP: 100000, MaxHP: 100000,
        AttackPower: 1, Defense: 100,
        SkillMultiplier: 1.0,
    }

    result := rule.ResolveCombat(a, b, deterministicRand)
    // Should hit max rounds (100)
    assert.True(t, result.Rounds <= 100)
}

func TestResolveCombat_TurnOrder(t *testing.T) {
    rule := NewCombatRule()
    deterministicRand := func() float64 { return 0.5 }

    // First attacker has lower speed, defender has higher
    // "speed" = penetration + crit_rate*100
    a := CombatEntity{
        ID: "slow", HP: 200, MaxHP: 200,
        AttackPower: 100, Penetration: 0.0,
        SkillMultiplier: 1.0,
    }
    b := CombatEntity{
        ID: "fast", HP: 200, MaxHP: 200,
        AttackPower: 100, Penetration: 0.1,
        SkillMultiplier: 1.0,
    }

    result := rule.ResolveCombat(a, b, deterministicRand)
    // First turn should be from the faster entity
    assert.Equal(t, "fast", result.Turns[0].AttackerID)
}

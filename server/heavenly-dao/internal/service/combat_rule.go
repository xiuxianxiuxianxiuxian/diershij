package service

import "math"

// CombatRule handles damage calculation, element interaction, and combat resolution.
type CombatRule struct {
    realmSuppressionPerLevel float64
    baseCritRate             float64
    baseCritDamage           float64
}

// NewCombatRule creates a new CombatRule with default configuration.
func NewCombatRule() *CombatRule {
    return &CombatRule{
        realmSuppressionPerLevel: 0.15,
        baseCritRate:             0.05,
        baseCritDamage:           1.5,
    }
}

// elementCounters defines which attacker element has advantage over which defender element.
// Format: attacker_element_defender_element -> multiplier
var elementCounters = map[string]float64{
    "fire_metal": 1.5,
    "fire_wood":  1.3,
    "water_fire": 1.5,
    "wood_earth": 1.5,
    "metal_wood": 1.5,
    "earth_water": 1.5,
}

// CalculateElementInteraction returns the element advantage factor.
// 1.0 means no advantage, > 1.0 means attacker has advantage.
func (c *CombatRule) CalculateElementInteraction(attackerElement string, defenderElement string) float64 {
    if attackerElement == "" || defenderElement == "" {
        return 1.0
    }
    if attackerElement == defenderElement {
        return 1.0
    }

    key := attackerElement + "_" + defenderElement
    if factor, ok := elementCounters[key]; ok {
        return factor
    }

    // If defender has advantage over attacker, attacker is weakened (reciprocal)
    reverseKey := defenderElement + "_" + attackerElement
    if factor, ok := elementCounters[reverseKey]; ok {
        return 1.0 / factor
    }

    return 1.0
}

// DamageInput holds the inputs needed to calculate combat damage.
type DamageInput struct {
    AttackerAttackPower float64
    DefenderDefense     float64
    DefenderDodgeRate   float64
    AttackerPenetration float64
    AttackerCritRate    float64
    AttackerCritDamage  float64
    SkillMultiplier     float64
    AttackerElement     string
    DefenderElement     string
    RealmDiff           int // positive = attacker higher realm, negative = attacker lower
    MethodDamageBonus   float64 // e.g., 0.20 for +20%
    DefenderReduction   float64 // flat damage reduction percentage (0-1)
}

// DamageResult holds the result of a damage calculation.
type DamageResult struct {
    FinalDamage     float64
    IsCrit          bool
    IsDodged        bool
    ElementFactor   float64
    RealmFactor     float64
    DefenseFactor   float64
}

// CalculateDamage computes the damage dealt from an attack.
//
// Formula:
//
//	damage = base_atk × skill_mult × realm_suppression × element × (1 - def_reduction) × method_bonus
//
// Where:
//   - base_atk = attacker.attack_power
//   - skill_mult = skill multiplier (default 1.0 for basic attack)
//   - realm_suppression = 1.0 + realm_diff × 0.15
//   - element = element interaction factor (1.0-1.5 or 0.67-1.0)
//   - def_reduction = defense / (defense + 100) × (1 - penetration)
//   - method_bonus = 1.0 + method_damage_bonus
func (c *CombatRule) CalculateDamage(input DamageInput, randFloat func() float64) DamageResult {
    // Check dodge first
    dodgeRate := math.Max(0, math.Min(1, input.DefenderDodgeRate))
    if randFloat() < dodgeRate {
        return DamageResult{
            FinalDamage: 0,
            IsDodged:    true,
        }
    }

    baseAtk := input.AttackerAttackPower
    if baseAtk <= 0 {
        return DamageResult{}
    }

    // Skill multiplier (default 1.0)
    skillMult := input.SkillMultiplier
    if skillMult <= 0 {
        skillMult = 1.0
    }

    // Realm suppression: positive diff = attacker advantage, negative = disadvantage
    realmFactor := 1.0 + float64(input.RealmDiff)*c.realmSuppressionPerLevel
    realmFactor = math.Max(0.1, realmFactor)

    // Element interaction
    elementFactor := c.CalculateElementInteraction(input.AttackerElement, input.DefenderElement)

    // Defense reduction: defense / (defense + 100), reduced by penetration
    defReduction := input.DefenderDefense / (input.DefenderDefense + 100.0)
    penetration := math.Max(0, math.Min(1, input.AttackerPenetration))
    effectiveDefReduction := defReduction * (1.0 - penetration)

    // Also apply flat damage reduction from the defender
    flatReduction := math.Max(0, math.Min(1, input.DefenderReduction))
    totalReduction := effectiveDefReduction + flatReduction*(1-effectiveDefReduction)
    totalReduction = math.Min(0.95, totalReduction) // cap at 95% reduction

    defenseFactor := 1.0 - totalReduction

    // Method bonus
    methodBonus := 1.0 + input.MethodDamageBonus

    rawDamage := baseAtk * skillMult * realmFactor * elementFactor * defenseFactor * methodBonus

    // Check crit
    critRate := math.Max(0, input.AttackerCritRate)
    if critRate <= 0 {
        critRate = c.baseCritRate
    }
    isCrit := randFloat() < critRate

    critMult := input.AttackerCritDamage
    if critMult <= 0 {
        critMult = c.baseCritDamage
    }

    finalDamage := rawDamage
    if isCrit {
        finalDamage = rawDamage * critMult
    }

    finalDamage = math.Max(0, finalDamage)

    return DamageResult{
        FinalDamage:   math.Round(finalDamage*100) / 100,
        IsCrit:        isCrit,
        IsDodged:      false,
        ElementFactor: elementFactor,
        RealmFactor:   realmFactor,
        DefenseFactor: defenseFactor,
    }
}

// CombatTurn represents a single turn in combat.
type CombatTurn struct {
    AttackerID  string
    DefenderID  string
    Damage      float64
    IsCrit      bool
    IsDodged    bool
}

// CombatResult represents the final result of a combat encounter.
type CombatResult struct {
    Winner       string
    Loser        string
    Turns        []CombatTurn
    AttackerHP   float64
    DefenderHP   float64
    TotalDamage  float64
    Rounds       int
}

// CombatEntity represents a participant in combat.
type CombatEntity struct {
    ID              string
    HP              float64
    MaxHP           float64
    AttackPower     float64
    Defense         float64
    DodgeRate       float64
    CritRate        float64
    CritDamage      float64
    Penetration     float64
    DamageReduction float64
    Element         string
    RealmLevel      int
    SkillMultiplier float64
    MethodBonus     float64
}

// ResolveCombat runs a turn-based combat simulation between two entities.
// Returns the combat result with winner, turns, and final state.
//
// Rules:
//   - Each round, both entities take one turn (attacker first, then defender counter-attacks)
//   - Speed determines who attacks first (higher speed = first)
//   - Combat ends when one entity's HP reaches 0
//   - Maximum 100 rounds to prevent infinite loops
func (c *CombatRule) ResolveCombat(attacker, defender CombatEntity, randFloat func() float64) CombatResult {
    var turns []CombatTurn
    totalDamage := 0.0
    maxRounds := 100

    // Determine who attacks first based on speed
    first, second := attacker, defender
    if defender.Speed() > attacker.Speed() {
        first, second = defender, attacker
    }

    for round := 0; round < maxRounds; round++ {
        // First entity attacks
        dmg := c.damageForTurn(first, second, randFloat)
        totalDamage += dmg.FinalDamage
        second.HP -= dmg.FinalDamage

        turns = append(turns, CombatTurn{
            AttackerID: first.ID,
            DefenderID: second.ID,
            Damage:     dmg.FinalDamage,
            IsCrit:     dmg.IsCrit,
            IsDodged:   dmg.IsDodged,
        })

        if second.HP <= 0 {
            return CombatResult{
                Winner:      first.ID,
                Loser:       second.ID,
                Turns:       turns,
                AttackerHP:  first.HP,
                DefenderHP:  0,
                TotalDamage: totalDamage,
                Rounds:      round + 1,
            }
        }

        // Second entity counter-attacks
        dmg = c.damageForTurn(second, first, randFloat)
        totalDamage += dmg.FinalDamage
        first.HP -= dmg.FinalDamage

        turns = append(turns, CombatTurn{
            AttackerID: second.ID,
            DefenderID: first.ID,
            Damage:     dmg.FinalDamage,
            IsCrit:     dmg.IsCrit,
            IsDodged:   dmg.IsDodged,
        })

        if first.HP <= 0 {
            return CombatResult{
                Winner:      second.ID,
                Loser:       first.ID,
                Turns:       turns,
                AttackerHP:  0,
                DefenderHP:  second.HP,
                TotalDamage: totalDamage,
                Rounds:      round + 1,
            }
        }
    }

    // Max rounds reached: higher HP wins
    if first.HP >= second.HP {
        return CombatResult{
            Winner:      first.ID,
            Loser:       second.ID,
            Turns:       turns,
            AttackerHP:  first.HP,
            DefenderHP:  second.HP,
            TotalDamage: totalDamage,
            Rounds:      maxRounds,
        }
    }

    return CombatResult{
        Winner:      second.ID,
        Loser:       first.ID,
        Turns:       turns,
        AttackerHP:  first.HP,
        DefenderHP:  second.HP,
        TotalDamage: totalDamage,
        Rounds:      maxRounds,
    }
}

func (c *CombatRule) damageForTurn(attacker, defender CombatEntity, randFloat func() float64) DamageResult {
    return c.CalculateDamage(DamageInput{
        AttackerAttackPower: attacker.AttackPower,
        DefenderDefense:     defender.Defense,
        DefenderDodgeRate:   defender.DodgeRate,
        AttackerPenetration: attacker.Penetration,
        AttackerCritRate:    attacker.CritRate,
        AttackerCritDamage:  attacker.CritDamage,
        SkillMultiplier:     attacker.SkillMultiplier,
        AttackerElement:     attacker.Element,
        DefenderElement:     defender.Element,
        RealmDiff:           attacker.RealmLevel - defender.RealmLevel,
        MethodDamageBonus:   attacker.MethodBonus,
        DefenderReduction:   defender.DamageReduction,
    }, randFloat)
}

// Speed returns the speed value for turn order determination.
// CombatEntity doesn't have a Speed field directly, we compute from attributes.
func (ce CombatEntity) Speed() float64 {
    return ce.Penetration + ce.CritRate*100 // simplified: use combined stats as proxy
}

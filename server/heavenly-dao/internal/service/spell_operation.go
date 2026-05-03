package service

import (
	"fmt"
	"math"
	"time"
)

// SpellOperation handles casting spells/skills.
type SpellOperation struct {
	cooldownMap    map[string]time.Time
	skillCooldowns map[string]map[string]time.Time // entityID -> skillID -> lastUsed
}

// NewSpellOperation creates a new SpellOperation.
func NewSpellOperation() *SpellOperation {
	return &SpellOperation{
		cooldownMap:    make(map[string]time.Time),
		skillCooldowns: make(map[string]map[string]time.Time),
	}
}

// CastSpellInput holds inputs for casting a spell.
type CastSpellInput struct {
	CasterID       string
	TargetID       string
	SkillID        string
	SkillName      string
	SkillDamage    float64
	SkillCooldown  int // seconds
	CasterAttack   float64
	CasterElement  string
	TargetElement  string
	CasterRealm    int
	TargetRealm    int
	IsTargetSelf   bool
}

// CastSpellResult holds the outcome of a spell cast.
type CastSpellResult struct {
	Success   bool
	Damage    float64
	Healed    float64
	Effect    string
	Message   string
}

// ExecuteCastSpell casts a spell/skill.
func (op *SpellOperation) ExecuteCastSpell(input CastSpellInput, now time.Time, randFloat func() float64) (*CastSpellResult, error) {
	// Check global cooldown
	if lastTime, ok := op.cooldownMap[input.CasterID]; ok {
		remaining := time.Second * 2 - now.Sub(lastTime)
		if remaining > 0 {
			return nil, fmt.Errorf("global cooldown: %v remaining", remaining)
		}
	}

	// Check skill-specific cooldown
	if skillCooldowns, exists := op.skillCooldowns[input.CasterID]; exists {
		if lastUsed, ok := skillCooldowns[input.SkillID]; ok {
			cooldownDuration := time.Duration(input.SkillCooldown) * time.Second
			remaining := cooldownDuration - now.Sub(lastUsed)
			if remaining > 0 {
				return nil, fmt.Errorf("技能冷却中：%v remaining", remaining)
			}
		}
	}

	// Calculate damage/effect
	damage := op.calculateSpellDamage(input, randFloat)

	// Set cooldowns
	op.cooldownMap[input.CasterID] = now
	if _, exists := op.skillCooldowns[input.CasterID]; !exists {
		op.skillCooldowns[input.CasterID] = make(map[string]time.Time)
	}
	op.skillCooldowns[input.CasterID][input.SkillID] = now

	return &CastSpellResult{
		Success: true,
		Damage:  damage,
		Message: fmt.Sprintf("%s 施展了 %s", input.CasterID, input.SkillName),
	}, nil
}

func (op *SpellOperation) calculateSpellDamage(input CastSpellInput, randFloat func() float64) float64 {
	// Base damage from skill
	baseDamage := input.SkillDamage

	// Attacker attack power modifier
	attackMult := 1.0 + input.CasterAttack/100.0

	// Element interaction
	elementMult := calculateElementMult(input.CasterElement, input.TargetElement)

	// Realm suppression
	realmDiff := input.CasterRealm - input.TargetRealm
	realmMult := 1.0
	if realmDiff > 0 {
		realmMult = 1.0 + float64(realmDiff)*0.2
	} else if realmDiff < 0 {
		realmMult = math.Max(0.5, 1.0+float64(realmDiff)*0.3)
	}

	// Random variance (+/- 10%)
	variance := 0.9 + randFloat()*0.2

	damage := baseDamage * attackMult * elementMult * realmMult * variance
	return math.Max(0, damage)
}

func calculateElementMult(casterElement, targetElement string) float64 {
	if casterElement == "" || targetElement == "" {
		return 1.0
	}

	// Element advantages
	advantages := map[string]string{
		"fire":  "wood",
		"wood":  "earth",
		"earth": "water",
		"water": "fire",
		"metal": "wood",
		"wind":  "metal",
		"ice":   "water",
		"lightning": "water",
	}

	if target, ok := advantages[casterElement]; ok && target == targetElement {
		return 1.5 // advantage
	}

	if caster, ok := advantages[targetElement]; ok && caster == casterElement {
		return 0.7 // disadvantage
	}

	return 1.0 // neutral
}

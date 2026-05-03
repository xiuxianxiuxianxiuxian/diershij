package service

import (
	"fmt"
	"time"

	"github.com/cultivation-world/shared/types"
)

// CombatOperation handles the high-level combat operation workflow.
type CombatOperation struct {
	combatRule     *CombatRule
	cooldownMap    map[string]time.Time
	cooldownPeriod time.Duration
}

// NewCombatOperation creates a new CombatOperation.
func NewCombatOperation(cooldownPeriod time.Duration) *CombatOperation {
	return &CombatOperation{
		combatRule:     NewCombatRule(),
		cooldownMap:    make(map[string]time.Time),
		cooldownPeriod: cooldownPeriod,
	}
}

// CombatOpInput holds all inputs needed to execute a combat operation.
type CombatOpInput struct {
	AttackerID        string
	DefenderID        string
	AttackerName      string
	DefenderName      string
	AttackerRealm     types.CultivationRealm
	DefenderRealm     types.CultivationRealm
	AttackerHP        float64
	DefenderHP        float64
	AttackerMaxHP     float64
	DefenderMaxHP     float64
	AttackerAttack    float64
	DefenderAttack    float64
	AttackerDefense   float64
	DefenderDefense   float64
	AttackerElement   string
	DefenderElement   string
	AttackerPenetrate float64
	DefenderPenetrate float64
	AttackerCritRate  float64
	DefenderCritRate  float64
	AttackerDodgeRate float64
	DefenderDodgeRate float64
	IsSelfDefense     bool
}

// CombatOpResult holds the full outcome of a combat operation.
type CombatOpResult struct {
	WinnerID         string
	LoserID          string
	Rounds           int
	TotalDamage      float64
	KarmaChange      int
	LootItems        []LootItem
	Injuries         []OpInjury
	IsDraw           bool
	Message          string
}

// LootItem represents an item dropped after combat.
type LootItem struct {
	Name     string
	Type     string // spirit_stone / material / pill / artifact
	Quantity int
	Quality  string
}

// OpInjury represents an injury sustained in combat.
type OpInjury struct {
	Type        string
	Severity    int
	Description string
	HealTime    int64
}

// ExecuteCombat runs a full combat operation between two entities.
func (op *CombatOperation) ExecuteCombat(input CombatOpInput, now time.Time, randFloat func() float64) (*CombatOpResult, error) {
	// Check cooldown
	if lastTime, ok := op.cooldownMap[input.AttackerID]; ok {
		if now.Sub(lastTime) < op.cooldownPeriod {
			return nil, fmt.Errorf("combat cooldown: %v remaining", op.cooldownPeriod-now.Sub(lastTime))
		}
	}

	// Build combat entities
	attacker := CombatEntity{
		ID:              input.AttackerID,
		HP:              input.AttackerHP,
		MaxHP:           input.AttackerMaxHP,
		AttackPower:     input.AttackerAttack,
		Defense:         input.AttackerDefense,
		Element:         input.AttackerElement,
		Penetration:     input.AttackerPenetrate,
		CritRate:        input.AttackerCritRate,
		CritDamage:      1.5,
		DodgeRate:       input.AttackerDodgeRate,
		RealmLevel:      types.CultivationRealmLevel(input.AttackerRealm),
		SkillMultiplier: 1.0,
		MethodBonus:     0,
	}

	defender := CombatEntity{
		ID:              input.DefenderID,
		HP:              input.DefenderHP,
		MaxHP:           input.DefenderMaxHP,
		AttackPower:     input.DefenderAttack,
		Defense:         input.DefenderDefense,
		Element:         input.DefenderElement,
		Penetration:     input.DefenderPenetrate,
		CritRate:        input.DefenderCritRate,
		CritDamage:      1.5,
		DodgeRate:       input.DefenderDodgeRate,
		RealmLevel:      types.CultivationRealmLevel(input.DefenderRealm),
		SkillMultiplier: 1.0,
		MethodBonus:     0,
	}

	// Run turn-based combat
	combatResult := op.combatRule.ResolveCombat(attacker, defender, randFloat)

	// Determine winner
	winnerID := combatResult.Winner
	loserID := combatResult.Loser

	attackerWon := winnerID == input.AttackerID
	draw := winnerID == "" && loserID == ""

	if draw {
		return &CombatOpResult{
			IsDraw:      true,
			Rounds:      combatResult.Rounds,
			TotalDamage: combatResult.TotalDamage,
			Message:     "战斗平局",
		}, nil
	}

	// Calculate karma change
	karmaChange := op.calculateCombatKarma(input, attackerWon)

	// Generate loot (only for attacker wins)
	var lootItems []LootItem
	if attackerWon {
		lootItems = op.generateLoot(input.DefenderRealm, randFloat)
	}

	// Generate injuries
	injuries := op.generateInjuries(input, combatResult, attackerWon)

	// Set cooldown
	op.cooldownMap[input.AttackerID] = now

	winnerName := input.AttackerName
	loserName := input.DefenderName
	if !attackerWon {
		winnerName = input.DefenderName
		loserName = input.AttackerName
	}

	return &CombatOpResult{
		WinnerID:     winnerID,
		LoserID:      loserID,
		Rounds:       combatResult.Rounds,
		TotalDamage:  combatResult.TotalDamage,
		KarmaChange:  karmaChange,
		LootItems:    lootItems,
		Injuries:     injuries,
		Message:      fmt.Sprintf("%s 在 %d 回合内击败了 %s", winnerName, combatResult.Rounds, loserName),
	}, nil
}

func (op *CombatOperation) calculateCombatKarma(input CombatOpInput, attackerWon bool) int {
	if input.IsSelfDefense {
		return 0
	}

	if attackerWon {
		realmDiff := types.CultivationRealmLevel(input.AttackerRealm) - types.CultivationRealmLevel(input.DefenderRealm)
		if realmDiff > 2 {
			return -20
		}
		return -10
	}

	return -5
}

func (op *CombatOperation) generateLoot(defenderRealm types.CultivationRealm, randFloat func() float64) []LootItem {
	realmLevel := types.CultivationRealmLevel(defenderRealm)
	if realmLevel < 0 {
		realmLevel = 0
	}

	baseStones := (realmLevel + 1) * 10
	stones := int(float64(baseStones) * (0.5 + randFloat()*1.0))

	loot := []LootItem{
		{
			Name:     "spirit_stones",
			Type:     "spirit_stone",
			Quantity: stones,
			Quality:  "standard",
		},
	}

	if randFloat() < 0.3 {
		loot = append(loot, LootItem{
			Name:     "spirit_herb",
			Type:     "material",
			Quantity: 1 + int(randFloat()*3),
			Quality:  "common",
		})
	}

	if randFloat() < 0.1 {
		loot = append(loot, LootItem{
			Name:     "healing_pill",
			Type:     "pill",
			Quantity: 1,
			Quality:  "fine",
		})
	}

	return loot
}

func (op *CombatOperation) generateInjuries(input CombatOpInput, result CombatResult, attackerWon bool) []OpInjury {
	var injuries []OpInjury

	// Loser takes injuries based on damage ratio
	loserMaxHP := input.DefenderMaxHP
	loserRemainingHP := result.DefenderHP
	if !attackerWon {
		loserMaxHP = input.AttackerMaxHP
		loserRemainingHP = result.AttackerHP
	}

	if loserMaxHP <= 0 {
		loserMaxHP = 100
	}

	damageRatio := 1.0 - (loserRemainingHP / loserMaxHP)

	if damageRatio > 0.5 {
		severity := int(damageRatio * 5)
		if severity > 5 {
			severity = 5
		}
		injuries = append(injuries, OpInjury{
			Type:        "internal",
			Severity:    severity,
			Description: "内伤",
			HealTime:    int64(severity * 2),
		})
	}

	if damageRatio > 0.8 {
		injuries = append(injuries, OpInjury{
			Type:        "external",
			Severity:    2,
			Description: "外伤",
			HealTime:    3,
		})
	}

	return injuries
}

// GetCooldownRemaining returns remaining cooldown time for an entity.
func (op *CombatOperation) GetCooldownRemaining(entityID string, now time.Time) time.Duration {
	if lastTime, ok := op.cooldownMap[entityID]; ok {
		remaining := op.cooldownPeriod - now.Sub(lastTime)
		if remaining > 0 {
			return remaining
		}
	}
	return 0
}

// ClearCooldown removes cooldown for an entity.
func (op *CombatOperation) ClearCooldown(entityID string) {
	delete(op.cooldownMap, entityID)
}

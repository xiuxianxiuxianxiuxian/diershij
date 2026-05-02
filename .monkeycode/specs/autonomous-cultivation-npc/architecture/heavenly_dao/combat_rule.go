package heavenlydao

import "fmt"

func (e *Engine) CalculateDamage(ctx *RuleContext) (DamageResult, error) {
	if ctx == nil || ctx.ActorState == nil || ctx.TargetState == nil {
		return DamageResult{}, fmt.Errorf("missing combat snapshots")
	}

	skillMultiplier, _ := ctx.Params["skill_multiplier"].(float64)
	if skillMultiplier == 0 {
		skillMultiplier = 1.0
	}

	baseDamage := ctx.ActorState.AttackPower * skillMultiplier
	realmDiff := float64(ctx.ActorState.RealmLevel - ctx.TargetState.RealmLevel)
	realmModifier := 1.0 + realmDiff*e.config.Combat.RealmSuppressionPerLevel
	if realmModifier < 0.2 {
		realmModifier = 0.2
	}

	elementKey, _ := ctx.Params["element_counter_key"].(string)
	elementModifier := e.config.Combat.ElementCounters[elementKey]
	if elementModifier == 0 {
		elementModifier = 1.0
	}

	defenseReduction := ctx.TargetState.Defense / (ctx.TargetState.Defense + 100.0)
	if defenseReduction > 0.9 {
		defenseReduction = 0.9
	}

	methodBonus := 1.0
	if bonus, ok := ctx.ActorState.ConsumedBonuses["method_damage_bonus"]; ok {
		methodBonus += bonus
	}

	finalDamage := baseDamage * realmModifier * elementModifier * (1.0 - defenseReduction) * methodBonus

	critRate := ctx.ActorState.CritRate / 100.0
	if critRate == 0 {
		critRate = e.config.Combat.BaseCritRate
	}
	critDamage := ctx.ActorState.CritDamage / 100.0
	if critDamage == 0 {
		critDamage = e.config.Combat.BaseCritDamage
	}

	critTriggered, _ := ctx.Params["force_crit"].(bool)
	if critTriggered {
		finalDamage *= critDamage
	}

	return DamageResult{
		BaseDamage:       baseDamage,
		FinalDamage:      finalDamage,
		RealmModifier:    realmModifier,
		ElementModifier:  elementModifier,
		DefenseReduction: defenseReduction,
		CritTriggered:    critTriggered,
		CritMultiplier:   critDamage,
	}, nil
}

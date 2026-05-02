package heavenlydao

import "fmt"

func (e *Engine) CalculateBreakthrough(ctx *RuleContext) (BreakthroughResult, error) {
	if ctx == nil || ctx.ActorState == nil {
		return BreakthroughResult{}, fmt.Errorf("invalid context")
	}

	targetRealm, _ := ctx.Params["target_realm"].(string)
	requiredTime, _ := ctx.Params["required_time"].(float64)
	cultivationTime, _ := ctx.Params["cultivation_time"].(float64)
	resourceBonus, _ := ctx.Params["resource_bonus"].(float64)

	baseSuccess := e.config.Breakthrough.BaseSuccessByRealm[targetRealm]
	accumulationFactor := 1.0
	if requiredTime > 0 {
		accumulationFactor = clamp(1.0+(cultivationTime/requiredTime), 1.0, 1.5)
	}
	methodQuality := clamp(ctx.ActorState.ActiveMethodQuality/100.0, 0.1, 2.0)
	mentalFactor := clamp(ctx.ActorState.MentalStability/100.0, 0.1, 1.0)
	luckFactor := 1.0 + ((ctx.ActorState.Luck - 50.0) / 200.0)

	successProbability := baseSuccess * accumulationFactor * methodQuality * (1.0 + resourceBonus) * mentalFactor * luckFactor
	successProbability = clamp(successProbability, e.config.Breakthrough.MinSuccessRate, e.config.Breakthrough.MaxSuccessRate)

	baseTribulationProbability := e.config.Tribulation.BaseProbabilityByRealm[targetRealm]
	karmaFactor := 1.0 + (ctx.ActorState.Karma / 1000.0)
	meritFactor := 1.0 - (ctx.ActorState.Merit / 2000.0)
	meritFactor = clamp(meritFactor, e.config.Tribulation.MeritFloorFactor, 1.0)
	luckResistance := 1.0 - ((ctx.ActorState.Luck / 100.0) * 0.2)

	tribulationProbability := baseTribulationProbability * karmaFactor * meritFactor * luckResistance
	tribulationProbability = clamp(tribulationProbability, e.config.Tribulation.MinProbability, e.config.Tribulation.MaxProbability)

	tribulationStrength := 100.0 * (1.0 + ctx.ActorState.Karma*e.config.Tribulation.StrengthPerKarma)
	tribulationStrength *= 1.0 + float64(ctx.ActorState.RecentBreakthroughs7d)*e.config.Tribulation.RecentStrengthBonus

	return BreakthroughResult{
		Success:                false,
		SuccessProbability:     successProbability,
		TribulationTriggered:   tribulationProbability >= 0.5,
		TribulationProbability: tribulationProbability,
		TribulationStrength:    tribulationStrength,
		CultivationLossPercent: e.config.Breakthrough.FailureCultivationLoss,
		CooldownHours:          e.config.Breakthrough.FailureCooldownPerRealm * maxInt(ctx.ActorState.RealmLevel, 1),
		MentalDamage:           e.config.Breakthrough.FailureMentalDamage,
	}, nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

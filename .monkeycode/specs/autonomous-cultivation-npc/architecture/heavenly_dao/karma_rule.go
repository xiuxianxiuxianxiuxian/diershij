package heavenlydao

import "fmt"

func (e *Engine) CalculateKarma(ctx *RuleContext) (KarmaResult, error) {
	if ctx == nil || ctx.ActorState == nil {
		return KarmaResult{}, fmt.Errorf("invalid context")
	}

	base := e.config.Karma.BaseValues[string(ctx.ActionType)]
	contextMultiplier := 1.0 - (ctx.ActorState.Karma / e.config.Karma.Cap)
	if contextMultiplier < 0.1 {
		contextMultiplier = 0.1
	}

	relationshipFactor := 1.0
	if ctx.TargetState != nil {
		for _, tag := range ctx.ActorState.RelationshipTags {
			switch tag {
			case "mentor_target":
				relationshipFactor = 2.0
			case "enemy_target":
				relationshipFactor = 0.5
			case "lover_target":
				relationshipFactor = 3.0
			}
		}
	}

	finalDelta := base * contextMultiplier * relationshipFactor

	return KarmaResult{
		BaseValue:          base,
		ContextMultiplier:  contextMultiplier,
		RelationshipFactor: relationshipFactor,
		FinalDelta:         finalDelta,
	}, nil
}

func (e *Engine) ApplyKarma(ctx *RuleContext) (KarmaResult, error) {
	result, err := e.CalculateKarma(ctx)
	if err != nil {
		return KarmaResult{}, err
	}

	if err := e.repo.AddKarma(ctx.ActorID, result.FinalDelta); err != nil {
		return KarmaResult{}, err
	}

	_ = e.bus.Publish(RuleEvent{
		Name:    "karma.applied",
		ActorID: ctx.ActorID,
		TargetID: ctx.TargetID,
		Payload: map[string]any{
			"action_type":          ctx.ActionType,
			"base_value":           result.BaseValue,
			"context_multiplier":   result.ContextMultiplier,
			"relationship_factor":  result.RelationshipFactor,
			"final_delta":          result.FinalDelta,
		},
	})

	return result, nil
}

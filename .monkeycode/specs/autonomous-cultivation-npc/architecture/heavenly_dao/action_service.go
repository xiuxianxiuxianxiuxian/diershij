package heavenlydao

import "fmt"

type ActionService struct {
	contextBuilder *ContextBuilder
	engine         HeavenlyDaoEngine
	repo           StateRepository
}

func NewActionService(contextBuilder *ContextBuilder, engine HeavenlyDaoEngine, repo StateRepository) *ActionService {
	return &ActionService{
		contextBuilder: contextBuilder,
		engine:         engine,
		repo:           repo,
	}
}

func (s *ActionService) HandleCombat(
	actorID EntityID,
	targetID EntityID,
	regionID RegionID,
	worldTime WorldTime,
	params map[string]any,
) (DamageResult, error) {
	ctx, err := s.contextBuilder.Build(actorID, targetID, ActionCombat, params, regionID, worldTime, worldTime.ServerTime.UnixNano())
	if err != nil {
		return DamageResult{}, err
	}

	damage, err := s.engine.CalculateDamage(ctx)
	if err != nil {
		return DamageResult{}, err
	}

	if err := s.repo.SaveOperationResult(actorID, ActionCombat, map[string]any{
		"target_id":     targetID,
		"final_damage":  damage.FinalDamage,
		"realm_modifier": damage.RealmModifier,
	}); err != nil {
		return DamageResult{}, err
	}

	_, _ = s.engine.ApplyKarma(ctx)
	return damage, nil
}

func (s *ActionService) HandleBreakthrough(
	actorID EntityID,
	regionID RegionID,
	worldTime WorldTime,
	params map[string]any,
) (BreakthroughResult, error) {
	ctx, err := s.contextBuilder.Build(actorID, "", ActionBreakthrough, params, regionID, worldTime, worldTime.ServerTime.UnixNano())
	if err != nil {
		return BreakthroughResult{}, err
	}

	result, err := s.engine.CalculateBreakthrough(ctx)
	if err != nil {
		return BreakthroughResult{}, err
	}

	if err := s.repo.SaveOperationResult(actorID, ActionBreakthrough, map[string]any{
		"success_probability":      result.SuccessProbability,
		"tribulation_probability":  result.TribulationProbability,
		"tribulation_strength":     result.TribulationStrength,
	}); err != nil {
		return BreakthroughResult{}, err
	}

	return result, nil
}

func (s *ActionService) HandleCreateMethod(
	actorID EntityID,
	regionID RegionID,
	worldTime WorldTime,
	params map[string]any,
) error {
	ctx, err := s.contextBuilder.Build(actorID, "", ActionCreateMethod, params, regionID, worldTime, worldTime.ServerTime.UnixNano())
	if err != nil {
		return err
	}

	requiredRealm, ok := params["required_realm_level"].(int)
	if !ok {
		return fmt.Errorf("missing required_realm_level")
	}
	requiredComprehension, ok := params["required_comprehension"].(float64)
	if !ok {
		return fmt.Errorf("missing required_comprehension")
	}

	return s.engine.CanCreateMethod(ctx, requiredRealm, requiredComprehension)
}

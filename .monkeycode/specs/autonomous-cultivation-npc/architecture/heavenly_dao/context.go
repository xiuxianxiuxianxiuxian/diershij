package heavenlydao

type SnapshotRepository interface {
	LoadEntitySnapshot(id EntityID) (*EntitySnapshot, error)
	LoadLocation(regionID RegionID) (*LocationInfo, error)
}

type ContextBuilder struct {
	repo SnapshotRepository
}

func NewContextBuilder(repo SnapshotRepository) *ContextBuilder {
	return &ContextBuilder{repo: repo}
}

func (b *ContextBuilder) Build(
	actorID EntityID,
	targetID EntityID,
	actionType ActionType,
	params map[string]any,
	regionID RegionID,
	worldTime WorldTime,
	seed int64,
) (*RuleContext, error) {
	actor, err := b.repo.LoadEntitySnapshot(actorID)
	if err != nil {
		return nil, err
	}

	var target *EntitySnapshot
	if targetID != "" {
		target, err = b.repo.LoadEntitySnapshot(targetID)
		if err != nil {
			return nil, err
		}
	}

	location, err := b.repo.LoadLocation(regionID)
	if err != nil {
		return nil, err
	}

	return &RuleContext{
		ActorID:     actorID,
		TargetID:    targetID,
		ActionType:  actionType,
		Params:      params,
		Location:    location,
		WorldTime:   worldTime,
		ActorState:  actor,
		TargetState: target,
		Seed:        seed,
	}, nil
}

package service

import (
    "context"
    "time"

    "github.com/cultivation-world/shared/errors"
    "github.com/cultivation-world/shared/types"
)

type OperationService struct {
    entityRepo EntityRepository
}

func NewOperationService(entityRepo EntityRepository) *OperationService {
    return &OperationService{entityRepo: entityRepo}
}

type EntityRepository interface {
    GetByID(ctx context.Context, id types.EntityID) (*types.Entity, error)
    Update(ctx context.Context, entity *types.Entity) error
    GetAttributes(ctx context.Context, entityID types.EntityID) (*types.Attributes, error)
    UpdateAttributes(ctx context.Context, entityID types.EntityID, attr *types.Attributes) error
}

func (s *OperationService) Execute(ctx context.Context, op *types.Operation) (*types.OperationResult, error) {
    entity, err := s.entityRepo.GetByID(ctx, op.ActorID)
    if err != nil {
        return nil, errors.ErrEntityNotFound_
    }

    if entity.Status == types.StatusDead {
        return nil, errors.NewGameError(errors.ErrInvalidOperation, "entity is dead")
    }

    switch op.ActionType {
    case types.ActionCultivate:
        return s.executeCultivate(ctx, entity, op)
    case types.ActionMove:
        return s.executeMove(ctx, entity, op)
    case types.ActionMeditate:
        return s.executeMeditate(ctx, entity, op)
    case types.ActionSleep:
        return s.executeSleep(ctx, entity, op)
    case types.ActionBreakthrough:
        return s.executeBreakthrough(ctx, entity, op)
    default:
        return nil, errors.ErrInvalidOperationType
    }
}

func (s *OperationService) executeCultivate(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
    attr, err := s.entityRepo.GetAttributes(ctx, entity.ID)
    if err != nil {
        attr = &types.Attributes{}
    }

    cultivationGain := 0.1 * float64(attr.Comprehension) / 100.0
    attr.CultivationProgress += cultivationGain

    if attr.CultivationProgress > 100 {
        attr.CultivationProgress = 100
    }

    if err := s.entityRepo.UpdateAttributes(ctx, entity.ID, attr); err != nil {
        return nil, err
    }

    entity.Status = types.StatusCultivating
    entity.UpdatedAt = time.Now()
    s.entityRepo.Update(ctx, entity)

    return &types.OperationResult{
        Success: true,
        Message: "修炼中，修为增加",
        Effects: map[string]interface{}{
            "cultivation_gain": cultivationGain,
            "progress":         attr.CultivationProgress,
        },
    }, nil
}

func (s *OperationService) executeMove(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
    regionID, ok := op.Params["region_id"].(string)
    if !ok {
        return nil, errors.ErrInvalidParams_
    }

    x, _ := op.Params["x"].(float64)
    y, _ := op.Params["y"].(float64)

    entity.Position = types.WorldPosition{
        RegionID: regionID,
        X:        x,
        Y:        y,
    }
    entity.UpdatedAt = time.Now()

    if err := s.entityRepo.Update(ctx, entity); err != nil {
        return nil, err
    }

    return &types.OperationResult{
        Success: true,
        Message: "移动成功",
        Effects: map[string]interface{}{
            "new_position": entity.Position,
        },
    }, nil
}

func (s *OperationService) executeMeditate(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
    attr, err := s.entityRepo.GetAttributes(ctx, entity.ID)
    if err != nil {
        attr = &types.Attributes{}
    }

    recovery := attr.MaxQi * 0.1
    attr.Qi += recovery
    if attr.Qi > attr.MaxQi {
        attr.Qi = attr.MaxQi
    }

    spRecovery := attr.MaxSpiritualPower * 0.1
    attr.SpiritualPower += spRecovery
    if attr.SpiritualPower > attr.MaxSpiritualPower {
        attr.SpiritualPower = attr.MaxSpiritualPower
    }

    if err := s.entityRepo.UpdateAttributes(ctx, entity.ID, attr); err != nil {
        return nil, err
    }

    entity.Status = types.StatusResting
    entity.UpdatedAt = time.Now()
    s.entityRepo.Update(ctx, entity)

    return &types.OperationResult{
        Success: true,
        Message: "打坐恢复中",
        Effects: map[string]interface{}{
            "qi_recovery":          recovery,
            "spiritual_recovery":   spRecovery,
        },
    }, nil
}

func (s *OperationService) executeSleep(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
    attr, err := s.entityRepo.GetAttributes(ctx, entity.ID)
    if err != nil {
        attr = &types.Attributes{}
    }

    attr.Qi = attr.MaxQi
    attr.SpiritualPower = attr.MaxSpiritualPower

    if err := s.entityRepo.UpdateAttributes(ctx, entity.ID, attr); err != nil {
        return nil, err
    }

    entity.Status = types.StatusNormal
    entity.UpdatedAt = time.Now()
    s.entityRepo.Update(ctx, entity)

    return &types.OperationResult{
        Success: true,
        Message: "休息完成，状态已恢复",
        Effects: map[string]interface{}{
            "qi":             attr.Qi,
            "spiritual_power": attr.SpiritualPower,
        },
    }, nil
}

func (s *OperationService) executeBreakthrough(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
    attr, err := s.entityRepo.GetAttributes(ctx, entity.ID)
    if err != nil {
        attr = &types.Attributes{}
    }

    if attr.CultivationProgress < 100 {
        return &types.OperationResult{
            Success: false,
            Message: "修为不足，无法突破",
        }, nil
    }

    successRate := 0.5 + float64(attr.Luck)/200.0
    if successRate > 0.8 {
        successRate = 0.8
    }

    newRealm := getNextRealm(entity.Realm)
    if newRealm == "" {
        return &types.OperationResult{
            Success: false,
            Message: "已达最高境界",
        }, nil
    }

    entity.Realm = newRealm
    attr.CultivationProgress = 0
    attr.MaxQi *= 1.5
    attr.MaxSpiritualPower *= 1.5
    attr.MaxLifespan = getRealmLifespan(newRealm)

    if err := s.entityRepo.UpdateAttributes(ctx, entity.ID, attr); err != nil {
        return nil, err
    }

    entity.UpdatedAt = time.Now()
    s.entityRepo.Update(ctx, entity)

    return &types.OperationResult{
        Success: true,
        Message: "突破成功！境界提升至" + string(newRealm),
        Effects: map[string]interface{}{
            "new_realm":       newRealm,
            "success_rate":    successRate,
            "max_qi":          attr.MaxQi,
            "max_spiritual":   attr.MaxSpiritualPower,
        },
    }, nil
}

func getNextRealm(current types.CultivationRealm) types.CultivationRealm {
    realms := []types.CultivationRealm{
        types.RealmMortal, types.RealmQiCondensation, types.RealmFoundation,
        types.RealmGoldenCore, types.RealmNascentSoul, types.RealmSoulTransform,
        types.RealmVoidRefinement, types.RealmIntegration, types.RealmMahayana,
        types.RealmTribulation,
    }

    for i, r := range realms {
        if r == current && i < len(realms)-1 {
            return realms[i+1]
        }
    }
    return ""
}

func getRealmLifespan(realm types.CultivationRealm) int {
    lifespans := map[types.CultivationRealm]int{
        types.RealmMortal:         80,
        types.RealmQiCondensation: 120,
        types.RealmFoundation:     200,
        types.RealmGoldenCore:     500,
        types.RealmNascentSoul:    1000,
        types.RealmSoulTransform:  3000,
        types.RealmVoidRefinement: 5000,
        types.RealmIntegration:    8000,
        types.RealmMahayana:       10000,
        types.RealmTribulation:    15000,
    }
    return lifespans[realm]
}

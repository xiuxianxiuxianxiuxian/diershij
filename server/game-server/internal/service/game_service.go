package service

import (
    "context"
    "time"

    "github.com/cultivation-world/shared/types"
    "google.golang.org/grpc"
)

type GameService struct {
    entityRepo    EntityRepository
    operationSvc  *OperationService
}

func NewGameService(entityRepo EntityRepository, operationSvc *OperationService) *GameService {
    return &GameService{
        entityRepo:   entityRepo,
        operationSvc: operationSvc,
    }
}

func (s *GameService) ExecuteOperation(ctx context.Context, req *game.OperationRequest) (*game.OperationResponse, error) {
    params := make(map[string]interface{})
    for k, v := range req.Params {
        params[k] = v
    }

    op := &types.Operation{
        ID:         req.OperationId,
        ActorID:    types.EntityID(req.ActorId),
        ActionType: types.ActionType(req.ActionType),
        Params:     params,
        Timestamp:  req.Timestamp,
    }

    result, err := s.operationSvc.Execute(ctx, op)
    if err != nil {
        return &game.OperationResponse{
            Success:   false,
            Message:   err.Error(),
            Timestamp: time.Now().UnixNano(),
        }, nil
    }

    effects := make(map[string]string)
    for k, v := range result.Effects {
        effects[k] = fmt.Sprintf("%v", v)
    }

    return &game.OperationResponse{
        Success:   result.Success,
        Message:   result.Message,
        Effects:   effects,
        Timestamp: result.Timestamp,
    }, nil
}

func (s *GameService) GetEntity(ctx context.Context, req *game.EntityRequest) (*game.EntityResponse, error) {
    entity, err := s.entityRepo.GetByID(ctx, types.EntityID(req.EntityId))
    if err != nil {
        return &game.EntityResponse{Found: false}, nil
    }

    return &game.EntityResponse{
        Entity: entityToProto(entity),
        Found:  true,
    }, nil
}

func (s *GameService) SyncState(ctx context.Context, req *game.SyncRequest) (*game.SyncResponse, error) {
    entity, err := s.entityRepo.GetByID(ctx, types.EntityID(req.EntityId))
    if err != nil {
        return nil, err
    }

    return &game.SyncResponse{
        Entity:    entityToProto(entity),
        WorldTime: time.Now().Unix(),
    }, nil
}

func (s *GameService) StreamEntityUpdates(req *game.EntityStreamRequest, stream grpc.ServerStreamingServer[game.EntityUpdate]) error {
    return nil
}

func (s *GameService) CreateEntity(ctx context.Context, req *game.CreateEntityRequest) (*game.CreateEntityResponse, error) {
    entity := &types.Entity{
        ID:         types.GenerateEntityID(),
        EntityType: types.EntityType(req.EntityType),
        Name:       req.Name,
        Realm:      types.RealmMortal,
        Position:   types.WorldPosition{RegionID: "qingyun_town", X: 0, Y: 0},
        Status:     types.StatusNormal,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    if err := s.entityRepo.Create(ctx, entity); err != nil {
        return nil, err
    }

    attr := &types.Attributes{
        Qi:                  100,
        MaxQi:               100,
        SpiritualPower:      100,
        MaxSpiritualPower:   100,
        DivineSense:         10,
        Comprehension:       50,
        Constitution:        50,
        Luck:                50,
        CultivationProgress: 0,
        AttackPower:         10,
        Defense:             10,
        Speed:               10,
        MentalStability:     50,
        RemainingLifespan:   80,
        MaxLifespan:         80,
    }

    if err := s.entityRepo.UpdateAttributes(ctx, entity.ID, attr); err != nil {
        return nil, err
    }

    return &game.CreateEntityResponse{
        Entity: entityToProto(entity),
    }, nil
}

func (s *GameService) AuthenticateEntity(ctx context.Context, req *game.AuthRequest) (*game.AuthResponse, error) {
    entity, err := s.entityRepo.GetByName(ctx, req.Username)
    if err != nil {
        return nil, err
    }

    return &game.AuthResponse{
        Entity: entityToProto(entity),
    }, nil
}

func entityToProto(e *types.Entity) *game.Entity {
    return &game.Entity{
        Id:         string(e.ID),
        EntityType: string(e.EntityType),
        Name:       e.Name,
        Realm:      string(e.Realm),
        Position: &game.WorldPosition{
            RegionId: e.Position.RegionID,
            X:        e.Position.X,
            Y:        e.Position.Y,
        },
        Attributes: &game.Attributes{
            Qi:                  e.Attributes.Qi,
            MaxQi:               e.Attributes.MaxQi,
            SpiritualPower:      e.Attributes.SpiritualPower,
            MaxSpiritualPower:   e.Attributes.MaxSpiritualPower,
            DivineSense:         e.Attributes.DivineSense,
            Comprehension:       int32(e.Attributes.Comprehension),
            Constitution:        int32(e.Attributes.Constitution),
            Luck:                int32(e.Attributes.Luck),
            CultivationProgress: e.Attributes.CultivationProgress,
            AttackPower:         e.Attributes.AttackPower,
            Defense:             e.Attributes.Defense,
            Speed:               e.Attributes.Speed,
            MentalStability:     int32(e.Attributes.MentalStability),
            RemainingLifespan:   int32(e.Attributes.RemainingLifespan),
            MaxLifespan:         int32(e.Attributes.MaxLifespan),
        },
        Karma: &game.Karma{
            KarmaValue:   int32(e.Karma.KarmaValue),
            Merit:        int32(e.Karma.Merit),
            HeavenlyMark: e.Karma.HeavenlyMark,
        },
        Status:    string(e.Status),
        CreatedAt: e.CreatedAt.Unix(),
        UpdatedAt: e.UpdatedAt.Unix(),
    }
}

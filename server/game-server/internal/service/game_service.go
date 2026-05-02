package service

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/cultivation-world/shared/proto/pb"
    "github.com/cultivation-world/shared/types"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

type GameService struct {
    pb.UnimplementedGameServiceServer
    entityRepo    EntityRepository
    operationSvc  *OperationService
    mu            sync.RWMutex
}

func NewGameService(entityRepo EntityRepository, operationSvc *OperationService) *GameService {
    return &GameService{
        entityRepo:    entityRepo,
        operationSvc: operationSvc,
    }
}

func (s *GameService) CreateEntity(ctx context.Context, req *pb.CreateEntityRequest) (*pb.CreateEntityResponse, error) {
    entityID := types.GenerateEntityID()
    now := time.Now()
    
    entity := &types.Entity{
        ID:         entityID,
        EntityType: types.EntityType(req.EntityType),
        Name:       req.Name,
        Realm:      types.RealmMortal,
        Position: types.WorldPosition{
            RegionID: "qingyun_town",
            X:        0,
            Y:        0,
        },
        Attributes: types.Attributes{
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
        },
        Karma: types.Karma{
            KarmaValue:   0,
            Merit:        0,
            HeavenlyMark: "clear",
        },
        Status:    types.StatusNormal,
        CreatedAt: now,
        UpdatedAt: now,
    }
    
    if err := s.entityRepo.Create(ctx, entity); err != nil {
        return nil, status.Errorf(codes.Internal, "failed to create entity: %v", err)
    }
    
    if err := s.entityRepo.UpdateAttributes(ctx, entityID, &entity.Attributes); err != nil {
        return nil, status.Errorf(codes.Internal, "failed to save attributes: %v", err)
    }
    
    return &pb.CreateEntityResponse{
        Entity: entityToProto(entity),
    }, nil
}

func (s *GameService) AuthenticateEntity(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
    entity, err := s.entityRepo.GetByName(ctx, req.Username)
    if err != nil {
        return &pb.AuthResponse{
            Success: false,
        }, nil
    }
    
    return &pb.AuthResponse{
        Entity:  entityToProto(entity),
        Success: true,
    }, nil
}

func (s *GameService) ExecuteOperation(ctx context.Context, req *pb.OperationRequest) (*pb.OperationResponse, error) {
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
        return &pb.OperationResponse{
            Success: false,
            Message: err.Error(),
            Timestamp: time.Now().UnixNano(),
        }, nil
    }
    
    effects := make(map[string]string)
    for k, v := range result.Effects {
        effects[k] = fmt.Sprintf("%v", v)
    }
    
    return &pb.OperationResponse{
        Success: result.Success,
        Message: result.Message,
        Effects: effects,
        Timestamp: result.Timestamp,
    }, nil
}

func (s *GameService) GetEntity(ctx context.Context, req *pb.EntityRequest) (*pb.EntityResponse, error) {
    entity, err := s.entityRepo.GetByID(ctx, types.EntityID(req.EntityId))
    if err != nil {
        return &pb.EntityResponse{
            Found: false,
        }, nil
    }
    
    return &pb.EntityResponse{
        Entity: entityToProto(entity),
        Found:  true,
    }, nil
}

func (s *GameService) SyncState(ctx context.Context, req *pb.SyncRequest) (*pb.SyncResponse, error) {
    entity, err := s.entityRepo.GetByID(ctx, types.EntityID(req.EntityId))
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "entity not found")
    }
    
    return &pb.SyncResponse{
        Entity: entityToProto(entity),
        WorldTime: time.Now().Unix(),
        NearbyEntityIds: []string{},
    }, nil
}

func (s *GameService) StreamEntityUpdates(req *pb.EntityStreamRequest, stream grpc.ServerStreamingServer[pb.EntityUpdate]) error {
    return nil
}

func entityToProto(e *types.Entity) *pb.Entity {
    if e == nil {
        return nil
    }
    
    return &pb.Entity{
        Id:         string(e.ID),
        EntityType: string(e.EntityType),
        Name:       e.Name,
        Realm:      string(e.Realm),
        Position: &pb.WorldPosition{
            RegionId: e.Position.RegionID,
            X:        e.Position.X,
            Y:        e.Position.Y,
        },
        Attributes: &pb.Attributes{
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
        Karma: &pb.Karma{
            KarmaValue:   int32(e.Karma.KarmaValue),
            Merit:        int32(e.Karma.Merit),
            HeavenlyMark: e.Karma.HeavenlyMark,
        },
        Status:    string(e.Status),
        CreatedAt: e.CreatedAt.Unix(),
        UpdatedAt: e.UpdatedAt.Unix(),
    }
}

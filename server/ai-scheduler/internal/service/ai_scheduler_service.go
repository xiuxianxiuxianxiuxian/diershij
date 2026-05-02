package service

import (
    "context"
    "fmt"
    "math/rand"
    "sync"
    "time"

    "github.com/cultivation-world/shared/proto/pb"
    "github.com/cultivation-world/shared/types"
    "google.golang.org/grpc"
)

type AISchedulerService struct {
    pb.UnimplementedAISchedulerServiceServer
    entityActions map[types.EntityID]*EntityState
    mu          sync.RWMutex
    scheduler   *Scheduler
}

type EntityState struct {
    LastAction time.Time
    Behavior   BehaviorType
}

type BehaviorType string

const (
    BehaviorIdle     BehaviorType = "idle"
    BehaviorCultivate BehaviorType = "cultivate"
    BehaviorExplore  BehaviorType = "explore"
)

type Scheduler struct {
    mu sync.RWMutex
}

func NewScheduler() *Scheduler {
    return &Scheduler{}
}

func (s *Scheduler) GenerateAction(ctx context.Context, entityID types.EntityID, behavior types.EntityType) (*types.AIAction, error) {
    actions := []types.ActionType{
        types.ActionCultivate,
        types.ActionMove,
        types.ActionExplore,
    }
    
    idx := rand.Intn(len(actions))
    
    return &types.AIAction{
        ID:         types.GenerateOperationID(),
        ActorID:    entityID,
        ActionType: actions[idx],
        Params: map[string]interface{}{
            "target": "somewhere",
        },
        Priority: rand.Intn(10),
    }, nil
}

func (s *Scheduler) ScheduleDecision(ctx context.Context, world *types.WorldState) ([]*types.AIAction, error) {
    return []*types.AIAction{}, nil
}

func NewAISchedulerService() *AISchedulerService {
    return &AISchedulerService{
        entityActions: make(map[types.EntityID]*EntityState),
        scheduler:     NewScheduler(),
    }
}

func (s *AISchedulerService) GetAIAction(ctx context.Context, req *pb.AIActionRequest) (*pb.AIActionResponse, error) {
    action, err := s.scheduler.GenerateAction(ctx, types.EntityID(req.EntityId), types.EntityType(req.EntityType))
    if err != nil {
        return nil, err
    }
    
    params := make(map[string]string)
    for k, v := range action.Params {
        params[k] = fmt.Sprintf("%v", v)
    }
    
    return &pb.AIActionResponse{
        ActionId:   string(action.ID),
        ActionType: string(action.ActionType),
        Params: params,
        Priority: int32(action.Priority),
    }, nil
}

func (s *AISchedulerService) SetEntityBehavior(ctx context.Context, req *pb.BehaviorRequest) (*pb.BehaviorResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.entityActions[types.EntityID(req.EntityId)] = &EntityState{
        LastAction: time.Now(),
        Behavior:   BehaviorType(req.BehaviorType),
    }
    
    return &pb.BehaviorResponse{
        Success: true,
    }, nil
}

func (s *AISchedulerService) StreamAIActions(req *pb.AIStreamRequest, stream grpc.ServerStreamingServer[pb.AIAction]) error {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            action := &pb.AIAction{
                ActionId: "auto-" + time.Now().Format("20060102150405"),
                EntityId: req.EntityId,
                ActionType: string(types.ActionCultivate),
            }
            if err := stream.Send(action); err != nil {
                return err
            }
        case <-stream.Context().Done():
            return stream.Context().Err()
        }
    }
}

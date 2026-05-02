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

type HeavenlyDaoService struct {
    pb.UnimplementedHeavenlyDaoServiceServer
    worldState     *types.WorldState
    karmaEngine    *KarmaEngine
    tribulationMgr *TribulationManager
    lawManager     *LawManager
    mu             sync.RWMutex
}

type KarmaEngine struct{}

func (e *KarmaEngine) UpdateKarma(ctx context.Context, entityID types.EntityID, action types.ActionType) error {
    return nil
}

type TribulationManager struct{}

func (m *TribulationManager) TriggerTribulation(ctx context.Context, entityID types.EntityID) (*types.TribulationResult, error) {
    return &types.TribulationResult{
        Success: true,
        Message: "Tribulation passed",
        NewRealm: types.RealmQiRefining,
    }, nil
}

func (m *TribulationManager) CheckTribulation(ctx context.Context, entityID types.EntityID) (bool, error) {
    return false, nil
}

type LawManager struct{}

func (m *LawManager) EnforceLaws(ctx context.Context, state *types.WorldState) error {
    return nil
}

func NewHeavenlyDaoService() *HeavenlyDaoService {
    return &HeavenlyDaoService{
        worldState: &types.WorldState{
            WorldTime: 0,
            Cycle:     1,
            Rules:     make(map[string]interface{}),
        },
        karmaEngine:    &KarmaEngine{},
        tribulationMgr: &TribulationManager{},
        lawManager:     &LawManager{},
    }
}

func (s *HeavenlyDaoService) UpdateEntityKarma(ctx context.Context, req *pb.KarmaUpdateRequest) (*pb.KarmaUpdateResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.karmaEngine.UpdateKarma(ctx, types.EntityID(req.EntityId), types.ActionType(req.ActionType))
    
    return &pb.KarmaUpdateResponse{
        KarmaValue: int32(req.KarmaChange),
    }, nil
}

func (s *HeavenlyDaoService) TriggerTribulation(ctx context.Context, req *pb.TribulationRequest) (*pb.TribulationResponse, error) {
    result, err := s.tribulationMgr.TriggerTribulation(ctx, types.EntityID(req.EntityId))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to trigger tribulation: %v", err)
    }
    
    return &pb.TribulationResponse{
        Success: result.Success,
        Message: result.Message,
        NewRealm: string(result.NewRealm),
    }, nil
}

func (s *HeavenlyDaoService) GetHeavenlyState(ctx context.Context, req *pb.HeavenlyStateRequest) (*pb.HeavenlyStateResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    return &pb.HeavenlyStateResponse{
        WorldTime: s.worldState.WorldTime,
        Cycle:     int32(s.worldState.Cycle),
        HeavenlyAura: 100.0,
    }, nil
}

func (s *HeavenlyDaoService) WatchHeavenlyEvents(req *pb.HeavenlyEventRequest, stream grpc.ServerStreamingServer[pb.HeavenlyEvent]) error {
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()
    
    eventCount := 0
    
    for {
        select {
        case <-ticker.C:
            eventCount++
            event := &pb.HeavenlyEvent{
                EventId:   fmt.Sprintf("event-%d", eventCount),
                EventType: "tick",
                Timestamp: time.Now().Unix(),
            }
            if err := stream.Send(event); err != nil {
                return err
            }
        case <-stream.Context().Done():
            return stream.Context().Err()
        }
    }
}

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

type WorldService struct {
    pb.UnimplementedWorldServiceServer
    worldState *types.WorldState
    regions    map[string]*types.Region
    resources  map[string]*types.Resource
    mu         sync.RWMutex
}

func NewWorldService() *WorldService {
    regions := make(map[string]*types.Region)
    
    regions["qingyun_town"] = &types.Region{
        ID:          "qingyun_town",
        Name:        "青云镇",
        Description: "修仙新手起点",
        Type:        "town",
        DangerLevel: 1,
        SpiritualRichness: 10.0,
    }
    
    return &WorldService{
        worldState: &types.WorldState{
            WorldTime: time.Now().Unix(),
            Cycle:     1,
            Rules:     make(map[string]interface{}),
        },
        regions: regions,
        resources: make(map[string]*types.Resource),
    }
}

func (s *WorldService) GetWorldState(ctx context.Context, req *pb.WorldStateRequest) (*pb.WorldStateResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    regions := make([]*pb.Region, 0, len(s.regions))
    for _, r := range s.regions {
        regions = append(regions, &pb.Region{
            Id:          r.ID,
            Name:        r.Name,
            Description: r.Description,
            Type:        r.Type,
            DangerLevel: int32(r.DangerLevel),
        })
    }
    
    return &pb.WorldStateResponse{
        WorldTime: s.worldState.WorldTime,
        Cycle:     int32(s.worldState.Cycle),
        Regions:   regions,
    }, nil
}

func (s *WorldService) GetRegion(ctx context.Context, req *pb.RegionRequest) (*pb.RegionResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    region, exists := s.regions[req.RegionId]
    if !exists {
        return nil, status.Errorf(codes.NotFound, "region not found")
    }
    
    return &pb.RegionResponse{
        Region: &pb.Region{
            Id:          region.ID,
            Name:        region.Name,
            Description: region.Description,
            Type:        region.Type,
            DangerLevel: int32(region.DangerLevel),
        },
    }, nil
}

func (s *WorldService) StreamWorldEvents(req *pb.WorldEventRequest, stream grpc.ServerStreamingServer[pb.WorldEvent]) error {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    eventCount := 0
    
    for {
        select {
        case <-ticker.C:
            eventCount++
            event := &pb.WorldEvent{
                EventId:   fmt.Sprintf("world-%d", eventCount),
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

func (s *WorldService) ModifyWorld(ctx context.Context, req *pb.ModifyWorldRequest) (*pb.ModifyWorldResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    s.worldState.WorldTime = time.Now().Unix()
    
    return &pb.ModifyWorldResponse{
        Success: true,
    }, nil
}

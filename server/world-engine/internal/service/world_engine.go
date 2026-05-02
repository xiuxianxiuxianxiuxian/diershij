package service

import (
    "context"
    "sync"
    "time"

    "github.com/cultivation-world/shared/types"
    "github.com/google/uuid"
)

type WorldEngineService struct {
    mu            sync.RWMutex
    regions       map[string]*types.Region
    worldEvents   []*types.WorldEvent
    worldState    *types.WorldState
    epoch         int64
}

func NewWorldEngineService() *WorldEngineService {
    svc := &WorldEngineService{
        regions:     make(map[string]*types.Region),
        worldEvents: make([]*types.WorldEvent, 0),
        epoch:       0,
    }

    svc.initializeWorld()

    return svc
}

func (s *WorldEngineService) initializeWorld() {
    s.regions = map[string]*types.Region{
        "east_wilderness": {
            ID:               "east_wilderness",
            Name:             "东荒域",
            SpiritualDensity: 50,
            SpiritualTier:    3,
            DangerLevel:      2,
            Description:      "东荒大地，灵气充沛，是新手修士的聚集地",
            Resources: []types.Resource{
                {ID: "res_1", Name: "灵草", Type: "herb", Rarity: 1, Quantity: 100, RespawnRate: 0.1},
                {ID: "res_2", Name: "灵石矿", Type: "ore", Rarity: 2, Quantity: 50, RespawnRate: 0.05},
            },
        },
        "qingyun_town": {
            ID:               "qingyun_town",
            Name:             "青云镇",
            ParentRegionID:   ptr("east_wilderness"),
            SpiritualDensity: 30,
            SpiritualTier:    1,
            DangerLevel:      0,
            Description:      "凡人城镇，修士的起点",
            Resources:        []types.Resource{},
        },
        "spirit_mist_mountain": {
            ID:               "spirit_mist_mountain",
            Name:             "灵雾山脉",
            ParentRegionID:   ptr("east_wilderness"),
            SpiritualDensity: 60,
            SpiritualTier:    4,
            DangerLevel:      3,
            Description:      "灵气浓郁的山脉，适合修炼",
            Resources: []types.Resource{
                {ID: "res_3", Name: "千年灵草", Type: "herb", Rarity: 3, Quantity: 20, RespawnRate: 0.02},
            },
        },
        "south_ridge": {
            ID:               "south_ridge",
            Name:             "南岭域",
            SpiritualDensity: 70,
            SpiritualTier:    5,
            DangerLevel:      4,
            Description:      "南岭大地，火属性灵气浓郁",
            Resources: []types.Resource{
                {ID: "res_4", Name: "火灵石", Type: "ore", Rarity: 4, Quantity: 30, RespawnRate: 0.03},
            },
        },
        "central_state": {
            ID:               "central_state",
            Name:             "中州域",
            SpiritualDensity: 90,
            SpiritualTier:    8,
            DangerLevel:      6,
            Description:      "世界中心，灵气最浓郁之地",
            Resources: []types.Resource{
                {ID: "res_5", Name: "天材地宝", Type: "treasure", Rarity: 5, Quantity: 10, RespawnRate: 0.01},
            },
        },
    }

    s.worldState = &types.WorldState{
        Epoch:        0,
        Regions:      make(map[types.RegionID]types.Region),
        ActiveEvents: []types.WorldEvent{},
        BalanceMetrics: types.BalanceMetrics{
            PowerDistribution:   0.5,
            ResourceCirculation: 0.3,
            SectDiversity:       0.4,
            KarmaDistribution:   0.5,
        },
        LastUpdated: time.Now(),
    }

    for id, region := range s.regions {
        s.worldState.Regions[types.RegionID(id)] = *region
    }
}

func (s *WorldEngineService) GetRegion(ctx context.Context, req *game.RegionRequest) (*game.RegionResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    region, exists := s.regions[req.RegionId]
    if !exists {
        return &game.RegionResponse{Found: false}, nil
    }

    return &game.RegionResponse{
        Region: regionToProto(region),
        Found:  true,
    }, nil
}

func (s *WorldEngineService) SpawnResources(ctx context.Context, req *game.SpawnRequest) (*game.SpawnResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    region, exists := s.regions[req.RegionId]
    if !exists {
        return &game.SpawnResponse{Success: false}, nil
    }

    spawned := make([]*game.Resource, 0, req.Quantity)
    for i := 0; i < int(req.Quantity); i++ {
        resource := &types.Resource{
            ID:          uuid.New().String(),
            Name:        req.ResourceType,
            Type:        req.ResourceType,
            Rarity:      1,
            Quantity:    1,
            RespawnRate: 0.1,
        }
        region.Resources = append(region.Resources, *resource)
        spawned = append(spawned, resourceToProto(resource))
    }

    s.regions[req.RegionId] = region

    return &game.SpawnResponse{
        Success: true,
        Spawned: spawned,
    }, nil
}

func (s *WorldEngineService) TriggerEvent(ctx context.Context, req *game.EventRequest) (*game.EventResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    event := &types.WorldEvent{
        ID:           uuid.New().String(),
        Name:         req.EventType,
        Type:         req.EventType,
        Description:  "World event triggered",
        RegionID:     types.RegionID(req.RegionId),
        StartTime:    time.Now(),
        Participants: []types.EntityID{},
        Status:       "active",
    }

    s.worldEvents = append(s.worldEvents, event)

    return &game.EventResponse{
        Success: true,
        Event:   worldEventToProto(event),
    }, nil
}

func (s *WorldEngineService) GetWorldState(ctx context.Context, req *game.WorldStateRequest) (*game.WorldStateResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    return &game.WorldStateResponse{
        State: worldStateToProto(s.worldState),
    }, nil
}

func (s *WorldEngineService) AdvanceEpoch() {
    s.mu.Lock()
    defer s.mu.Unlock()

    s.epoch++
    s.worldState.Epoch = s.epoch
    s.worldState.LastUpdated = time.Now()

    for id, region := range s.regions {
        for i, res := range region.Resources {
            if res.Quantity < 100 {
                region.Resources[i].Quantity += int(res.RespawnRate * 100)
            }
        }
        s.regions[string(id)] = region
        s.worldState.Regions[types.RegionID(id)] = *region
    }
}

func ptr(s string) *string {
    return &s
}

func regionToProto(r *types.Region) *game.Region {
    resources := make([]*game.Resource, 0, len(r.Resources))
    for _, res := range r.Resources {
        resources = append(resources, resourceToProto(&res))
    }

    var parentID string
    if r.ParentRegionID != nil {
        parentID = *r.ParentRegionID
    }

    return &game.Region{
        Id:               r.ID,
        Name:             r.Name,
        ParentRegionId:   parentID,
        SpiritualDensity: r.SpiritualDensity,
        SpiritualTier:    int32(r.SpiritualTier),
        DangerLevel:      int32(r.DangerLevel),
        Resources:        resources,
        Description:      r.Description,
        Lore:             r.Lore,
    }
}

func resourceToProto(r *types.Resource) *game.Resource {
    var lastHarvested int64
    if r.LastHarvested != nil {
        lastHarvested = r.LastHarvested.Unix()
    }

    return &game.Resource{
        Id:            r.ID,
        Name:          r.Name,
        Type:          r.Type,
        Rarity:        int32(r.Rarity),
        Quantity:      int32(r.Quantity),
        RespawnRate:   r.RespawnRate,
        LastHarvested: lastHarvested,
    }
}

func worldEventToProto(e *types.WorldEvent) *game.WorldEvent {
    participants := make([]string, 0, len(e.Participants))
    for _, p := range e.Participants {
        participants = append(participants, string(p))
    }

    var endTime int64
    if e.EndTime != nil {
        endTime = e.EndTime.Unix()
    }

    return &game.WorldEvent{
        Id:            e.ID,
        Name:          e.Name,
        Type:          e.Type,
        Description:   e.Description,
        RegionId:      string(e.RegionID),
        StartTime:     e.StartTime.Unix(),
        EndTime:       endTime,
        ParticipantIds: participants,
        Status:        e.Status,
    }
}

func worldStateToProto(s *types.WorldState) *game.WorldState {
    regions := make(map[string]*game.Region)
    for id, r := range s.Regions {
        region := r
        regions[string(id)] = regionToProto(&region)
    }

    events := make([]*game.WorldEvent, 0, len(s.ActiveEvents))
    for _, e := range s.ActiveEvents {
        events = append(events, worldEventToProto(&e))
    }

    return &game.WorldState{
        Epoch:    s.Epoch,
        Regions:  regions,
        ActiveEvents: events,
        BalanceMetrics: &game.BalanceMetrics{
            PowerDistribution:    s.BalanceMetrics.PowerDistribution,
            ResourceCirculation:  s.BalanceMetrics.ResourceCirculation,
            SectDiversity:        s.BalanceMetrics.SectDiversity,
            KarmaDistribution:    s.BalanceMetrics.KarmaDistribution,
        },
        LastUpdated: s.LastUpdated.Unix(),
    }
}

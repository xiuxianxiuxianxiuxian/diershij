package service

import (
	"context"
	"math/rand"
	"sync"
	"time"

	cultivation "github.com/cultivation-world/shared/proto/pb"
	"github.com/cultivation-world/shared/types"
	"github.com/google/uuid"
)

// ── 世界事件类型 ──

const (
	EventTypeBeastTide    = "beast_tide"     // 妖兽潮
	EventTypeSecretRealm  = "secret_realm"    // 秘境开启
	EventTypeTreasure     = "treasure_appear" // 天材地宝出世
	EventTypeTribulation  = "tribulation_omen" // 天劫异象
)

// EventEffects 定义事件对区域的具体影响
type EventEffects struct {
	SpiritualDensityBonus float64 // 灵气密度加成
	DangerLevelBonus      int     // 危险等级加成
	CultivationMultiplier float64 // 修炼效率倍率（1.0 = 正常）
	CombatProbability     float64 // 战斗概率增加（0-1）
	ExplorationBonus      float64 // 探索收益加成
	ResourceRespawnBonus  float64 // 资源刷新倍率
	TribulationModifier   float64 // 天劫概率修正
	SpecialResource       string  // 特产资源名称
	SpecialResourceRate   float64 // 特产出现概率
}

// EventDefinition 事件模板定义
type EventDefinition struct {
	Type          string
	Name          string
	Description   string
	Duration      time.Duration
	Cooldown      time.Duration
	TriggerChance float64       // 每次调度 tick 触发概率
	MinTier       int           // 最低区域灵气等级
	Intensity     float64       // 事件强度 0-1
	Effects       EventEffects
}

// ActiveEvent 活跃事件实例
type ActiveEvent struct {
	ID          string
	Type        string
	Name        string
	Description string
	RegionID    types.RegionID
	StartTime   time.Time
	EndTime     time.Time
	Intensity   float64
	Effects     EventEffects
	Status      string // "active", "ending"
}

// EventCooldown 区域事件冷却
type EventCooldown struct {
	EventType string
	RegionID  types.RegionID
	EndTime   time.Time
}

// 预定义事件模板
var eventTemplates = map[string]EventDefinition{
	EventTypeBeastTide: {
		Type:          EventTypeBeastTide,
		Name:          "妖兽潮",
		Description:   "大量妖兽聚集，区域危险度上升，但击杀收益增加",
		Duration:      5 * time.Minute,
		Cooldown:      15 * time.Minute,
		TriggerChance: 0.08,
		MinTier:       2,
		Intensity:     1.0,
		Effects: EventEffects{
			DangerLevelBonus:      2,
			CombatProbability:     0.3,
			CultivationMultiplier: 0.7,  // 妖兽潮期间修炼效率降低
			ExplorationBonus:      0.0,
			ResourceRespawnBonus:  0.0,
		},
	},
	EventTypeSecretRealm: {
		Type:          EventTypeSecretRealm,
		Name:          "秘境开启",
		Description:   "古老秘境现世，探索收益大幅提升，可获得稀有物资",
		Duration:      8 * time.Minute,
		Cooldown:      20 * time.Minute,
		TriggerChance: 0.05,
		MinTier:       3,
		Intensity:     1.0,
		Effects: EventEffects{
			SpiritualDensityBonus: 20,
			ExplorationBonus:      0.5,
			CultivationMultiplier: 1.3,
			SpecialResource:       "秘境灵石",
			SpecialResourceRate:   0.15,
		},
	},
	EventTypeTreasure: {
		Type:          EventTypeTreasure,
		Name:          "天材地宝出世",
		Description:   "天地异象，珍稀宝物在区域中出现，采集收益大增",
		Duration:      6 * time.Minute,
		Cooldown:      25 * time.Minute,
		TriggerChance: 0.04,
		MinTier:       4,
		Intensity:     1.0,
		Effects: EventEffects{
			SpiritualDensityBonus: 15,
			ResourceRespawnBonus:  0.5,
			ExplorationBonus:      0.3,
			CultivationMultiplier: 1.1,
			SpecialResource:       "天材地宝",
			SpecialResourceRate:   0.1,
		},
	},
	EventTypeTribulation: {
		Type:          EventTypeTribulation,
		Name:          "天劫异象",
		Description:   "天道波动，突破时天劫概率增加，但成功后的收益也更大",
		Duration:      10 * time.Minute,
		Cooldown:      30 * time.Minute,
		TriggerChance: 0.03,
		MinTier:       5,
		Intensity:     1.0,
		Effects: EventEffects{
			SpiritualDensityBonus: 30,
			CultivationMultiplier: 1.5,
			TribulationModifier:   0.2,
		},
	},
}

type WorldEngineService struct {
	cultivation.UnimplementedWorldServiceServer
	mu            sync.RWMutex
	regions       map[string]*types.Region
	activeEvents  []*ActiveEvent
	cooldowns     []EventCooldown
	worldState    *types.WorldState
	epoch         int64
	notifyChan    chan *ActiveEvent // 新触发事件通知通道
}

func NewWorldEngineService() *WorldEngineService {
	svc := &WorldEngineService{
		regions:      make(map[string]*types.Region),
		activeEvents: make([]*ActiveEvent, 0),
		cooldowns:    make([]EventCooldown, 0),
		epoch:        0,
		notifyChan:   make(chan *ActiveEvent, 64),
	}

	svc.initializeWorld()

	return svc
}

// GetNotifyChan 返回事件通知通道（供 scheduler 使用）
func (s *WorldEngineService) GetNotifyChan() <-chan *ActiveEvent {
	return s.notifyChan
}

func ptrRegionID(s string) *types.RegionID {
	rid := types.RegionID(s)
	return &rid
}

func (s *WorldEngineService) initializeWorld() {
	s.regions = map[string]*types.Region{
		"east_wilderness": {
			ID:               types.RegionID("east_wilderness"),
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
			ID:               types.RegionID("qingyun_town"),
			Name:             "青云镇",
			ParentRegionID:   ptrRegionID("east_wilderness"),
			SpiritualDensity: 30,
			SpiritualTier:    1,
			DangerLevel:      0,
			Description:      "凡人城镇，修士的起点",
			Resources:        []types.Resource{},
		},
		"spirit_mist_mountain": {
			ID:               types.RegionID("spirit_mist_mountain"),
			Name:             "灵雾山脉",
			ParentRegionID:   ptrRegionID("east_wilderness"),
			SpiritualDensity: 60,
			SpiritualTier:    4,
			DangerLevel:      3,
			Description:      "灵气浓郁的山脉，适合修炼",
			Resources: []types.Resource{
				{ID: "res_3", Name: "千年灵草", Type: "herb", Rarity: 3, Quantity: 20, RespawnRate: 0.02},
			},
		},
		"south_ridge": {
			ID:               types.RegionID("south_ridge"),
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
			ID:               types.RegionID("central_state"),
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

// ── 事件系统核心方法 ──

// TryTriggerEvents 检查所有区域的事件触发条件，返回新触发的事件列表
func (s *WorldEngineService) TryTriggerEvents() []*ActiveEvent {
	s.mu.Lock()
	defer s.mu.Unlock()

	var newEvents []*ActiveEvent
	now := time.Now()

	for regionID, region := range s.regions {
		// 跳过已有活跃事件的区域
		if s.hasActiveEventLocked(types.RegionID(regionID)) {
			continue
		}

		for _, tmpl := range eventTemplates {
			// 检查区域等级要求
			if region.SpiritualTier < tmpl.MinTier {
				continue
			}

			// 检查冷却
			if s.isOnCooldownLocked(tmpl.Type, types.RegionID(regionID), now) {
				continue
			}

			// 概率触发
			if rand.Float64() >= tmpl.TriggerChance {
				continue
			}

			event := &ActiveEvent{
				ID:          uuid.New().String(),
				Type:        tmpl.Type,
				Name:        tmpl.Name,
				Description: tmpl.Description,
				RegionID:    types.RegionID(regionID),
				StartTime:   now,
				EndTime:     now.Add(tmpl.Duration),
				Intensity:   tmpl.Intensity,
				Effects:     tmpl.Effects,
				Status:      "active",
			}

			s.activeEvents = append(s.activeEvents, event)
			newEvents = append(newEvents, event)

			// 放入通知通道（非阻塞）
			select {
			case s.notifyChan <- event:
			default:
			}

			// 每区域每次 tick 只触发一个事件
			break
		}
	}

	return newEvents
}

// UpdateActiveEvents 更新活跃事件状态，清理过期事件并设为冷却
func (s *WorldEngineService) UpdateActiveEvents() []*ActiveEvent {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	var expired []*ActiveEvent
	var remaining []*ActiveEvent

	for _, event := range s.activeEvents {
		if now.After(event.EndTime) {
			event.Status = "ended"
			expired = append(expired, event)

			// 添加冷却
			if tmpl, ok := eventTemplates[event.Type]; ok {
				s.cooldowns = append(s.cooldowns, EventCooldown{
					EventType: event.Type,
					RegionID:  event.RegionID,
					EndTime:   now.Add(tmpl.Cooldown),
				})
			}
		} else {
			remaining = append(remaining, event)
		}
	}

	s.activeEvents = remaining

	// 清理过期冷却
	var activeCooldowns []EventCooldown
	for _, cd := range s.cooldowns {
		if now.Before(cd.EndTime) {
			activeCooldowns = append(activeCooldowns, cd)
		}
	}
	s.cooldowns = activeCooldowns

	return expired
}

func (s *WorldEngineService) hasActiveEventLocked(regionID types.RegionID) bool {
	for _, e := range s.activeEvents {
		if e.RegionID == regionID && e.Status == "active" {
			return true
		}
	}
	return false
}

func (s *WorldEngineService) isOnCooldownLocked(eventType string, regionID types.RegionID, now time.Time) bool {
	for _, cd := range s.cooldowns {
		if cd.EventType == eventType && cd.RegionID == regionID && now.Before(cd.EndTime) {
			return true
		}
	}
	return false
}

// GetActiveEventsForRegion 获取指定区域的活跃事件
func (s *WorldEngineService) GetActiveEventsForRegion(regionID types.RegionID) []*ActiveEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*ActiveEvent
	for _, e := range s.activeEvents {
		if e.RegionID == regionID && e.Status == "active" {
			result = append(result, e)
		}
	}
	return result
}

// GetRegionEffects 获取指定区域所有活跃事件的叠加效果
func (s *WorldEngineService) GetRegionEffects(regionID types.RegionID) EventEffects {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var combined EventEffects
	for _, e := range s.activeEvents {
		if e.RegionID == regionID && e.Status == "active" {
			combined.SpiritualDensityBonus += e.Effects.SpiritualDensityBonus
			combined.DangerLevelBonus += e.Effects.DangerLevelBonus
			if e.Effects.CultivationMultiplier > combined.CultivationMultiplier {
				combined.CultivationMultiplier = e.Effects.CultivationMultiplier
			}
			if e.Effects.CombatProbability > combined.CombatProbability {
				combined.CombatProbability = e.Effects.CombatProbability
			}
			if e.Effects.ExplorationBonus > combined.ExplorationBonus {
				combined.ExplorationBonus = e.Effects.ExplorationBonus
			}
			if e.Effects.ResourceRespawnBonus > combined.ResourceRespawnBonus {
				combined.ResourceRespawnBonus = e.Effects.ResourceRespawnBonus
			}
			if e.Effects.TribulationModifier > combined.TribulationModifier {
				combined.TribulationModifier = e.Effects.TribulationModifier
			}
			if e.Effects.SpecialResource != "" {
				combined.SpecialResource = e.Effects.SpecialResource
				combined.SpecialResourceRate = e.Effects.SpecialResourceRate
			}
		}
	}
	return combined
}

// GetAllActiveEvents 返回所有活跃事件（用于 RPC）
func (s *WorldEngineService) GetAllActiveEvents() []*ActiveEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*ActiveEvent, len(s.activeEvents))
	copy(result, s.activeEvents)
	return result
}

// ── gRPC 接口实现 ──

func (s *WorldEngineService) GetRegion(ctx context.Context, req *cultivation.RegionRequest) (*cultivation.RegionResponse, error) {
	s.mu.RLock()
	region, exists := s.regions[req.RegionId]
	if !exists {
		s.mu.RUnlock()
		return &cultivation.RegionResponse{Found: false}, nil
	}

	// 复制区域并应用事件效果
	modified := *region
	for _, e := range s.activeEvents {
		if e.RegionID == region.ID && e.Status == "active" {
			modified.SpiritualDensity += e.Effects.SpiritualDensityBonus
			modified.DangerLevel += e.Effects.DangerLevelBonus
		}
	}
	s.mu.RUnlock()

	return &cultivation.RegionResponse{
		Region: regionToProto(&modified),
		Found:  true,
	}, nil
}

func (s *WorldEngineService) SpawnResources(ctx context.Context, req *cultivation.SpawnRequest) (*cultivation.SpawnResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	region, exists := s.regions[req.RegionId]
	if !exists {
		return &cultivation.SpawnResponse{Success: false}, nil
	}

	spawned := make([]*cultivation.Resource, 0, req.Quantity)
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

	return &cultivation.SpawnResponse{
		Success: true,
		Spawned: spawned,
	}, nil
}

func (s *WorldEngineService) TriggerEvent(ctx context.Context, req *cultivation.EventRequest) (*cultivation.EventResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 查找事件模板
	tmpl, ok := eventTemplates[req.EventType]
	if !ok {
		// 未知事件类型，使用默认值
		tmpl = EventDefinition{
			Duration: 5 * time.Minute,
		}
	}

	now := time.Now()
	event := &ActiveEvent{
		ID:          uuid.New().String(),
		Type:        req.EventType,
		Name:        tmpl.Name,
		Description: tmpl.Description,
		RegionID:    types.RegionID(req.RegionId),
		StartTime:   now,
		EndTime:     now.Add(tmpl.Duration),
		Intensity:   tmpl.Intensity,
		Effects:     tmpl.Effects,
		Status:      "active",
	}

	// 如果名称未定义，使用请求中的 type
	if event.Name == "" {
		event.Name = req.EventType
	}
	if event.Description == "" {
		event.Description = req.Params["description"]
	}

	s.activeEvents = append(s.activeEvents, event)

	// 尝试通知
	select {
	case s.notifyChan <- event:
	default:
	}

	// 返回 proto 格式事件
	protoEvent := &cultivation.WorldEvent{
		Id:          event.ID,
		Name:        event.Name,
		Type:        event.Type,
		Description: event.Description,
		RegionId:    string(event.RegionID),
		StartTime:   event.StartTime.Unix(),
		EndTime:     event.EndTime.Unix(),
		Status:      event.Status,
	}

	return &cultivation.EventResponse{
		Success: true,
		Event:   protoEvent,
	}, nil
}

func (s *WorldEngineService) GetWorldState(ctx context.Context, req *cultivation.WorldStateRequest) (*cultivation.WorldStateResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &cultivation.WorldStateResponse{
		State: worldStateToProto(s.worldState, s.activeEvents),
	}, nil
}

func (s *WorldEngineService) AdvanceEpoch() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.epoch++
	s.worldState.Epoch = s.epoch
	s.worldState.LastUpdated = time.Now()

	for id, region := range s.regions {
		respawnMultiplier := 1.0
		for _, e := range s.activeEvents {
			if e.RegionID == region.ID && e.Status == "active" {
				if e.Effects.ResourceRespawnBonus > 0 {
					respawnMultiplier = 1.0 + e.Effects.ResourceRespawnBonus
				}
			}
		}

		for i, res := range region.Resources {
			regenerated := int(res.RespawnRate * 100 * respawnMultiplier)
			if regenerated > 0 && res.Quantity < 100 {
				region.Resources[i].Quantity += regenerated
				if region.Resources[i].Quantity > 100 {
					region.Resources[i].Quantity = 100
				}
			}
		}
		s.regions[id] = region
		s.worldState.Regions[types.RegionID(id)] = *region
	}
}

// ── 持久化支持方法 ──

// WorldPersistenceState 持久化所需的世界状态快照
type WorldPersistenceState struct {
	Epoch   int64
	Metrics types.BalanceMetrics
}

// RestoreState 从数据库恢复世界状态
func (s *WorldEngineService) RestoreState(epoch int64, metrics *types.BalanceMetrics) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.epoch = epoch
	s.worldState.Epoch = epoch
	if metrics != nil {
		s.worldState.BalanceMetrics = *metrics
	}
}

// RestoreRegionResources 从数据库恢复区域资源数量
func (s *WorldEngineService) RestoreRegionResources(regionID types.RegionID, resources map[string]int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if region, ok := s.regions[string(regionID)]; ok {
		for i, res := range region.Resources {
			if qty, found := resources[res.ID]; found {
				region.Resources[i].Quantity = qty
			}
		}
		s.regions[string(regionID)] = region
		s.worldState.Regions[regionID] = *region
	}
}

// GetAllRegions 返回所有区域的切片副本
func (s *WorldEngineService) GetAllRegions() []*types.Region {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*types.Region, 0, len(s.regions))
	for _, r := range s.regions {
		cp := *r
		result = append(result, &cp)
	}
	return result
}

// GetAllRegionsMap 返回所有区域的映射（用于持久化保存）
func (s *WorldEngineService) GetAllRegionsMap() map[string]*types.Region {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*types.Region, len(s.regions))
	for id, r := range s.regions {
		cp := *r
		result[id] = &cp
	}
	return result
}

// GetStateForPersistence 返回当前世界状态快照（用于持久化保存）
func (s *WorldEngineService) GetStateForPersistence() *WorldPersistenceState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &WorldPersistenceState{
		Epoch:   s.epoch,
		Metrics: s.worldState.BalanceMetrics,
	}
}

// ── Proto 转换 ──

func regionToProto(r *types.Region) *cultivation.Region {
	resources := make([]*cultivation.Resource, 0, len(r.Resources))
	for _, res := range r.Resources {
		resources = append(resources, resourceToProto(&res))
	}

	var parentID string
	if r.ParentRegionID != nil {
		parentID = string(*r.ParentRegionID)
	}

	return &cultivation.Region{
		Id:               string(r.ID),
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

func resourceToProto(r *types.Resource) *cultivation.Resource {
	var lastHarvested int64
	if r.LastHarvested != nil {
		lastHarvested = r.LastHarvested.Unix()
	}

	return &cultivation.Resource{
		Id:            r.ID,
		Name:          r.Name,
		Type:          r.Type,
		Rarity:        int32(r.Rarity),
		Quantity:      int32(r.Quantity),
		RespawnRate:   r.RespawnRate,
		LastHarvested: lastHarvested,
	}
}

func activeEventToProtoEvent(e *ActiveEvent) *cultivation.WorldEvent {
	participants := make([]string, 0)

	return &cultivation.WorldEvent{
		Id:             e.ID,
		Name:           e.Name,
		Type:           e.Type,
		Description:    e.Description,
		RegionId:       string(e.RegionID),
		StartTime:      e.StartTime.Unix(),
		EndTime:        e.EndTime.Unix(),
		ParticipantIds: participants,
		Status:         e.Status,
	}
}

func worldStateToProto(s *types.WorldState, activeEvents []*ActiveEvent) *cultivation.WorldState {
	regions := make(map[string]*cultivation.Region)
	for id, r := range s.Regions {
		region := r
		regions[string(id)] = regionToProto(&region)
	}

	events := make([]*cultivation.WorldEvent, 0, len(activeEvents))
	for _, e := range activeEvents {
		events = append(events, activeEventToProtoEvent(e))
	}

	return &cultivation.WorldState{
		Epoch:        s.Epoch,
		Regions:      regions,
		ActiveEvents: events,
		BalanceMetrics: &cultivation.BalanceMetrics{
			PowerDistribution:   s.BalanceMetrics.PowerDistribution,
			ResourceCirculation: s.BalanceMetrics.ResourceCirculation,
			SectDiversity:       s.BalanceMetrics.SectDiversity,
			KarmaDistribution:   s.BalanceMetrics.KarmaDistribution,
		},
		LastUpdated: s.LastUpdated.Unix(),
	}
}

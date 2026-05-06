package service

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/cultivation-world/shared/types"
)

// LeaderboardEntry 排行榜条目
type LeaderboardEntry struct {
	Rank     int     `json:"rank"`
	EntityID string  `json:"entity_id"`
	Name     string  `json:"name"`
	Value    float64 `json:"value"`
	Realm    string  `json:"realm"`
}

// LeaderboardData 排行榜缓存数据
type LeaderboardData struct {
	CultivationRank []*LeaderboardEntry `json:"cultivation_rank"`
	CombatRank      []*LeaderboardEntry `json:"combat_rank"`
	WealthRank      []*LeaderboardEntry `json:"wealth_rank"`
	MeritRank       []*LeaderboardEntry `json:"merit_rank"`
	UpdatedAt       time.Time           `json:"updated_at"`
	mu              sync.RWMutex
}

// LeaderboardCache 排行榜缓存（定时刷新）
var LeaderboardCache = &LeaderboardData{}

// LeaderboardStartRefresh 启动排行榜定时刷新（默认每5分钟）
func LeaderboardStartRefresh(svc *OperationService, interval time.Duration) {
	if interval <= 0 {
		interval = 5 * time.Minute
	}
	go func() {
		for {
			svc.refreshLeaderboard(context.Background())
			time.Sleep(interval)
		}
	}()
}

// executeLeaderboard 查看排行榜
func (s *OperationService) executeLeaderboard(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	boardType, _ := op.Params["board_type"].(string)

	LeaderboardCache.mu.RLock()
	defer LeaderboardCache.mu.RUnlock()

	var entries []*LeaderboardEntry
	var title string

	switch boardType {
	case "cultivation", "修为":
		entries = LeaderboardCache.CultivationRank
		title = "修为榜"
	case "combat", "战力":
		entries = LeaderboardCache.CombatRank
		title = "战力榜"
	case "wealth", "财富":
		entries = LeaderboardCache.WealthRank
		title = "财富榜"
	case "merit", "功德":
		entries = LeaderboardCache.MeritRank
		title = "功德榜"
	default:
		// 默认显示修为榜
		entries = LeaderboardCache.CultivationRank
		title = "修为榜"
	}

	if entries == nil {
		return &types.OperationResult{
			Success: true,
			Message: "排行榜数据尚未生成",
			Effects: map[string]interface{}{
				"board_type": boardType,
				"entries":    []interface{}{},
			},
		}, nil
	}

	rankList := make([]map[string]interface{}, 0, len(entries))
	for _, e := range entries {
		rankList = append(rankList, map[string]interface{}{
			"rank":      e.Rank,
			"entity_id": e.EntityID,
			"name":      e.Name,
			"value":     e.Value,
			"realm":     e.Realm,
		})
	}

	return &types.OperationResult{
		Success: true,
		Message: title,
		Effects: map[string]interface{}{
			"board_type": boardType,
			"title":      title,
			"entries":    rankList,
			"count":      len(rankList),
			"updated_at": LeaderboardCache.UpdatedAt.Unix(),
		},
	}, nil
}

// refreshLeaderboard 刷新排行榜数据
func (s *OperationService) refreshLeaderboard(ctx context.Context) {
	// 获取所有玩家实体
	players, err := s.getAllPlayers(ctx)
	if err != nil || len(players) == 0 {
		return
	}

	entries := make([]*LeaderboardEntry, 0, len(players))
	for _, p := range players {
		attr, _ := s.entityRepo.GetAttributes(ctx, p.ID)
		if attr == nil {
			continue
		}
		entries = append(entries, &LeaderboardEntry{
			EntityID: string(p.ID),
			Name:     p.Name,
			Realm:    string(p.Realm),
			Value:    attr.CultivationProgress + float64(types.CultivationRealmLevel(p.Realm))*100,
		})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Value > entries[j].Value })
	for i, e := range entries {
		e.Rank = i + 1
	}

	// 战力榜（深拷贝，避免共享指针）
	combatEntries := make([]*LeaderboardEntry, len(players))
	for i, p := range players {
		val := 0.0
		attr, _ := s.entityRepo.GetAttributes(ctx, p.ID)
		if attr != nil {
			bonuses := s.GetEquipmentBonuses(ctx, p.ID)
			val = float64(attr.AttackPower+attr.Defense+attr.Speed) +
				bonuses.AttackPower + bonuses.Defense + bonuses.Speed
		}
		combatEntries[i] = &LeaderboardEntry{
			EntityID: string(p.ID),
			Name:     p.Name,
			Realm:    string(p.Realm),
			Value:    val,
		}
	}
	sort.Slice(combatEntries, func(i, j int) bool { return combatEntries[i].Value > combatEntries[j].Value })
	for rank, e := range combatEntries {
		e.Rank = rank + 1
	}

	// 财富榜（深拷贝，避免共享指针）
	wealthEntries := make([]*LeaderboardEntry, len(players))
	for i, p := range players {
		val := 0.0
		attr, _ := s.entityRepo.GetAttributes(ctx, p.ID)
		if attr != nil {
			val = float64(attr.SpiritStones.LowGrade +
				attr.SpiritStones.MediumGrade*10 +
				attr.SpiritStones.HighGrade*100 +
				attr.SpiritStones.PremiumGrade*1000)
		}
		wealthEntries[i] = &LeaderboardEntry{
			EntityID: string(p.ID),
			Name:     p.Name,
			Realm:    string(p.Realm),
			Value:    val,
		}
	}
	sort.Slice(wealthEntries, func(i, j int) bool { return wealthEntries[i].Value > wealthEntries[j].Value })
	for i, e := range wealthEntries {
		e.Rank = i + 1
	}

	// 功德榜
	meritEntries := make([]*LeaderboardEntry, 0, len(players))
	for _, p := range players {
		meritEntries = append(meritEntries, &LeaderboardEntry{
			EntityID: string(p.ID),
			Name:     p.Name,
			Realm:    string(p.Realm),
			Value:    float64(p.Karma.Merit),
		})
	}
	sort.Slice(meritEntries, func(i, j int) bool { return meritEntries[i].Value > meritEntries[j].Value })
	for i, e := range meritEntries {
		e.Rank = i + 1
	}

	LeaderboardCache.mu.Lock()
	LeaderboardCache.CultivationRank = entries
	LeaderboardCache.CombatRank = combatEntries
	LeaderboardCache.WealthRank = wealthEntries
	LeaderboardCache.MeritRank = meritEntries
	LeaderboardCache.UpdatedAt = time.Now()
	LeaderboardCache.mu.Unlock()
}

// getAllPlayers 获取所有玩家实体（简化实现：通过数据库查询）
func (s *OperationService) getAllPlayers(ctx context.Context) ([]*types.Entity, error) {
	// 从游戏服务器注册表中获取所有玩家
	// 简化：实际应从数据库查询 entities WHERE entity_type = 'player'
	// 这里通过 entityRepo 扩展
	if entityLister, ok := s.entityRepo.(interface {
		GetAllPlayers(ctx context.Context) ([]*types.Entity, error)
	}); ok {
		return entityLister.GetAllPlayers(ctx)
	}
	return nil, nil
}

package service

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	cultivation "github.com/cultivation-world/shared/proto/pb"
	"github.com/cultivation-world/shared/config"
)

type AISchedulerService struct {
	cultivation.UnimplementedAISchedulerServiceServer
	cfg          *config.Config
	npcRegistry  map[string]*NPCProfile
	memoryStores map[string]*NPCMemoryStore
	templates    *BehaviorTemplateLibrary
	gameClient   GameServiceClient
	mu           sync.RWMutex
}

type NPCProfile struct {
	NPCID           string
	EntityID        string
	PersonalityType string
	MoralAlignment  string
	AmbitionLevel   int
	RiskTolerance   float64
	BackgroundStory string
	CurrentGoal     string
	CurrentRegion   string
	Realm           string
	Status          string
}

type BehaviorTemplateLibrary struct {
	templates map[string][]BehaviorTemplate
}

type BehaviorTemplate struct {
	Name      string
	Condition string
	Action    string
	Weight    float64
}

// 性格 → 动作基础权重（数值越高倾向越大）
var personalityWeights = map[string]map[string]float64{
	"aggressive": {
		"cultivate": 20, "meditate": 5, "explore": 20,
		"combat": 40, "gather": 10, "craft": 3, "trade": 2,
	},
	"cautious": {
		"cultivate": 35, "meditate": 30, "explore": 10,
		"combat": 5, "gather": 10, "craft": 5, "trade": 5,
	},
	"scholarly": {
		"cultivate": 35, "meditate": 20, "explore": 15,
		"combat": 5, "gather": 5, "craft": 15, "trade": 5,
	},
	"gregarious": {
		"cultivate": 20, "meditate": 10, "explore": 20,
		"combat": 15, "gather": 10, "craft": 5, "trade": 20,
	},
	"balanced": {
		"cultivate": 25, "meditate": 15, "explore": 20,
		"combat": 15, "gather": 10, "craft": 5, "trade": 10,
	},
}

// 目标 → 动作偏向修正
var goalAffinity = map[string]map[string]float64{
	"seek_power":      {"combat": 15, "explore": 10, "cultivate": 5},
	"collect_resources": {"gather": 20, "explore": 10, "trade": 5},
	"cultivate_realm":   {"cultivate": 25, "meditate": 10, "breakthrough": 15},
	"socialize":         {"trade": 15, "explore": 10, "combat": 5},
	"craft_mastery":     {"craft": 25, "gather": 10, "trade": 5},
}

func NewAISchedulerService(cfg *config.Config, gameClient GameServiceClient) *AISchedulerService {
	if cfg == nil {
		cfg = config.LoadConfigFromEnv()
	}
	return &AISchedulerService{
		cfg:          cfg,
		npcRegistry:  make(map[string]*NPCProfile),
		memoryStores: make(map[string]*NPCMemoryStore),
		templates:    NewBehaviorTemplateLibrary(),
		gameClient:   gameClient,
	}
}

func NewBehaviorTemplateLibrary() *BehaviorTemplateLibrary {
	return &BehaviorTemplateLibrary{
		templates: map[string][]BehaviorTemplate{
			"cultivate": {
				{Name: "daily_cultivation", Condition: "qi<50%", Action: "cultivate", Weight: 0.8},
				{Name: "breakthrough_attempt", Condition: "progress>=100%", Action: "breakthrough", Weight: 0.6},
			},
			"explore": {
				{Name: "resource_gathering", Condition: "low_resources", Action: "gather", Weight: 0.7},
				{Name: "region_exploration", Condition: "curious", Action: "explore", Weight: 0.5},
			},
			"social": {
				{Name: "seek_alliance", Condition: "weak", Action: "form_alliance", Weight: 0.4},
				{Name: "trade", Condition: "surplus_resources", Action: "trade", Weight: 0.6},
			},
		},
	}
}

// GetNPCIDs returns all registered NPC IDs (for the background loop).
func (s *AISchedulerService) GetNPCIDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := make([]string, 0, len(s.npcRegistry))
	for id := range s.npcRegistry {
		ids = append(ids, id)
	}
	return ids
}

// GetNPC returns a copy of an NPC profile.
func (s *AISchedulerService) GetNPC(npcID string) *NPCProfile {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, exists := s.npcRegistry[npcID]
	if !exists {
		return nil
	}
	cp := *p
	return &cp
}

// GetMemoryStore returns the memory store for an NPC.
func (s *AISchedulerService) GetMemoryStore(npcID string) *NPCMemoryStore {
	s.mu.RLock()
	store, exists := s.memoryStores[npcID]
	s.mu.RUnlock()
	if !exists {
		s.mu.Lock()
		// Double-check: another goroutine may have created it while we unlocked
		if store, exists = s.memoryStores[npcID]; !exists {
			store = NewNPCMemoryStore(npcID)
			s.memoryStores[npcID] = store
		}
		s.mu.Unlock()
	}
	return store
}

// UpdateNPCState directly updates runtime fields on a registered NPC (no copy).
func (s *AISchedulerService) UpdateNPCState(npcID, realm, region, status string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p, ok := s.npcRegistry[npcID]; ok {
		if realm != "" {
			p.Realm = realm
		}
		if region != "" {
			p.CurrentRegion = region
		}
		if status != "" {
			p.Status = status
		}
	}
}

// UpdateNPCGoal directly updates the current goal on a registered NPC (no copy).
func (s *AISchedulerService) UpdateNPCGoal(npcID, goal string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p, ok := s.npcRegistry[npcID]; ok {
		p.CurrentGoal = goal
	}
}

// RecordInteraction records a player-NPC interaction in memory.
func (s *AISchedulerService) RecordInteraction(npcID string, playerID string, playerName string, content string, affinityDelta int) {
	store := s.GetMemoryStore(npcID)
	store.RememberInteraction(playerID, playerName, content, affinityDelta)
}

// ScheduleDecision 两步决策：
//
//	① 行为树模板匹配（~70% 基础行为）
//	② 加权算法引擎（综合性格/目标/记忆/关系/背景/境界/风险/野心）
func (s *AISchedulerService) ScheduleDecision(ctx context.Context, req *cultivation.DecisionRequest) (*cultivation.DecisionResponse, error) {
	s.mu.RLock()
	profile, exists := s.npcRegistry[req.NpcId]
	s.mu.RUnlock()

	if !exists {
		profile = &NPCProfile{
			NPCID:           req.NpcId,
			PersonalityType: "balanced",
			MoralAlignment:  "neutral",
			AmbitionLevel:   50,
			RiskTolerance:   0.5,
		}
	}

	// 第一步：行为树模板优先（处理日常/习惯性行为）
	template := s.matchTemplate(req.Context, req.AvailableActions)
	if template != nil && rand.Float64() < 0.7 {
		return &cultivation.DecisionResponse{
			Action:    template.Action,
			Params:    make(map[string]string),
			Reasoning: fmt.Sprintf("Template matched: %s", template.Name),
			Source:    "behavior_tree",
			TokenCost: 0,
		}, nil
	}

	// 第二步：加权算法决策
	decision := s.algorithmicDecision(profile, req.Context, req.AvailableActions)
	return decision, nil
}

// algorithmicDecision 核心算法：加权随机选择
//
//	权重 = 性格基础 + 背景修正 + 目标匹配 + 境界影响 + 记忆影响 + 关系修正 + 风险修正 + 野心修正 + 随机波动
func (s *AISchedulerService) algorithmicDecision(profile *NPCProfile, _ string, available []string) *cultivation.DecisionResponse {
	weights, factors := s.calculateWeights(profile, "", available)
	action := s.weightedSelect(available, weights)

	params := s.buildActionParams(action, profile)
	reasoning := s.buildReasoning(action, profile, factors)

	return &cultivation.DecisionResponse{
		Action:    action,
		Params:    params,
		Reasoning: reasoning,
		Source:    "algorithm",
		TokenCost: 0,
	}
}

// buildActionParams 为动作补充上下文参数
func (s *AISchedulerService) buildActionParams(action string, _ *NPCProfile) map[string]string {
	params := make(map[string]string)
	switch action {
	case "cultivate":
		params["duration"] = "30m"
	case "meditate":
		params["duration"] = "30m"
	case "explore":
		params["direction"] = "random"
	case "move":
		params["target_region"] = "random"
	case "combat":
		params["target_type"] = "wandering_monster"
	case "gather":
		params["resource_type"] = "any"
	}
	return params
}

// buildReasoning 生成人类可读的中文决策理由
func (s *AISchedulerService) buildReasoning(action string, profile *NPCProfile, factors map[string]float64) string {
	actionNames := map[string]string{
		"cultivate": "修炼", "meditate": "冥想", "explore": "探索",
		"combat": "战斗", "gather": "采集", "craft": "炼制",
		"trade": "交易", "breakthrough": "突破", "move": "移动",
	}

	actionName := actionNames[action]
	if actionName == "" {
		actionName = action
	}

	personalityDesc := map[string]string{
		"aggressive": "激进好战", "cautious": "谨慎小心",
		"scholarly": "学者风范", "gregarious": "喜好交际",
		"balanced": "性格平和",
	}

	parts := []string{
		fmt.Sprintf("[%s]", personalityDesc[profile.PersonalityType]),
		fmt.Sprintf("决定%s", actionName),
	}

	if profile.CurrentGoal != "" {
		parts = append(parts, fmt.Sprintf("目标[%s]", profile.CurrentGoal))
	}

	if f, ok := factors["memory"]; ok && f > 10 {
		parts = append(parts, "受近期记忆影响")
	}
	if f, ok := factors["enemy"]; ok && f > 0 {
		parts = append(parts, "感应到敌对目标")
	}
	if f, ok := factors["ambition"]; ok && f > 5 {
		parts = append(parts, "野心驱使")
	}
	if f, ok := factors["risk"]; ok && f > 5 {
		parts = append(parts, "偏好冒险")
	}
	if f, ok := factors["background"]; ok && f > 0 {
		parts = append(parts, "背景经历影响")
	}
	if f, ok := factors["realm"]; ok && f > 0 {
		parts = append(parts, "境界自信")
	} else if f, ok := factors["realm"]; ok && f < 0 {
		parts = append(parts, "境界不足需谨慎")
	}
	if f, ok := factors["goal"]; ok && f > 10 {
		parts = append(parts, "目标驱动")
	}

	return strings.Join(parts, "，")
}

// calculateWeights 计算所有可用动作的最终权重，同时返回各因素的影响值
func (s *AISchedulerService) calculateWeights(profile *NPCProfile, _ string, available []string) (map[string]float64, map[string]float64) {
	weights := make(map[string]float64, len(available))
	factors := make(map[string]float64) // 记录各因素的影响，用于生成推理

	// 获取性格基础权重表
	baseWeights, ok := personalityWeights[profile.PersonalityType]
	if !ok {
		baseWeights = personalityWeights["balanced"]
	}

	// 获取目标匹配表
	goalMap := s.matchGoal(profile.CurrentGoal)

	// 获取记忆影响
	memoryStore := s.GetMemoryStore(profile.NPCID)
	memories := memoryStore.GetRecentMemories(5, "")
	hasNegativeMemory := false
	for _, m := range memories {
		if m.Importance >= 0.7 {
			hasNegativeMemory = true
			break
		}
	}
	if hasNegativeMemory {
		factors["memory"] = 20
	}

	// 附近有敌对关系？
	rels := memoryStore.GetAllRelationships()
	hasEnemyNearby := false
	for _, rel := range rels {
		if rel.Affinity < -30 {
			hasEnemyNearby = true
			break
		}
	}
	if hasEnemyNearby {
		factors["enemy"] = 15
	}

	// 背景故事影响（按动作区分）
	bgModifiers := s.calcBackgroundModifiers(profile.BackgroundStory)
	hasBgMod := false
	for _, v := range bgModifiers {
		if v != 0 {
			hasBgMod = true
			break
		}
	}
	if hasBgMod {
		factors["background"] = 5 // 只标记有影响，实际值按动作
	}

	// 境界影响
	realmFactor := s.calcRealmFactor(profile.Realm)
	if realmFactor != 0 {
		factors["realm"] = realmFactor
	}

	// 目标影响
	if goalMap != nil {
		totalGoal := 0.0
		for _, v := range goalMap {
			totalGoal += v
		}
		if totalGoal > 0 {
			factors["goal"] = totalGoal
		}
	}

	// 风险偏好
	riskBonus := (profile.RiskTolerance - 0.5) * 30
	factors["risk"] = riskBonus

	// 野心修正
	ambitionFactor := float64(profile.AmbitionLevel-50) / 5
	factors["ambition"] = ambitionFactor

	for _, action := range available {
		w := 0.0

		// ① 性格基础权重（0-40）
		w += baseWeights[action]

		// ② 目标匹配度（0-25）
		if goalMap != nil {
			w += goalMap[action]
		}

		// ③ 记忆影响（-30 ～ +30）
		if hasNegativeMemory {
			switch action {
			case "combat":
				w += 20
			case "meditate":
				w += 5
			}
		}

		// ④ 关系修正（±20）
		if hasEnemyNearby {
			switch action {
			case "combat":
				w += 15
			case "flee", "meditate":
				w += 10
			}
		}

		// ⑤ 背景故事修正（±15 按动作区分）
		w += bgModifiers[action]

		// ⑥ 境界影响（±10）
		w += realmFactor

		// ⑦ 风险偏好（±15）
		switch action {
		case "combat", "explore":
			w += riskBonus
		case "meditate", "cultivate":
			w -= riskBonus
		}

		// ⑧ 野心修正（±10）
		switch action {
		case "breakthrough", "combat", "explore":
			w += ambitionFactor
		case "meditate":
			w -= ambitionFactor
		}

		// ⑨ 随机波动（0-20），防止行为固化
		w += rand.Float64() * 20

		weights[action] = w
	}

	return weights, factors
}

// calcBackgroundModifiers 从背景故事中提取按动作区分的偏向修正值
func (s *AISchedulerService) calcBackgroundModifiers(background string) map[string]float64 {
	mods := map[string]float64{
		"cultivate": 0, "meditate": 0, "explore": 0,
		"combat": 0, "gather": 0, "craft": 0, "trade": 0,
	}
	if background == "" {
		return mods
	}

	type bgRule struct {
		words  []string
		mods   map[string]float64 // 各动作的修正值
	}
	rules := []bgRule{
		{[]string{"散修", "独行", "孤", "流浪"}, map[string]float64{"cultivate": 5, "explore": 5, "trade": -8, "craft": -5}},
		{[]string{"宗门", "大派", "名门", "正道"}, map[string]float64{"trade": 5, "craft": 3, "combat": 3}},
		{[]string{"魔道", "邪修", "嗜血", "杀戮"}, map[string]float64{"combat": 12, "explore": 5, "meditate": -5, "cultivate": -3}},
		{[]string{"丹", "药", "医", "草"}, map[string]float64{"gather": 8, "craft": 6, "trade": 3}},
		{[]string{"战", "斗", "武", "剑", "杀"}, map[string]float64{"combat": 10, "explore": 4, "meditate": -3}},
		{[]string{"商", "交易", "市", "贾"}, map[string]float64{"trade": 10, "gather": 3, "combat": -5}},
		{[]string{"古", "遗迹", "上古", "秘境"}, map[string]float64{"explore": 10, "gather": 5, "meditate": -3}},
	}

	for _, rule := range rules {
		for _, word := range rule.words {
			if strings.Contains(background, word) {
				for action, delta := range rule.mods {
					mods[action] += delta
				}
				break
			}
		}
	}

	return mods
}

// calcRealmFactor 境界影响：高阶更自信，低阶更谨慎
func (s *AISchedulerService) calcRealmFactor(realm string) float64 {
	realmLevels := map[string]int{
		"mortal":           0,
		"qi_condensation":  1,
		"foundation":       2,
		"golden_core":      3,
		"nascent_soul":     4,
		"soul_transformation": 5,
		"void_refinement":  6,
		"integration":      7,
		"mahayana":         8,
		"tribulation":      9,
	}
	level, exists := realmLevels[realm]
	if !exists {
		return 0
	}
	// 炼气以下 → 负修正（谨慎），金丹以上 → 正修正（自信）
	return float64(level-2) * 2 // -4 ~ +14, 基数低阶负、高阶正
}

// weightedSelect 按权重随机选择一个动作
func (s *AISchedulerService) weightedSelect(available []string, weights map[string]float64) string {
	if len(available) == 0 {
		return "meditate"
	}
	if len(available) == 1 {
		return available[0]
	}

	total := 0.0
	for _, a := range available {
		w := weights[a]
		if w < 0 {
			w = 0
		}
		total += w
	}

	if total <= 0 {
		return available[rand.Intn(len(available))]
	}

	r := rand.Float64() * total
	cumulative := 0.0
	for _, a := range available {
		w := weights[a]
		if w < 0 {
			w = 0
		}
		cumulative += w
		if r <= cumulative {
			return a
		}
	}

	return available[len(available)-1]
}

// matchGoal 将文本目标映射到权重修正表（按优先级顺序匹配，首个命中返回）
func (s *AISchedulerService) matchGoal(goal string) map[string]float64 {
	if goal == "" {
		return nil
	}

	// 优先级从高到低排序，避免 map 随机遍历导致不同次返回不同结果
	type kwMapping struct {
		keyword  string
		goalType string
	}
	keywords := []kwMapping{
		{"炼丹", "craft_mastery"},
		{"炼器", "craft_mastery"},
		{"收集", "collect_resources"},
		{"结交", "socialize"},
		{"突破", "seek_power"},
		{"对手", "seek_power"},
		{"修炼", "cultivate_realm"},
		{"境界", "cultivate_realm"},
		{"修为", "cultivate_realm"},
		{"资源", "collect_resources"},
		{"材料", "collect_resources"},
		{"道友", "socialize"},
	}

	for _, k := range keywords {
		if strings.Contains(goal, k.keyword) {
			if m, ok := goalAffinity[k.goalType]; ok {
				return m
			}
		}
	}
	return nil
}

func (s *AISchedulerService) ExecuteBehaviorTree(ctx context.Context, req *cultivation.BehaviorTreeRequest) (*cultivation.BehaviorTreeResponse, error) {
	action := s.executeBehaviorTreeLogic(req.TreeName, req.Context)

	return &cultivation.BehaviorTreeResponse{
		Action:  action.Action,
		Params:  action.Params,
		Success: true,
	}, nil
}

func (s *AISchedulerService) RegisterNPC(ctx context.Context, req *cultivation.NPCRegisterRequest) (*cultivation.NPCRegisterResponse, error) {
	profile := &NPCProfile{
		NPCID:           req.NpcId,
		PersonalityType: req.PersonalityType,
		MoralAlignment:  req.MoralAlignment,
		AmbitionLevel:   int(req.AmbitionLevel),
		RiskTolerance:   req.RiskTolerance,
		BackgroundStory: req.BackgroundStory,
		CurrentGoal:     req.CurrentGoal,
	}

	s.mu.Lock()
	s.npcRegistry[req.NpcId] = profile
	if _, exists := s.memoryStores[req.NpcId]; !exists {
		s.memoryStores[req.NpcId] = NewNPCMemoryStore(req.NpcId)
	}
	s.mu.Unlock()

	return &cultivation.NPCRegisterResponse{
		Success: true,
		Message: "NPC registered successfully",
	}, nil
}

func (s *AISchedulerService) UnregisterNPC(ctx context.Context, req *cultivation.NPCUnregisterRequest) (*cultivation.NPCUnregisterResponse, error) {
	s.mu.Lock()
	delete(s.npcRegistry, req.NpcId)
	delete(s.memoryStores, req.NpcId)
	s.mu.Unlock()

	return &cultivation.NPCUnregisterResponse{
		Success: true,
	}, nil
}

func (s *AISchedulerService) matchTemplate(_ string, availableActions []string) *BehaviorTemplate {
	for _, actionType := range []string{"cultivate", "explore", "social"} {
		for _, template := range s.templates.templates[actionType] {
			for _, action := range availableActions {
				if template.Action == action {
					return &template
				}
			}
		}
	}
	return nil
}

func (s *AISchedulerService) executeBehaviorTreeLogic(treeName string, context map[string]string) *BehaviorTreeResult {
	switch treeName {
	case "daily_routine":
		return &BehaviorTreeResult{
			Action: "cultivate",
			Params: map[string]string{"duration": "1h"},
		}
	case "combat":
		return &BehaviorTreeResult{
			Action: "attack",
			Params: map[string]string{"target": context["enemy_id"]},
		}
	case "exploration":
		return &BehaviorTreeResult{
			Action: "explore",
			Params: map[string]string{"direction": "random"},
		}
	default:
		return &BehaviorTreeResult{
			Action: "meditate",
			Params: make(map[string]string),
		}
	}
}

// ExecuteNPCAction executes an NPC's chosen action through the game-server.
func (s *AISchedulerService) ExecuteNPCAction(npcID string, action string, params map[string]string) *OperationResult {
	if s.gameClient == nil {
		return &OperationResult{Success: false, Message: "game client not connected"}
	}

	s.mu.RLock()
	profile, exists := s.npcRegistry[npcID]
	s.mu.RUnlock()

	if !exists {
		return &OperationResult{Success: false, Message: "NPC not found"}
	}

	entityID := profile.EntityID
	if entityID == "" {
		entityID = npcID
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.gameClient.ExecuteOperation(ctx, entityID, action, params)
	if err != nil {
		return &OperationResult{Success: false, Message: err.Error()}
	}

	s.mu.Lock()
	if p, ok := s.npcRegistry[npcID]; ok {
		p.Status = action
	}
	s.mu.Unlock()

	return result
}

type BehaviorTreeResult struct {
	Action string
	Params map[string]string
}


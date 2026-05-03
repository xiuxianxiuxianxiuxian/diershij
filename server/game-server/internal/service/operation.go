package service

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/cultivation-world/shared/errors"
	"github.com/cultivation-world/shared/types"
)

type OperationService struct {
	entityRepo    EntityRepository
	itemRepo      ItemRepository
	inventoryRepo InventoryRepository
	spellRepo     SpellRepository
	messageRepo   MessageRepository
	worldService  WorldServiceClient
	daoService    HeavenlyDaoClient
	sectRepo      SectRepository
	recipeRepo    RecipeRepository
}

type EntityRepository interface {
	GetByID(ctx context.Context, id types.EntityID) (*types.Entity, error)
	GetByName(ctx context.Context, name string) (*types.Entity, error)
	Create(ctx context.Context, entity *types.Entity) error
	Update(ctx context.Context, entity *types.Entity) error
	GetAttributes(ctx context.Context, entityID types.EntityID) (*types.Attributes, error)
	UpdateAttributes(ctx context.Context, entityID types.EntityID, attr *types.Attributes) error
}

type ItemRepository interface {
	GetByID(ctx context.Context, id types.ItemID) (*types.Item, error)
	GetByName(ctx context.Context, name string) (*types.Item, error)
}

type InventoryRepository interface {
	GetByEntityID(ctx context.Context, entityID types.EntityID) ([]*types.InventoryItem, error)
	GetItem(ctx context.Context, entityID types.EntityID, itemID types.ItemID) (*types.InventoryItem, error)
	AddItem(ctx context.Context, entityID types.EntityID, itemID types.ItemID, quantity int) error
	RemoveItem(ctx context.Context, entityID types.EntityID, itemID types.ItemID, quantity int) error
	EquipItem(ctx context.Context, entityID types.EntityID, itemID types.ItemID, slot string) error
}

type SpellRepository interface {
	GetByID(ctx context.Context, id types.SpellID) (*types.Spell, error)
	GetEntitySpell(ctx context.Context, entityID types.EntityID, spellID types.SpellID) (*types.EntitySpell, error)
	UpdateSpellCastTime(ctx context.Context, entityID types.EntityID, spellID types.SpellID) error
}

type MessageRepository interface {
	Create(ctx context.Context, message *types.DBMessage) error
}

type WorldServiceClient interface {
	GetRegion(ctx context.Context, regionID string) (*types.Region, error)
	SpawnResources(ctx context.Context, regionID string) error
}

type SectInfo struct {
	ID        string
	Name      string
	FounderID string
	Alignment string
}

type SectRepository interface {
	Create(ctx context.Context, sectID string, name string, founderID string) error
	GetByID(ctx context.Context, id string) (*SectInfo, error)
	GetByName(ctx context.Context, name string) (*SectInfo, error)
	AddMember(ctx context.Context, sectID string, entityID string, rank string) error
	GetMember(ctx context.Context, sectID string, entityID string) (bool, error)
	RemoveMember(ctx context.Context, sectID string, entityID string) error
}

type RecipeInfo struct {
	ID         string
	Type       string
	Difficulty int
	Name       string
}

type RecipeRepository interface {
	GetByID(ctx context.Context, id string) (*RecipeInfo, error)
}

func NewOperationService(
	entityRepo EntityRepository,
	itemRepo ItemRepository,
	inventoryRepo InventoryRepository,
	spellRepo SpellRepository,
	messageRepo MessageRepository,
	worldService WorldServiceClient,
	daoService HeavenlyDaoClient,
	sectRepo SectRepository,
	recipeRepo RecipeRepository,
) *OperationService {
	return &OperationService{
		entityRepo:    entityRepo,
		itemRepo:      itemRepo,
		inventoryRepo: inventoryRepo,
		spellRepo:     spellRepo,
		messageRepo:   messageRepo,
		worldService:  worldService,
		daoService:    daoService,
		sectRepo:      sectRepo,
		recipeRepo:    recipeRepo,
	}
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
	case types.ActionCombat:
		return s.executeCombat(ctx, entity, op)
	case types.ActionExplore:
		return s.executeExplore(ctx, entity, op)
	case types.ActionGather:
		return s.executeGather(ctx, entity, op)
	case types.ActionCraft:
		return s.executeCraft(ctx, entity, op)
	case types.ActionCreateMethod:
		return s.executeCreateMethod(ctx, entity, op)
	case types.ActionTrade:
		return s.executeTrade(ctx, entity, op)
	case types.ActionFormSect:
		return s.executeFormSect(ctx, entity, op)
	case types.ActionJoinSect:
		return s.executeJoinSect(ctx, entity, op)
	case types.ActionSendMessage:
		return s.executeSendMessage(ctx, entity, op)
	case types.ActionCastSpell:
		return s.executeCastSpell(ctx, entity, op)
	case types.ActionLeaveSect:
		return s.executeLeaveSect(ctx, entity, op)
	case types.ActionAddFriend:
		return s.executeAddFriend(ctx, entity, op)
	case types.ActionRemoveFriend:
		return s.executeRemoveFriend(ctx, entity, op)
	case types.ActionAcceptFriend:
		return s.executeAcceptFriend(ctx, entity, op)
	case types.ActionFlee:
		return s.executeFlee(ctx, entity, op)
	case types.ActionUseSkill:
		return s.executeUseSkill(ctx, entity, op)
	default:
		return nil, errors.ErrInvalidOperationType
	}
}

// 修炼
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

// 移动
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

// 打坐
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
			"qi_recovery":        recovery,
			"spiritual_recovery": spRecovery,
		},
	}, nil
}

// 休息
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
			"qi":              attr.Qi,
			"spiritual_power": attr.SpiritualPower,
		},
	}, nil
}

// 突破
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

	newRealm := getNextRealm(entity.Realm)
	if newRealm == "" {
		return &types.OperationResult{
			Success: false,
			Message: "已达最高境界",
		}, nil
	}

	// 天道劫数检查
	if s.daoService != nil {
		tribResult, err := s.daoService.CheckTribulation(ctx, string(entity.ID), string(newRealm))
		if err == nil && !tribResult.Success {
			return &types.OperationResult{
				Success: false,
				Message: "突破失败：" + tribResult.Reason,
				Effects: map[string]interface{}{
					"tribulation": true,
					"severity":    tribResult.Severity,
				},
			}, nil
		}
	}

	successRate := 0.5 + float64(attr.Luck)/200.0
	if successRate > 0.8 {
		successRate = 0.8
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
			"new_realm":     newRealm,
			"success_rate":  successRate,
			"max_qi":        attr.MaxQi,
			"max_spiritual": attr.MaxSpiritualPower,
		},
	}, nil
}

// 战斗
func (s *OperationService) executeCombat(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	targetID, ok := op.Params["target_id"].(string)
	if !ok {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "target_id is required")
	}

	// 获取目标
	target, err := s.entityRepo.GetByID(ctx, types.EntityID(targetID))
	if err != nil {
		return nil, errors.NewGameError(errors.ErrEntityNotFound, "target not found")
	}

	// 检查目标状态
	if target.Status == types.StatusDead {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "target is already dead")
	}

	// 检查是否在同一区域
	if entity.Position.RegionID != target.Position.RegionID {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "target is not in the same region")
	}

	// 获取属性
	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)
	targetAttr, _ := s.entityRepo.GetAttributes(ctx, target.ID)

	// 计算距离
	distance := math.Sqrt(math.Pow(entity.Position.X-target.Position.X, 2) + math.Pow(entity.Position.Y-target.Position.Y, 2))
	if distance > 10 {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "target is too far away")
	}

	// 战斗计算
	result := s.calculateCombat(attr, targetAttr)

	// 更新状态
	entity.Status = types.StatusCombat
	target.Status = types.StatusCombat

	// 应用伤害
	if result.DamageDealt > 0 {
		targetAttr.Qi -= float64(result.DamageDealt)
		if targetAttr.Qi <= 0 {
			targetAttr.Qi = 0
			target.Status = types.StatusDead
			result.Killed = true
		}
	}

	// 保存更新
	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)
	s.entityRepo.UpdateAttributes(ctx, target.ID, targetAttr)
	s.entityRepo.Update(ctx, entity)
	s.entityRepo.Update(ctx, target)

	return &types.OperationResult{
		Success: true,
		Message: result.Message,
		Effects: map[string]interface{}{
			"damage_dealt":  result.DamageDealt,
			"is_crit":       result.IsCrit,
			"is_dodge":      result.IsDodge,
			"target_status": target.Status,
		},
	}, nil
}

type CombatResult struct {
	DamageDealt int
	IsCrit      bool
	IsDodge     bool
	Killed      bool
	Message     string
}

func (s *OperationService) calculateCombat(attacker, defender *types.Attributes) CombatResult {
	result := CombatResult{}

	// 闪避判定
	dodgeRoll := rand.Float64() * 100
	if dodgeRoll < float64(defender.DodgeRate) {
		result.IsDodge = true
		result.Message = "攻击被闪避！"
		return result
	}

	// 基础伤害计算
	baseDamage := float64(attacker.AttackPower) - float64(defender.Defense)*0.5
	if baseDamage < 1 {
		baseDamage = 1
	}

	// 暴击判定
	critRoll := rand.Float64() * 100
	if critRoll < float64(attacker.CritRate) {
		result.IsCrit = true
		baseDamage *= 1.5
		result.Message = "暴击！"
	} else {
		result.Message = "攻击命中！"
	}

	result.DamageDealt = int(baseDamage)
	return result
}

// 探索
func (s *OperationService) executeExplore(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)

	// 获取区域信息
	region, err := s.worldService.GetRegion(ctx, entity.Position.RegionID)
	if err != nil {
		return nil, errors.NewGameError(errors.ErrRegionNotFound, "region not found")
	}

	// 消耗灵气
	qiCost := float64(region.SpiritualTier) * 5
	if attr.Qi < qiCost {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "not enough qi")
	}
	attr.Qi -= qiCost

	// 探索成功率
	successRate := 0.3 + float64(region.SpiritualDensity)/200.0
	if rand.Float64() > successRate {
		s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)
		return &types.OperationResult{
			Success: true,
			Message: "探索了一番，但一无所获",
			Effects: map[string]interface{}{
				"qi_cost": qiCost,
			},
		}, nil
	}

	// 探索结果
	discoveries := []string{}

	// 发现资源
	if len(region.Resources) > 0 && rand.Float64() < 0.4 {
		resource := region.Resources[rand.Intn(len(region.Resources))]
		discoveries = append(discoveries, fmt.Sprintf("发现了%s", resource.Name))
	}

	// 奇遇事件
	if rand.Float64() < 0.1 {
		discoveries = append(discoveries, "触发奇遇！")
	}

	entity.Status = types.StatusExploring
	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)
	s.entityRepo.Update(ctx, entity)

	return &types.OperationResult{
		Success: true,
		Message: "探索成功！",
		Effects: map[string]interface{}{
			"qi_cost":     qiCost,
			"discoveries": discoveries,
		},
	}, nil
}

// 采集
func (s *OperationService) executeGather(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	resourceType, ok := op.Params["resource_type"].(string)
	if !ok {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "resource_type is required")
	}

	quantity, _ := op.Params["quantity"].(float64)
	if quantity <= 0 {
		quantity = 1
	}

	// 获取区域信息
	region, err := s.worldService.GetRegion(ctx, entity.Position.RegionID)
	if err != nil {
		return nil, errors.NewGameError(errors.ErrRegionNotFound, "region not found")
	}

	// 查找资源
	var targetResource *types.Resource
	for i := range region.Resources {
		if region.Resources[i].Type == resourceType {
			targetResource = &region.Resources[i]
			break
		}
	}

	if targetResource == nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "resource not found in this region")
	}

	if targetResource.Quantity < int(quantity) {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "not enough resources")
	}

	// 获取属性
	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)

	// 消耗灵气
	qiCost := float64(targetResource.Rarity) * 5 * quantity
	if attr.Qi < qiCost {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "not enough qi")
	}
	attr.Qi -= qiCost

	// 增加采集技能经验
	attr.HerbKnowledge += int(quantity)

	// 创建物品（简化处理，实际应该根据资源类型创建对应物品）
	// 这里假设资源直接作为物品添加到背包

	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)

	// 刷新资源
	s.worldService.SpawnResources(ctx, string(region.ID))

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("采集成功！获得 %s x%.0f", targetResource.Name, quantity),
		Effects: map[string]interface{}{
			"resource":  targetResource.Name,
			"quantity":  quantity,
			"qi_cost":   qiCost,
			"skill_exp": quantity,
		},
	}, nil
}

// 炼器/炼丹
func (s *OperationService) executeCraft(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	recipeID, ok := op.Params["recipe_id"].(string)
	if !ok {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "recipe_id is required")
	}

	// 获取属性
	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)

	// 查询配方
	recipeDifficulty := 3
	if s.recipeRepo != nil {
		recipe, err := s.recipeRepo.GetByID(ctx, recipeID)
		if err == nil && recipe != nil {
			recipeDifficulty = recipe.Difficulty
		}
	}

	// 计算成功率
	successRate := 0.5 + float64(attr.AlchemyLevel-recipeDifficulty)*0.1
	successRate += float64(attr.Comprehension) / 200.0
	if successRate > 0.95 {
		successRate = 0.95
	}
	if successRate < 0.1 {
		successRate = 0.1
	}

	// 消耗灵气
	qiCost := float64(recipeDifficulty) * 20
	if attr.Qi < qiCost {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "not enough qi")
	}
	attr.Qi -= qiCost

	// 判定成功
	if rand.Float64() > successRate {
		// 失败
		s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)
		return &types.OperationResult{
			Success: false,
			Message: "制作失败，材料损失",
			Effects: map[string]interface{}{
				"qi_cost":      qiCost,
				"success_rate": successRate,
			},
		}, nil
	}

	// 成功
	attr.AlchemyLevel += 1
	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)

	return &types.OperationResult{
		Success: true,
		Message: "制作成功！",
		Effects: map[string]interface{}{
			"qi_cost":      qiCost,
			"success_rate": successRate,
			"skill_exp":    1,
		},
	}, nil
}

// 自创功法
func (s *OperationService) executeCreateMethod(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	// 验证境界
	if types.CultivationRealmLevel(entity.Realm) < types.CultivationRealmLevel(types.RealmNascentSoul) {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "requires Nascent Soul realm or higher")
	}

	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)

	// 验证悟性
	if attr.Comprehension < 80 {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "comprehension must be at least 80")
	}

	// 消耗资源
	qiCost := attr.MaxQi * 0.5
	spCost := float64(attr.DivineSense) * 0.5

	if attr.Qi < qiCost {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "not enough qi")
	}
	if attr.SpiritualPower < spCost {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "not enough spiritual power")
	}

	attr.Qi -= qiCost
	attr.SpiritualPower -= spCost

	// 计算功法品质
	quality := s.calculateMethodQuality(attr)

	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("成功创造%s功法！", quality),
		Effects: map[string]interface{}{
			"qi_cost":        qiCost,
			"sp_cost":        spCost,
			"method_quality": quality,
		},
	}, nil
}

func (s *OperationService) calculateMethodQuality(attr *types.Attributes) string {
	score := float64(attr.Comprehension)*0.4 + float64(attr.DivineSense)*0.3 + float64(attr.Luck)*0.3

	if score >= 150 {
		return "天级"
	} else if score >= 120 {
		return "地级"
	} else if score >= 90 {
		return "玄级"
	}
	return "黄级"
}

// 交易
func (s *OperationService) executeTrade(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	targetID, ok := op.Params["target_id"].(string)
	if !ok {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "target_id is required")
	}

	itemID, ok := op.Params["item_id"].(string)
	if !ok {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "item_id is required")
	}

	price, ok := op.Params["price"].(float64)
	if !ok || price <= 0 {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "invalid price")
	}

	// 获取目标
	target, err := s.entityRepo.GetByID(ctx, types.EntityID(targetID))
	if err != nil {
		return nil, errors.NewGameError(errors.ErrEntityNotFound, "target not found")
	}

	// 检查是否在同一区域
	if entity.Position.RegionID != target.Position.RegionID {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "target is not in the same region")
	}

	// 获取属性
	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)
	targetAttr, _ := s.entityRepo.GetAttributes(ctx, target.ID)

	// 检查灵石
	if targetAttr.SpiritStones.LowGrade < int64(price) {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "target does not have enough spirit stones")
	}

	// 执行交易（简化处理）
	targetAttr.SpiritStones.LowGrade -= int64(price)
	attr.SpiritStones.LowGrade += int64(price)

	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)
	s.entityRepo.UpdateAttributes(ctx, target.ID, targetAttr)

	return &types.OperationResult{
		Success: true,
		Message: "交易成功！",
		Effects: map[string]interface{}{
			"price":   price,
			"item_id": itemID,
		},
	}, nil
}

// 创建宗门
func (s *OperationService) executeFormSect(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	sectName, ok := op.Params["sect_name"].(string)
	if !ok || sectName == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "sect_name is required")
	}

	// 验证境界
	if types.CultivationRealmLevel(entity.Realm) < types.CultivationRealmLevel(types.RealmGoldenCore) {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "requires Golden Core realm or higher")
	}

	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)

	// 检查费用
	cost := int64(10000)
	if attr.SpiritStones.LowGrade < cost {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "not enough spirit stones (need 10000)")
	}

	attr.SpiritStones.LowGrade -= cost
	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)

	// 持久化宗门到数据库
	sectID := string(types.GenerateEntityID())
	if err := s.sectRepo.Create(ctx, sectID, sectName, string(entity.ID)); err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "failed to create sect")
	}
	s.sectRepo.AddMember(ctx, sectID, string(entity.ID), "sect_leader")

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("成功创建宗门：%s！", sectName),
		Effects: map[string]interface{}{
			"sect_name": sectName,
			"sect_id":   sectID,
			"cost":      cost,
		},
	}, nil
}

// 加入宗门
func (s *OperationService) executeJoinSect(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	sectID, ok := op.Params["sect_id"].(string)
	if !ok {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "sect_id is required")
	}

	// 验证宗门存在
	sect, err := s.sectRepo.GetByID(ctx, sectID)
	if err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "failed to query sect")
	}
	if sect == nil {
		return nil, errors.NewGameError(errors.ErrEntityNotFound, "sect not found")
	}

	// 检查是否已经加入
	existing, err := s.sectRepo.GetMember(ctx, sectID, string(entity.ID))
	if err == nil && existing {
		return &types.OperationResult{
			Success: false,
			Message: "已经在该宗门中",
		}, nil
	}

	if err := s.sectRepo.AddMember(ctx, sectID, string(entity.ID), "member"); err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "failed to join sect")
	}

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("成功加入宗门：%s！", sect.Name),
		Effects: map[string]interface{}{
			"sect_id":   sectID,
			"sect_name": sect.Name,
		},
	}, nil
}

// 发送消息
func (s *OperationService) executeSendMessage(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	content, ok := op.Params["content"].(string)
	if !ok || content == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "content is required")
	}

	msgType, _ := op.Params["message_type"].(string)
	if msgType == "" {
		msgType = "private"
	}

	receiverID, _ := op.Params["receiver_id"].(string)

	message := &types.DBMessage{
		ID:         generateMessageID(),
		SenderID:   string(entity.ID),
		ReceiverID: receiverID,
		Type:       msgType,
		Content:    content,
		IsRead:     false,
		CreatedAt:  time.Now().Unix(),
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	return &types.OperationResult{
		Success: true,
		Message: "消息发送成功！",
		Effects: map[string]interface{}{
			"message_id": message.ID,
			"type":       msgType,
		},
	}, nil
}

// 施法
func (s *OperationService) executeCastSpell(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	spellID, ok := op.Params["spell_id"].(string)
	if !ok {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "spell_id is required")
	}

	// 获取法术
	spell, err := s.spellRepo.GetByID(ctx, types.SpellID(spellID))
	if err != nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "spell not found")
	}

	// 检查是否已学习
	entitySpell, err := s.spellRepo.GetEntitySpell(ctx, entity.ID, spell.ID)
	if err != nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "spell not learned")
	}

	// 检查冷却
	if !entitySpell.CanCast(time.Now()) {
		return nil, errors.NewGameError(errors.ErrCooldownActive, "spell is on cooldown")
	}

	// 获取属性
	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)

	// 检查灵气
	if attr.Qi < float64(spell.Cost) {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "not enough qi")
	}

	attr.Qi -= float64(spell.Cost)

	// 计算效果
	damage := float64(spell.BaseDamage) * (1 + float64(attr.DivineSense)/100)

	// 更新冷却时间
	s.spellRepo.UpdateSpellCastTime(ctx, entity.ID, spell.ID)
	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("施放%s成功！", spell.Name),
		Effects: map[string]interface{}{
			"spell_name": spell.Name,
			"damage":     damage,
			"qi_cost":    spell.Cost,
		},
	}, nil
}
// 离开宗门
func (s *OperationService) executeLeaveSect(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	if s.sectRepo == nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "sect system not available")
	}
	return &types.OperationResult{
		Success: false,
		Message: "功能开发中",
	}, nil
}

// 添加好友
func (s *OperationService) executeAddFriend(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	name, ok := op.Params["name"].(string)
	if !ok || name == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "name is required")
	}
	target, err := s.entityRepo.GetByName(ctx, name)
	if err != nil {
		return nil, errors.NewGameError(errors.ErrEntityNotFound, "player not found")
	}
	if target == nil {
		return nil, errors.NewGameError(errors.ErrEntityNotFound, "player not found")
	}
	if target.ID == entity.ID {
		return &types.OperationResult{
			Success: false,
			Message: "不能添加自己为好友",
		}, nil
	}
	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("已向 %s 发送好友请求", name),
		Effects: map[string]interface{}{
			"target_id":   string(target.ID),
			"target_name": target.Name,
		},
	}, nil
}

// 删除好友
func (s *OperationService) executeRemoveFriend(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	friendID, ok := op.Params["friend_id"].(string)
	if !ok || friendID == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "friend_id is required")
	}
	return &types.OperationResult{
		Success: true,
		Message: "好友已删除",
		Effects: map[string]interface{}{
			"friend_id": friendID,
		},
	}, nil
}

// 接受好友请求
func (s *OperationService) executeAcceptFriend(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	requestID, ok := op.Params["request_id"].(string)
	if !ok || requestID == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "request_id is required")
	}
	return &types.OperationResult{
		Success: true,
		Message: "好友请求已接受",
		Effects: map[string]interface{}{
			"request_id": requestID,
		},
	}, nil
}

// 逃跑
func (s *OperationService) executeFlee(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	if entity.Status != types.StatusCombat {
		return &types.OperationResult{
			Success: false,
			Message: "当前不在战斗中",
		}, nil
	}
	entity.Status = types.StatusNormal
	entity.UpdatedAt = time.Now()
	s.entityRepo.Update(ctx, entity)
	return &types.OperationResult{
		Success: true,
		Message: "成功逃离战斗！",
		Effects: map[string]interface{}{},
	}, nil
}

// 使用技能
func (s *OperationService) executeUseSkill(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	if entity.Status != types.StatusCombat {
		return &types.OperationResult{
			Success: false,
			Message: "当前不在战斗中",
		}, nil
	}
	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)
	skillDamage := 10 + float64(attr.AttackPower)*1.2
	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)
	return &types.OperationResult{
		Success: true,
		Message: "使用技能攻击！",
		Effects: map[string]interface{}{
			"damage": skillDamage,
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

func generateMessageID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}

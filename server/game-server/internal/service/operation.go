package service

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"
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
	friendRepo    FriendRepository
	methodRepo    MethodRepository
}

type EntityRepository interface {
	GetByID(ctx context.Context, id types.EntityID) (*types.Entity, error)
	GetByName(ctx context.Context, name string) (*types.Entity, error)
	Create(ctx context.Context, entity *types.Entity) error
	Update(ctx context.Context, entity *types.Entity) error
	GetAttributes(ctx context.Context, entityID types.EntityID) (*types.Attributes, error)
	UpdateAttributes(ctx context.Context, entityID types.EntityID, attr *types.Attributes) error
	UpdateKarma(ctx context.Context, entityID types.EntityID, karma *types.Karma) error
	SetPasswordHash(ctx context.Context, entityID types.EntityID, hash string) error
	GetPasswordHash(ctx context.Context, entityID types.EntityID) (string, error)
	UpdateSpiritualRoots(ctx context.Context, entityID types.EntityID, roots []types.SpiritualRoot) error
	GetSpiritualRoots(ctx context.Context, entityID types.EntityID) ([]types.SpiritualRoot, error)
}

type ItemRepository interface {
	GetByID(ctx context.Context, id types.ItemID) (*types.Item, error)
	GetByName(ctx context.Context, name string) (*types.Item, error)
	Create(ctx context.Context, item *types.Item) error
}

type InventoryRepository interface {
	UnequipItem(ctx context.Context, entityID types.EntityID, slot string) error
	GetEquippedItems(ctx context.Context, entityID types.EntityID) ([]*types.InventoryItem, error)
	GetByEntityID(ctx context.Context, entityID types.EntityID) ([]*types.InventoryItem, error)
	GetItem(ctx context.Context, entityID types.EntityID, itemID types.ItemID) (*types.InventoryItem, error)
	AddItem(ctx context.Context, entityID types.EntityID, itemID types.ItemID, quantity int) error
	RemoveItem(ctx context.Context, entityID types.EntityID, itemID types.ItemID, quantity int) error
	EquipItem(ctx context.Context, entityID types.EntityID, itemID types.ItemID, slot string) error
	UpdateDurability(ctx context.Context, entityID types.EntityID, itemID types.ItemID, durability int) error
}

type SpellRepository interface {
	GetByID(ctx context.Context, id types.SpellID) (*types.Spell, error)
	GetEntitySpells(ctx context.Context, entityID types.EntityID) ([]*types.EntitySpell, error)
	GetEntitySpell(ctx context.Context, entityID types.EntityID, spellID types.SpellID) (*types.EntitySpell, error)
	UpdateSpellCastTime(ctx context.Context, entityID types.EntityID, spellID types.SpellID) error
	LearnSpell(ctx context.Context, entityID types.EntityID, spellID types.SpellID) error
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
	ListMembers(ctx context.Context, sectID string) ([]*SectMemberInfo, error)
}

type SectMemberInfo struct {
	EntityID     string  `json:"entity_id"`
	Name         string  `json:"name"`
	Rank         string  `json:"rank"`
	Contribution float64 `json:"contribution"`
	JoinedAt     int64   `json:"joined_at"`
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

type FriendRepository interface {
	AddFriend(ctx context.Context, entityID, friendID string) error
	RemoveFriend(ctx context.Context, entityID, friendID string) error
	AreFriends(ctx context.Context, entityID, friendID string) (bool, error)
	CreateRequest(ctx context.Context, fromID, toID string) (string, error)
	GetPendingRequest(ctx context.Context, fromID, toID string) (*FriendInfo, error)
	GetRequestByID(ctx context.Context, requestID string) (*FriendRequestInfo, error)
	AcceptRequest(ctx context.Context, requestID string) error
	GetFriends(ctx context.Context, entityID string) ([]*FriendshipInfo, error)
}

type FriendshipInfo struct {
	FriendID  string `json:"friend_id"`
	CreatedAt int64  `json:"created_at"`
}

type FriendInfo struct {
	ID        string
	FromID    string
	ToID      string
}

// ── 功法系统 ──

type MethodInfo struct {
	ID                    string
	Name                  string
	Quality               string
	RealmRequirement      string
	ElementAffinity       string
	CultivationMultiplier float64
	BreakthroughBonus     float64
	Description           string
}

type EntityMethodInfo struct {
	MethodID     string
	Method       *MethodInfo
	MasteryLevel float64
	IsMainMethod bool
	LearnedAt    int64
}

type MethodRepository interface {
	GetByID(ctx context.Context, id string) (*MethodInfo, error)
	GetByRealm(ctx context.Context, realm string) ([]*MethodInfo, error)
	GetEntityMethods(ctx context.Context, entityID types.EntityID) ([]*EntityMethodInfo, error)
	LearnMethod(ctx context.Context, entityID types.EntityID, methodID string) error
	SetMainMethod(ctx context.Context, entityID types.EntityID, methodID string) error
	GetMainMethod(ctx context.Context, entityID types.EntityID) (*EntityMethodInfo, error)
}

type FriendRequestInfo struct {
	ID     string
	FromID string
	ToID   string
	Status string
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
	friendRepo FriendRepository,
	methodRepo MethodRepository,
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
		friendRepo:    friendRepo,
		methodRepo:    methodRepo,
	}
}

func (s *OperationService) Execute(ctx context.Context, op *types.Operation) (*types.OperationResult, error) {
	entity, err := s.entityRepo.GetByID(ctx, op.ActorID)
	if err != nil {
		return nil, errors.ErrEntityNotFound_
	}

	if entity.Status == types.StatusDead {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "角色已死亡")
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
	case types.ActionUseItem:
		return s.executeUseItem(ctx, entity, op)
	case types.ActionDropItem:
		return s.executeDropItem(ctx, entity, op)
	case types.ActionEquipItem:
		return s.executeEquipItem(ctx, entity, op)
	case types.ActionUnequipItem:
		return s.executeUnequipItem(ctx, entity, op)
	case types.ActionLearnSpell:
		return s.executeLearnSpell(ctx, entity, op)
	case types.ActionListFriends:
		return s.executeListFriends(ctx, entity, op)
	case types.ActionSectInfo:
		return s.executeSectInfo(ctx, entity, op)
	case types.ActionLearnMethod:
		return s.executeLearnMethod(ctx, entity, op)
	case types.ActionSetMainMethod:
		return s.executeSetMainMethod(ctx, entity, op)
	case types.ActionListMethods:
		return s.executeListMethods(ctx, entity, op)
	default:
		return nil, errors.ErrInvalidOperationType
	}
}

// modifyKarma 增加实体业力值并持久化
func (s *OperationService) modifyKarma(ctx context.Context, entityID types.EntityID, delta int, reason string) {
	entity, err := s.entityRepo.GetByID(ctx, entityID)
	if err != nil || entity == nil {
		return
	}
	entity.Karma.KarmaValue += delta
	s.entityRepo.UpdateKarma(ctx, entityID, &entity.Karma)
}

// 修炼（使用 heavenly-dao 效率计算）
func (s *OperationService) executeCultivate(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	attr, err := s.entityRepo.GetAttributes(ctx, entity.ID)
	if err != nil {
		attr = &types.Attributes{}
	}

	cultivationGain := 0.1 * float64(attr.Comprehension) / 100.0

	// 尝试调用 heavenly-dao 获取修炼效率
	if s.daoService != nil {
		daoInput := &CultivateEfficiencyInput{
			EntityID:        string(entity.ID),
			Realm:           entity.Realm,
			SpiritualRoots:  attr.SpiritualRoots,
			SpiritualDensity: 0.5,
			Comprehension:   attr.Comprehension,
			MentalStability: attr.MentalStability,
			BaseLifespan:    attr.MaxLifespan,
			CurrentAge:      attr.Age,
		}
		daoResult, err := s.daoService.ExecuteCultivate(ctx, daoInput)
		if err == nil && daoResult != nil && daoResult.Success {
			cultivationGain = daoResult.CultivationGained
			attr.CultivationProgress += cultivationGain
		} else {
			// heavenly-dao 不可用时使用简化计算
			attr.CultivationProgress += cultivationGain
		}
	} else {
		attr.CultivationProgress += cultivationGain
	}

	if attr.CultivationProgress > 100 {
		attr.CultivationProgress = 100
	}

	if err := s.entityRepo.UpdateAttributes(ctx, entity.ID, attr); err != nil {
		return nil, err
	}

	entity.Status = types.StatusCultivating
	entity.UpdatedAt = time.Now()
	s.entityRepo.Update(ctx, entity)

	s.modifyKarma(ctx, entity.ID, 1, "修炼增长业力")

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

	// 验证区域存在
	region, err := s.worldService.GetRegion(ctx, regionID)
	if err != nil || region == nil {
		return nil, errors.NewGameError(errors.ErrRegionNotFound, "区域不存在")
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

	s.modifyKarma(ctx, entity.ID, 2, "打坐恢复")

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

	s.modifyKarma(ctx, entity.ID, 3, "休息恢复")

	return &types.OperationResult{
		Success: true,
		Message: "休息完成，状态已恢复",
		Effects: map[string]interface{}{
			"qi":              attr.Qi,
			"spiritual_power": attr.SpiritualPower,
		},
	}, nil
}

// 突破（使用 heavenly-dao 完整突破系统）
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

	// 查询主修功法品质
	methodQuality := 50.0
	if s.methodRepo != nil {
		mainMethod, err := s.methodRepo.GetMainMethod(ctx, entity.ID)
		if err == nil && mainMethod != nil && mainMethod.Method != nil {
			switch mainMethod.Method.Quality {
			case "天级":
				methodQuality = 90
			case "地级":
				methodQuality = 75
			case "玄级":
				methodQuality = 60
			default:
				methodQuality = 45
			}
			// 突破加成
			methodQuality += mainMethod.Method.BreakthroughBonus * 100
		}
	}

	// 使用 heavenly-dao 完整突破系统
	if s.daoService != nil {
		daoInput := &BreakthroughInput{
			EntityID:        string(entity.ID),
			CurrentRealm:    entity.Realm,
			TargetRealm:     newRealm,
			CultivationTime: attr.CultivationProgress,
			MethodQuality:   methodQuality,
			ResourceBonus:   0,
			MentalStability: attr.MentalStability,
			Luck:            attr.Luck,
			Karma:           entity.Karma.KarmaValue,
			Merit:           entity.Karma.Merit,
		}
		daoResult, err := s.daoService.ExecuteBreakthrough(ctx, daoInput)
		if err != nil {
			return nil, err
		}

		if !daoResult.Success {
			// 应用失败惩罚
			progressLoss := daoResult.PenaltyProgressLoss
			if progressLoss <= 0 {
				progressLoss = 0.1
			}
			attr.CultivationProgress *= (1 - progressLoss)
			attr.MentalStability -= daoResult.PenaltyMentalDamage
			if attr.MentalStability < 0 {
				attr.MentalStability = 0
			}

			s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)

			return &types.OperationResult{
				Success: false,
				Message: daoResult.Message,
				Effects: map[string]interface{}{
					"success_rate":       daoResult.SuccessRate,
					"progress_loss":      progressLoss,
					"mental_damage":      daoResult.PenaltyMentalDamage,
					"cooldown_hours":     daoResult.PenaltyCooldownHours,
				},
			}, nil
		}

		// 突破成功 - 处理天劫
		if daoResult.TribulationTriggered {
			return &types.OperationResult{
				Success: false,
				Message: fmt.Sprintf("天劫降临！强度: %.1f, 类型: %s, 突破失败", daoResult.TribulationStrength, daoResult.TribulationType),
				Effects: map[string]interface{}{
					"tribulation":        true,
					"tribulation_strength": daoResult.TribulationStrength,
					"tribulation_type":   daoResult.TribulationType,
				},
			}, nil
		}

		// 突破成功 - 应用属性变化
		entity.Realm = daoResult.NewRealm
		attr.CultivationProgress = 0
		attr.MaxQi *= 1.5
		attr.MaxSpiritualPower *= 1.5
		attr.MaxLifespan = getRealmLifespan(daoResult.NewRealm)

		s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)
		entity.UpdatedAt = time.Now()
		s.entityRepo.Update(ctx, entity)
		s.modifyKarma(ctx, entity.ID, 5, "突破成功")

		return &types.OperationResult{
			Success: true,
			Message: daoResult.Message,
			Effects: map[string]interface{}{
				"new_realm":           daoResult.NewRealm,
				"success_rate":        daoResult.SuccessRate,
				"max_qi":              attr.MaxQi,
				"max_spiritual":       attr.MaxSpiritualPower,
				"tribulation_triggered": daoResult.TribulationTriggered,
			},
		}, nil
	}

	// fallback: daoService 不可用时的简化突破
	successRate := 0.5 + float64(attr.Luck)/200.0
	if successRate > 0.8 {
		successRate = 0.8
	}

	if rand.Float64() > successRate {
		attr.CultivationProgress = 0
		s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)
		return &types.OperationResult{
			Success: false,
			Message: "突破失败，修为尽失！",
			Effects: map[string]interface{}{
				"success_rate": successRate,
			},
		}, nil
	}

	entity.Realm = newRealm
	attr.CultivationProgress = 0
	attr.MaxQi *= 1.5
	attr.MaxSpiritualPower *= 1.5
	attr.MaxLifespan = getRealmLifespan(newRealm)

	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)
	entity.UpdatedAt = time.Now()
	s.entityRepo.Update(ctx, entity)
	s.modifyKarma(ctx, entity.ID, 5, "突破成功")

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

// EquipmentBonuses 表示装备提供的完整属性加成
type EquipmentBonuses struct {
	AttackPower       float64
	Defense           float64
	Speed             float64
	CritRate          float64
	CritDamage        float64
	DodgeRate         float64
	HitRate           float64
	Penetration       float64
	DamageReduction   float64
	MaxQi             float64
	MaxSpiritualPower float64
	DivineSense       float64
	Comprehension     float64
	Constitution      float64
	Luck              float64
}

// 读取装备所有属性加成
func (s *OperationService) GetEquipmentBonuses(ctx context.Context, entityID types.EntityID) *EquipmentBonuses {
	bonuses := &EquipmentBonuses{}
	if s.inventoryRepo == nil {
		return bonuses
	}
	equipped, err := s.inventoryRepo.GetEquippedItems(ctx, entityID)
	if err != nil || len(equipped) == 0 {
		return bonuses
	}
	for _, invItem := range equipped {
		if invItem.Item == nil || invItem.Item.Attributes == nil {
			continue
		}
		attrs := invItem.Item.Attributes
		bonuses.AttackPower += getFloat64FromMap(attrs, "attack_power")
		bonuses.Defense += getFloat64FromMap(attrs, "defense")
		bonuses.Speed += getFloat64FromMap(attrs, "speed")
		bonuses.CritRate += getFloat64FromMap(attrs, "crit_rate")
		bonuses.CritDamage += getFloat64FromMap(attrs, "crit_damage")
		bonuses.DodgeRate += getFloat64FromMap(attrs, "dodge_rate")
		bonuses.HitRate += getFloat64FromMap(attrs, "hit_rate")
		bonuses.Penetration += getFloat64FromMap(attrs, "penetration")
		bonuses.DamageReduction += getFloat64FromMap(attrs, "damage_reduction")
		bonuses.MaxQi += getFloat64FromMap(attrs, "max_qi")
		bonuses.MaxSpiritualPower += getFloat64FromMap(attrs, "max_spiritual_power")
		bonuses.DivineSense += getFloat64FromMap(attrs, "divine_sense")
		bonuses.Comprehension += getFloat64FromMap(attrs, "comprehension")
		bonuses.Constitution += getFloat64FromMap(attrs, "constitution")
		bonuses.Luck += getFloat64FromMap(attrs, "luck")
	}
	return bonuses
}

func getFloat64FromMap(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		if f, ok2 := v.(float64); ok2 {
			return f
		}
	}
	return 0
}
	// reduceEquipmentDurability 减少装备耐久度，归零自动卸下
func (s *OperationService) reduceEquipmentDurability(ctx context.Context, entityID types.EntityID) {
	if s.inventoryRepo == nil {
		return
	}
	equipped, err := s.inventoryRepo.GetEquippedItems(ctx, entityID)
	if err != nil || len(equipped) == 0 {
		return
	}
	for _, invItem := range equipped {
		newDur := invItem.Durability - 1
		if newDur < 0 {
			newDur = 0
		}
		s.inventoryRepo.UpdateDurability(ctx, entityID, invItem.ItemID, newDur)
		if newDur <= 0 {
			s.inventoryRepo.UnequipItem(ctx, entityID, invItem.Slot)
		}
	}
}

// 战斗
func (s *OperationService) executeCombat(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	targetID, ok := op.Params["target_id"].(string)
	if !ok {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少目标ID")
	}

	// 获取目标
	target, err := s.entityRepo.GetByID(ctx, types.EntityID(targetID))
	if err != nil {
		return nil, errors.NewGameError(errors.ErrEntityNotFound, "目标不存在")
	}

	// 检查目标状态
	if target.Status == types.StatusDead {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "目标已死亡")
	}

	// 检查是否在同一区域
	if entity.Position.RegionID != target.Position.RegionID {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "目标不在同一区域")
	}

	// 获取属性
	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)
	targetAttr, _ := s.entityRepo.GetAttributes(ctx, target.ID)

	// 应用完整装备属性加成
	bonuses := s.GetEquipmentBonuses(ctx, entity.ID)
	tgtBonuses := s.GetEquipmentBonuses(ctx, target.ID)

	attr.AttackPower += bonuses.AttackPower
	attr.Defense += bonuses.Defense
	attr.Speed += bonuses.Speed
	attr.CritRate += bonuses.CritRate
	attr.CritDamage += bonuses.CritDamage
	attr.DodgeRate += bonuses.DodgeRate
	attr.HitRate += bonuses.HitRate
	attr.Penetration += bonuses.Penetration
	attr.DamageReduction += bonuses.DamageReduction
	targetAttr.AttackPower += tgtBonuses.AttackPower
	targetAttr.Defense += tgtBonuses.Defense
	targetAttr.Speed += tgtBonuses.Speed
	targetAttr.CritRate += tgtBonuses.CritRate
	targetAttr.CritDamage += tgtBonuses.CritDamage
	targetAttr.DodgeRate += tgtBonuses.DodgeRate
	targetAttr.HitRate += tgtBonuses.HitRate
	targetAttr.Penetration += tgtBonuses.Penetration
	targetAttr.DamageReduction += tgtBonuses.DamageReduction

	// 计算距离
	distance := math.Sqrt(math.Pow(entity.Position.X-target.Position.X, 2) + math.Pow(entity.Position.Y-target.Position.Y, 2))
	if distance > 10 {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "目标太远")
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

	karmaChange := -2
	if result.Killed {
		karmaChange = -10
		// NPC 击杀掉落
		if target.EntityType == types.EntityTypeNPC {
			drops := s.handleNPCDrops(ctx, entity.ID, target)
			result.Message += " " + drops
			karmaChange = -5
		}
	}
	s.modifyKarma(ctx, entity.ID, karmaChange, "战斗")
	// 减少装备耐久度
	s.reduceEquipmentDurability(ctx, entity.ID)
	s.reduceEquipmentDurability(ctx, target.ID)

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

// handleNPCDrops generates loot drops when an NPC is killed.
func (s *OperationService) handleNPCDrops(ctx context.Context, entityID types.EntityID, npc *types.Entity) string {
	realmLevel := types.CultivationRealmLevel(npc.Realm)
	drops := []struct {
		Name     string
		ItemType types.ItemType
		Rarity   int
		Chance   float64
	}{
		// 基础掉落
		{"妖兽精血", types.ItemTypeMaterial, 1, 0.60},
		{"妖兽皮毛", types.ItemTypeMaterial, 1, 0.50},
		{"妖兽骨", types.ItemTypeMaterial, 1, 0.30},
		{"灵石", types.ItemTypeMaterial, 1, 0.80},
	}

	// 境界>=3 (金丹) 额外掉落
	if realmLevel >= 3 {
		drops = append(drops,
			struct {
				Name     string
				ItemType types.ItemType
				Rarity   int
				Chance   float64
			}{"内丹", types.ItemTypeMaterial, 2, 0.30},
			struct {
				Name     string
				ItemType types.ItemType
				Rarity   int
				Chance   float64
			}{"妖丹", types.ItemTypeMaterial, 3, 0.20},
		)
	}
	// 境界>=4 (元婴) 额外掉落
	if realmLevel >= 4 {
		drops = append(drops,
			struct {
				Name     string
				ItemType types.ItemType
				Rarity   int
				Chance   float64
			}{"妖核", types.ItemTypeMaterial, 4, 0.10},
			struct {
				Name     string
				ItemType types.ItemType
				Rarity   int
				Chance   float64
			}{"妖兽材料", types.ItemTypeMaterial, 3, 0.25},
		)
	}
	// 境界>=5 (化神) 额外掉落
	if realmLevel >= 5 {
		drops = append(drops, struct {
			Name     string
			ItemType types.ItemType
			Rarity   int
			Chance   float64
		}{"妖兽精魄", types.ItemTypeMaterial, 5, 0.05})
	}

	var dropped []string
	for _, d := range drops {
		if rand.Float64() < d.Chance {
			s.ensureItemInInventory(ctx, entityID, d.Name, d.ItemType, d.Rarity)
			dropped = append(dropped, d.Name)
		}
	}

	if len(dropped) == 0 {
		return "未掉落任何物品"
	}
	return "掉落: " + strings.Join(dropped, ", ")
}

// 探索
func (s *OperationService) executeExplore(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)

	// 获取区域信息
	region, err := s.worldService.GetRegion(ctx, entity.Position.RegionID)
	if err != nil {
		return nil, errors.NewGameError(errors.ErrRegionNotFound, "区域不存在")
	}

	// 消耗灵气
	qiCost := float64(region.SpiritualTier) * 5
	if attr.Qi < qiCost {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "灵力不足")
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
		s.ensureItemInInventory(ctx, entity.ID, resource.Name, types.ItemTypeMaterial, resource.Rarity)
	}

	// 奇遇事件
	if rand.Float64() < 0.1 {
		discoveries = append(discoveries, "触发奇遇！")
	}

	entity.Status = types.StatusExploring
	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)
	s.entityRepo.Update(ctx, entity)

	s.modifyKarma(ctx, entity.ID, 1, "探索")

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
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少资源类型")
	}

	quantity, _ := op.Params["quantity"].(float64)
	if quantity <= 0 {
		quantity = 1
	}

	// 获取区域信息
	region, err := s.worldService.GetRegion(ctx, entity.Position.RegionID)
	if err != nil {
		return nil, errors.NewGameError(errors.ErrRegionNotFound, "区域不存在")
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
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "该区域没有此资源")
	}

	if targetResource.Quantity < int(quantity) {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "资源不足")
	}

	// 获取属性
	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)

	// 消耗灵气
	qiCost := float64(targetResource.Rarity) * 5 * quantity
	if attr.Qi < qiCost {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "灵力不足")
	}
	attr.Qi -= qiCost

	// 增加采集技能经验
	attr.HerbKnowledge += int(quantity)

	// 创建物品（简化处理，实际应该根据资源类型创建对应物品）
	// 这里假设资源直接作为物品添加到背包

	s.ensureItemInInventory(ctx, entity.ID, targetResource.Name, types.ItemTypeMaterial, targetResource.Rarity)
	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)

	// 刷新资源
	s.worldService.SpawnResources(ctx, string(region.ID))

	s.modifyKarma(ctx, entity.ID, 1, "采集")

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
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少配方ID")
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
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "灵力不足")
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

	s.ensureItemInInventory(ctx, entity.ID, "crafted_"+recipeID, types.ItemTypePill, recipeDifficulty)

	s.modifyKarma(ctx, entity.ID, 2, "炼制成功")

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
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "需要元婴期或更高境界")
	}

	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)

	// 验证悟性
	if attr.Comprehension < 80 {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "悟性需要达到80")
	}

	// 消耗资源
	qiCost := attr.MaxQi * 0.5
	spCost := float64(attr.DivineSense) * 0.5

	if attr.Qi < qiCost {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "灵力不足")
	}
	if attr.SpiritualPower < spCost {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "神识不足")
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
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少目标ID")
	}

	itemID, ok := op.Params["item_id"].(string)
	if !ok {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少物品ID")
	}

	price, ok := op.Params["price"].(float64)
	if !ok || price <= 0 {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "价格无效")
	}

	// 获取目标
	target, err := s.entityRepo.GetByID(ctx, types.EntityID(targetID))
	if err != nil {
		return nil, errors.NewGameError(errors.ErrEntityNotFound, "目标不存在")
	}

	// 检查是否在同一区域
	if entity.Position.RegionID != target.Position.RegionID {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "目标不在同一区域")
	}

	// 获取属性
	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)
	targetAttr, _ := s.entityRepo.GetAttributes(ctx, target.ID)

	// 检查灵石
	if targetAttr.SpiritStones.LowGrade < int64(price) {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "目标灵石不足")
	}

	// 执行交易（简化处理）
	targetAttr.SpiritStones.LowGrade -= int64(price)
	attr.SpiritStones.LowGrade += int64(price)

	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)
	s.entityRepo.UpdateAttributes(ctx, target.ID, targetAttr)

	s.modifyKarma(ctx, entity.ID, 2, "交易")

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
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少宗门名称")
	}

	// 验证境界
	if types.CultivationRealmLevel(entity.Realm) < types.CultivationRealmLevel(types.RealmGoldenCore) {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "需要金丹期或更高境界")
	}

	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)

	// 检查费用
	cost := int64(10000)
	if attr.SpiritStones.LowGrade < cost {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "灵石不足（需要10000）")
	}

	attr.SpiritStones.LowGrade -= cost
	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)

	// 持久化宗门到数据库
	sectID := string(types.GenerateEntityID())
	if err := s.sectRepo.Create(ctx, sectID, sectName, string(entity.ID)); err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "创建宗门失败")
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
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少宗门ID")
	}

	// 验证宗门存在
	sect, err := s.sectRepo.GetByID(ctx, sectID)
	if err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "查询宗门失败")
	}
	if sect == nil {
		return nil, errors.NewGameError(errors.ErrEntityNotFound, "宗门不存在")
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
		return nil, errors.NewGameError(errors.ErrInternalError, "加入宗门失败")
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
		return nil, errors.NewGameError(errors.ErrInvalidParams, "请输入内容")
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

// 施法（含装备/灵根加成 + 实际伤害目标）
func (s *OperationService) executeCastSpell(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	spellID, ok := op.Params["spell_id"].(string)
	if !ok {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少法术ID")
	}

	// 获取法术
	spell, err := s.spellRepo.GetByID(ctx, types.SpellID(spellID))
	if err != nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "法术不存在")
	}

	// 检查是否已学习
	entitySpell, err := s.spellRepo.GetEntitySpell(ctx, entity.ID, spell.ID)
	if err != nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "未学习此法术")
	}

	// 检查冷却
	if !entitySpell.CanCast(time.Now()) {
		return nil, errors.NewGameError(errors.ErrCooldownActive, "法术冷却中")
	}

	// 获取属性
	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)

	// 检查灵气
	if attr.Qi < float64(spell.Cost) {
		return nil, errors.NewGameError(errors.ErrInsufficientResources, "灵力不足")
	}

	attr.Qi -= float64(spell.Cost)

	// 装备法术加成
	eq := s.GetEquipmentBonuses(ctx, entity.ID)

	// 灵根元素匹配加成
	elementBonus := 1.0
	spellElement := string(spell.Element)
	for _, root := range attr.SpiritualRoots {
		if root.Element == spellElement {
			elementBonus = 1.5
			break
		}
	}

	// 计算伤害：基础 * 神识 * 灵根加成 * 装备加成
	damage := float64(spell.BaseDamage) * (1 + float64(attr.DivineSense)/100) * elementBonus
	damage *= (1 + eq.AttackPower/100) // 装备攻击%叠乘

	// 如果法术是治疗类型
	healAmount := float64(spell.BaseHeal)

	// 处理目标
	targetID, hasTarget := op.Params["target_id"].(string)
	if hasTarget && targetID != "" {
		target, err := s.entityRepo.GetByID(ctx, types.EntityID(targetID))
		if err == nil && target.Status != types.StatusDead {
			targetAttr, _ := s.entityRepo.GetAttributes(ctx, target.ID)

			if spell.Type == types.SpellTypeAttack {
				targetAttr.Qi -= damage
				if targetAttr.Qi <= 0 {
					targetAttr.Qi = 0
					target.Status = types.StatusDead
				}
				s.entityRepo.UpdateAttributes(ctx, target.ID, targetAttr)
				s.entityRepo.Update(ctx, target)
			} else if spell.Type == types.SpellTypeHeal {
				targetAttr.Qi += healAmount
				if targetAttr.Qi > targetAttr.MaxQi {
					targetAttr.Qi = targetAttr.MaxQi
				}
				s.entityRepo.UpdateAttributes(ctx, target.ID, targetAttr)
			}
		}
	}

	// 更新冷却时间
	s.spellRepo.UpdateSpellCastTime(ctx, entity.ID, spell.ID)
	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)

	s.modifyKarma(ctx, entity.ID, -1, "施法")

	effects := map[string]interface{}{
		"spell_name":    spell.Name,
		"damage":        damage,
		"qi_cost":       spell.Cost,
		"element_match": elementBonus > 1.0,
	}
	if hasTarget && targetID != "" {
		effects["target_id"] = targetID
	}

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("施放%s成功！", spell.Name),
		Effects: effects,
	}, nil
}
// 离开宗门
func (s *OperationService) executeLeaveSect(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	if s.sectRepo == nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "宗门系统不可用")
	}

	sectID, ok := op.Params["sect_id"].(string)
	if !ok || sectID == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少宗门ID")
	}

	// Verify the entity is a member
	existing, err := s.sectRepo.GetMember(ctx, sectID, string(entity.ID))
	if err != nil || !existing {
		return &types.OperationResult{
			Success: false,
			Message: "您不在此宗门中",
		}, nil
	}

	if err := s.sectRepo.RemoveMember(ctx, sectID, string(entity.ID)); err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "退出宗门失败")
	}

	return &types.OperationResult{
		Success: true,
		Message: "成功退出宗门",
		Effects: map[string]interface{}{
			"sect_id": sectID,
		},
	}, nil
}

// 添加好友
func (s *OperationService) executeAddFriend(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	name, ok := op.Params["name"].(string)
	if !ok || name == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "请输入名称")
	}
	target, err := s.entityRepo.GetByName(ctx, name)
	if err != nil {
		return nil, errors.NewGameError(errors.ErrEntityNotFound, "玩家不存在")
	}
	if target == nil {
		return nil, errors.NewGameError(errors.ErrEntityNotFound, "玩家不存在")
	}
	if target.ID == entity.ID {
		return &types.OperationResult{
			Success: false,
			Message: "不能添加自己为好友",
		}, nil
	}

	// Check if already friends
	if s.friendRepo != nil {
		areFriends, _ := s.friendRepo.AreFriends(ctx, string(entity.ID), string(target.ID))
		if areFriends {
			return &types.OperationResult{
				Success: false,
				Message: "已经是好友了",
			}, nil
		}

		// Check for existing pending request
		existingReq, _ := s.friendRepo.GetPendingRequest(ctx, string(entity.ID), string(target.ID))
		if existingReq != nil {
			return &types.OperationResult{
				Success: false,
				Message: "已经发送过好友请求了",
			}, nil
		}

		// Create friend request
		requestID, err := s.friendRepo.CreateRequest(ctx, string(entity.ID), string(target.ID))
		if err != nil {
			return nil, errors.NewGameError(errors.ErrInternalError, "创建好友请求失败")
		}

		return &types.OperationResult{
			Success: true,
			Message: fmt.Sprintf("已向 %s 发送好友请求", name),
			Effects: map[string]interface{}{
				"target_id":   string(target.ID),
				"target_name": target.Name,
				"request_id":  requestID,
			},
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
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少好友ID")
	}

	if s.friendRepo != nil {
		s.friendRepo.RemoveFriend(ctx, string(entity.ID), friendID)
		s.friendRepo.RemoveFriend(ctx, friendID, string(entity.ID))
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
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少请求ID")
	}

	if s.friendRepo != nil {
		// Get request details
		fr, err := s.friendRepo.GetRequestByID(ctx, requestID)
		if err != nil {
			return nil, errors.NewGameError(errors.ErrInternalError, "查询好友请求失败")
		}
		if fr == nil {
			return &types.OperationResult{
				Success: false,
				Message: "好友请求不存在",
			}, nil
		}

		// Verify it's for this entity
		if fr.ToID != string(entity.ID) {
			return &types.OperationResult{
				Success: false,
				Message: "这不是发给您的好友请求",
			}, nil
		}

		// Accept the request
		if err := s.friendRepo.AcceptRequest(ctx, requestID); err != nil {
			return nil, errors.NewGameError(errors.ErrInternalError, "接受好友请求失败")
		}

		// Create bidirectional friendship
		s.friendRepo.AddFriend(ctx, fr.FromID, fr.ToID)
		s.friendRepo.AddFriend(ctx, fr.ToID, fr.FromID)
	}

	return &types.OperationResult{
		Success: true,
		Message: "好友请求已接受",
		Effects: map[string]interface{}{
			"request_id": requestID,
		},
	}, nil
}

// 逃跑（同时恢复目标战斗状态）
func (s *OperationService) executeFlee(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	if entity.Status != types.StatusCombat {
		return &types.OperationResult{
			Success: false,
			Message: "当前不在战斗中",
		}, nil
	}

	// 同时恢复目标战斗状态
	targetID, ok := op.Params["target_id"].(string)
	if ok && targetID != "" {
		target, err := s.entityRepo.GetByID(ctx, types.EntityID(targetID))
		if err == nil && target != nil && target.Status == types.StatusCombat {
			target.Status = types.StatusNormal
			target.UpdatedAt = time.Now()
			s.entityRepo.Update(ctx, target)
		}
	}

	entity.Status = types.StatusNormal
	entity.UpdatedAt = time.Now()
	s.entityRepo.Update(ctx, entity)
	s.modifyKarma(ctx, entity.ID, 1, "逃离战斗")
	return &types.OperationResult{
		Success: true,
		Message: "成功逃离战斗！",
		Effects: map[string]interface{}{},
	}, nil
}

// 使用技能（实际伤害目标）
func (s *OperationService) executeUseSkill(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	if entity.Status != types.StatusCombat {
		return &types.OperationResult{
			Success: false,
			Message: "当前不在战斗中",
		}, nil
	}

	targetID, ok := op.Params["target_id"].(string)
	if !ok || targetID == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少目标ID")
	}

	target, err := s.entityRepo.GetByID(ctx, types.EntityID(targetID))
	if err != nil {
		return nil, errors.NewGameError(errors.ErrEntityNotFound, "目标不存在")
	}

	if target.Status == types.StatusDead {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "目标已死亡")
	}

	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)
	targetAttr, _ := s.entityRepo.GetAttributes(ctx, target.ID)

	// 装备加成
	eq := s.GetEquipmentBonuses(ctx, entity.ID)

	skillDamage := 10 + (float64(attr.AttackPower)+eq.AttackPower)*1.2

	// 实际伤害目标
	targetAttr.Qi -= skillDamage
	killed := false
	if targetAttr.Qi <= 0 {
		targetAttr.Qi = 0
		target.Status = types.StatusDead
		killed = true
	}

	s.entityRepo.UpdateAttributes(ctx, target.ID, targetAttr)
	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)
	s.entityRepo.Update(ctx, entity)
	if killed {
		s.entityRepo.Update(ctx, target)
	}

	s.modifyKarma(ctx, entity.ID, -1, "使用技能")

	return &types.OperationResult{
		Success: true,
		Message: "使用技能攻击！",
		Effects: map[string]interface{}{
			"damage":    skillDamage,
			"target_id": targetID,
			"killed":    killed,
			"target_qi": targetAttr.Qi,
		},
	}, nil
}

func (s *OperationService) executeUseItem(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	itemName, ok := op.Params["item_name"].(string)
	if !ok || itemName == "" {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "请指定物品名")
	}

	item, err := s.itemRepo.GetByName(ctx, itemName)
	if err != nil || item == nil {
		return &types.OperationResult{
			Success: false,
			Message: "找不到该物品: " + itemName,
		}, nil
	}

	invItem, err := s.inventoryRepo.GetItem(ctx, entity.ID, item.ID)
	if err != nil || invItem == nil || invItem.Quantity <= 0 {
		return &types.OperationResult{
			Success: false,
			Message: "背包中没有该物品",
		}, nil
	}

	if !item.Usable {
		return &types.OperationResult{
			Success: false,
			Message: "该物品无法使用",
		}, nil
	}

	if item.RealmRequirement != "" && types.CultivationRealmLevel(entity.Realm) < types.CultivationRealmLevel(item.RealmRequirement) {
		return &types.OperationResult{
			Success: false,
			Message: "境界不足，无法使用该物品",
		}, nil
	}

	// Use the item: heal effects from attributes
	attr, _ := s.entityRepo.GetAttributes(ctx, entity.ID)
	effects := map[string]interface{}{
		"item_name":  item.Name,
	}

	if item.Attributes != nil {
		if qi, ok := item.Attributes["qi"]; ok {
			if v, ok := qi.(float64); ok {
				attr.Qi += v
				if attr.Qi > attr.MaxQi {
					attr.Qi = attr.MaxQi
				}
				effects["qi_recovery"] = v
			}
		}
		if sp, ok := item.Attributes["spiritual_power"]; ok {
			if v, ok := sp.(float64); ok {
				attr.SpiritualPower += v
				if attr.SpiritualPower > attr.MaxSpiritualPower {
					attr.SpiritualPower = attr.MaxSpiritualPower
				}
				effects["spiritual_recovery"] = v
			}
		}
	}

	s.entityRepo.UpdateAttributes(ctx, entity.ID, attr)
	s.inventoryRepo.RemoveItem(ctx, entity.ID, item.ID, 1)

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("使用 %s 成功", item.Name),
		Effects: effects,
	}, nil
}

func (s *OperationService) executeDropItem(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	itemID, ok := op.Params["item_id"].(string)
	if !ok || itemID == "" {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "请指定物品ID")
	}

	invItem, err := s.inventoryRepo.GetItem(ctx, entity.ID, types.ItemID(itemID))
	if err != nil || invItem == nil {
		return &types.OperationResult{
			Success: false,
			Message: "背包中没有该物品",
		}, nil
	}

	s.inventoryRepo.RemoveItem(ctx, entity.ID, types.ItemID(itemID), 1)

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("丢弃 %s 成功", itemID),
		Effects: map[string]interface{}{
			"item_id": itemID,
		},
	}, nil
}

// 装备物品
func (s *OperationService) executeEquipItem(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	itemID, ok := op.Params["item_id"].(string)
	if !ok || itemID == "" {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "请指定物品ID")
	}

	slot, _ := op.Params["slot"].(string)

	invItem, err := s.inventoryRepo.GetItem(ctx, entity.ID, types.ItemID(itemID))
	if err != nil || invItem == nil {
		return &types.OperationResult{
			Success: false,
			Message: "背包中没有该物品",
		}, nil
	}

	if invItem.Equipped {
		return &types.OperationResult{
			Success: false,
			Message: "该物品已装备",
		}, nil
	}

	if slot == "" {
		slot = itemTypeToSlot(string(invItem.Item.Type))
		if slot == "" {
			return &types.OperationResult{
				Success: false,
				Message: "无法确定装备位，请手动指定（weapon/armor/helmet/boots/necklace/ring）",
			}, nil
		}
	}

	if err := s.inventoryRepo.EquipItem(ctx, entity.ID, types.ItemID(itemID), slot); err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "装备失败")
	}

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("已装备 %s 到 %s 位", invItem.Item.Name, slot),
		Effects: map[string]interface{}{
			"item_name": invItem.Item.Name,
			"slot":      slot,
		},
	}, nil
}

// 卸下装备
func (s *OperationService) executeUnequipItem(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	slot, ok := op.Params["slot"].(string)
	if !ok || slot == "" {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "请指定装备位")
	}

	if err := s.inventoryRepo.UnequipItem(ctx, entity.ID, slot); err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "卸下装备失败")
	}

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("已卸下 %s 位的装备", slot),
		Effects: map[string]interface{}{
			"slot": slot,
		},
	}, nil
}

// 学习法术
func (s *OperationService) executeLearnSpell(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	spellID, ok := op.Params["spell_id"].(string)
	if !ok || spellID == "" {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "请指定法术ID")
	}

	spell, err := s.spellRepo.GetByID(ctx, types.SpellID(spellID))
	if err != nil {
		return &types.OperationResult{
			Success: false,
			Message: "法术不存在",
		}, nil
	}

	if spell.RealmRequirement != "" && types.CultivationRealmLevel(entity.Realm) < types.CultivationRealmLevel(spell.RealmRequirement) {
		return &types.OperationResult{
			Success: false,
			Message: fmt.Sprintf("境界不足，需要 %s", string(spell.RealmRequirement)),
		}, nil
	}

	existing, err := s.spellRepo.GetEntitySpell(ctx, entity.ID, types.SpellID(spellID))
	if err == nil && existing != nil {
		return &types.OperationResult{
			Success: false,
			Message: "已学习该法术",
		}, nil
	}

	if err := s.spellRepo.LearnSpell(ctx, entity.ID, types.SpellID(spellID)); err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "学习法术失败")
	}

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("成功学习法术：%s", spell.Name),
		Effects: map[string]interface{}{
			"spell_id":   spellID,
			"spell_name": spell.Name,
		},
	}, nil
}

// 好友列表
func (s *OperationService) executeListFriends(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	if s.friendRepo == nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "好友系统不可用")
	}

	friends, err := s.friendRepo.GetFriends(ctx, string(entity.ID))
	if err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "查询好友列表失败")
	}

	friendList := make([]map[string]interface{}, 0, len(friends))
	for _, f := range friends {
		friendEntity, err := s.entityRepo.GetByID(ctx, types.EntityID(f.FriendID))
		name := f.FriendID
		if err == nil && friendEntity != nil {
			name = friendEntity.Name
		}
		friendList = append(friendList, map[string]interface{}{
			"friend_id":   f.FriendID,
			"friend_name": name,
			"created_at":  f.CreatedAt,
		})
	}

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("共有 %d 位好友", len(friendList)),
		Effects: map[string]interface{}{
			"friends": friendList,
			"count":   len(friendList),
		},
	}, nil
}

// 宗门信息
func (s *OperationService) executeSectInfo(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	sectID, ok := op.Params["sect_id"].(string)
	if !ok || sectID == "" {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "请指定宗门ID")
	}

	sect, err := s.sectRepo.GetByID(ctx, sectID)
	if err != nil || sect == nil {
		return &types.OperationResult{
			Success: false,
			Message: "宗门不存在",
		}, nil
	}

	founderName := sect.FounderID
	if founderEntity, err := s.entityRepo.GetByID(ctx, types.EntityID(sect.FounderID)); err == nil {
		founderName = founderEntity.Name
	}

	members, _ := s.sectRepo.ListMembers(ctx, sectID)
	memberList := make([]map[string]interface{}, 0, len(members))
	for _, m := range members {
		memberList = append(memberList, map[string]interface{}{
			"entity_id":    m.EntityID,
			"name":         m.Name,
			"rank":         m.Rank,
			"contribution": m.Contribution,
			"joined_at":    m.JoinedAt,
		})
	}

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("宗门：%s", sect.Name),
		Effects: map[string]interface{}{
			"sect_id":      sect.ID,
			"sect_name":    sect.Name,
			"founder_id":   sect.FounderID,
			"founder_name": founderName,
			"alignment":    sect.Alignment,
			"members":      memberList,
			"member_count": len(memberList),
		},
	}, nil
}

// ── 功法系统操作 ──

// executeLearnMethod 学习功法
func (s *OperationService) executeLearnMethod(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	methodID, ok := op.Params["method_id"].(string)
	if !ok || methodID == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少功法ID")
	}

	if s.methodRepo == nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "功法系统不可用")
	}

	// 查询功法是否存在
	method, err := s.methodRepo.GetByID(ctx, methodID)
	if err != nil || method == nil {
		return &types.OperationResult{
			Success: false,
			Message: "功法不存在",
		}, nil
	}

	// 验证境界要求
	if method.RealmRequirement != "" && types.CultivationRealmLevel(entity.Realm) < types.CultivationRealmLevel(types.CultivationRealm(method.RealmRequirement)) {
		return &types.OperationResult{
			Success: false,
			Message: fmt.Sprintf("境界不足，需要 %s", method.RealmRequirement),
		}, nil
	}

	// 检查是否已学习
	entityMethods, _ := s.methodRepo.GetEntityMethods(ctx, entity.ID)
	for _, em := range entityMethods {
		if em.MethodID == methodID {
			return &types.OperationResult{
				Success: false,
				Message: "已学习该功法",
			}, nil
		}
	}

	// 学习功法
	if err := s.methodRepo.LearnMethod(ctx, entity.ID, methodID); err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "学习功法失败")
	}

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("成功学习功法：%s（%s）", method.Name, method.Quality),
		Effects: map[string]interface{}{
			"method_id":   methodID,
			"method_name": method.Name,
			"quality":     method.Quality,
		},
	}, nil
}

// executeSetMainMethod 设置主修功法
func (s *OperationService) executeSetMainMethod(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	methodID, ok := op.Params["method_id"].(string)
	if !ok || methodID == "" {
		return nil, errors.NewGameError(errors.ErrInvalidParams, "缺少功法ID")
	}

	if s.methodRepo == nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "功法系统不可用")
	}

	// 验证已学习
	entityMethods, err := s.methodRepo.GetEntityMethods(ctx, entity.ID)
	if err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "查询功法失败")
	}

	learned := false
	for _, em := range entityMethods {
		if em.MethodID == methodID {
			learned = true
			break
		}
	}
	if !learned {
		return &types.OperationResult{
			Success: false,
			Message: "未学习该功法",
		}, nil
	}

	method, _ := s.methodRepo.GetByID(ctx, methodID)

	// 设置为主修
	if err := s.methodRepo.SetMainMethod(ctx, entity.ID, methodID); err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "设置主修功法失败")
	}

	msg := "成功设置主修功法"
	if method != nil {
		msg = fmt.Sprintf("已将主修功法设为：%s（%s）", method.Name, method.Quality)
	}

	methodName := ""
	if method != nil {
		methodName = method.Name
	}

	return &types.OperationResult{
		Success: true,
		Message: msg,
		Effects: map[string]interface{}{
			"method_id":   methodID,
			"method_name": methodName,
		},
	}, nil
}

// executeListMethods 查看可学功法列表
func (s *OperationService) executeListMethods(ctx context.Context, entity *types.Entity, op *types.Operation) (*types.OperationResult, error) {
	if s.methodRepo == nil {
		return nil, errors.NewGameError(errors.ErrInvalidOperation, "功法系统不可用")
	}

	// 根据当前境界筛选可学功法
	methods, err := s.methodRepo.GetByRealm(ctx, string(entity.Realm))
	if err != nil {
		return nil, errors.NewGameError(errors.ErrInternalError, "查询功法失败")
	}

	// 已学习的功法
	entityMethods, _ := s.methodRepo.GetEntityMethods(ctx, entity.ID)
	learnedMap := make(map[string]bool)
	mainMethodID := ""
	for _, em := range entityMethods {
		learnedMap[em.MethodID] = true
		if em.IsMainMethod {
			mainMethodID = em.MethodID
		}
	}

	// 也获取更高境界的功法预览
	allRealms := []types.CultivationRealm{
		types.RealmMortal, types.RealmQiCondensation, types.RealmFoundation,
		types.RealmGoldenCore, types.RealmNascentSoul, types.RealmSoulTransform,
	}
	var previewMethods []*MethodInfo
	current := entity.Realm
	found := false
	for _, r := range allRealms {
		if r == current {
			found = true
			continue
		}
		if found {
			higherMethods, _ := s.methodRepo.GetByRealm(ctx, string(r))
			previewMethods = append(previewMethods, higherMethods...)
			if len(previewMethods) >= 3 {
				break
			}
		}
	}

	methodList := make([]map[string]interface{}, 0, len(methods))
	for _, m := range methods {
		isLearned := learnedMap[m.ID]
		isMain := m.ID == mainMethodID
		methodList = append(methodList, map[string]interface{}{
			"id":                     m.ID,
			"name":                   m.Name,
			"quality":                m.Quality,
			"element_affinity":       m.ElementAffinity,
			"cultivation_multiplier": m.CultivationMultiplier,
			"breakthrough_bonus":     m.BreakthroughBonus,
			"description":            m.Description,
			"learned":                isLearned,
			"is_main":                isMain,
		})
	}

	previewList := make([]map[string]interface{}, 0, len(previewMethods))
	for _, m := range previewMethods {
		previewList = append(previewList, map[string]interface{}{
			"id":      m.ID,
			"name":    m.Name,
			"quality": m.Quality,
		})
	}

	return &types.OperationResult{
		Success: true,
		Message: fmt.Sprintf("当前境界可学功法 %d 部", len(methodList)),
		Effects: map[string]interface{}{
			"methods":          methodList,
			"count":            len(methodList),
			"higher_preview":   previewList,
			"main_method_id":   mainMethodID,
		},
	}, nil
}

func itemTypeToSlot(itemType string) string {
	switch itemType {
	case "weapon":
		return "weapon"
	case "armor":
		return "armor"
	case "helmet":
		return "helmet"
	case "boots":
		return "boots"
	case "necklace":
		return "necklace"
	case "ring":
		return "ring1"
	}
	return ""
}
// ensureItemInInventory ensures a named material item exists in the DB and adds it to entity's inventory
func (s *OperationService) ensureItemInInventory(ctx context.Context, entityID types.EntityID, name string, itemType types.ItemType, rarity int) {
	item, err := s.itemRepo.GetByName(ctx, name)
	if err != nil || item == nil {
		item = &types.Item{
			ID:          types.ItemID(types.GenerateEntityID()),
			Name:        name,
			Type:        itemType,
			Rarity:      rarity,
			Description: "采集获得的材料",
			Attributes:  map[string]interface{}{},
			Stackable:   true,
			MaxStack:    99,
			Usable:      false,
		}
		s.itemRepo.Create(ctx, item)
	}
	s.inventoryRepo.AddItem(ctx, entityID, item.ID, 1)
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
	id := string(types.GenerateEntityID())
	// Format as UUID: 8-4-4-4-12
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		id[0:8], id[8:12], id[12:16], id[16:20], id[20:32])
}

package service

import (
	"context"
	"encoding/json"
	"math/rand"
	"sync"
	"time"

	cultivation "github.com/cultivation-world/shared/proto/pb"
	"github.com/cultivation-world/shared/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"golang.org/x/crypto/bcrypt"

	gameerrors "github.com/cultivation-world/shared/errors"
)

type GameService struct {
	cultivation.UnimplementedGameServiceServer
	entityRepo    EntityRepository
	operationSvc  *OperationService
	spellRepo     SpellRepository
	itemRepo      ItemRepository
	inventoryRepo InventoryRepository
	mu            sync.RWMutex
}

func NewGameService(entityRepo EntityRepository, operationSvc *OperationService, spellRepo SpellRepository, itemRepo ItemRepository, inventoryRepo InventoryRepository) *GameService {
	return &GameService{
		entityRepo:    entityRepo,
		operationSvc:  operationSvc,
		spellRepo:     spellRepo,
		itemRepo:      itemRepo,
		inventoryRepo: inventoryRepo,
	}
}

func (s *GameService) CreateEntity(ctx context.Context, req *cultivation.CreateEntityRequest) (*cultivation.CreateEntityResponse, error) {
	// Check if username already exists
	if existing, err := s.entityRepo.GetByName(ctx, req.Name); err == nil && existing != nil {
		return nil, status.Errorf(codes.AlreadyExists, "username already exists")
	}

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
			CritRate:            5,
			CritDamage:          150,
			DodgeRate:           5,
			HitRate:             90,
			Penetration:         0,
			DamageReduction:     0,
			MentalStability:     50,
			RemainingLifespan:   80,
			MaxLifespan:         80,
			AlchemyLevel:        1,
			ArtificingLevel:     1,
			FormationLevel:      1,
			FireControl:         1,
			HerbKnowledge:       1,
			MiningSkill:         1,
			TalismanSkill:       1,
			BeastTaming:         1,
			Reputation:          0,
			SectContribution:    0,
			DaoHeart:            10,
			Enlightenment:       5,
			RootPurity:          50,
			PoisonLevel:         0,
			CurseLevel:          0,
		},
		Karma: types.Karma{
			KarmaValue:   0,
			Merit:        0,
			KarmicDebt:   0,
			HeavenlyMark: "clear",
		},
		Status:    types.StatusNormal,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 生成灵根
	roots := generateSpiritualRoots()
	entity.Attributes.SpiritualRoots = roots
	// 计算平均纯度作为 root_purity
	if len(roots) > 0 {
		total := 0
		for _, r := range roots {
			total += r.Purity
		}
		entity.Attributes.RootPurity = total / len(roots)
	}

	if err := s.entityRepo.Create(ctx, entity); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create entity: %v", err)
	}

	if err := s.entityRepo.UpdateAttributes(ctx, entityID, &entity.Attributes); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save attributes: %v", err)
	}

	// 持久化灵根
	if len(roots) > 0 {
		if err := s.entityRepo.UpdateSpiritualRoots(ctx, entityID, roots); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to save spiritual roots: %v", err)
		}
	}

	// Hash and store password
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
		}
		if err := s.entityRepo.SetPasswordHash(ctx, entityID, string(hash)); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to store password: %v", err)
		}
	}

	return &cultivation.CreateEntityResponse{
		Entity: entityToProto(entity),
	}, nil
}

// generateSpiritualRoots generates 1-3 random spiritual roots for a new entity.
// Primary root: purity 60-90, secondary roots: purity 20-50.
// 5% chance of a mutated root (ice/lightning/poison).
func generateSpiritualRoots() []types.SpiritualRoot {
	elements := []string{"fire", "water", "wood", "metal", "earth"}
	rareElements := []string{"wind", "thunder", "ice"}

	count := rand.Intn(3) + 1 // 1-3 roots

	// Shuffle elements
	perm := rand.Perm(len(elements))
	roots := make([]types.SpiritualRoot, 0, count)

	for i := 0; i < count && i < len(elements); i++ {
		elem := elements[perm[i]]
		purity := 0

		if i == 0 {
			// Primary root: high purity
			purity = rand.Intn(31) + 60 // 60-90
			// 5% mutation chance on primary root
			if rand.Float64() < 0.05 {
				rareElem := rareElements[rand.Intn(len(rareElements))]
				elem = rareElem
			}
		} else {
			// Secondary roots: lower purity
			purity = rand.Intn(31) + 20 // 20-50
		}

		roots = append(roots, types.SpiritualRoot{
			Element: elem,
			Purity:  purity,
		})
	}

	return roots
}

func (s *GameService) AuthenticateEntity(ctx context.Context, req *cultivation.AuthRequest) (*cultivation.AuthResponse, error) {
	entity, err := s.entityRepo.GetByName(ctx, req.Username)
	if err != nil {
		return &cultivation.AuthResponse{
			Success: false,
		}, nil
	}

	// Verify password hash
	storedHash, err := s.entityRepo.GetPasswordHash(ctx, entity.ID)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)) != nil {
		return &cultivation.AuthResponse{
			Success: false,
		}, nil
	}

	return &cultivation.AuthResponse{
		Entity:  entityToProto(entity),
		Success: true,
	}, nil
}

func (s *GameService) ExecuteOperation(ctx context.Context, req *cultivation.OperationRequest) (*cultivation.OperationResponse, error) {
	params := make(map[string]interface{})
	for k, v := range req.Params {
		// JSON-decode values that were JSON-encoded by the gateway
		var decoded interface{}
		if err := json.Unmarshal([]byte(v), &decoded); err == nil {
			// Convert json.Number to float64 if needed (though Unmarshal uses float64 by default)
			params[k] = decoded
		} else {
			params[k] = v
		}
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
		msg := err.Error()
		if ge, ok := err.(*gameerrors.GameError); ok {
			msg = ge.Message
		}
		return &cultivation.OperationResponse{
			Success:   false,
			Message:   msg,
			Timestamp: time.Now().UnixNano(),
		}, nil
	}

	// JSON-encode effects to preserve types through protobuf map<string,string>
	effects := make(map[string]string)
	for k, v := range result.Effects {
		b, _ := json.Marshal(v)
		effects[k] = string(b)
	}

	return &cultivation.OperationResponse{
		Success:   result.Success,
		Message:   result.Message,
		Effects:   effects,
		Timestamp: result.Timestamp,
	}, nil
}

func (s *GameService) GetEntity(ctx context.Context, req *cultivation.EntityRequest) (*cultivation.EntityResponse, error) {
	entity, err := s.entityRepo.GetByID(ctx, types.EntityID(req.EntityId))
	if err != nil {
		return &cultivation.EntityResponse{
			Found: false,
		}, nil
	}

	proto := entityToProto(entity)
	s.populateSpellsAndItems(ctx, proto, entity.ID)
	s.applyEquipmentBonuses(ctx, proto)

	return &cultivation.EntityResponse{
		Entity: proto,
		Found:  true,
	}, nil
}

func (s *GameService) SyncState(ctx context.Context, req *cultivation.SyncRequest) (*cultivation.SyncResponse, error) {
	entity, err := s.entityRepo.GetByID(ctx, types.EntityID(req.EntityId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "entity not found")
	}

	proto := entityToProto(entity)
	s.populateSpellsAndItems(ctx, proto, entity.ID)

	s.applyEquipmentBonuses(ctx, proto)
	return &cultivation.SyncResponse{
		Entity:          proto,
		WorldTime:       time.Now().Unix(),
		NearbyEntityIds: []string{},
	}, nil
}

func (s *GameService) StreamEntityUpdates(req *cultivation.EntityStreamRequest, stream cultivation.GameService_StreamEntityUpdatesServer) error {
	return nil
}

func (s *GameService) populateSpellsAndItems(ctx context.Context, proto *cultivation.Entity, entityID types.EntityID) {
	// Fetch and populate entity spells
	if s.spellRepo != nil {
		entitySpells, err := s.spellRepo.GetEntitySpells(ctx, entityID)
		if err == nil {
			proto.Spells = make([]*cultivation.SpellData, 0, len(entitySpells))
			now := time.Now()
			for _, es := range entitySpells {
				sd := &cultivation.SpellData{
					SpellId:       string(es.SpellID),
					Proficiency:   int32(es.Proficiency),
				}
				if es.LastCastAt != nil {
					sd.LastCastAt = es.LastCastAt.Unix()
				}
				if es.Spell != nil {
					sd.Name = es.Spell.Name
					sd.SpellType = string(es.Spell.Type)
					sd.Element = string(es.Spell.Element)
					sd.Cost = int32(es.Spell.Cost)
					sd.BaseDamage = int32(es.Spell.BaseDamage)
					sd.BaseHeal = int32(es.Spell.BaseHeal)
					sd.Cooldown = int32(es.Spell.Cooldown)
					sd.Description = es.Spell.Description
					sd.RealmRequirement = string(es.Spell.RealmRequirement)
					sd.CooldownRemaining = int64(es.GetCooldownRemaining(now).Seconds())
				}
				proto.Spells = append(proto.Spells, sd)
			}
		}
	}

	// Fetch and populate inventory items
	if s.inventoryRepo != nil {
		invItems, err := s.inventoryRepo.GetByEntityID(ctx, entityID)
		if err == nil {
			proto.Items = make([]*cultivation.ItemData, 0, len(invItems))
			for _, inv := range invItems {
				id := &cultivation.ItemData{
					InventoryId: inv.ID,
					ItemId:      string(inv.ItemID),
					Quantity:    int32(inv.Quantity),
					Equipped:    inv.Equipped,
					Slot:        inv.Slot,
					Durability:  int32(inv.Durability),
					Bound:       inv.Bound,
				}
				if inv.Item != nil {
					id.Name = inv.Item.Name
					id.ItemType = string(inv.Item.Type)
					id.Rarity = int32(inv.Item.Rarity)
					id.Description = inv.Item.Description
					if len(inv.Item.Attributes) > 0 {
					attrBytes, _ := json.Marshal(inv.Item.Attributes)
					id.AttributesJson = string(attrBytes)
					}
				}
				proto.Items = append(proto.Items, id)
			}
		}
	}
}

	// applyEquipmentBonuses 将装备属性加成合并到角色面板显示（不持久化）
func (s *GameService) applyEquipmentBonuses(ctx context.Context, proto *cultivation.Entity) {
	if s.inventoryRepo == nil || proto == nil || proto.Attributes == nil {
		return
	}
	equipped, err := s.inventoryRepo.GetEquippedItems(ctx, types.EntityID(proto.Id))
	if err != nil || len(equipped) == 0 {
		return
	}
	pa := proto.Attributes
	for _, invItem := range equipped {
		if invItem.Item == nil || invItem.Item.Attributes == nil {
			continue
		}
		attrs := invItem.Item.Attributes
		if v, ok := attrs["attack_power"].(float64); ok { pa.AttackPower += float64(v) }
		if v, ok := attrs["defense"].(float64); ok { pa.Defense += float64(v) }
		if v, ok := attrs["speed"].(float64); ok { pa.Speed += float64(v) }
		if v, ok := attrs["max_qi"].(float64); ok { pa.MaxQi += float64(v) }
		if v, ok := attrs["max_spiritual_power"].(float64); ok { pa.MaxSpiritualPower += float64(v) }
		if v, ok := attrs["crit_rate"].(float64); ok { pa.CritRate += float64(v) }
		if v, ok := attrs["crit_damage"].(float64); ok { pa.CritDamage += float64(v) }
		if v, ok := attrs["dodge_rate"].(float64); ok { pa.DodgeRate += float64(v) }
		if v, ok := attrs["hit_rate"].(float64); ok { pa.HitRate += float64(v) }
		if v, ok := attrs["penetration"].(float64); ok { pa.Penetration += float64(v) }
		if v, ok := attrs["damage_reduction"].(float64); ok { pa.DamageReduction += float64(v) }
		if v, ok := attrs["divine_sense"].(float64); ok { pa.DivineSense += float64(v) }
		if v, ok := attrs["comprehension"].(float64); ok { pa.Comprehension += int32(v) }
		if v, ok := attrs["constitution"].(float64); ok { pa.Constitution += int32(v) }
		if v, ok := attrs["luck"].(float64); ok { pa.Luck += int32(v) }
	}
}
func entityToProto(e *types.Entity) *cultivation.Entity {
	if e == nil {
		return nil
	}

	return &cultivation.Entity{
		Id:         string(e.ID),
		EntityType: string(e.EntityType),
		Name:       e.Name,
		Realm:      string(e.Realm),
		Position: &cultivation.WorldPosition{
			RegionId: e.Position.RegionID,
			X:        e.Position.X,
			Y:        e.Position.Y,
		},
		Attributes: &cultivation.Attributes{
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
			CritRate:            e.Attributes.CritRate,
			CritDamage:          e.Attributes.CritDamage,
			DodgeRate:           e.Attributes.DodgeRate,
			HitRate:             e.Attributes.HitRate,
			Penetration:         e.Attributes.Penetration,
			DamageReduction:     e.Attributes.DamageReduction,
			AlchemyLevel:        int32(e.Attributes.AlchemyLevel),
			ArtificingLevel:     int32(e.Attributes.ArtificingLevel),
			FormationLevel:      int32(e.Attributes.FormationLevel),
			FireControl:         int32(e.Attributes.FireControl),
			HerbKnowledge:       int32(e.Attributes.HerbKnowledge),
			MiningSkill:         int32(e.Attributes.MiningSkill),
			TalismanSkill:       int32(e.Attributes.TalismanSkill),
			BeastTaming:         int32(e.Attributes.BeastTaming),
			Reputation:          int32(e.Attributes.Reputation),
			SectContribution:    int32(e.Attributes.SectContribution),
			DaoHeart:            int32(e.Attributes.DaoHeart),
			Enlightenment:       int32(e.Attributes.Enlightenment),
			RootPurity:          int32(e.Attributes.RootPurity),
			PoisonLevel:         int32(e.Attributes.PoisonLevel),
			CurseLevel:          int32(e.Attributes.CurseLevel),
			SpiritualRoots: spiritualRootsToProto(e.Attributes.SpiritualRoots),
			RootAwakened:        e.Attributes.RootAwakened,
			MutatedRoot:         e.Attributes.MutatedRoot,
		},
		Karma: &cultivation.Karma{
			KarmaValue:   int32(e.Karma.KarmaValue),
			Merit:        int32(e.Karma.Merit),
			KarmicDebt:   int32(e.Karma.KarmicDebt),
			HeavenlyMark: e.Karma.HeavenlyMark,
		},
		SpiritStones: &cultivation.SpiritStones{
			LowGrade:     e.Attributes.SpiritStones.LowGrade,
			MediumGrade:  e.Attributes.SpiritStones.MediumGrade,
			HighGrade:    e.Attributes.SpiritStones.HighGrade,
			PremiumGrade: e.Attributes.SpiritStones.PremiumGrade,
		},
		Status:    string(e.Status),
		CreatedAt: e.CreatedAt.Unix(),
		UpdatedAt: e.UpdatedAt.Unix(),
	}
}

func spiritualRootsToProto(roots []types.SpiritualRoot) []*cultivation.SpiritualRoot {
	if len(roots) == 0 {
		return nil
	}
	result := make([]*cultivation.SpiritualRoot, 0, len(roots))
	for _, r := range roots {
		result = append(result, &cultivation.SpiritualRoot{
			Element: r.Element,
			Purity:  int32(r.Purity),
		})
	}
	return result
}

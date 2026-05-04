package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cultivation "github.com/cultivation-world/shared/proto/pb"
	"github.com/cultivation-world/shared/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GameServiceClient struct {
	conn       *grpc.ClientConn
	gameClient cultivation.GameServiceClient
}

func NewGameServiceClient(host string, port int) (*GameServiceClient, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GameServiceClient{
		conn:       conn,
		gameClient: cultivation.NewGameServiceClient(conn),
	}, nil
}

func (c *GameServiceClient) Close() error {
	return c.conn.Close()
}

func (c *GameServiceClient) ExecuteOperation(op *types.Operation) (*types.OperationResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// JSON-encode params to preserve types through protobuf map<string,string>
	params := make(map[string]string)
	for k, v := range op.Params {
		b, _ := json.Marshal(v)
		params[k] = string(b)
	}

	resp, err := c.gameClient.ExecuteOperation(ctx, &cultivation.OperationRequest{
		OperationId: op.ID,
		ActorId:     string(op.ActorID),
		ActionType:  string(op.ActionType),
		Params:      params,
		Timestamp:   op.Timestamp,
	})

	if err != nil {
		return nil, err
	}

	// JSON-decode effects back to proper types
	effects := make(map[string]interface{})
	for k, v := range resp.Effects {
		var decoded interface{}
		if err := json.Unmarshal([]byte(v), &decoded); err == nil {
			effects[k] = decoded
		} else {
			effects[k] = v
		}
	}

	return &types.OperationResult{
		Success:   resp.Success,
		Message:   resp.Message,
		Effects:   effects,
		Timestamp: resp.Timestamp,
	}, nil
}

func (c *GameServiceClient) CreateEntity(ctx context.Context, username, password string, entityType types.EntityType) (*types.Entity, error) {
	resp, err := c.gameClient.CreateEntity(ctx, &cultivation.CreateEntityRequest{
		Name:       username,
		EntityType: string(entityType),
		Password:   password,
	})

	if err != nil {
		return nil, err
	}

	return protoToEntity(resp.Entity), nil
}

func (c *GameServiceClient) AuthenticateEntity(ctx context.Context, username, password string) (*types.Entity, error) {
	resp, err := c.gameClient.AuthenticateEntity(ctx, &cultivation.AuthRequest{
		Username: username,
		Password: password,
	})

	if err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("invalid credentials")
	}

	return protoToEntity(resp.Entity), nil
}

func (c *GameServiceClient) GetEntity(ctx context.Context, entityID types.EntityID) (*types.Entity, error) {
	resp, err := c.gameClient.GetEntity(ctx, &cultivation.EntityRequest{
		EntityId: string(entityID),
	})

	if err != nil {
		return nil, err
	}

	return protoToEntity(resp.Entity), nil
}

// GetEntityExtras returns entity, spells (as maps), and items (as maps) for state sync.
func (c *GameServiceClient) GetEntityExtras(ctx context.Context, entityID types.EntityID) (*types.Entity, []map[string]interface{}, []map[string]interface{}, error) {
	resp, err := c.gameClient.GetEntity(ctx, &cultivation.EntityRequest{
		EntityId: string(entityID),
	})

	if err != nil {
		return nil, nil, nil, err
	}

	entity := protoToEntity(resp.Entity)

	// Convert spells to maps
	spells := make([]map[string]interface{}, 0, len(resp.Entity.Spells))
	for _, s := range resp.Entity.Spells {
		m := map[string]interface{}{
			"spell_id":           s.SpellId,
			"name":               s.Name,
			"spell_type":         s.SpellType,
			"element":            s.Element,
			"cost":               s.Cost,
			"base_damage":        s.BaseDamage,
			"base_heal":          s.BaseHeal,
			"cooldown":           s.Cooldown,
			"description":        s.Description,
			"realm_requirement":  s.RealmRequirement,
			"proficiency":        s.Proficiency,
			"last_cast_at":       s.LastCastAt,
			"cooldown_remaining": s.CooldownRemaining,
		}
		spells = append(spells, m)
	}

	// Convert items to maps
	items := make([]map[string]interface{}, 0, len(resp.Entity.Items))
	for _, it := range resp.Entity.Items {
		m := map[string]interface{}{
			"inventory_id": it.InventoryId,
			"item_id":      it.ItemId,
			"name":         it.Name,
			"item_type":    it.ItemType,
			"rarity":       it.Rarity,
			"description":  it.Description,
			"quantity":     it.Quantity,
			"equipped":     it.Equipped,
			"slot":         it.Slot,
			"durability":   it.Durability,
			"bound":        it.Bound,
		}
		items = append(items, m)
	}

	return entity, spells, items, nil
}

func protoToEntity(e *cultivation.Entity) *types.Entity {
	if e == nil {
		return nil
	}

	entity := &types.Entity{
		ID:         types.EntityID(e.Id),
		EntityType: types.EntityType(e.EntityType),
		Name:       e.Name,
		Realm:      types.CultivationRealm(e.Realm),
		Position: types.WorldPosition{
			RegionID: e.Position.RegionId,
			X:        e.Position.X,
			Y:        e.Position.Y,
		},
		Attributes: types.Attributes{
			Qi:                  e.Attributes.Qi,
			MaxQi:               e.Attributes.MaxQi,
			SpiritualPower:      e.Attributes.SpiritualPower,
			MaxSpiritualPower:   e.Attributes.MaxSpiritualPower,
			DivineSense:         e.Attributes.DivineSense,
			Comprehension:       int(e.Attributes.Comprehension),
			Constitution:        int(e.Attributes.Constitution),
			Luck:                int(e.Attributes.Luck),
			CultivationProgress: e.Attributes.CultivationProgress,
			AttackPower:         e.Attributes.AttackPower,
			Defense:             e.Attributes.Defense,
			Speed:               e.Attributes.Speed,
			MentalStability:     int(e.Attributes.MentalStability),
			RemainingLifespan:   int(e.Attributes.RemainingLifespan),
			MaxLifespan:         int(e.Attributes.MaxLifespan),
			CritRate:            e.Attributes.CritRate,
			CritDamage:          e.Attributes.CritDamage,
			DodgeRate:           e.Attributes.DodgeRate,
			HitRate:             e.Attributes.HitRate,
			Penetration:         e.Attributes.Penetration,
			DamageReduction:     e.Attributes.DamageReduction,
			AlchemyLevel:        int(e.Attributes.AlchemyLevel),
			ArtificingLevel:     int(e.Attributes.ArtificingLevel),
			FormationLevel:      int(e.Attributes.FormationLevel),
			FireControl:         int(e.Attributes.FireControl),
			HerbKnowledge:       int(e.Attributes.HerbKnowledge),
			MiningSkill:         int(e.Attributes.MiningSkill),
			TalismanSkill:       int(e.Attributes.TalismanSkill),
			BeastTaming:         int(e.Attributes.BeastTaming),
			Reputation:          int(e.Attributes.Reputation),
			SectContribution:    int(e.Attributes.SectContribution),
			DaoHeart:            int(e.Attributes.DaoHeart),
			Enlightenment:       int(e.Attributes.Enlightenment),
			RootPurity:          int(e.Attributes.RootPurity),
			PoisonLevel:         int(e.Attributes.PoisonLevel),
			CurseLevel:          int(e.Attributes.CurseLevel),
		},
		Karma: types.Karma{
			KarmaValue:   int(e.Karma.KarmaValue),
			Merit:        int(e.Karma.Merit),
			KarmicDebt:   int(e.Karma.KarmicDebt),
			HeavenlyMark: e.Karma.HeavenlyMark,
		},
		Status:    types.EntityStatus(e.Status),
		CreatedAt: time.Unix(e.CreatedAt, 0),
		UpdatedAt: time.Unix(e.UpdatedAt, 0),
	}

	// Map SpiritStones from proto Entity field into Attributes
	if e.SpiritStones != nil {
		entity.Attributes.SpiritStones = types.SpiritStones{
			LowGrade:     e.SpiritStones.LowGrade,
			MediumGrade:  e.SpiritStones.MediumGrade,
			HighGrade:    e.SpiritStones.HighGrade,
			PremiumGrade: e.SpiritStones.PremiumGrade,
		}
	}

	return entity
}

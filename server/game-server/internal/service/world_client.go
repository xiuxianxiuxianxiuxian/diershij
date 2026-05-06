package service

import (
	"context"
	"time"

	pb "github.com/cultivation-world/shared/proto/pb"
	"github.com/cultivation-world/shared/types"
	"google.golang.org/grpc"
)

type worldGrpcClient struct {
	client pb.WorldServiceClient
}

func NewWorldGrpcClient(cc grpc.ClientConnInterface) WorldServiceClient {
	return &worldGrpcClient{client: pb.NewWorldServiceClient(cc)}
}

func (w *worldGrpcClient) GetRegion(ctx context.Context, regionID string) (*types.Region, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := w.client.GetRegion(ctx, &pb.RegionRequest{RegionId: regionID})
	if err != nil {
		return nil, err
	}
	if !resp.Found || resp.Region == nil {
		return nil, nil
	}
	return protoRegionToTypes(resp.Region), nil
}

func (w *worldGrpcClient) SpawnResources(ctx context.Context, regionID string) error {
	_, err := w.client.SpawnResources(ctx, &pb.SpawnRequest{RegionId: regionID})
	return err
}

func (w *worldGrpcClient) GetEventModifiers(ctx context.Context, regionID string) (*EventModifiers, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	resp, err := w.client.GetWorldState(ctx, &pb.WorldStateRequest{})
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.State == nil {
		return &EventModifiers{CultivationMultiplier: 1.0}, nil
	}

	mods := &EventModifiers{CultivationMultiplier: 1.0}
	now := time.Now()

	for _, pe := range resp.State.ActiveEvents {
		if pe.RegionId != regionID || pe.Status != "active" {
			continue
		}
		endTime := time.Unix(pe.EndTime, 0)
		if now.After(endTime) {
			continue
		}

		mods.ActiveEvents = append(mods.ActiveEvents, EventInfo{
			Type:        pe.Type,
			Name:        pe.Name,
			Description: pe.Description,
			RegionID:    pe.RegionId,
			EndTime:     pe.EndTime,
		})

		switch pe.Type {
		case "beast_tide":
			mods.CultivationMultiplier = 0.7
			mods.CombatProbability = 0.3
		case "secret_realm":
			mods.CultivationMultiplier = 1.3
			mods.ExplorationBonus = 0.5
		case "treasure_appear":
			mods.ExplorationBonus = 0.3
			mods.ResourceRespawnBonus = 0.5
		case "tribulation_omen":
			mods.CultivationMultiplier = 1.5
			mods.TribulationModifier = 0.2
		}
	}

	return mods, nil
}

func protoRegionToTypes(r *pb.Region) *types.Region {
	var parentID *types.RegionID
	if r.ParentRegionId != "" {
		id := types.RegionID(r.ParentRegionId)
		parentID = &id
	}

	resources := make([]types.Resource, len(r.Resources))
	for i, res := range r.Resources {
		lastHarvested := time.Unix(0, res.LastHarvested)
		resources[i] = types.Resource{
			ID:           res.Id,
			Name:         res.Name,
			Type:         res.Type,
			Rarity:       int(res.Rarity),
			Quantity:     int(res.Quantity),
			RespawnRate:  res.RespawnRate,
			LastHarvested: &lastHarvested,
		}
	}

	var rules types.RegionRules
	if r.Rules != nil {
		rules = types.RegionRules{
			IsRestricted:     r.Rules.IsRestricted,
			RestrictedBy:     r.Rules.RestrictedBy,
			TaxRate:         r.Rules.TaxRate,
			ForbiddenActions: r.Rules.ForbiddenActions,
		}
	}

	return &types.Region{
		ID:               types.RegionID(r.Id),
		Name:             r.Name,
		ParentRegionID:   parentID,
		SpiritualDensity: r.SpiritualDensity,
		SpiritualTier:    int(r.SpiritualTier),
		DangerLevel:      int(r.DangerLevel),
		Resources:        resources,
		Rules:            rules,
		Description:      r.Description,
		Lore:             r.Lore,
	}
}

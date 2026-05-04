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

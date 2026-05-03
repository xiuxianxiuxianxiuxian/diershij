package service

import (
	"context"
	"fmt"

	pb "github.com/cultivation-world/shared/proto/pb"
	"google.golang.org/grpc"
)

// HeavenlyDaoClient defines the interface for heavenly-dao service calls.
type HeavenlyDaoClient interface {
	CheckTribulation(ctx context.Context, entityID string, targetRealm string) (*TribulationResult, error)
}

// TribulationResult represents the outcome of a tribulation check.
type TribulationResult struct {
	Success  bool
	Reason   string
	Severity int
}

type heavenlyDaoGrpcClient struct {
	client pb.HeavenlyDaoServiceClient
}

func NewHeavenlyDaoGrpcClient(cc grpc.ClientConnInterface) HeavenlyDaoClient {
	return &heavenlyDaoGrpcClient{client: pb.NewHeavenlyDaoServiceClient(cc)}
}

func (h *heavenlyDaoGrpcClient) CheckTribulation(ctx context.Context, entityID string, targetRealm string) (*TribulationResult, error) {
	resp, err := h.client.CheckTribulation(ctx, &pb.TribulationRequest{
		EntityId:    entityID,
		TargetRealm: targetRealm,
	})
	if err != nil {
		return nil, fmt.Errorf("tribulation check failed: %w", err)
	}
	if resp == nil {
		return &TribulationResult{Success: true}, nil
	}
	reason := "天劫未触发"
	if resp.WillTrigger {
		reason = fmt.Sprintf("天劫降临！强度: %.1f, 类型: %s", resp.Strength, resp.TribulationType)
	}
	return &TribulationResult{
		Success:  !resp.WillTrigger,
		Reason:   reason,
		Severity: int(resp.Strength),
	}, nil
}

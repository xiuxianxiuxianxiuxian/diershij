package service

import (
	"context"
	"fmt"
	"time"

	pb "github.com/cultivation-world/shared/proto/pb"
	"github.com/cultivation-world/shared/types"
	"google.golang.org/grpc"
)

// HeavenlyDaoClient defines the interface for heavenly-dao service calls.
type HeavenlyDaoClient interface {
	CheckTribulation(ctx context.Context, entityID string, targetRealm string) (*TribulationResult, error)
	ExecuteBreakthrough(ctx context.Context, input *BreakthroughInput) (*BreakthroughOutput, error)
	ExecuteCultivate(ctx context.Context, input *CultivateEfficiencyInput) (*CultivateEfficiencyOutput, error)
}

// TribulationResult represents the outcome of a tribulation check.
type TribulationResult struct {
	Success  bool
	Reason   string
	Severity int
}

// BreakthroughInput holds the data needed for heavenly-dao breakthrough calculation.
type BreakthroughInput struct {
	EntityID           string
	CurrentRealm       types.CultivationRealm
	TargetRealm        types.CultivationRealm
	CultivationTime    float64
	MethodQuality      float64
	ResourceBonus      float64
	MentalStability    int
	Luck               int
	Karma              int
	Merit              int
}

// BreakthroughOutput holds the result from heavenly-dao breakthrough.
type BreakthroughOutput struct {
	Success              bool
	NewRealm             types.CultivationRealm
	SuccessRate          float64
	TribulationTriggered bool
	TribulationStrength  float64
	TribulationType      string
	PenaltyProgressLoss  float64
	PenaltyCooldownHours float64
	PenaltyMentalDamage  int
	Message              string
}

// CultivateEfficiencyInput holds the data for cultivation efficiency calculation.
type CultivateEfficiencyInput struct {
	EntityID         string
	Realm            types.CultivationRealm
	SpiritualRoots   []types.SpiritualRoot
	SpiritualDensity float64
	Comprehension    int
	MentalStability  int
	BaseLifespan     int
	CurrentAge       int
}

// CultivateEfficiencyOutput holds the result from heavenly-dao cultivation.
type CultivateEfficiencyOutput struct {
	Success           bool
	CultivationGained float64
	Message           string
}

type heavenlyDaoGrpcClient struct {
	client pb.HeavenlyDaoServiceClient
}

func NewHeavenlyDaoGrpcClient(cc grpc.ClientConnInterface) HeavenlyDaoClient {
	return &heavenlyDaoGrpcClient{client: pb.NewHeavenlyDaoServiceClient(cc)}
}

func (h *heavenlyDaoGrpcClient) CheckTribulation(ctx context.Context, entityID string, targetRealm string) (*TribulationResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
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

func (h *heavenlyDaoGrpcClient) ExecuteBreakthrough(ctx context.Context, input *BreakthroughInput) (*BreakthroughOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	resp, err := h.client.ExecuteBreakthrough(ctx, &pb.BreakthroughRequest{
		EntityId:           input.EntityID,
		CurrentRealm:       string(input.CurrentRealm),
		TargetRealm:        string(input.TargetRealm),
		CultivationProgress: input.CultivationTime,
		Comprehension:      0,
		MethodQuality:      input.MethodQuality,
		ResourceBonus:      input.ResourceBonus,
		MentalStability:    int32(input.MentalStability),
		Luck:               int32(input.Luck),
		Karma:              int32(input.Karma),
		Merit:              int32(input.Merit),
	})
	if err != nil {
		return nil, fmt.Errorf("breakthrough failed: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("breakthrough: empty response")
	}
	return &BreakthroughOutput{
		Success:              resp.Success,
		NewRealm:             types.CultivationRealm(resp.NewRealm),
		SuccessRate:          resp.SuccessRate,
		TribulationTriggered: resp.TribulationTriggered,
		TribulationStrength:  resp.TribulationStrength,
		TribulationType:      resp.TribulationType,
		PenaltyProgressLoss:  resp.PenaltyProgressLoss,
		PenaltyCooldownHours: resp.PenaltyCooldownHours,
		PenaltyMentalDamage:  int(resp.PenaltyMentalDamage),
		Message:              resp.Message,
	}, nil
}

func (h *heavenlyDaoGrpcClient) ExecuteCultivate(ctx context.Context, input *CultivateEfficiencyInput) (*CultivateEfficiencyOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	roots := make([]*pb.SpiritualRootInput, 0, len(input.SpiritualRoots))
	for _, r := range input.SpiritualRoots {
		roots = append(roots, &pb.SpiritualRootInput{
			Element: r.Element,
			Purity:  int32(r.Purity),
		})
	}

	resp, err := h.client.ExecuteCultivate(ctx, &pb.CultivateRequest{
		EntityId:        input.EntityID,
		Realm:           string(input.Realm),
		SpiritualRoots:  roots,
		SpiritualDensity: input.SpiritualDensity,
		Comprehension:   int32(input.Comprehension),
		MentalStability: int32(input.MentalStability),
		BaseLifespan:    int32(input.BaseLifespan),
		CurrentAge:      int32(input.CurrentAge),
	})
	if err != nil {
		return nil, fmt.Errorf("cultivate failed: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("cultivate: empty response")
	}
	return &CultivateEfficiencyOutput{
		Success:           resp.Success,
		CultivationGained: resp.CultivationGained,
		Message:           resp.Message,
	}, nil
}

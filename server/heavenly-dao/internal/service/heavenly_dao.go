package service

import (
	"context"
	"strconv"

	cultivation "github.com/cultivation-world/shared/proto/pb"
	"github.com/cultivation-world/shared/types"
)

// HeavenlyDaoService 实现天道服务
type HeavenlyDaoService struct {
	cultivation.UnimplementedHeavenlyDaoServiceServer
	karmaRule        *KarmaRule
	tribulationRule  *TribulationRule
	worldBalanceRule *WorldBalanceRule
}

// NewHeavenlyDaoService 创建新的天道服务
func NewHeavenlyDaoService() *HeavenlyDaoService {
	return &HeavenlyDaoService{
		karmaRule:        NewKarmaRule(nil),
		tribulationRule:  NewTribulationRule(),
		worldBalanceRule: NewWorldBalanceRule(),
	}
}

// EvaluateKarma 评估业力
func (s *HeavenlyDaoService) EvaluateKarma(ctx context.Context, req *cultivation.KarmaRequest) (*cultivation.KarmaResponse, error) {
	// 构建 karma context
	karmaCtx := &KarmaContext{
		ActionType: req.ActionType,
	}

	result := s.karmaRule.CalculateKarmaChange(karmaCtx)

	return &cultivation.KarmaResponse{
		KarmaChange:  int32(result.KarmaChange),
		MeritChange:  0, // 简化处理
		HeavenlyMark: result.NewHeavenlyMark,
		Reason:       result.Reason,
	}, nil
}

// CheckTribulation 检查天劫
func (s *HeavenlyDaoService) CheckTribulation(ctx context.Context, req *cultivation.TribulationRequest) (*cultivation.TribulationResponse, error) {
	// 解析 realm
	realmInt, _ := strconv.Atoi(req.TargetRealm)
	realm := types.CultivationRealm(realmInt)

	input := TribulationInput{
		TargetRealm: realm,
		Karma:       100, // 默认值
		Merit:       100,
		Luck:        50,
	}

	result := s.tribulationRule.Assess(input)

	// 确定天劫类型
	tribulationType := "minor"
	if result.Strength > 200 {
		tribulationType = "major"
	}
	if result.Strength > 500 {
		tribulationType = "heavenly"
	}

	return &cultivation.TribulationResponse{
		WillTrigger:     result.Triggered,
		Probability:     result.Probability,
		Strength:        result.Strength,
		TribulationType: tribulationType,
	}, nil
}

// BalanceCheck 检查世界平衡
func (s *HeavenlyDaoService) BalanceCheck(ctx context.Context, req *cultivation.BalanceCheckRequest) (*cultivation.BalanceCheckResponse, error) {
	// 构建世界指标 - BalanceCheckRequest 没有字段，使用默认值
	metrics := WorldMetrics{
		ActiveSectCount: 5, // 默认值
	}

	health := s.worldBalanceRule.EvaluateWorldHealth(metrics)

	// 构建调整建议
	adjustments := s.worldBalanceRule.ApplyBalanceAdjustment(health, metrics)
	var adjustmentStrs []string
	for _, adj := range adjustments {
		adjustmentStrs = append(adjustmentStrs, adj.Reason)
	}

	return &cultivation.BalanceCheckResponse{
		NeedsAdjustment: health.Status != "healthy",
		Adjustments:     adjustmentStrs,
	}, nil
}

// ApplyKarmaDecay 应用业力衰减
func (s *HeavenlyDaoService) ApplyKarmaDecay(ctx context.Context, req *cultivation.DecayRequest) (*cultivation.DecayResponse, error) {
	// 使用 karma rule 的衰减计算 - 需要假设一个当前 karma 值
	currentKarma := 1000 // 假设默认值
	newKarma := s.karmaRule.ApplyKarmaDecay(currentKarma, float64(req.TimeElapsed))

	return &cultivation.DecayResponse{
		OldKarma:    int32(currentKarma),
		NewKarma:    int32(newKarma),
		DecayAmount: int32(currentKarma - newKarma),
	}, nil
}

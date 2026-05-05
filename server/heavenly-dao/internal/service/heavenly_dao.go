package service

import (
	"context"
	"time"

	cultivation "github.com/cultivation-world/shared/proto/pb"
	"github.com/cultivation-world/shared/types"
)

// HeavenlyDaoService 实现天道服务
type HeavenlyDaoService struct {
	cultivation.UnimplementedHeavenlyDaoServiceServer
	karmaRule        *KarmaRule
	tribulationRule  *TribulationRule
	worldBalanceRule *WorldBalanceRule
	operationSvc     *OperationService
}

// NewHeavenlyDaoService 创建新的天道服务
func NewHeavenlyDaoService() *HeavenlyDaoService {
	return &HeavenlyDaoService{
		karmaRule:        NewKarmaRule(nil),
		tribulationRule:  NewTribulationRule(),
		worldBalanceRule: NewWorldBalanceRule(),
		operationSvc:     NewOperationService(0),
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

// CheckTribulation 检查天劫（使用真实属性值）
func (s *HeavenlyDaoService) CheckTribulation(ctx context.Context, req *cultivation.TribulationRequest) (*cultivation.TribulationResponse, error) {
	// 解析 realm
	realm := types.CultivationRealm(req.TargetRealm)

	// 使用请求中传入的真实属性值（不再使用硬编码默认值）
	karma := int(req.GetKarma())
	merit := int(req.GetMerit())
	luck := int(req.GetLuck())
	if karma == 0 {
		karma = 100
	}
	if merit == 0 {
		merit = 100
	}
	if luck == 0 {
		luck = 50
	}

	input := TribulationInput{
		TargetRealm: realm,
		Karma:       karma,
		Merit:       merit,
		Luck:        luck,
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

// ExecuteBreakthrough 执行完整突破计算（调用 heavenly-dao 的突破规则）
func (s *HeavenlyDaoService) ExecuteBreakthrough(ctx context.Context, req *cultivation.BreakthroughRequest) (*cultivation.BreakthroughResponse, error) {
	currentRealm := types.CultivationRealm(req.GetCurrentRealm())
	targetRealm := types.CultivationRealm(req.GetTargetRealm())

	// 使用 heavenly-dao 的突破规则
	input := OpBreakthroughInput{
		EntityID:        req.GetEntityId(),
		CurrentRealm:    currentRealm,
		TargetRealm:     targetRealm,
		CultivationTime: req.GetCultivationProgress(),
		RequiredTime:    100.0,
		MethodQuality:   req.GetMethodQuality(),
		ResourceBonus:   req.GetResourceBonus(),
		MentalStability: int(req.GetMentalStability()),
		Luck:            int(req.GetLuck()),
	}

	result, err := s.operationSvc.ExecuteBreakthrough(input, time.Now(), DefaultRand())
	if err != nil {
		return nil, err
	}

	resp := &cultivation.BreakthroughResponse{
		Success:     result.Success,
		NewRealm:    string(result.NewRealm),
		SuccessRate: result.SuccessRate,
		Message:     result.Message,
	}

	// 填充天劫信息
	if result.Tribulation != nil {
		resp.TribulationTriggered = result.Tribulation.Triggered
		resp.TribulationStrength = result.Tribulation.Strength
		tribType := "minor"
		if result.Tribulation.Strength > 200 {
			tribType = "major"
		}
		if result.Tribulation.Strength > 500 {
			tribType = "heavenly"
		}
		resp.TribulationType = tribType
	}

	// 填充失败惩罚
	if result.Penalty != nil {
		resp.PenaltyProgressLoss = result.Penalty.ProgressLoss
		resp.PenaltyCooldownHours = result.Penalty.CooldownHours
		resp.PenaltyMentalDamage = int32(result.Penalty.MentalDamage)
	}

	return resp, nil
}

// ExecuteCultivate 执行修炼效率计算
func (s *HeavenlyDaoService) ExecuteCultivate(ctx context.Context, req *cultivation.CultivateRequest) (*cultivation.CultivateResponse, error) {
	// 转换灵根数据
	roots := make([]types.SpiritualRoot, 0, len(req.GetSpiritualRoots()))
	for _, r := range req.GetSpiritualRoots() {
		roots = append(roots, types.SpiritualRoot{
			Element: r.GetElement(),
			Purity:  int(r.GetPurity()),
		})
	}

	input := CultivateInput{
		EntityID:         req.GetEntityId(),
		Realm:            types.CultivationRealm(req.GetRealm()),
		SpiritualRoots:   roots,
		SpiritualDensity: req.GetSpiritualDensity(),
		Comprehension:    int(req.GetComprehension()),
		MentalStability:  int(req.GetMentalStability()),
		BaseLifespan:     int(req.GetBaseLifespan()),
		CurrentAge:       int(req.GetCurrentAge()),
	}

	result, err := s.operationSvc.ExecuteCultivate(input, time.Now(), DefaultRand())
	if err != nil {
		return nil, err
	}

	return &cultivation.CultivateResponse{
		Success:           result.Success,
		CultivationGained: result.CultivationGained,
		Rate:              result.Rate,
		RealmBonus:        result.RealmBonus,
		SpiritualMult:     result.SpiritualMult,
		MethodMatch:       result.MethodMatch,
		MentalFactor:      result.MentalFactor,
		AgingPenalty:      result.AgingPenalty,
		Message:           result.Message,
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

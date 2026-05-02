package service

import (
    "context"
    "math"

    "github.com/cultivation-world/shared/types"
)

type HeavenlyDaoService struct {
    karmaThresholds  map[string]int
    realmLifespan    map[string]int
    tribulationBase  map[string]float64
    karmaDecayRate   float64
}

func NewHeavenlyDaoService() *HeavenlyDaoService {
    return &HeavenlyDaoService{
        karmaDecayRate: 0.01,
        karmaThresholds: map[string]int{
            "clear":      100,
            "slight":     500,
            "heavy":      1000,
            "notorious":  5000,
            "heaven_fury": 10000,
        },
        realmLifespan: map[string]int{
            "mortal":           80,
            "qi_condensation":  120,
            "foundation":       200,
            "golden_core":      500,
            "nascent_soul":     1000,
            "soul_transformation": 3000,
            "void_refinement":  5000,
            "integration":      8000,
            "mahayana":         10000,
            "tribulation":      15000,
        },
        tribulationBase: map[string]float64{
            "qi_condensation":  0.1,
            "foundation":       0.2,
            "golden_core":      0.3,
            "nascent_soul":     0.5,
            "soul_transformation": 0.6,
            "void_refinement":  0.7,
            "integration":      0.8,
            "mahayana":         0.9,
            "tribulation":      0.95,
        },
    }
}

func (s *HeavenlyDaoService) EvaluateKarma(ctx context.Context, req *game.KarmaRequest) (*game.KarmaResponse, error) {
    baseKarma := s.getActionKarmaBase(req.ActionType)

    contextMultiplier := 1.0

    karmaChange := int(float64(baseKarma) * contextMultiplier)
    meritChange := 0

    if karmaChange < 0 {
        meritChange = -karmaChange / 2
        karmaChange = 0
    }

    heavenlyMark := s.calculateHeavenlyMark(karmaChange)

    return &game.KarmaResponse{
        KarmaChange:   int32(karmaChange),
        MeritChange:   int32(meritChange),
        HeavenlyMark:  heavenlyMark,
        Reason:        s.getKarmaReason(req.ActionType),
    }, nil
}

func (s *HeavenlyDaoService) CheckTribulation(ctx context.Context, req *game.TribulationRequest) (*game.TribulationResponse, error) {
    baseProb := s.tribulationBase[req.TargetRealm]

    karmaFactor := 1.0
    meritFactor := 1.0
    luckFactor := 1.0

    prob := baseProb * karmaFactor * meritFactor * luckFactor
    prob = math.Max(0.1, math.Min(0.95, prob))

    strength := s.calculateTribulationStrength(req.TargetRealm, 0)

    return &game.TribulationResponse{
        WillTrigger:     prob > 0.3,
        Probability:     prob,
        Strength:        strength,
        TribulationType: s.getTribulationType(req.TargetRealm),
    }, nil
}

func (s *HeavenlyDaoService) BalanceCheck(ctx context.Context, req *game.BalanceCheckRequest) (*game.BalanceCheckResponse, error) {
    metrics := &game.BalanceMetrics{
        PowerDistribution:   0.5,
        ResourceCirculation: 0.3,
        SectDiversity:       0.4,
        KarmaDistribution:   0.6,
    }

    adjustments := []string{}
    needsAdjustment := false

    if metrics.PowerDistribution > 0.7 {
        adjustments = append(adjustments, "spawn_opportunity_for_weak")
        needsAdjustment = true
    }

    if metrics.ResourceCirculation < 0.1 {
        adjustments = append(adjustments, "increase_resource_spawn")
        needsAdjustment = true
    }

    if metrics.KarmaDistribution > 500 {
        adjustments = append(adjustments, "trigger_heavenly_cleansing")
        needsAdjustment = true
    }

    return &game.BalanceCheckResponse{
        NeedsAdjustment: needsAdjustment,
        Adjustments:     adjustments,
        Metrics:         metrics,
    }, nil
}

func (s *HeavenlyDaoService) ApplyKarmaDecay(ctx context.Context, req *game.DecayRequest) (*game.DecayResponse, error) {
    decayAmount := int(float64(req.OldKarma) * s.karmaDecayRate)
    newKarma := req.OldKarma - decayAmount
    if newKarma < 0 {
        newKarma = 0
    }

    return &game.DecayResponse{
        OldKarma:     req.OldKarma,
        NewKarma:     int32(newKarma),
        DecayAmount:  int32(decayAmount),
    }, nil
}

func (s *HeavenlyDaoService) getActionKarmaBase(actionType string) int {
    karmaMap := map[string]int{
        "kill_innocent":     500,
        "kill_cultivator":   200,
        "kill_demon":        -50,
        "save_life":         -100,
        "teach_method":      -200,
        "betray_master":     1000,
        "break_oath":        300,
        "destroy_sect":      800,
        "create_method":     -500,
    }

    if karma, ok := karmaMap[actionType]; ok {
        return karma
    }
    return 0
}

func (s *HeavenlyDaoService) calculateHeavenlyMark(karma int) string {
    if karma < 100 {
        return "clear"
    } else if karma < 500 {
        return "slight"
    } else if karma < 1000 {
        return "heavy"
    } else if karma < 5000 {
        return "notorious"
    }
    return "heaven_fury"
}

func (s *HeavenlyDaoService) calculateTribulationStrength(realm string, karma int) float64 {
    baseStrength := 100.0

    realmMultiplier := map[string]float64{
        "qi_condensation":  1.0,
        "foundation":       2.0,
        "golden_core":      5.0,
        "nascent_soul":     10.0,
        "soul_transformation": 20.0,
        "void_refinement":  50.0,
        "integration":      100.0,
        "mahayana":         200.0,
        "tribulation":      500.0,
    }

    multiplier := realmMultiplier[realm]
    karmaFactor := 1.0 + float64(karma)/500.0

    return baseStrength * multiplier * karmaFactor
}

func (s *HeavenlyDaoService) getTribulationType(realm string) string {
    types := map[string]string{
        "qi_condensation":  "thunder",
        "foundation":       "thunder",
        "golden_core":      "thunder_fire",
        "nascent_soul":     "thunder_fire_wind",
        "soul_transformation": "five_element",
        "void_refinement":  "heart_demon",
        "integration":      "dao_tribulation",
        "mahayana":         "extinction",
        "tribulation":      "ascension",
    }

    if t, ok := types[realm]; ok {
        return t
    }
    return "thunder"
}

func (s *HeavenlyDaoService) getKarmaReason(actionType string) string {
    reasons := map[string]string{
        "kill_innocent":   "杀害无辜，业力缠身",
        "kill_cultivator": "同修相残，有损天道",
        "kill_demon":      "斩妖除魔，功德无量",
        "save_life":       "救人性命，功德加身",
        "teach_method":    "传道授业，功德无量",
        "betray_master":   "欺师灭祖，业力滔天",
        "break_oath":      "背信弃义，业力加身",
        "destroy_sect":    "毁人道统，业力深重",
        "create_method":   "开宗立派，功德无量",
    }

    if reason, ok := reasons[actionType]; ok {
        return reason
    }
    return "因果循环，天道昭昭"
}

package service

import (
    "context"
    "fmt"
    "math/rand"
    "time"

    "github.com/cultivation-world/shared/config"
)

type AISchedulerService struct {
    cfg          *config.Config
    npcRegistry  map[string]*NPCProfile
    templates    *BehaviorTemplateLibrary
    rateLimiter  *RateLimiter
}

type NPCProfile struct {
    NPCID           string
    PersonalityType string
    MoralAlignment  string
    AmbitionLevel   int
    RiskTolerance   float64
    BackgroundStory string
    CurrentGoal     string
}

type BehaviorTemplateLibrary struct {
    templates map[string][]BehaviorTemplate
}

type BehaviorTemplate struct {
    Name      string
    Condition string
    Action    string
    Weight    float64
}

type RateLimiter struct {
    tokens     int
    maxTokens  int
    refillRate int
    lastRefill time.Time
}

func NewAISchedulerService(cfg *config.Config) *AISchedulerService {
    return &AISchedulerService{
        cfg:         cfg,
        npcRegistry: make(map[string]*NPCProfile),
        templates:   NewBehaviorTemplateLibrary(),
        rateLimiter: NewRateLimiter(cfg.LLM.RateLimit),
    }
}

func NewBehaviorTemplateLibrary() *BehaviorTemplateLibrary {
    return &BehaviorTemplateLibrary{
        templates: map[string][]BehaviorTemplate{
            "cultivate": {
                {Name: "daily_cultivation", Condition: "qi<50%", Action: "cultivate", Weight: 0.8},
                {Name: "breakthrough_attempt", Condition: "progress>=100%", Action: "breakthrough", Weight: 0.6},
            },
            "explore": {
                {Name: "resource_gathering", Condition: "low_resources", Action: "gather", Weight: 0.7},
                {Name: "region_exploration", Condition: "curious", Action: "explore", Weight: 0.5},
            },
            "social": {
                {Name: "seek_alliance", Condition: "weak", Action: "form_alliance", Weight: 0.4},
                {Name: "trade", Condition: "surplus_resources", Action: "trade", Weight: 0.6},
            },
        },
    }
}

func NewRateLimiter(maxTokens int) *RateLimiter {
    return &RateLimiter{
        tokens:     maxTokens,
        maxTokens:  maxTokens,
        refillRate: 10,
        lastRefill: time.Now(),
    }
}

func (s *AISchedulerService) ScheduleDecision(ctx context.Context, req *game.DecisionRequest) (*game.DecisionResponse, error) {
    profile, exists := s.npcRegistry[req.NpcId]
    if !exists {
        profile = &NPCProfile{
            NPCID:          req.NpcId,
           	PersonalityType: "balanced",
            MoralAlignment:  "neutral",
            AmbitionLevel:   50,
            RiskTolerance:   0.5,
        }
    }

    template := s.matchTemplate(req.Context, req.AvailableActions)
    if template != nil && rand.Float64() < 0.7 {
        return &game.DecisionResponse{
            Action:    template.Action,
            Params:    make(map[string]string),
            Reasoning: fmt.Sprintf("Template matched: %s", template.Name),
            Source:    "behavior_tree",
            TokenCost: 0,
        }, nil
    }

    if !s.rateLimiter.Allow() {
        return s.defaultDecision(req.AvailableActions), nil
    }

    decision := s.callLLM(ctx, req, profile)
    return decision, nil
}

func (s *AISchedulerService) ExecuteBehaviorTree(ctx context.Context, req *game.BehaviorTreeRequest) (*game.BehaviorTreeResponse, error) {
    action := s.executeBehaviorTreeLogic(req.TreeName, req.Context)

    return &game.BehaviorTreeResponse{
        Action:  action.Action,
        Params:  action.Params,
        Success: true,
    }, nil
}

func (s *AISchedulerService) RegisterNPC(ctx context.Context, req *game.NPCRegisterRequest) (*game.NPCRegisterResponse, error) {
    profile := &NPCProfile{
        NPCID:           req.NpcId,
        PersonalityType: req.PersonalityType,
        MoralAlignment:  req.MoralAlignment,
        AmbitionLevel:   int(req.AmbitionLevel),
        RiskTolerance:   req.RiskTolerance,
        BackgroundStory: req.BackgroundStory,
        CurrentGoal:     req.CurrentGoal,
    }

    s.npcRegistry[req.NpcId] = profile

    return &game.NPCRegisterResponse{
        Success: true,
        Message: "NPC registered successfully",
    }, nil
}

func (s *AISchedulerService) UnregisterNPC(ctx context.Context, req *game.NPCUnregisterRequest) (*game.NPCUnregisterResponse, error) {
    delete(s.npcRegistry, req.NpcId)

    return &game.NPCUnregisterResponse{
        Success: true,
    }, nil
}

func (s *AISchedulerService) matchTemplate(context string, availableActions []string) *BehaviorTemplate {
    for _, actionType := range []string{"cultivate", "explore", "social"} {
        for _, template := range s.templates.templates[actionType] {
            for _, action := range availableActions {
                if template.Action == action {
                    return &template
                }
            }
        }
    }
    return nil
}

func (s *AISchedulerService) callLLM(ctx context.Context, req *game.DecisionRequest, profile *NPCProfile) *game.DecisionResponse {
    action := req.AvailableActions[rand.Intn(len(req.AvailableActions))]

    return &game.DecisionResponse{
        Action:    action,
        Params:    make(map[string]string),
        Reasoning: fmt.Sprintf("LLM decision based on %s personality", profile.PersonalityType),
        Source:    "llm",
        TokenCost: 100,
    }
}

func (s *AISchedulerService) defaultDecision(availableActions []string) *game.DecisionResponse {
    if len(availableActions) == 0 {
        return &game.DecisionResponse{
            Action:    "meditate",
            Params:    make(map[string]string),
            Reasoning: "Default action: meditate",
            Source:    "fallback",
        }
    }

    return &game.DecisionResponse{
        Action:    availableActions[0],
        Params:    make(map[string]string),
        Reasoning: "Rate limited, using default",
        Source:    "fallback",
    }
}

func (s *AISchedulerService) executeBehaviorTreeLogic(treeName string, context map[string]string) *BehaviorTreeResult {
    switch treeName {
    case "daily_routine":
        return &BehaviorTreeResult{
            Action: "cultivate",
            Params: map[string]string{"duration": "1h"},
        }
    case "combat":
        return &BehaviorTreeResult{
            Action: "attack",
            Params: map[string]string{"target": context["enemy_id"]},
        }
    case "exploration":
        return &BehaviorTreeResult{
            Action: "explore",
            Params: map[string]string{"direction": "random"},
        }
    default:
        return &BehaviorTreeResult{
            Action: "meditate",
            Params: make(map[string]string),
        }
    }
}

type BehaviorTreeResult struct {
    Action string
    Params map[string]string
}

func (rl *RateLimiter) Allow() bool {
    now := time.Now()
    elapsed := now.Sub(rl.lastRefill).Seconds()
    rl.tokens += int(elapsed * float64(rl.refillRate))
    if rl.tokens > rl.maxTokens {
        rl.tokens = rl.maxTokens
    }
    rl.lastRefill = now

    if rl.tokens > 0 {
        rl.tokens--
        return true
    }
    return false
}

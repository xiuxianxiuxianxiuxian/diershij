package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLLMClient(t *testing.T) {
	client := NewLLMClient("test-key", "https://api.deepseek.com")
	assert.NotNil(t, client)
	assert.Equal(t, "https://api.deepseek.com", client.BaseURL)
	assert.Equal(t, "test-key", client.APIKey)
	assert.Equal(t, "deepseek-chat", client.DefaultModel)
	assert.NotNil(t, client.RateLimiter)
	assert.NotNil(t, client.CircuitBreaker)
}

func TestNewLLMClientWithConfig(t *testing.T) {
	cfg := LLMClientConfig{
		APIKey:       "test-key",
		BaseURL:      "https://api.deepseek.com",
		Model:        "deepseek-reasoner",
		RateLimit:    100,
		RatePeriod:   time.Minute,
		FailureThreshold: 3,
		RecoveryTimeout: 60 * time.Second,
	}

	client := NewLLMClientWithConfig(cfg)
	assert.NotNil(t, client)
	assert.Equal(t, "deepseek-reasoner", client.DefaultModel)
}

func TestLLMClient_Call(t *testing.T) {
	client := NewLLMClient("test-key", "https://api.deepseek.com")

	ctx := context.Background()
	req := LLMRequest{
		Model: "deepseek-chat",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant"},
			{Role: "user", Content: "Hello"},
		},
	}

	resp, err := client.Call(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Contains(t, resp.Content, "Hello")
	assert.Equal(t, 150, resp.TokenUsage.TotalTokens)
}

func TestLLMClient_CallWithContextCancellation(t *testing.T) {
	client := NewLLMClient("test-key", "https://api.deepseek.com")

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req := LLMRequest{
		Messages: []Message{
			{Role: "user", Content: "test"},
		},
	}

	_, err := client.Call(ctx, req)
	assert.Error(t, err)
}

func TestLLMClient_MakeDecision(t *testing.T) {
	client := NewLLMClient("test-key", "https://api.deepseek.com")

	ctx := context.Background()
	systemPrompt := "You are an NPC making decisions"
	userPrompt := `{"action": "cultivate", "reasoning": "time to cultivate", "confidence": 0.9}`

	result, err := client.MakeDecision(ctx, systemPrompt, userPrompt)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "cultivate", result.Action)
	assert.Equal(t, "time to cultivate", result.Reasoning)
	assert.InDelta(t, 0.9, result.Confidence, 0.01)
}

func TestParseDecisionResult_ValidJSON(t *testing.T) {
	content := `{"action": "explore", "target": "forest", "reasoning": "need resources", "confidence": 0.8}`

	result, err := ParseDecisionResult(content)
	assert.NoError(t, err)
	assert.Equal(t, "explore", result.Action)
	assert.Equal(t, "forest", result.Target)
	assert.Equal(t, "need resources", result.Reasoning)
	assert.InDelta(t, 0.8, result.Confidence, 0.01)
}

func TestParseDecisionResult_InvalidJSON(t *testing.T) {
	content := "Just explore the forest"

	result, err := ParseDecisionResult(content)
	assert.NoError(t, err)
	assert.Equal(t, "Just explore the forest", result.Action)
	assert.InDelta(t, 0.5, result.Confidence, 0.01)
}

func TestParseDecisionResult_MissingAction(t *testing.T) {
	content := `{"reasoning": "no action", "confidence": 0.5}`

	_, err := ParseDecisionResult(content)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "缺少 action")
}

func TestParseDecisionResult_ConfidenceBounds(t *testing.T) {
	// Too high
	content := `{"action": "test", "confidence": 1.5}`
	result, _ := ParseDecisionResult(content)
	assert.LessOrEqual(t, result.Confidence, 1.0)

	// Too low
	content = `{"action": "test", "confidence": -0.5}`
	result, _ = ParseDecisionResult(content)
	assert.GreaterOrEqual(t, result.Confidence, 0.0)
}

func TestTokenBucketLimiter_Allow(t *testing.T) {
	// 10 tokens per second
	limiter := NewTokenBucketLimiter(10, time.Second)

	// Should allow 10 requests
	for i := 0; i < 10; i++ {
		assert.True(t, limiter.Allow())
	}

	// 11th should fail
	assert.False(t, limiter.Allow())
}

func TestTokenBucketLimiter_Refill(t *testing.T) {
	limiter := NewTokenBucketLimiter(5, 100*time.Millisecond)

	// Use all tokens
	for i := 0; i < 5; i++ {
		limiter.Allow()
	}
	assert.False(t, limiter.Allow())

	// Wait for refill
	time.Sleep(120 * time.Millisecond)
	assert.True(t, limiter.Allow())
}

func TestCircuitBreaker_NormalOperation(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 3,
		RecoveryTimeout:  time.Second,
	})

	assert.Equal(t, StateClosed, cb.GetState())
	assert.True(t, cb.Allow())
}

func TestCircuitBreaker_TripsAfterFailures(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 3,
		RecoveryTimeout:  time.Second,
	})

	cb.RecordFailure()
	cb.RecordFailure()
	assert.Equal(t, StateClosed, cb.GetState())

	cb.RecordFailure()
	assert.Equal(t, StateOpen, cb.GetState())
	assert.False(t, cb.Allow())
}

func TestCircuitBreaker_Recovery(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 2,
		RecoveryTimeout:  100 * time.Millisecond,
	})

	// Trip the breaker
	cb.RecordFailure()
	cb.RecordFailure()
	assert.Equal(t, StateOpen, cb.GetState())
	assert.False(t, cb.Allow())

	// Wait for recovery
	time.Sleep(120 * time.Millisecond)
	assert.True(t, cb.Allow())
	assert.Equal(t, StateHalfOpen, cb.GetState())

	// Success should close the circuit
	cb.RecordSuccess()
	assert.Equal(t, StateClosed, cb.GetState())
}

func TestCircuitBreaker_SuccessResetsFailureCount(t *testing.T) {
	cb := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: 3,
		RecoveryTimeout:  time.Second,
	})

	cb.RecordFailure()
	cb.RecordFailure()
	assert.Equal(t, 2, cb.GetFailureCount())

	cb.RecordSuccess()
	assert.Equal(t, 0, cb.GetFailureCount())
}

func TestSystemPromptTemplate(t *testing.T) {
	prompt := SystemPromptTemplate("张三", "善良", "金丹期")
	assert.Contains(t, prompt, "张三")
	assert.Contains(t, prompt, "善良")
	assert.Contains(t, prompt, "金丹期")
	assert.Contains(t, prompt, "cultivate")
	assert.Contains(t, prompt, "breakthrough")
}

func TestDecisionFallback_LLMFallback(t *testing.T) {
	client := NewLLMClient("test-key", "https://api.deepseek.com")
	templates := NewTemplateLibrary()
	tree := NewBehaviorTree("fallback", &LeafNode{
		BaseNode: BaseNode{Name: "default"},
		Action: func(ctx *NPCContext) NodeStatus {
			ctx.Log("fallback")
			return StatusSuccess
		},
	})

	fallback := NewDecisionFallback(client, templates, tree)

	ctx := context.Background()
	npcCtx := NewNPCContext("test_npc")
	systemPrompt := "You are an NPC"
	userPrompt := `{"action": "cultivate", "confidence": 0.9}`

	result, source := fallback.Decide(ctx, npcCtx, systemPrompt, userPrompt)
	assert.NotNil(t, result)
	assert.Equal(t, "llm", source)
}

func TestDecisionFallback_TemplateFallback(t *testing.T) {
	client := NewLLMClient("test-key", "https://api.deepseek.com")
	// Trip circuit breaker to skip LLM
	for i := 0; i < 5; i++ {
		client.CircuitBreaker.RecordFailure()
	}

	templates := NewTemplateLibrary()
	templates.AddTemplate(&BehaviorTemplate{
		ID:          "t1",
		Type:        TemplateTypeDecision,
		Pattern:     "combat has_target",
		Action:      "attack",
		Parameters:  map[string]string{"type": "melee"},
		Priority:    10,
		SuccessRate: 0.8,
		Tags:        []string{"combat"},
	})

	tree := NewBehaviorTree("fallback", &LeafNode{
		BaseNode: BaseNode{Name: "default"},
		Action: func(ctx *NPCContext) NodeStatus {
			return StatusSuccess
		},
	})

	fallback := NewDecisionFallback(client, templates, tree)

	ctx := context.Background()
	npcCtx := NewNPCContext("test_npc")
	npcCtx.IsInCombat = true
	npcCtx.HasTarget = true

	result, source := fallback.Decide(ctx, npcCtx, "", "")
	assert.NotNil(t, result)
	assert.Equal(t, "attack", result.Action)
	assert.Equal(t, "template", source)
}

func TestDecisionFallback_BehaviorTreeFallback(t *testing.T) {
	client := NewLLMClient("test-key", "https://api.deepseek.com")
	// Trip circuit breaker
	for i := 0; i < 5; i++ {
		client.CircuitBreaker.RecordFailure()
	}

	templates := NewTemplateLibrary() // empty

	tree := NewBehaviorTree("fallback", &LeafNode{
		BaseNode: BaseNode{Name: "default"},
		Action: func(ctx *NPCContext) NodeStatus {
			return StatusSuccess
		},
	})

	fallback := NewDecisionFallback(client, templates, tree)

	ctx := context.Background()
	npcCtx := NewNPCContext("test_npc")

	result, source := fallback.Decide(ctx, npcCtx, "", "")
	assert.NotNil(t, result)
	assert.Equal(t, "behavior_tree", source)
}

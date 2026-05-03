package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// LLMProvider represents the LLM provider type.
type LLMProvider string

const (
	ProviderDeepSeekChat    LLMProvider = "deepseek-chat"
	ProviderDeepSeekReasoner LLMProvider = "deepseek-reasoner"
)

// LLMRequest represents a request to the LLM API.
type LLMRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	Timeout     time.Duration
}

// Message represents a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMResponse represents the response from the LLM API.
type LLMResponse struct {
	Content     string
	TokenUsage  TokenUsage
	Model       string
	Duration    time.Duration
}

// TokenUsage tracks token consumption.
type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// DecisionResult represents an NPC decision from LLM.
type DecisionResult struct {
	Action      string                 `json:"action"`
	Target      string                 `json:"target,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Reasoning   string                 `json:"reasoning"`
	Confidence  float64                `json:"confidence"`
}

// LLMClient handles communication with the DeepSeek API.
type LLMClient struct {
	BaseURL     string
	APIKey      string
	DefaultModel string
	RateLimiter *TokenBucketLimiter
	CircuitBreaker *CircuitBreaker
}

// NewLLMClient creates a new LLM client.
func NewLLMClient(apiKey string, baseURL string) *LLMClient {
	return &LLMClient{
		BaseURL:      baseURL,
		APIKey:       apiKey,
		DefaultModel: "deepseek-chat",
		RateLimiter:  NewTokenBucketLimiter(600, time.Minute), // 600 RPM
		CircuitBreaker: NewCircuitBreaker(CircuitBreakerConfig{
			FailureThreshold: 5,
			RecoveryTimeout:  30 * time.Second,
		}),
	}
}

// LLMClientConfig holds configuration for the LLM client.
type LLMClientConfig struct {
	APIKey       string
	BaseURL      string
	Model        string
	RateLimit    int
	RatePeriod   time.Duration
	FailureThreshold int
	RecoveryTimeout  time.Duration
}

// NewLLMClientWithConfig creates a new LLM client with custom config.
func NewLLMClientWithConfig(cfg LLMClientConfig) *LLMClient {
	client := &LLMClient{
		BaseURL:      cfg.BaseURL,
		APIKey:       cfg.APIKey,
		DefaultModel: cfg.Model,
		RateLimiter:  NewTokenBucketLimiter(cfg.RateLimit, cfg.RatePeriod),
		CircuitBreaker: NewCircuitBreaker(CircuitBreakerConfig{
			FailureThreshold: cfg.FailureThreshold,
			RecoveryTimeout:  cfg.RecoveryTimeout,
		}),
	}

	if client.DefaultModel == "" {
		client.DefaultModel = "deepseek-chat"
	}

	return client
}

// Call sends a request to the LLM API.
func (c *LLMClient) Call(ctx context.Context, req LLMRequest) (*LLMResponse, error) {
	// Check rate limiter
	if !c.RateLimiter.Allow() {
		return nil, fmt.Errorf("速率限制: 请稍后再试")
	}

	// Check circuit breaker
	if !c.CircuitBreaker.Allow() {
		return nil, fmt.Errorf("熔断器开启: API 暂时不可用")
	}

	// Record start time
	start := time.Now()

	// Simulate API call (in production, this would be an HTTP request)
	response, err := c.simulateAPICall(ctx, req)
	if err != nil {
		c.CircuitBreaker.RecordFailure()
		return nil, err
	}

	c.CircuitBreaker.RecordSuccess()
	response.Duration = time.Since(start)
	response.Model = req.Model
	if response.Model == "" {
		response.Model = c.DefaultModel
	}

	return response, nil
}

// simulateAPICall simulates an API call for testing.
func (c *LLMClient) simulateAPICall(ctx context.Context, req LLMRequest) (*LLMResponse, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// For testing, return a mock response based on the request content
	lastMessage := ""
	if len(req.Messages) > 0 {
		lastMessage = req.Messages[len(req.Messages)-1].Content
	}

	// If the message looks like JSON decision, return it directly
	if len(lastMessage) > 0 && lastMessage[0] == '{' {
		return &LLMResponse{
			Content: lastMessage,
			TokenUsage: TokenUsage{
				PromptTokens:     100,
				CompletionTokens: 50,
				TotalTokens:      150,
			},
		}, nil
	}

	return &LLMResponse{
		Content: fmt.Sprintf("模拟响应: %s", lastMessage),
		TokenUsage: TokenUsage{
			PromptTokens:     100,
			CompletionTokens: 50,
			TotalTokens:      150,
		},
	}, nil
}

// MakeDecision sends a decision request to the LLM.
func (c *LLMClient) MakeDecision(ctx context.Context, systemPrompt string, userPrompt string) (*DecisionResult, error) {
	req := LLMRequest{
		Model: c.DefaultModel,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		MaxTokens:   500,
		Temperature: 0.7,
	}

	response, err := c.Call(ctx, req)
	if err != nil {
		return nil, err
	}

	return ParseDecisionResult(response.Content)
}

// ParseDecisionResult parses the LLM response into a DecisionResult.
func ParseDecisionResult(content string) (*DecisionResult, error) {
	var result DecisionResult

	// Try to parse as JSON
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		// If not valid JSON, create a basic result
		result = DecisionResult{
			Action:    content,
			Reasoning: "非 JSON 响应",
			Confidence: 0.5,
		}
	}

	// Validate the result
	if result.Action == "" {
		return nil, fmt.Errorf("决策结果缺少 action 字段")
	}

	if result.Confidence < 0 || result.Confidence > 1 {
		result.Confidence = 0.5
	}

	return &result, nil
}

// TokenBucketLimiter implements a token bucket rate limiter.
type TokenBucketLimiter struct {
	mu         sync.Mutex
	tokens     float64
	maxTokens  float64
	refillRate float64 // tokens per period
	lastRefill time.Time
	period     time.Duration
}

// NewTokenBucketLimiter creates a new token bucket limiter.
func NewTokenBucketLimiter(maxTokens int, period time.Duration) *TokenBucketLimiter {
	now := time.Now()
	return &TokenBucketLimiter{
		tokens:     float64(maxTokens),
		maxTokens:  float64(maxTokens),
		refillRate: float64(maxTokens),
		lastRefill: now,
		period:     period,
	}
}

// Allow checks if a request is allowed.
func (l *TokenBucketLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.refill()

	if l.tokens >= 1 {
		l.tokens--
		return true
	}
	return false
}

func (l *TokenBucketLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(l.lastRefill)

	if elapsed >= l.period {
		l.tokens = l.maxTokens
		l.lastRefill = now
	} else if elapsed > 0 {
		tokensToAdd := elapsed.Seconds() / l.period.Seconds() * l.refillRate
		l.tokens = min(l.maxTokens, l.tokens+tokensToAdd)
		l.lastRefill = now
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// CircuitBreaker implements the circuit breaker pattern.
type CircuitBreaker struct {
	mu               sync.Mutex
	state            CircuitState
	failureCount     int
	successCount     int
	failureThreshold int
	recoveryTimeout  time.Duration
	lastFailureTime  time.Time
}

// CircuitState represents the state of the circuit breaker.
type CircuitState string

const (
	StateClosed   CircuitState = "closed"   // normal operation
	StateOpen     CircuitState = "open"     // tripped, requests blocked
	StateHalfOpen CircuitState = "half_open" // testing recovery
)

// CircuitBreakerConfig holds configuration for the circuit breaker.
type CircuitBreakerConfig struct {
	FailureThreshold int
	RecoveryTimeout  time.Duration
}

// NewCircuitBreaker creates a new circuit breaker.
func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	return &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: cfg.FailureThreshold,
		recoveryTimeout:  cfg.RecoveryTimeout,
	}
}

// Allow checks if a request is allowed through the circuit breaker.
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastFailureTime) > cb.recoveryTimeout {
			cb.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return true
	default:
		return false
	}
}

// RecordSuccess records a successful request.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateHalfOpen:
		cb.state = StateClosed
		cb.failureCount = 0
	case StateClosed:
		cb.failureCount = 0
	}
}

// RecordFailure records a failed request.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failureCount++
	cb.lastFailureTime = time.Now()

	if cb.failureCount >= cb.failureThreshold {
		cb.state = StateOpen
	}
}

// GetState returns the current circuit breaker state.
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// GetFailureCount returns the current failure count.
func (cb *CircuitBreaker) GetFailureCount() int {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.failureCount
}

// SystemPromptTemplate generates a system prompt for NPC decision making.
func SystemPromptTemplate(npcName string, personality string, realm string) string {
	return fmt.Sprintf(`你是 %s，一个修仙世界中的 NPC。

性格: %s
境界: %s

你正在做出决策。请返回 JSON 格式的决策结果:
{
  "action": "动作名称",
  "target": "目标 (可选)",
  "parameters": {},
  "reasoning": "决策理由",
  "confidence": 0.0-1.0
}

可用动作: cultivate, breakthrough, explore, gather, combat, trade, craft, rest, flee
`, npcName, personality, realm)
}

// DecisionFallback implements the fallback chain: LLM -> Template -> Behavior Tree.
type DecisionFallback struct {
	LLMClient    *LLMClient
	TemplateLib  *TemplateLibrary
	BehaviorTree *BehaviorTree
}

// NewDecisionFallback creates a new decision fallback chain.
func NewDecisionFallback(llm *LLMClient, templates *TemplateLibrary, tree *BehaviorTree) *DecisionFallback {
	return &DecisionFallback{
		LLMClient:    llm,
		TemplateLib:  templates,
		BehaviorTree: tree,
	}
}

// Decide attempts to get a decision through the fallback chain.
func (f *DecisionFallback) Decide(ctx context.Context, npcCtx *NPCContext, systemPrompt string, userPrompt string) (*DecisionResult, string) {
	// Try LLM first
	if f.LLMClient != nil && f.LLMClient.CircuitBreaker.Allow() {
		result, err := f.LLMClient.MakeDecision(ctx, systemPrompt, userPrompt)
		if err == nil && result.Confidence >= 0.7 {
			return result, "llm"
		}
	}

	// Try template matching
	if f.TemplateLib != nil {
		result := f.TemplateLib.MatchAndDecide(npcCtx)
		if result != nil {
			return result, "template"
		}
	}

	// Fall back to behavior tree
	if f.BehaviorTree != nil {
		status := f.BehaviorTree.Evaluate(npcCtx)
		return &DecisionResult{
			Action:     string(status),
			Reasoning:  "行为树决策",
			Confidence: 0.5,
		}, "behavior_tree"
	}

	return nil, "none"
}

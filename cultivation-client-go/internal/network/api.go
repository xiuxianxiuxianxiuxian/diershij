package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cultivation-client/internal/store"
	"cultivation-client/internal/types"
)

type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

var apiInstance *APIClient

func GetAPIClient() *APIClient {
	if apiInstance == nil {
		apiInstance = &APIClient{
			baseURL: "http://localhost:8081",
			httpClient: &http.Client{
				Timeout: 30 * time.Second,
			},
		}
	}
	return apiInstance
}

func (c *APIClient) SetBaseURL(url string) {
	c.baseURL = url
}

func (c *APIClient) request(method, path string, body interface{}, result interface{}) error {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, c.baseURL+path, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	token := store.GetAuthStore().GetToken()
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			if errMsg, ok := errResp["error"].(string); ok {
				return fmt.Errorf("API error: %s", errMsg)
			}
		}
		return fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Login 用户登录
func (c *APIClient) Login(username, password string) (*types.LoginResponse, error) {
	req := types.LoginRequest{
		Username: username,
		Password: password,
	}

	var resp types.LoginResponse
	if err := c.request("POST", "/auth/login", req, &resp); err != nil {
		return nil, err
	}

	if resp.Success {
		store.GetAuthStore().SetToken(resp.Token)
		store.GetAuthStore().SetEntity(&resp.Entity)
	}

	return &resp, nil
}

// Register 用户注册
func (c *APIClient) Register(username, password string) (*types.RegisterResponse, error) {
	req := types.RegisterRequest{
		Username: username,
		Password: password,
	}

	var resp types.RegisterResponse
	if err := c.request("POST", "/auth/register", req, &resp); err != nil {
		return nil, err
	}

	if resp.Success {
		store.GetAuthStore().SetToken(resp.Token)
		store.GetAuthStore().SetEntity(&resp.Entity)
	}

	return &resp, nil
}

// Logout 用户登出
func (c *APIClient) Logout() error {
	store.GetAuthStore().Logout()
	store.GetGameStore().Clear()
	return nil
}

// GetCharacter 获取角色信息（通过WebSocket获取，这里返回空）
func (c *APIClient) GetCharacter() (*types.Character, error) {
	// 服务端目前没有实现此API，返回默认值
	char := &types.Character{
		ID:               store.GetAuthStore().GetEntityID(),
		Name:             store.GetAuthStore().GetEntityName(),
		Level:            1,
		Health:           100,
		MaxHealth:        100,
		Energy:           50,
		MaxEnergy:        50,
		Attack:           10,
		Defense:          5,
		Speed:            8,
		CultivationRealm: "炼气期",
	}
	store.GetGameStore().SetCharacter(char)
	return char, nil
}

// GetWorldState 获取世界状态（通过WebSocket获取，这里返回空）
func (c *APIClient) GetWorldState() (*types.WorldState, error) {
	// 服务端目前没有实现此API，返回默认值
	world := &types.WorldState{
		CurrentMap:    "新手村",
		PlayersOnline: 0,
		Events:        make([]types.WorldEvent, 0),
		Announcements: make([]types.Announcement, 0),
	}
	store.GetGameStore().SetWorld(world)
	return world, nil
}

// GetSocialInfo 获取社交信息（通过WebSocket获取，这里返回空）
func (c *APIClient) GetSocialInfo() (*types.SocialInfo, error) {
	// 服务端目前没有实现此API，返回默认值
	social := &types.SocialInfo{
		Friends:  make([]types.Friend, 0),
		Messages: make([]types.Message, 0),
		Requests: make([]types.FriendRequest, 0),
	}
	store.GetGameStore().SetSocial(social)
	return social, nil
}

// GetSettings 获取设置
func (c *APIClient) GetSettings() (*types.Settings, error) {
	settings := store.GetGameStore().GetSettings()
	return settings, nil
}

// UpdateSettings 更新设置
func (c *APIClient) UpdateSettings(settings *types.Settings) error {
	store.GetGameStore().SetSettings(settings)
	return nil
}

// StartCombat 开始战斗（通过WebSocket实现）
func (c *APIClient) StartCombat(enemyID string) (*types.CombatState, error) {
	combat := &types.CombatState{
		InCombat:   false,
		BattleLog:  make([]types.CombatLog, 0),
		TurnNumber: 0,
	}
	store.GetGameStore().SetCombat(combat)
	return combat, nil
}

// UseSkill 使用技能（通过WebSocket实现）
func (c *APIClient) UseSkill(skillID string) (*types.CombatState, error) {
	return store.GetGameStore().GetCombat(), nil
}

// FleeCombat 逃跑（通过WebSocket实现）
func (c *APIClient) FleeCombat() (*types.CombatState, error) {
	combat := &types.CombatState{
		InCombat:   false,
		BattleLog:  make([]types.CombatLog, 0),
		TurnNumber: 0,
	}
	store.GetGameStore().SetCombat(combat)
	return combat, nil
}

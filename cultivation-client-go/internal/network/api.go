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
func (c *APIClient) Login(username, password string) (*types.AuthResponse, error) {
	req := types.LoginRequest{Username: username, Password: password}
	var resp types.AuthResponse

	if err := c.request("POST", "/auth/login", req, &resp); err != nil {
		return nil, err
	}

	if resp.Success && resp.Entity != nil {
		store.GetAuthStore().SetToken(resp.Token)
		store.GetAuthStore().SetEntity(&types.Entity{
			ID:         resp.Entity.ID,
			Name:       resp.Entity.Name,
			Realm:      resp.Entity.Realm,
			EntityType: resp.Entity.EntityType,
		})
		// 立即用服务端数据填充角色信息
		store.GetGameStore().SetCharacterFromServerEntity(resp.Entity)
	}

	return &resp, nil
}

// Register 用户注册
func (c *APIClient) Register(username, password string) (*types.AuthResponse, error) {
	req := types.RegisterRequest{Username: username, Password: password}
	var resp types.AuthResponse

	if err := c.request("POST", "/auth/register", req, &resp); err != nil {
		return nil, err
	}

	if resp.Success && resp.Entity != nil {
		store.GetAuthStore().SetToken(resp.Token)
		store.GetAuthStore().SetEntity(&types.Entity{
			ID:         resp.Entity.ID,
			Name:       resp.Entity.Name,
			Realm:      resp.Entity.Realm,
			EntityType: resp.Entity.EntityType,
		})
		// 立即用服务端数据填充角色信息
		store.GetGameStore().SetCharacterFromServerEntity(resp.Entity)
	}

	return &resp, nil
}

// Logout 用户登出
func (c *APIClient) Logout() error {
	store.GetAuthStore().Logout()
	store.GetGameStore().Clear()
	return nil
}

// GetSettings 获取设置
func (c *APIClient) GetSettings() *types.Settings {
	return store.GetGameStore().GetSettings()
}

// UpdateSettings 更新设置
func (c *APIClient) UpdateSettings(settings *types.Settings) error {
	store.GetGameStore().SetSettings(settings)
	return nil
}

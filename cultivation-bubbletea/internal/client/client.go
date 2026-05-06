package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const gatewayBase = "http://localhost:8081"

// Login authenticates with username/password via REST and returns the JWT token and entity ID.
func Login(username, password string) (string, string, error) {
	body, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	resp, err := http.Post(gatewayBase+"/auth/login", "application/json",
		bytes.NewReader(body))
	if err != nil {
		return "", "", fmt.Errorf("登录请求失败: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("登录失败 (%d): %s", resp.StatusCode, string(raw))
	}

	var result struct {
		Token    string `json:"token"`
		EntityID string `json:"entity_id"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", "", fmt.Errorf("解析响应失败: %w", err)
	}
	return result.Token, result.EntityID, nil
}

// Register creates a new account via REST and returns the JWT token and entity ID.
func Register(username, password string) (string, string, error) {
	body, _ := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	resp, err := http.Post(gatewayBase+"/auth/register", "application/json",
		bytes.NewReader(body))
	if err != nil {
		return "", "", fmt.Errorf("注册请求失败: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("注册失败 (%d): %s", resp.StatusCode, string(raw))
	}

	var result struct {
		Token    string `json:"token"`
		EntityID string `json:"entity_id"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", "", fmt.Errorf("解析响应失败: %w", err)
	}
	return result.Token, result.EntityID, nil
}

// SendRESTRequest is a helper for any future REST calls.
func SendRESTRequest(method, path string, token string, payload interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if payload != nil {
		b, _ := json.Marshal(payload)
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, gatewayBase+path, bodyReader)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

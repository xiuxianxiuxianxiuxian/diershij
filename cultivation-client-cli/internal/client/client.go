package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	httpClient = &http.Client{Timeout: 30 * time.Second}
	apiBase    = "http://localhost:8081"
	authToken  string
	entityID   string
)

type authResponse struct {
	Success bool                   `json:"success"`
	Token   string                 `json:"token"`
	Entity  map[string]interface{} `json:"entity"`
	Error   string                 `json:"error,omitempty"`
}

func postJSON(path string, reqBody, respTarget interface{}) error {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", apiBase+path, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Try to extract error from body
		var errResp struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(raw, &errResp) == nil && errResp.Error != "" {
			return fmt.Errorf("%s", errResp.Error)
		}
		if len(raw) > 0 {
			return fmt.Errorf("服务器错误 (%d): %s", resp.StatusCode, string(raw))
		}
		return fmt.Errorf("服务器错误 (%d)", resp.StatusCode)
	}

	return json.Unmarshal(raw, respTarget)
}

func Login(username, password string) (string, string, error) {
	var resp authResponse
	err := postJSON("/auth/login", map[string]string{
		"username": username,
		"password": password,
	}, &resp)
	if err != nil {
		return "", "", err
	}
	if !resp.Success {
		msg := resp.Error
		if msg == "" {
			msg = "登录失败"
		}
		return "", "", fmt.Errorf("%s", msg)
	}
	authToken = resp.Token
	if id, ok := resp.Entity["id"].(string); ok {
		entityID = id
	}
	fmt.Printf("欢迎回来, %s!", getStr(resp.Entity, "name"))
	return resp.Token, entityID, nil
}

func Register(username, password string) (string, string, error) {
	var resp authResponse
	err := postJSON("/auth/register", map[string]string{
		"username": username,
		"password": password,
	}, &resp)
	if err != nil {
		return "", "", err
	}
	if !resp.Success {
		msg := resp.Error
		if msg == "" {
			msg = "注册失败"
		}
		return "", "", fmt.Errorf("%s", msg)
	}
	authToken = resp.Token
	if id, ok := resp.Entity["id"].(string); ok {
		entityID = id
	}
	fmt.Printf("欢迎, %s!", getStr(resp.Entity, "name"))
	return resp.Token, entityID, nil
}

func GetToken() string {
	return authToken
}

package main

import (
    "fmt"
    "log"
    "os"

    "github.com/cultivation-world/gateway/internal/handler"
    "github.com/cultivation-world/gateway/internal/service"
    "github.com/cultivation-world/shared/config"
)

func main() {
    cfg := config.LoadConfigFromEnv()

    gameClient, err := service.NewGameServiceClient(
        getEnv("GAME_SERVER_HOST", "localhost"),
        getEnvInt("GAME_SERVER_PORT", 50051),
    )
    if err != nil {
        log.Fatalf("Failed to connect to game server: %v", err)
    }
    defer gameClient.Close()

    wsHub := handler.NewWebSocketHub()
    go wsHub.Run()

    authSvc := service.NewAuthService(
        getEnv("JWT_SECRET", "cultivation-secret-key"),
        gameClient,
    )

    server := handler.NewServer(cfg, wsHub, authSvc, gameClient)

    log.Printf("Gateway starting on %s:%d", cfg.Server.Host, cfg.Server.Port)
    if err := server.Start(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        var result int
        if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
            return result
        }
    }
    return defaultValue
}

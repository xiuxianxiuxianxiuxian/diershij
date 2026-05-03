package config

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

type DatabaseConfig struct {
    Host            string        `json:"host"`
    Port            int           `json:"port"`
    User            string        `json:"user"`
    Password        string        `json:"password"`
    Database        string        `json:"database"`
    MaxConnections  int           `json:"max_connections"`
    ConnTimeout     time.Duration `json:"conn_timeout"`
}

type RedisConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Password string `json:"password"`
    DB       int    `json:"db"`
}

type ServerConfig struct {
    Name         string        `json:"name"`
    Host         string        `json:"host"`
    Port         int           `json:"port"`
    ReadTimeout  time.Duration `json:"read_timeout"`
    WriteTimeout time.Duration `json:"write_timeout"`
}

type GRPCConfig struct {
    Host string `json:"host"`
    Port int    `json:"port"`
}

type LLMConfig struct {
    Provider    string `json:"provider"`
    APIKey      string `json:"api_key"`
    DailyModel  string `json:"daily_model"`
    ReasonModel string `json:"reason_model"`
    RateLimit   int    `json:"rate_limit"`
    Timeout     int    `json:"timeout"`
}

type HeavenlyDaoConfig struct {
    KarmaDecayRate     float64              `json:"karma_decay_rate"`
    TribulationBase    map[string]float64   `json:"tribulation_base"`
    RealmLifespan      map[string]int       `json:"realm_lifespan"`
    KarmaThresholds    map[string]int       `json:"karma_thresholds"`
}

type Config struct {
    Database     DatabaseConfig     `json:"database"`
    Redis        RedisConfig        `json:"redis"`
    Server       ServerConfig       `json:"server"`
    GRPC         GRPCConfig         `json:"grpc"`
    LLM          LLMConfig          `json:"llm"`
    HeavenlyDao  HeavenlyDaoConfig  `json:"heavenly_dao"`
    WorldEngine  GRPCConfig         `json:"world_engine"`
    DaoService   GRPCConfig         `json:"dao_service"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}

func LoadConfigFromEnv() *Config {
    return &Config{
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnvInt("DB_PORT", 5432),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", "postgres"),
            Database: getEnv("DB_NAME", "cultivation"),
        },
        Redis: RedisConfig{
            Host:     getEnv("REDIS_HOST", "localhost"),
            Port:     getEnvInt("REDIS_PORT", 6379),
            Password: getEnv("REDIS_PASSWORD", ""),
            DB:       getEnvInt("REDIS_DB", 0),
        },
        Server: ServerConfig{
            Host: getEnv("SERVER_HOST", "0.0.0.0"),
            Port: getEnvInt("SERVER_PORT", 8081),
        },
        GRPC: GRPCConfig{
            Host: getEnv("GRPC_HOST", "0.0.0.0"),
            Port: getEnvInt("GRPC_PORT", 50051),
        },
        LLM: LLMConfig{
            Provider:    getEnv("LLM_PROVIDER", "deepseek"),
            APIKey:      getEnv("LLM_API_KEY", ""),
            DailyModel:  getEnv("LLM_DAILY_MODEL", "deepseek-chat"),
            ReasonModel: getEnv("LLM_REASON_MODEL", "deepseek-reasoner"),
            RateLimit:   getEnvInt("LLM_RATE_LIMIT", 600),
            Timeout:     getEnvInt("LLM_TIMEOUT", 10),
        },
        WorldEngine: GRPCConfig{
            Host: getEnv("WORLD_ENGINE_HOST", "localhost"),
            Port: getEnvInt("WORLD_ENGINE_PORT", 50054),
        },
        DaoService: GRPCConfig{
            Host: getEnv("HEAVENLY_DAO_HOST", "localhost"),
            Port: getEnvInt("HEAVENLY_DAO_PORT", 50053),
        },
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

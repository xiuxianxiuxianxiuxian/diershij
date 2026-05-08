<p align="center">
  <a href="README.md">🇨🇳 中文</a> | <a href="README.en.md">🇬🇧 English</a>
</p>

# Cultivation World — Microservice MUD Game

A fully self-evolving multiplayer online text-based cultivation world (MUD). Five-service microservice architecture with AI-driven NPC behavior, dynamic world events, and a heavenly dao karma system.

## Tech Stack

- **Backend**: Go microservices (Go 1.21 / 1.25), gRPC inter-service communication, Gin HTTP framework
- **Database**: PostgreSQL (pgxpool) + Redis (go-redis)
- **Clients**: Terminal CLI (primary) + Bubble Tea TUI (new) + Gio desktop GUI (retained)
- **AI**: DeepSeek LLM-driven NPC decisions + behavior trees + memory system
- **Protocols**: WebSocket JSON + Protocol Buffers (protoc v4.25.3)
- **Infrastructure**: Docker Compose, GitHub Actions CI

## Service Architecture

```
┌─────────────────┐    HTTP/WS     ┌────────────┐
│  CLI / TUI /    │ ◄────────────►│  Gateway   │ :8080 / :8081
│  Gio Clients    │               └─────┬──────┘
└─────────────────┘                      │ gRPC
                  ┌──────────────────────┼──────────────────────┐
                  ▼                      ▼                      ▼
            ┌──────────┐          ┌──────────┐          ┌──────────────┐
            │  Game    │          │ Heavenly │          │     AI       │
            │  Server  │◄────────►│   Dao    │          │  Scheduler   │
            │  :50051  │  gRPC    │  :50053  │          │   :50052     │
            └────┬─────┘          └──────────┘          └──────────────┘
                 │                                         │ LLM
                 ▼                                         ▼
            ┌──────────┐                            ┌────────────┐
            │  World   │                            │  DeepSeek  │
            │  Engine  │ :50054                     │    API     │
            └──────────┘                            └────────────┘
```

### Service Responsibilities

| Service | Port | Responsibilities |
|---------|------|-----------------|
| **Gateway** | 8080 (HTTP/WS), 8081 | JWT auth, WebSocket connection management, message routing |
| **Game Server** | 50051 | Entity management, 31+ operation types, state sync, equipment/items/spells/friends/mail/shop/leaderboard |
| **Heavenly Dao** | 50053 | Heavenly engine: cultivation formulas, breakthrough probability, tribulation judgment, karma |
| **AI Scheduler** | 50052 | NPC decisions: behavior tree + LLM dual loop, memory system, NPC personality, autonomous behavior |
| **World Engine** | 50054 | Region management, resource respawn, world events, world state persistence |

## Quick Start

### Prerequisites

- Go 1.21+ (some modules require 1.25)
- Docker & Docker Compose
- PostgreSQL 15 / Redis 7 (started via Docker)

### One-Click Start (Recommended)

```bash
cd server
docker-compose up -d
```

This starts PostgreSQL, Redis, and all 5 microservices.

### Local Development

```bash
# 1. Start databases
docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=123456 -e POSTGRES_DB=cultivation postgres:15-alpine
docker run -d --name redis -p 6379:6379 redis:7-alpine

# 2. Initialize database (run in order)
PGPASSWORD=123456 psql -h localhost -U postgres -d cultivation -f server/init-db/01_init.sql
# ... run 02-09 migration files in sequence

# 3. Start each service individually
cd server/game-server && go run ./cmd
cd server/heavenly-dao && go run ./cmd
cd server/ai-scheduler && go run ./cmd
cd server/world-engine && go run ./cmd
cd server/gateway && go run ./cmd
```

### Windows Quick Start

```powershell
# Build & start all services
.\start-all.ps1 -rebuild

# Stop all services
.\stop-all.ps1
```

### Launching Clients

```bash
# CLI client (primary)
cd cultivation-client-cli && go run ./cmd

# Bubble Tea TUI client (new)
cd cultivation-bubbletea && go run ./cmd

# Gio desktop client (retained)
cd cultivation-client-go && go run ./cmd
```

The CLI client supports 21+ commands: `cult` (cultivate), `bt` (breakthrough), `exp` (explore), `combat`, `skill` (use skill), `flee`, `spell` (cast spell), `move`, `gather`, `craft`, `cm` (create method), `chat`, `msg` (message), `friend` (add friend), `trade`, `st` (status), `create_sect`, `join_sect`, `leave_sect`.

## Project Layout

```
diershij/
├── server/                        # Server
│   ├── gateway/                   # API Gateway
│   ├── game-server/               # Game Core
│   │   └── internal/
│   │       ├── service/           # Operation dispatcher (operation.go)
│   │       └── repository/        # Persistence layer
│   ├── heavenly-dao/              # Heavenly Dao Engine
│   ├── ai-scheduler/              # AI NPC Scheduler
│   ├── world-engine/              # World Engine
│   ├── shared/                    # Shared Library
│   │   ├── types/                 # Go type definitions
│   │   ├── proto/                 # Protobuf definitions + generated code
│   │   ├── config/                # Configuration
│   │   └── errors/                # Error definitions
│   ├── init-db/                   # SQL migrations (01-09)
│   ├── test_workflows.py          # Integration test script
│   ├── docker-compose.yml
│   └── config.json                # Default config
├── cultivation-client-cli/        # Terminal CLI Client
├── cultivation-bubbletea/         # Bubble Tea TUI Client
├── cultivation-client-go/         # Gio Desktop Client
└── .github/workflows/             # CI configuration
```

## WebSocket Protocol

```json
{"type": "operation", "payload": {"action_type": "cultivate", "params": {}}, "timestamp": 123}
```

**Client → Server**: `operation` (action_type + params)
**Server → Client**: `op_result` / `state_sync` / `entity_update` / `world_event` / `chat` / `error`

## Game Systems

### Realm System

Mortal → Qi Refining → Foundation Establishment → Golden Core → Nascent Soul → Spirit Transformation → Void Refining → Unity → Mahayana → Tribulation Transcendence

### Spiritual Root System

Randomly generates 1-3 spiritual roots on registration, 5% chance of mutation:
- Basic elements: Metal, Wood, Water, Fire, Earth
- Mutated roots: Wind, Lightning, Ice
- Primary root purity: 60-90, secondary roots: 20-50

### Key Formulas

**Cultivation Efficiency** = Base Rate × Root Bonus × Spell Match × Realm Decay × Mental State × (1 - Aging Penalty)

**Breakthrough Probability** = Base Success × Accumulation × Spell Quality × Resource Bonus × Mental State × Luck, clamped to [5%, 80%]

### Mail System

Offline messaging between players and system notifications. Supports send, receive, inbox listing, and auto-cleanup of expired mail.

### Shop System

NPC shop trading with item buy/sell and currency settlement. Interact via the `trade` command.

### Leaderboard

Global player rankings by realm, combat power, etc., updated periodically.

### NPC System

AI Scheduler-driven NPC behavior engine:
- Behavior tree + LLM dual-loop decision making
- NPC personality system (personality traits)
- Memory system (long-term/short-term)
- Autonomous behavior (explore, cultivate, socialize, combat)

## Configuration

| Env Var | Default | Description |
|---------|---------|-------------|
| DB_HOST | localhost | PostgreSQL host |
| DB_PORT | 5432 | PostgreSQL port |
| DB_PASSWORD | postgres | Database password |
| REDIS_HOST | localhost | Redis host |
| REDIS_PORT | 6379 | Redis port |
| JWT_SECRET | cultivation-jwt-secret-key-2024 | JWT signing key |
| LLM_API_KEY | - | DeepSeek API key (AI NPC) |

## Development Status

- [x] Five-service microservice architecture
- [x] Gateway — JWT auth + WebSocket routing
- [x] Game Server — 31 operation types + entity management
- [x] Heavenly Dao — cultivation/breakthrough/tribulation/karma rules engine
- [x] AI Scheduler — behavior tree + LLM decision pipeline
- [x] World Engine — region/resource/event management
- [x] CLI client — 21+ commands, full interaction
- [x] Spiritual root system — random generation + purity + mutation
- [x] Spell system — learn/major/quality affects breakthrough
- [x] Equipment system — 13 stat bonuses + durability
- [x] Combat system — NPC drops + spells + flee
- [x] Shop/trade system — NPC trading, currency settlement
- [x] Mail system — offline messages, system notifications
- [x] Leaderboard — realm/combat power rankings
- [x] World events — beast tides, heavenly anomalies, spirit tides, secret realms, sect wars
- [x] Bubble Tea TUI client — terminal UI
- [ ] Cross-server party dungeons
- [ ] Sect wars / territory control

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A multiplayer text-based MUD (Multi-User Dungeon) with Chinese xianxia/cultivation themes, using Go microservices communicating via gRPC. Players cultivate immortality, form sects, and interact in an AI-driven persistent world.

## Architecture

### Service Topology

```
Client (CLI terminal / WebSocket)
        |
    Gateway (port 8080) — WebSocket + HTTP, JWT auth, message routing
        |
    Game Server (port 50051) — entity management, operation validation, state sync
        |
    Heavenly Dao (port 50053) — karma/retribution, tribulation, world balance
    AI Scheduler (port 50052) — NPC decisions via behavior trees + DeepSeek LLM
    World Engine (port 50054) — regions, resources, world events
```

All backend services use PostgreSQL for persistence and Redis for caching.

### Repository Layout

```
server/
  gateway/        — WebSocket/HTTP gateway, JWT auth, request routing (Gin)
  game-server/    — core game logic, operation dispatch, entity management
  heavenly-dao/   — karma evaluation, tribulation checks, balance (most tested)
  ai-scheduler/   — NPC AI via behavior trees and LLM (DeepSeek)
  world-engine/   — region management, resource spawning, events
  shared/         — shared types, config, protobuf definitions, errors
  init-db/        — SQL schema and migrations (01-08 in order)
  config.json     — default server config
  docker-compose.yml — full orchestration

cultivation-client-go/  — Gio (gioui.org) desktop GUI client (retained)
cultivation-client-cli/ — Terminal CLI client (primary, 21 commands + interactive loop)
```

## Development Commands

```bash
# Start all services via Docker
cd server && docker-compose up -d

# Run individual service (requires PostgreSQL + Redis running locally)
cd server/game-server && go run ./cmd
cd server/gateway && go run ./cmd
cd server/heavenly-dao && go run ./cmd
cd server/ai-scheduler && go run ./cmd
cd server/world-engine && go run ./cmd

# Run CLI client (connects to gateway at localhost:8081)
cd cultivation-client-cli && go run ./cmd

# Build CLI binary
cd cultivation-client-cli && go build -o ../cultivation-cli.exe ./cmd

# Run CLI binary directly
./cultivation-cli.exe

# Run Gio desktop client (retained, unused)
cd cultivation-client-go && go run ./cmd

# Run tests
cd server/heavenly-dao && go test ./internal/service/...
cd server/shared && go test ./types/...

# Run all heavenly-dao tests verbosely
cd server/heavenly-dao && go test -v ./internal/service/...

# Database migrations (run in order)
PGPASSWORD=123456 psql -h localhost -U postgres -d cultivation -f server/init-db/01_init.sql
PGPASSWORD=123456 psql -h localhost -U postgres -d cultivation -f server/init-db/02_game_operations.sql
PGPASSWORD=123456 psql -h localhost -U postgres -d cultivation -f server/init-db/03_fix_missing.sql
PGPASSWORD=123456 psql -h localhost -U postgres -d cultivation -f server/init-db/04_add_extra_attributes.sql
PGPASSWORD=123456 psql -h localhost -U postgres -d cultivation -f server/init-db/05_friends.sql
PGPASSWORD=123456 psql -h localhost -U postgres -d cultivation -f server/init-db/06_add_password_hash.sql
PGPASSWORD=123456 psql -h localhost -U postgres -d cultivation -f server/init-db/07_add_name_unique.sql
PGPASSWORD=123456 psql -h localhost -U postgres -d cultivation -f server/init-db/08_add_spiritual_roots_methods.sql

# Windows: start/stop all backend services locally (double-click .bat or PowerShell)
.\启动游戏.bat
.\停止游戏.bat
# Or from PowerShell:
.\start-all.ps1
.\stop-all.ps1
# start-all.ps1 supports -rebuild flag to recompile all binaries before launching:
.\start-all.ps1 -rebuild

# Integration test (requires running server with test data)
cd server && python test_ws.py
```

### Module Names

- `github.com/cultivation-world/gateway` — gateway (Go 1.25)
- `github.com/cultivation-world/game-server` — game server (Go 1.21)
- `github.com/cultivation-world/heavenly-dao` — heavenly dao (Go 1.21)
- `github.com/cultivation-world/ai-scheduler` — AI scheduler (Go 1.21)
- `github.com/cultivation-world/world-engine` — world engine (Go 1.21)
- `github.com/cultivation-world/shared` — shared library (Go 1.21)
- `cultivation-client` — Gio desktop client (Go 1.23)
- `cultivation-client-cli` — Terminal CLI client (Go 1.23)

All server modules use `replace github.com/cultivation-world/shared => ../shared` in go.mod. No go.work file.

## Key Patterns & Conventions

### Operation Flow
Client sends WebSocket JSON → Gateway `handleOperation()` → gRPC `ExecuteOperation()` → GameServer `OperationService` dispatches by `ActionType` → returns `OperationResult{Success, Message, Effects}` → Gateway wraps in `op_result` message → Client receives in WS handler.

### WebSocket Protocol
Messages use JSON envelope: `{"type": "...", "payload": {...}, "timestamp": 123}`.
- Client → Gateway: `operation` (with `action_type` + `params` in payload), `chat`
- Gateway → Client: `op_result`, `state_sync`, `entity_update`, `world_event`, `chat`, `error`

Key message types: auth, auth_result, operation, op_result, state_sync, entity_update, world_event, chat, system, error.

### Game Actions (21 types)
Actions are dispatched by string ActionType in `game-server/internal/service/operation.go`:
cultivate, breakthrough, combat, explore, gather, craft, create_method, trade, form_sect, join_sect, leave_sect, send_message, cast_spell, meditate, sleep, move, add_friend, remove_friend, accept_friend, flee, use_skill.

### Repository Pattern (game-server)
Each repository is defined as a Go interface (EntityRepository, ItemRepository, InventoryRepository, SpellRepository, etc.) in `game-server/internal/service/operation.go`. Implementations live in `game-server/internal/repository/`. Dependencies are wired explicitly in `cmd/main.go`.

### Configuration
Runtime: `config.LoadConfigFromEnv()` reads env vars (DB_HOST, DB_PORT, etc.).
Fallback: `server/config.json` loaded via `config.LoadConfig(path)`.
Config struct in `server/shared/config/config.go`.

### Database
Core tables: entities, base_attributes, karma_attributes, spirit_stones, world_regions, sects, sect_members, npc_personalities, operation_logs, transactions.
Gameplay tables: items, inventory, recipes, entity_recipes, spells, entity_spells, messages, combat_logs, explore_logs, gather_logs, craft_logs, method_creation_logs, cast_logs.

### Testing
- **heavenly-dao** uses `stretchr/testify` with table-driven tests and factory helpers (`DefaultTribulationInput()`, `DefaultKarmaConfig()`)
- **shared/types** uses standard `testing.T` with `t.Error()`/`t.Errorf()`
- **gateway, game-server, ai-scheduler, world-engine** have no Go tests — only the Python integration test (`test_ws.py`)
- Add new heavenly-dao tests alongside existing rule test files

### Client Architecture

**CLI (primary):** `cultivation-client-cli/` — Interactive terminal app. Single WebSocket connection, async message display in background goroutine. REST auth (login/register), then interactive command loop. 21 commands with aliases (cultivate/cult, meditate/med, breakthrough/bt, explore/exp, gather, move, craft, create_method/cm, combat, use_skill/skill, flee, cast_spell/spell, chat, msg, add_friend/friend, remove_friend/unfriend, accept_friend/accept, form_sect/create_sect, join_sect, leave_sect, trade, status/st, help, clear/cls, exit). Status displays formatted character info from cached entity state.

**Gio (retained):** `cultivation-client-go/` — Original Gio (gioui.org) desktop GUI client. Singleton stores (`store/auth.go`, `store/game.go`) hold JWT token and entity state. WebSocket auto-reconnects with 5-second delay. Operation wrappers (Cultivate, Move, Breakthrough, etc.) in `network/websocket.go` marshal actions as typed JSON messages. Gio GUI pages in `gui/pages/tabs/` (character, combat, world, social, settings, common).

### Service Ports (gRPC)
- Gateway: 8080 (HTTP/WS), 8081 (gRPC)
- Game Server: 50051
- AI Scheduler: 50052
- Heavenly Dao: 50053
- World Engine: 50054

### Proto Definitions
All `.proto` files in `server/shared/proto/`, generated Go code in `server/shared/proto/pb/`.
Regenerate: `protoc --go_out=. --go-grpc_out=. *.proto` from the proto directory.

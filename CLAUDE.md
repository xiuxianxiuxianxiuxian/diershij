# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A multiplayer text-based MUD (Multi-User Dungeon) with Chinese xianxia/cultivation themes, using Go microservices communicating via gRPC. Players cultivate immortality, form sects, and interact in an AI-driven persistent world.

## Architecture

### Service Topology

```
Client (Gio GUI / WebSocket)
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
  gateway/        — WebSocket/HTTP gateway, JWT auth, request routing
  game-server/    — core game logic, CRUD for entities/items/spells
  heavenly-dao/   — karma evaluation, tribulation checks, balance
  ai-scheduler/   — NPC AI via behavior trees and LLM (DeepSeek)
  world-engine/   — region management, resource spawning, events
  shared/         — shared types, config, protobuf definitions, errors
  init-db/        — SQL schema and migrations
  config.json     — default server config
  docker-compose.yml — full orchestration

cultivation-client-go/  — Gio (gioui.org) desktop GUI client
  cmd/main.go           — entry point
  internal/app/         — app lifecycle, view routing, WS handlers
  internal/network/     — HTTP + WebSocket clients
  internal/store/       — auth/game state stores
  internal/types/       — client-side type definitions
```

### Communication Patterns

- **Client ↔ Gateway**: WebSocket JSON messages with type/payload format
- **Gateway ↔ Backend**: gRPC (protobuf-defined services)
- **Service ↔ Service**: gRPC via shared proto definitions
- **Authentication**: JWT obtained via HTTP POST /auth/login, then sent as WebSocket query param

### Proto Services (in `server/shared/proto/`)

| Service | Key RPCs |
|---------|----------|
| GameService | ExecuteOperation, GetEntity, CreateEntity, AuthenticateEntity, SyncState, StreamEntityUpdates |
| HeavenlyDaoService | EvaluateKarma, CheckTribulation, BalanceCheck, ApplyKarmaDecay |
| AISchedulerService | ScheduleDecision, ExecuteBehaviorTree, RegisterNPC, UnregisterNPC |
| WorldService | GetRegion, SpawnResources, TriggerEvent, GetWorldState |

### Core Domain Model

- **Entity** — player or NPC with 80+ attributes across categories (basic, cultivation, combat, spiritual roots, mental state, life skills, social, wealth, special, law, Dao, lifespan, status effects)
- **Cultivation Realm** — 10-tier hierarchy (mortal → qi_condensation → foundation → golden_core → nascent_soul → soul_transformation → void_refinement → integration → mahayana → tribulation)
- **Karma** — influences tribulation probability, world favor, NPC reactions
- **Region** — world areas with spiritual density, danger level, resources

### Database Schema

Main tables: `entities`, `base_attributes`, `karma_attributes`, `spirit_stones`, `world_regions`, `sects`, `sect_members`, `npc_personalities`, `operation_logs`, `transactions`

## Development Commands

```bash
# Start all server services via Docker
cd server && docker-compose up -d

# Or start services individually (requires PostgreSQL + Redis running)
cd server/game-server && go run ./cmd
cd server/gateway && go run ./cmd
cd server/heavenly-dao && go run ./cmd
cd server/ai-scheduler && go run ./cmd
cd server/world-engine && go run ./cmd

# Start desktop client
cd cultivation-client-go && go run ./cmd

# Initialise database (standalone)
psql -h localhost -U postgres -d cultivation -f server/init-db/01_init.sql

# Windows: start everything (PowerShell)
.\start-all.ps1
```

All microservices share a pattern: `go run ./cmd` from the service directory. Each service has its own `go.mod`.

### Project Module Names

- `github.com/cultivation-world/gateway` — gateway
- `github.com/cultivation-world/game-server` — game server
- `github.com/cultivation-world/heavenly-dao` — heavenly dao
- `github.com/cultivation-world/ai-scheduler` — AI scheduler
- `github.com/cultivation-world/world-engine` — world engine
- `github.com/cultivation-world/shared` — shared library
- `cultivation-client` — Gio desktop client

## Key Conventions

- All inter-service communication uses gRPC with protobuf (proto files in `server/shared/proto/`)
- Gateway is the only public-facing service; backend services are internal
- Client-server real-time communication uses typed JSON WebSocket messages with a `type`/`payload` envelope
- Configs use env vars (via `config.LoadConfigFromEnv()`) or `config.json`
- Shared types are defined in `server/shared/types/` and mirrored in proto definitions

# Implementation Task List: autonomous-cultivation-npc

Feature Name: autonomous-cultivation-npc
Created: 2026-05-03
Updated: 2026-05-03 (Phase 1 complete)

## Overview

This task list covers the implementation of business logic and tests for the Cultivation World MUD game. The project has a completed basic framework; this list focuses on filling in the business logic gaps and adding comprehensive tests.

## Phase 1: Core Type System Extension

- [x] 1.1 Extend Attributes type to cover all 83+ properties (entity.go)
  - [x] 1.1.1 Add spiritual roots system (SpiritualRoot struct + list)
  - [x] 1.1.2 Add combat attributes (crit_rate, dodge_rate, penetration, etc.)
  - [x] 1.1.3 Add mental state attributes (obsession_count, dao_heart, inner_demon_resistance, enlightenment)
  - [x] 1.1.4 Add life skills (alchemy_level, artificing_level, formation_level, etc.)
  - [x] 1.1.5 Add social attributes (reputation, sect_contribution, relationship lists)
  - [x] 1.1.6 Add wealth attributes (spirit_stones struct, property_value)
  - [x] 1.1.7 Add special attributes (bloodline, physique, destiny, world_favor)
  - [x] 1.1.8 Add law attributes (laws map, law_resonance, domain_power, etc.)
  - [x] 1.1.9 Add dao attributes (dao_seed_type, dao_seed_level, dao_marks, etc.)
  - [x] 1.1.10 Add status effects (injuries, buffs, debuffs, poison_level, curse_level)

- [x] 1.2 Add CultivationMethod type (60+ properties)
  - [x] 1.2.1 Define CultivationMethod struct with all attributes
  - [x] 1.2.2 Define Skill struct (active/ultimate skills)
  - [x] 1.2.3 Add method rank, category, element affinity constants
  - [x] 1.2.4 Write tests for method type initialization and validation

- [x] 1.3 Add Item/Equipment types
  - [x] 1.3.1 Define ItemTemplate struct
  - [x] 1.3.2 Define EntityInventory struct
  - [x] 1.3.3 Define EntityEquipment struct
  - [x] 1.3.4 Define Pill, Artifact, Talisman, Recipe, Material types
  - [x] 1.3.5 Write tests

- [x] 1.4 Add Sect and Social types
  - [x] 1.4.1 Define Sect struct (name, founder, philosophy, territory, rules, alignment)
  - [x] 1.4.2 Define SectMember struct (rank, contribution)
  - [x] 1.4.3 Define Relationship struct (type, strength, history)
  - [x] 1.4.4 Define NPCPersonality struct
  - [x] 1.4.5 Write tests

## Phase 2: Heavenly Dao Engine (16 Algorithms)

- [x] 2.1 Karma Algorithm (因果业力算法)
  - [x] 2.1.1 Implement calculateKarmaChange with context multiplier, relationship modifier
  - [x] 2.1.2 Implement applyKarmaDecay with exponential decay
  - [x] 2.1.3 Implement calculateHeavenlyMark from karma thresholds
  - [x] 2.1.4 Write tests: karma change for different actions, decay over time, mark transitions

- [x] 2.2 Tribulation Algorithm (天劫触发算法)
  - [x] 2.2.1 Implement calculateTribulationProbability with karma/merit/luck factors
  - [x] 2.2.2 Implement calculateTribulationStrength with recent breakthrough multiplier
  - [x] 2.2.3 Write tests: probability ranges, strength calculations, edge cases

- [x] 2.3 Lifespan Algorithm (寿命与衰老算法)
  - [x] 2.3.1 Implement calculateRemainingLifespan
  - [x] 2.3.2 Implement calculateAgingPenalty (tiered by remaining percentage)
  - [x] 2.3.3 Implement handleLifespanDepletion (trigger breakthrough or death)
  - [x] 2.3.4 Write tests: aging at different life stages, lifespan per realm

- [x] 2.4 Cultivation Efficiency Algorithm (修炼效率算法)
  - [x] 2.4.1 Implement calculateCultivationRate with all factors (comprehension, spiritual, method match, realm, mental state, aging)
  - [x] 2.4.2 Implement calculateMethodCompatibility (spiritual root vs method requirements)
  - [x] 2.4.3 Write tests: rate calculation with various factors, compatibility scoring

- [x] 2.5 Breakthrough Algorithm (突破成功率算法)
  - [x] 2.5.1 Implement calculateBreakthroughSuccess with all modifiers
  - [x] 2.5.2 Implement calculateBreakthroughFailurePenalty
  - [x] 2.5.3 Integrate with tribulation check before breakthrough
  - [x] 2.5.4 Write tests: success probability, failure penalties, tribulation integration

- [x] 2.6 Combat Damage Algorithm (战斗伤害算法)
  - [x] 2.6.1 Implement calculateDamage with realm suppression, element interaction, defense reduction, method bonus, crit
  - [x] 2.6.2 Implement calculateElementInteraction
  - [x] 2.6.3 Implement resolveCombat (turn-based combat loop)
  - [x] 2.6.4 Write tests: damage calculation, element interactions, combat resolution

- [x] 2.7 Method Compatibility Algorithm (功法冲突与兼容算法)
  - [x] 2.7.1 Implement calculateMethodConflict (element opposite, alignment, realm diff)
  - [x] 2.7.2 Implement calculateMethodBacklash probability
  - [x] 2.7.3 Write tests: conflict scoring, backlash probability tiers

- [x] 2.8 Alchemy Algorithm (丹药炼制算法)
  - [x] 2.8.1 Implement calculatePillSuccessRate
  - [x] 2.8.2 Implement generatePillQuality
  - [x] 2.8.3 Implement handleAlchemyFailure
  - [x] 2.8.4 Write tests: success rate factors, quality distribution, failure consequences

- [x] 2.9 Artifact Forging Algorithm (法宝炼制算法)
  - [x] 2.9.1 Implement calculateForgingSuccessRate
  - [x] 2.9.2 Implement generateArtifactGrade
  - [x] 2.9.3 Write tests

- [x] 2.10 Formation Algorithm (阵法效果算法)
  - [x] 2.10.1 Implement calculateFormationPower
  - [x] 2.10.2 Implement calculateBreakFormationRate
  - [x] 2.10.3 Write tests

- [x] 2.11 Fortune Algorithm (气运与机缘算法)
  - [x] 2.11.1 Implement calculateFortune
  - [x] 2.11.2 Implement calculateOpportunityRate
  - [x] 2.11.3 Write tests

- [x] 2.12 Spiritual Root Algorithm (灵根觉醒算法)
  - [x] 2.12.1 Implement calculateAwakeningRate
  - [x] 2.12.2 Implement generateMutatedRoot
  - [x] 2.12.3 Write tests

- [x] 2.13 Spiritual Tide Algorithm (灵气潮汐算法)
  - [x] 2.13.1 Implement calculateCurrentTide
  - [x] 2.13.2 Implement adjustSpiritualDensity
  - [x] 2.13.3 Write tests

- [x] 2.14 Demon Beast Algorithm (妖兽生成算法)
  - [x] 2.14.1 Implement calculateBeastSpawnRate
  - [x] 2.14.2 Implement generateBeastLevelDistribution
  - [x] 2.14.3 Write tests

- [x] 2.15 Sect Fortune Algorithm (宗门气运算法)
  - [x] 2.15.1 Implement calculateSectFortune
  - [x] 2.15.2 Implement predictSectTrajectory
  - [x] 2.15.3 Write tests

- [x] 2.16 World Balance Algorithm (世界平衡算法)
  - [x] 2.16.1 Implement evaluateWorldHealth (Gini coefficient, resource circulation, sect diversity, karma stddev)
  - [x] 2.16.2 Implement applyBalanceAdjustment
  - [x] 2.16.3 Write tests

## Phase 3: Game Server Operations

- [x] 3.1 Combat Operation
  - [x] 3.1.1 Implement executeCombat (initiate combat, turn-based resolution)
  - [x] 3.1.2 Implement combat result processing (karma, loot, injuries)
  - [x] 3.1.3 Write tests: combat initiation, resolution, results

- [x] 3.2 Explore Operation
  - [x] 3.2.1 Implement executeExplore (region exploration, random events)
  - [x] 3.2.2 Implement opportunity trigger (fortune-based)
  - [x] 3.2.3 Write tests

- [x] 3.3 Gather Operation
  - [x] 3.3.1 Implement executeGather (resource collection)
  - [x] 3.3.2 Implement resource depletion tracking
  - [x] 3.3.3 Write tests

- [x] 3.4 Craft Operation (炼丹/炼器/阵法)
  - [x] 3.4.1 Implement executeCraft (dispatch to alchemy/forging/formation)
  - [x] 3.4.2 Integrate with Heavenly Dao algorithms
  - [x] 3.4.3 Write tests

- [x] 3.5 Create Method Operation (功法自创)
  - [x] 3.5.1 Implement executeCreateMethod with premium spirit stone validation (10000)
  - [x] 3.5.2 Implement method validation and generation
  - [x] 3.5.3 Implement creator reward tracking (5-10 premium stones at 10 learners)
  - [x] 3.5.4 Write tests

- [x] 3.6 Trade Operation
  - [x] 3.6.1 Implement executeTrade (buyer/seller exchange)
  - [x] 3.6.2 Implement heavenly tax calculation (progressive rate)
  - [x] 3.6.3 Implement spirit stone exchange (with fees)
  - [x] 3.6.4 Write tests

- [x] 3.7 Sect Operations (创建/加入/退出)
  - [x] 3.7.1 Implement executeFormSect
  - [x] 3.7.2 Implement executeJoinSect
  - [x] 3.7.3 Implement sect contribution tracking
  - [x] 3.7.4 Write tests

- [x] 3.8 Cast Spell Operation
  - [x] 3.8.1 Implement executeCastSpell
  - [x] 3.8.2 Integrate with skill system
  - [x] 3.8.3 Write tests

- [x] 3.9 Operation Service Refactoring
  - [x] 3.9.1 Refactor executeBreakthrough to use Heavenly Dao algorithms
  - [x] 3.9.2 Refactor executeCultivate to use cultivation efficiency algorithm
  - [x] 3.9.3 Add cooldown system for all operations
  - [x] 3.9.4 Write integration tests for operation service

## Phase 4: World Engine Enhancement

- [x] 4.1 Complete World Initialization
  - [x] 4.1.1 Add remaining regions (all 15+ from design doc)
  - [x] 4.1.2 Add initial NPC spawning (50-100 with distribution)
  - [x] 4.1.3 Add initial sects (青云宗, 血煞殿, 天机阁)
  - [x] 4.1.4 Add world lore/history entries
  - [x] 4.1.5 Write tests

- [x] 4.2 Resource Spawn Algorithm
  - [x] 4.2.1 Implement calculateSpawnRate with spiritual/pressure/balance factors
  - [x] 4.2.2 Implement calculateRareSpawnChance
  - [x] 4.2.3 Implement periodic resource refresh ticker
  - [x] 4.2.4 Write tests

- [x] 4.3 World Event System
  - [x] 4.3.1 Implement world crisis generation (demon beast tide, heavenly anomaly)
  - [x] 4.3.2 Implement event lifecycle (start, progress, resolve)
  - [x] 4.3.3 Implement event participant rewards
  - [x] 4.3.4 Write tests

- [x] 4.4 Epoch System
  - [x] 4.4.1 Implement advanceWorldEpoch
  - [x] 4.4.2 Implement epoch transition criteria check
  - [x] 4.4.3 Implement new epoch generation (spiritual reset, resource redistribution, new secret realms)
  - [x] 4.4.4 Write tests

- [x] 4.5 Premium Spirit Stone System
  - [x] 4.5.1 Implement allowed sources hardcoding
  - [x] 4.5.2 Implement grant with duplicate prevention
  - [x] 4.5.3 Implement annual production cap
  - [x] 4.5.4 Write tests

## Phase 5: AI Scheduler Enhancement

- [x] 5.1 Behavior Tree Engine
  - [x] 5.1.1 Implement proper behavior tree node types (sequence, selector, decorator, leaf)
  - [x] 5.1.2 Implement tree evaluation with deterministic output
  - [x] 5.1.3 Implement NPC state perception module
  - [x] 5.1.4 Write tests: tree evaluation, node behavior, determinism

- [x] 5.2 LLM Integration (DeepSeek)
  - [x] 5.2.1 Implement DeepSeek API client (OpenAI-compatible format)
  - [x] 5.2.1.1 Implement token bucket rate limiter (600 RPM)
  - [x] 5.2.1.2 Implement circuit breaker for API failures
  - [x] 5.2.2 Implement LLM decision call with system prompt template
  - [x] 5.2.3 Implement response parsing and validation
  - [x] 5.2.4 Implement timeout fallback to behavior tree
  - [x] 5.2.5 Write tests: API client, rate limiter, circuit breaker, fallback chain

- [x] 5.3 Template Matching Enhancement
  - [x] 5.3.1 Implement similarity calculation (context vs template pattern)
  - [x] 5.3.2 Implement template library with 500+ behavior, 1000+ dialogue, 200+ decision templates
  - [x] 5.3.3 Implement learning mechanism (add LLM results to template library)
  - [x] 5.3.4 Write tests

- [x] 5.4 NPC Decision Pipeline
  - [x] 5.4.1 Implement decision cycle loop (behavior tree 1-5s, LLM 30-120s)
  - [x] 5.4.2 Implement context building from entity state
  - [x] 5.4.3 Implement action execution from decision result
  - [x] 5.4.4 Write integration tests

## Phase 6: Database Layer

- [ ] 6.1 Entity Repository Extension
  - [ ] 6.1.1 Add CRUD for all new attribute tables
  - [ ] 6.1.2 Add methods for spiritual roots, methods, inventory, equipment
  - [ ] 6.1.3 Add methods for relationships, sect memberships
  - [ ] 6.1.4 Write repository tests

- [ ] 6.2 New Repositories
  - [ ] 6.2.1 Implement SectRepository
  - [ ] 6.2.2 Implement TransactionRepository
  - [ ] 6.2.3 Implement MethodRepository
  - [ ] 6.2.4 Implement ItemRepository
  - [ ] 6.2.5 Implement RelationshipRepository
  - [ ] 6.2.6 Implement NPCPersonalityRepository
  - [ ] 6.2.7 Implement WorldEventRepository
  - [ ] 6.2.8 Write repository tests for each

- [ ] 6.3 Redis Cache Layer
  - [ ] 6.3.1 Implement entity state caching with TTL
  - [ ] 6.3.2 Implement cache invalidation on updates
  - [ ] 6.3.3 Implement distributed lock for concurrent operations
  - [ ] 6.3.4 Write cache tests

## Phase 7: Gateway Enhancement

- [ ] 7.1 WebSocket Protocol Enhancement
  - [ ] 7.1.1 Implement message validation and schema enforcement
  - [ ] 7.1.2 Implement region-based broadcasting (nearby entities only)
  - [ ] 7.1.3 Implement message rate limiting per client
  - [ ] 7.1.4 Write tests

- [ ] 7.2 Authentication Enhancement
  - [ ] 7.2.1 Implement proper JWT token generation and validation
  - [ ] 7.2.2 Implement token refresh
  - [ ] 7.2.3 Implement session management
  - [ ] 7.2.4 Write tests

- [ ] 7.3 gRPC Client Pool
  - [ ] 7.3.1 Implement connection pooling for game server
  - [ ] 7.3.2 Implement health check and reconnection
  - [ ] 7.3.3 Implement load balancing across server nodes
  - [ ] 7.3.4 Write tests

## Phase 8: Configuration

- [ ] 8.1 Heavenly Dao Config
  - [ ] 8.1.1 Create configs/heavenly-dao-config.yaml with all algorithm parameters
  - [ ] 8.1.2 Implement YAML config loader
  - [ ] 8.1.3 Implement hot-reload support
  - [ ] 8.1.4 Write tests

- [ ] 8.2 World Config
  - [ ] 8.2.1 Create configs/world-config.yaml (regions, resources, initial NPCs)
  - [ ] 8.2.2 Implement config loader
  - [ ] 8.2.3 Write tests

## Phase 9: Integration & System Tests

- [ ] 9.1 Player/NPC Operation Consistency Test
  - [ ] 9.1.1 Test same operation by player and NPC yields same result
  - [ ] 9.1.2 Test operation validation logic is identical

- [ ] 9.2 World Evolution Test
  - [ ] 9.2.1 Test 24-hour NPC behavior cycle
  - [ ] 9.2.2 Test resource spawn and consumption cycle
  - [ ] 9.2.3 Test karma accumulation and tribulation trigger

- [ ] 9.3 Economy Test
  - [ ] 9.3.1 Test trade system with no system pricing
  - [ ] 9.3.2 Test spirit stone exchange with fees
  - [ ] 9.3.3 Test premium spirit stone strict acquisition

- [ ] 9.4 Concurrency Test
  - [ ] 9.4.1 Test concurrent operations on same entity
  - [ ] 9.4.2 Test concurrent trades between entities
  - [ ] 9.4.3 Test concurrent combat scenarios

- [ ] 9.5 Chaos Test
  - [ ] 9.5.1 Test entity state recovery after server restart
  - [ ] 9.5.2 Test cache consistency after Redis failure

## Phase 10: Documentation & Cleanup

- [ ] 10.1 Update API Documentation
- [ ] 10.2 Update Database Schema Documentation
- [ ] 10.3 Add GoDoc comments to all public functions
- [ ] 10.4 Code review and cleanup (remove dead code, fix lints)

## Dependencies Between Phases

```
Phase 1 (Types) → Phase 2 (Heavenly Dao) → Phase 3 (Operations)
       ↓                                      ↓
Phase 4 (World) → Phase 5 (AI) ←─────────────┘
       ↓                    ↓
Phase 6 (Database) ───────→ Phase 7 (Gateway)
       ↓
Phase 8 (Config)
       ↓
Phase 9 (Integration Tests)
       ↓
Phase 10 (Documentation)
```

## Testing Framework

- **Unit Tests**: Go `testing` package + `testify/assert` for assertions
- **Test Location**: `_test.go` files alongside source files
- **Naming Convention**: `Test{FunctionName}{Scenario}` (e.g., `TestCalculateKarmaChange_KillInnocent`)
- **Coverage Target**: >80% for all new code

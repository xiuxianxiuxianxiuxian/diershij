package store

import (
	"sync"

	"cultivation-client/internal/types"
)

type GameStore struct {
	mu                  sync.RWMutex
	character           *types.Character
	world               *types.WorldState
	combat              *types.CombatState
	social              *types.SocialInfo
	settings            *types.Settings
	lastOperationResult *types.OperationResult
}

var gameInstance *GameStore
var gameOnce sync.Once

func GetGameStore() *GameStore {
	gameOnce.Do(func() {
		gameInstance = &GameStore{
			settings: &types.Settings{
				AudioVolume:       0.8,
				MusicVolume:       0.6,
				ShowDamageNumbers: true,
				AutoPlay:          false,
				ShowFPS:           false,
				Language:          "zh_CN",
			},
			combat: &types.CombatState{
				InCombat:   false,
				BattleLog:  make([]types.CombatLog, 0),
				TurnNumber: 0,
			},
			social: &types.SocialInfo{
				Friends:  make([]types.Friend, 0),
				Messages: make([]types.Message, 0),
				Requests: make([]types.FriendRequest, 0),
			},
		}
	})
	return gameInstance
}

func (s *GameStore) SetCharacter(character *types.Character) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.character = character
}

func (s *GameStore) GetCharacter() *types.Character {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.character
}

func (s *GameStore) SetWorld(world *types.WorldState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.world = world
}

func (s *GameStore) GetWorld() *types.WorldState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.world
}

func (s *GameStore) SetCombat(combat *types.CombatState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.combat = combat
}

func (s *GameStore) GetCombat() *types.CombatState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.combat
}

func (s *GameStore) SetSocial(social *types.SocialInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.social = social
}

func (s *GameStore) GetSocial() *types.SocialInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.social
}

func (s *GameStore) SetSettings(settings *types.Settings) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.settings = settings
}

func (s *GameStore) GetSettings() *types.Settings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.settings
}

func (s *GameStore) UpdateSettings(f func(*types.Settings)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f(s.settings)
}

func (s *GameStore) AddCombatLog(log types.CombatLog) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.combat.BattleLog = append(s.combat.BattleLog, log)
	if len(s.combat.BattleLog) > 100 {
		s.combat.BattleLog = s.combat.BattleLog[len(s.combat.BattleLog)-100:]
	}
}

func (s *GameStore) AddMessage(msg types.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.social.Messages = append(s.social.Messages, msg)
}

func (s *GameStore) SetLastOperationResult(result *types.OperationResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastOperationResult = result
}

func (s *GameStore) GetLastOperationResult() *types.OperationResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastOperationResult
}

// SetCharacterFromServerMap 从 state_sync/entity_update 的 map 数据填充角色
func (s *GameStore) SetCharacterFromServerMap(rawEntity map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.setCharacterFromMap(rawEntity)
}

// SetCharacterFromServerEntity 从 API 登录/注册响应的 ServerEntity 填充角色
func (s *GameStore) SetCharacterFromServerEntity(se *types.ServerEntity) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.character == nil {
		s.character = &types.Character{}
	}
	s.character.ID = se.ID
	s.character.Name = se.Name
	s.character.CultivationRealm = se.Realm
	s.character.Level = realmToLevel(se.Realm)
	if se.Attributes != nil {
		s.applyAttributes(se.Attributes)
	}
}

func (s *GameStore) setCharacterFromMap(rawEntity map[string]interface{}) {
	if s.character == nil {
		s.character = &types.Character{}
	}
	if id, ok := rawEntity["id"].(string); ok {
		s.character.ID = id
	}
	if name, ok := rawEntity["name"].(string); ok {
		s.character.Name = name
	}
	if realm, ok := rawEntity["realm"].(string); ok {
		s.character.CultivationRealm = realm
		s.character.Level = realmToLevel(realm)
	}
	if status, ok := rawEntity["status"].(string); ok {
		s.character.Status = status
	}
	if attrs, ok := rawEntity["attributes"].(map[string]interface{}); ok {
		s.applyAttributes(attrs)
	}
	if pos, ok := rawEntity["position"].(map[string]interface{}); ok {
		if regionID, ok := pos["region_id"].(string); ok {
			s.character.RegionID = regionID
		}
		if x, ok := getFloat64(pos, "x"); ok {
			s.character.PositionX = x
		}
		if y, ok := getFloat64(pos, "y"); ok {
			s.character.PositionY = y
		}
	}
	if karma, ok := rawEntity["karma"].(map[string]interface{}); ok {
		if kv, ok := getInt(karma, "karma_value"); ok {
			s.character.KarmaValue = kv
		}
		if m, ok := getInt(karma, "merit"); ok {
			s.character.Merit = m
		}
		if hm, ok := karma["heavenly_mark"].(string); ok {
			s.character.HeavenlyMark = hm
		}
		if kd, ok := getInt(karma, "karmic_debt"); ok {
			s.character.KarmicDebt = kd
		}
	}
}

func (s *GameStore) applyAttributes(attrs map[string]interface{}) {
		if ss, ok := attrs["spirit_stones"].(map[string]interface{}); ok {
			if lg, ok := getInt64(ss, "low_grade"); ok {
				s.character.LowGradeStones = lg
			}
			if mg, ok := getInt64(ss, "medium_grade"); ok {
				s.character.MediumGradeStones = mg
			}
			if hg, ok := getInt64(ss, "high_grade"); ok {
				s.character.HighGradeStones = hg
			}
			if pg, ok := getInt64(ss, "premium_grade"); ok {
				s.character.PremiumGradeStones = pg
			}
		}
	if qi, ok := getFloat64(attrs, "qi"); ok {
		s.character.Qi = qi
		s.character.Health = int(qi)
	}
	if maxQi, ok := getFloat64(attrs, "max_qi"); ok {
		s.character.MaxQi = maxQi
		s.character.MaxHealth = int(maxQi)
	}
	if sp, ok := getFloat64(attrs, "spiritual_power"); ok {
		s.character.SpiritualPower = sp
		s.character.Energy = int(sp)
	}
	if maxSp, ok := getFloat64(attrs, "max_spiritual_power"); ok {
		s.character.MaxSpiritualPower = maxSp
		s.character.MaxEnergy = int(maxSp)
	}
	if atk, ok := getFloat64(attrs, "attack_power"); ok {
		s.character.Attack = int(atk)
	}
	if def, ok := getFloat64(attrs, "defense"); ok {
		s.character.Defense = int(def)
	}
	if speed, ok := getFloat64(attrs, "speed"); ok {
		s.character.Speed = int(speed)
	}
	if prog, ok := getFloat64(attrs, "cultivation_progress"); ok {
		s.character.CultivationProgress = prog
	}
	if comp, ok := getInt(attrs, "comprehension"); ok {
		s.character.Comprehension = comp
	}
	if cons, ok := getInt(attrs, "constitution"); ok {
		s.character.Constitution = cons
	}
	if luck, ok := getInt(attrs, "luck"); ok {
		s.character.Luck = luck
	}
	if ds, ok := getFloat64(attrs, "divine_sense"); ok {
		s.character.DivineSense = ds
	}
	if ms, ok := getInt(attrs, "mental_stability"); ok {
		s.character.MentalStability = ms
	}
	if rl, ok := getInt(attrs, "remaining_lifespan"); ok {
		s.character.RemainingLifespan = rl
	}
	if ml, ok := getInt(attrs, "max_lifespan"); ok {
		s.character.MaxLifespan = ml
	}
	// 战斗属性
	if cr, ok := getFloat64(attrs, "crit_rate"); ok {
		s.character.CritRate = cr
	}
	if cd, ok := getFloat64(attrs, "crit_damage"); ok {
		s.character.CritDamage = cd
	}
	if dr, ok := getFloat64(attrs, "dodge_rate"); ok {
		s.character.DodgeRate = dr
	}
	if hr, ok := getFloat64(attrs, "hit_rate"); ok {
		s.character.HitRate = hr
	}
	if pen, ok := getFloat64(attrs, "penetration"); ok {
		s.character.Penetration = pen
	}
	if dmgRed, ok := getFloat64(attrs, "damage_reduction"); ok {
		s.character.DamageReduction = dmgRed
	}
	// 生活技能
	if al, ok := getInt(attrs, "alchemy_level"); ok {
		s.character.AlchemyLevel = al
	}
	if arl, ok := getInt(attrs, "artificing_level"); ok {
		s.character.ArtificingLevel = arl
	}
	if fl, ok := getInt(attrs, "formation_level"); ok {
		s.character.FormationLevel = fl
	}
	if fc, ok := getInt(attrs, "fire_control"); ok {
		s.character.FireControl = fc
	}
	if hk, ok := getInt(attrs, "herb_knowledge"); ok {
		s.character.HerbKnowledge = hk
	}
	if ms, ok := getInt(attrs, "mining_skill"); ok {
		s.character.MiningSkill = ms
	}
	if ts, ok := getInt(attrs, "talisman_skill"); ok {
		s.character.TalismanSkill = ts
	}
	if bt, ok := getInt(attrs, "beast_taming"); ok {
		s.character.BeastTaming = bt
	}
	// 社交
	if rep, ok := getInt(attrs, "reputation"); ok {
		s.character.Reputation = rep
	}
	if sc, ok := getInt(attrs, "sect_contribution"); ok {
		s.character.SectContribution = sc
	}
	// 心境 / 灵性
	if dh, ok := getInt(attrs, "dao_heart"); ok {
		s.character.DaoHeart = dh
	}
	if en, ok := getInt(attrs, "enlightenment"); ok {
		s.character.Enlightenment = en
	}
	if rp, ok := getInt(attrs, "root_purity"); ok {
		s.character.RootPurity = rp
	}
	if pl, ok := getInt(attrs, "poison_level"); ok {
		s.character.PoisonLevel = pl
	}
	if cl, ok := getInt(attrs, "curse_level"); ok {
		s.character.CurseLevel = cl
	}
}

func getFloat64(m map[string]interface{}, key string) (float64, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint64:
		return float64(n), true
	default:
		return 0, false
	}
}

func getInt64(m map[string]interface{}, key string) (int64, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return int64(n), true
	case int:
		return int64(n), true
	case int64:
		return n, true
	case uint64:
		return int64(n), true
	default:
		return 0, false
	}
}

func getInt(m map[string]interface{}, key string) (int, bool) {
	v, ok := m[key]
	if !ok {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	case uint64:
		return int(n), true
	default:
		return 0, false
	}
}

func (s *GameStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.character = nil
	s.world = nil
	s.combat = &types.CombatState{
		InCombat:   false,
		BattleLog:  make([]types.CombatLog, 0),
		TurnNumber: 0,
	}
	s.social = &types.SocialInfo{
		Friends:  make([]types.Friend, 0),
		Messages: make([]types.Message, 0),
		Requests: make([]types.FriendRequest, 0),
	}
	s.lastOperationResult = nil
}

// realmToLevel 根据修仙境界返回等级（1-100）
func realmToLevel(realm string) int {
	switch realm {
	case "mortal":
		return 1
	case "qi_condensation":
		return 10
	case "foundation":
		return 20
	case "golden_core":
		return 30
	case "nascent_soul":
		return 40
	case "soul_transformation":
		return 50
	case "void_refinement":
		return 60
	case "integration":
		return 70
	case "mahayana":
		return 80
	case "tribulation":
		return 90
	default:
		return 1
	}
}

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

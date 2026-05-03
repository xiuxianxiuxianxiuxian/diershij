package store

import (
	"sync"

	"cultivation-client/internal/types"
)

type AuthStore struct {
	mu         sync.RWMutex
	token      string
	user       *types.User
	entity     *types.Entity
	isLoggedIn bool
}

var authInstance *AuthStore
var authOnce sync.Once

func GetAuthStore() *AuthStore {
	authOnce.Do(func() {
		authInstance = &AuthStore{
			isLoggedIn: false,
		}
	})
	return authInstance
}

func (s *AuthStore) SetToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = token
}

func (s *AuthStore) GetToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.token
}

func (s *AuthStore) SetUser(user *types.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.user = user
	s.isLoggedIn = true
}

func (s *AuthStore) GetUser() *types.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.user
}

func (s *AuthStore) SetEntity(entity *types.Entity) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entity = entity
	s.isLoggedIn = true
}

func (s *AuthStore) GetEntity() *types.Entity {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.entity
}

func (s *AuthStore) GetEntityID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.entity != nil {
		return s.entity.ID
	}
	return ""
}

func (s *AuthStore) GetEntityName() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.entity != nil {
		return s.entity.Name
	}
	return ""
}

func (s *AuthStore) IsLoggedIn() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isLoggedIn
}

func (s *AuthStore) Logout() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.token = ""
	s.user = nil
	s.entity = nil
	s.isLoggedIn = false
}

func (s *AuthStore) IsAuthenticated() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.token != "" && s.isLoggedIn
}

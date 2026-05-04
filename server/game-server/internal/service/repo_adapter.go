package service

import (
	"context"

	"github.com/cultivation-world/game-server/internal/repository"
)

type sectRepoAdapter struct {
	repo *repository.SectRepository
}

func NewSectRepoAdapter(repo *repository.SectRepository) SectRepository {
	return &sectRepoAdapter{repo: repo}
}

func (a *sectRepoAdapter) Create(ctx context.Context, sectID string, name string, founderID string) error {
	return a.repo.Create(ctx, &repository.Sect{
		ID:        sectID,
		Name:      name,
		FounderID: founderID,
	})
}

func (a *sectRepoAdapter) GetByID(ctx context.Context, id string) (*SectInfo, error) {
	sect, err := a.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sect == nil {
		return nil, nil
	}
	return &SectInfo{
		ID:        sect.ID,
		Name:      sect.Name,
		FounderID: sect.FounderID,
		Alignment: sect.Alignment,
	}, nil
}

func (a *sectRepoAdapter) GetByName(ctx context.Context, name string) (*SectInfo, error) {
	sect, err := a.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if sect == nil {
		return nil, nil
	}
	return &SectInfo{
		ID:        sect.ID,
		Name:      sect.Name,
		FounderID: sect.FounderID,
		Alignment: sect.Alignment,
	}, nil
}

func (a *sectRepoAdapter) AddMember(ctx context.Context, sectID string, entityID string, rank string) error {
	return a.repo.AddMember(ctx, &repository.SectMember{
		SectID:   sectID,
		EntityID: entityID,
		Rank:     rank,
	})
}

func (a *sectRepoAdapter) GetMember(ctx context.Context, sectID string, entityID string) (bool, error) {
	m, err := a.repo.GetMember(ctx, sectID, entityID)
	if err != nil {
		return false, err
	}
	return m != nil, nil
}

func (a *sectRepoAdapter) RemoveMember(ctx context.Context, sectID string, entityID string) error {
	return a.repo.RemoveMember(ctx, sectID, entityID)
}

func (a *sectRepoAdapter) ListMembers(ctx context.Context, sectID string) ([]*SectMemberInfo, error) {
	members, err := a.repo.GetMembers(ctx, sectID)
	if err != nil {
		return nil, err
	}
	result := make([]*SectMemberInfo, 0, len(members))
	for _, m := range members {
		result = append(result, &SectMemberInfo{
			EntityID:     m.EntityID,
			Rank:         m.Rank,
			Contribution: m.Contribution,
			JoinedAt:     m.JoinedAt.Unix(),
		})
	}
	return result, nil
}

type recipeRepoAdapter struct {
	repo *repository.RecipeRepository
}

func NewRecipeRepoAdapter(repo *repository.RecipeRepository) RecipeRepository {
	return &recipeRepoAdapter{repo: repo}
}

func (a *recipeRepoAdapter) GetByID(ctx context.Context, id string) (*RecipeInfo, error) {
	recipe, err := a.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if recipe == nil {
		return nil, nil
	}
	return &RecipeInfo{
		ID:         recipe.ID,
		Type:       recipe.Type,
		Difficulty: recipe.Difficulty,
		Name:       recipe.Name,
	}, nil
}

type friendRepoAdapter struct {
	repo *repository.FriendRepository
}

func NewFriendRepoAdapter(repo *repository.FriendRepository) FriendRepository {
	return &friendRepoAdapter{repo: repo}
}

func (a *friendRepoAdapter) AddFriend(ctx context.Context, entityID, friendID string) error {
	return a.repo.AddFriend(ctx, entityID, friendID)
}

func (a *friendRepoAdapter) RemoveFriend(ctx context.Context, entityID, friendID string) error {
	return a.repo.RemoveFriend(ctx, entityID, friendID)
}

func (a *friendRepoAdapter) AreFriends(ctx context.Context, entityID, friendID string) (bool, error) {
	return a.repo.AreFriends(ctx, entityID, friendID)
}

func (a *friendRepoAdapter) CreateRequest(ctx context.Context, fromID, toID string) (string, error) {
	return a.repo.CreateRequest(ctx, fromID, toID)
}

func (a *friendRepoAdapter) GetPendingRequest(ctx context.Context, fromID, toID string) (*FriendInfo, error) {
	fr, err := a.repo.GetPendingRequest(ctx, fromID, toID)
	if err != nil {
		return nil, err
	}
	if fr == nil {
		return nil, nil
	}
	return &FriendInfo{ID: fr.ID, FromID: fr.FromID, ToID: fr.ToID}, nil
}

func (a *friendRepoAdapter) GetRequestByID(ctx context.Context, requestID string) (*FriendRequestInfo, error) {
	fr, err := a.repo.GetRequestByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if fr == nil {
		return nil, nil
	}
	return &FriendRequestInfo{ID: fr.ID, FromID: fr.FromID, ToID: fr.ToID, Status: fr.Status}, nil
}

func (a *friendRepoAdapter) AcceptRequest(ctx context.Context, requestID string) error {
	return a.repo.AcceptRequest(ctx, requestID)
}

func (a *friendRepoAdapter) GetFriends(ctx context.Context, entityID string) ([]*FriendshipInfo, error) {
	friends, err := a.repo.GetFriends(ctx, entityID)
	if err != nil {
		return nil, err
	}
	result := make([]*FriendshipInfo, 0, len(friends))
	for _, f := range friends {
		result = append(result, &FriendshipInfo{
			FriendID:  f.FriendID,
			CreatedAt: f.CreatedAt.Unix(),
		})
	}
	return result, nil
}

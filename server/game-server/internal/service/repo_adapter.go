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

// Package repo contains Player repository layer
package repo

import (
	"context"

	"github.com/arsnazarenko/devops-basketball/api/gen"
	"github.com/arsnazarenko/devops-basketball/internal/usecase"
	"github.com/arsnazarenko/devops-basketball/pkg/postgres"
)

var _ usecase.PlayerRepo = (*PlayerRepo)(nil)

type PlayerRepo struct {
	pg *postgres.Postgres
}

func NewPlayerRepo(pg *postgres.Postgres) *PlayerRepo {
	return &PlayerRepo{
		pg: pg,
	}
}

// CreatePlayer implements usecase.PlayerRepo.
func (p *PlayerRepo) CreatePlayer(ctx context.Context, request gen.CreatePlayerRequestObject) (gen.CreatePlayerResponseObject, error) {
	return gen.CreatePlayer201JSONResponse{
		Age:         0,
		Citizenship: "",
		Height:      0,
		Id:          0,
		Name:        "",
		Role:        "",
		Surname:     "",
		TeamId:      0,
		Weight:      0,
	}, nil
}

// DeletePlayer implements usecase.PlayerRepo.
func (p *PlayerRepo) DeletePlayer(ctx context.Context, request gen.DeletePlayerRequestObject) (gen.DeletePlayerResponseObject, error) {
	panic("unimplemented")
}

// GetPlayer implements usecase.PlayerRepo.
func (p *PlayerRepo) GetPlayer(ctx context.Context, request gen.GetPlayerRequestObject) (gen.GetPlayerResponseObject, error) {
	panic("unimplemented")
}

// ListPlayers implements usecase.PlayerRepo.
func (p *PlayerRepo) ListPlayers(ctx context.Context, request gen.ListPlayersRequestObject) (gen.ListPlayersResponseObject, error) {
	panic("unimplemented")
}

// UpdatePlayer implements usecase.PlayerRepo.
func (p *PlayerRepo) UpdatePlayer(ctx context.Context, request gen.UpdatePlayerRequestObject) (gen.UpdatePlayerResponseObject, error) {
	panic("unimplemented")
}

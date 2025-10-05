package usecase

import (
	"context"

	"github.com/arsnazarenko/devops-basketball/api/gen"
)

type PlayerUC struct {
	r PlayerRp
}

func NewPlayerUsecase(repo PlayerRp) *PlayerUC {
	return &PlayerUC{
		r: repo,
	}
}

var _ Player = (*PlayerUC)(nil)

// CreatePlayer implements Player.
func (p *PlayerUC) CreatePlayer(ctx context.Context, player *gen.PlayerCreate) (*gen.Player, error) {
	return p.r.CreatePlayer(ctx, player)
}

// DeletePlayer implements Player.
func (p *PlayerUC) DeletePlayer(ctx context.Context, playerID int64) error {
	return p.r.DeletePlayer(ctx, playerID)
}

// GetPlayer implements Player.
func (p *PlayerUC) GetPlayer(ctx context.Context, playerID int64) (*gen.Player, error) {
	return p.r.GetPlayer(ctx, playerID)
}

// GetPlayerList implements Player.
func (p *PlayerUC) GetPlayerList(ctx context.Context, count uint64, offset uint64) ([]gen.Player, error) {
	return p.r.GetPlayerList(ctx, count, offset)
}

// UpdatePlayer implements Player.
func (p *PlayerUC) UpdatePlayer(ctx context.Context, playerID int64, player *gen.PlayerUpdate) (*gen.Player, error) {
	return p.r.UpdatePlayer(ctx, playerID, player)
}

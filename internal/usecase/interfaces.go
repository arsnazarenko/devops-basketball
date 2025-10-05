// Package usecase contains repository layer interfaces
package usecase

import (
	"context"

	"github.com/arsnazarenko/devops-basketball/api/gen"
)

type (
	// Player - use case
	Player interface {
		CreatePlayer(ctx context.Context, player *gen.PlayerCreate) (*gen.Player, error)
		UpdatePlayer(ctx context.Context, playerID int64, player *gen.PlayerUpdate) (*gen.Player, error)
		DeletePlayer(ctx context.Context, playerID int64) error
		GetPlayer(ctx context.Context, playerID int64) (*gen.Player, error)
		GetPlayerList(ctx context.Context, pageSize, pageNumber uint64) ([]gen.Player, error)
	}

	// PlayerRp - mongodb
	PlayerRp interface {
		CreatePlayer(ctx context.Context, player *gen.PlayerCreate) (*gen.Player, error)
		UpdatePlayer(ctx context.Context, playerID int64, player *gen.PlayerUpdate) (*gen.Player, error)
		DeletePlayer(ctx context.Context, playerID int64) error
		GetPlayer(ctx context.Context, playerID int64) (*gen.Player, error)
		GetPlayerList(ctx context.Context, pageSize, pageNumber uint64) ([]gen.Player, error)
	}
)

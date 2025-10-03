// Package usecase contains repository layer interfaces
package usecase

import (
	"context"

	"github.com/arsnazarenko/devops-basketball/api/gen"
)

type PlayerRepo interface {
	UpdatePlayer(ctx context.Context, request gen.UpdatePlayerRequestObject) (gen.UpdatePlayerResponseObject, error)
	CreatePlayer(ctx context.Context, request gen.CreatePlayerRequestObject) (gen.CreatePlayerResponseObject, error)
	DeletePlayer(ctx context.Context, request gen.DeletePlayerRequestObject) (gen.DeletePlayerResponseObject, error)
	GetPlayer(ctx context.Context, request gen.GetPlayerRequestObject) (gen.GetPlayerResponseObject, error)
	ListPlayers(ctx context.Context, request gen.ListPlayersRequestObject) (gen.ListPlayersResponseObject, error)
}

package v1

import (
	"context"
	"errors"

	"github.com/arsnazarenko/devops-basketball/api/gen"
	"github.com/arsnazarenko/devops-basketball/internal/apperrors"
	"github.com/arsnazarenko/devops-basketball/internal/usecase"
)

const (
	defaultPageSize   = 50
	defaultPageNumber = 1
)

var _ gen.StrictServerInterface = (*PlayersServerImpl)(nil)

type PlayersServerImpl struct {
	uc usecase.Player
}

func NewPlayersServerImpl(uc usecase.Player) *PlayersServerImpl {
	return &PlayersServerImpl{
		uc: uc,
	}
}

// CreatePlayer implements gen.StrictServerInterface.
func (p *PlayersServerImpl) CreatePlayer(ctx context.Context, request gen.CreatePlayerRequestObject) (gen.CreatePlayerResponseObject, error) {
	created, err := p.uc.CreatePlayer(ctx, request.Body)
	if errors.Is(err, apperrors.ErrTeamNotFound) {
		return gen.CreatePlayer400Response{}, nil
	}
	if err != nil {
		return nil, err
	}
	return gen.CreatePlayer201JSONResponse(*created), nil
}

// DeletePlayer implements gen.StrictServerInterface.
func (p *PlayersServerImpl) DeletePlayer(ctx context.Context, request gen.DeletePlayerRequestObject) (gen.DeletePlayerResponseObject, error) {
	err := p.uc.DeletePlayer(ctx, request.Id)
	if errors.Is(err, apperrors.ErrPlayerNotFound) {
		return gen.DeletePlayer404Response{}, nil
	}
	if err != nil {
		return nil, err // internal Server Error (500)
	}
	return gen.DeletePlayer204Response{}, nil
}

// GetPlayer implements gen.StrictServerInterface.
func (p *PlayersServerImpl) GetPlayer(ctx context.Context, request gen.GetPlayerRequestObject) (gen.GetPlayerResponseObject, error) {
	player, err := p.uc.GetPlayer(ctx, request.Id)
	if errors.Is(err, apperrors.ErrPlayerNotFound) {
		return gen.GetPlayer404Response{}, nil
	}
	if err != nil {
		return nil, err
	}
	return gen.GetPlayer200JSONResponse(*player), nil
}

// ListPlayers implements gen.StrictServerInterface.
func (p *PlayersServerImpl) ListPlayers(ctx context.Context, request gen.ListPlayersRequestObject) (gen.ListPlayersResponseObject, error) {
	var (
		pageSize   uint64 = defaultPageSize
		pageNumber uint64 = defaultPageNumber
	)

	if request.Params.PageNumber != nil {
		pageNumber = uint64(*request.Params.PageNumber)
	}
	if request.Params.PageSize != nil {
		pageSize = uint64(*request.Params.PageSize)
	}

	list, err := p.uc.GetPlayerList(ctx, pageSize, pageNumber)
	if errors.Is(err, apperrors.ErrInvalidPlayerPageNumber) || errors.Is(err, apperrors.ErrInvalidPlayerPageSize) {
		return gen.ListPlayers400Response{}, nil
	}
	if err != nil {
		return nil, err
	}
	return gen.ListPlayers200JSONResponse(list), nil
}

// UpdatePlayer implements gen.StrictServerInterface.
func (p *PlayersServerImpl) UpdatePlayer(ctx context.Context, request gen.UpdatePlayerRequestObject) (gen.UpdatePlayerResponseObject, error) {
	updated, err := p.uc.UpdatePlayer(ctx, request.Id, request.Body)
	if errors.Is(err, apperrors.ErrPlayerNotFound) {
		return gen.UpdatePlayer404Response{}, nil
	}
	if errors.Is(err, apperrors.ErrTeamNotFound) {
		return gen.UpdatePlayer400Response{}, nil
	}
	if err != nil {
		return nil, err
	}
	return gen.UpdatePlayer200JSONResponse(*updated), nil
}

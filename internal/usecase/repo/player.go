// Package repo contains Player repository layer
package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/arsnazarenko/devops-basketball/api/gen"
	"github.com/arsnazarenko/devops-basketball/internal/apperrors"
	"github.com/arsnazarenko/devops-basketball/internal/usecase"
	"github.com/arsnazarenko/devops-basketball/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

var _ usecase.PlayerRp = (*PlayerRepo)(nil)

type PlayerRepo struct {
	pg *postgres.Postgres
}

func NewPlayerRepo(pg *postgres.Postgres) *PlayerRepo {
	return &PlayerRepo{
		pg: pg,
	}
}

// CreatePlayer implements usecase.PlayerRp.
func (p *PlayerRepo) CreatePlayer(ctx context.Context, player *gen.PlayerCreate) (*gen.Player, error) {
	query := "INSERT INTO players (name, surname, age, height, weight, citizenship, role, team_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	var id int64
	err := p.pg.Pool.QueryRow(ctx, query,
		player.Name,
		player.Surname,
		player.Age,
		player.Height,
		player.Weight,
		player.Citizenship,
		player.Role,
		player.TeamId,
	).Scan(&id)
	// TODO: must check that Team with passed id exists, overwise - return apperrors.ErrTeamNotFound
	if err != nil {
		return nil, fmt.Errorf("repo.CreatePlayer: create player error: %w", err)
	}

	return &gen.Player{
		Id:          id,
		Age:         player.Age,
		Citizenship: player.Citizenship,
		Height:      player.Height,
		Weight:      player.Weight,
		Name:        player.Name,
		Role:        gen.PlayerRole(player.Role),
		Surname:     player.Surname,
		TeamId:      player.TeamId,
	}, nil
}

// DeletePlayer implements usecase.PlayerRp.
func (p *PlayerRepo) DeletePlayer(ctx context.Context, playerID int64) error {
	query := "DELETE FROM players WHERE id = $1"
	res, err := p.pg.Pool.Exec(ctx, query, playerID)
	if err != nil {
		return fmt.Errorf("repo.DeletePlayer: error: %w", err)
	}
	if res.RowsAffected() == 0 {
		return apperrors.ErrPlayerNotFound
	}
	return nil
}

// GetPlayer implements usecase.PlayerRp.
func (p *PlayerRepo) GetPlayer(ctx context.Context, playerID int64) (*gen.Player, error) {
	query := "SELECT id, name, surname, age, height, weight, citizenship, role, team_id FROM players WHERE id = $1"

	var player gen.Player
	if err := p.pg.Pool.QueryRow(ctx, query, playerID).Scan(
		&player.Id,
		&player.Name,
		&player.Surname,
		&player.Age,
		&player.Height,
		&player.Weight,
		&player.Citizenship,
		&player.Role,
		&player.TeamId,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrPlayerNotFound
		}
		return nil, fmt.Errorf("repo.GetPlayer: error: %w", err)
	}
	return &player, nil
}

// GetPlayerList implements usecase.PlayerRp.
func (p *PlayerRepo) GetPlayerList(ctx context.Context, pageSize uint64, pageNumber uint64) ([]gen.Player, error) {
	if pageNumber < 1 {
		return nil, apperrors.ErrInvalidPlayerPageNumber
	}
	if pageSize < 1 {
		return nil, apperrors.ErrInvalidPlayerPageSize
	}
	limit, offset := pageSize, (pageNumber-1)*pageSize
	query := "SELECT id, name, surname, age, height, weight, citizenship, role, team_id FROM players LIMIT $1 OFFSET $2"
	rows, err := p.pg.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrPlayerNotFound
		}
	}
	defer rows.Close()
	list := []gen.Player{}
	for rows.Next() {
		var player gen.Player
		if err := rows.Scan(
			&player.Id,
			&player.Name,
			&player.Surname,
			&player.Age,
			&player.Height,
			&player.Weight,
			&player.Citizenship,
			&player.Role,
			&player.TeamId,
		); err != nil {
			return nil, fmt.Errorf("repo.GetPlayerList: error: %w", err)
		}
		list = append(list, player)
	}
	return list, nil
}

// UpdatePlayer implements usecase.PlayerRp.
func (p *PlayerRepo) UpdatePlayer(ctx context.Context, playerID int64, player *gen.PlayerUpdate) (*gen.Player, error) {
	query := "UPDATE players SET name = $1, surname = $2, age = $3, height = $4, weight = $5, citizenship = $6, role = $7, team_id = $8 WHERE id = $9 RETURNING *"

	var updated gen.Player
	row := p.pg.Pool.QueryRow(ctx, query,
		player.Name,
		player.Surname,
		player.Age,
		player.Height,
		player.Weight,
		player.Citizenship,
		player.Role,
		player.TeamId,
		playerID,
	)
	if err := row.Scan(
		&updated.Id,
		&updated.Name,
		&updated.Surname,
		&updated.Age,
		&updated.Height,
		&updated.Weight,
		&updated.Citizenship,
		&updated.Role,
		&updated.TeamId,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrPlayerNotFound
		}
		// TODO: must check that Team with passed id exists, overwise - return apperrors.ErrTeamNotFound
		return nil, fmt.Errorf("repo.UpdatePlayer: error: %w", err)
	}
	return &updated, nil
}

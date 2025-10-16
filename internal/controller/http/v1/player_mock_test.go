package v1

import (
	"context"

	"github.com/arsnazarenko/devops-basketball/api/gen"
	"github.com/stretchr/testify/mock"
)

// MockPlayer is a mock implementation of usecase.Player interface
type MockPlayer struct {
	mock.Mock
}

func (m *MockPlayer) CreatePlayer(ctx context.Context, player *gen.PlayerCreate) (*gen.Player, error) {
	args := m.Called(ctx, player)
	return args.Get(0).(*gen.Player), args.Error(1)
}

func (m *MockPlayer) UpdatePlayer(ctx context.Context, playerID int64, player *gen.PlayerUpdate) (*gen.Player, error) {
	args := m.Called(ctx, playerID, player)
	return args.Get(0).(*gen.Player), args.Error(1)
}

func (m *MockPlayer) DeletePlayer(ctx context.Context, playerID int64) error {
	args := m.Called(ctx, playerID)
	return args.Error(0)
}

func (m *MockPlayer) GetPlayer(ctx context.Context, playerID int64) (*gen.Player, error) {
	args := m.Called(ctx, playerID)
	return args.Get(0).(*gen.Player), args.Error(1)
}

func (m *MockPlayer) GetPlayerList(ctx context.Context, pageSize, pageNumber uint64) ([]gen.Player, error) {
	args := m.Called(ctx, pageSize, pageNumber)
	return args.Get(0).([]gen.Player), args.Error(1)
}

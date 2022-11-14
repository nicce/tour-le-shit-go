package mock

import (
	"tour-le-shit-go/internal/players/model"
)

type MockedRepository struct {
	Players []model.Player
}

func NewRepository(players []model.Player) *MockedRepository {
	return &MockedRepository{Players: players}
}

func (r *MockedRepository) GetScore(season int) ([]model.Player, error) {
	return r.Players, nil
}

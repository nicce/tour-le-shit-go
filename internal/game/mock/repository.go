package mock

import (
	"tour-le-shit-go/internal/game/model"
)

type MockedRepository struct {
	Players []model.PlayerScore
}

func NewRepository(players []model.PlayerScore) *MockedRepository {
	return &MockedRepository{Players: players}
}

func (r *MockedRepository) GetScore(season int) ([]model.PlayerScore, error) {
	return r.Players, nil
}

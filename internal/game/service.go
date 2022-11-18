package game

import (
	"fmt"
	"tour-le-shit-go/internal/game/model"
)

type Repository interface {
	GetScore(season int) ([]model.PlayerScore, error)
}

type Service interface {
	GetScoreBySeason(season int) ([]model.PlayerScore, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r: r}
}

func (s *service) GetScoreBySeason(season int) ([]model.PlayerScore, error) {
	p, err := s.r.GetScore(season)
	if err != nil {
		return nil, fmt.Errorf("error fetching score from repository %w", err)
	}

	return p, nil
}

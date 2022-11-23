package score

import (
	"fmt"
	"tour-le-shit-go/internal/score/model"
)

type Repository interface {
	GetPlayerScore(id string, season int) ([]model.Score, error)
	DeleteScore(id string) error
	AddScore(score model.ScoreInput) (*model.Score, error)
	GetScoreboard(season int) (model.Scoreboard, error)
}

type Service interface {
	GetPlayerScoreBySeason(id string, season int) ([]model.Score, error)
	DeleteScore(id string) error
	AddScore(score model.ScoreInput) (*model.Score, error)
	GetScoreboard(season int) (model.Scoreboard, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r: r}
}

func (s *service) GetPlayerScoreBySeason(id string, season int) ([]model.Score, error) {
	scores, err := s.r.GetPlayerScore(id, season)
	if err != nil {
		return nil, fmt.Errorf("error fetching player with id %s from repository %w", id, err)
	}

	return scores, nil
}

func (s *service) DeleteScore(id string) error {
	err := s.r.DeleteScore(id)
	if err != nil {
		return fmt.Errorf("error deleting player with id %s %w", id, err)
	}

	return nil
}

func (s *service) AddScore(scoreInput model.ScoreInput) (*model.Score, error) {
	score, err := s.r.AddScore(scoreInput)
	if err != nil {
		return score, fmt.Errorf("error ading scoreInput to player with id %s %w", scoreInput.PlayerId, err)
	}

	return score, nil
}

func (s *service) GetScoreboard(season int) (model.Scoreboard, error) {
	sb, err := s.r.GetScoreboard(season)
	if err != nil {
		return sb, fmt.Errorf("error fetching scoreboard %w", err)
	}

	return sb, nil
}

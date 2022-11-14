package players

import "tour-le-shit-go/internal/players/model"

type Repository interface {
	GetScore(season int) ([]model.Player, error)
}

type Service interface {
	GetScoreBySeason(season int) ([]model.Player, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) GetScoreBySeason(season int) ([]model.Player, error) {
	scoreboard, err := s.r.GetScore(season)
	return scoreboard, err
}

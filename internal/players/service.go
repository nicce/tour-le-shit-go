package players

import (
	"fmt"
	"tour-le-shit-go/internal/players/model"
)

type Repository interface {
	GetPlayerById(id string) (*model.Player, error)
	GetPlayers() ([]model.Player, error)
	CreatePlayer(name string) ([]model.Player, error)
	UpdatePlayer(id, name string) ([]model.Player, error)
	DeletePlayer(id string) ([]model.Player, error)
}

type Service interface {
	GetMembers() ([]model.Player, error)
	CreateMember(name string) ([]model.Player, error)
	UpdateMember(id, name string) ([]model.Player, error)
	DeleteMember(id string) ([]model.Player, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r: r}
}

func (s *service) GetMembers() ([]model.Player, error) {
	p, err := s.r.GetPlayers()
	if err != nil {
		return nil, fmt.Errorf("error fetching players from repository %w", err)
	}

	return p, nil
}

func (s *service) CreateMember(name string) ([]model.Player, error) {
	p, err := s.r.CreatePlayer(name)
	if err != nil {
		return nil, fmt.Errorf("error creating player from repository %w", err)
	}

	return p, nil
}

func (s *service) UpdateMember(id, name string) ([]model.Player, error) {
	p, err := s.r.UpdatePlayer(id, name)
	if err != nil {
		return nil, fmt.Errorf("error updating player from repository %w", err)
	}

	return p, nil
}

func (s *service) DeleteMember(id string) ([]model.Player, error) {
	p, err := s.r.DeletePlayer(id)
	if err != nil {
		return nil, fmt.Errorf("error deleting player from repository %w", err)
	}

	return p, nil
}

package mock

import (
	"fmt"
	"tour-le-shit-go/internal/ierrors"
	"tour-le-shit-go/internal/players/model"

	"github.com/google/uuid"
)

type MockedRepository struct {
	members []model.Player
}

func NewRepository(members []model.Player) *MockedRepository {
	return &MockedRepository{members: members}
}

func (r *MockedRepository) GetPlayerById(id string) (*model.Player, error) {
	for _, m := range r.members {
		if m.Id == id {
			return &model.Player{
				Id:   m.Id,
				Name: m.Name,
			}, nil
		}
	}

	return nil, nil
}

func (r *MockedRepository) GetPlayers() ([]model.Player, error) {
	return r.members, nil
}

func (r *MockedRepository) CreatePlayer(name string) ([]model.Player, error) {
	for _, m := range r.members {
		if name == m.Name {
			return nil, ierrors.HttpError{
				Code:       ierrors.BadRequestStatusCode,
				Message:    fmt.Sprintf("name %s already exists.", name),
				InnerError: "",
			}
		}
	}

	r.members = append(r.members, model.Player{
		Id:   uuid.New().String(),
		Name: name,
	})

	return r.members, nil
}

func (r *MockedRepository) UpdatePlayer(id, name string) ([]model.Player, error) {
	indexToUpdate := -1

	for i, m := range r.members {
		if id == m.Id {
			indexToUpdate = i
		}
	}

	if indexToUpdate < 0 {
		return nil, ierrors.HttpError{
			Code:       ierrors.BadRequestStatusCode,
			Message:    fmt.Sprintf("player with id %s does not exist", id),
			InnerError: "",
		}
	}

	r.members[indexToUpdate] = model.Player{
		Id:   id,
		Name: name,
	}

	return r.members, nil
}

func (r *MockedRepository) DeletePlayer(id string) ([]model.Player, error) {
	updatedMembers := make([]model.Player, 0)

	for _, m := range r.members {
		if m.Id == id {
			continue
		}

		updatedMembers = append(updatedMembers, model.Player{
			Id:   m.Id,
			Name: m.Name,
		})
	}

	r.members = updatedMembers

	return r.members, nil
}

package file

import (
	"encoding/json"
	"io"
	"os"
	"tour-le-shit-go/internal/game/model"
	"tour-le-shit-go/internal/ierrors"
)

type scoreboard struct {
	Season  int      `json:"season"`
	Players []player `json:"game"`
}

type player struct {
	Name       string `json:"name"`
	Points     int    `json:"points"`
	LastPlayed string `json:"lastPlayed"`
}

type JSONFileRepository struct {
	path string
}

func NewRepository(filePath string) *JSONFileRepository {
	return &JSONFileRepository{path: filePath}
}

func (r *JSONFileRepository) GetScore(season int) ([]model.Player, error) {
	file, err := os.Open(r.path)

	if err != nil {
		return nil, ierrors.DbError{Message: "failed to fetch file: " + err.Error()}
	}

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, ierrors.DbError{Message: "failed to read file: " + err.Error()}
	}

	var sbs []scoreboard

	err = json.Unmarshal(byteValue, &sbs)
	if err != nil {
		return nil, ierrors.DbError{Message: "failed to convert file content to struct: " + err.Error()}
	}

	result := make([]model.Player, 0)

	for _, sb := range sbs {
		if sb.Season == season {
			for _, p := range sb.Players {
				result = append(result, model.Player{
					Name:       p.Name,
					Points:     p.Points,
					LastPlayed: p.LastPlayed,
				})
			}
		}
	}

	return result, nil
}

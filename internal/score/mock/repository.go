package mock

import (
	"fmt"
	"strings"
	"tour-le-shit-go/internal/score/model"
	"tour-le-shit-go/internal/utils"

	"github.com/google/uuid"
)

const KeyDelimiter = "_"
const BirdieMultiplier = 2
const EagleMultiplier = 3
const MuliganDiminisher = 3

type MockedRepository struct {
	scores []model.Score
}

func NewRepository(scores []model.Score) *MockedRepository {
	return &MockedRepository{scores: scores}
}

func (r *MockedRepository) GetPlayerScore(id string, season int) ([]model.Score, error) {
	result := make([]model.Score, 0)

	for _, s := range r.scores {
		if s.PlayerId == id && s.Season == season {
			result = append(result, model.Score{
				Id:         s.Id,
				PlayerId:   s.PlayerId,
				PlayerName: s.PlayerName,
				Points:     s.Points,
				Birdies:    s.Birdies,
				Eagles:     s.Eagles,
				Muligans:   s.Muligans,
				Season:     s.Season,
				Day:        s.Day,
			})
		}
	}

	return result, nil
}

func (r *MockedRepository) DeleteScore(id string) error {
	updatedScore := make([]model.Score, 0)

	for _, s := range r.scores {
		if s.Id != id {
			updatedScore = append(updatedScore, model.Score{
				Id:         s.Id,
				PlayerId:   s.PlayerId,
				PlayerName: s.PlayerName,
				Points:     s.Points,
				Birdies:    s.Birdies,
				Eagles:     s.Eagles,
				Muligans:   s.Eagles,
				Season:     s.Season,
				Day:        s.Day,
			})
		}
	}

	r.scores = updatedScore

	return nil
}

func (r *MockedRepository) AddScore(input model.ScoreInput) (*model.Score, error) {
	id := uuid.New().String()
	addedScore := model.Score{
		Id:         id,
		PlayerId:   input.PlayerId,
		PlayerName: "Mock user",
		Points:     input.Points,
		Birdies:    input.Birdies,
		Eagles:     input.Eagles,
		Muligans:   input.Muligans,
		Season:     input.Season,
		Day:        utils.GetToday(),
	}

	r.scores = append(r.scores, addedScore)

	return &addedScore, nil
}

func (r *MockedRepository) GetScoreboard(season int) (model.Scoreboard, error) {
	points := make(map[string]int, 0)
	lastPlayeds := make(map[string]string, 0)

	for _, s := range r.scores {
		if s.Season == season {
			key := fmt.Sprintf("%s%s%s", s.PlayerId, KeyDelimiter, s.PlayerName)
			points[key] += s.Points + BirdieMultiplier*s.Birdies + EagleMultiplier*s.Eagles - MuliganDiminisher*s.Muligans

			if lastPlayeds[key] < s.Day {
				lastPlayeds[key] = s.Day
			}
		}
	}

	sbp := make([]model.ScoreboardPlayer, 0)
	for k, v := range points {
		sbp = append(sbp, model.ScoreboardPlayer{
			Id:         strings.Split(k, KeyDelimiter)[0],
			Name:       strings.Split(k, KeyDelimiter)[1],
			Points:     v,
			LastPlayed: lastPlayeds[k],
		})
	}

	return model.Scoreboard{
		Players: sbp,
		Season:  season,
	}, nil
}

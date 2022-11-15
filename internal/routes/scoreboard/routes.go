package scoreboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"tour-le-shit-go/internal/errors"
	"tour-le-shit-go/internal/players"
)

type Scoreboard struct {
	Season  int      `json:"season"`
	Players []Player `json:"players"`
}

type Player struct {
	Name       string `json:"name"`
	Position   int    `json:"position"`
	Points     int    `json:"points"`
	LastPlayed string `json:"lastPlayed"`
}

func NewScoreboardRoute(s players.Service) Route {
	return Route{s}
}

type Route struct {
	s players.Service
}

func (route *Route) ScoreboardRouteHandler(w http.ResponseWriter, r *http.Request) error {
	season := r.URL.Query().Get("season")
	sint, err := strconv.Atoi(season)
	if err != nil {
		return errors.HttpError{Code: 400, Message: fmt.Sprintf("invalid season query param, expected integer got %s", season)}
	}
	p, err := route.s.GetScoreBySeason(sint)

	if err != nil {
		return errors.HttpError{Code: 500, Message: "server error, please contact support", InnerError: err.Error()}
	}

	sort.Slice(p, func(i, j int) bool {
		return p[i].Points > p[j].Points
	})

	slice := make([]Player, 0)
	for i, player := range p {
		slice = append(slice, Player{
			Name:       player.Name,
			Points:     player.Points,
			Position:   i + 1,
			LastPlayed: player.LastPlayed,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Scoreboard{Season: sint, Players: slice})
	return nil
}

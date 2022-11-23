package scoreboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"tour-le-shit-go/internal/ierrors"
	"tour-le-shit-go/internal/score"
)

type Scoreboard struct {
	Season  int      `json:"season"`
	Players []Player `json:"players"`
}

type Player struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Position   int    `json:"position"`
	Points     int    `json:"points"`
	LastPlayed string `json:"lastPlayed"`
}

func NewScoreboardRoute(s score.Service) Route {
	return Route{s}
}

type Route struct {
	s score.Service
}

func (route *Route) ScoreboardRouteHandler(w http.ResponseWriter, r *http.Request) error {
	season := r.URL.Query().Get("season")

	sint, err := strconv.Atoi(season)
	if err != nil {
		return ierrors.HttpError{Code: ierrors.BadRequestStatusCode, Message: fmt.Sprintf("invalid season query param, expected integer got %s", season)}
	}

	sb, err := route.s.GetScoreboard(sint)

	if err != nil {
		return ierrors.HttpError{Code: ierrors.ServerErrorStatusCode, Message: "server error, please contact support", InnerError: err.Error()}
	}

	sortedPlayerList := sb.Players
	sort.Slice(sortedPlayerList, func(i, j int) bool {
		return sortedPlayerList[i].Points > sortedPlayerList[j].Points
	})

	slice := make([]Player, 0)
	for i, playerScore := range sortedPlayerList {
		slice = append(slice, Player{
			Id:         playerScore.Id,
			Name:       playerScore.Name,
			Points:     playerScore.Points,
			Position:   i + 1,
			LastPlayed: playerScore.LastPlayed,
		})
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(Scoreboard{Season: sb.Season, Players: slice})
	if err != nil {
		return fmt.Errorf("unknown error %w", err)
	}

	return nil
}

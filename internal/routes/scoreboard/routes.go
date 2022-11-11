package scoreboard

import (
	"encoding/json"
	"net/http"
)

type Scoreboard struct {
	Name     string `json:"name"`
	Position int    `json:"position"`
}

func FetchScoreboardRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	scoreboard := make([]Scoreboard, 0)
	scoreboard = append(scoreboard, Scoreboard{
		Name:     "Niclas",
		Position: 1,
	})
	scoreboard = append(scoreboard, Scoreboard{
		Name:     "Andreas",
		Position: 2,
	})
	json.NewEncoder(w).Encode(scoreboard)

}

package server

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"tour-le-shit-go/internal/players"
	"tour-le-shit-go/internal/players/mock"
	"tour-le-shit-go/internal/players/model"
	"tour-le-shit-go/internal/routes/scoreboard"
)

func TestScoreboardRoute(t *testing.T) {
	t.Run("returns 400 due to no query param", func(t *testing.T) {
		// arrange
		srv := BeforeEach([]model.Player{})
		client := httptest.NewServer(srv)
		defer client.Close()

		// act
		res, err := client.Client().Get(client.URL + "/scoreboard")

		// assert
		got := res.StatusCode
		expected := 400
		if err != nil {
			t.Fatal("got error expected none")
		}
		if got != expected {
			t.Errorf("expected %d got %d", expected, got)
		}
	})
	t.Run("returns 400 due to invalid query param", func(t *testing.T) {
		// arrange
		srv := BeforeEach([]model.Player{})
		client := httptest.NewServer(srv)
		defer client.Close()

		// act
		res, err := client.Client().Get(client.URL + "/scoreboard?season=one")

		// assert
		got := res.StatusCode
		expected := 400
		if err != nil {
			t.Fatal("got error expected none")
		}
		if got != expected {
			t.Errorf("expected %d got %d", expected, got)
		}
	})
	t.Run("returns 200", func(t *testing.T) {
		// arrange
		p := make([]model.Player, 0)
		p = append(p, model.Player{
			Name:       "Terminator",
			Points:     444,
			LastPlayed: "2022-11-17",
		})
		p = append(p, model.Player{
			Name:       "Chuck Norris",
			Points:     666,
			LastPlayed: "2022-11-17",
		})
		srv := BeforeEach(p)
		client := httptest.NewServer(srv)
		defer client.Close()

		// act
		res, err := client.Client().Get(client.URL + "/scoreboard?season=1")

		// assert
		got := res.StatusCode
		expected := 200
		if err != nil {
			t.Fatal("got error expected none")
		}
		if got != expected {
			t.Errorf("expected %d got %d", expected, got)
		}
		body, _ := io.ReadAll(res.Body)
		var result scoreboard.Scoreboard
		_ = json.Unmarshal(body, &result)

		firstPosition := result.Players[0].Position
		if result.Players[0].Position != 1 {
			t.Errorf("expected 1 got %d", firstPosition)
		}
	})
}

func BeforeEach(p []model.Player) *Server {
	repository := mock.NewRepository(p)
	service := players.NewService(repository)
	route := scoreboard.NewScoreboardRoute(service)
	cfg := Config{
		ScoreboardRoute: route,
	}
	return New(cfg)
}

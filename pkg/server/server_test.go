package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tour-le-shit-go/internal/game"
	gameMock "tour-le-shit-go/internal/game/mock"
	gameModel "tour-le-shit-go/internal/game/model"
	"tour-le-shit-go/internal/players"
	playersMock "tour-le-shit-go/internal/players/mock"
	playersModel "tour-le-shit-go/internal/players/model"
	"tour-le-shit-go/internal/routes/members"
	"tour-le-shit-go/internal/routes/scoreboard"
	"tour-le-shit-go/pkg/server"
)

const MemberName = "Test"

func TestScoreboardRoute(t *testing.T) {
	t.Parallel()

	beforeEach := func(p []gameModel.PlayerScore) *httptest.Server {
		gameRepository := gameMock.NewRepository(p)
		gameService := game.NewService(gameRepository)
		scoreboardRoute := scoreboard.NewScoreboardRoute(gameService)

		cfg := server.Config{
			ScoreboardRoute: scoreboardRoute,
		}

		return httptest.NewServer(server.New(cfg).Handler)
	}

	t.Run("returns 400 due to no query param", func(t *testing.T) {
		t.Parallel()
		// arrange
		srv := beforeEach([]gameModel.PlayerScore{})
		defer srv.Close()

		// act
		request, _ := http.NewRequestWithContext(context.Background(), "GET", srv.URL+"/scoreboard", strings.NewReader(""))
		res, err := srv.Client().Do(request)

		// assert
		got := res.StatusCode
		expected := 400
		if err != nil {
			t.Fatal("got error expected none")
		}
		if got != expected {
			t.Errorf("expected %d got %d", expected, got)
		}

		_ = res.Body.Close()
	})
	t.Run("returns 400 due to invalid query param", func(t *testing.T) {
		t.Parallel()
		// arrange
		srv := beforeEach([]gameModel.PlayerScore{})
		defer srv.Close()

		// act
		request, _ := http.NewRequestWithContext(context.Background(), "GET", srv.URL+"/scoreboard?season=one", strings.NewReader(""))
		res, err := srv.Client().Do(request)
		// assert
		got := res.StatusCode
		expected := 400
		if err != nil {
			t.Fatal("got error expected none")
		}
		if got != expected {
			t.Errorf("expected %d got %d", expected, got)
		}

		_ = res.Body.Close()
	})
	t.Run("returns 200", func(t *testing.T) {
		t.Parallel()
		// arrange
		p := make([]gameModel.PlayerScore, 0)
		p = append(p, gameModel.PlayerScore{
			Player: playersModel.Player{
				Id:   "1",
				Name: "Terminator",
			},
			Points:     444,
			LastPlayed: "2022-11-17",
		})
		p = append(p, gameModel.PlayerScore{
			Player: playersModel.Player{
				Id:   "1",
				Name: "Chuck Norris",
			},
			Points:     666,
			LastPlayed: "2022-11-17",
		})
		srv := beforeEach(p)
		defer srv.Close()

		// act
		request, _ := http.NewRequestWithContext(context.Background(), "GET", srv.URL+"/scoreboard?season=1", strings.NewReader(""))
		res, err := srv.Client().Do(request)

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

		_ = res.Body.Close()
	})
}

func TestCreateMembersRoute(t *testing.T) {
	t.Parallel()

	beforeEach := func(m []playersModel.Player) *httptest.Server {
		playerRepository := playersMock.NewRepository(m)
		playerService := players.NewService(playerRepository)
		membersRoute := members.NewMemberRoute(playerService)

		cfg := server.Config{
			MembersRoute: membersRoute,
		}

		return httptest.NewServer(server.New(cfg).Handler)
	}

	t.Run("invalid body returns 400", func(t *testing.T) {
		t.Parallel()
		// arrange
		srv := beforeEach([]playersModel.Player{})
		defer srv.Close()

		// act
		request, _ := http.NewRequestWithContext(context.Background(), "PUT", srv.URL+"/members", strings.NewReader("invalid body"))
		res, err := srv.Client().Do(request)

		// assert
		expected := 400
		if err != nil {
			t.Fatal("got error expected none")
		}

		if res.StatusCode != expected {
			t.Errorf("expected %d got %d", expected, res.StatusCode)
		}

		_ = res.Body.Close()
	})

	t.Run("member already exists returns 400", func(t *testing.T) {
		t.Parallel()
		// arrange
		m := make([]playersModel.Player, 0)
		m = append(m, playersModel.Player{
			Id:   "abc-123",
			Name: MemberName,
		})
		srv := beforeEach(m)
		defer srv.Close()

		input := members.MemberInput{
			Name: MemberName,
		}
		b, _ := json.Marshal(input)
		request, _ := http.NewRequestWithContext(context.Background(), "PUT", srv.URL+"/members", bytes.NewReader(b))

		// act
		res, err := srv.Client().Do(request)

		// assert
		expected := 400
		if err != nil {
			t.Fatal("got error expected none")
		}

		if res.StatusCode != expected {
			t.Errorf("expected %d got %d", expected, res.StatusCode)
		}

		_ = res.Body.Close()
	})
	t.Run("returns 200", func(t *testing.T) {
		t.Parallel()
		// arrange
		srv := beforeEach([]playersModel.Player{})
		defer srv.Close()
		input := members.MemberInput{
			Name: "Niclas",
		}
		b, _ := json.Marshal(input)
		request, _ := http.NewRequestWithContext(context.Background(), "PUT", srv.URL+"/members", bytes.NewReader(b))

		// act
		res, err := srv.Client().Do(request)

		// assert
		expected := 200
		if err != nil {
			t.Fatal("got error expected none")
		}

		if res.StatusCode != expected {
			t.Errorf("expected %d got %d", expected, res.StatusCode)
		}

		var result []members.Member
		output, _ := io.ReadAll(res.Body)
		_ = json.Unmarshal(output, &result)

		if len(result) != 1 {
			t.Errorf("expected len %d got %d", 1, len(result))
		}

		_ = res.Body.Close()
	})
}

func TestUpdateMembersRoute(t *testing.T) {
	t.Parallel()

	beforeEach := func(m []playersModel.Player) *httptest.Server {
		playerRepository := playersMock.NewRepository(m)
		playerService := players.NewService(playerRepository)
		membersRoute := members.NewMemberRoute(playerService)

		cfg := server.Config{
			MembersRoute: membersRoute,
		}

		return httptest.NewServer(server.New(cfg).Handler)
	}

	t.Run("member does not exist return 400", func(t *testing.T) {
		t.Parallel()
		// arrange
		srv := beforeEach([]playersModel.Player{})
		defer srv.Close()
		input := members.MemberInput{
			Name: "New name",
		}
		b, _ := json.Marshal(input)
		request, _ := http.NewRequestWithContext(context.Background(), "POST", srv.URL+"/members/new-id", bytes.NewReader(b))

		// arrange
		res, err := srv.Client().Do(request)

		// assert
		if err != nil {
			t.Fatal("got error expected none")
		}

		expected := 400

		if res.StatusCode != expected {
			t.Errorf("expected %d got %d", expected, res.StatusCode)
		}

		_ = res.Body.Close()
	})
	t.Run("invalid body return 400", func(t *testing.T) {
		t.Parallel()
		// arrange
		srv := beforeEach([]playersModel.Player{})
		defer srv.Close()
		request, _ := http.NewRequestWithContext(context.Background(), "POST", srv.URL+"/members/new-id", strings.NewReader("invalid body"))

		// arrange
		res, err := srv.Client().Do(request)

		// assert
		if err != nil {
			t.Fatal("got error expected none")
		}

		expected := 400

		if res.StatusCode != expected {
			t.Errorf("expected %d got %d", expected, res.StatusCode)
		}

		_ = res.Body.Close()
	})
	t.Run("return 200", func(t *testing.T) {
		t.Parallel()
		// arrange

		m := make([]playersModel.Player, 0)
		m = append(m, playersModel.Player{
			Id:   "abc-123",
			Name: MemberName,
		})
		srv := beforeEach(m)
		defer srv.Close()

		input := members.MemberInput{
			Name: "New Test",
		}
		b, _ := json.Marshal(input)
		request, _ := http.NewRequestWithContext(context.Background(), "PUT", srv.URL+"/members", bytes.NewReader(b))

		// act
		res, err := srv.Client().Do(request)

		// assert
		if err != nil {
			t.Fatal("got error expected none")
		}

		expected := 200

		if res.StatusCode != expected {
			t.Errorf("expected %d got %d", expected, res.StatusCode)
		}

		_ = res.Body.Close()
	})
}

func TestDeleteMembersRoute(t *testing.T) {
	t.Parallel()

	beforeEach := func(m []playersModel.Player) *httptest.Server {
		playerRepository := playersMock.NewRepository(m)
		playerService := players.NewService(playerRepository)
		membersRoute := members.NewMemberRoute(playerService)

		cfg := server.Config{
			MembersRoute: membersRoute,
		}

		return httptest.NewServer(server.New(cfg).Handler)
	}

	t.Run("return 200", func(t *testing.T) {
		t.Parallel()

		// arrange
		memberId := "abc-123"
		m := make([]playersModel.Player, 0)

		m = append(m, playersModel.Player{
			Id:   memberId,
			Name: MemberName,
		})

		srv := beforeEach(m)
		defer srv.Close()

		// act
		request, _ := http.NewRequestWithContext(context.Background(), "DELETE", srv.URL+"/members/"+memberId, strings.NewReader(""))
		res, err := srv.Client().Do(request)

		// assert
		if err != nil {
			t.Fatal("got error expected none")
		}

		expected := 200

		if res.StatusCode != expected {
			t.Errorf("expected %d got %d", expected, res.StatusCode)
		}

		var expectedDeletedMembers []members.Member
		b, _ := io.ReadAll(res.Body)
		_ = json.Unmarshal(b, &expectedDeletedMembers)
		if len(expectedDeletedMembers) != 0 {
			t.Errorf("expected len 0 got %d", len(expectedDeletedMembers))
		}

		_ = res.Body.Close()
	})
}

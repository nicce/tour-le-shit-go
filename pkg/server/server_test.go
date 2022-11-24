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
	"tour-le-shit-go/internal/players"
	playersMock "tour-le-shit-go/internal/players/mock"
	playersModel "tour-le-shit-go/internal/players/model"
	"tour-le-shit-go/internal/routes/members"
	"tour-le-shit-go/internal/routes/scoreboard"
	"tour-le-shit-go/internal/score"
	scoreMock "tour-le-shit-go/internal/score/mock"
	scoreModel "tour-le-shit-go/internal/score/model"
	"tour-le-shit-go/pkg/server"
)

const MemberName = "Test"

func TestScoreboardRoute(t *testing.T) {
	t.Parallel()

	beforeEach := func(s []scoreModel.Score) *httptest.Server {
		scoreRepository := scoreMock.NewRepository(s)
		scoreService := score.NewService(scoreRepository)
		scoreboardRoute := scoreboard.NewScoreboardRoute(scoreService)

		cfg := server.Config{
			ScoreboardRoute: scoreboardRoute,
		}

		return httptest.NewServer(server.New(cfg).Handler)
	}

	t.Run("returns 200 with correct position based on points", func(t *testing.T) {
		t.Parallel()

		// arrange
		scores := make([]scoreModel.Score, 0)

		scores = append(scores, scoreModel.Score{
			Id:         "id1",
			PlayerId:   "Player1",
			PlayerName: "Player1",
			Points:     30,
			Birdies:    0,
			Eagles:     0,
			Muligans:   0,
			Season:     1,
			Day:        "2022-01-01",
		})
		scores = append(scores, scoreModel.Score{
			Id:         "id2",
			PlayerId:   "Player2",
			PlayerName: "Player2",
			Points:     31,
			Birdies:    0,
			Eagles:     0,
			Muligans:   0,
			Season:     1,
			Day:        "2022-01-01",
		})

		srv := beforeEach(scores)
		defer srv.Close()

		// act
		request, _ := http.NewRequestWithContext(context.Background(), "GET", srv.URL+"/scoreboard?season=1", strings.NewReader(""))
		res, err := srv.Client().Do(request)

		// assert
		if err != nil {
			t.Fatalf("got error: %v expected none", err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("got error: %v expected none", err)
		}

		var scoreboardResponse scoreboard.Scoreboard

		err = json.Unmarshal(body, &scoreboardResponse)
		if err != nil {
			t.Fatalf("got error: %v expected none", err)
		}

		expectedPlayerLength := 2
		if len(scoreboardResponse.Players) != expectedPlayerLength {
			t.Fatalf("expected %d got %d", expectedPlayerLength, len(scoreboardResponse.Players))
		}

		expectedPlayerIdAsFirst := "Player2"
		if scoreboardResponse.Players[0].Id != expectedPlayerIdAsFirst {
			t.Errorf("expected %s got %s", expectedPlayerIdAsFirst, scoreboardResponse.Players[0].Id)
		}

		_ = res.Body.Close()
	})
	t.Run("returns 200 with empty array due to different season", func(t *testing.T) {
		t.Parallel()

		// arrange
		scores := make([]scoreModel.Score, 0)

		scores = append(scores, scoreModel.Score{
			Id:         "id1",
			PlayerId:   "Player1",
			PlayerName: "Player1",
			Points:     30,
			Birdies:    0,
			Eagles:     0,
			Muligans:   0,
			Season:     1,
			Day:        "2022-01-01",
		})
		scores = append(scores, scoreModel.Score{
			Id:         "id2",
			PlayerId:   "Player2",
			PlayerName: "Player2",
			Points:     31,
			Birdies:    0,
			Eagles:     0,
			Muligans:   0,
			Season:     1,
			Day:        "2022-01-01",
		})

		srv := beforeEach(scores)
		defer srv.Close()

		// act
		request, _ := http.NewRequestWithContext(context.Background(), "GET", srv.URL+"/scoreboard?season=2", strings.NewReader(""))
		res, err := srv.Client().Do(request)

		// assert
		if err != nil {
			t.Fatalf("got error: %v expected none", err)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("got error: %v expected none", err)
		}

		var scoreboardResponse scoreboard.Scoreboard

		err = json.Unmarshal(body, &scoreboardResponse)
		if err != nil {
			t.Fatalf("got error: %v expected none", err)
		}

		expectedSeason := 2
		if scoreboardResponse.Season != expectedSeason {
			t.Errorf("expected %d got %d", expectedSeason, scoreboardResponse.Season)
		}

		expectedPlayerLength := 0
		if len(scoreboardResponse.Players) != expectedPlayerLength {
			t.Errorf("expected %d got %d", expectedPlayerLength, len(scoreboardResponse.Players))
		}

		_ = res.Body.Close()
	})
	t.Run("returns 400 due to no season query param", func(t *testing.T) {
		t.Parallel()

		// arrange
		srv := beforeEach(make([]scoreModel.Score, 0))
		defer srv.Close()

		// act
		request, _ := http.NewRequestWithContext(context.Background(), "GET", srv.URL+"/scoreboard", strings.NewReader(""))
		res, err := srv.Client().Do(request)

		// assert
		if err != nil {
			t.Fatalf("got error: %v expected none", err)
		}

		expectedStatusCode := 400
		if res.StatusCode != expectedStatusCode {
			t.Errorf("expected %d got %d", expectedStatusCode, res.StatusCode)
		}

		_ = res.Body.Close()
	})
	t.Run("returns 400 due to non integer season query param", func(t *testing.T) {
		t.Parallel()

		// arrange
		srv := beforeEach(make([]scoreModel.Score, 0))
		defer srv.Close()

		// act
		request, _ := http.NewRequestWithContext(context.Background(), "GET", srv.URL+"/scoreboard?season=thisiswrong", strings.NewReader(""))
		res, err := srv.Client().Do(request)

		// assert
		if err != nil {
			t.Fatalf("got error: %v expected none", err)
		}

		expectedStatusCode := 400
		if res.StatusCode != expectedStatusCode {
			t.Errorf("expected %d got %d", expectedStatusCode, res.StatusCode)
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

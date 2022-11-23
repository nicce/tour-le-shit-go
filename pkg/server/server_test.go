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
	"tour-le-shit-go/pkg/server"
)

const MemberName = "Test"

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

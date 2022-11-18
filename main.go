package main

import (
	"database/sql"
	"fmt"
	"log"
	"tour-le-shit-go/internal/env"
	"tour-le-shit-go/internal/game"
	gameDb "tour-le-shit-go/internal/game/db"
	gameMock "tour-le-shit-go/internal/game/mock"
	gameModel "tour-le-shit-go/internal/game/model"
	"tour-le-shit-go/internal/players"
	playersMock "tour-le-shit-go/internal/players/mock"
	playersModel "tour-le-shit-go/internal/players/model"
	"tour-le-shit-go/internal/routes/members"
	"tour-le-shit-go/internal/routes/scoreboard"
	"tour-le-shit-go/pkg/server"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env", ".env.default")
	if err != nil {
		log.Fatal(err.Error())
	}

	appEnv := env.GetAppEnv()

	scoreboardRoute := createScoreboardRoute(appEnv)
	membersRoute := createMembersRoute(appEnv)

	config := server.Config{
		ScoreboardRoute: scoreboardRoute,
		Port:            appEnv.Port,
		MembersRoute:    membersRoute,
	}

	srv := server.New(config)

	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func createScoreboardRoute(appEnv env.AppEnv) scoreboard.Route {
	var repository game.Repository

	switch appEnv.ScoreboardMode {
	case "PSQL":
		database, err := sql.Open("postgres", "user=user dbname=tourleshit password=password sslmode=disable")
		if err != nil {
			panic(err)
		}

		repository = gameDb.NewRepository(database)
	case "MOCK":
		repository = gameMock.NewRepository([]gameModel.Player{})
	default:
		panic(fmt.Sprintf("invalid scoreboard mode %s", appEnv.ScoreboardMode))
	}

	service := game.NewService(repository)

	return scoreboard.NewScoreboardRoute(service)
}

func createMembersRoute(appEnv env.AppEnv) members.Route {
	var repository players.Repository

	switch appEnv.MembersMode {
	case "PSQL":
		panic("not yet implemented")
	case "MOCK":
		repository = playersMock.NewRepository([]playersModel.Player{})
	default:
		panic(fmt.Sprintf("invalid members mode %s", appEnv.MembersMode))
	}

	service := players.NewService(repository)

	return members.NewMemberRoute(service)
}

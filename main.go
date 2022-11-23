package main

import (
	"database/sql"
	"fmt"
	"log"
	"tour-le-shit-go/internal/env"
	"tour-le-shit-go/internal/players"
	playersDb "tour-le-shit-go/internal/players/db"
	playersMock "tour-le-shit-go/internal/players/mock"
	playersModel "tour-le-shit-go/internal/players/model"
	"tour-le-shit-go/internal/routes/members"
	"tour-le-shit-go/internal/routes/scoreboard"
	"tour-le-shit-go/internal/routes/scores"
	"tour-le-shit-go/internal/score"
	scoreDb "tour-le-shit-go/internal/score/db"
	scoreMock "tour-le-shit-go/internal/score/mock"
	scoreModel "tour-le-shit-go/internal/score/model"
	"tour-le-shit-go/pkg/server"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const MockMode = "MOCK"
const PsqlMode = "PSQL"

func main() {
	err := godotenv.Load(".env", ".env.default")
	if err != nil {
		log.Fatal(err.Error())
	}

	appEnv := env.GetAppEnv()

	var playersRepository players.Repository

	switch appEnv.MembersMode {
	case PsqlMode:
		database, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", appEnv.Db.Username, appEnv.Db.Name, appEnv.Db.Password))
		if err != nil {
			panic(err)
		}

		playersRepository = playersDb.NewRepository(database)
	case MockMode:
		playersRepository = playersMock.NewRepository([]playersModel.Player{})
	default:
		panic(fmt.Sprintf("invalid members mode %s", appEnv.MembersMode))
	}

	var scoreRepository score.Repository

	switch appEnv.ScoreMode {
	case PsqlMode:
		database, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", appEnv.Db.Username, appEnv.Db.Name, appEnv.Db.Password))
		if err != nil {
			panic(err)
		}

		scoreRepository = scoreDb.NewRepository(database, playersRepository)
	case MockMode:
		scoreRepository = scoreMock.NewRepository([]scoreModel.Score{})
	}

	scoreService := score.NewService(scoreRepository)

	playersService := players.NewService(playersRepository)

	config := server.Config{
		ScoresRoute:     scores.NewScoresRoute(scoreService),
		ScoreboardRoute: scoreboard.NewScoreboardRoute(scoreService),
		Port:            appEnv.Port,
		MembersRoute:    members.NewMemberRoute(playersService),
	}

	srv := server.New(config)

	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

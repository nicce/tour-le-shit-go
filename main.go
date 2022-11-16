package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"tour-le-shit-go/internal/env"
	"tour-le-shit-go/internal/players"
	"tour-le-shit-go/internal/players/file"
	"tour-le-shit-go/internal/routes/scoreboard"
	"tour-le-shit-go/pkg/server"
)

func main() {
	err := godotenv.Load(".env", ".env.default")
	if err != nil {
		log.Fatal(err.Error())
	}

	appEnv := env.GetAppEnv()

	var repository players.Repository
	if appEnv.ScoreboardMode == "FILE" {
		repository = file.NewRepository("players.json")
	} else if appEnv.ScoreboardMode == "PSQL" {
		repository = file.NewRepository("players.json")
	}

	service := players.NewService(repository)
	scoreboardRoute := scoreboard.NewScoreboardRoute(service)

	config := server.Config{
		ScoreboardRoute: scoreboardRoute,
	}
	srv := server.New(config)
	log.Fatal(http.ListenAndServe(":"+appEnv.Port, srv))
}

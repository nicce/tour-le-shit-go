package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"tour-le-shit-go/internal/env"
	"tour-le-shit-go/internal/players"
	"tour-le-shit-go/internal/players/file"
	"tour-le-shit-go/internal/routes/scoreboard"
	"tour-le-shit-go/pkg/server"
)

const port = ":4000"

func main() {
	err := godotenv.Load(".env", ".env.default")
	if err != nil {
		log.Fatal(err.Error())
	}

	appEnv := env.GetAppEnv()
	fmt.Print(appEnv.ScoreboardMode)

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
	log.Fatal(http.ListenAndServe(port, srv))
}

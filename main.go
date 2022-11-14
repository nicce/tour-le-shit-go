package main

import (
	"log"
	"net/http"
	"tour-le-shit-go/internal/players"
	"tour-le-shit-go/internal/players/file"
	"tour-le-shit-go/internal/routes/scoreboard"
	"tour-le-shit-go/pkg/server"
)

const port = ":4000"

func main() {
	repository := file.NewRepository("players.json")
	service := players.NewService(repository)
	scoreboardRoute := scoreboard.NewScoreboardRoute(service)

	config := server.Config{
		ScoreboardRoute: scoreboardRoute,
	}
	srv := server.New(config)
	log.Fatal(http.ListenAndServe(port, srv))
}

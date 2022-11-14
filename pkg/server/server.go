package server

import (
	"net/http"
	"tour-le-shit-go/internal/routes/scoreboard"
)

type Server struct {
	http.Handler
}

type Config struct {
	ScoreboardRoute scoreboard.Route
}

// New creates a server with routes configured
func New(cfg Config) *Server {
	s := new(Server)
	router := http.NewServeMux()

	router.HandleFunc("/scoreboard", cfg.ScoreboardRoute.ScoreboardRouteHandler)

	s.Handler = router
	return s
}

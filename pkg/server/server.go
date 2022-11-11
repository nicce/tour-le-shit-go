package server

import (
	"net/http"
	"tour-le-shit-go/internal/routes/scoreboard"
)

type Server struct {
	http.Handler
}

// New creates a server with routes configured
func New() *Server {
	s := new(Server)
	router := http.NewServeMux()

	createScoreboardRoutes("/scoreboard", router)
	s.Handler = router

	return s
}

func createScoreboardRoutes(route string, router *http.ServeMux) {
	router.HandleFunc(route, scoreboard.FetchScoreboardRoute)

}

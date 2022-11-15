package server

import (
	"log"
	"net/http"
	"tour-le-shit-go/internal/errors"
	"tour-le-shit-go/internal/logger"
	"tour-le-shit-go/internal/routes/scoreboard"
)

type Server struct {
	http.Handler
}

type Config struct {
	ScoreboardRoute scoreboard.Route
}

type rootHandler func(http.ResponseWriter, *http.Request) error

// New creates a server with routes configured
func New(cfg Config) *Server {
	s := new(Server)
	router := http.NewServeMux()

	router.Handle("/scoreboard", rootHandler(cfg.ScoreboardRoute.ScoreboardRouteHandler))

	s.Handler = logger.RequestLogger(router)
	return s
}

func (fn rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err == nil {
		return
	}
	log.Printf("an error occured %v", err)
	w.Header().Set("content-type", "application/problem+json")
	httpError, ok := err.(errors.HttpError)
	if !ok {
		w.WriteHeader(500)
		return

	}
	w.WriteHeader(httpError.Code)
	w.Write(httpError.PrintBody())
	return
}

package server

import (
	"errors"
	"log"
	"net/http"
	"time"
	"tour-le-shit-go/internal/ierrors"
	"tour-le-shit-go/internal/logger"
	"tour-le-shit-go/internal/routes/members"
	"tour-le-shit-go/internal/routes/scoreboard"

	"github.com/gorilla/mux"
)

type Server struct {
	http.Handler
}

type Config struct {
	MembersRoute    members.Route
	Port            string
	ScoreboardRoute scoreboard.Route
}

type rootHandler func(http.ResponseWriter, *http.Request) error

const Timeout = 5 * time.Second

// New creates a http.Server with routes configured.
func New(cfg Config) *http.Server {
	s := new(Server)
	router := mux.NewRouter()

	router.Handle("/scoreboard", rootHandler(cfg.ScoreboardRoute.ScoreboardRouteHandler))
	router.Handle("/members/{id}", rootHandler(cfg.MembersRoute.MemberRouteHandler))
	router.Handle("/members", rootHandler(cfg.MembersRoute.MembersRouteHandler))

	s.Handler = logger.RequestLogger(router)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           s.Handler,
		ReadHeaderTimeout: Timeout,
	}

	return srv
}

func (fn rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err == nil {
		return
	}

	log.Printf("an error occurred %v", err)
	w.Header().Set("content-type", "application/problem+json")

	var httpError ierrors.HttpError
	ok := errors.As(err, &httpError)

	if !ok {
		w.WriteHeader(ierrors.ServerErrorStatusCode)

		return
	}

	w.WriteHeader(httpError.Code)

	_, err = w.Write(httpError.PrintBody())
	if err != nil {
		panic(err)
	}
}

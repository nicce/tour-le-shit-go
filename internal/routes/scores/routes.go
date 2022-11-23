package scores

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"tour-le-shit-go/internal/ierrors"
	"tour-le-shit-go/internal/score"
	"tour-le-shit-go/internal/score/model"

	"github.com/gorilla/mux"
)

type Response struct {
	Player PlayerScoreResponse `json:"player"`
	Season int                 `json:"season"`
	Scores []ScoreResponse     `json:"scores"`
}

type PlayerScoreResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ScoreResponse struct {
	Id       string `json:"id"`
	Points   int    `json:"points"`
	Birdies  int    `json:"birdies"`
	Eagles   int    `json:"eagles"`
	Muligans int    `json:"muligans"`
	Day      string `json:"day"`
}

type ScoreRequest struct {
	PlayerId string `json:"playerId"`
	Points   int    `json:"points"`
	Birdies  int    `json:"birdies"`
	Eagles   int    `json:"eagles"`
	Muligans int    `json:"muligans"`
	Season   int    `json:"season"`
}

const ContentTypeKey = "Content-Type"
const ContentTypeValue = "application/json"
const CreatedStatusCode = 201
const NoContentStatusCode = 204

type Route struct {
	s score.Service
}

func NewScoresRoute(s score.Service) Route {
	return Route{s: s}
}

func (r *Route) ScoresRouteHandler(w http.ResponseWriter, req *http.Request) error {
	switch req.Method {
	case "GET":
		return r.handleGetRequest(w, req)
	case "PUT":
		return r.handlePutRequest(w, req)
	}

	return ierrors.HttpError{
		Code:       ierrors.BadRequestStatusCode,
		Message:    "Unsupported method type",
		InnerError: "",
	}
}

func (r *Route) ScoreRouteHandler(w http.ResponseWriter, req *http.Request) error {
	if req.Method == "DELETE" {
		return r.handleDeleteRequest(w, req)
	}

	return ierrors.HttpError{
		Code:       ierrors.BadRequestStatusCode,
		Message:    "Unsupported method type",
		InnerError: "",
	}
}

func (r *Route) handleGetRequest(w http.ResponseWriter, req *http.Request) error {
	season := req.URL.Query().Get("season")
	playerId := req.URL.Query().Get("playerId")

	sint, err := strconv.Atoi(season)
	if err != nil {
		return ierrors.HttpError{Code: ierrors.BadRequestStatusCode, Message: fmt.Sprintf("invalid season query param, expected integer got %s", season)}
	}

	scores, err := r.s.GetPlayerScoreBySeason(playerId, sint)
	if err != nil {
		return fmt.Errorf("error fetching player score: %w", err)
	}

	if len(scores) == 0 {
		w.Header().Set(ContentTypeKey, ContentTypeValue)
		_ = json.NewEncoder(w).Encode([]Response{})

		return nil
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Day > scores[j].Day
	})

	result := make([]ScoreResponse, 0)

	for _, s := range scores {
		result = append(result, ScoreResponse{
			Id:       s.Id,
			Points:   s.Points,
			Birdies:  s.Birdies,
			Eagles:   s.Eagles,
			Muligans: s.Muligans,
			Day:      s.Day,
		})
	}

	response := Response{
		Player: PlayerScoreResponse{
			Id:   scores[0].PlayerId,
			Name: scores[0].PlayerName,
		},
		Season: sint,
		Scores: result,
	}

	w.Header().Set(ContentTypeKey, ContentTypeValue)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		return fmt.Errorf("unknown error %w", err)
	}

	return nil
}

func (r *Route) handlePutRequest(w http.ResponseWriter, req *http.Request) error {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return ierrors.HttpError{
			Code:       ierrors.BadRequestStatusCode,
			Message:    "invalid request body",
			InnerError: err.Error(),
		}
	}

	var scoreRequest ScoreRequest

	err = json.Unmarshal(b, &scoreRequest)
	if err != nil {
		return ierrors.HttpError{
			Code:       ierrors.BadRequestStatusCode,
			Message:    "invalid request body",
			InnerError: err.Error(),
		}
	}

	_, err = r.s.AddScore(model.ScoreInput{
		PlayerId: scoreRequest.PlayerId,
		Points:   scoreRequest.Points,
		Birdies:  scoreRequest.Birdies,
		Eagles:   scoreRequest.Eagles,
		Muligans: scoreRequest.Muligans,
		Season:   scoreRequest.Season,
	})

	if err != nil {
		return fmt.Errorf("error adding score: %w", err)
	}

	w.WriteHeader(CreatedStatusCode)

	return nil
}

func (r *Route) handleDeleteRequest(w http.ResponseWriter, req *http.Request) error {
	id := mux.Vars(req)["id"]
	err := r.s.DeleteScore(id)

	if err != nil {
		return fmt.Errorf("error deleting score %w", err)
	}

	w.WriteHeader(NoContentStatusCode)

	return nil
}

package members

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"tour-le-shit-go/internal/ierrors"
	"tour-le-shit-go/internal/players"

	"github.com/gorilla/mux"
)

// Member the tour le shit tour.
type Member struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type MemberInput struct {
	Name string `json:"name"`
}

type Route struct {
	s players.Service
}

func NewMemberRoute(s players.Service) Route {
	return Route{s: s}
}

const ContentTypeKey = "Content-Type"
const ContentTypeValue = "application/json"

func (r Route) MembersRouteHandler(w http.ResponseWriter, req *http.Request) error {
	switch req.Method {
	case "GET":
		return r.handleGetRequest(w)
	case "PUT":
		return r.handlePutReqeuest(w, req)
	}

	return ierrors.HttpError{
		Code:       ierrors.BadRequestStatusCode,
		Message:    "Unsupported method type",
		InnerError: "",
	}
}

func (r Route) MemberRouteHandler(w http.ResponseWriter, req *http.Request) error {
	switch req.Method {
	case "POST":
		return r.handlePostRequest(w, req)
	case "DELETE":
		return r.handleDeleteRequest(w, req)
	}

	return ierrors.HttpError{
		Code:       ierrors.BadRequestStatusCode,
		Message:    "Unsupported method type",
		InnerError: "",
	}
}

func (r Route) handleGetRequest(w http.ResponseWriter) error {
	members, err := r.s.GetMembers()
	if err != nil {
		return fmt.Errorf("error fetching members %w", err)
	}

	result := make([]Member, 0)
	for _, m := range members {
		result = append(result, Member{
			Id:   m.Id,
			Name: m.Name,
		})
	}

	w.Header().Set(ContentTypeKey, ContentTypeValue)

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		return fmt.Errorf("unknown error %w", err)
	}

	return nil
}

func (r Route) handlePutReqeuest(w http.ResponseWriter, req *http.Request) error {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return ierrors.HttpError{
			Code:       ierrors.BadRequestStatusCode,
			Message:    "invalid body",
			InnerError: err.Error(),
		}
	}

	var cm MemberInput

	err = json.Unmarshal(b, &cm)
	if err != nil {
		return ierrors.HttpError{
			Code:       ierrors.BadRequestStatusCode,
			Message:    "invalid request body",
			InnerError: err.Error(),
		}
	}

	members, err := r.s.CreateMember(cm.Name)
	if err != nil {
		return fmt.Errorf("error creating member %w", err)
	}

	w.Header().Set(ContentTypeKey, ContentTypeValue)

	err = json.NewEncoder(w).Encode(members)
	if err != nil {
		return fmt.Errorf("unknown error %w", err)
	}

	return nil
}

func (r Route) handlePostRequest(w http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)

	b, err := io.ReadAll(req.Body)
	if err != nil {
		return ierrors.HttpError{
			Code:       ierrors.BadRequestStatusCode,
			Message:    "invalid body",
			InnerError: err.Error(),
		}
	}

	var m MemberInput

	err = json.Unmarshal(b, &m)
	if err != nil {
		return ierrors.HttpError{
			Code:       ierrors.BadRequestStatusCode,
			Message:    "invalid body",
			InnerError: err.Error(),
		}
	}

	members, err := r.s.UpdateMember(vars["id"], m.Name)
	if err != nil {
		return fmt.Errorf("error updating member %w", err)
	}

	w.Header().Set(ContentTypeKey, ContentTypeValue)

	err = json.NewEncoder(w).Encode(members)
	if err != nil {
		return fmt.Errorf("unknown error %w", err)
	}

	return nil
}

func (r Route) handleDeleteRequest(w http.ResponseWriter, req *http.Request) error {
	id := mux.Vars(req)["id"]
	members, err := r.s.DeleteMember(id)

	if err != nil {
		return fmt.Errorf("error deleting member %w", err)
	}

	w.Header().Set(ContentTypeKey, ContentTypeValue)

	err = json.NewEncoder(w).Encode(members)
	if err != nil {
		return fmt.Errorf("unknown error %w", err)
	}

	return nil
}

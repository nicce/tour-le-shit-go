package errors

import (
	"encoding/json"
	"fmt"
)

type HttpError struct {
	Code       int    `json:"status"`
	Message    string `json:"detail"`
	InnerError string `json:"-"`
}

func (h HttpError) Error() string {
	return fmt.Sprintf("message: %s innerError: %s", h.Message, h.InnerError)
}

func (h HttpError) PrintBody() []byte {
	b, err := json.Marshal(h)

	if err != nil {
		return []byte("unknown error")
	}
	return b
}

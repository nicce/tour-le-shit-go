package ierrors

import (
	"encoding/json"
	"fmt"
)

const BadRequestStatusCode = 400
const ServerErrorStatusCode = 500

type HttpError struct {
	Code       int    `json:"status"`
	Message    string `json:"detail"`
	InnerError string `json:"-"`
}

type DbError struct {
	Message string
}

func (h HttpError) Error() string {
	return fmt.Sprintf("message: %s innerError: %s", h.Message, h.InnerError)
}

func (d DbError) Error() string {
	return d.Message
}

func (h HttpError) PrintBody() []byte {
	b, err := json.Marshal(h)

	if err != nil {
		return []byte("unknown error")
	}

	return b
}

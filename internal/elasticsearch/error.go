package elasticsearch

import (
	"fmt"
	"net/http"
)

type Error struct {
	Status      int    `json:"-"`
	ErrorType   string `json:"error.type"`
	Reason      string `json:"error.reason"`
	Index       string `json:"error.index,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s: %s", e.Status, e.ErrorType, e.Reason)
}

func NewError(status int, errType, reason string) *Error {
	return &Error{Status: status, ErrorType: errType, Reason: reason}
}

func NewNotFound(index string) *Error {
	return &Error{
		Status:    http.StatusNotFound,
		ErrorType: "index_not_found_exception",
		Reason:    fmt.Sprintf("no such index [%s]", index),
		Index:     index,
	}
}

func NewBadRequest(reason string) *Error {
	return &Error{
		Status:    http.StatusBadRequest,
		ErrorType: "action_request_validation_exception",
		Reason:    reason,
	}
}

func NewNotImplemented() *Error {
	return &Error{
		Status:    http.StatusNotImplemented,
		ErrorType: "not_implemented_exception",
		Reason:    "This API is not implemented",
	}
}

type ErrorResponse struct {
	Error  ErrorBody `json:"error"`
	Status int       `json:"status"`
}

type ErrorBody struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
	Index  string `json:"index,omitempty"`
}

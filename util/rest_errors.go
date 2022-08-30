package util

import (
	"fmt"
	"net/http"
)

type RestErr interface {
	Message() string
	Status() int
	Error() string
	Causes() []interface{}
}

type restErr struct {
	Messages string        `json:"message"`
	Statuss  int           `json:"status"`
	Errors   string        `json:"error"`
	Cause    []interface{} `json:"causes"`
}

func (e restErr) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s - causes: [%v]",
		e.Messages, e.Statuss, e.Errors, e.Cause)
}

func (e restErr) Message() string {
	return e.Messages
}
func (e restErr) Status() int {
	return e.Statuss
}

func (e restErr) Causes() []interface{} {
	return e.Cause
}

func NewRestError(message string, status int, err string, causes []interface{}) RestErr {
	return restErr{
		Messages: message,
		Statuss:  status,
		Errors:   err,
		Cause:    causes,
	}
}

func NewBadRequestError(message string) RestErr {
	return &restErr{
		Messages: message,
		Statuss:  http.StatusBadRequest,
		Errors:   "bad_request",
	}
}

func NewInternalServerError(message string, err error) RestErr {
	result := &restErr{
		Messages: message,
		Statuss:  http.StatusInternalServerError,
		Errors:   "internal_server_error",
	}
	if err != nil {
		result.Cause = append(result.Cause, err.Error())
	}
	return result
}
func NewNotFoundError(message string) RestErr {
	return &restErr{
		Messages: message,
		Statuss:  http.StatusNotFound,
		Errors:   "not_found",
	}
}

func NewUnauthorizedError(message string) RestErr {
	return &restErr{
		Messages: message,
		Statuss:  http.StatusUnauthorized,
		Errors:   "unauthorized",
	}
}

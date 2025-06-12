package services

import "net/http"

type APIError struct {
	Code    int    // HTTP status code
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}

func NewAPIError(code int, msg string) *APIError {
	return &APIError{
		Code:    code,
		Message: msg,
	}
}

func NewBadRequestError(msg string) *APIError {
	return NewAPIError(http.StatusBadRequest, msg)
}

func NewInternalServerError(msg string) *APIError {
	return NewAPIError(http.StatusInternalServerError, msg)
}

func NewNotFoundError(msg string) *APIError {
	return NewAPIError(http.StatusNotFound, msg)
}

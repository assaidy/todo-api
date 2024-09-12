package utils

import (
	"fmt"
	"net/http"
)

type ApiError struct {
	StatusCode int `json:"statusCode"`
	Msg        any `json:"msg"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("api error: %d - %v", e.StatusCode, e.Msg)
}

func NewApiError(statusCode int, msg any) ApiError {
	return ApiError{
		StatusCode: statusCode,
		Msg:        msg,
	}
}

func InvalidJSONError() ApiError {
	return NewApiError(http.StatusBadRequest, "invalid JSON request data")
}

func InvalidRequestData(msg any) ApiError {
	return NewApiError(http.StatusBadRequest, msg)
}

func NotFoundError(msg string) ApiError {
	return NewApiError(http.StatusNotFound, msg)
}

func AlreadyExistsError(msg string) ApiError {
	return NewApiError(http.StatusBadRequest, msg)
}

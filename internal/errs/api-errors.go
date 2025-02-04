package errs

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int `json:"statusCode"`
	Msg        any `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("%v", e.Msg)
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
}

func InvalidJSON() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid JSON request"))
}

func InvalidID() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid song ID"))
}

func InvalidPage() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid page"))
}

func InvalidPageSize() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid page size"))
}

func InvalidDate() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid date format"))
}

func EndOfText() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("end of song text"))
}

func NoText() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("song does not have text yet"))
}

func NoSongs() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("songs not found"))
}

func APICallTimeout() APIError {
	return NewAPIError(http.StatusRequestTimeout, fmt.Errorf("request timeout"))
}

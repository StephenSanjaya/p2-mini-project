package httputil

import "fmt"

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (h *HTTPError) Error() string {
	return fmt.Sprintf("Error %d: %s", h.Code, h.Message)
}

func NewError(status int, message string, err error) *HTTPError {
	return &HTTPError{
		Message: message,
		Detail:  err.Error(),
	}
}

package httputil

import (
	"errors"
	"fmt"
	"net/http"
)

type HTTPError struct {
	code    int
	message string
	cause   error
}

func NewHTTPError(code int) *HTTPError {
	return &HTTPError{
		code:    code,
		message: http.StatusText(code),
	}
}

func (e *HTTPError) Message() string {
	return e.message
}

func (e *HTTPError) Code() int {
	return e.code
}

func (e *HTTPError) Unwrap() error {
	return e.cause
}

func (e *HTTPError) Error() string {
	if e.cause == nil {
		return fmt.Sprintf("[%d] %s", e.code, e.message)
	}

	return fmt.Sprintf("[%d] %s: %v", e.code, e.message, e.cause)
}

func (e *HTTPError) WithMessage(m string) *HTTPError {
	return &HTTPError{
		code:    e.code,
		message: m,
		cause:   e.cause,
	}
}

func (e *HTTPError) WithCause(err error) *HTTPError {
	return &HTTPError{
		code:    e.code,
		message: e.message,
		cause:   err,
	}
}

func ErrorFrom(err error) *HTTPError {
	if err == nil {
		return nil
	}

	if herr := new(HTTPError); errors.As(err, &herr) {
		return herr
	}

	return NewHTTPError(http.StatusInternalServerError).WithCause(err)
}

type ErrorEncoderFunc func(w http.ResponseWriter, r *http.Request, err error)

func defaultErrorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	herr := ErrorFrom(err)

	res := errorResponse{
		Code:  herr.Code(),
		Error: herr.Message(),
	}

	code := httpErrorStatusCode(herr.Code())

	if err := WriteJSON(w, code, res); err != nil {
		panic(err)
	}
}

type errorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

func httpErrorStatusCode(code int) int {
	if code >= 400 && code <= 599 {
		return code
	}

	return http.StatusInternalServerError
}

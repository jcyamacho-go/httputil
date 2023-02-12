package httputil

import (
	"errors"
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

func (e *HTTPError) Code() int {
	return e.code
}

func (e *HTTPError) Unwrap() error {
	return e.cause
}

func (e *HTTPError) Error() string {
	return e.message
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

type ErrorWriter interface {
	Write(w http.ResponseWriter, err error)
}

type ErrorWriterFunc func(w http.ResponseWriter, err error)

func (ew ErrorWriterFunc) Write(w http.ResponseWriter, err error) {
	ew(w, err)
}

var DefaultErrorWriter ErrorWriter = ErrorWriterFunc(func(w http.ResponseWriter, err error) {
	res := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	code := http.StatusInternalServerError

	he := &HTTPError{}
	if errors.As(err, &he) {
		code = he.Code()
	}

	if err := WriteJSON(w, code, res); err != nil {
		panic(err)
	}
})

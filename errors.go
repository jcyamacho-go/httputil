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

type ErrorEncoder interface {
	Encode(w http.ResponseWriter, r *http.Request, err error)
}

type ErrorEncoderFunc func(w http.ResponseWriter, r *http.Request, err error)

func (f ErrorEncoderFunc) Encode(w http.ResponseWriter, r *http.Request, err error) {
	f(w, r, err)
}

type errorResponse struct {
	Error string `json:"error"`
}

var DefaultErrorEncoder ErrorEncoder = ErrorEncoderFunc(func(w http.ResponseWriter, r *http.Request, err error) {
	res := errorResponse{
		Error: err.Error(),
	}

	code := http.StatusInternalServerError

	if herr := new(HTTPError); errors.As(err, &herr) {
		code = httpErrorStatusCode(herr.Code())
	}

	if err := WriteJSON(w, code, res); err != nil {
		panic(err)
	}
})

func httpErrorStatusCode(code int) int {
	if code >= 400 && code <= 599 {
		return code
	}

	return http.StatusInternalServerError
}

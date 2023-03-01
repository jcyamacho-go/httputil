package httputil

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrBadRequest                    = NewStatusError(http.StatusBadRequest)
	ErrUnauthorized                  = NewStatusError(http.StatusUnauthorized)
	ErrPaymentRequired               = NewStatusError(http.StatusPaymentRequired)
	ErrForbidden                     = NewStatusError(http.StatusForbidden)
	ErrNotFound                      = NewStatusError(http.StatusNotFound)
	ErrMethodNotAllowed              = NewStatusError(http.StatusMethodNotAllowed)
	ErrNotAcceptable                 = NewStatusError(http.StatusNotAcceptable)
	ErrProxyAuthRequired             = NewStatusError(http.StatusProxyAuthRequired)
	ErrRequestTimeout                = NewStatusError(http.StatusRequestTimeout)
	ErrConflict                      = NewStatusError(http.StatusConflict)
	ErrGone                          = NewStatusError(http.StatusGone)
	ErrLengthRequired                = NewStatusError(http.StatusLengthRequired)
	ErrPreconditionFailed            = NewStatusError(http.StatusPreconditionFailed)
	ErrRequestEntityTooLarge         = NewStatusError(http.StatusRequestEntityTooLarge)
	ErrRequestURITooLong             = NewStatusError(http.StatusRequestURITooLong)
	ErrUnsupportedMediaType          = NewStatusError(http.StatusUnsupportedMediaType)
	ErrRequestedRangeNotSatisfiable  = NewStatusError(http.StatusRequestedRangeNotSatisfiable)
	ErrExpectationFailed             = NewStatusError(http.StatusExpectationFailed)
	ErrTeapot                        = NewStatusError(http.StatusTeapot)
	ErrMisdirectedRequest            = NewStatusError(http.StatusMisdirectedRequest)
	ErrUnprocessableEntity           = NewStatusError(http.StatusUnprocessableEntity)
	ErrLocked                        = NewStatusError(http.StatusLocked)
	ErrFailedDependency              = NewStatusError(http.StatusFailedDependency)
	ErrTooEarly                      = NewStatusError(http.StatusTooEarly)
	ErrUpgradeRequired               = NewStatusError(http.StatusUpgradeRequired)
	ErrPreconditionRequired          = NewStatusError(http.StatusPreconditionRequired)
	ErrTooManyRequests               = NewStatusError(http.StatusTooManyRequests)
	ErrRequestHeaderFieldsTooLarge   = NewStatusError(http.StatusRequestHeaderFieldsTooLarge)
	ErrUnavailableForLegalReasons    = NewStatusError(http.StatusUnavailableForLegalReasons)
	ErrInternalServerError           = NewStatusError(http.StatusInternalServerError)
	ErrNotImplemented                = NewStatusError(http.StatusNotImplemented)
	ErrBadGateway                    = NewStatusError(http.StatusBadGateway)
	ErrServiceUnavailable            = NewStatusError(http.StatusServiceUnavailable)
	ErrGatewayTimeout                = NewStatusError(http.StatusGatewayTimeout)
	ErrHTTPVersionNotSupported       = NewStatusError(http.StatusHTTPVersionNotSupported)
	ErrVariantAlsoNegotiates         = NewStatusError(http.StatusVariantAlsoNegotiates)
	ErrInsufficientStorage           = NewStatusError(http.StatusInsufficientStorage)
	ErrLoopDetected                  = NewStatusError(http.StatusLoopDetected)
	ErrNotExtended                   = NewStatusError(http.StatusNotExtended)
	ErrNetworkAuthenticationRequired = NewStatusError(http.StatusNetworkAuthenticationRequired)
)

type HTTPError struct {
	code    int
	message string
	cause   error
}

func NewStatusError(code int) *HTTPError {
	return &HTTPError{
		code:    code,
		message: http.StatusText(code),
	}
}

func NewError(code int, message string) *HTTPError {
	return &HTTPError{
		code:    code,
		message: message,
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

	return ErrInternalServerError.WithCause(err)
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

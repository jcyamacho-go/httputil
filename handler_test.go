package httputil

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_ServeHTTP(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	code := http.StatusOK
	bodyText := "response body text"

	handler := NewHandler(func(w http.ResponseWriter, r *http.Request) error {
		return WriteString(w, code, bodyText)
	})

	handler.ServeHTTP(w, r)

	if w.Code != code {
		t.Errorf("invalid response code, got: %d, want: %d", w.Code, code)
	}

	if res := w.Body.String(); res != bodyText {
		t.Errorf("invalid response body, got: %s, want: %s", res, bodyText)
	}
}

func TestHandler_WithMiddleware(t *testing.T) {
	var calls []string

	middleware := func(name string) Middleware {
		return func(next HandlerFunc) HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) error {
				calls = append(calls, fmt.Sprintf("%s-start", name))

				defer func() {
					calls = append(calls, fmt.Sprintf("%s-end", name))
				}()

				return next(w, r)
			}
		}
	}

	handler := NewHandler(func(w http.ResponseWriter, r *http.Request) error {
		calls = append(calls, "handler")
		return nil
	}).WithMiddlewares(
		middleware("mdw-1"),
		middleware("mdw-2"),
	)

	handler.ServeHTTP(nil, nil)

	got := strings.Join(calls, ",")
	want := "mdw-1-start,mdw-2-start,handler,mdw-2-end,mdw-1-end"

	if got != want {
		t.Errorf("invalid middleware calls, got: %s, want: %s", got, want)
	}
}

func TestHandler_WithErrorEncoder(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	code := http.StatusBadRequest
	handlerError := errors.New("handler error")
	errorEncoderCalls := 0

	var errorEncoder ErrorEncoderFunc = func(w http.ResponseWriter, r *http.Request, err error) {
		errorEncoderCalls++
		_ = WriteString(w, code, err.Error())
	}

	handler := NewHandler(func(w http.ResponseWriter, r *http.Request) error {
		return handlerError
	}).WithErrorEncoder(errorEncoder)

	handler.ServeHTTP(w, r)

	if errorEncoderCalls != 1 {
		t.Errorf("invalid error encoder calls, got: %d, want: 1", errorEncoderCalls)
	}

	if w.Code != code {
		t.Errorf("invalid response code, got: %d, want: %d", w.Code, code)
	}

	if res := w.Body.String(); res != handlerError.Error() {
		t.Errorf("invalid response body, got: %s, want: %v", res, handlerError)
	}
}

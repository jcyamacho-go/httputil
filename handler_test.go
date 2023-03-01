package httputil

import (
	"errors"
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

	var middleware1 Middleware = func(next HandlerFunc) HandlerFunc {
		calls = append(calls, "mdw-1")
		return next
	}

	var middleware2 Middleware = func(next HandlerFunc) HandlerFunc {
		calls = append(calls, "mdw-2")
		return next
	}

	handler := NewHandler(func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}).WithMiddlewares(
		middleware1,
		middleware2,
	)

	handler.ServeHTTP(nil, nil)

	got := strings.Join(calls, ",")
	want := "mdw-2,mdw-1"

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

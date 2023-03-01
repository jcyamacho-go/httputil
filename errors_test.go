package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestErrorFrom_Std_Error(t *testing.T) {
	err := errors.New("demo error")

	herr := ErrorFrom(err)

	msg := fmt.Sprintf("[500] Internal Server Error: %v", err)
	if got := herr.Error(); got != msg {
		t.Errorf("invalid error string value, got: %s, want %s", got, msg)
	}
}

func TestErrorFrom_Wrapped(t *testing.T) {
	err := NewStatusError(http.StatusBadRequest).
		WithMessage("error message")

	nerr := fmt.Errorf("demo error: %w", err)
	herr := ErrorFrom(nerr)

	if err != herr {
		t.Errorf("invalid error, got: %v, want %v", herr, err)
	}
}

func FuzzHTTPError_Code(f *testing.F) {
	f.Add(400)
	f.Add(402)
	f.Add(410)
	f.Add(500)
	f.Add(503)

	f.Fuzz(func(t *testing.T, code int) {
		err := NewStatusError(code)
		if c := err.Code(); c != code {
			t.Errorf("invalid code, want: %d, got: %d", code, c)
		}
	})
}

func TestDefaultErrorEncoder_StatusCode500(t *testing.T) {
	err := errors.New("error message")

	w := httptest.NewRecorder()
	defaultErrorEncoder(w, nil, err)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("invalid code, want: 500, got: %d", w.Code)
	}

	body := `{"code":500,"error":"Internal Server Error"}`
	if v := strings.TrimRight(w.Body.String(), "\n"); v != body {
		t.Errorf("invalid body, want: %s, got: %s!", body, v)
	}
}

func FuzzDefaultErrorEncoder(f *testing.F) {
	f.Add(400, http.StatusText(400))
	f.Add(403, http.StatusText(403))
	f.Add(404, "error-1")
	f.Add(500, "error-2")
	f.Add(503, "error: cause")

	f.Fuzz(func(t *testing.T, code int, message string) {
		w := httptest.NewRecorder()
		err := NewError(code, message)

		defaultErrorEncoder(w, nil, err)

		if v := httpErrorStatusCode(code); w.Code != v {
			t.Errorf("invalid code, want: %d, got: %d", v, w.Code)
		}

		body, _ := json.Marshal(map[string]any{
			"code":  code,
			"error": message,
		})

		if v := strings.TrimRight(w.Body.String(), "\n"); v != string(body) {
			t.Errorf("invalid body, want: %s, got: %s!", body, v)
		}
	})
}

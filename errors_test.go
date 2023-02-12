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

func TestDefaultErrorWriter_StatusCode500(t *testing.T) {
	message := "Unauthorized"
	err := errors.New(message)

	w := httptest.NewRecorder()
	DefaultErrorWriter.Write(w, err)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("invalid code, want: 500, got: %d", w.Code)
	}

	body := fmt.Sprintf(`{"error":"%s"}`, message)
	if v := strings.TrimRight(w.Body.String(), "\n"); v != body {
		t.Errorf("invalid body, want: %s, got: %s!", body, v)
	}
}

func FuzzDefaultErrorWriter(f *testing.F) {
	f.Add(400, http.StatusText(400))
	f.Add(403, http.StatusText(403))
	f.Add(404, "error-1")
	f.Add(500, "error-2")
	f.Add(503, "error: error message")

	f.Fuzz(func(t *testing.T, code int, message string) {
		w := httptest.NewRecorder()
		err := NewHTTPError(code).WithMessage(message)

		DefaultErrorWriter.Write(w, err)

		if v := httpErrorStatusCode(code); w.Code != v {
			t.Errorf("invalid code, want: %d, got: %d", v, w.Code)
		}

		body, _ := json.Marshal(map[string]any{
			"error": message,
		})

		if v := strings.TrimRight(w.Body.String(), "\n"); v != string(body) {
			t.Errorf("invalid body, want: %s, got: %s!", body, v)
		}
	})
}

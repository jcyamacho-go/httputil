package pprof

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", DefaultPrefix, nil)

	h := Handler()
	h.ServeHTTP(w, r)

	text := w.Body.String()
	if !strings.Contains(text, "profiles available") {
		t.Errorf("expected available profiles, got:\n%s", text)
	}
}

func TestMiddleware(t *testing.T) {
	original := handler

	t.Cleanup(func() {
		handler = original
	})

	tests := []struct {
		pprofCalls int
		nextCalls  int
		path       string
	}{
		{
			pprofCalls: 1,
			nextCalls:  0,
			path:       DefaultPrefix,
		},
		{
			pprofCalls: 0,
			nextCalls:  1,
			path:       "/debug/pproff",
		},
		{
			pprofCalls: 1,
			nextCalls:  0,
			path:       DefaultPrefix + "heap",
		},
		{
			pprofCalls: 0,
			nextCalls:  1,
			path:       "/v1/" + DefaultPrefix + "heap",
		},
		{
			pprofCalls: 0,
			nextCalls:  1,
			path:       "/api/path",
		},
		{
			pprofCalls: 0,
			nextCalls:  1,
			path:       "/debug/demo",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.path, func(t *testing.T) {
			ph := new(handlerMock)
			nh := new(handlerMock)

			handler = ph
			h := Middleware(nh)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", test.path, nil)

			h.ServeHTTP(w, r)

			if test.pprofCalls != ph.calls {
				t.Errorf("invalid pprof handler calls, got: %d, want: %d", ph.calls, test.pprofCalls)
			}

			if test.nextCalls != nh.calls {
				t.Errorf("invalid next handler calls, got: %d, want: %d", nh.calls, test.nextCalls)
			}
		})
	}
}

type handlerMock struct {
	calls int
}

func (h *handlerMock) ServeHTTP(http.ResponseWriter, *http.Request) {
	h.calls++
}

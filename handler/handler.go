package handler

import "net/http"

type Handler func(w http.ResponseWriter, r *http.Request) error

func New(h Handler, options ...Option) http.HandlerFunc {
	cfg := Config{
		errorWriter: ErrorWriterFunc(DefaultErrorWriter),
	}

	for _, option := range options {
		option(&cfg)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			cfg.errorWriter.Write(w, err)
		}
	}
}

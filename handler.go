package httputil

import "net/http"

type Config struct {
	errorWriter ErrorWriter
}

type Option func(*Config)

func WithErrorWriter(ew ErrorWriter) Option {
	return func(c *Config) {
		c.errorWriter = ew
	}
}

type Handler func(w http.ResponseWriter, r *http.Request) error

func NewHandler(h Handler, options ...Option) http.HandlerFunc {
	cfg := Config{
		errorWriter: DefaultErrorWriter,
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

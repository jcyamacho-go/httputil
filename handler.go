package httputil

import "net/http"

type Config struct {
	errorEncoder ErrorEncoder
}

type Option func(*Config)

func WithErrorWriter(ew ErrorEncoder) Option {
	return func(c *Config) {
		c.errorEncoder = ew
	}
}

type Handler func(w http.ResponseWriter, r *http.Request) error

func NewHandler(h Handler, options ...Option) http.HandlerFunc {
	cfg := Config{
		errorEncoder: DefaultErrorEncoder,
	}

	for _, option := range options {
		option(&cfg)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			cfg.errorEncoder.Encode(w, r, err)
		}
	}
}

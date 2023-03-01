package httputil

import "net/http"

type (
	HandlerFunc func(w http.ResponseWriter, r *http.Request) error
	Middleware  func(next HandlerFunc) HandlerFunc
)

type Handler struct {
	errorEncoder ErrorEncoderFunc
	middlewares  []Middleware
	handler      HandlerFunc
}

func NewHandler(h HandlerFunc) *Handler {
	return &Handler{
		handler:      h,
		errorEncoder: defaultErrorEncoder,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := h.handler

	for i := len(h.middlewares) - 1; i >= 0; i-- {
		handler = h.middlewares[i](handler)
	}

	if err := handler(w, r); err != nil {
		h.errorEncoder(w, r, err)
	}
}

func (h *Handler) WithErrorEncoder(errorEncoder ErrorEncoderFunc) *Handler {
	h.errorEncoder = errorEncoder
	return h
}

func (h *Handler) WithMiddlewares(middlewares ...Middleware) *Handler {
	h.middlewares = middlewares
	return h
}

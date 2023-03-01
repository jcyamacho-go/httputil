package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jcyamacho-go/httputil"
	"github.com/jcyamacho-go/httputil/pprof"
)

func logger(next httputil.HandlerFunc) httputil.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		start := time.Now()

		defer func() {
			duration := time.Since(start)
			log.Printf("[%s] %s (%v): %v \n", r.Method, r.URL, duration, err)
		}()

		return next(w, r)
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(pprof.Middleware)

	r.Route("/hello", func(r chi.Router) {
		handler := httputil.NewHandler(helloHandler).
			WithMiddlewares(logger).
			ServeHTTP

		r.Get("/", handler)
		r.Get("/{name}", handler)
	})

	r.Route("/error", func(r chi.Router) {
		handler := httputil.NewHandler(errorHandler).
			WithMiddlewares(logger).
			ServeHTTP

		r.Get("/", handler)
		r.Get("/{code}", handler)
		r.Get("/{code}/{message}", handler)
	})

	addr := "0.0.0.0:3000"

	log.Printf("server running in: %s \n", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request) error {
	code := http.StatusInternalServerError

	if sc := chi.URLParam(r, "code"); sc != "" {
		if c, err := strconv.Atoi(sc); err == nil {
			code = c
		}
	}

	if message := chi.URLParam(r, "message"); message != "" {
		return httputil.NewError(code, message)
	}

	return httputil.NewStatusError(code)
}

func helloHandler(w http.ResponseWriter, r *http.Request) error {
	var in helloInput

	if err := httputil.BindQuery(r, &in); err != nil {
		return httputil.ErrBadRequest.WithCause(err)
	}

	if name := chi.URLParam(r, "name"); name != "" {
		in.Name = name
	}

	reply := helloReply{
		Message: fmt.Sprintf("hello %s!", in.Name),
	}

	return httputil.WriteJSON(w, http.StatusOK, reply)
}

type helloInput struct {
	Name string `query:"name"`
}

type helloReply struct {
	Message string `json:"message"`
}

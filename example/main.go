package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jcyamacho-go/httputil"
	"github.com/jcyamacho-go/httputil/pprof"
)

func logger(next httputil.HandlerFunc) httputil.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		start := time.Now()
		err := next(w, r)
		duration := time.Since(start)

		log.Printf("[%s] %s (%v) %v \n", r.Method, r.URL, duration, err)

		return err
	}
}

func main() {
	r := chi.NewRouter()
	r.Use(pprof.Middleware)

	r.Route("/hello", func(r chi.Router) {
		handler := httputil.NewHandler(hello).
			WithMiddlewares(logger).
			ServeHTTP

		r.Get("/", handler)
		r.Get("/{name}", handler)
	})

	addr := "0.0.0.0:3000"

	log.Printf("server running in: %s \n", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) error {
	var in helloInput

	if err := httputil.BindQuery(r, &in); err != nil {
		return httputil.NewHTTPError(http.StatusBadRequest).
			WithCause(err)
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
	Message string `json:"name"`
}

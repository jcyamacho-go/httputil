package main

import (
	"fmt"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/jcyamacho/httputil"
	"github.com/jcyamacho/httputil/pprof"
)

func main() {
	r := chi.NewRouter()
	r.Use(pprof.Middleware)

	r.Get("/hello", httputil.NewHandler(hello))

	http.ListenAndServe("0.0.0.0:3000", r)
}

func hello(w http.ResponseWriter, r *http.Request) error {
	var in helloInput

	if err := httputil.BindQuery(r, &in); err != nil {
		return httputil.NewHTTPError(http.StatusBadRequest).
			WithCause(err)
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

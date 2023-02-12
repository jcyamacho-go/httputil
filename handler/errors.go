package handler

import (
	"net/http"

	"github.com/jcyamacho/httputil/response"
)

type ErrorWriter interface {
	Write(w http.ResponseWriter, err error)
}

type ErrorWriterFunc func(w http.ResponseWriter, err error)

func (ew ErrorWriterFunc) Write(w http.ResponseWriter, err error) {
	ew(w, err)
}

func DefaultErrorWriter(w http.ResponseWriter, err error) {
	res := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	if err := response.JSON(w, http.StatusInternalServerError, res); err != nil {
		panic(err)
	}
}

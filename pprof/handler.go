package pprof

import (
	"net/http"
	"net/http/pprof"
)

// NewHandler returns a a pprof handler for the routes: /debug/pprof/*
func NewHandler() http.Handler {
	m := http.NewServeMux()

	m.HandleFunc("/debug/pprof/", pprof.Index)
	m.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	m.HandleFunc("/debug/pprof/profile", pprof.Profile)
	m.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	m.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return m
}

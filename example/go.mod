module github.com/jcyamacho/httputil/example

go 1.18

replace github.com/jcyamacho/httputil => ../

require (
	github.com/go-chi/chi/v5 v5.0.8
	github.com/jcyamacho/httputil v0.0.0-00010101000000-000000000000
)

require github.com/go-playground/form/v4 v4.2.0 // indirect

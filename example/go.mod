module github.com/jcyamacho-go/httputil/example

go 1.18

replace github.com/jcyamacho-go/httputil => ../

require (
	github.com/go-chi/chi/v5 v5.0.8
	github.com/jcyamacho-go/httputil v0.0.6
)

require github.com/go-playground/form/v4 v4.2.0 // indirect

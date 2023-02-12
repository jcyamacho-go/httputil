package bind

import (
	"net/http"

	"github.com/go-playground/form/v4"
)

func Query(r *http.Request, v any, options ...FormOption) error {
	decoder := form.NewDecoder()
	decoder.SetTagName("query")

	for _, option := range options {
		option(decoder)
	}

	query := r.URL.Query()

	return decoder.Decode(v, query)
}

package bind

import (
	"encoding/xml"
	"net/http"
)

type XMLOption func(*xml.Decoder)

func XML(r *http.Request, v any, options ...XMLOption) error {
	decoder := xml.NewDecoder(r.Body)
	for _, option := range options {
		option(decoder)
	}

	return decoder.Decode(v)
}

package bind

import (
	"net/url"

	"github.com/go-playground/form/v4"
)

func URLValues(values url.Values, v any, options ...FormOption) error {
	decoder := form.NewDecoder()
	decoder.SetTagName("json")

	for _, option := range options {
		option(decoder)
	}

	return decoder.Decode(v, values)
}

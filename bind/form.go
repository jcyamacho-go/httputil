package bind

import (
	"net/http"

	"github.com/go-playground/form/v4"
)

type FormOption func(*form.Decoder)

func Form(r *http.Request, v any, options ...FormOption) error {
	decoder := form.NewDecoder()
	for _, option := range options {
		option(decoder)
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	return decoder.Decode(v, r.Form)
}

func TagName(tagName string) FormOption {
	return func(d *form.Decoder) {
		d.SetTagName(tagName)
	}
}

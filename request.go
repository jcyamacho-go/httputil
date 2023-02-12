package httputil

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/url"

	"github.com/go-playground/form/v4"
)

const (
	TagNameForm  = "form"
	TagNameQuery = "query"
	TagNameJSON  = "json"
)

type FormOption func(*form.Decoder)

func WithTagName(tagName string) FormOption {
	return func(d *form.Decoder) {
		d.SetTagName(tagName)
	}
}

func BindForm(r *http.Request, v any, options ...FormOption) error {
	decoder := form.NewDecoder()
	decoder.SetTagName(TagNameForm)

	for _, option := range options {
		option(decoder)
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	return decoder.Decode(v, r.Form)
}

func BindJSON(r *http.Request, v any) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

func BindQuery(r *http.Request, v any, options ...FormOption) error {
	decoder := form.NewDecoder()
	decoder.SetTagName(TagNameQuery)

	for _, option := range options {
		option(decoder)
	}

	query := r.URL.Query()

	return decoder.Decode(v, query)
}

func BindValues(values url.Values, v any, options ...FormOption) error {
	decoder := form.NewDecoder()
	decoder.SetTagName(TagNameJSON)

	for _, option := range options {
		option(decoder)
	}

	return decoder.Decode(v, values)
}

type XMLOption func(*xml.Decoder)

func BindXML(r *http.Request, v any, options ...XMLOption) error {
	decoder := xml.NewDecoder(r.Body)
	for _, option := range options {
		option(decoder)
	}

	return decoder.Decode(v)
}

package response

import (
	"net/http"

	"github.com/jcyamacho/httputil/header"
)

func Blob(w http.ResponseWriter, code int, contentType string, data []byte) error {
	w.Header().Set(header.ContentType, contentType)
	w.WriteHeader(code)
	_, err := w.Write(data)

	return err
}

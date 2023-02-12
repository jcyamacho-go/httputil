package response

import (
	"io"
	"net/http"

	"github.com/jcyamacho/httputil/header"
)

func Stream(w http.ResponseWriter, code int, contentType string, reader io.Reader) error {
	w.Header().Set(header.ContentType, contentType)
	w.WriteHeader(code)
	_, err := io.Copy(w, reader)

	return err
}

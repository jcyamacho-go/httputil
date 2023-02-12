package response

import (
	"encoding/xml"
	"net/http"

	"github.com/jcyamacho/httputil/header"
	"github.com/jcyamacho/httputil/mime"
)

func XML(w http.ResponseWriter, code int, v any) error {
	w.Header().Set(header.ContentType, mime.ApplicationXMLCharsetUTF8)
	w.WriteHeader(code)

	return xml.NewEncoder(w).Encode(v)
}

func XMLBlob(w http.ResponseWriter, code int, data []byte) error {
	return Blob(w, code, mime.ApplicationXMLCharsetUTF8, data)
}

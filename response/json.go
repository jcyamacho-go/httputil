package response

import (
	"encoding/json"
	"net/http"

	"github.com/jcyamacho/httputil/header"
	"github.com/jcyamacho/httputil/mime"
)

func JSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set(header.ContentType, mime.ApplicationJSONCharsetUTF8)
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(v)
}

func JSONBlob(w http.ResponseWriter, code int, data []byte) error {
	return Blob(w, code, mime.ApplicationJSONCharsetUTF8, data)
}

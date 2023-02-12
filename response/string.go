package response

import (
	"net/http"

	"github.com/jcyamacho/httputil/mime"
)

func String(w http.ResponseWriter, code int, text string) error {
	return Blob(w, code, mime.TextPlainCharsetUTF8, []byte(text))
}

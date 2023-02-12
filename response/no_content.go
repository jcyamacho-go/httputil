package response

import "net/http"

func NoContent(w http.ResponseWriter, code int) error {
	w.WriteHeader(code)
	return nil
}

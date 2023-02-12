package bind

import (
	"encoding/json"
	"net/http"
)

func JSON(r *http.Request, v any) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

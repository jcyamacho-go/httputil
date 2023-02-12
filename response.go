package httputil

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

func WriteBlob(w http.ResponseWriter, code int, contentType string, data []byte) error {
	w.Header().Set(HeaderContentType, contentType)
	w.WriteHeader(code)
	_, err := w.Write(data)

	return err
}

func WriteJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set(HeaderContentType, MimeApplicationJSONCharsetUTF8)
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(v)
}

func WriteJSONBlob(w http.ResponseWriter, code int, data []byte) error {
	return WriteBlob(w, code, MimeApplicationJSONCharsetUTF8, data)
}

func WriteNoContent(w http.ResponseWriter, code int) error {
	w.WriteHeader(code)
	return nil
}

func WriteStream(w http.ResponseWriter, code int, contentType string, reader io.Reader) error {
	w.Header().Set(HeaderContentType, contentType)
	w.WriteHeader(code)
	_, err := io.Copy(w, reader)

	return err
}

func WriteString(w http.ResponseWriter, code int, text string) error {
	return WriteBlob(w, code, MimeTextPlainCharsetUTF8, []byte(text))
}

func WriteXML(w http.ResponseWriter, code int, v any) error {
	w.Header().Set(HeaderContentType, MimeApplicationXMLCharsetUTF8)
	w.WriteHeader(code)

	return xml.NewEncoder(w).Encode(v)
}

func WritXMLBlob(w http.ResponseWriter, code int, data []byte) error {
	return WriteBlob(w, code, MimeApplicationXMLCharsetUTF8, data)
}

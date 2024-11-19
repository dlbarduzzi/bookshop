package jsontil

import (
	"encoding/json"
	"net/http"
)

func Marshal(w http.ResponseWriter, data interface{}, code int, headers http.Header) error {
	res, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	res = append(res, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)

	_, err = w.Write(res)
	return err
}

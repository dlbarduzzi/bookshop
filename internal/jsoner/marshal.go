package jsoner

import (
	"encoding/json"
	"net/http"
)

type Envelope map[string]interface{}

func Marshal(w http.ResponseWriter, data Envelope, status int, headers http.Header) error {
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}

	res = append(res, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(res)
	return err
}

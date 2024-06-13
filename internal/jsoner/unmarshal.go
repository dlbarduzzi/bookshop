package jsoner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// maxBodyBytes is the maximum number of bytes a json payload request size can be.
const maxBodyBytes = 1_048_576 // 1MB

func Unmarshal(w http.ResponseWriter, r *http.Request, data interface{}) (int, error) {
	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	if err := d.Decode(data); err != nil {
		var (
			syntaxError           *json.SyntaxError
			maxBytesError         *http.MaxBytesError
			unmarshalTypeError    *json.UnmarshalTypeError
			invalidUnmarshalError *json.InvalidUnmarshalError
		)
		switch {
		case errors.As(err, &syntaxError):
			return http.StatusBadRequest,
				fmt.Errorf("request body contains malformed json at position %d",
					syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return http.StatusBadRequest,
				errors.New("request body contains malformed json")

		case errors.As(err, &unmarshalTypeError):
			return http.StatusBadRequest,
				fmt.Errorf("request body contains incorrect json type %q at position %d",
					unmarshalTypeError.Field, unmarshalTypeError.Offset)

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			if field == "\"\"" {
				return http.StatusBadRequest,
					errors.New("request body contains empty json field name")
			} else {
				return http.StatusBadRequest,
					fmt.Errorf("request body contains unknown json field %s", field)
			}

		case errors.Is(err, io.EOF):
			return http.StatusBadRequest,
				errors.New("request body cannot be empty")

		case errors.As(err, &maxBytesError):
			return http.StatusRequestEntityTooLarge,
				fmt.Errorf("request body cannot be larger than %d bytes",
					maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return http.StatusInternalServerError,
				fmt.Errorf("failed to decode json request body; %w", err)
		}
	}

	if d.More() {
		return http.StatusBadRequest,
			errors.New("request body cannot have more than 1 json object")
	}

	return http.StatusOK, nil
}

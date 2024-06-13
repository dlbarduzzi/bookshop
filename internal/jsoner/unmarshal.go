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
		var syntaxErr *json.SyntaxError
		var maxBytesErr *http.MaxBytesError
		var unmarshalErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxErr):
			return malformedPosition(syntaxErr.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return malformed()
		case errors.As(err, &unmarshalErr):
			return incorrectType(unmarshalErr.Field, unmarshalErr.Offset)
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			return invalidField(err)
		case errors.Is(err, io.EOF):
			return bodyEmpty()
		case errors.As(err, &maxBytesErr):
			return bodyTooLarge(maxBytesErr.Limit)
		default:
			return serverError(err)
		}
	}

	if d.More() {
		return multipleBodies()
	}

	return http.StatusOK, nil
}

func malformedPosition(position int64) (int, error) {
	err := fmt.Errorf("request body contains malformed json at position %d", position)
	code := http.StatusBadRequest
	return code, err
}

func malformed() (int, error) {
	err := errors.New("request body contains malformed json")
	code := http.StatusBadRequest
	return code, err
}

func incorrectType(field string, position int64) (int, error) {
	err := fmt.Errorf("request body contains incorrect json type %q at position %d", field, position)
	code := http.StatusBadRequest
	return code, err
}

func invalidField(err error) (int, error) {
	code := http.StatusBadRequest
	field := strings.TrimPrefix(err.Error(), "json: unknown field ")

	if field == "\"\"" {
		err = errors.New("request body contains empty json field name")
	} else {
		err = fmt.Errorf("request body contains unknown json field %s", field)
	}

	return code, err
}

func bodyEmpty() (int, error) {
	err := errors.New("request body cannot be empty")
	code := http.StatusBadRequest
	return code, err
}

func bodyTooLarge(limit int64) (int, error) {
	err := fmt.Errorf("request body cannot be larger than %d bytes", limit)
	code := http.StatusRequestEntityTooLarge
	return code, err
}

func serverError(err error) (int, error) {
	msg := fmt.Errorf("failed to decode json request body; %w", err)
	code := http.StatusInternalServerError
	return code, msg
}

func multipleBodies() (int, error) {
	err := errors.New("request body cannot have more than 1 json object")
	code := http.StatusBadRequest
	return code, err
}

package jsoner

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarshal(t *testing.T) {
	t.Parallel()

	data := Envelope{"foo": "bar"}

	headers := http.Header{}
	headers.Add("Foo", "bar")

	w := httptest.NewRecorder()
	err := Marshal(w, data, http.StatusOK, headers)
	want := http.StatusOK

	if w.Code != want {
		t.Errorf("expected status code to be %d; got %d", want, w.Code)
	}

	if err != nil {
		t.Errorf("expected error to be nil; got %v", err)
	}
}

func TestMarshalError(t *testing.T) {
	t.Parallel()

	data := Envelope{"foo": "bar"}
	data["new-foo"] = data

	w := httptest.NewRecorder()
	err := Marshal(w, data, 0, nil)
	want := "json: unsupported value: encountered a cycle via jsoner.Envelope"

	if err == nil || err.Error() != want {
		t.Errorf("expected error to be %v; got %v", want, err)
	}
}

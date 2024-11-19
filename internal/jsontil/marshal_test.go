package jsontil

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarshal(t *testing.T) {
	t.Parallel()

	res := map[string]interface{}{"foo": "bar"}

	entryKey := "Foo"
	entryValue := "Bar"

	headers := http.Header{}
	headers.Add(entryKey, entryValue)

	w := httptest.NewRecorder()

	err := Marshal(w, res, http.StatusOK, headers)
	if err != nil {
		t.Errorf("expected error to be nil; got %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected status code to be %d; got %d", http.StatusOK, w.Code)
	}

	headerValue := w.Header().Get(entryKey)
	if headerValue != entryValue {
		t.Errorf("expected header %s to be %s; got %s", entryKey, entryValue, headerValue)
	}
}

package validator

import "testing"

func TestIsValid(t *testing.T) {
	t.Parallel()

	v := NewValidator()
	v.AddError("foo", "bar")

	if v.IsValid() {
		t.Error("expected validator to be false")
	}
}

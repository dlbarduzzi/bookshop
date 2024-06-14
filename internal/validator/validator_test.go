package validator

import (
	"strings"
	"testing"
)

func TestIsValid(t *testing.T) {
	t.Parallel()

	v := NewValidator()
	v.AddError("foo", "bar")

	if v.IsValid() {
		t.Error("expected validator to be false")
	}
}

func TestCheck(t *testing.T) {
	t.Parallel()

	v := NewValidator()
	v.Check(false, "foo", "bar")

	if v.IsValid() {
		t.Error("expected validator to be false")
	}

	if v.Errors["foo"] != "bar" {
		t.Errorf("expected foo error to be bar; got %v", v.Errors["foo"])
	}
}

func TestMatchRegex(t *testing.T) {
	t.Parallel()

	email := "test.user@email.com"
	if !MatchRegex(email, EmailRegex) {
		t.Errorf("expected email %s to be valid", email)
	}

	email = "test.user"
	if MatchRegex(email, EmailRegex) {
		t.Errorf("expected email %s not to be valid", email)
	}
}

func TestValueInList(t *testing.T) {
	t.Parallel()

	names := []string{"foo", "bar"}

	name := "foo"
	if !ValueInList(name, names...) {
		t.Errorf("expected value %s to be in list", name)
	}

	name = "baz"
	if ValueInList(name, names...) {
		t.Errorf("expected value %s not to be in list", name)
	}
}

func TestValuesInList(t *testing.T) {
	t.Parallel()

	unique := []string{"foo", "bar", "baz"}
	if !ValuesAreUnique(unique) {
		t.Errorf("expected values %s to be unique", strings.Join(unique, ","))
	}

	notUnique := append(unique, "foo")
	if ValuesAreUnique(notUnique) {
		t.Errorf("expected values %s not to be unique", strings.Join(unique, ","))
	}
}

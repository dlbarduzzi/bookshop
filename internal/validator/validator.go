package validator

import (
	"regexp"
	"slices"
)

var EmailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Errors map[string]string

type Validator struct {
	Errors
}

func NewValidator() *Validator {
	return &Validator{Errors: make(Errors)}
}

func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(field string, value string) {
	if _, exists := v.Errors[field]; !exists {
		v.Errors[field] = value
	}
}

func (v *Validator) Check(ok bool, field string, value string) {
	if !ok {
		v.AddError(field, value)
	}
}

func MatchRegex(value string, re *regexp.Regexp) bool {
	return re.MatchString(value)
}

func ValueInList[T comparable](value T, list ...T) bool {
	return slices.Contains(list, value)
}

func ValuesAreUnique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}

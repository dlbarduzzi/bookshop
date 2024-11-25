package validator

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

func (v *Validator) AddError(key string, value string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = value
	}
}

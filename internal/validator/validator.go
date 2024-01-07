package validator

import (
	"strings"
	"unicode/utf8"
)

// Define a new Validator type which contains a map of validation errors for our
// form fields.
type Validator struct {
	FieldErrors map[string]string
}

// return true if the FieldErrors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// add an error message to the FieldErrors map (so long as no
// entry already exists for the given key).
func (v *Validator) AddFieldError(key, message string) {
	// init the map first, if it does not exist
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// add an error message to the FieldErrors map only if a
// validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// return true if a value is not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// return true if a value contains no more than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// returns true if a value is in a list of permitted integers.
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

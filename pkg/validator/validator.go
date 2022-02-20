package validator

import (
	"fmt"
	"regexp"
)

var (
	SpecialCharRX = func(min, max int) *regexp.Regexp {
		return regexp.MustCompile(`^(.*?[\*\.!@#\$%\^&\(\)\{\}\[\]:;<>,.\?\\/~_\+\-=\|]){1,}.*$`)
	}
	DigitRX = func(min, max int) *regexp.Regexp {
		return regexp.MustCompile(fmt.Sprintf("^(.*?[0-9]){%d,%d}.*$", min, max))
	}
	LowercaseRX = func(min, max int) *regexp.Regexp {
		return regexp.MustCompile(fmt.Sprintf("^(.*?[a-z]){%d,%d}.*$", min, max))
	}
	UppercaseRX = func(min, max int) *regexp.Regexp {
		return regexp.MustCompile(fmt.Sprintf("^(.*?[A-Z]){%d,%d}.*$", min, max))
	}
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	Errors Errors
	tmp    struct {
		keyError string
	}
}

type Errors map[string][]string

func New() *Validator {
	return &Validator{
		Errors: make(Errors),
	}
}

// Valid checks if the validator has errors.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) SetTmpKeyError(keyError string) *Validator {
	v.tmp.keyError = keyError
	return v
}

func (v *Validator) AddError(message string) *Validator {
	v.Errors[v.tmp.keyError] = append(v.Errors[v.tmp.keyError], message)
	return nil
}

// Check adds an error message to the map only if a validation check is not 'ok'.
func (v *Validator) Check(ok bool, message string) bool {
	if !ok {
		v.AddError(message)
		return false
	}
	return true
}

// In returns true if a specific value is in a list of strings.
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}

	return false
}

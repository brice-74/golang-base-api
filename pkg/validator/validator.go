package validator

import (
	"fmt"
	"regexp"
)

func setMinMaxStr(min, max int) (minstr, maxstr string) {
	minstr = fmt.Sprint(min)
	maxstr = fmt.Sprint(max)
	if min == 0 {
		minstr = ""
	}
	if max == 0 {
		maxstr = ""
	}
	return minstr, maxstr
}

var (
	SpecialCharRX = func(min, max int) *regexp.Regexp {
		minstr, maxstr := setMinMaxStr(min, max)
		return regexp.MustCompile(`^(.*?[\*\.!@#\$%\^&\(\)\{\}\[\]:;<>,.\?\\/~_\+\-=\|]){` + minstr + `,` + maxstr + `}.*$`)
	}
	DigitRX = func(min, max int) *regexp.Regexp {
		minstr, maxstr := setMinMaxStr(min, max)
		return regexp.MustCompile(fmt.Sprintf("^(.*?[0-9]){%s,%s}.*$", minstr, maxstr))
	}
	LowercaseRX = func(min, max int) *regexp.Regexp {
		minstr, maxstr := setMinMaxStr(min, max)
		return regexp.MustCompile(fmt.Sprintf("^(.*?[a-z]){%s,%s}.*$", minstr, maxstr))
	}
	UppercaseRX = func(min, max int) *regexp.Regexp {
		minstr, maxstr := setMinMaxStr(min, max)
		return regexp.MustCompile(fmt.Sprintf("^(.*?[A-Z]){%s,%s}.*$", minstr, maxstr))
	}
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	Errors   Errors
	KeyError string
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

func (v *Validator) AddError(message string) *Validator {
	v.Errors[v.KeyError] = append(v.Errors[v.KeyError], message)
	return nil
}

// Check adds an error message to the map only if a validation check is not 'ok'.
func (v *Validator) Check(ok bool, message string) {
	if !ok {
		v.AddError(message)
	}
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

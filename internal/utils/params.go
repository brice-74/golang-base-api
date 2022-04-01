package utils

import (
	"regexp"
	"strings"

	"github.com/brice-74/golang-base-api/pkg/validator"
)

type QueryParams struct {
	Offset         int
	Limit          int
	Sort           string
	SortableFields []string
}

// SortColumn checks if the Sort value is allowed from the SortableFields list
// and returns the value without the specified direction.
func (p QueryParams) SortColumn() string {
	for _, f := range p.SortableFields {
		if p.Sort == f {
			return strings.TrimPrefix(toSnakeCase(p.Sort), "-")
		}
	}

	panic("sort parameter " + p.Sort + " not allowed")
}

// Converts potential camelCase values to snake case.
func toSnakeCase(value string) string {
	snake := firstCapRX.ReplaceAllString(value, "${1}_${2}")
	snake = allCapRX.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

var (
	firstCapRX = regexp.MustCompile("(.)([A-Z][a-z]+)")
	allCapRX   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// SortDirection parses the Sort value and returns the corresponding SQL direction.
func (p QueryParams) SortDirection() string {
	if strings.HasPrefix(p.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

// Validate QueryParams values. Should be called before using values.
func (p QueryParams) Validate(v *validator.Validator) {
	v.Check(p.Offset <= 10_000_000, "offset", "must be a maximum of 10 million")
	v.Check(p.Limit > 0, "limit", "must be greater than zero")
	v.Check(p.Limit <= 100, "limit", "must be a maximum of 100")
	v.Check(validator.In(p.Sort, p.SortableFields...), "sort", "invalid sort value")
}

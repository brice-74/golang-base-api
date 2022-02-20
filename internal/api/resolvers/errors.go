package resolvers

import (
	"fmt"

	"github.com/brice-74/golang-base-api/pkg/validator"
)

const (
	errServer             = "ServerError"
	errInvalidCredentials = "InvalidCredentials"
	errUnauthorized       = "Unauthorized"
	errValidator          = "ValidatorError"
	errDatabaseOperation  = "DatabaseOperationError"
	errNotAuthenticated   = "NotAuthenticatedError"
	errNotFound           = "NotFoundError"
)

var (
	resolverErrUnauthorized = resolverError{
		Code:       errUnauthorized,
		StatusCode: 401,
		Message:    "Access Unauthorized",
	}
)

// resolverError is a general resolver error helper.
type resolverError struct {
	Code       string `json:"code"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (e resolverError) Error() string {
	return fmt.Sprintf("error [%s]: %s", e.Code, e.Message)
}

func (e resolverError) Extensions() map[string]interface{} {
	return map[string]interface{}{
		"statusCode": e.StatusCode,
		"code":       e.Code,
		"message":    e.Message,
	}
}

// validatorError can be used to returns errors from the validator package.
type validatorError struct {
	Errors validator.Errors `json:"errors"`
}

func (e validatorError) Error() string {
	return fmt.Sprintf("validation error [%s]", errValidator)
}

func (e validatorError) Extensions() map[string]interface{} {
	return map[string]interface{}{
		"code":   errValidator,
		"errors": e.Errors,
	}
}

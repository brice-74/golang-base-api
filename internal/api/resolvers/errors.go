package resolvers

import (
	"fmt"

	"github.com/brice-74/golang-base-api/pkg/validator"
)

const (
	errInvalidCredentials = "InvalidCredentials"
	errUnauthorized       = "Unauthorized"
	errValidator          = "ValidatorError"
	errDatabaseOperation  = "DatabaseOperationError"
	errNotFound           = "NotFoundError"
)

func resolverErrNotFound(err error) resolverError {
	msg := "Ressource coul not be found"
	if err != nil {
		msg = err.Error()
	}

	return resolverError{
		Code:       errNotFound,
		StatusCode: 404,
		Message:    msg,
	}
}

func resolverErrInvalidCredentials(err error) resolverError {
	msg := "Invalid credentials"
	if err != nil {
		msg = err.Error()
	}

	return resolverError{
		Code:       errInvalidCredentials,
		StatusCode: 403,
		Message:    msg,
	}
}

func resolverErrUnauthorized(err error) resolverError {
	msg := "Unauthorized access"
	if err != nil {
		msg = err.Error()
	}

	return resolverError{
		Code:       errUnauthorized,
		StatusCode: 401,
		Message:    msg,
	}
}

func resolverErrDatabaseOperation(err error) resolverError {
	msg := "Database operation error"
	if err != nil {
		msg = err.Error()
	}

	return resolverError{
		Code:       errDatabaseOperation,
		StatusCode: 500,
		Message:    msg,
	}
}

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
		"statusCode": 422,
		"code":       errValidator,
		"errors":     e.Errors,
	}
}

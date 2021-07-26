package fault

import (
	"database/sql"
	"errors"

	"github.com/go-playground/validator/v10"
	"regeet.io/api/internal/conf"
)

var (
	// ErrMissingCookie
	ErrMissingCookie = errors.New("there is no such cookie")

	// ErrMissingJwtToken
	ErrMissingJwtToken = errors.New("the jwt token is missing or malformed")

	// ErrInvalidJwtToken
	ErrInvalidJwtToken = errors.New("the jwt token is invalid or expired")
)

// ValidationError
type ValidationError struct {
	Namespace string `json:"namespace"`
	Tag       string `json:"tag"`
	Message   string `json:"message"`
}

// UserInputError
type UserInputError struct {
	Errors map[string]ValidationError
}

// Error
func (e UserInputError) Error() string {
	return "" // TODO
}

// AddError
func (e UserInputError) AddError(namespace string, tag string, message string) {
	e.Errors[namespace] = ValidationError{
		Namespace: namespace,
		Tag:       tag,
		Message:   message,
	}
}

// UserInputErrorFrom
func UserInputErrorFrom(err error) UserInputError {
	ret := UserInputError{
		Errors: map[string]ValidationError{},
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, err := range errs {
			ret.AddError(
				err.Namespace(),
				err.Tag(),
				err.Translate(conf.EnTrans),
			)
		}
	}

	return ret
}

var (
	ErrUnauthenticated  = errors.New("you need to authenticate to be able to perform this operation")
	ErrForbidden        = errors.New("you don't have enough permissions to perform this operation")
	ErrResourceNotFound = sql.ErrNoRows
	ErrUserInput        = errors.New("the provided input is not correct")
)

// IsForbiddenError
func IsForbiddenError(err error) bool {
	return err == ErrForbidden
}

// IsUnauthenticatedError
func IsUnauthenticatedError(err error) bool {
	return err == ErrUnauthenticated
}

// IsResourceNotFoundError
func IsResourceNotFoundError(err error) bool {
	return err == ErrResourceNotFound
}

// IsNonResourceNotFoundError
func IsNonResourceNotFoundError(err error) bool {
	return err != nil && err != ErrResourceNotFound
}

// IsUserInputError
func IsUserInputError(err error) bool {
	// ErrUserInput
	if err == ErrUserInput {
		return true
	}

	// Structured
	if _, ok := err.(*UserInputError); ok {
		return true
	}

	return false
}

// IsNonUserInputError
func IsNonUserInputError(err error) bool {
	return err != nil && !IsUserInputError(err)
}

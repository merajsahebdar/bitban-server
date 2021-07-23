package common

import (
	"database/sql"
	"errors"

	"github.com/go-playground/validator/v10"
)

// ErrMissingJwtToken
var ErrMissingJwtToken = errors.New("the jwt token is missing or malformed")

// ErrInvalidJwtToken
var ErrInvalidJwtToken = errors.New("the jwt token is invalid or expired")

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
				err.Translate(EnTrans),
			)
		}
	}

	return ret
}

var (
	ErrResourceNotFound = sql.ErrNoRows
	ErrUserInput        = errors.New("the provided input is not correct")
)

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

package resolver

import (
	"github.com/vektah/gqlparser/v2/gqlerror"
	"regeet.io/api/internal/fault"
)

// ErrorExtensions
type ErrorExtensions = map[string]interface{}

// newNotFound
func newNotFoundErrorExtensions() ErrorExtensions {
	return ErrorExtensions{
		"code": "NOT_FOUND",
	}
}

// newAuthenticationErrorExtensions
func newAuthenticationErrorExtensions() ErrorExtensions {
	return ErrorExtensions{
		"code": "UNAUTHENTICATED",
	}
}

// newForbiddenErrorExtensions
func newForbiddenErrorExtensions() ErrorExtensions {
	return ErrorExtensions{
		"code": "FORBIDDEN",
	}
}

// newUserInputErrorExtensions
func newUserInputErrorExtensions() ErrorExtensions {
	return ErrorExtensions{
		"code": "BAD_USER_INPUT",
	}
}

// newInternalServerErrorExtensions
func newInternalServerErrorExtensions() ErrorExtensions {
	return ErrorExtensions{
		"code": "INTERNAL_SERVER_ERROR",
	}
}

// NotFoundErrorFrom
func NotFoundErrorFrom(err error) *gqlerror.Error {
	return &gqlerror.Error{
		Message:    "no such resource found",
		Extensions: newNotFoundErrorExtensions(),
	}
}

// AuthenticationErrorFrom
func AuthenticationErrorFrom(err error) *gqlerror.Error {
	return &gqlerror.Error{
		Message:    "you need to authenticate to be able to perform this operation",
		Extensions: newAuthenticationErrorExtensions(),
	}
}

// ForbiddenErrorFrom
func ForbiddenErrorFrom(err error) *gqlerror.Error {
	return &gqlerror.Error{
		Message:    "you don't have enough permissions to perform this operation",
		Extensions: newForbiddenErrorExtensions(),
	}
}

// UserInputErrorFrom
func UserInputErrorFrom(err error) *gqlerror.Error {
	ext := newUserInputErrorExtensions()

	if errUserInput, ok := err.(fault.UserInputError); ok {
		ext["errors"] = errUserInput.Errors
	}

	return &gqlerror.Error{
		Message:    "the provided input is not valid",
		Extensions: ext,
	}
}

// InternalServerErrorFrom
func InternalServerErrorFrom(err error) *gqlerror.Error {
	return &gqlerror.Error{
		Message:    "got an internal server error, try again later",
		Extensions: newInternalServerErrorExtensions(),
	}
}

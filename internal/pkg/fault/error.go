/*
 * Copyright 2021 Meraj Sahebdar
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package fault

import (
	"database/sql"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"regeet.io/api/internal/cfg"
)

// ValidationError
type ValidationError struct {
	Namespace string `json:"namespace"`
	Tag       string `json:"tag"`
	Message   string `json:"message"`
}

// UserInputError
type UserInputError struct {
	error
	Errors map[string]ValidationError
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
		error:  err,
		Errors: map[string]ValidationError{},
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, err := range errs {
			ret.AddError(
				err.Namespace(),
				err.Tag(),
				err.Translate(cfg.EnTrans),
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
	if _, ok := err.(UserInputError); ok {
		return true
	}

	return false
}

// IsNonUserInputError
func IsNonUserInputError(err error) bool {
	return err != nil && !IsUserInputError(err)
}

// IsPqUniqueViolationError
func IsPqUniqueViolationError(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code.Name() == "unique_violation"
	}

	return false
}

// IsNonPqUniqueViolationError
func IsNonPqUniqueViolationError(err error) bool {
	return err != nil && !IsPqUniqueViolationError(err)
}

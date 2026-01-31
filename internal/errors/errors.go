// Package errors provides a structured error type and functions to work with it.
// Adapted from gitlab.com/loungeup/go-loungeup (MIT, LoungeUp) and Ben Johnson.
// Reference: https://www.gobeyond.dev/failure-is-your-domain/
package errors

import (
	"database/sql"
	"errors"
	"fmt"
)

// As and Is re-export standard errors package for convenience.
var (
	As = errors.As
	Is = errors.Is
)

// Pre-defined error codes.
const (
	CodeConflict = "conflict"
	CodeInternal = "internal"
	CodeInvalid  = "invalid"
	CodeNotFound = "notFound"

	errorMessageConflict = "Conflict"
	errorMessageInternal = "An internal error has occurred. Please contact technical support."
	errorMessageInvalid  = "Invalid"
	errorMessageNotFound = "Not found"
)

// Error represents a structured application error.
type Error struct {
	// Code is machine-readable.
	Code string

	// Message is human-readable.
	Message string

	// Operation that caused the error.
	Operation string

	// UnderlyingError that caused this error, if any.
	UnderlyingError error
}

func (e *Error) Error() string {
	if e.UnderlyingError != nil {
		if e.Operation != "" {
			return fmt.Sprintf("%s: %v", e.Operation, e.UnderlyingError)
		}

		return e.UnderlyingError.Error()
	}

	result := ""
	if e.Code != "" {
		result += "<" + e.Code + "> "
	}
	if e.Message != "" {
		result += e.Message
	}
	if e.Operation != "" {
		result = e.Operation + ": " + result
	}

	return result
}

func (e *Error) Unwrap() error {
	return e.UnderlyingError
}

func (e *Error) defaultMessage() string {
	switch e.Code {
	case CodeConflict:
		return errorMessageConflict
	case CodeInvalid:
		return errorMessageInvalid
	case CodeNotFound:
		return errorMessageNotFound
	default:
		return errorMessageInternal
	}
}

// ErrorCode returns the machine-readable code for err.
func ErrorCode(err error) string {
	if err == nil {
		return ""
	}

	if e, ok := err.(*Error); ok && e.Code != "" {
		return e.Code
	}
	if e, ok := err.(*Error); ok && e.UnderlyingError != nil {
		return ErrorCode(e.UnderlyingError)
	}

	return CodeInternal
}

// ErrorMessage returns the human-readable message for err.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	switch e, ok := err.(*Error); {
	case ok && e.Message != "":
		return e.Message
	case ok && e.Code != "":
		return e.defaultMessage()
	case ok && e.UnderlyingError != nil:
		return ErrorMessage(e.UnderlyingError)
	default:
		return errorMessageInternal
	}
}

// MapSQLError maps SQL errors to structured errors (SQLite-specific).
func MapSQLError(err error) error {
	if err == nil {
		return nil
	}

	if err == sql.ErrNoRows {
		return &Error{
			Code:            CodeNotFound,
			UnderlyingError: err,
		}
	}

	return &Error{
		Code:            CodeInternal,
		UnderlyingError: err,
	}
}

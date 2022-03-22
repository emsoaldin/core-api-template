package apperror

import "fmt"

// Error is error structure specific to web based projects
type Error struct {
	*stack
	cause error
	code  string
	err   error
}

// Error interface implementation
func (h *Error) Error() string { return h.err.Error() }

// Code returns error code
func (h *Error) Code() string { return h.code }

// Cause returns original error
func (h *Error) Cause() error {
	return h.cause
}

// Unwrap implements erooro unwrapping
func (h *Error) Unwrap() error {
	return h.cause
}

// Trace returns error stack trace string
func (h *Error) Trace() string {
	return fmt.Sprintf("%+v", h)
}

// New returns an error with the supplied message and code.
//
// function also records the stack trace at the point it was called.
func New(code string, err error, cause ...error) error {

	e := &Error{
		code:  code,
		err:   err,
		stack: callers(),
	}

	if len(cause) > 0 {
		e.cause = cause[0]
	}

	return e
}

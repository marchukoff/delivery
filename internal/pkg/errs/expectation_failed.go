package errs

import (
	"errors"
	"fmt"
	"strings"
)

var ErrExpectationFailed = errors.New("expectation failed")

type ExpectationFailedError struct {
	ParamName string
	Got       any
	Want      []any
	Cause     error
}

func NewExpectationFailedErrorWithCause(paramName string, got any, cause error, want ...any) *ExpectationFailedError {
	return &ExpectationFailedError{
		ParamName: paramName,
		Got:       got,
		Want:      want,
		Cause:     cause,
	}
}

func NewExpectationFailedError(paramName string, got any, want ...any) *ExpectationFailedError {
	return &ExpectationFailedError{
		ParamName: paramName,
		Got:       got,
		Want:      want,
	}
}

func (e *ExpectationFailedError) Error() string {
	want := make([]string, 0, len(e.Want))
	for i := range e.Want {
		want = append(want, sanitize(e.Want[i]))
	}

	if e.Cause != nil {
		return fmt.Sprintf("%s: %s got: %v, want: %s (cause: %v)",
			ErrExpectationFailed, sanitize(e.Got), e.ParamName, strings.Join(want, ", "), e.Cause)
	}
	return fmt.Sprintf("%s: %s got: %v, want: %s",
		ErrExpectationFailed, sanitize(e.Got), e.ParamName, strings.Join(want, ", "))
}

func (e *ExpectationFailedError) Unwrap() error {
	return ErrExpectationFailed
}

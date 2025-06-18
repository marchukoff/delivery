package errs

import (
	"errors"
	"fmt"
)

var ErrObjectNotFound = errors.New("object not found")

type ObjectNotFoundError struct {
	ParamName string
	ID        any
	Cause     error
}

func NewObjectNotFoundErrorWithCause(paramName string, ID any, cause error) *ObjectNotFoundError {
	return &ObjectNotFoundError{
		ParamName: paramName,
		ID:        ID,
		Cause:     cause,
	}
}

func NewObjectNotFoundError(paramName string, ID any) *ObjectNotFoundError {
	return &ObjectNotFoundError{
		ParamName: paramName,
		ID:        ID,
	}
}

func (e *ObjectNotFoundError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: param is: %s, ID is: %v (cause: %v)",
			ErrObjectNotFound, e.ParamName, e.ID, e.Cause)
	}
	return fmt.Sprintf("%s: param is: %s, ID is: %v", ErrObjectNotFound, e.ParamName, e.ID)
}

func (e *ObjectNotFoundError) Unwrap() error {
	return ErrObjectNotFound
}

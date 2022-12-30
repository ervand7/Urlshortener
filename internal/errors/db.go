package errors

import (
	"errors"
	"fmt"
)

// ShortAlreadyExistsError error.
type ShortAlreadyExistsError struct {
	Err error
}

// NewShortAlreadyExistsError constructor of ShortAlreadyExistsError.
func NewShortAlreadyExistsError(short string) error {
	return &ShortAlreadyExistsError{
		Err: fmt.Errorf(`%w`, errors.New(short)),
	}
}

// Unwrap for implementation of Wrapper.
func (s ShortAlreadyExistsError) Unwrap() error {
	return s.Err
}

// Error for implementation of Error.
func (s ShortAlreadyExistsError) Error() string {
	return fmt.Sprintf("%s", s.Err)
}

// URLNotActiveError error.
type URLNotActiveError struct {
	Err error
}

// NewURLNotActiveError constructor of URLNotActiveError.
func NewURLNotActiveError(short string) error {
	return &URLNotActiveError{
		Err: fmt.Errorf(`%w`, errors.New(short+" is not active")),
	}
}

// Unwrap for implementation of Wrapper.
func (u URLNotActiveError) Unwrap() error {
	return u.Err
}

// Error for implementation of Error.
func (u URLNotActiveError) Error() string {
	return fmt.Sprintf("%s", u.Err)
}

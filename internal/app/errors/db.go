package errors

import (
	"errors"
	"fmt"
)

// ShortAlreadyExistsError error
type ShortAlreadyExistsError struct {
	Err error
}

func NewShortAlreadyExistsError(short string) error {
	return &ShortAlreadyExistsError{
		Err: fmt.Errorf(`%w`, errors.New(short)),
	}
}

func (s ShortAlreadyExistsError) Unwrap() error {
	return s.Err
}

func (s ShortAlreadyExistsError) Error() string {
	return fmt.Sprintf("%s", s.Err)
}

// URLNotActiveError error
type URLNotActiveError struct {
	Err error
}

func NewURLNotActiveError(short string) error {
	return &URLNotActiveError{
		Err: fmt.Errorf(`%w`, errors.New(short+" is not active")),
	}
}

func (u URLNotActiveError) Unwrap() error {
	return u.Err
}

func (u URLNotActiveError) Error() string {
	return fmt.Sprintf("%s", u.Err)
}

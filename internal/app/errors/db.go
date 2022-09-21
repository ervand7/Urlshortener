package errors

import (
	"errors"
	"fmt"
)

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

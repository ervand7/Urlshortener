package apperrors

import (
	"errors"
	"fmt"
)

type ShortAlreadyExistsError struct {
	Err error
}

func NewShortAlreadyExistsError(origin string) error {
	return &ShortAlreadyExistsError{
		Err: fmt.Errorf(`%w`, errors.New("url was already shortened: "+origin)),
	}
}

func (s ShortAlreadyExistsError) Unwrap() error {
	return s.Err
}

func (s ShortAlreadyExistsError) Error() string {
	return fmt.Sprintf("%s", s.Err)
}

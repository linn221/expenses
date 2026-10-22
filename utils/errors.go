package utils

import (
	"errors"
	"fmt"
)

var ErrBadRequest = errors.New("bad request")
var ErrInternalServerError = errors.New("internal server error")

func NewBadRequest(msg string) error {
	return fmt.Errorf("%w: %s", ErrBadRequest, msg)
}

func NewInternalServerError(msg string) error {
	return fmt.Errorf("%w: %s", ErrInternalServerError, msg)
}

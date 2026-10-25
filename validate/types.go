package validate

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type FilterFunc func(db *gorm.DB)

type Action int

const (
	ActionCreate Action = iota
	ActionUpdate
	ActionDelete
	ActionActiveStatusChange
)

var ErrClient = errors.New("validation failed")
var ErrServer = errors.New("internal server error")

func clientErrorf(placeholder string, args ...any) error {
	message := placeholder
	if len(args) > 0 {
		message = fmt.Sprintf(placeholder, args...)
	}
	return fmt.Errorf("%w: %s", ErrClient, message)
}

func clientError(message string) error {
	return fmt.Errorf("%w: %s", ErrClient, errors.New(message))
}

func serverError(err error) error {
	return fmt.Errorf("%w: %s", ErrServer, err.Error())
}

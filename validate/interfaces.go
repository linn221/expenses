package validate

import (
	"context"

	"github.com/linn221/expenses/app"
)

type Rule interface {
	Skipped() bool
	Validate(ctx context.Context, c *app.Container) error
}

type Input interface {
	ValidationRules(*ValidationContext)
}

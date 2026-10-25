package models

import "github.com/linn221/expenses/app"

type ModelConverter[M any] interface {
	ToModel(*app.Session) M
}

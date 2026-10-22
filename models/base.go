package models

import "github.com/linn221/expenses/app"

type NewModel[M any] interface {
	ToModel(*app.Session) M
}

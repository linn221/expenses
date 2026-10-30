package controllers

import (
	"net/http"

	"github.com/linn221/expenses/app"
	"github.com/linn221/expenses/dtos"
	"gorm.io/gorm"
)

/*
for Services
*/

type Getter[T any] interface {
	Get(*app.Session) (*T, error)
}
type Lister[T, F any] interface {
	List(*app.Session, F) ([]T, error)
}

// ApplyFilters returns the model query with list conditions already applied.
type ApplyFilters[F any] interface {
	ApplyFilters(*app.Session, F) *gorm.DB
}
type Creater[T any] interface {
	Create(*app.Session, T) (*T, error)
}
type Updater[T any] interface {
	Getter[T]
	Update(*app.Session, *T, *T) (*T, error)
}
type Deleter[T any] interface {
	Getter[T]
	Delete(*app.Session, *T) (*T, error)
}

/*
for Filters
*/
type QueryParser[T any] interface {
	Parse(*http.Request) (*T, error)
}

// for services
type ModelResponseGenerator[T any] interface {
	HttpResponse() dtos.HttpResponseFromModel[T]
}
type ModelResponseGeneratorSelf interface {
	HttpResponse() any
}

// for model service
type CollectionResponseGenerator[T any] interface {
	HttpResponseCollection() dtos.HttpResponseFromModelCollection[T]
}

package routing

import "net/http"

type ResourceController interface {
	Create() http.HandlerFunc
	Update() http.HandlerFunc
	Delete() http.HandlerFunc
	Get() http.HandlerFunc
}

type PatchController interface {
	Patch() http.HandlerFunc
}

type ListController interface {
	List() http.HandlerFunc
}

type PaginateController interface {
	Paginate() http.HandlerFunc
}

type IsActiveToggleController interface {
	ToggleIsActive() http.HandlerFunc
}

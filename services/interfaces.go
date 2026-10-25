package services

type Togglable interface {
	GetIsActive() bool
}

type UpdatableModel[T any] interface {
	UpdateValues(input T) map[string]any
}

package dtos

/*
for Response
the model's service does not have to have Resource method
but the model can have Resource method itself as a fallback
*/

type HttpResponseFromModel[T any] func(*T) any
type HttpResponseFromModelCollection[T any] func([]T) any

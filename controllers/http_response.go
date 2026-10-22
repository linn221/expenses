package controllers

type HasResource interface {
	Resource(any) (any, error)
}
type HasResourceSelf interface {
	Resource() any
}

// if the service has Resource(T) (any, error) method, prefer it for responses
// if not find Resource() any method on the model
// only when none exists, return the model
// which may call db
func getResource(service any, i any) (any, error) {
	res, ok := service.(HasResource)
	if ok {
		return res.Resource(i)
	}
	resSelf, ok := i.(HasResourceSelf)
	if ok {
		return resSelf.Resource(), nil
	}
	return i, nil
}

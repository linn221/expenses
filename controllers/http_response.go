package controllers

/*
the model's service does not have to have Resource method
but the model can have Resource method itself as a fallback
*/
type HasResource interface {
	Resource(any) any
}
type HasResourceSelf interface {
	Resource() any
}

/*
in order to use CollectionResource for customizable http responses,
a model has to have a collection resource type that implements AddInstances([]model) error method (ie, CategoryCollection for Category)
the model's servie must have CollectionResource() method that returns the newly instantiated CollectionResource type (ie, CollectionResource() method that returns CategoryCollection type on CategoryService)
*/
// for collectionModel type
type CollectionResource[T any] interface {
	AddInstances([]T)
}

// for model service
type HasCollectionResource[T any] interface {
	CollectionResource() CollectionResource[T]
}

// if the service has Resource(T) (any, error) method, prefer it for responses
// if not find Resource() any method on the model
// only when none exists, return the model
// which may call db
func getResource(service any, i any) any {
	res, ok := service.(HasResource)
	if ok {
		return res.Resource(i)
	}
	resSelf, ok := i.(HasResourceSelf)
	if ok {
		return resSelf.Resource()
	}
	return i
}

func getCollectionResource[T any](service any, results []T) any {
	collectionResService, ok := service.(HasCollectionResource[T])
	if ok {
		collectionResource := collectionResService.CollectionResource()
		collectionResource.AddInstances(results)
		return collectionResource
	}
	return results
}

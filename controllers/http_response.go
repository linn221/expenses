package controllers

// if the service has Resource(T) (any, error) method, prefer it for responses
// if not find Resource() any method on the model
// only when none exists, return the model
// which may call db
func getResource[T any](service any, i *T) any {
	responseGenerator, ok := service.(ModelResponseGenerator[T])
	if ok {
		return responseGenerator.HttpResponse()(i)
	}
	resSelf, ok := any(i).(ModelResponseGeneratorSelf)
	if ok {
		return resSelf.HttpResponse()
	}
	return i
}

func getCollectionResource[T any](service any, results []T) any {
	responseGenerator, ok := service.(CollectionResponseGenerator[T])
	if ok {
		return responseGenerator.HttpResponseCollection()(results)
	}
	return results
}

package dtos

func addToCollection[Model, Instance any](results []Model, convert func(Model) Instance) []Instance {
	instances := make([]Instance, 0, len(results))
	for _, r := range results {
		instances = append(instances, convert(r))
	}
	return instances
}

package dtos

import (
	"net/http"
	"strconv"

	"github.com/linn221/expenses/utils"
)

func addToCollection[Model, Instance any](results []Model, convert func(Model) Instance) []Instance {
	instances := make([]Instance, 0, len(results))
	for _, r := range results {
		instances = append(instances, convert(r))
	}
	return instances
}

func walletIDFromPath(r *http.Request) (int, error) {
	idValue := r.PathValue("walletId")
	if idValue == "" {
		return 0, utils.NewBadRequest("resource id not found")
	}
	walletID, err := strconv.Atoi(idValue)
	if err != nil {
		return 0, utils.NewBadRequest("resource id not found")
	}
	return walletID, nil
}

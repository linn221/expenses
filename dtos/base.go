package dtos

import (
	"net/http"
	"strconv"

	"github.com/linn221/expenses/utils"
	"github.com/shopspring/decimal"
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

// NaturalDecimal is a decimal strictly greater than zero.
type NaturalDecimal struct {
	decimal.Decimal
}

func (n *NaturalDecimal) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return utils.NewBadRequest("value must be greater than zero")
	}
	if err := n.Decimal.UnmarshalJSON(data); err != nil {
		return utils.NewBadRequest("value is not a valid decimal")
	}
	if !n.GreaterThan(decimal.Zero) {
		return utils.NewBadRequest("value must be greater than zero")
	}
	return nil
}

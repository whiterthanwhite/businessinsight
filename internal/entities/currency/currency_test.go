package currency

import (
	"encoding/json"
	"testing"
)

func TestCurrency(t *testing.T) {
	currencies := []*Currency{
		&Currency{
			Code:        "GEL",
			Description: "Georgian currency",
		},
		&Currency{
			Code:        "RUB",
			Description: "Russian currency",
		},
	}
	currenciesJSON, err := json.Marshal(currencies)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(currenciesJSON))
}

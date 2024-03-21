package currency

import "encoding/json"

type Currency struct {
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
}

func ParseJSON(currenciesJSON []byte) ([]Currency, error) {
	var currencies []Currency
	if err := json.Unmarshal(currenciesJSON, &currencies); err != nil {
		return nil, err
	}

	return currencies, nil
}

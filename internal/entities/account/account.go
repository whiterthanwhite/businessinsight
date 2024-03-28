package account

import "encoding/json"

type Account struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	CurrencyCode string `json:"currency_code"`
}

func ParseJSON(dataJSON []byte) ([]Account, error) {
	var accounts []Account
	err := json.Unmarshal(dataJSON, &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

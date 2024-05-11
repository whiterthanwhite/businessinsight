package operation

import (
	"encoding/json"
	"time"

	"github.com/whiterthanwhite/businessinsight/internal/entities/operation_type"
)

type Operation struct {
	EntryNo       int                          `json:"entryNo"`
	DateTime      time.Time                    `json:"dateTime"`
	CreationDate  time.Time                    `json:"creation_date,omitempty"`
	CreationTime  time.Time                    `json:"cretion_time,omitempty"`
	Type          operation_type.OperationType `json:"type"`
	Amount        float64                      `json:"amount"`
	SourceId      int                          `json:"sourceId"`
	CurrencyCode  string                       `json:"currencyCode"`
	CategoryId    int                          `json:"categoryId"`
	TransactionNo int                          `json:"transactionNo"`
	Description   string                       `json:"description"`
}

type operationJSON struct {
	EntryNo       int                          `json:"entryNo"`
	DateTime      string                       `json:"dateTime"`
	CreationDate  time.Time                    `json:"creation_date,omitempty"`
	CreationTime  time.Time                    `json:"cretion_time,omitempty"`
	Type          operation_type.OperationType `json:"type"`
	Amount        float64                      `json:"amount"`
	SourceId      int                          `json:"sourceId"`
	CurrencyCode  string                       `json:"currencyCode"`
	CategoryId    int                          `json:"categoryId"`
	TransactionNo int                          `json:"transactionNo"`
	Description   string                       `json:"description"`
}

func (o *Operation) MarshalJSON() ([]byte, error) {
	oJSON := operationJSON{
		EntryNo:       o.EntryNo,
		DateTime:      o.DateTime.Format("2006-01-02T15:04"),
		Type:          o.Type,
		Amount:        o.Amount,
		SourceId:      o.SourceId,
		CurrencyCode:  o.CurrencyCode,
		CategoryId:    o.CategoryId,
		TransactionNo: o.TransactionNo,
		Description:   o.Description,
	}
	body, err := json.Marshal(&oJSON)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (o *Operation) UnmarshalJSON(body []byte) error {
	var oJSON operationJSON
	var err error
	if err = json.Unmarshal(body, &oJSON); err != nil {
		return err
	}
	o.EntryNo = oJSON.EntryNo
	if o.DateTime, err = time.Parse("2006-01-02T15:04", oJSON.DateTime); err != nil {
		return err
	}
	o.Type = oJSON.Type
	o.Amount = oJSON.Amount
	o.SourceId = oJSON.SourceId
	o.CurrencyCode = oJSON.CurrencyCode
	o.CategoryId = oJSON.CategoryId
	o.TransactionNo = oJSON.TransactionNo
	o.Description = oJSON.Description
	return nil
}

/*
func ParseJSON(body []byte) ([]Operation, error) {
	var operations []Operation
	err := json.Unmarshal(body, &operations)
	if err != nil {
		return nil, err
	}
	return operations, nil
}
*/

func (o *Operation) Compare(with *Operation) bool {
	if o.DateTime.Compare(with.DateTime) != 0 ||
		o.Type != with.Type ||
		o.Amount != with.Amount ||
		o.SourceId != with.SourceId ||
		o.CurrencyCode != with.CurrencyCode ||
		o.CategoryId != with.CategoryId ||
		o.TransactionNo != with.TransactionNo ||
		o.Description != with.Description {

		return false
	}

	return true
}

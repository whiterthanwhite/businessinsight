package operation

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/whiterthanwhite/businessinsight/internal/entities/operation_type"
)

// [{"entryNo":0,"dateTime":"2024-04-07T00:36","type":"Income","amount":0,"sourceId":0,"currencyCode":"","categoryId":0,"transactionNo":0,"description":""}]

func TestMarshalOperation(t *testing.T) {
	operations := []Operation{
		{
			EntryNo:       1,
			DateTime:      time.Now(),
			Type:          operation_type.Expense,
			Amount:        100.00,
			SourceId:      1,
			CurrencyCode:  "GEL",
			CategoryId:    1,
			TransactionNo: 0,
			Description:   "Test operation 1",
		},
		{
			EntryNo:       2,
			DateTime:      time.Now(),
			Type:          operation_type.Expense,
			Amount:        101.00,
			SourceId:      2,
			CurrencyCode:  "GEL",
			CategoryId:    2,
			TransactionNo: 0,
			Description:   "Test operation 2",
		},
	}
	body, err := json.Marshal(operations)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(body))
}

func TestUnmarshalOperation(t *testing.T) {
	requestBody := `[{"entryNo":0,"dateTime":"2024-04-07T00:36","type":"Income","amount":0,"sourceId":0,"currencyCode":"","categoryId":0,"transactionNo":0,"description":""}]`
	var operations []*Operation
	if err := json.Unmarshal([]byte(requestBody), &operations); err != nil {
		t.Fatal(err.Error())
	}
	t.Log(operations[0])
}

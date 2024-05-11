package db

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/whiterthanwhite/businessinsight/internal/entities/account"
	"github.com/whiterthanwhite/businessinsight/internal/entities/category"
	"github.com/whiterthanwhite/businessinsight/internal/entities/currency"
	"github.com/whiterthanwhite/businessinsight/internal/entities/operation"
	"github.com/whiterthanwhite/businessinsight/internal/entities/operation_type"
)

var (
	dbConnStr = flag.String("db", "", "")
)

func TestOperationBefore(t *testing.T) {
	flag.Parse()
	_, err := Connect(context.TODO(), *dbConnStr)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestPrepareCurrencies(t *testing.T) {
	conn, err := GetInstance()
	if err != nil {
		t.Fatal(err.Error())
	}
	curr := &currency.Currency{
		Code:        "RUB",
		Description: "Russian currency",
	}
	err = conn.InsertCurrency(context.TODO(), curr)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestPrepareAccounts(t *testing.T) {
	conn, err := GetInstance()
	if err != nil {
		t.Fatal(err.Error())
	}
	acc := &account.Account{
		Name:         "Test account",
		CurrencyCode: "RUB",
	}
	err = conn.InsertAccount(context.TODO(), acc)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestPrepareCategories(t *testing.T) {
	conn, err := GetInstance()
	if err != nil {
		t.Fatal(err.Error())
	}
	cat := &category.Category{
		Type:        operation_type.Expense,
		Name:        "Test Expense",
		Description: "Test expense cateogory",
	}
	err = conn.InsertCategory(context.TODO(), cat)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestOperationInsert(t *testing.T) {
	conn, err := GetInstance()
	if err != nil {
		t.Fatal(err.Error())
	}

	dateTime := time.Now()
	operation := &operation.Operation{
		EntryNo:      0,
		DateTime:     dateTime,
		Type:         operation_type.Expense,
		Amount:       -100.001001001,
		SourceId:     2,
		CurrencyCode: "RUB",
		CategoryId:   1,
		Description:  "Hello. This is test operation",
		CreationDate: dateTime,
		CreationTime: dateTime,
	}

	err = conn.InsertOperation(context.TODO(), operation)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestOperationUpdate(t *testing.T) {
	conn, err := GetInstance()
	if err != nil {
		t.Fatal(err.Error())
	}

	operations, err := conn.GetOperations(context.TODO())
	if err != nil {
		t.Fatal(err.Error())
	}

	for _, uOperation := range operations {
		newTime := time.Now()
		uOperation.CreationDate = newTime
		uOperation.CreationTime = newTime
		o, err := conn.GetOperation(context.TODO(), &uOperation)
		if err != nil {
			t.Error(err.Error())
			continue
		}
		if *o != uOperation {
			err = conn.UpdateOperation(context.TODO(), &uOperation)
			if err != nil {
				t.Error(err.Error())
			}
		}
	}
}

func TestOperationDelete(t *testing.T) {
	conn, err := GetInstance()
	if err != nil {
		t.Fatal(err.Error())
	}

	operations, err := conn.GetOperations(context.TODO())
	if err != nil {
		t.Fatal(err.Error())
	}

	for _, delOperation := range operations {
		err = conn.DeleteOperation(context.TODO(), &delOperation)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestOperationAfter(t *testing.T) {
	conn, err := GetInstance()
	if err != nil {
		t.Fatal(err.Error())
	}

	err = conn.Close(context.TODO())
	if err != nil {
		t.Fatal(err.Error())
	}
}

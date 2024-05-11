package helper

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/whiterthanwhite/businessinsight/internal/db"
)

func ExportDataToCSV(parentCtx context.Context) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	if err := exportAccountsToCSV(ctx); err != nil {
		return err
	}
	if err := exportCurrenciesCSV(ctx); err != nil {
		return err
	}
	if err := exportCategoryCSV(ctx); err != nil {
		return err
	}
	if err := exportOperationToCSV(ctx); err != nil {
		return err
	}

	log.Println("export finished")
	return nil
}

func exportAccountsToCSV(parentCtx context.Context) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	conn, err := db.GetInstance()
	if err != nil {
		return err
	}

	accounts, err := conn.GetAccounts(ctx)
	if err != nil {
		return err
	}

	values := make([][]string, len(accounts))
	for i, account := range accounts {
		values[i] = []string{fmt.Sprint(account.Id), account.CurrencyCode, account.Name}
	}

	err = exportToCSV("accounts", "accounts exported", values)
	if err != nil {
		return err
	}

	return nil
}

func exportCurrenciesCSV(parentCtx context.Context) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	conn, err := db.GetInstance()
	if err != nil {
		return err
	}

	currencies, err := conn.GetCurrencies(ctx)
	if err != nil {
		return err
	}

	values := make([][]string, len(currencies))
	for i, curr := range currencies {
		values[i] = []string{fmt.Sprint(curr.Code), curr.Description}
	}

	err = exportToCSV("currencies", "currencies exported", values)
	if err != nil {
		return err
	}
	return nil
}

func exportCategoryCSV(parentCtx context.Context) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	conn, err := db.GetInstance()
	if err != nil {
		return err
	}

	currencies, err := conn.GetCurrencies(ctx)
	if err != nil {
		return err
	}

	values := make([][]string, len(currencies))
	for i, curr := range currencies {
		values[i] = []string{curr.Code, curr.Description}
	}

	err = exportToCSV("categories", "categories exported", values)
	if err != nil {
		return err
	}

	return nil
}

func exportOperationToCSV(parentCtx context.Context) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	conn, err := db.GetInstance()
	if err != nil {
		return err
	}

	currencies, err := conn.GetOperations(ctx)
	if err != nil {
		return err
	}

	values := make([][]string, len(currencies))
	for i, oper := range currencies {
		values[i] = []string{
			fmt.Sprint(oper.EntryNo),
			oper.DateTime.Format(time.DateTime),
			fmt.Sprint(oper.Type),
			fmt.Sprint(oper.Amount),
			fmt.Sprint(oper.SourceId),
			oper.CurrencyCode,
			fmt.Sprint(oper.CategoryId),
			fmt.Sprint(oper.TransactionNo),
			oper.Description,
		}
	}

	err = exportToCSV("operations", "operations exported", values)
	if err != nil {
		return err
	}

	return nil
}

func exportToCSV(fileName, msgExportFinished string, values [][]string) error {
	f, err := os.Create(fmt.Sprintf("%s.csv", fileName))
	if err != nil {
		return err
	}
	defer f.Close()

	scvWriter := csv.NewWriter(f)
	if err := scvWriter.WriteAll(values); err != nil {
		return err
	}

	log.Println(msgExportFinished)
	return nil
}

package db

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/whiterthanwhite/businessinsight/internal/entities/currency"
)

func (d *databaseConnection) GetCurrencies(parentCtx context.Context) ([]currency.Currency, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	rows, err := d.conn.Query(ctx, "SELECT * FROM currency;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var currencies []currency.Currency
	for rows.Next() {
		curr := currency.Currency{}
		err = rows.Scan(&curr.Code, &curr.Description)
		if err != nil {
			return nil, err
		}
		currencies = append(currencies, curr)
	}

	return currencies, nil
}

func (d *databaseConnection) GetCurrency(parentCtx context.Context, newCurrency *currency.Currency) (*currency.Currency, error) {
	log.Println("start function GetCurrency")
	defer log.Println("end function GetCurrency")

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	xCurrency := new(currency.Currency)
	err := d.conn.QueryRow(ctx, "SELECT * FROM currency WHERE code = $1;", newCurrency.Code).Scan(&xCurrency.Code, &xCurrency.Description)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	if err == pgx.ErrNoRows {
		return nil, nil
	}

	return xCurrency, nil
}

func (d *databaseConnection) InsertCurrency(parentCtx context.Context, newCurrency *currency.Currency) error {
	log.Println("start function InsertCurrency")
	defer log.Println("end function InsertCurrency")

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	newCurrency.Code = strings.ToUpper(newCurrency.Code)
	_, err := d.conn.Exec(ctx, "INSERT INTO currency VALUES ($1, $2);", &newCurrency.Code, &newCurrency.Description)
	if err != nil {
		return err
	}

	return nil
}

func (d *databaseConnection) DeleteCurrencies(parentCtx context.Context, currencies []currency.Currency) error {
	if len(currencies) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}

	for _, curr := range currencies {
		_, err := tx.Exec(ctx, "DELETE FROM currency WHERE code = $1;", curr.Code)
		if err != nil {
			tErr := errors.Join(err)
			err = tx.Rollback(ctx)
			tErr = errors.Join(tErr, err)
			return tErr
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (d *databaseConnection) UpdateCurrency(parentCtx context.Context, newCurrency *currency.Currency) error {
	log.Println("start function UpdateCurrency")
	defer log.Println("end function UpdateCurrency")

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, err := d.conn.Exec(ctx, "UPDATE currency SET description = $1 WHERE code = $2;", &newCurrency.Description, &newCurrency.Code)
	if err != nil {
		return err
	}

	return nil
}

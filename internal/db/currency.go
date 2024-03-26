package db

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/whiterthanwhite/businessinsight/internal/entities/currency"
)

func (c *databaseConnection) InsertCurrency(parentCtx context.Context, currencies []currency.Currency) error {
	if len(currencies) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	tx, err := c.conn.Begin(ctx)
	if err != nil {
		return err
	}

	log.Println(currencies)
	for _, curr := range currencies {
		xCurrency, err := c.currencyExist(ctx, &curr)
		if err != nil {
			tErr := errors.Join(err)
			err = tx.Rollback(ctx)
			tErr = errors.Join(tErr, err)
			return tErr
		}
		if xCurrency != nil {
			needUpdate := false
			if xCurrency.Description != curr.Description {
				needUpdate = true
			}
			if needUpdate {
				log.Println("Update")
				_, err = tx.Exec(ctx, "UPDATE currency SET description = $1 WHERE code = $2;", curr.Description, curr.Code)
				if err != nil {
					tErr := errors.Join(err)
					err = tx.Rollback(ctx)
					tErr = errors.Join(tErr, err)
					return tErr
				}
			}
		} else {
			log.Println("Insert")
			_, err = tx.Exec(ctx, "INSERT INTO currency VALUES ($1, $2);", curr.Code, curr.Description)
			if err != nil {
				tErr := errors.Join(err)
				err = tx.Rollback(ctx)
				tErr = errors.Join(tErr, err)
				return tErr
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *databaseConnection) GetCurrencies(parentCtx context.Context) ([]currency.Currency, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	rows, err := c.conn.Query(ctx, "SELECT * FROM currency;")
	if err != nil {
		return nil, err
	}

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

func (c *databaseConnection) DeleteCurrencies(parentCtx context.Context, currencies []currency.Currency) error {
	if len(currencies) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	tx, err := c.conn.Begin(ctx)
	if err != nil {
		return err
	}

	for _, curr := range currencies {
		xCurrency, err := c.currencyExist(ctx, &curr)
		if err != nil {
			tErr := errors.Join(err)
			err = tx.Rollback(ctx)
			tErr = errors.Join(tErr, err)
			return tErr
		}
		if xCurrency != nil {
			_, err := tx.Exec(ctx, "DELETE FROM currency WHERE code = $1;", curr.Code)
			if err != nil {
				tErr := errors.Join(err)
				err = tx.Rollback(ctx)
				tErr = errors.Join(tErr, err)
				return tErr
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *databaseConnection) currencyExist(parentCtx context.Context, newCurrency *currency.Currency) (*currency.Currency, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	var xCurrency *currency.Currency
	var code, description string
	err := c.conn.QueryRow(ctx, "SELECT * FROM currency WHERE code = $1;", newCurrency.Code).Scan(&code, &description)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	if err != pgx.ErrNoRows {
		xCurrency = &currency.Currency{
			Code:        code,
			Description: description,
		}
	}

	return xCurrency, nil
}

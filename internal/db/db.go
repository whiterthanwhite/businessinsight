package db

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/whiterthanwhite/businessinsight/internal/entities/currency"
)

var dbConn *databaseConnection

type databaseConnection struct {
	conn *pgx.Conn
}

func (c *databaseConnection) Close(parentCtx context.Context) error {
	ctx, _ := context.WithCancel(parentCtx)
	return c.conn.Close(ctx)
}

func (c *databaseConnection) HelloWorld(parentCtx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*10)
	defer cancel()

	var message string
	err := c.conn.QueryRow(ctx, "select 'Hello, World!'").Scan(&message)
	if err != nil {
		return "", err
	}
	return message, nil
}

func (c *databaseConnection) CurrentTime(parentCtx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*10)
	defer cancel()

	var message string
	err := c.conn.QueryRow(ctx, "select $1", time.Now().Format(time.DateTime)).Scan(&message)
	if err != nil {
		return "", err
	}
	return message, nil
}

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

	for _, curr := range currencies {
		_, err = tx.Exec(ctx, "INSERT INTO currency VALUES ($1, $2);", curr.Code, curr.Description)
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

func (c *databaseConnection) InitTables(parentCtx context.Context) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	log.Println("Inititialize source tables")

	var count int
	err := c.conn.QueryRow(ctx, "SELECT COUNT(*) FROM pg_type WHERE typname = 'operation_type';").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = c.conn.Exec(ctx, operationType)
		if err != nil {
			return err
		}
		log.Println("operation_type created")
	}

	err = c.conn.QueryRow(ctx, "SELECT COUNT(*) FROM pg_class WHERE relname = 'currency';").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = c.conn.Exec(ctx, currencySQL)
		if err != nil {
			return err
		}
		log.Println("currency created")
	}

	err = c.conn.QueryRow(ctx, "SELECT COUNT(*) FROM pg_class WHERE relname = 'account';").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = c.conn.Exec(ctx, account)
		if err != nil {
			return err
		}
		log.Println("account created")
	}

	err = c.conn.QueryRow(ctx, "SELECT COUNT(*) FROM pg_class WHERE relname = 'category';").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = c.conn.Exec(ctx, category)
		if err != nil {
			return err
		}
		log.Println("category created")
	}

	err = c.conn.QueryRow(ctx, "SELECT COUNT(*) FROM pg_class WHERE relname = 'operation';").Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = c.conn.Exec(ctx, operation)
		if err != nil {
			return err
		}
		log.Println("operation created")
	}

	return nil
}

func Connect(parentCtx context.Context, connectionStr string) (*databaseConnection, error) {
	if dbConn != nil {
		return dbConn, nil
	}
	ctx, _ := context.WithCancel(parentCtx)
	dbConn := new(databaseConnection)
	var err error
	dbConn.conn, err = pgx.Connect(ctx, connectionStr)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

func GetInstance() *databaseConnection {
	return dbConn
}

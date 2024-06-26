package db

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

var dbConn *databaseConnection

type databaseConnection struct {
	conn  *pgx.Conn
	mutex sync.Mutex
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
		_, err = c.conn.Exec(ctx, QUERY_CREATE_OPERATION_TYPE)
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
		_, err = c.conn.Exec(ctx, QUERY_CREATE_TABLE_CURRENCY)
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
		_, err = c.conn.Exec(ctx, QUERY_CREATE_TABLE_ACCOUNT)
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
		_, err = c.conn.Exec(ctx, QUERY_CREATE_TABLE_CATEGORY)
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
		_, err = c.conn.Exec(ctx, QUERY_CREATE_TABLE_OPERATION)
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
	conn := new(databaseConnection)
	var err error
	conn.conn, err = pgx.Connect(ctx, connectionStr)
	if err != nil {
		return nil, err
	}
	dbConn = conn
	return dbConn, nil
}

func GetInstance() (*databaseConnection, error) {
	if dbConn == nil {
		return nil, errors.New("database variable was not initialized")
	}
	return dbConn, nil
}

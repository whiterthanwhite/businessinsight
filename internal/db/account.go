package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/whiterthanwhite/businessinsight/internal/entities/account"
)

func (d *databaseConnection) GetAccounts(parentCtx context.Context) ([]account.Account, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	rows, err := d.conn.Query(ctx, "SELECT * FROM account ORDER BY id;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []account.Account
	for rows.Next() {
		newAccount := account.Account{}
		err := rows.Scan(&newAccount.Id, &newAccount.Name, &newAccount.CurrencyCode)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, newAccount)
	}

	return accounts, nil
}

func (d *databaseConnection) GetAccount(parentCtx context.Context, newAccount *account.Account) (*account.Account, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	xAccount := new(account.Account)
	err := d.conn.QueryRow(ctx, `SELECT * FROM account WHERE id = $1 LIMIT 1;`, &newAccount.Id).
		Scan(&xAccount.Id, &xAccount.Name, &xAccount.CurrencyCode)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	if err == pgx.ErrNoRows {
		return nil, nil
	}

	return xAccount, nil
}

func (d *databaseConnection) InsertAccount(parentCtx context.Context, newAccount *account.Account) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, err := d.conn.Exec(ctx, `INSERT INTO account (name, currency_code) VALUES ($1, $2);`, newAccount.Name, newAccount.CurrencyCode)
	if err != nil {
		return err
	}

	return nil
}

func (d *databaseConnection) DeleteAccount(parentCtx context.Context, deleteAccount *account.Account) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	log.Println(deleteAccount)
	_, err := d.conn.Exec(ctx, `DELETE FROM account WHERE id = $1;`, deleteAccount.Id)
	if err != nil {
		return err
	}

	return nil
}

func (d *databaseConnection) UpdateAccount(parentCtx context.Context, newAccount *account.Account) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, err := d.conn.Exec(ctx, `UPDATE account SET name = $1 currency_code = $2 WHERE id = $3;`, newAccount.Name,
		newAccount.CurrencyCode, newAccount.Id)
	if err != nil {
		return err
	}

	return nil
}

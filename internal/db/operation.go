package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/whiterthanwhite/businessinsight/internal/entities/operation"
)

func (d *databaseConnection) InsertOperation(parentCtx context.Context, newOperation *operation.Operation) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	ct, err := d.conn.Exec(ctx,
		`
		INSERT INTO operation (date_time, type, amount, source_id, currency_code, category_id, transaction_no, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
		`,
		&newOperation.DateTime,
		&newOperation.Type,
		&newOperation.Amount,
		&newOperation.SourceId,
		&newOperation.CurrencyCode,
		&newOperation.CategoryId,
		&newOperation.TransactionNo,
		&newOperation.Description,
	)
	if err != nil {
		return err
	}

	log.Printf("Insert: %v; Rows affected: %v\n", ct.Insert(), ct.RowsAffected())
	return nil
}

func (d *databaseConnection) GetOperations(parentCtx context.Context) ([]operation.Operation, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	rows, err := d.conn.Query(ctx, `SELECT * FROM operation ORDER BY date_time DESC, entry_no;`)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	operations := make([]operation.Operation, 0)
	for rows.Next() {
		operation := new(operation.Operation)
		if err = rows.Scan(
			&operation.EntryNo,
			&operation.DateTime,
			&operation.Type,
			&operation.Amount,
			&operation.SourceId,
			&operation.CurrencyCode,
			&operation.CategoryId,
			&operation.TransactionNo,
			&operation.Description,
		); err != nil {
			return nil, err
		}
		operations = append(operations, *operation)
	}

	return operations, nil
}

func (d *databaseConnection) GetOperation(parentCtx context.Context, newOpeartion *operation.Operation) (*operation.Operation, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	operation := new(operation.Operation)
	err := d.conn.QueryRow(ctx, `SELECT * FROM operation WHERE entry_no = $1 LIMIT 1;`, &newOpeartion.EntryNo).Scan(
		&operation.EntryNo,
		&operation.DateTime,
		&operation.Type,
		&operation.Amount,
		&operation.SourceId,
		&operation.CurrencyCode,
		&operation.CategoryId,
		&operation.TransactionNo,
		&operation.Description,
	)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	} else if err == pgx.ErrNoRows {
		return nil, nil
	}
	return operation, nil
}

func (d *databaseConnection) UpdateOperation(parentCtx context.Context, newOperation *operation.Operation) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	ct, err := d.conn.Exec(ctx,
		`
		UPDATE operation
		SET date_time = $1, type = $2, amount = $3, source_id = $4, currency_code = $5, category_id = $6, transaction_no = $7, description = $8
		WHERE entry_no = $9;
		`,
		&newOperation.DateTime,
		&newOperation.Type,
		&newOperation.Amount,
		&newOperation.SourceId,
		&newOperation.CurrencyCode,
		&newOperation.CategoryId,
		&newOperation.TransactionNo,
		&newOperation.Description,
		&newOperation.EntryNo,
	)
	if err != nil {
		return err
	}

	log.Printf("Update: %v; Row affected: %v\n", ct.Update(), ct.RowsAffected())
	return nil
}

func (d *databaseConnection) DeleteOperation(parentCtx context.Context, deleteOperation *operation.Operation) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	ct, err := d.conn.Exec(ctx, `DELETE FROM operation WHERE entry_no = $1;`, &deleteOperation.EntryNo)
	if err != nil {
		return err
	}

	log.Printf("Delete: %v; Row affected: %v\n", ct.Delete(), ct.RowsAffected())
	return nil
}

func (d *databaseConnection) GetMaxTransactionNo(parentCtx context.Context) (int, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	var lastTransactionNo int
	err := d.conn.QueryRow(ctx, `SELECT max(transaction_no) FROM operation;`).Scan(&lastTransactionNo)
	if err != nil && err != pgx.ErrNoRows {
		return 0, err
	} else if err == pgx.ErrNoRows {
		return 0, nil
	}

	return lastTransactionNo, nil
}

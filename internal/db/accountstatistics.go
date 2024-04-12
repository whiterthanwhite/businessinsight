package db

import (
	"context"

	"github.com/whiterthanwhite/businessinsight/internal/entities/accountstatistics"
)

func (d *databaseConnection) GetAccountStatistics(parentCtx context.Context) ([]accountstatistics.AccountStatistics, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	rows, err := d.conn.Query(ctx,
		`
			SELECT account.name, SUM(operation.amount), account.currency_code
			FROM operation JOIN account ON account.id = operation.source_id
			GROUP BY account.name, account.currency_code;
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accountsStatistics []accountstatistics.AccountStatistics
	for rows.Next() {
		accountStatistics := new(accountstatistics.AccountStatistics)
		err = rows.Scan(&accountStatistics.Name, &accountStatistics.Total, nil)
		if err != nil {
			return nil, err
		}
		accountsStatistics = append(accountsStatistics, *accountStatistics)
	}
	return accountsStatistics, nil
}

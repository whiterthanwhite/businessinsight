package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/whiterthanwhite/businessinsight/internal/entities/category"
)

func (d *databaseConnection) InsertCategory(parentCtx context.Context, newCategory *category.Category) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, err := d.conn.Exec(ctx, `INSERT INTO category (type, name, description) VALUES ($1, $2, $3);`, &newCategory.Type,
		&newCategory.Name, &newCategory.Description)
	if err != nil {
		return err
	}

	return nil
}

func (d *databaseConnection) GetCategory(parentCtx context.Context, newCategory *category.Category) (*category.Category, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	xCategory := new(category.Category)
	err := d.conn.QueryRow(ctx, `SELECT * FROM category WHERE id = $1;`, &newCategory.Id).
		Scan(&xCategory.Id, &xCategory.Type, &xCategory.Name, &xCategory.Description)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	if err == pgx.ErrNoRows {
		return nil, nil
	}

	return xCategory, nil
}

func (d *databaseConnection) GetCategories(parentCtx context.Context) ([]category.Category, error) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	rows, err := d.conn.Query(ctx, `SELECT * FROM category ORDER BY id;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []category.Category
	for rows.Next() {
		category := category.Category{}
		err = rows.Scan(&category.Id, &category.Type, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (d *databaseConnection) UpdateCategory(parentCtx context.Context, category *category.Category) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, err := d.conn.Exec(ctx, `UPDATE category SET type = $1, name = $2, description = $3 WHERE id = $4`,
		&category.Type, &category.Name, &category.Description, &category.Id)
	if err != nil {
		return err
	}
	return nil
}

func (d *databaseConnection) DeleteCategory(parentCtx context.Context, deleteCategory *category.Category) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	d.mutex.Lock()
	defer d.mutex.Unlock()

	_, err := d.conn.Exec(ctx, `DELETE FROM category WHERE id = $1;`, &deleteCategory.Id)
	if err != nil {
		return err
	}
	return nil
}

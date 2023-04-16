package db

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Connect struct {
	sql *sql.DB
}

func NewConnect(driverName string, dsn string) (Connect, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return Connect{}, err
	}

	if err = db.Ping(); err != nil {
		return Connect{}, err
	}

	return Connect{
		sql: db,
	}, nil
}

func (c *Connect) Close() error {
	return c.sql.Close()
}

func (c *Connect) CreateBrand(ctx context.Context, name string) (int, error) {
	stmt, err := c.sql.PrepareContext(ctx, CreateBrand())
	if err != nil {
		return 0, err
	}

	var brandId int
	if err := stmt.QueryRowContext(ctx, name).Scan(&brandId); err != nil {
		return 0, err
	}

	return brandId, nil
}

func (c *Connect) CreateCategory(ctx context.Context, name string) (int, error) {
	stmt, err := c.sql.PrepareContext(ctx, CreateCategory())
	if err != nil {
		return 0, err
	}

	var categoryId int
	if err := stmt.QueryRowContext(ctx, name).Scan(&categoryId); err != nil {
		return 0, err
	}

	return categoryId, nil
}

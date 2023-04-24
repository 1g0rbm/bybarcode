package db

import (
	"bybarcode/internal/auth"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"

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

func (c *Connect) CreateAccountIfNotExist(ctx context.Context, id int, username string, firstName string, lastName string) error {
	stmt, err := c.sql.PrepareContext(ctx, CreateAccountIfNotExist())
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, id, username, firstName, lastName)
	return err
}

func (c *Connect) CreateSession(ctx context.Context, session auth.Session) error {
	tx, err := c.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	stmtCheck, err := c.sql.PrepareContext(ctx, findAccount())
	if err != nil {
		return err
	}

	acc := auth.Account{}
	err = stmtCheck.
		QueryRowContext(ctx, session.AccountID).
		Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Username)

	if errors.As(err, &pgx.ErrNoRows) {
		return fmt.Errorf("there is no account with id %d", session.AccountID)
	} else if err != nil {
		return err
	}

	stmtInsert, err := c.sql.PrepareContext(ctx, CreateSession())
	if err != nil {
		return err
	}

	_, err = stmtInsert.ExecContext(
		ctx,
		session.ID,
		session.Token,
		session.RefreshToken,
		session.AccountID,
		session.ExpireAt,
		session.CreatedAt,
		session.UpdatedAt,
	)

	return err
}

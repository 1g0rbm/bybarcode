package db

import (
	"bybarcode/internal/auth"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"time"

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

func (c *Connect) FindAccountById(ctx context.Context, accountId int64) (auth.Account, error) {
	acc := auth.Account{}
	stmtCheck, err := c.sql.PrepareContext(ctx, findAccount())
	if err != nil {
		return acc, err
	}

	err = stmtCheck.
		QueryRowContext(ctx, accountId).
		Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Username)
	if err != nil {
		return acc, err
	}

	return acc, nil
}

func (c *Connect) CreateSession(ctx context.Context, accountId int64) (auth.Session, error) {
	session := auth.Session{}

	tx, err := c.sql.BeginTx(ctx, nil)
	if err != nil {
		return session, err
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
		return session, err
	}

	acc := auth.Account{}
	err = stmtCheck.
		QueryRowContext(ctx, accountId).
		Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Username)

	if errors.As(err, &pgx.ErrNoRows) {
		return session, fmt.Errorf("there is no account with id %d", accountId)
	} else if err != nil {
		return session, err
	}

	stmtSessionCheck, err := c.sql.PrepareContext(ctx, findNotExpiredSessionByAccountId())
	if err != nil {
		return session, err
	}

	err = stmtSessionCheck.
		QueryRowContext(ctx, accountId).
		Scan(
			&session.ID,
			&session.Token,
			&session.RefreshToken,
			&session.AccountID,
			&session.ExpireAt,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
	if err != nil && !errors.As(err, &pgx.ErrNoRows) {
		return session, err
	}

	if session.ID != uuid.Nil {
		return session, nil
	}

	stmtInsert, err := c.sql.PrepareContext(ctx, CreateSession())
	if err != nil {
		return session, err
	}

	session.ID = uuid.New()
	session.Token = uuid.New()
	session.RefreshToken = uuid.New()
	session.AccountID = accountId
	session.ExpireAt = time.Now().Add(24 * time.Hour)
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

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

	return session, err
}

func (c *Connect) FindNotExpiredSession(ctx context.Context, token string) (auth.Session, error) {
	session := auth.Session{}

	stmt, err := c.sql.PrepareContext(ctx, findNotExpiredSessionByToken())
	if err != nil {
		return session, err
	}

	err = stmt.
		QueryRowContext(ctx, token).
		Scan(
			&session.ID,
			&session.Token,
			&session.RefreshToken,
			&session.AccountID,
			&session.ExpireAt,
			&session.CreatedAt,
			&session.UpdatedAt,
		)

	return session, err
}

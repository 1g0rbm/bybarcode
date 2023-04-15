package db

import (
	"database/sql"
)

type Connect struct {
	sql *sql.DB
}

func NewConnect(driverName string, dsn string) (Connect, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return Connect{}, err
	}

	return Connect{
		sql: db,
	}, nil
}

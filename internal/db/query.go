package db

import "strings"

func CreateBrand() string {
	query := `
INSERT INTO brands
	(name)
	VALUES ($1)
	RETURNING id;
`
	return strings.Trim(query, " ")
}

func CreateCategory() string {
	query := `
INSERT INTO brands
	(name)
	VALUES ($1)
	RETURNING id;
`
	return strings.Trim(query, " ")
}

func CreateAccountIfNotExist() string {
	query := `
INSERT INTO account (id, username, first_name, last_name)
SELECT $1, $2, $3, $4
WHERE NOT EXISTS (
  SELECT * FROM account 
  WHERE id = $1
)
`
	return strings.Trim(query, " ")
}

func findAccount() string {
	return `SELECT * FROM account where id = $1`
}

func CreateSession() string {
	query := `
        INSERT INTO sessions (id, token, refresh_token, account_id, expire_at, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

	return strings.Trim(query, " ")
}

func findNotExpiredSession() string {
	query := `
	SELECT * FROM sessions WHERE account_id = $1 AND expire_at > now()
`
	return strings.Trim(query, " ")
}

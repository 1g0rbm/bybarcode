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

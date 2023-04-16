package db

import "strings"

const createCategory = `
INSERT INTO categories
	(name)
	VALUES ($1)
	RETURNING id;
`

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

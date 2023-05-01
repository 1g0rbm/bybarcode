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

func CreateProduct() string {
	query := `
INSERT INTO products 
    (name, upcean, category_id, brand_id) 
	VALUES ($1, $2, $3, $4)
	RETURNING id
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

func findNotExpiredSessionByAccountId() string {
	query := `
	SELECT * FROM sessions WHERE account_id = $1 AND expire_at > now()
`
	return strings.Trim(query, " ")
}

func findNotExpiredSessionByToken() string {
	query := `
	SELECT * FROM sessions WHERE token = $1 AND expire_at > now()
`
	return strings.Trim(query, " ")
}

func findProductByBarcode() string {
	query := `
	SELECT p.id, p.name, p.upcean, p.category_id, p.brand_id, c.id, c.name , b.id, b.name FROM products p
	LEFT JOIN categories c on p.category_id = c.id
	LEFT JOIN brands b on b.id = p.brand_id
	WHERE p.upcean = $1
`
	return strings.Trim(query, " ")
}

func findProductByIdOrBarcode() string {
	return `SELECT id, name, upcean, category_id, brand_id FROM products WHERE id = $1 OR upcean = $1`
}

func updateProductById() string {
	query := `
	UPDATE products
	SET name = $1, upcean = $2, category_id = $3, brand_id = $4
	WHERE id = $5
	RETURNING id;
`
	return strings.Trim(query, " ")
}

func deleteProductById() string {
	query := `
	DELETE FROM products
	WHERE id = $1
	RETURNING id;
`
	return strings.Trim(query, " ")
}

func createShoppingList() string {
	query := `
INSERT INTO shopping_lists 
    (name, account_id) 
	VALUES ($1, $2)
	RETURNING id
`
	return strings.Trim(query, " ")
}

func updateShoppingListById() string {
	query := `
	UPDATE shopping_lists
	SET name = $1
	WHERE id = $2
	RETURNING id;
`
	return strings.Trim(query, " ")
}

func findShoppingListById() string {
	return `SELECT id, name, account_id FROM shopping_lists WHERE id = $1`
}

func getShoppingListsByAccount() string {
	return `SELECT id, name FROM shopping_lists WHERE account_id = $1`
}

func deleteShoppingListById() string {
	query := `
	DELETE FROM shopping_lists
	WHERE id = $1
	RETURNING id;
`
	return strings.Trim(query, " ")
}

func addProductToShoppingList() string {
	return `INSERT INTO shopping_list__products (shopping_list_id, product_id) VALUES ($1, $2)`
}

func getShoppingListProducts() string {
	query := `
	SELECT p.id, p.name, p.upcean, p.category_id, p.brand_id FROM products p
	JOIN shopping_list__products slp on p.id = slp.product_id
	WHERE slp.shopping_list_id = $1
`

	return strings.Trim(query, " ")
}

func deleteProductFromShoppingList() string {
	query := `
	DELETE FROM shopping_list__products 
	WHERE shopping_list_id = $1 AND product_id = $2
`
	return strings.Trim(query, " ")
}

func toggleProductStateInList() string {
	query := `
	UPDATE shopping_list__products SET checked = NOT checked
	WHERE shopping_list_id = $1 AND product_id = $2
	RETURNING shopping_list_id, product_id
`

	return strings.Trim(query, " ")
}

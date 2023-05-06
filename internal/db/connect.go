package db

import (
	"bybarcode/internal/auth"
	"bybarcode/internal/products"
	"bybarcode/internal/stat"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Events struct {
	ids chan int64
}

type Connect struct {
	sql *sql.DB
}

var ErrDuplicateKey error

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

func (c *Connect) UpdateProduct(ctx context.Context, p products.Product) (products.Product, error) {
	stmt, err := c.sql.PrepareContext(ctx, updateProductById())
	if err != nil {
		return p, err
	}

	var productId int64
	err = stmt.
		QueryRowContext(ctx, p.Name, p.Upcean, p.CategoryId, p.BrandId, p.ID).
		Scan(&productId)

	return p, err
}

func (c *Connect) DeleteProduct(ctx context.Context, id int64) error {
	stmt, err := c.sql.PrepareContext(ctx, deleteProductById())
	if err != nil {
		return err
	}

	var productId int64
	err = stmt.QueryRowContext(ctx, id).Scan(&productId)

	return err
}

func (c *Connect) CreateProduct(ctx context.Context, p products.Product) (int64, error) {
	stmt, err := c.sql.PrepareContext(ctx, CreateProduct())
	if err != nil {
		return 0, err
	}

	var productId int64
	if err := stmt.QueryRowContext(ctx, p.Name, p.Upcean, p.CategoryId, p.BrandId).Scan(&productId); err != nil {
		return 0, err
	}

	return productId, nil
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

func (c *Connect) FindProductByBarcode(ctx context.Context, barcode int64) (products.Product, error) {
	p := products.Product{
		Category: &products.Category{},
		Brand:    &products.Brand{},
	}

	stmt, err := c.sql.PrepareContext(ctx, findProductByBarcode())
	if err != nil {
		return p, err
	}

	err = stmt.
		QueryRowContext(ctx, barcode).
		Scan(
			&p.ID,
			&p.Name,
			&p.Upcean,
			&p.CategoryId,
			&p.BrandId,
			&p.Category.ID,
			&p.Category.Name,
			&p.Brand.ID,
			&p.Brand.Name,
		)

	return p, err
}

func (c *Connect) CreateShoppingList(ctx context.Context, sl products.ShoppingList) (int64, error) {
	stmt, err := c.sql.PrepareContext(ctx, createShoppingList())
	if err != nil {
		return 0, err
	}

	var listId int64
	if err := stmt.QueryRowContext(ctx, sl.Name, sl.AccountId).Scan(&listId); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			ErrDuplicateKey = fmt.Errorf("duplciate shopping list name %s", sl.Name)
			return 0, ErrDuplicateKey
		}
		return 0, err
	}

	return listId, nil
}

func (c *Connect) UpdateShoppingList(ctx context.Context, sl products.ShoppingList) (products.ShoppingList, error) {
	stmt, err := c.sql.PrepareContext(ctx, updateShoppingListById())
	if err != nil {
		return sl, err
	}

	var slId int64
	err = stmt.
		QueryRowContext(ctx, sl.Name, sl.ID).
		Scan(&slId)

	return sl, err
}

func (c *Connect) GetShoppingListsByAccount(ctx context.Context, accId int64) ([]products.ShoppingList, error) {
	var (
		id   int64
		name string
	)

	stmt, err := c.sql.PrepareContext(ctx, getShoppingListsByAccount())
	if err != nil {
		return nil, err
	}

	r, err := stmt.QueryContext(ctx, accId)
	if err != nil {
		return nil, err
	}

	defer func(r *sql.Rows) {
		if rErr := r.Close(); rErr != nil {
			err = rErr
		}
	}(r)

	var slm []products.ShoppingList
	for r.Next() {
		if err = r.Scan(&id, &name); err != nil {
			return nil, err
		}

		sl := products.ShoppingList{
			ID:        id,
			Name:      name,
			AccountId: accId,
		}

		slm = append(slm, sl)
	}

	if err = r.Err(); err != nil {
		return nil, err
	}

	return slm, err
}

func (c *Connect) DeleteShoppingList(ctx context.Context, id int64) error {
	stmt, err := c.sql.PrepareContext(ctx, deleteShoppingListById())
	if err != nil {
		return err
	}

	var slId int64
	err = stmt.QueryRowContext(ctx, id).Scan(&slId)

	return err
}

func (c *Connect) AddProductToShoppingListByIds(ctx context.Context, productId int64, listId int64) error {
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

	productStmt, err := tx.PrepareContext(ctx, findProductByIdOrBarcode())
	if err != nil {
		return err
	}

	p := products.Product{}
	err = productStmt.QueryRowContext(ctx, productId).Scan(&p.ID, &p.Name, &p.Upcean, &p.CategoryId, &p.BrandId)
	if err != nil {
		return err
	}

	slStmt, err := tx.PrepareContext(ctx, findShoppingListById())
	if err != nil {
		return err
	}

	sl := products.ShoppingList{}
	err = slStmt.QueryRowContext(ctx, listId).Scan(&sl.ID, &sl.Name, &sl.AccountId)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, addProductToShoppingList())
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, sl.ID, p.ID)
	if err != nil && strings.Contains(err.Error(), "shopping_list__products_pkey") {
		ErrDuplicateKey = fmt.Errorf("there is already exist product %s in list %s", p.Name, sl.Name)
		return ErrDuplicateKey
	}

	err = tx.Commit()

	return err
}

func (c *Connect) GetShoppingListProducts(ctx context.Context, slId int64) ([]products.ProductInList, error) {
	stmt, err := c.sql.PrepareContext(ctx, getShoppingListProducts())
	if err != nil {
		return nil, err
	}

	r, err := stmt.QueryContext(ctx, slId)
	if err != nil {
		return nil, err
	}

	defer func(r *sql.Rows) {
		if rErr := r.Close(); rErr != nil {
			err = rErr
		}
	}(r)

	var (
		pList      []products.ProductInList
		id         int64
		name       string
		barcode    int64
		categoryId int64
		brandId    int64
		checked    bool
	)
	for r.Next() {
		if err = r.Scan(&id, &name, &barcode, &categoryId, &brandId, &checked); err != nil {
			return nil, err
		}

		p := products.ProductInList{
			Product: products.Product{
				ID:         id,
				Name:       name,
				Upcean:     barcode,
				CategoryId: categoryId,
				BrandId:    brandId,
			},
			Checked: checked,
		}

		pList = append(pList, p)
	}

	if err = r.Err(); err != nil {
		return nil, err
	}

	return pList, err
}

func (c *Connect) DeleteProductFromShoppingList(ctx context.Context, slId int64, pId int64) error {
	stmt, err := c.sql.PrepareContext(ctx, deleteProductFromShoppingList())
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, slId, pId)

	return err
}

func (c *Connect) ToggleProductStateInShoppingList(ctx context.Context, slId int64, pId int64) error {
	stmt, err := c.sql.PrepareContext(ctx, toggleProductStateInList())
	if err != nil {
		return err
	}

	var (
		_slId int64
		_pId  int64
	)
	return stmt.QueryRowContext(ctx, slId, pId).Scan(&_slId, &_pId)
}

func (c *Connect) GetStatistic(ctx context.Context, from time.Time, to time.Time) ([]stat.Statistic, error) {
	stmt, err := c.sql.PrepareContext(ctx, getStatistic())
	if err != nil {
		return nil, err
	}

	r, err := stmt.QueryContext(ctx, from, to)
	if err != nil {
		return nil, err
	}

	defer func(r *sql.Rows) {
		if rErr := r.Close(); rErr != nil {
			err = rErr
		}
	}(r)

	var (
		sList                []stat.Statistic
		id                   int64
		name                 string
		slId                 int64
		createdAt            time.Time
		productsCount        int
		checkedProductsCount int
	)
	for r.Next() {
		if err = r.Scan(&id, &name, &slId, &createdAt, &productsCount, &checkedProductsCount); err != nil {
			return nil, err
		}

		p := stat.Statistic{
			ID:                   id,
			ShoppingListName:     name,
			ShoppingListId:       slId,
			CreatedAt:            createdAt,
			ProductsCount:        productsCount,
			CheckedProductsCount: checkedProductsCount,
		}

		sList = append(sList, p)
	}

	if err = r.Err(); err != nil {
		return nil, err
	}

	return sList, err
}

func (c *Connect) AddedUpdStatisticByShoppingList(ctx context.Context, listId int64) error {
	statStmt, err := c.sql.PrepareContext(ctx, updateStatByAddingProduct())
	if err != nil {
		return err
	}

	_, err = statStmt.ExecContext(ctx, listId)

	return err
}

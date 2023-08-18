package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// NewImplStorageProductMySQL returns new ImplStorageProductMySQL
func NewImplStorageProductMySQL(db *sql.DB) *ImplStorageProductMySQL {
	return &ImplStorageProductMySQL{db: db}
}

// ProductMySQL is a product model for MySQL
type ProductMySQL struct {
	Id			sql.NullInt32
	Name		sql.NullString
	Quantity	sql.NullInt32
	CodeValue	sql.NullString
	IsPublished	sql.NullBool
	Expiration	sql.NullTime
	Price		sql.NullFloat64
}

// ImplStorageProductMySQL is an implementation of StorageProduct interface
type ImplStorageProductMySQL struct {
	db *sql.DB
}

// GetOne returns one product by id
func (impl *ImplStorageProductMySQL) GetOne(id int) (p *Product, err error) {
	// query
	query := "SELECT id, name, quantity, code_value, is_published, expiration, price FROM products WHERE id = ?"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = impl.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}
	defer stmt.Close()

	// execute query
	row := stmt.QueryRow(id)
	if row.Err() != nil {
		err = row.Err()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			err = fmt.Errorf("%w. %v", ErrStorageProductNotFound, row.Err())
		default:
			err = fmt.Errorf("%w. %v", ErrStorageProductInternal, row.Err())
		}

		return
	}

	// scan row
	var pr ProductMySQL
	err = row.Scan(&pr.Id, &pr.Name, &pr.Quantity, &pr.CodeValue, &pr.IsPublished, &pr.Expiration, &pr.Price)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	// serialization
	p = new(Product)
	if pr.Id.Valid {
		(*p).Id = int(pr.Id.Int32)
	}
	if pr.Name.Valid {
		(*p).Name = pr.Name.String
	}
	if pr.Quantity.Valid {
		(*p).Quantity = int(pr.Quantity.Int32)
	}
	if pr.CodeValue.Valid {
		(*p).CodeValue = pr.CodeValue.String
	}
	if pr.IsPublished.Valid {
		(*p).IsPublished = pr.IsPublished.Bool
	}
	if pr.Expiration.Valid {
		(*p).Expiration = pr.Expiration.Time
	}
	if pr.Price.Valid {
		(*p).Price = pr.Price.Float64
	}

	return
}

// Store stores product
func (impl *ImplStorageProductMySQL) Store(p *Product) (err error) {
	// deserialize
	var pr ProductMySQL
	if (*p).Name != "" {
		pr.Name.Valid = true
		pr.Name.String = (*p).Name
	}
	if (*p).Quantity != 0 {
		pr.Quantity.Valid = true
		pr.Quantity.Int32 = int32((*p).Quantity)
	}
	if (*p).CodeValue != "" {
		pr.CodeValue.Valid = true
		pr.CodeValue.String = (*p).CodeValue
	}
	if !(*p).IsPublished {
		pr.IsPublished.Valid = true
		pr.IsPublished.Bool = (*p).IsPublished
	}
	if !(*p).Expiration.IsZero() {
		pr.Expiration.Valid = true
		pr.Expiration.Time = (*p).Expiration
	}
	if (*p).Price != 0 {
		pr.Price.Valid = true
		pr.Price.Float64 = (*p).Price
	}

	// query
	query := "INSERT INTO products (name, quantity, code_value, is_published, expiration, price) VALUES (?, ?, ?, ?, ?, ?)"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = impl.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}
	defer stmt.Close()

	// execute query
	result, err := stmt.Exec(pr.Name, pr.Quantity, pr.CodeValue, pr.IsPublished, pr.Expiration, pr.Price)
	if err != nil {
		errMySQL, ok := err.(*mysql.MySQLError); if ok {
			switch errMySQL.Number {
			case 1062:
				err = fmt.Errorf("%w. %v", ErrStorageProductNotUnique, err)
			default:
				err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
			}

			return
		}

		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("%w. %s", ErrStorageProductInternal, "rows affected != 1")
		return
	}

	// get last insert id
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	(*p).Id = int(lastInsertID)
	
	return
}

// Update updates product
func (impl *ImplStorageProductMySQL) Update(p *Product) (err error) {
	// deserialize
	var pr ProductMySQL
	if (*p).Name != "" {
		pr.Name.Valid = true
		pr.Name.String = (*p).Name
	}
	if (*p).Quantity != 0 {
		pr.Quantity.Valid = true
		pr.Quantity.Int32 = int32((*p).Quantity)
	}
	if (*p).CodeValue != "" {
		pr.CodeValue.Valid = true
		pr.CodeValue.String = (*p).CodeValue
	}
	if !(*p).IsPublished {
		pr.IsPublished.Valid = true
		pr.IsPublished.Bool = (*p).IsPublished
	}
	if !(*p).Expiration.IsZero() {
		pr.Expiration.Valid = true
		pr.Expiration.Time = (*p).Expiration
	}
	if (*p).Price != 0 {
		pr.Price.Valid = true
		pr.Price.Float64 = (*p).Price
	}

	// query
	query := "UPDATE products SET name = ?, quantity = ?, code_value = ?, is_published = ?, expiration = ?, price = ? WHERE id = ?"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = impl.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}
	defer stmt.Close()

	// execute query
	result, err := stmt.Exec(pr.Name, pr.Quantity, pr.CodeValue, pr.IsPublished, pr.Expiration, pr.Price, (*p).Id)
	if err != nil {
		errMySQL, ok := err.(*mysql.MySQLError); if ok {
			switch errMySQL.Number {
			case 1062:
				err = fmt.Errorf("%w. %v", ErrStorageProductNotUnique, err)
			default:
				err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
			}

			return
		}

		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("%w. %s", ErrStorageProductInternal, "rows affected != 1")
		return
	}

	return
}

// Delete deletes product by id
func (impl *ImplStorageProductMySQL) Delete(id int) (err error) {
	// query
	query := "DELETE FROM products WHERE id = ?"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = impl.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}
	defer stmt.Close()

	// execute query
	result, err := stmt.Exec(id)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageProductInternal, err)
		return
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("%w. %s", ErrStorageProductInternal, "rows affected != 1")
		return
	}

	return
}
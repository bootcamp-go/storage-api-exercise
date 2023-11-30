package repository

import (
	"app/internal"
	"database/sql"
)

// NewProductsDefault returns a new instance of ProductsDefault
func NewProductsDefault(db *sql.DB) *ProductsDefault {
	return &ProductsDefault{
		db: db,
	}
}

// ProductsDefault is a struct that represents a product repository
type ProductsDefault struct {
	// db is the database connection
	db *sql.DB
}

// GetOne returns a product by id
func (r *ProductsDefault) GetOne(id int) (p internal.Product, err error) {
	// execute the query
	row := r.db.QueryRow(
		"SELECT `id`, `name`, `quantity`, `code_value`, `is_published`, `expiration`, `price` " +
		"FROM `products` WHERE `id` = ?",
		id,
	)
	if err = row.Err(); err != nil {
		return
	}

	// scan the row into the product
	err = row.Scan(&p.ID, &p.Name, &p.Quantity, &p.CodeValue, &p.IsPublished, &p.Expiration, &p.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrProductNotFound
		}
		return
	}

	return
}

// Store stores a product
func (r *ProductsDefault) Store(p *internal.Product) (err error) {
	// execute the query
	result, err := r.db.Exec(
		"INSERT INTO `products` (`name`, `quantity`, `code_value`, `is_published`, `expiration`, `price`) " +
		"VALUES (?, ?, ?, ?, ?, ?)",
		p.Name, p.Quantity, p.CodeValue, p.IsPublished, p.Expiration, p.Price,
	)
	if err != nil {
		return
	}

	// get the last inserted id
	id, err := result.LastInsertId()
	if err != nil {
		return
	}
	p.ID = int(id)

	return
}

// Update updates a product
func (r *ProductsDefault) Update(p *internal.Product) (err error) {
	// execute the query
	_, err = r.db.Exec(
		"UPDATE `products` SET `name` = ?, `quantity` = ?, `code_value` = ?, `is_published` = ?, `expiration` = ?, `price` = ? " +
		"WHERE `id` = ?",
		p.Name, p.Quantity, p.CodeValue, p.IsPublished, p.Expiration, p.Price, p.ID,
	)
	if err != nil {
		return
	}

	return
}

// Delete deletes a product by id
func (r *ProductsDefault) Delete(id int) (err error) {
	// execute the query
	_, err = r.db.Exec(
		"DELETE FROM `products` WHERE `id` = ?",
		id,
	)
	if err != nil {
		return
	}

	return
}
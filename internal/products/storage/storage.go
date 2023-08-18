package storage

import (
	"errors"
	"time"
)

// Product is a product model
type Product struct {
	Id			int
	Name    	string
	Quantity	int
	CodeValue	string
	IsPublished bool
	Expiration  time.Time
	Price       float64
}

// StorageProduct is an interface for product storage
type StorageProduct interface {
	// GetOne returns one product by id
	GetOne(id int) (p *Product, err error)

	// Store stores product
	Store(p *Product) (err error)

	// Update updates product
	Update(p *Product) (err error)

	// Delete deletes product by id
	Delete(id int) (err error)
}

var (
	ErrStorageProductInternal = errors.New("internal storage product error")
	ErrStorageProductNotFound = errors.New("storage product not found")
	ErrStorageProductNotUnique = errors.New("storage product not unique")
	ErrStorageProductRelation = errors.New("storage product relation error")
)
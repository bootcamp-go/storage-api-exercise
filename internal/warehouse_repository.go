package internal

import "errors"

var (
	// ErrWarehouseNotFound is an error that will be returned when a warehouse is not found
	ErrWarehouseNotFound = errors.New("repository: warehouse not found")
	// ErrWarehouseNotUnique is an error that will be returned when a warehouse is not unique
	ErrWarehouseNotUnique = errors.New("repository: warehouse not unique")
	// ErrWarehouseFilter is an error that will be returned when a warehouse filter is not valid
	ErrWarehouseFilter = errors.New("repository: warehouse filter not valid")
)

// RepositoryWarehouses is an interface that represents a warehouse repository
type RepositoryWarehouses interface {
	// GetOne returns a warehouse by id
	GetOne(id int) (w Warehouse, err error)
	// GetAll returns all warehouses
	GetAll() (w []Warehouse, err error)
	// GetReportProducts returns a report of the amount of products in the warehouses
	GetReportProducts(filter map[string]any) (w []WarehouseReportProducts, err error)
	// Store stores a warehouse
	Store(w *Warehouse) (err error)
}
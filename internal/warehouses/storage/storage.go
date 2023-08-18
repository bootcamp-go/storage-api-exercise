package storage

import "errors"

// Warehouse is a warehouse model
type Warehouse struct {
	Id        int
	Name      string
	Address   string
	Telephone string
	Capacity  int
}

// ReportProducts is a report of products in a warehouse
type ReportProducts struct {
	WarehouseName string
	ProductsCount int
}

// StorageWarehouse is an interface for warehouse storage
type StorageWarehouse interface {
	// GetOne returns one warehouse by id
	GetOne(id int) (w *Warehouse, err error)

	// GetAll returns all warehouses
	GetAll() (w []*Warehouse, err error)

	// Create creates new warehouse
	Create(w *Warehouse) (err error)

	// ReportProducts returns report of products in warehouses
	ReportProducts(id int) (rs []*ReportProducts, err error)
}

var (
	ErrStorageWarehouseInternal = errors.New("storage warehouse internal error")
	ErrStorageWarehouseNotFound = errors.New("storage warehouse not found")
)
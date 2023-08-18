package storage

import (
	"database/sql"
	"errors"
	"fmt"
)

// NewImplStorageWarehouseMySQL returns new StorageMySQL
func NewImplStorageWarehouseMySQL(db *sql.DB) *StorageMySQL {
	return &StorageMySQL{db: db}
}

// WarehouseMySQL is a warehouse model for MySQL
type WarehouseMySQL struct {
	Id        sql.NullInt32
	Name      sql.NullString
	Address   sql.NullString
	Telephone sql.NullString
	Capacity  sql.NullInt32
}

// ReportProductsMySQL is a report of products in a warehouse for MySQL
type ReportProductsMySQL struct {
	WarehouseName sql.NullString
	ProductsCount sql.NullInt32
}

// StorageMySQL is an implementation of StorageWarehouse interface
type StorageMySQL struct {
	db *sql.DB
}

// GetOne returns one warehouse by id
func (s *StorageMySQL) GetOne(id int) (w *Warehouse, err error) {
	// query
	query := "SELECT id, name, address, telephone, capacity FROM warehouses WHERE id = ?"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageWarehouseInternal, err)
		return
	}
	defer stmt.Close()

	// execute query
	row := stmt.QueryRow(id)
	if row.Err() != nil {
		err = row.Err()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			err = fmt.Errorf("%w. %v", ErrStorageWarehouseNotFound, err)
		default:
			err = fmt.Errorf("%w. %v", ErrStorageWarehouseInternal, err)
		}

		return
	}

	// scan row
	var wh WarehouseMySQL
	err = row.Scan(&w.Id, &w.Name, &w.Address, &w.Telephone, &w.Capacity)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageWarehouseInternal, err)
		return
	}

	// serialization
	w = new(Warehouse)
	if wh.Id.Valid {
		w.Id = int(wh.Id.Int32)
	}
	if wh.Name.Valid {
		w.Name = wh.Name.String
	}
	if wh.Address.Valid {
		w.Address = wh.Address.String
	}
	if wh.Telephone.Valid {
		w.Telephone = wh.Telephone.String
	}
	if wh.Capacity.Valid {
		w.Capacity = int(wh.Capacity.Int32)
	}

	return
}

// GetAll returns all warehouses
func (s *StorageMySQL) GetAll() (w []*Warehouse, err error) {
	// query
	query := "SELECT id, name, address, telephone, capacity FROM warehouses"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageWarehouseInternal, err)
		return
	}
	defer stmt.Close()

	// execute query
	rows, err := stmt.Query()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageWarehouseInternal, err)
		return
	}
	defer rows.Close()

	// scan rows
	for rows.Next() {
		var wh WarehouseMySQL
		err = rows.Scan(&wh.Id, &wh.Name, &wh.Address, &wh.Telephone, &wh.Capacity)
		if err != nil {
			err = fmt.Errorf("%w. %v", ErrStorageWarehouseInternal, err)
			return
		}

		// serialization
		wr := new(Warehouse)
		if wh.Id.Valid {
			wr.Id = int(wh.Id.Int32)
		}
		if wh.Name.Valid {
			wr.Name = wh.Name.String
		}
		if wh.Address.Valid {
			wr.Address = wh.Address.String
		}
		if wh.Telephone.Valid {
			wr.Telephone = wh.Telephone.String
		}
		if wh.Capacity.Valid {
			wr.Capacity = int(wh.Capacity.Int32)
		}

		w = append(w, wr)
	}

	return
}

// Create creates new warehouse
func (s *StorageMySQL) Create(w *Warehouse) (err error) {
	// deserialization
	var wh WarehouseMySQL
	if w.Id != 0 {
		wh.Id.Valid = true
		wh.Id.Int32 = int32(w.Id)
	}
	if w.Name != "" {
		wh.Name.Valid = true
		wh.Name.String = w.Name
	}
	if w.Address != "" {
		wh.Address.Valid = true
		wh.Address.String = w.Address
	}
	if w.Telephone != "" {
		wh.Telephone.Valid = true
		wh.Telephone.String = w.Telephone
	}
	if w.Capacity != 0 {
		wh.Capacity.Valid = true
		wh.Capacity.Int32 = int32(w.Capacity)
	}

	// query
	query := "INSERT INTO warehouses(name, address, telephone, capacity) VALUES(?, ?, ?, ?)"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageWarehouseInternal, err)
		return
	}
	defer stmt.Close()

	// execute query
	var result sql.Result
	result, err = stmt.Exec(wh.Name, wh.Address, wh.Telephone, wh.Capacity)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageWarehouseInternal, err)
		return
	}

	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageWarehouseInternal, err)
		return
	}
	if rowsAffected != 1 {
		err = fmt.Errorf("%w. %v", ErrStorageWarehouseInternal, err)
		return
	}

	// get last insert id
	var lastInsertID int64
	lastInsertID, err = result.LastInsertId()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageWarehouseInternal, err)
		return
	}

	// set id
	(*w).Id = int(lastInsertID)

	return
}

// ReportProducts returns report of products in warehouse
func (s *StorageMySQL) ReportProducts(id int) (rs []*ReportProducts, err error) {
	// query
	query := "SELECT warehouses.name, COUNT(products.id) FROM warehouses LEFT JOIN products ON warehouses.id = products.warehouse_id GROUP BY warehouses.id HAVING IF(? = 0 or ? IS NULL, TRUE, warehouses.id = ?)"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageWarehouseInternal, err)
		return
	}

	// execute query
	var rows *sql.Rows
	rows, err = stmt.Query(id, id, id)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrStorageWarehouseInternal, err)
		return
	}
	defer rows.Close()

	// scan row
	if rows.Next() {
		var rp ReportProductsMySQL
		err = rows.Scan(&rp.WarehouseName, &rp.ProductsCount)
		if err != nil {
			err = fmt.Errorf("%w. %v", ErrStorageWarehouseInternal, err)
			return
		}

		// serialization
		r := new(ReportProducts)
		if rp.WarehouseName.Valid {
			r.WarehouseName = rp.WarehouseName.String
		}
		if rp.ProductsCount.Valid {
			r.ProductsCount = int(rp.ProductsCount.Int32)
		}

		rs = append(rs, r)
	}

	return
}


	


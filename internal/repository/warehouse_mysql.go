package repository

import (
	"app/internal"
	"database/sql"
)

// NewWarehousesMySQL returns a new instance of WarehousesMySQL
func NewWarehousesMySQL(db *sql.DB) *WarehousesMySQL {
	return &WarehousesMySQL{
		db: db,
	}
}

// WarehousesMySQL is a struct that represents a warehouse repository
type WarehousesMySQL struct {
	// db is the database connection
	db *sql.DB
}

// GetOne returns a warehouse by id
func (r *WarehousesMySQL) GetOne(id int) (w internal.Warehouse, err error) {
	// execute the query
	row := r.db.QueryRow(
		"SELECT `id`, `name`, `address`, `telephone`, `capacity` " +
		"FROM `warehouses` WHERE `id` = ?",
		id,
	)
	if err = row.Err(); err != nil {
		return
	}

	// scan the row into the warehouse
	err = row.Scan(&w.ID, &w.Name, &w.Address, &w.Telephone, &w.Capacity)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrWarehouseNotFound
		}
		return
	}

	return
}

// GetAll returns all warehouses
func (r *WarehousesMySQL) GetAll() (w []internal.Warehouse, err error) {
	// execute the query
	rows, err := r.db.Query(
		"SELECT `id`, `name`, `address`, `telephone`, `capacity` " +
		"FROM `warehouses`",
	)
	if err != nil {
		return
	}
	defer rows.Close()

	// scan the rows into the warehouses
	for rows.Next() {
		var warehouse internal.Warehouse
		err = rows.Scan(&warehouse.ID, &warehouse.Name, &warehouse.Address, &warehouse.Telephone, &warehouse.Capacity)
		if err != nil {
			return
		}

		w = append(w, warehouse)
	}

	return
}

// GetReportProducts returns a report of the amount of products in the warehouses
func (r *WarehousesMySQL) GetReportProducts(filter map[string]any) (w []internal.WarehouseReportProducts, err error) {
	// query: default
	query := "SELECT `w.name`, COUNT(`p.id`) AS `products_count` FROM `warehouses` AS `w` " +
			 "LEFT JOIN `products` AS `p` ON `w.id` = `p.warehouse_id` " +
			 "GROUP BY `w.id`"

	// query: build
	values := make([]any, 0)
	if len(filter) > 0 {
		// check id
		if id, ok := filter["id"]; ok {
			idInt, ok := id.(int)
			if !ok {
				err = internal.ErrWarehouseFilter
				return
			}

			query += " WHERE `w.id` = ?"
			values = append(values, idInt)
		}
	}

	// query: execute
	rows, err := r.db.Query(query, values...)
	if err != nil {
		return
	}

	// scan the rows into the warehouses
	for rows.Next() {
		var warehouse internal.WarehouseReportProducts
		err = rows.Scan(&warehouse.Name, &warehouse.ProductsCount)
		if err != nil {
			return
		}

		w = append(w, warehouse)
	}
	if err = rows.Err(); err != nil {
		return
	}

	return
}

// Store stores a warehouse
func (r *WarehousesMySQL) Store(w *internal.Warehouse) (err error) {
	// execute the query
	result, err := r.db.Exec(
		"INSERT INTO `warehouses` (`name`, `address`, `telephone`, `capacity`) VALUES (?, ?, ?, ?)",
		w.Name, w.Address, w.Telephone, w.Capacity,
	)
	if err != nil {
		return
	}

	// get the last inserted id
	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	// set the id
	(*w).ID = int(id)

	return
}
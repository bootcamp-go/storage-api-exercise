package handlers

import (
	"app/internal/warehouses/storage"
	"app/pkg/web/request"
	"app/pkg/web/response"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// NewControllerWarehouses creates new controller for warehouses
func NewControllerWarehouses(storage storage.StorageWarehouse) *ControllerWarehouses {
	return &ControllerWarehouses{storage: storage}
}

type ControllerWarehouses struct {
	// storage is an storage for warehouses
	storage storage.StorageWarehouse
}

// GetOne returns one warehouse by id
type ResponseWarehouseGetOne struct {
	Id			int    `json:"id"`
	Name		string `json:"name"`
	Address		string `json:"address"`
	Telephone	string `json:"telephone"`
	Capacity	int    `json:"capacity"`
}
type ResponseBodyWarehouseGetOne struct {
	Message string					 `json:"message"`
	Data    *ResponseWarehouseGetOne `json:"data"`
	Error   bool					 `json:"error"`
}
func (c *ControllerWarehouses) GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			code := http.StatusBadRequest
			body := &ResponseBodyWarehouseGetOne{Message: "invalid id", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}

		// process
		warehouse, err := c.storage.GetOne(id)
		if err != nil {
			var code int; var body *ResponseBodyWarehouseGetOne
			switch {
			case errors.Is(err, storage.ErrStorageWarehouseNotFound):
				code = http.StatusNotFound
				body = &ResponseBodyWarehouseGetOne{Message: "warehouse not found", Data: nil, Error: true}
			default:
				code = http.StatusInternalServerError
				body = &ResponseBodyWarehouseGetOne{Message: "internal error", Data: nil, Error: true}
			}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusOK
		body := &ResponseBodyWarehouseGetOne{
			Message: "success",
			Data: &ResponseWarehouseGetOne{
				Id:        warehouse.Id,
				Name:      warehouse.Name,
				Address:   warehouse.Address,
				Telephone: warehouse.Telephone,
				Capacity:  warehouse.Capacity,
			},
			Error: false,
		}

		response.JSON(w, code, body)
	}
}

// GetAll returns all warehouses
type ResponseWarehouseGetAll struct {
	Id			int    `json:"id"`
	Name		string `json:"name"`
	Address		string `json:"address"`
	Telephone	string `json:"telephone"`
	Capacity	int    `json:"capacity"`
}
type ResponseBodyWarehouseGetAll struct {
	Message string					 `json:"message"`
	Data    []*ResponseWarehouseGetAll `json:"data"`
	Error   bool					 `json:"error"`
}
func (c *ControllerWarehouses) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// ...

		// process
		warehouses, err := c.storage.GetAll()
		if err != nil {
			code := http.StatusInternalServerError
			body := &ResponseBodyWarehouseGetAll{Message: "internal error", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusOK
		body := &ResponseBodyWarehouseGetAll{
			Message: "success",
			Data:    make([]*ResponseWarehouseGetAll, 0),
			Error:   false,
		}
		for _, warehouse := range warehouses {
			body.Data = append(body.Data, &ResponseWarehouseGetAll{
				Id:        warehouse.Id,
				Name:      warehouse.Name,
				Address:   warehouse.Address,
				Telephone: warehouse.Telephone,
				Capacity:  warehouse.Capacity,
			})
		}

		response.JSON(w, code, body)
	}
}

// Create creates new warehouse
type RequestWarehouseCreate struct {
	Name		string `json:"name"`
	Address		string `json:"address"`
	Telephone	string `json:"telephone"`
	Capacity	int    `json:"capacity"`
}
type ResponseWarehouseCreate struct {
	Id			int    `json:"id"`
	Name		string `json:"name"`
	Address		string `json:"address"`
	Telephone	string `json:"telephone"`
	Capacity	int    `json:"capacity"`
}
type ResponseBodyWarehouseCreate struct {
	Message string					 `json:"message"`
	Data    *ResponseWarehouseCreate `json:"data"`
	Error   bool					 `json:"error"`
}
func (c *ControllerWarehouses) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		var req RequestWarehouseCreate
		if err := request.JSON(r, &req); err != nil {
			code := http.StatusBadRequest
			body := &ResponseBodyWarehouseCreate{Message: "invalid request", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}

		// process
		warehouse := &storage.Warehouse{
			Name:      req.Name,
			Address:   req.Address,
			Telephone: req.Telephone,
			Capacity:  req.Capacity,
		}
		if err := c.storage.Create(warehouse); err != nil {
			code := http.StatusInternalServerError
			body := &ResponseBodyWarehouseCreate{Message: "internal error", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusOK
		body := &ResponseBodyWarehouseCreate{
			Message: "success",
			Data: &ResponseWarehouseCreate{
				Id:        warehouse.Id,
				Name:      warehouse.Name,
				Address:   warehouse.Address,
				Telephone: warehouse.Telephone,
				Capacity:  warehouse.Capacity,
			},
			Error: false,
		}

		response.JSON(w, code, body)
	}
}

// ReportProducts returns report of products in warehouses
type ResponseReportProducts struct {
	WarehouseName string `json:"warehouse_name"`
	ProductsCount int    `json:"products_count"`
}
type ResponseBodyReportProducts struct {
	Message string					   `json:"message"`
	Data    []*ResponseReportProducts  `json:"data"`
	Error   bool					   `json:"error"`
}
func (c *ControllerWarehouses) ReportProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))

		// process
		report, err := c.storage.ReportProducts(id)
		if err != nil {
			code := http.StatusInternalServerError
			body := &ResponseBodyReportProducts{Message: "internal error", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusOK
		body := &ResponseBodyReportProducts{
			Message: "success",
			Data:    make([]*ResponseReportProducts, 0),
			Error:   false,
		}
		for _, r := range report {
			body.Data = append(body.Data, &ResponseReportProducts{
				WarehouseName: r.WarehouseName,
				ProductsCount: r.ProductsCount,
			})
		}

		response.JSON(w, code, body)
	}
}
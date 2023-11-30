package handler

import (
	"app/internal"
	"app/platform/web/request"
	"app/platform/web/response"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// NewWarehousesDefault returns a new instance of WarehousesDefault
func NewWarehousesDefault(rp internal.RepositoryWarehouses) *WarehousesDefault {
	return &WarehousesDefault{
		rp: rp,
	}
}

// WarehousesDefault is a struct that represents the default warehouse handler
type WarehousesDefault struct {
	// rp is the warehouse repository
	rp internal.RepositoryWarehouses
}

// WarehouseJSON is an struct that represents a warehouse in JSON format
type WarehouseJSON struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Address string `json:"address"`
	Telephone string `json:"telephone"`
	Capacity int `json:"capacity"`
}

// WarehouseReportProductsJSON is an struct that represents a warehouse report of products in JSON format
type WarehouseReportProductsJSON struct {
	Name string `json:"name"`
	ProductsCount int `json:"products_count"`
}

// GetOne returns a warehouse by id
func (h *WarehousesDefault) GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid id")
			return
		}

		// process
		wh, err := h.rp.GetOne(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrWarehouseNotFound):
				response.Error(w, http.StatusNotFound, "warehouse not found")
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		// - serialize
		data := WarehouseJSON{
			ID: wh.ID,
			Name: wh.Name,
			Address: wh.Address,
			Telephone: wh.Telephone,
			Capacity: wh.Capacity,
		}
		response.JSON(w, http.StatusOK, map[string]any{"message": "warehouse found", "data": data})
	}
}

// GetAll returns all warehouses
func (h *WarehousesDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// ...

		// process
		wh, err := h.rp.GetAll()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// response
		// - serialize
		var data []WarehouseJSON
		for _, v := range wh {
			data = append(data, WarehouseJSON{
				ID: v.ID,
				Name: v.Name,
				Address: v.Address,
				Telephone: v.Telephone,
				Capacity: v.Capacity,
			})
		}
		response.JSON(w, http.StatusOK, map[string]any{"message": "warehouses found", "data": data})
	}
}

// GetReportProducts returns a report of the amount of products in the warehouses
func (h *WarehousesDefault) GetReportProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		filter := make(map[string]any)
		query := r.URL.Query()
		for k, v := range query {
			filter[k] = v[0]
		}

		// process
		wh, err := h.rp.GetReportProducts(filter)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// response
		// - serialize
		var data []internal.WarehouseReportProducts
		for _, v := range wh {
			data = append(data, internal.WarehouseReportProducts{
				Name: v.Name,
				ProductsCount: v.ProductsCount,
			})
		}
		response.JSON(w, http.StatusOK, map[string]any{"message": "warehouses report of products found", "data": data})
	}
}

// RequestBodyWarehouseCreate is an struct that represents the request body of a warehouse to create
type RequestBodyWarehouseCreate struct {
	Name string `json:"name"`
	Address string `json:"address"`
	Telephone string `json:"telephone"`
	Capacity int `json:"capacity"`
}

// Create creates a warehouse
func (h *WarehousesDefault) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		var body RequestBodyWarehouseCreate
		if err := request.JSON(r, &body); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// process
		wh := internal.Warehouse{
			Name: body.Name,
			Address: body.Address,
			Telephone: body.Telephone,
			Capacity: body.Capacity,
		}
		if err := h.rp.Store(&wh); err != nil {
			switch {
			case errors.Is(err, internal.ErrWarehouseNotUnique):
				response.Error(w, http.StatusConflict, "warehouse not unique")
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		// - serialize
		data := WarehouseJSON{
			ID: wh.ID,
			Name: wh.Name,
			Address: wh.Address,
			Telephone: wh.Telephone,
			Capacity: wh.Capacity,
		}
		response.JSON(w, http.StatusCreated, map[string]any{"message": "warehouse created", "data": data})
	}
}
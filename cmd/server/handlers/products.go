package handlers

import (
	"app/internal/products/storage"
	"app/pkg/web/request"
	"app/pkg/web/response"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// NewControllerProduct returns new ControllerProduct
func NewControllerProduct(storage storage.StorageProduct) *ControllerProduct {
	return &ControllerProduct{storage: storage}
}

// ControllerProduct is a controller for products
type ControllerProduct struct {
	// storage is a storage for products
	storage storage.StorageProduct
}

// GetOne returns one product by id
type ResponseProduct struct {
	Name    	string		`json:"name"`
	Quantity	int			`json:"quantity"`
	CodeValue	string		`json:"code_value"`
	IsPublished bool		`json:"is_published"`
	Expiration  time.Time	`json:"expiration"`
	Price       float64		`json:"price"`
}
type ResponseBody struct {
	Message string			 `json:"message"`
	Data    *ResponseProduct `json:"data"`
	Error   bool			 `json:"error"`
}
func (c *ControllerProduct) GetOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			code := http.StatusBadRequest
			body := &ResponseBody{Message: "parameter must be int", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}

		// process
		product, err := c.storage.GetOne(id)
		if err != nil {
			var code int; var body *ResponseBody
			switch {
			case errors.Is(err, storage.ErrStorageProductNotFound):
				code = http.StatusNotFound
				body = &ResponseBody{Message: "product not found", Data: nil, Error: true}
			default:
				code = http.StatusInternalServerError
				body = &ResponseBody{Message: "internal error", Data: nil, Error: true}
			}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusOK
		body := &ResponseBody{
			Message: "success",
			Data: &ResponseProduct{			// serialization
				Name:   product.Name,
				Quantity: product.Quantity,
				CodeValue: product.CodeValue,
				IsPublished: product.IsPublished,
				Expiration: product.Expiration,
				Price: product.Price,
			},
			Error: false,
		}

		response.JSON(w, code, body)
	}
}

// Store stores product
type RequestProductStore struct {
	Name    	string		`json:"name"`
	Quantity	int			`json:"quantity"`
	CodeValue	string		`json:"code_value"`
	IsPublished bool		`json:"is_published"`
	Expiration  time.Time	`json:"expiration"`
	Price       float64		`json:"price"`
}
type ResponseProductStore struct {
	Name    	string		`json:"name"`
	Quantity	int			`json:"quantity"`
	CodeValue	string		`json:"code_value"`
	IsPublished bool		`json:"is_published"`
	Expiration  time.Time	`json:"expiration"`
	Price       float64		`json:"price"`
}
type ResponseBodyStore struct {
	Message string					`json:"message"`
	Data    *ResponseProductStore	`json:"data"`
	Error   bool					`json:"error"`
}
func (c *ControllerProduct) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		var req RequestProductStore
		err := request.JSON(r, &req)
		if err != nil {
			code := http.StatusBadRequest
			body := &ResponseBody{Message: "invalid json", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}

		// process
		// -> deserialization
		product := &storage.Product{
			Name:   req.Name,
			Quantity: req.Quantity,
			CodeValue: req.CodeValue,
			IsPublished: req.IsPublished,
			Expiration: req.Expiration,
			Price: req.Price,
		}
		err = c.storage.Store(product)
		if err != nil {
			var code int; var body *ResponseBody
			switch {
			case errors.Is(err, storage.ErrStorageProductNotUnique):
				code = http.StatusBadRequest
				body = &ResponseBody{Message: "product not unique", Data: nil, Error: true}
			default:
				code = http.StatusInternalServerError
				body = &ResponseBody{Message: "internal error", Data: nil, Error: true}
			}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusCreated
		body := &ResponseBodyStore{
			Message: "success",
			Data: &ResponseProductStore{	// serialization
				Name:   product.Name,
				Quantity: product.Quantity,
				CodeValue: product.CodeValue,
				IsPublished: product.IsPublished,
				Expiration: product.Expiration,
				Price: product.Price,
			},
			Error: false,
		}

		response.JSON(w, code, body)
	}
}


// Update updates product
type RequestProductUpdate struct {
	Name    	string	`json:"name"`
	Quantity	int		`json:"quantity"`
	CodeValue	string	`json:"code_value"`
	IsPublished bool	`json:"is_published"`
	Expiration  time.Time	`json:"expiration"`
	Price       float64	`json:"price"`
}
type ResponseProductUpdate struct {
	Name    	string	`json:"name"`
	Quantity	int		`json:"quantity"`
	CodeValue	string	`json:"code_value"`
	IsPublished bool	`json:"is_published"`
	Expiration  time.Time	`json:"expiration"`
	Price       float64	`json:"price"`
}
type ResponseBodyUpdate struct {
	Message string					`json:"message"`
	Data    *ResponseProductUpdate	`json:"data"`
	Error   bool					`json:"error"`
}
func (c *ControllerProduct) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			code := http.StatusBadRequest
			body := &ResponseBody{Message: "parameter must be int", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}

		// process
		// -> get searched product by id
		p, err := c.storage.GetOne(id)
		if err != nil {
			var code int; var body *ResponseBodyUpdate
			switch {
			case errors.Is(err, storage.ErrStorageProductNotFound):
				code = http.StatusNotFound
				body = &ResponseBodyUpdate{Message: "product not found", Data: nil, Error: true}
			default:
				code = http.StatusInternalServerError
				body = &ResponseBodyUpdate{Message: "internal error", Data: nil, Error: true}
			}

			response.JSON(w, code, body)
			return
		}
		// -- serialization
		req := &RequestProductUpdate{
			Name:   p.Name,
			Quantity: p.Quantity,
			CodeValue: p.CodeValue,
			IsPublished: p.IsPublished,
			Expiration: p.Expiration,
			Price: p.Price,
		}

		// -> patch product to RequestProductUpdate(filled with original data)
		err = request.JSON(r, req)
		if err != nil {
			code := http.StatusBadRequest
			body := &ResponseBody{Message: "invalid json", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}
		// -- deserialization
		product := &storage.Product{
			Id: id,
			Name:   req.Name,
			Quantity: req.Quantity,
			CodeValue: req.CodeValue,
			Expiration: req.Expiration,
			Price: req.Price,
		}
		// -- update product
		err = c.storage.Update(product)
		if err != nil {
			var code int; var body *ResponseBody
			switch {
			case errors.Is(err, storage.ErrStorageProductNotFound):
				code = http.StatusNotFound
				body = &ResponseBody{Message: "product not found", Data: nil, Error: true}
			case errors.Is(err, storage.ErrStorageProductNotUnique):
				code = http.StatusBadRequest
				body = &ResponseBody{Message: "product not unique", Data: nil, Error: true}
			default:
				code = http.StatusInternalServerError
				body = &ResponseBody{Message: "internal error", Data: nil, Error: true}
			}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusOK
		body := &ResponseBodyUpdate{
			Message: "success",
			Data: &ResponseProductUpdate{	// serialization
				Name:   product.Name,
				Quantity: product.Quantity,
				CodeValue: product.CodeValue,
				IsPublished: product.IsPublished,
				Expiration: product.Expiration,
				Price: product.Price,
			},
			Error: false,
		}

		response.JSON(w, code, body)
	}
}


// Delete deletes product by id
type ResponseBodyDelete struct {
	Message string	`json:"message"`
	Data    any		`json:"data"`
	Error   bool	`json:"error"`
}
func (c *ControllerProduct) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			code := http.StatusBadRequest
			body := &ResponseBody{Message: "parameter must be int", Data: nil, Error: true}

			response.JSON(w, code, body)
			return
		}

		// process
		// -> delete product by id
		err = c.storage.Delete(id)
		if err != nil {
			var code int; var body *ResponseBody
			switch {
			case errors.Is(err, storage.ErrStorageProductNotFound):
				code = http.StatusNotFound
				body = &ResponseBody{Message: "product not found", Data: nil, Error: true}
			default:
				code = http.StatusInternalServerError
				body = &ResponseBody{Message: "internal error", Data: nil, Error: true}
			}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusNoContent
		body := any(nil)

		response.JSON(w, code, body)
	}
}
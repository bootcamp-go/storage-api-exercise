package handler_test

import (
	"app/internal/handler"
	"app/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

// TestWarehousesDefault_GetOne tests the GetOne method
func TestWarehousesDefault_GetOne(t *testing.T) {
	t.Run("case 1: success - returns a warehouse", func(t *testing.T) {
		// arrange
		// - database: connection
		db, err := sql.Open("txdb", "")
		require.NoError(t, err)
		defer db.Close()
		// - database: teardown
		defer func(db *sql.DB) {
			// delete warehouses
			_, err = db.Exec("DELETE FROM `warehouses`;")
			if err != nil {
				panic(err)
			}
			_, err = db.Exec("ALTER TABLE `warehouses` AUTO_INCREMENT = 1;")
			if err != nil {
				panic(err)
			}
		}(db)
		// - database: setup
		err = func(db *sql.DB) error {
			// insert warehouses
			_, err := db.Exec(
				"INSERT INTO `warehouses` (`id`, `name`, `address`, `telephone`, `capacity`) VALUES " +
				"(1, 'warehouse 1', 'address 1', 'telephone 1', 100)," +
				"(2, 'warehouse 2', 'address 2', 'telephone 2', 200);",
			)
			if err != nil {
				return err
			}
			return nil
		}(db)
		require.NoError(t, err)

		// - repository: mysql
		rp := repository.NewWarehousesMySQL(db)
		// - handler: default
		hd := handler.NewWarehousesDefault(rp)
		hdFunc := hd.GetOne()

		// act
		// - request
		request := httptest.NewRequest(http.MethodGet, "/warehouses/1", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx))
		// - response
		response := httptest.NewRecorder()
		hdFunc(response, request)

		// assert
		expectedCode := http.StatusOK
		expectedBody := `{"message":"warehouse found","data":{"id":1,"name":"warehouse 1","address":"address 1","telephone":"telephone 1","capacity":100}}`
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})

	t.Run("case 2: failure - warehouse not found", func(t *testing.T) {
		// arrange
		// - database: connection
		db, err := sql.Open("txdb", "")
		require.NoError(t, err)
		defer db.Close()
		// - database: teardown
		defer func(db *sql.DB) {
			// delete warehouses
			_, err = db.Exec("DELETE FROM `warehouses`;")
			if err != nil {
				panic(err)
			}
			_, err = db.Exec("ALTER TABLE `warehouses` AUTO_INCREMENT = 1;")
			if err != nil {
				panic(err)
			}
		}(db)
		// - database: setup
		// ...

		// - repository: mysql
		rp := repository.NewWarehousesMySQL(db)
		// - handler: default
		hd := handler.NewWarehousesDefault(rp)
		hdFunc := hd.GetOne()

		// act
		// - request
		request := httptest.NewRequest(http.MethodGet, "/warehouses/1", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "1")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx))
		// - response
		response := httptest.NewRecorder()
		hdFunc(response, request)

		// assert
		expectedCode := http.StatusNotFound
		expectedBody := fmt.Sprintf(
			`{"status":"%s","message":"%s"}`,
			http.StatusText(expectedCode),
			"warehouse not found",
		)
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})

	t.Run("case 3: failure - invalid id", func(t *testing.T) {
		// arrange
		// - database: connection
		// - database: teardown
		// - database: setup

		// - repository: mysql
		// - handler: default
		hd := handler.NewWarehousesDefault(nil)
		hdFunc := hd.GetOne()

		// act
		// - request
		request := httptest.NewRequest(http.MethodGet, "/warehouses/invalid", nil)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("id", "invalid")
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, chiCtx))
		// - response
		response := httptest.NewRecorder()
		hdFunc(response, request)

		// assert
		expectedCode := http.StatusBadRequest
		expectedBody := fmt.Sprintf(
			`{"status":"%s","message":"%s"}`,
			http.StatusText(expectedCode),
			"invalid id",
		)
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})
}
// TestWarehousesDefault_GetAll tests the GetAll method
func TestWarehousesDefault_GetAll(t *testing.T) {
	t.Run("case 1: success - returns all warehouses", func(t *testing.T) {
		// arrange
		// - database: connection
		db, err := sql.Open("txdb", "")
		require.NoError(t, err)
		defer db.Close()
		// - database: teardown
		defer func(db *sql.DB) {
			// delete warehouses
			_, err = db.Exec("DELETE FROM `warehouses`;")
			if err != nil {
				panic(err)
			}
			_, err = db.Exec("ALTER TABLE `warehouses` AUTO_INCREMENT = 1;")
			if err != nil {
				panic(err)
			}
		}(db)
		// - database: setup
		err = func(db *sql.DB) error {
			// insert warehouses
			_, err := db.Exec(
				"INSERT INTO `warehouses` (`id`, `name`, `address`, `telephone`, `capacity`) VALUES " +
				"(1, 'warehouse 1', 'address 1', 'telephone 1', 100)," +
				"(2, 'warehouse 2', 'address 2', 'telephone 2', 200);",
			)
			if err != nil {
				return err
			}
			return nil
		}(db)
		require.NoError(t, err)

		// - repository: mysql
		rp := repository.NewWarehousesMySQL(db)
		// - handler: default
		hd := handler.NewWarehousesDefault(rp)
		hdFunc := hd.GetAll()

		// act
		request := httptest.NewRequest(http.MethodGet, "/warehouses", nil)
		response := httptest.NewRecorder()
		hdFunc(response, request)

		// assert
		expectedCode := http.StatusOK
		expectedBody := `{"message":"warehouses found","data":[{"id":1,"name":"warehouse 1","address":"address 1","telephone":"telephone 1","capacity":100},{"id":2,"name":"warehouse 2","address":"address 2","telephone":"telephone 2","capacity":200}]}`
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})

	t.Run("case 2: success - returns empty array", func(t *testing.T) {
		// arrange
		// - database: connection
		db, err := sql.Open("txdb", "")
		require.NoError(t, err)
		defer db.Close()
		// - database: teardown
		// ...
		// - database: setup
		// ...

		// - repository: mysql
		rp := repository.NewWarehousesMySQL(db)
		// - handler: default
		hd := handler.NewWarehousesDefault(rp)
		hdFunc := hd.GetAll()

		// act
		request := httptest.NewRequest(http.MethodGet, "/warehouses", nil)
		response := httptest.NewRecorder()
		hdFunc(response, request)

		// assert
		expectedCode := http.StatusOK
		expectedBody := `{"message":"warehouses found","data":[]}`
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})
}
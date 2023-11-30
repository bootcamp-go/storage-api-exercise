package handler_test

import (
	"app/internal/handler"
	"app/internal/repository"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

// init registers txdb driver
func init() {
	// db connection
	cfg := mysql.Config{
		User:      "root",
		Passwd:    "",
		Net:       "tcp",
		Addr:      "127.0.0.1:3306",
		DBName:    "storage_api_test_db",
		ParseTime: true,
	}
	// register txdb driver
	txdb.Register("txdb", "mysql", cfg.FormatDSN())
}

// TestProductsDefault_GetAll tests the GetAll method
func TestProductsDefault_GetAll(t *testing.T) {
	t.Run("case 1: success - returns all products", func(t *testing.T) {
		// arrange
		// - database: connection
		db, err := sql.Open("txdb", "")
		require.NoError(t, err)
		defer db.Close()
		// - database: teardown
		defer func(db *sql.DB) {
			// delete products
			_, err := db.Exec("DELETE FROM `products`;")
			if err != nil {
				panic(err)
			}
			// delete warehouses
			_, err = db.Exec("DELETE FROM `warehouses`;")
			if err != nil {
				panic(err)
			}
			// reset auto increment
			_, err = db.Exec("ALTER TABLE `products` AUTO_INCREMENT = 1;")
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
					"(1, 'warehouse 1', 'address 1', 'telephone 1', 100);",
			)
			if err != nil {
				return err
			}
			// insert products
			_, err = db.Exec(
				"INSERT INTO `products` (`id`, `name`, `quantity`, `code_value`, `is_published`, `expiration`, `price`, `warehouse_id`) VALUES" +
					"(1, 'product 1', 10, 'code 1', 1, '2021-01-01', 10.00, 1)," +
					"(2, 'product 2', 20, 'code 2', 1, '2021-01-02', 20.00, 1);",
			)
			if err != nil {
				return err
			}
			return nil
		}(db)
		require.NoError(t, err)

		// - repository: mysql
		rp := repository.NewProductsMySQL(db)
		// - handler: default
		hd := handler.NewProductsDefault(rp)
		hdFunc := hd.GetAll()

		// act
		request := httptest.NewRequest(http.MethodGet, "/products", nil)
		response := httptest.NewRecorder()
		hdFunc(response, request)

		// assert
		expectedCode := http.StatusOK
		expectedBody := `{"message":"products found","data":[{"id":1,"name":"product 1","quantity":10,"code_value":"code 1","is_published":true,"expiration":"2021-01-01","price":10,"warehouse_id":1},{"id":2,"name":"product 2","quantity":20,"code_value":"code 2","is_published":true,"expiration":"2021-01-02","price":20,"warehouse_id":1}]}`
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
		rp := repository.NewProductsMySQL(db)
		// - handler: default
		hd := handler.NewProductsDefault(rp)
		hdFunc := hd.GetAll()

		// act
		request := httptest.NewRequest(http.MethodGet, "/products", nil)
		response := httptest.NewRecorder()
		hdFunc(response, request)

		// assert
		expectedCode := http.StatusOK
		expectedBody := `{"message":"products found","data":[]}`
		require.Equal(t, expectedCode, response.Code)
		require.JSONEq(t, expectedBody, response.Body.String())
	})
}

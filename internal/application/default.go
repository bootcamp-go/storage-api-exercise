package application

import (
	"app/internal/handler"
	"app/internal/repository"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-sql-driver/mysql"
)

// ConfigDefault is a struct that represents the default application configuration
type ConfigDefault struct {
	// Database is the database configuration
	Database mysql.Config
	// Address is the address of the application
	Address string
}

// NewDefault returns a new default application
func NewDefault(cfg *ConfigDefault) *Default {
	// default
	cfgDefault := &ConfigDefault{
		Address: ":8080",
	}
	if cfg != nil {
		cfgDefault.Database = cfg.Database
		if cfg.Address != "" {
			cfgDefault.Address = cfg.Address
		}
	}

	return &Default{
		cfgDb: cfgDefault.Database,
		addr:  cfgDefault.Address,
	}
}

// Default is a struct that represents the default application
type Default struct {
	// cfgDb is the database configuration
	cfgDb mysql.Config
	// addr is the address of the application
	addr string
}

// Run runs the default application
func (d *Default) Run() (err error) {
	// dependencies
	// - database: connection
	db, err := sql.Open("mysql", d.cfgDb.FormatDSN())
	if err != nil {
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		return
	}
	
	// - repository: default
	rpPr := repository.NewProductsMySQL(db)
	rpWh := repository.NewWarehousesMySQL(db)
	
	// - handler: default
	hdPr := handler.NewProductsDefault(rpPr)
	hdWh := handler.NewWarehousesDefault(rpWh)

	// - router: chi
	rt := chi.NewRouter()
	// - router: middlewares
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)
	// - router: routes
	rt.Route("/products", func(r chi.Router) {
		// - GET /products/{id}
		r.Get("/{id}", hdPr.GetOne())
		// - GET /products
		r.Get("/", hdPr.GetAll())
		// - POST /products
		r.Post("/", hdPr.Create())
		// - PUT /products/{id}
		r.Patch("/{id}", hdPr.Update())
		// - DELETE /products/{id}
		r.Delete("/{id}", hdPr.Delete())
	})
	rt.Route("/warehouses", func(r chi.Router) {
		// - GET /warehouses/{id}
		r.Get("/{id}", hdWh.GetOne())
		// - GET /warehouses
		r.Get("/", hdWh.GetAll())
		// - GET /warehouses/report-products
		r.Get("/reportProducts", hdWh.GetReportProducts())
		// - POST /warehouses
		r.Post("/", hdWh.Create())
	})

	// run
	err = http.ListenAndServe(d.addr, rt)
	if err != nil {
		return
	}
	return
}

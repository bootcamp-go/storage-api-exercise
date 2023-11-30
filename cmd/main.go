package main

import (
	"app/internal/application"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

func main() {
	// env
	// ...

	// application
	// - config
	cfg := &application.ConfigDefault{
		Database: mysql.Config{
			User:      "root",
			Passwd:    "",
			Net:       "tcp",
			Addr:      "127.0.0.1:3306",
			DBName:    "storage_api_db",
			ParseTime: true,
		},
		Address: "127.0.0.1:8080",
	}
	app := application.NewDefault(cfg)
	// - run
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}

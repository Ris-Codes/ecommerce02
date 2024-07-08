package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	var err error

	dsn := os.Getenv("DB_URL")
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Creating Tables

	// createTable("admin")
	// createTable("users")
	// createTable("address")
	// createTable("categories")
	// createTable("products")
	// createTable("cart")
	// createTable("orders_status")
	// createTable("orders")
	// createTable("order_items")

}

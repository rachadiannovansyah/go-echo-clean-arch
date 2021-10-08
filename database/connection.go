package database

import (
	"database/sql"
	"log"

	config "../config"
)

// InitDB ...
func InitDB(cfg *config.Config) *sql.DB {
	dbConn, err := sql.Open(`mysql`, cfg.DB.DSN)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return dbConn
}

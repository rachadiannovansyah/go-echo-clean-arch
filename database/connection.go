package database

import (
	"database/sql"
	"log"

	"github.com/rachadiannovansyah/go-echo-clean-arch/config"
)

// InitDB ...
func InitDB(cfg *config.Config) *sql.DB {
	dbConn, err := sql.Open("mysql", cfg.DB.DSN)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return dbConn
}

package main

import (
	"database/sql"

	"github.com/MhdFiras-3/gofeed/internal/database"
)

type apiConfig struct {
	db     *database.Queries
	dbConn *sql.DB
}

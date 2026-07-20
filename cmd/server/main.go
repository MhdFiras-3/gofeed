package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/MhdFiras-3/gofeed/internal/database"
	"github.com/MhdFiras-3/gofeed/internal/handlers"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	dbConnection, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("err connecting to database: %v", err)
	}
	dbQueries := database.New(dbConnection)

	apicfg := &handlers.APIConfig{
		DB:     dbQueries,
		DBConn: dbConnection,
	}

}

package handlers

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/MhdFiras-3/gofeed/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var testDB *sql.DB
var testCfg *APIConfig

func TestMain(m *testing.M) {
	godotenv.Load()

	testDBURL := os.Getenv("TEST_DB_URL")

	db, err := sql.Open("postgres", testDBURL)
	if err != nil {
		log.Fatalf("failed to open test db: %v", err)
	}
	testDB = db

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}
	if err := goose.Up(testDB, "sql/migrations"); err != nil {
		log.Fatalf("failed to run up migrations: %v", err)
	}

	testCfg = &APIConfig{
		DB:     database.New(testDB),
		DBConn: testDB,
	}

	code := m.Run()

	if err := goose.DownTo(testDB, "sql/migrations", 0); err != nil {
		log.Fatalf("failed to run down migrations: %v", err)
	}

	testDB.Close()
	os.Exit(code)
}

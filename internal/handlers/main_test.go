package handlers

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/MhdFiras-3/gofeed/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

var testCfg *APIConfig

func TestMain(m *testing.M) {
	if err := godotenv.Load(filepath.Join(projectRootDir(), ".env")); err != nil {
		log.Printf("could not load .env: %v", err)
	}

	testDBURL := os.Getenv("TEST_DB_URL")
	log.Printf("using DB URL: %s", testDBURL)

	db, err := sql.Open("postgres", testDBURL)
	if err != nil {
		log.Fatalf("failed to open test db: %v", err)
	}
	migrationsPath := filepath.Join(projectRootDir(), "sql/migrations")
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}
	if err := goose.Up(db, migrationsPath); err != nil {
		log.Fatalf("failed to run up migrations: %v", err)
	}

	testCfg = &APIConfig{
		DB:     database.New(db),
		DBConn: db,
	}

	code := m.Run()

	if err := goose.DownTo(db, migrationsPath, 0); err != nil {
		log.Fatalf("failed to run down migrations: %v", err)
	}

	db.Close()
	os.Exit(code)
}

func projectRootDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..")
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MhdFiras-3/gofeed/internal/database"
	"github.com/MhdFiras-3/gofeed/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	const port = "8080"
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
		DB:        dbQueries,
		DBConn:    dbConnection,
		JWTSecret: os.Getenv("JWT_SECRET"),
		JWTExpiry: time.Hour,
	}

	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/users", apicfg.HandlerCreateUser)
		r.Post("/login", apicfg.HandlerLogin)
		r.Post("/refresh", apicfg.HandlerRefresh)

	})
	fmt.Printf("serving on %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

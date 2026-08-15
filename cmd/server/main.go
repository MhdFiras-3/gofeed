package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MhdFiras-3/gofeed/internal/database"
	"github.com/MhdFiras-3/gofeed/internal/handlers"
	"github.com/MhdFiras-3/gofeed/internal/scraper"

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
		Ticker:    10 * time.Second,
	}

	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", apicfg.HandlerCreateUser)
		r.Post("/login", apicfg.HandlerLogin)
		r.Post("/refresh", apicfg.HandlerRefresh)
		r.Post("/logout", apicfg.HandlerLogOut)
		r.Get("/feeds", apicfg.HandlerGetAllFeeds)
		r.Get("/feeds/{feedID}", apicfg.HandlerGetFeedByID)

		r.Group(func(r chi.Router) {
			r.Use(apicfg.MiddlewareAuth)
			r.Get("/me", apicfg.HandlerGetCurrentUser)
			r.Patch("/me", apicfg.HandlerUpdateUser)
			r.Delete("/me", apicfg.HandlerDeleteUser)
			r.Post("/feeds", apicfg.HandlerCreateFeed)
			r.Get("/follows", apicfg.HandlerGetFeedFollows)
			r.Delete("/follows/{feedID}", apicfg.HandlerDeleteFeedFollow)
			r.Get("/posts", apicfg.HandlerGetPostsForUser)
			r.Post("/posts/{postID}/read", apicfg.HandlerMarkPostRead)
			r.Get("posts/read", apicfg.HandlerGetReadPostsForUser)
		})

	})
	go scraper.StartScraping(context.Background(), apicfg.DB, apicfg.Ticker)
	fmt.Printf("serving on %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

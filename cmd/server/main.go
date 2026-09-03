package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	defer dbConnection.Close()

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
			r.Post("/follows", apicfg.HandlerCreateFeedFollow)
			r.Delete("/follows/{feedID}", apicfg.HandlerDeleteFeedFollow)
			r.Get("/posts", apicfg.HandlerGetPostsForUser)
			r.Post("/posts/{postID}/read", apicfg.HandlerMarkPostRead)
			r.Get("/posts/read", apicfg.HandlerGetReadPostsForUser)
		})

	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()
	go scraper.StartScraping(ctx, apicfg.DB, apicfg.Ticker)
	go func() {
		fmt.Printf("serving on %s\n", port)
		log.Fatal(server.ListenAndServe())
	}()
	<-ctx.Done()
	log.Println("shutdown signal recieved, exiting...")

	shutdownCtx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server forced to shutdown, %v", err)
	}
}

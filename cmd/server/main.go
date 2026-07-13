package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MhdFiras-3/gofeed/internal/auth"
	"github.com/MhdFiras-3/gofeed/internal/database"
	"github.com/google/uuid"

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

	apicfg := &apiConfig{
		db:     dbQueries,
		dbConn: dbConnection,
	}

}
func rotateRefreshToken(api *apiConfig, userID uuid.UUID, refreshtoken string) (database.RefreshToken, error) {
	tx, err := api.dbConn.BeginTx(context.Background(), &sql.TxOptions{})

	if err != nil {
		return database.RefreshToken{}, err
	}
	defer tx.Rollback()
	dbTx := api.db.WithTx(tx)

	err = dbTx.RevokeRefreshToken(context.Background(), refreshtoken)
	if err != nil {
		return database.RefreshToken{}, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}
	token, err := auth.MakeRefreshToken()
	if err != nil {
		return database.RefreshToken{}, err
	}
	newRefreshToken, err := dbTx.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		RefreshToken: token,
		ExpiresAt:    time.Now().Add(24 * time.Hour * 7),
		UserID:       userID,
	})

	if err != nil {
		return database.RefreshToken{}, err
	}
	err = tx.Commit()
	if err != nil {
		return database.RefreshToken{}, err
	}

	return newRefreshToken, nil
}

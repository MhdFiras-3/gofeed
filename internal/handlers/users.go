package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MhdFiras-3/gofeed/internal/auth"
	"github.com/MhdFiras-3/gofeed/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserLoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func rotateRefreshToken(api *APIConfig, userID uuid.UUID, refreshtoken string) (database.RefreshToken, error) {
	tx, err := api.DBConn.BeginTx(context.Background(), &sql.TxOptions{})

	if err != nil {
		return database.RefreshToken{}, err
	}
	defer tx.Rollback()
	dbTx := api.DB.WithTx(tx)

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

func (cfg *APIConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type requestParam struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type responseErr struct {
		Errors []inputError `json:"errors"`
	}
	var reqData requestParam
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqData); err != nil {
		respWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	if errs := validateInput(reqData.Name, reqData.Email, reqData.Password); len(errs) > 0 {
		respWithJson(w, http.StatusBadRequest, responseErr{Errors: errs})
		return
	}
	hashedPassword, err := auth.HashPassword(reqData.Password)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	createdUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Name:           reqData.Name,
		Email:          reqData.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			respWithError(w, http.StatusConflict, "email already registered")
			return
		}
		respWithError(w, http.StatusBadRequest, "something went wrong")
		return
	}
	respWithJson(w, http.StatusCreated, User{
		ID:        createdUser.ID,
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	})
}

func (cfg *APIConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type requestParam struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		AccessToken  string
		RefreshToken string
	}
	var reqData requestParam
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqData); err != nil {
		respWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	DBUser, err := cfg.DB.GetUserByEmail(r.Context(), reqData.Email)
	if err != nil {
		respWithError(w, http.StatusUnauthorized, "wrong email or password")
	}
	match, err := auth.ValidatePasswordHash(reqData.Password, DBUser.HashedPassword)
	if !match || err != nil {
		respWithError(w, http.StatusUnauthorized, "wrong email or password")
		return
	}
	accessToken, err := auth.MakeJWT(DBUser.ID, cfg.JWTSecret, time.Hour)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "failed to issue token")
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "failed to issue token")
		return
	}
	_, err = cfg.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Hour * 24 * 60),
		UserID:       DBUser.ID,
	})
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "failed to save refresh token")
		return
	}

	respWithJson(w, http.StatusOK, response{
		User: User{
			ID:        DBUser.ID,
			Name:      DBUser.Name,
			Email:     DBUser.Email,
			CreatedAt: DBUser.CreatedAt,
			UpdatedAt: DBUser.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

}

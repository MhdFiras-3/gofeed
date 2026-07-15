package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/MhdFiras-3/gofeed/internal/auth"
	"github.com/MhdFiras-3/gofeed/internal/database"
	"github.com/google/uuid"
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

func (cfg *apiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type requestParam struct {
		Name     string
		Email    string
		Password string
	}
	var reqData requestParam
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqData); err != nil {
		respWithError(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	hashedPassword, err := auth.HashPassword(reqData.Password)
	if err != nil {
		respWithError(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	createdUser, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		Name:           reqData.Name,
		Email:          reqData.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respWithError(w, "something went wrong", http.StatusInternalServerError)
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

func (cfg *apiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {

}

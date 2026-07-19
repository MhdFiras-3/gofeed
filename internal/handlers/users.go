package handlers

import (
	"encoding/json"
	"errors"
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

func (cfg *apiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type requestParam struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type errResp struct {
		Errors []inputError `json:"errors"`
	}
	var reqData requestParam
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqData); err != nil {
		respWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	if errs := validateInput(reqData.Name, reqData.Email, reqData.Password); len(errs) > 0 {
		respWithJson(w, http.StatusBadRequest, errResp{Errors: errs})
		return
	}
	hashedPassword, err := auth.HashPassword(reqData.Password)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	createdUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
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

func (cfg *apiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {

}

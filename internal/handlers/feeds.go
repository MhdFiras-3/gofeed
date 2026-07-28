package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/MhdFiras-3/gofeed/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (cfg *APIConfig) HandlerCreateFeed(w http.ResponseWriter, r *http.Request) {
	type requestParam struct {
		URL  *string `json:"url"`
		Name *string `json:"name"`
	}
	var reqData requestParam
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqData); err != nil {
		respWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if reqData.URL == nil || reqData.Name == nil {
		respWithError(w, http.StatusBadRequest, "URL must not be empty")
		return
	}
	err := parseURL(*reqData.URL)
	if err != nil {
		respWithError(w, http.StatusBadRequest, "invalid URL")
		return
	}
	err = parseName(*reqData.Name)
	if err != nil {
		respWithError(w, http.StatusBadRequest, "invalid feed name")
		return
	}

	_, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		respWithError(w, http.StatusInternalServerError, "failed missing user id in context")
		log.Println("failed to get user id from context to create feed")
		return
	}

	_, err = cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		Url:  *reqData.URL,
		Name: *reqData.Name,
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			respWithError(w, http.StatusConflict, "feed already exist")
			log.Printf("failed to create feed: %v", err)
			return
		}
		respWithError(w, http.StatusInternalServerError, "something went wrong")
		log.Printf("failed to create feed: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

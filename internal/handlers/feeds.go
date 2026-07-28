package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MhdFiras-3/gofeed/internal/database"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandlerCreateFeed(w http.ResponseWriter, r *http.Request) {
	type requestParam struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	}
	var reqData requestParam
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&reqData); err != nil {
		respWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err := parseURL(reqData.URL)
	if err != nil {
		respWithError(w, http.StatusBadRequest, "invalid URL")
		return
	}
	err = parseName(reqData.Name)
	if err != nil {
		respWithError(w, http.StatusBadRequest, "invalid feed name")
		return
	}

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		respWithError(w, http.StatusInternalServerError, "failed missing user id in context")
		log.Println("failed to get user id from context to create feed")
		return
	}

	feedDB, err := cfg.DB.CreateFeed(r.Context(), reqData.URL)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "something went wrong")
		log.Printf("failed to create feed: %v", err)
		return
	}
	_, err = cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		UserID: userID,
		FeedID: feedDB.ID,
		Name:   reqData.Name,
	})
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "something went wrong")
		log.Printf("failed to create feed follow: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

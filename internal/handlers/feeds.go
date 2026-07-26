package handlers

import (
	"log"
	"net/http"

	"github.com/MhdFiras-3/gofeed/internal/database"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandlerCreateFeed(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		respWithError(w, http.StatusInternalServerError, "failed missing user id in context")
		log.Println("failed to get user id from context to create feed")
		return
	}

	cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{})
}

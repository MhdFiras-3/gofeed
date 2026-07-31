package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/MhdFiras-3/gofeed/internal/database"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandlerGetPostsForUser(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID          uuid.UUID  `json:"id"`
		Title       string     `json:"title"`
		Url         string     `json:"url"`
		Description *string    `json:"description"`
		FeedID      uuid.UUID  `json:"feed_id"`
		CreatedAt   time.Time  `json:"created_at"`
		UpdatedAt   time.Time  `json:"updated_at"`
		PublishedAt *time.Time `json:"published_at"`
	}

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		respWithError(w, http.StatusInternalServerError, "missing user id in context")
		log.Println("failed to get user id from context to get posts for user")
		return
	}

	limit := parseQueryParamLimit(r.URL.Query().Get("limit"))
	postsDB, err := cfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{UserID: userID, Limit: limit})
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "failed to get posts")
		log.Printf("failed to get posts for user from database: %v", err)
		return
	}

	responses := make([]response, 0, len(postsDB))
	for _, post := range postsDB {
		var description *string
		var publishedAt *time.Time
		if post.Description.Valid {
			description = &post.Description.String
		}
		if post.PublishedAt.Valid {
			publishedAt = &post.PublishedAt.Time
		}
		responses = append(responses, response{
			ID:          post.ID,
			Title:       post.Title,
			Url:         post.Url,
			Description: description,
			FeedID:      post.FeedID,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
			PublishedAt: publishedAt,
		})
	}
	respWithJson(w, http.StatusOK, responses)
}

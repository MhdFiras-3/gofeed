package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/MhdFiras-3/gofeed/internal/database"
	"github.com/google/uuid"
)

func (cfg *APIConfig) HandlerCreateFeed(w http.ResponseWriter, r *http.Request) {
	type requestParam struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		UserID    uuid.UUID `json:"user_id"`
		FeedID    uuid.UUID `json:"feed_id"`
		Name      string    `json:"name"`
		UserName  string    `json:"user_name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
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
	tx, err := cfg.DBConn.BeginTx(r.Context(), &sql.TxOptions{})
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "something went wrong")
		log.Printf("failed to begin tx: %v", err)
		return
	}
	defer tx.Rollback()
	dbTx := cfg.DB.WithTx(tx)
	feedDB, err := dbTx.CreateFeed(r.Context(), reqData.URL)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "something went wrong")
		log.Printf("failed to create feed: %v", err)
		return
	}
	feedFollowDB, err := dbTx.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		UserID: userID,
		FeedID: feedDB.ID,
		Name:   reqData.Name,
	})
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "something went wrong")
		log.Printf("failed to create feed follow: %v", err)
		return
	}
	if err := tx.Commit(); err != nil {
		respWithError(w, http.StatusInternalServerError, "something went wrong")
		log.Printf("failed to commit tx: %v", err)
		return
	}
	respWithJson(w, http.StatusCreated, response{
		ID:        feedFollowDB.ID,
		UserID:    feedFollowDB.UserID,
		FeedID:    feedFollowDB.FeedID,
		Name:      feedFollowDB.Name,
		UserName:  feedFollowDB.UserName,
		CreatedAt: feedFollowDB.CreatedAt,
		UpdatedAt: feedFollowDB.UpdatedAt,
	})
}

func (cfg *APIConfig) HandlerGetFeedFollows(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID        uuid.UUID `json:"id"`
		UserID    uuid.UUID `json:"user_id"`
		FeedID    uuid.UUID `json:"feed_id"`
		Name      string    `json:"name"`
		UserName  string    `json:"user_name"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		respWithError(w, http.StatusInternalServerError, "failed missing user id in context")
		log.Println("failed to get user id from context to get feed follows")
		return
	}
	userFeedFollows, err := cfg.DB.GetFeedFollowsForUser(r.Context(), userID)
	if err != nil {
		respWithError(w, http.StatusInternalServerError, "something went wrong")
		log.Printf("failed to get feed follows from database: %v", err)
		return
	}
	responses := make([]response, 0, len(userFeedFollows))

	for _, userFollows := range userFeedFollows {
		responses = append(responses, response{
			ID:        userFollows.ID,
			UserID:    userFollows.UserID,
			FeedID:    userFollows.FeedID,
			Name:      userFollows.Name,
			UserName:  userFollows.UserName,
			CreatedAt: userFollows.CreatedAt,
			UpdatedAt: userFollows.UpdatedAt,
		})
	}

	respWithJson(w, http.StatusOK, responses)

}

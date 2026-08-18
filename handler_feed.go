package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mata-codes/rssagg/internal/database"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name   string        `json:"name"`
		Url    string        `json:"url"`
		UserID uuid.NullUUID `json:"user_id"`
	}
	var decoder json.Decoder = *json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json %v", err))
		return
	}
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create feed: %v", err))
		return
	}
	respondWithJson(w, 201, databaseToFeed(feed))
}
func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get feeds: %v", err))
		return
	}
	respondWithJson(w, 201, databaseToFeeds(feeds))
}

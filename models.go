package main

import (
	"github.com/google/uuid"
	"github.com/mata-codes/rssagg/internal/database"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func databaseToUser(db database.User) User {
	return User{
		ID:        db.ID,
		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		Name:      db.Name,
		ApiKey:    db.ApiKey,
	}
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseToFeed(db database.Feed) Feed {
	return Feed{
		ID:        db.ID,
		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		Name:      db.Name,
		Url:       db.Url,
		UserID:    db.UserID,
	}
}

func databaseToFeeds(dbFeeds []database.Feed) []Feed {
	feeds := []Feed{}
	for _, dbfeed := range dbFeeds {
		feeds = append(feeds, databaseToFeed(dbfeed))
	}
	return feeds
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func databaseToFeedFollow(db database.FeedFollow) FeedFollow {
	return FeedFollow{
		ID:        db.ID,
		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		UserID:    db.UserID,
		FeedID:    db.FeedID,
	}
}

func databaseToFeedFollows(dbFeedFollows []database.FeedFollow) []FeedFollow {
	feeds := []FeedFollow{}
	for _, dbfeed := range dbFeedFollows {
		feeds = append(feeds, databaseToFeedFollow(dbfeed))
	}
	return feeds
}

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	Url         string    `json:"url"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func databaseToPost(dbPost database.Post) Post {
	var description *string
	if dbPost.Description.Valid {
		description = &dbPost.Description.String
	}
	return Post{
		ID:          dbPost.ID,
		CreatedAt:   dbPost.CreatedAt,
		UpdatedAt:   dbPost.UpdatedAt,
		Title:       dbPost.Title,
		Description: description,
		PublishedAt: dbPost.PublishedAt,
		Url:         dbPost.Url,
		FeedID:      dbPost.FeedID,
	}

}

func databaseToPosts(dbPosts []database.Post) []Post {

	posts := []Post{}
	for _, dbpost := range dbPosts {
		posts = append(posts, databaseToPost(dbpost))
	}
	return posts
}

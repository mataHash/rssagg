package main

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/mata-codes/rssagg/internal/database"
	"log"
	"strings"
	"sync"
	"time"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBeetwenRequest time.Duration,
) {
	log.Printf("Scrapping on %v gorutines every %s duration\n", concurrency, timeBeetwenRequest)
	ticker := time.NewTicker(timeBeetwenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching feeds: ", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()

	}

}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedASFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Couldn't mark feed as fetched:", err)
		return
	}
	feedRss, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}
	for _, item := range feedRss.Channel.Item {
		description := sql.NullString{}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing published date %v with err %v ", item.PubDate, err)
			continue
		}

		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "llave duplicada") {
				continue
			}
			log.Println("Failed to create post:", err)

		}
	}
	log.Printf("Feed %s collected, %d posts found", feed.Name, len(feedRss.Channel.Item))
}

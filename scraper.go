package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"log"
	"github.com/KillerBeast69/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch next feed: %v", err)
	}

	params := database.MarkFeedFetchedParams{
		ID: feed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
	}

	err = s.db.MarkFeedFetched(context.Background(), params)
	if err != nil {
		return fmt.Errorf("failed to mark feed as fetched")
	}

	rssfeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("could not fetch RSS feed %s: %v", feed.Name, err)
	}

	fmt.Printf("\n--- found %d posts for %s ---\n", len(rssfeed.Channel.Item), feed.Name)
	for _, item := range rssfeed.Channel.Item {
		pudDate, err := parseData(item.PubDate)
		if err != nil {
			log.Printf("could not parse date for post %s: %v", item.Title, err)
			continue
		}

		fmt.Printf("* %s\n", item.Title)
		post_params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: pudDate,
			FeedID:      feed.ID,
		}

		_, err = s.db.CreatePost(context.Background(), post_params)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("error while creating a post: %v", err)
		}
	}
	return nil
}

func parseData(dataStr string) (time.Time, error) {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC3339,
	}
		
	for _, layout := range layouts {
		if t, err := time.Parse(layout, dataStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("could not parse date: %s", dataStr)
}

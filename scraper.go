package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/KillerBeast69/blog-aggregator/internal/database"
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
		fmt.Printf("* %s\n", item.Title)
	}

	return nil
}

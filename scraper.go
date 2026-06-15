package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
		fmt.Printf("* %s\n", item.Title)
		post_params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: item.PubDate,
			FeedID:      feed.ID,
		}

		_, err := s.db.CreatePost(context.Background(), post_params)
		if err != nil {
			// where should I log my error?
			return fmt.Errorf("error while creating a post: %v", err)
		}

		// how do I convert the data into database sql type

	}

	return nil
}

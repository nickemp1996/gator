package main

import (
	"os"
	"fmt"
	"time"
	"strings"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/nickemp1996/gator/internal/rss"
	"github.com/nickemp1996/gator/internal/database"
)

func scrapeFeeds(s *state) {
	nextFeed, err := s.queries.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println("Error getting feed url:", err)
		os.Exit(1)
	}

	err = s.queries.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		fmt.Println("Error marking feed as fetched:", err)
		os.Exit(1)
	}

	feed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		fmt.Println("Error getting feed:", err)
		os.Exit(1)
	}

	fmt.Printf("Processing feed from %s\n", feed.Channel.Title)

	for _, item := range feed.Channel.Item {
		var nullString sql.NullString
		nullString.Valid = false
		if len(item.Description) != 0 {
			nullString = sql.NullString{
				String: item.Description,
				Valid:  true,
			}
		}

		var nullTime sql.NullTime
		nullTime.Valid = false
		if len(item.PubDate) != 0 {
			t, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				fmt.Println("Error formating PubDate to time:", err)
			} else {
				nullTime = sql.NullTime{
					Time:	t,
					Valid:	true,
				}
			}
		}

		postParams := database.CreatePostParams{
			ID:          	uuid.New(),
			Title:       	item.Title,
			Url:         	item.Link,
			Description:	nullString,
			PublishedAt:	nullTime,
			FeedID:      	nextFeed.ID,
		}

		post, err := s.queries.CreatePost(context.Background(), postParams)
		if err != nil {
			if strings.Contains(err.Error(), "violates unique constraint") {
				continue
			}
			fmt.Println("Error creating post:", err)
			os.Exit(1)
		}

		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)
		fmt.Printf("Description: %s\n", post.Description.String)
		fmt.Printf("Published at: %v\n\n", post.PublishedAt.Time)
	}
}
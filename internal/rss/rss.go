package rss

import (
	"io"
	"fmt"
	"html"
	"context"
	"net/http"
	"encoding/xml"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating rquest: %v\n", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("Error performing GET request due to context error: %v\n", ctx.Err())
		}
		return nil, fmt.Errorf("Error performing GET request: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { // Check if the status code is not 200 OK
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("Warning: Could not read response body for status %d: %v\n", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("Received non-200 status code: %d %s. Response body: %s\n", resp.StatusCode, resp.Status, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading GET response body: %v\n", err)
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling XML: %v\n", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}
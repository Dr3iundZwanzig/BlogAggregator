package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Dr3iundZwanzig/BlogAggregator/internal/database"
	"github.com/google/uuid"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return &RSSFeed{}, err
	}

	return helperUnescapeString(&feed), nil
}

func helperUnescapeString(feed *RSSFeed) *RSSFeed {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}
	return feed
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("error getting next feed to fetch")
	}

	currentTime := time.Now().UTC()
	mark := database.MarkFeedFetchedParams{
		ID: nextFeed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  currentTime,
			Valid: true,
		},
		UpdatedAt: currentTime,
	}
	err = s.db.MarkFeedFetched(context.Background(), mark)
	if err != nil {
		return fmt.Errorf("error marking feed: %v", err)
	}

	fetchedFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("Error fetching feed")
	}

	for _, item := range fetchedFeed.Channel.Item {
		description := sql.NullString{
			Valid: false,
		}
		fmt.Println(item.PubDate)
		published, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			return fmt.Errorf("error parsing puplish date: %v", err)
		}
		if item.Description != "" {
			description.Valid = true
			description.String = item.Description
		}
		postParam := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   currentTime,
			UpdatedAt:   currentTime,
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: published,
			FeedID:      nextFeed.ID,
		}
		newPost, err := s.db.CreatePost(context.Background(), postParam)
		if err != nil {
			if !strings.Contains(err.Error(), "posts_url_key") {
				return fmt.Errorf("error ceating post: %v", err)
			}

		}
		fmt.Printf("post createt: %v", newPost.Title)
	}
	return nil
}

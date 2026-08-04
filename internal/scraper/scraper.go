package scraper

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"html"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/MhdFiras-3/gofeed/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gofeed")
	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	rssFeedData := &RSSFeed{}
	if err := xml.Unmarshal(resData, rssFeedData); err != nil {
		return nil, err
	}
	rssFeedData.Channel.Title = html.UnescapeString(rssFeedData.Channel.Title)
	rssFeedData.Channel.Description = html.UnescapeString(rssFeedData.Channel.Description)
	for i := range rssFeedData.Channel.Item {
		rssFeedData.Channel.Item[i].Title = html.UnescapeString(rssFeedData.Channel.Item[i].Title)
		rssFeedData.Channel.Item[i].Description = html.UnescapeString(rssFeedData.Channel.Item[i].Description)
	}
	return rssFeedData, nil
}

func scrapeFeed(ctx context.Context, DB *database.Queries, feedURL string, feedID uuid.UUID) error {
	rssDataFeed, err := fetchFeed(ctx, feedURL)
	if err != nil {
		return err
	}

	if err = DB.MarkFeedFetched(ctx, feedID); err != nil {
		return err
	}
	for _, item := range rssDataFeed.Channel.Item {

		itemNullDescrip := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}
		itemNullPub := sql.NullTime{}
		if item.PubDate != "" {
			time, err := time.Parse("2006-01-02 15:04:05", item.PubDate)
			if err != nil {
				log.Printf("failed to parse time for post url: %s", item.Link)
			} else {
				itemNullPub = sql.NullTime{
					Time:  time,
					Valid: true,
				}
			}

		}

		_, err = DB.CreatePost(ctx, database.CreatePostParams{
			Title:       item.Title,
			Url:         item.Link,
			Description: itemNullDescrip,
			FeedID:      feedID,
			PublishedAt: itemNullPub,
		})
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				log.Printf("unique violation on constraint: post already exist id: %s", feedID)
				continue
			}
			log.Printf("failed to create post: %v, url: %s", err, item.Link)
		}

	}
	return nil
}

func StartScraping(ctx context.Context, DB *database.Queries, tick time.Duration) {
	var wg sync.WaitGroup
	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		feeds, err := DB.GetNextFeedsToFetch(ctx, 2)
		if err != nil {
			log.Printf("failed to fetch feeds to scrape: %v", err)
			continue
		}
		for _, feed := range feeds {
			wg.Add(1)
			go func(feedURL string, feedID uuid.UUID) {
				defer wg.Done()
				scrapeFeed(ctx, DB, feedURL, feedID)

			}(feed.Url, feed.ID)

		}
		wg.Wait()
	}

}

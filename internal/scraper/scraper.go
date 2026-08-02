package scraper

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"
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

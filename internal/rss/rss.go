package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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
	client := &http.Client{}

	// Create a new HTTP request
	req, err := http.NewRequest("GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error while creating a request from URL: %w", err)
	}
	req.Header.Add("User-Agent", "gator")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error while receiving answer from URL: %w", err)
	}
	defer resp.Body.Close()

	// Read the received data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error while reading data from URL: %w", err)
	}

	// Unmarshal the (XML) data to RSS structs
	rss := &RSSFeed{}
	if err = xml.Unmarshal(data, &rss); err != nil {
		return &RSSFeed{}, fmt.Errorf("error while unmarshaling XML data from URL: %w", err)
	}

	// Clean escaped HTML entities
	for _, item := range rss.Channel.Item {
		item.unescapeHTML()
	}
	rss.unescapeHTML()

	return rss, nil
}

func (feed *RSSFeed) unescapeHTML() {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
}

func (item *RSSItem) unescapeHTML() {
	item.Title = html.UnescapeString(item.Title)
	item.Description = html.UnescapeString(item.Description)
}

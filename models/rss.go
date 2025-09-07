package models

import "time"

// Normalized RSS item you can later insert into `feeds` or pass to sentiment
type RSSItem struct {
	SourceID       int64     `json:"source_id"`
	Lang           string    `json:"lang"` // e.g., "hun"
	Title          string    `json:"title"`
	Link           string    `json:"link"`
	SentimentKey   string    `json:"sentiment_key"`
	SentimentValue float64   `json:"sentiment_value"`
	Category       string    `json:"category"`
	Published      time.Time `json:"published"` // UTC
	FeedDate       string    `json:"feed_date"` // YYYY-MM-DD
}

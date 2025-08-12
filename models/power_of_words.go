package models

import "time"

// FeedWord represents a single word entry from the feeds table.
type FeedWord struct {
	Word string `json:"word"`
}

// Feed represents a row in the feeds table.
type Feed struct {
	ID        int       `json:"id"`
	Published time.Time `json:"published"`
	Words     []string  `json:"words"`
	Title     string    `json:"title"`
	SourceID  int       `json:"source_id"`
}

// FeedSentiment represents a row in the feed_sentiments table.
type FeedSentiment struct {
	ID             int     `json:"id"`
	SentimentKey   string  `json:"sentiment_key"`
	SentimentValue float32 `json:"sentiment_value"`
	Sentiments     string  `json:"sentiments"`
}

// Source represents a row in the sources table.
type Source struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// SentimentGroupRow
type SentimentGroupRow struct {
	GroupBy      string `json:"group_by"`
	SentimentKey string `json:"sentiment_key"`
	Count        int    `json:"count"`
}

// FeedEnvelope bundles a feed with its sentiment and source.
type FeedEnvelope struct {
	Feed      Feed          `json:"feeds"`
	Sentiment FeedSentiment `json:"feed_sentiments"`
	Source    Source        `json:"sources"`
}

type Sentiments struct {
	Positive int `json:"positive_sentiments"`
	Negative int `json:"negative_sentiments"`
	Neutral  int `json:"neutral_sentiments"`
}

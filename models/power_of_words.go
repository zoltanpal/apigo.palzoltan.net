package models

import "time"

// FeedWord represents a single word entry from the feeds table.
type FeedWord struct {
	Word string `json:"word"`
}

type WordCount struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

// Feed represents a row in the feeds table.
type Feed struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Link      string    `json:"link"`
	SourceID  int       `json:"source_id"`
	Words     []string  `json:"words"`
	Published time.Time `json:"published"`
}

// FeedSentiment represents a row in the feed_sentiments table.
type FeedSentiment struct {
	ID                int     `json:"id"`
	SentimentKey      string  `json:"sentiment_key"`
	SentimentValue    float32 `json:"sentiment_value"`
	Sentiments        string  `json:"sentiments"`
	SentimentCompound float32 `json:"sentiment_compound"`
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

type TopFeedRow struct {
	Title             string    `json:"title"`
	Published         time.Time `json:"published"`
	SourceName        string    `json:"source_name"`
	SentimentValue    float64   `json:"sentiment_value"`
	SentimentCompound float64   `json:"sentiment_compound"`
}

type BiasDetectionRow struct {
	SourceName        string  `json:"source_name"`
	Keyword           string  `json:"keyword"`
	MentionCount      int     `json:"mention_count"`
	NetSentimentScore float64 `json:"net_sentiment_score"`
	SentimentStdDev   float64 `json:"sentiment_std_dev"`
}

type CorrelationRow struct {
	SourceName  string  `json:"sourcename"`
	Month       string  `json:"month"`
	AvgCompound float64 `json:"avg_compound"`
}

type WordCoOccurrenceRow struct {
	CoWord        string `json:"co_word"`
	CoOccurrence  int    `json:"co_occurrence"`
	PositiveCount int    `json:"positive_count"`
	NegativeCount int    `json:"negative_count"`
	NeutralCount  int    `json:"neutral_count"`
}

type FeedWithDetails struct {
	Feed          Feed          `json:"feed"`
	FeedSentiment FeedSentiment `json:"feed_sentiment"`
	Source        Source        `json:"source"`
}

type FeedResponse struct {
	Total int               `json:"total"`
	Page  int               `json:"page"`
	Feeds []FeedWithDetails `json:"feeds"`
}

type PhraseFrequencyRow struct {
	Source    string `json:"source"`
	Year      int    `json:"year"`
	Phrase    string `json:"phrase"`
	Month     string `json:"date_ts"`
	Frequency int    `json:"freq"`
	Ranked    int    `json:"rnk"`
}

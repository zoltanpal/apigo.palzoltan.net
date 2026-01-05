package models

import "time"

// FeedWord represents a single word entry from the feeds table.
type FeedWord struct {
	Word string `json:"word"`
}

// WordCount represents a word and its occurrence count.
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

type SourceDetails struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Rss  string `json:"rss"`
	Lang string `json:"lang"`
}

// SentimentGroupRow represents aggregated sentiment data grouped by a specific field.
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

// Sentiments holds counts of different sentiment types.
type Sentiments struct {
	Positive int `json:"positive_sentiments"`
	Negative int `json:"negative_sentiments"`
	Neutral  int `json:"neutral_sentiments"`
}

// TopFeedRow represents a top feed entry with sentiment details.
type TopFeedRow struct {
	Title             string    `json:"title"`
	Published         time.Time `json:"published"`
	SourceName        string    `json:"source_name"`
	SentimentValue    float64   `json:"sentiment_value"`
	SentimentCompound float64   `json:"sentiment_compound"`
}

// BiasDetectionRow represents bias detection metrics for a keyword from a source.
type BiasDetectionRow struct {
	SourceName        string  `json:"source_name"`
	Keyword           string  `json:"keyword"`
	MentionCount      int     `json:"mention_count"`
	NetSentimentScore float64 `json:"net_sentiment_score"`
	SentimentStdDev   float64 `json:"sentiment_std_dev"`
}

// CorrelationRow represents correlation data between source and sentiment over time.
type CorrelationRow struct {
	SourceName  string  `json:"sourcename"`
	Month       string  `json:"month"`
	AvgCompound float64 `json:"avg_compound"`
}

// WordCoOccurrenceRow represents co-occurrence data of words with sentiment counts.
type WordCoOccurrenceRow struct {
	CoWord        string `json:"co_word"`
	CoOccurrence  int    `json:"co_occurrence"`
	PositiveCount int    `json:"positive_count"`
	NegativeCount int    `json:"negative_count"`
	NeutralCount  int    `json:"neutral_count"`
}

// FeedWithDetails bundles a feed with its sentiment and source details.
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

type GNewsItem struct {
	Title     string `json:"title"`
	Published string `json:"published"`
	Source    string `json:"source"`
}

type GNewsResponse struct {
	Title          string  `json:"title"`
	Source         string  `json:"source"`
	Published      string  `json:"published"`
	SentimentKey   string  `json:"sentiment_key"`
	SentimentValue float32 `json:"sentiment_value"`
}

type PhraseFrequencyRow struct {
	Source    string `json:"source"`
	Phrase    string `json:"phrase"`
	Year      int    `json:"year"`
	DateGroup int    `json:"date_group"`
	Frequency int    `json:"freq"`
	Ranked    int    `json:"rnk"`
}

type Statistics struct {
	FirstFeedDate        string  `json:"first_feed_date"`
	LastFeedDate         string  `json:"last_feed_date"`
	TimeSpanDays         int     `json:"time_span_days"`
	TotalFeeds           int     `json:"total_feeds"`
	TotalSources         int     `json:"total_sources"`
	TotalPositive        int     `json:"total_positive"`
	TotalNegative        int     `json:"total_negative"`
	TotalNeutral         int     `json:"total_neutral"`
	PctPositive          float64 `json:"pct_positive"`
	PctNegative          float64 `json:"pct_negative"`
	PctNeutral           float64 `json:"pct_neutral"`
	AvgFeedsPerDay       float64 `json:"avg_feeds_per_day"`
	MostActiveSourceName string  `json:"most_active_source_name"`
}

// RSSItem represents a single RSS feed item with sentiment and category information.
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

package models

// SentimentSeries
type SentimentSeries struct {
	Negative []int    `json:"Negative"`
	Neutral  []int    `json:"Neutral"`
	Positive []int    `json:"Positive"`
	Keys     []string `json:"keys"`
}

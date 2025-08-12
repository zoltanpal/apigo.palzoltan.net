package utils

import (
	"sort"
	"strings"

	"golang-restapi/models"
)

// SentimentSeries matches your Python output shape (plus Keys for ordering).
type SentimentSeries struct {
	Negative []int    `json:"Negative"`
	Neutral  []int    `json:"Neutral"`
	Positive []int    `json:"Positive"`
	Keys     []string `json:"keys"` // sorted list of group keys (e.g., source IDs or dates)
}

// GenerateSentimentSeries organizes rows into 3 parallel series,
// aligned by a sorted list of group keys.
func GenerateSentimentSeries(rows []models.SentimentGroupRow) SentimentSeries {
	neg := map[string]int{}
	neu := map[string]int{}
	pos := map[string]int{}
	keySet := map[string]struct{}{}

	for _, r := range rows {
		key := r.GroupBy
		keySet[key] = struct{}{}

		switch strings.ToLower(r.SentimentKey) {
		case "negative":
			neg[key] += r.Count
		case "neutral":
			neu[key] += r.Count
		case "positive":
			pos[key] += r.Count
		default:
			// no
		}
	}

	// Collect and sort keys for stable ordering
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	series := SentimentSeries{
		Negative: make([]int, len(keys)),
		Neutral:  make([]int, len(keys)),
		Positive: make([]int, len(keys)),
		Keys:     keys,
	}

	for i, k := range keys {
		series.Negative[i] = neg[k]
		series.Neutral[i] = neu[k]
		series.Positive[i] = pos[k]
	}

	return series
}

package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"golang-restapi/db"
	"golang-restapi/models"
	"golang-restapi/queries"
)

// GetFeedsWords returns all feed words between startDate and endDate.
func GetFeedsWords(startDate, endDate string) ([]models.FeedWord, error) {
	rows, err := db.DB.Query(queries.GetWordsByDateRange, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("GetFeedsWords: query error: %w", err)
	}
	defer rows.Close()

	var words []models.FeedWord
	for rows.Next() {
		var fw models.FeedWord
		if err := rows.Scan(&fw.Word); err != nil {
			return nil, fmt.Errorf("GetFeedsWords: scan error: %w", err)
		}
		words = append(words, fw)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetFeedsWords: rows iteration error: %w", err)
	}
	return words, nil
}

// GetSentimentGrouped returns grouped sentiment counts
// groupBy: "source" or "date"
func GetSentimentGrouped(
	ctx context.Context,
	startDate, endDate string,
	freeText string,
	groupBy string,
) ([]models.SentimentGroupRow, error) {

	base := queries.GetSentimentGroupedBySource
	if strings.EqualFold(groupBy, "date") {
		base = queries.GetSentimentGroupedByDate
	}

	conds := []string{"f.feed_date BETWEEN $1 AND $2"}
	args := []interface{}{startDate, endDate}
	idx := 3

	if freeText != "" {
		conds = append(conds, fmt.Sprintf("(f.title ILIKE '%%' || $%d || '%%')", idx))
		args = append(args, freeText)
		idx++
	}

	where := "WHERE " + strings.Join(conds, " AND ")

	sql := base + `
        ` + where + `
        GROUP BY group_by, sentiment_key
        ORDER BY group_by ASC
    `

	rows, err := db.DB.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("GetSentimentGrouped: query error: %w", err)
	}
	defer rows.Close()

	out := make([]models.SentimentGroupRow, 0, 128)
	for rows.Next() {
		var r models.SentimentGroupRow
		if err := rows.Scan(&r.GroupBy, &r.SentimentKey, &r.Count); err != nil {
			return nil, fmt.Errorf("GetSentimentGrouped: scan error: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetSentimentGrouped: rows iteration error: %w", err)
	}
	return out, nil
}

// CountSentiments returns the aggregated counts for positive/negative/neutral.
func CountSentiments(ctx context.Context, startDate, endDate string) (models.Sentiments, error) {
	var out models.Sentiments

	// Append time boundaries to match your Python version
	start := startDate + " 00:00:00"
	end := endDate + " 23:59:59"

	row := db.DB.QueryRowContext(ctx, queries.CountSentiments, start, end)
	if err := row.Scan(&out.Positive, &out.Negative, &out.Neutral); err != nil {
		return models.Sentiments{}, fmt.Errorf("CountSentiments: scan error: %w", err)
	}
	return out, nil
}

// Fetch the Feeds with the max positive & negative sentiments
func TopFeeds(ctx context.Context, startDate, endDate, posNeg string, limit int) ([]models.TopFeedRow, error) {
	start := startDate + " 00:00:00"
	end := endDate + " 23:59:59"

	rows, err := db.DB.QueryContext(ctx, queries.TopFeeds, start, end, strings.ToLower(posNeg), limit)
	if err != nil {
		return nil, fmt.Errorf("TopFeeds: query error: %w", err)
	}
	defer rows.Close()

	out := make([]models.TopFeedRow, 0, limit)
	for rows.Next() {
		var r models.TopFeedRow
		if err := rows.Scan(
			&r.Title,
			&r.Published,
			&r.SourceName,
			&r.SentimentValue,
			&r.SentimentCompound,
		); err != nil {
			return nil, fmt.Errorf("TopFeeds: scan error: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("TopFeeds: rows iteration error: %w", err)
	}
	return out, nil
}

// BiasDetection
func BiasDetection(
	ctx context.Context,
	startDate, endDate string,
	words []string,
) ([]models.BiasDetectionRow, error) {

	if len(words) == 0 {
		return nil, fmt.Errorf("BiasDetection: words must not be empty")
	}

	// $1 words[], $2 start, $3 end
	args := []any{
		pq.Array(words),
		startDate + " 00:00:00",
		endDate + " 23:59:59",
	}

	sql := fmt.Sprintf(queries.BiasDetection)

	rows, err := db.DB.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("BiasDetection: query error: %w", err)
	}
	defer rows.Close()

	out := make([]models.BiasDetectionRow, 0, 128)
	for rows.Next() {
		var r models.BiasDetectionRow
		if err := rows.Scan(
			&r.SourceName,
			&r.Keyword,
			&r.MentionCount,
			&r.NetSentimentScore,
			&r.SentimentStdDev,
		); err != nil {
			return nil, fmt.Errorf("BiasDetection: scan error: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("BiasDetection: rows iteration error: %w", err)
	}
	return out, nil
}

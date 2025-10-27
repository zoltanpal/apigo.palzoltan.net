package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"golang-restapi/db"
	"golang-restapi/models"
	"golang-restapi/queries"
	"golang-restapi/utils"

	"github.com/lib/pq"

	"golang-restapi/sentimentpb"
)

// MostCommonWords fetches feed words, filters stopwords, and counts occurrences.
func MostCommonWords(ctx context.Context, startDate, endDate string, n int) ([]models.WordCount, error) {

	rows, err := db.DB.QueryContext(ctx, queries.GetWordsByDateRange, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("MostCommonWords: query error: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)

	for rows.Next() {
		var words []sql.NullString
		if err := rows.Scan(pq.Array(&words)); err != nil {
			return nil, fmt.Errorf("MostCommonWords: scan error: %w", err)
		}
		for _, w := range words {
			if !w.Valid {
				continue
			}
			word := w.String
			if utils.IsStopword(word) {
				continue
			}
			counts[word]++
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("MostCommonWords: rows iteration error: %w", err)
	}

	// convert map -> slice
	result := make([]models.WordCount, 0, len(counts))
	for word, count := range counts {
		if count < 2 {
			continue // skip rare words
		}
		result = append(result, models.WordCount{Word: word, Count: count})
	}

	// sort by frequency
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	// take top N
	if n < len(result) {
		result = result[:n]
	}
	return result, nil
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
	word string,
) ([]models.BiasDetectionRow, error) {

	w := utils.SanitizeTSWord(word)
	if w == "" || strings.Contains(w, " ") {
		return nil, fmt.Errorf("BiasDetection: 'word' must be a single non-empty string")
	}

	args := []any{
		w,                       // $1
		startDate + " 00:00:00", // $2
		endDate + " 23:59:59",   // $3
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

// CorrelationBetweenSourcesAvgCompound computes the monthly average sentiment_compound per source.
func CorrelationBetweenSourcesAvgCompound(
	ctx context.Context,
	startDate, endDate string,
	word string,
	sources []int, // optional
) ([]models.CorrelationRow, error) {

	w := utils.SanitizeTSWord(word)
	if w == "" || strings.Contains(w, " ") {
		return nil, fmt.Errorf("correlation: 'word' must be a single non-empty string")
	}

	args := []any{
		w,                       // $1
		startDate + " 00:00:00", // $2
		endDate + " 23:59:59",   // $3
	}
	extra := ""
	if len(sources) > 0 {
		extra = " AND f.source_id = ANY($4)"
		args = append(args, pq.Array(sources))
	}

	sql := fmt.Sprintf(queries.CorrelationBetweenSourcesAvgCompound, extra)

	rows, err := db.DB.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("correlation(single): query error: %w", err)
	}
	defer rows.Close()

	out := make([]models.CorrelationRow, 0, 64)
	for rows.Next() {
		var r models.CorrelationRow
		var month time.Time
		if err := rows.Scan(&r.SourceName, &month, &r.AvgCompound); err != nil {
			return nil, fmt.Errorf("correlation(single): scan error: %w", err)
		}
		r.Month = month.Format("2006-01-02")
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("correlation(single): rows iteration error: %w", err)
	}
	return out, nil
}

// WordCoOccurrences returns the top co-occurring words with a given `word`.
func WordCoOccurrences(
	ctx context.Context,
	startDate, endDate string,
	word string,
	sources []int, // optional
) ([]models.WordCoOccurrenceRow, error) {
	if word == "" {
		return nil, fmt.Errorf("word is required")
	}

	args := []any{
		word,                          // $1
		startDate + " 00:00:00",       // $2
		endDate + " 23:59:59",         // $3
		pq.Array(utils.StopWordsList), // $4
	}

	extra := ""
	if len(sources) > 0 {
		extra = " AND f.source_id = ANY($5)"
		args = append(args, pq.Array(sources))
	}

	sql := fmt.Sprintf(queries.WordCoOccurrences, extra)

	rows, err := db.DB.QueryContext(ctx, sql, args...)

	if err != nil {
		return nil, fmt.Errorf("WordCoOccurrences: query error: %w", err)
	}
	defer rows.Close()

	out := make([]models.WordCoOccurrenceRow, 0, 64)
	for rows.Next() {
		var r models.WordCoOccurrenceRow
		if err := rows.Scan(
			&r.CoWord,
			&r.CoOccurrence,
			&r.PositiveCount,
			&r.NegativeCount,
			&r.NeutralCount,
		); err != nil {
			return nil, fmt.Errorf("WordCoOccurrences: scan error: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("WordCoOccurrences: rows iteration error: %w", err)
	}
	return out, nil
}

func PhraseFrequencyTrends(
	ctx context.Context,
	startDate, endDate, dateGroup string,
	sources []int, // optional
	namesExcluded bool, // use bool instead of string
) ([]models.PhraseFrequencyRow, error) {

	// Ensure you always pass arrays (empty is fine).
	// Param order maps to $1..$7 in SQL above.
	args := []any{
		startDate + " 00:00:00",          // $1
		endDate + " 23:59:59",            // $2
		dateGroup,                        // $3 ("week" or "month")
		pq.Array(utils.StopwordsSimple),  // $4 ::text[]
		namesExcluded,                    // $5 ::boolean
		pq.Array(utils.StopPhrasessList), // $6 ::text[] (used only if $5 = true)
		pq.Array(sources),                // $7 ::int[]   (empty => no filter)
	}

	rows, err := db.DB.QueryContext(ctx, queries.PhraseFrequencyTrendsNew, args...)
	if err != nil {
		return nil, fmt.Errorf("PhraseFrequencyTrends: query error: %w", err)
	}
	defer rows.Close()

	out := make([]models.PhraseFrequencyRow, 0, 64)
	for rows.Next() {
		var r models.PhraseFrequencyRow
		if err := rows.Scan(
			&r.Source,
			&r.Phrase,
			&r.Year,
			&r.DateGroup,
			&r.Frequency,
			&r.Ranked,
		); err != nil {
			return nil, fmt.Errorf("PhraseFrequencyTrends: scan error: %w", err)
		}
		out = append(out, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("PhraseFrequencyTrends: rows iteration error: %w", err)
	}

	return out, nil

}

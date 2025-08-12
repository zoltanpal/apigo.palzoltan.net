package queries

const (
	// GetWordsByDateRange returns only the 'words' column
	// for all feeds whose feed_date is between two given dates.
	GetWordsByDateRange = `
        SELECT words
        FROM feeds_partitioned
        WHERE feed_date BETWEEN $1 AND $2
    `

	// GetSentimentGroupedBySource
	GetSentimentGroupedBySource = `
        SELECT
            f.source_id::text AS group_by,
            COALESCE(fs.sentiment_key, 'none') AS sentiment_key,
            COUNT(f.id) AS count
        FROM feeds_partitioned f
        LEFT JOIN feed_sentiments fs
          ON fs.feed_id = f.id AND fs.model_id = 1
    `
	// GetSentimentGroupedByDate
	GetSentimentGroupedByDate = `
        SELECT
            to_char(f.feed_date, 'YYYY-MM-DD') AS group_by,
            COALESCE(fs.sentiment_key, 'none') AS sentiment_key,
            COUNT(f.id) AS count
        FROM feeds_partitioned f
        LEFT JOIN feed_sentiments fs
          ON fs.feed_id = f.id AND fs.model_id = 1
    `
	// CountSentiments
	CountSentiments = `
        SELECT
            COALESCE(SUM(CASE WHEN fs.sentiment_key = 'positive' THEN 1 ELSE 0 END), 0) AS positive_sentiments,
            COALESCE(SUM(CASE WHEN fs.sentiment_key = 'negative' THEN 1 ELSE 0 END), 0) AS negative_sentiments,
            COALESCE(SUM(CASE WHEN fs.sentiment_key = 'neutral'  THEN 1 ELSE 0 END), 0) AS neutral_sentiments
        FROM feed_sentiments fs
        WHERE fs.model_id = 1
        AND fs.feed_date BETWEEN $1 AND $2
    `
)

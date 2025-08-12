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
	// TopFeeds
	TopFeeds = `
        SELECT
            f.title,
            f.published,
            s.name,
            fs.sentiment_value,
            fs.sentiment_compound
        FROM feed_sentiments fs
        JOIN feeds   f ON fs.feed_id  = f.id
        JOIN sources s ON f.source_id = s.id
        WHERE fs.model_id = 1
            AND f.feed_date BETWEEN $1 AND $2
            AND fs.sentiment_key = $3
        ORDER BY fs.sentiment_value DESC
        LIMIT $4
    `

	// BiasDetection
	BiasDetection = `
        WITH input_words AS (
            SELECT unnest($1::text[]) AS input_word
        )
        SELECT
            s.name AS source_name,
            iw.input_word AS keyword,
            COUNT(*) AS mention_count,
            (
                (SUM(CASE WHEN fs.sentiment_key = 'positive' THEN fs.sentiment_value ELSE 0 END)
            - SUM(CASE WHEN fs.sentiment_key = 'negative' THEN fs.sentiment_value ELSE 0 END))
                / NULLIF(COUNT(*), 0)::double precision
            )::double precision AS net_sentiment_score,
            COALESCE(STDDEV(fs.sentiment_value)::double precision, 0)::double precision AS sentiment_std_dev
        FROM feeds f
        JOIN feed_sentiments fs ON fs.feed_id = f.id AND fs.model_id = 1
        JOIN sources s          ON f.source_id = s.id
        CROSS JOIN input_words iw
        WHERE f.published BETWEEN $2 AND $3
            AND f.search_vector @@ to_tsquery('hungarian', iw.input_word || ':*')
        GROUP BY s.name, iw.input_word
        ORDER BY iw.input_word, net_sentiment_score DESC;
    `
)

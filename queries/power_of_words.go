package queries

const (
	FeedsBaseQuery = `
        SELECT
            f.id, f.title, f.link, f.source_id, f.words, f.published,
            fs.id, fs.sentiment_key, fs.sentiment_value, fs.sentiment_compound,
            s.id, s.name
        FROM feeds f
        JOIN feed_sentiments fs ON fs.feed_id = f.id AND fs.model_id = 1
        JOIN sources s ON f.source_id = s.id
        WHERE f.feed_date BETWEEN $1 AND $2
        %s
        ORDER BY f.published DESC
        LIMIT $%d OFFSET $%d
    `

	// FeedsCountQuery: same filters, just COUNT(*).
	FeedsCountQuery = `
            SELECT COUNT(*)
            FROM feeds f
            JOIN feed_sentiments fs ON fs.feed_id = f.id AND fs.model_id = 1
            JOIN sources s ON f.source_id = s.id
            WHERE f.feed_date BETWEEN $1 AND $2
            %s
        `

	// GetWordsByDateRange returns only the 'words' column
	// for all feeds whose feed_date is between two given dates.
	GetWordsByDateRange = `
        SELECT words
        FROM feeds
        WHERE feed_date BETWEEN $1 AND $2
    `

	// GetSentimentGroupedBySource
	GetSentimentGroupedBySource = `
        SELECT
            f.source_id::text AS group_by,
            COALESCE(fs.sentiment_key, 'none') AS sentiment_key,
            COUNT(f.id) AS count
        FROM feeds f
        LEFT JOIN feed_sentiments fs
          ON fs.feed_id = f.id AND fs.model_id = 1
    `
	// GetSentimentGroupedByDate
	GetSentimentGroupedByDate = `
        SELECT
            to_char(f.feed_date, 'YYYY-MM-DD') AS group_by,
            COALESCE(fs.sentiment_key, 'none') AS sentiment_key,
            COUNT(f.id) AS count
        FROM feeds f
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
            s.name as source_name,
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
        SELECT
        s.name          AS source_name,
        $1::text        AS keyword,
        COUNT(DISTINCT f.id) AS mention_count,
        (
            (SUM(CASE WHEN fs.sentiment_key = 'positive' THEN fs.sentiment_value ELSE 0 END)
        - SUM(CASE WHEN fs.sentiment_key = 'negative' THEN fs.sentiment_value ELSE 0 END))
            / NULLIF(COUNT(DISTINCT f.id), 0)::double precision
        )::double precision AS net_sentiment_score,
        COALESCE(STDDEV(fs.sentiment_value)::double precision, 0)::double precision AS sentiment_std_dev
        FROM feeds f
        JOIN feed_sentiments fs ON fs.feed_id = f.id AND fs.model_id = 1
        JOIN sources s          ON s.id = f.source_id
        WHERE f.feed_date BETWEEN $2 AND $3                          -- use feed_date for partition pruning
        AND f.search_vector @@ to_tsquery('public.hun_unaccent', $1 || ':*')
        GROUP BY s.name
        ORDER BY net_sentiment_score DESC;
    `

	// CorrelationBetweenSourcesAvgCompound
	CorrelationBetweenSourcesAvgCompound = `
        SELECT
            s.name AS sourcename,
            date_trunc('month', f.published)::date AS month,
            COALESCE(AVG(fs.sentiment_compound), 0) AS avg_compound
        FROM feeds f
        JOIN feed_sentiments fs ON fs.feed_id = f.id AND fs.model_id = 1 AND fs.feed_date BETWEEN $2 AND $3
        LEFT JOIN sources s ON f.source_id = s.id
        WHERE f.search_vector @@ to_tsquery('public.hun_unaccent', $1)
            AND f.feed_date BETWEEN $2 AND $3 AND fs.feed_date BETWEEN $2 AND $3
            %s
        GROUP BY s.name, month
        ORDER BY s.name, month
    `

	WordCoOccurrences = `
        WITH target_articles AS (
            SELECT f.id, f.words
            FROM feeds f
            WHERE $1 = ANY(f.words)
                AND f.feed_date BETWEEN $2 AND $3
                %s
        ),
        co_words AS (
            SELECT ta.id AS feed_id, w AS co_word
            FROM target_articles ta,
                unnest(ta.words) AS w
            WHERE w <> $1 AND NOT (w = ANY($4))
        ),
        sentiments AS (
            SELECT feed_id,
                    COUNT(*) FILTER (WHERE sentiment_key = 'positive') AS pos_count,
                    COUNT(*) FILTER (WHERE sentiment_key = 'negative') AS neg_count,
                    COUNT(*) FILTER (WHERE sentiment_key = 'neutral')  AS neu_count
            FROM feed_sentiments
            WHERE model_id = 1 AND feed_date BETWEEN $2 AND $3
            GROUP BY feed_id
        )
        SELECT
            cw.co_word,
            COUNT(*) AS co_occurrence,
            COALESCE(SUM(s.pos_count), 0) AS positive_count,
            COALESCE(SUM(s.neg_count), 0) AS negative_count,
            COALESCE(SUM(s.neu_count), 0) AS neutral_count
            FROM co_words cw
            LEFT JOIN sentiments s ON cw.feed_id = s.feed_id
        GROUP BY cw.co_word
        HAVING COUNT(*) > 1
        ORDER BY co_occurrence DESC
        LIMIT 30;
    `

	phraseExclusion = `
        AND NOT ()
    `
	PhraseFrequencyTrends = `
        WITH bigrams AS (
            SELECT
                s.name AS source,
                lower(f.words_masked[i]) || ' ' || lower(f.words_masked[i+1]) AS phrase,
                CASE
                    WHEN $3 = 'week'  THEN extract(isoyear FROM f.feed_date)::int
                    WHEN $3 = 'month' THEN extract(year FROM f.feed_date)::int
                END AS year,
                date_part($3, f.feed_date)::int AS date_group
            FROM feeds f
            JOIN sources s ON f.source_id = s.id
            CROSS JOIN LATERAL generate_subscripts(f.words_masked, 1) AS i
            WHERE f.feed_date BETWEEN $1 AND $2
            AND f.words_masked IS NOT NULL
            AND i < array_length(f.words_masked, 1)
            -- exclude if either word is a stopword
            AND lower(f.words_masked[i]) <> ALL($4)
            AND lower(f.words_masked[i+1]) <> ALL($4) 
            AND (lower(f.words_masked[i]) || ' ' || lower(f.words_masked[i+1])) <> ALL($5::text[])
            %s
        ),
        counts AS (
            SELECT
                source,
                year,
                date_group,
                phrase,
                COUNT(*)::int AS freq
            FROM bigrams
            GROUP BY 1,2,3,4
            HAVING COUNT(*) >= 3
        ),
        ranked AS (
            SELECT
                source,
                phrase,
                year,
                date_group,
                freq,
                RANK() OVER (
                    PARTITION BY source, year, date_group
                    ORDER BY freq DESC, phrase
                ) AS rnk
            FROM counts
        )
        SELECT *
        FROM ranked
        WHERE rnk <= 5
        ORDER BY year, date_group, source, rnk

    `
)

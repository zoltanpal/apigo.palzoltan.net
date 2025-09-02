package repositories

import (
	"context"
	"fmt"
	"strings"

	"golang-restapi/db"
	"golang-restapi/models"
	"golang-restapi/queries"

	"github.com/lib/pq"
)

func GetFeeds(
	ctx context.Context,
	startDate, endDate string,
	sources []int,
	freeText string,
	page, perPage int,
) (models.FeedResponse, error) {

	var resp models.FeedResponse
	var conds []string

	args := []any{startDate, endDate}
	idx := 3 // because $1=start, $2=end

	if len(sources) > 0 {
		conds = append(conds, fmt.Sprintf("AND f.source_id = ANY($%d)", idx))
		args = append(args, pq.Array(sources))
		idx++
	}

	if freeText != "" {
		conds = append(conds, fmt.Sprintf("AND f.title ILIKE '%%' || $%d || '%%'", idx))
		args = append(args, freeText)
		idx++
	}

	whereClause := ""
	if len(conds) > 0 {
		whereClause = " " + strings.Join(conds, " ")
	}

	// --- Count query ---
	countSQL := fmt.Sprintf(queries.FeedsCountQuery, whereClause)
	if err := db.DB.QueryRowContext(ctx, countSQL, args...).Scan(&resp.Total); err != nil {
		return resp, fmt.Errorf("GetFeeds: count query error: %w", err)
	}

	if resp.Total == 0 {
		return resp, nil
	}

	// --- Data query ---
	offset := int64((page - 1) * perPage)
	limit := int64(perPage)

	argsWithLimit := append(args, limit, offset)
	dataSQL := fmt.Sprintf(queries.FeedsBaseQuery, whereClause, idx, idx+1)

	rows, err := db.DB.QueryContext(ctx, dataSQL, argsWithLimit...)

	if err != nil {
		return resp, fmt.Errorf("GetFeeds: data query error: %w", err)
	}
	defer rows.Close()

	feeds := []models.FeedWithDetails{}
	for rows.Next() {
		var f models.Feed
		var fs models.FeedSentiment
		var s models.Source

		if err := rows.Scan(
			&f.ID, &f.Title, &f.Link, &f.SourceID, pq.Array(&f.Words), &f.Published,
			&fs.ID, &fs.SentimentKey, &fs.SentimentValue, &fs.SentimentCompound,
			&s.ID, &s.Name,
		); err != nil {
			return resp, fmt.Errorf("GetFeeds: scan error: %w", err)
		}

		feeds = append(feeds, models.FeedWithDetails{
			Feed:          f,
			FeedSentiment: fs,
			Source:        s,
		})
	}
	if err := rows.Err(); err != nil {
		return resp, fmt.Errorf("GetFeeds: rows iteration error: %w", err)
	}

	resp.Page = page
	resp.Feeds = feeds
	return resp, nil
}

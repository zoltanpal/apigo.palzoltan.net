package utils

import (
	"context"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"

	"golang-restapi/models"
)

type Reader struct {
	parser *gofeed.Parser
}

func NewReader() *Reader {
	return &Reader{parser: gofeed.NewParser()}
}

// FetchSource reads one RSS URL and returns normalized items.
func (r *Reader) FetchSource(ctx context.Context, sourceID int64, url string, lang string) ([]models.RSSItem, error) {
	feed, err := r.parser.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, err
	}

	items := make([]models.RSSItem, 0, len(feed.Items))
	for _, item := range feed.Items {
		if item == nil || item.Title == "" || item.Link == "" {
			continue
		}
		// published time
		ts := time.Now().UTC()
		if item.PublishedParsed != nil {
			ts = item.PublishedParsed.UTC()
		} else if item.UpdatedParsed != nil {
			ts = item.UpdatedParsed.UTC()
		}

		// first category if present
		cat := ""
		if len(item.Categories) > 0 {
			cat = strings.TrimSpace(item.Categories[0])
		}

		items = append(items, models.RSSItem{
			SourceID:  sourceID,
			Lang:      lang,
			Title:     strings.TrimSpace(item.Title),
			Link:      strings.TrimSpace(item.Link),
			Category:  cat,
			Published: ts,
			FeedDate:  ts.Format("2006-01-02"),
		})
	}
	return items, nil
}

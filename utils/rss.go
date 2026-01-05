// utils/rss.go
package utils

import (
	"context"
	"encoding/xml"
	"fmt"
	"golang-restapi/models"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"golang-restapi/sentimentpb"
)

// RSS struct for parsing Google News RSS
type RSS struct {
	Channel struct {
		Items []struct {
			Title   string `xml:"title"`
			Link    string `xml:"link"`
			PubDate string `xml:"pubDate"`
			Source  string `xml:"source"`
		} `xml:"item"`
	} `xml:"channel"`
}

func GetGoogleNews(q, period, lang, country string) ([]models.GNewsItem, error) {

	escapedQ := url.QueryEscape(q)
	url := fmt.Sprintf("https://news.google.com/rss/search?q=%s&hl=%s&gl=%s&ceid=%s:%s",
		escapedQ, lang, country, strings.ToUpper(country), lang)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GetGoogleNews: http error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("GetGoogleNews: read error: %w", err)
	}

	var rss RSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		return nil, fmt.Errorf("GetGoogleNews: parse error: %w", err)
	}

	feeds := make([]models.GNewsItem, 0, len(rss.Channel.Items))
	for _, item := range rss.Channel.Items {
		title := item.Title
		source := item.Source
		if strings.HasSuffix(title, " - "+source) {
			title = strings.TrimSuffix(title, " - "+source)
		}
		feeds = append(feeds, models.GNewsItem{
			Title:     title,
			Published: item.PubDate,
			Source:    source,
		})
	}

	return feeds, nil
}

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

// SentimentAnalyzeTitle calls the sentiment analysis gRPC service for the given title.
func SentimentAnalyzeTitles(ctx context.Context, items []string, lang string) ([]*sentimentpb.AnalyzeResponse, error) {
	if SentimentClient == nil {
		return nil, fmt.Errorf("SentimentClient not initialized (call InitSentimentClient first)")
	}

	// timeout per RPC
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	reqItems := make([]*sentimentpb.AnalyzeRequest, 0, len(items))
	for _, t := range items {
		reqItems = append(reqItems, &sentimentpb.AnalyzeRequest{
			Text:     t,
			Language: lang,
		})
	}

	resp, err := SentimentClient.BatchAnalyze(ctx, &sentimentpb.BatchAnalyzeRequest{
		Items: reqItems,
	})
	if err != nil {
		return nil, fmt.Errorf("BatchAnalyze failed: %w", err)
	}

	return resp.Results, nil
}

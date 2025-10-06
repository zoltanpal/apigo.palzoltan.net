// utils/rss.go
package utils

import (
	"encoding/xml"
	"fmt"
	"golang-restapi/models"
	"io"
	"net/http"
	"net/url"
	"strings"
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

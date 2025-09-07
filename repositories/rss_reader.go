package repositories

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"golang-restapi/config"
	"golang-restapi/db"
	"golang-restapi/models"
	"golang-restapi/queries"
	"golang-restapi/utils"
)

func RssReader(ctx context.Context, cfg config.Config, lang string) []models.RSSItem {
	// Query sources by language
	rows, err := db.DB.QueryContext(ctx, queries.SourcesByLanguage, lang)
	if err != nil {
		log.Printf("RssReader: query error: %v", err)
		return nil
	}
	defer rows.Close()

	// Load sources into memory (or stream them onto the work chan directly)
	var sources []models.Source
	for rows.Next() {
		var s models.Source
		// Adjust this scan order/types to your schema:
		// Expecting: id, rss, lang
		if err := rows.Scan(&s.ID, &s.Rss, &s.Lang); err != nil {
			log.Printf("RssReader: scan error: %v", err)
			continue
		}
		sources = append(sources, s)
	}
	if err := rows.Err(); err != nil {
		log.Printf("RssReader: rows error: %v", err)
	}

	reader := utils.NewReader()
	scli := utils.New(cfg.SentimentURL, cfg.SentimentToken, 36*time.Second)

	workers := 6
	work := make(chan models.Source)
	results := make(chan []models.RSSItem)

	var wg sync.WaitGroup
	wg.Add(workers)

	// Start workers
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for s := range work {
				fmt.Printf("Source: %s (lang: %s)\n", s.Rss, s.Lang)

				items, err := reader.FetchSource(ctx, int64(s.ID), s.Rss, s.Lang)
				if err != nil {
					log.Printf("RssReader: fetch error for source %d: %v\n", s.ID, err)
					continue
				}

				for i := range items {
					it := items[i]

					// Use the source language (or fallback to function arg `lang` if that’s your design)
					sentiments, err := scli.Analyze(ctx, s.Lang, it.Title)
					if err != nil {
						log.Printf("sentiment analyze error (source %d): %v", s.ID, err)
						continue
					}

					var bestKey string
					var bestVal float64
					first := true

					for k, v := range sentiments {
						if k == "compound" {
							continue
						}
						if first || v > bestVal {
							bestKey = k
							bestVal = v
							first = false
						}
					}

					// Normalize "very_positive"/"very_negative" -> "positive"/"negative"
					if strings.HasPrefix(bestKey, "very_") {
						bestKey = strings.TrimPrefix(bestKey, "very_")
					}

					items[i].SentimentKey = bestKey
					items[i].SentimentValue = bestVal
				}

				// Send processed items to collector
				if len(items) > 0 {
					results <- items
				}
			}
		}()
	}

	// Feed work
	go func() {
		for _, s := range sources {
			work <- s
		}
		close(work)
	}()

	// Close results after all workers exit
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var feeds []models.RSSItem
	for batch := range results {
		feeds = append(feeds, batch...)
	}

	return feeds
}

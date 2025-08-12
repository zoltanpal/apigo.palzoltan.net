package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang-restapi/repositories"
	"golang-restapi/utils"
)

// GetFeedsWords handles GET /pow/words?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func GetFeedsWords(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Query parameters start_date and end_date are required (format: YYYY-MM-DD)",
		})
		return
	}

	const layout = "2006-01-02"
	if _, err := time.Parse(layout, startDateStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		return
	}
	if _, err := time.Parse(layout, endDateStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
		return
	}

	words, err := repositories.GetFeedsWords(startDateStr, endDateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch feed words"})
		return
	}
	c.JSON(http.StatusOK, words)
}

// GetSentimentGrouped handles GET /pow/get_sentiment_grouped?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func GetSentimentGrouped(c *gin.Context) {
	const layout = "2006-01-02"

	start := c.Query("start_date")
	end := c.Query("end_date")
	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required (YYYY-MM-DD)"})
		return
	}
	if _, err := time.Parse(layout, start); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}
	if _, err := time.Parse(layout, end); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	freeText := c.DefaultQuery("free_text", "")
	groupBy := c.DefaultQuery("group_by", "source")

	rows, err := repositories.GetSentimentGrouped(
		c.Request.Context(),
		start, end,
		freeText,
		groupBy,
	)
	if err != nil {
		// In dev, you can expose err.Error(); in prod, return generic.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch grouped sentiments"})
		return
	}

	series := utils.GenerateSentimentSeries(rows)
	c.JSON(http.StatusOK, series)
}

// CountSentiments [GET] /pow/count_sentiments?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func CountSentiments(c *gin.Context) {
	const layout = "2006-01-02"
	start := c.Query("start_date")
	end := c.Query("end_date")

	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required (YYYY-MM-DD)"})
		return
	}
	if _, err := time.Parse(layout, start); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}
	if _, err := time.Parse(layout, end); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	counts, err := repositories.CountSentiments(c.Request.Context(), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count sentiments"})
		return
	}

	c.JSON(http.StatusOK, counts)
}

// GET /top_feeds?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&pos_neg=positive&limit=5
func TopFeeds(c *gin.Context) {
	const layout = "2006-01-02"

	start := c.Query("start_date")
	end := c.Query("end_date")
	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required (YYYY-MM-DD)"})
		return
	}
	if _, err := time.Parse(layout, start); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}
	if _, err := time.Parse(layout, end); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	posNeg := strings.ToLower(c.DefaultQuery("pos_neg", "positive"))
	switch posNeg {
	case "positive", "negative", "neutral":
		// ok
	default:
		posNeg = "positive"
	}

	limitVal, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil || limitVal < 1 {
		limitVal = 5
	}

	rows, repoErr := repositories.TopFeeds(c.Request.Context(), start, end, posNeg, limitVal)
	if repoErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch top feeds"})
		return
	}

	c.JSON(http.StatusOK, rows)
}

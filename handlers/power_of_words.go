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

// BiasDetection GET /bias_detection?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&words=a,b,c[&sources=1,2]
func BiasDetection(c *gin.Context) {
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

	words := c.QueryArray("words")
	if len(words) == 0 {
		words = utils.ParseStringList(c.Query("words"))
	}
	if len(words) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "words is required"})
		return
	}

	rows, err := repositories.BiasDetection(
		c.Request.Context(),
		start, end,
		words,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to compute bias detection"})
		return
	}

	c.JSON(http.StatusOK, rows)
}

// CorrelationBetweenSourcesAvgCompound GET /correlation_between_sources_avg_compound
func CorrelationBetweenSourcesAvgCompound(c *gin.Context) {
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

	word := strings.TrimSpace(c.Query("word"))
	if word == "" || len(strings.Fields(word)) != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "word is required and must be a single token"})
		return
	}

	// sources: optional, CSV or repeated ints
	srcIDs := []int{}
	if arr := c.QueryArray("sources"); len(arr) > 0 {
		for _, s := range arr {
			if v, err := strconv.Atoi(s); err == nil {
				srcIDs = append(srcIDs, v)
			}
		}
	} else {
		srcIDs = utils.ParseIntList(c.Query("sources"))
	}

	rows, err := repositories.CorrelationBetweenSourcesAvgCompound(
		c.Request.Context(),
		start, end,
		word,
		srcIDs,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to compute correlation"})
		return
	}
	c.JSON(http.StatusOK, rows)
}

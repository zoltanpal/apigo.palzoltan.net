package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang-restapi/repositories"
	"golang-restapi/utils"

	"github.com/gin-gonic/gin"
)

const feedDateLayout = "2006-01-02"

// GET /feeds?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&sources=1,2&free_text=word&page=1&items_per_page=30
func GetFeeds(c *gin.Context) {

	start := c.Query("start_date")
	end := c.Query("end_date")

	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}
	if _, err := time.Parse(feedDateLayout, start); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}
	if _, err := time.Parse(feedDateLayout, end); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	// query params
	srcIDs := utils.ParseIntList(c.Query("sources"))
	freeText := c.Query("free_text")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("items_per_page", "30"))

	resp, err := repositories.GetFeeds(
		c.Request.Context(),
		start, end,
		srcIDs,
		freeText,
		page, perPage,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GET /most_common_words?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&nm_common=20
func MostCommonWordsHandler(c *gin.Context) {

	start := c.Query("start_date")
	end := c.Query("end_date")
	nmStr := c.DefaultQuery("nm_common", "20")

	if _, err := time.Parse(feedDateLayout, start); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
		return
	}
	if _, err := time.Parse(feedDateLayout, end); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
		return
	}

	n, err := strconv.Atoi(nmStr)
	if err != nil || n <= 0 {
		n = 20
	}

	result, repoErr := repositories.MostCommonWords(c.Request.Context(), start, end, n)
	if repoErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": repoErr.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetSentimentGrouped handles GET /pow/get_sentiment_grouped?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
func GetSentimentGrouped(c *gin.Context) {

	start := c.Query("start_date")
	end := c.Query("end_date")
	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required (YYYY-MM-DD)"})
		return
	}
	if _, err := time.Parse(feedDateLayout, start); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}
	if _, err := time.Parse(feedDateLayout, end); err != nil {
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
	start := c.Query("start_date")
	end := c.Query("end_date")

	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required (YYYY-MM-DD)"})
		return
	}
	if _, err := time.Parse(feedDateLayout, start); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}
	if _, err := time.Parse(feedDateLayout, end); err != nil {
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

	start := c.Query("start_date")
	end := c.Query("end_date")
	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required (YYYY-MM-DD)"})
		return
	}
	if _, err := time.Parse(feedDateLayout, start); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}
	if _, err := time.Parse(feedDateLayout, end); err != nil {
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

	start := c.Query("start_date")
	end := c.Query("end_date")
	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required (YYYY-MM-DD)"})
		return
	}
	if _, err := time.Parse(feedDateLayout, start); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}
	if _, err := time.Parse(feedDateLayout, end); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	word := strings.TrimSpace(c.Query("word"))
	if word == "" || len(strings.Fields(word)) != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "word is required and must be a single token"})
		return
	}

	rows, err := repositories.BiasDetection(
		c.Request.Context(),
		start, end,
		word,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to compute bias detection"})
		return
	}

	c.JSON(http.StatusOK, rows)
}

// CorrelationBetweenSourcesAvgCompound GET /correlation_between_sources_avg_compound
func CorrelationBetweenSourcesAvgCompound(c *gin.Context) {

	start := c.Query("start_date")
	end := c.Query("end_date")
	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required (YYYY-MM-DD)"})
		return
	}
	if _, err := time.Parse(feedDateLayout, start); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}
	if _, err := time.Parse(feedDateLayout, end); err != nil {
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

// WordCoOccurrences GET /word_co_occurences?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD&word=foo[&sources=1,2]
func WordCoOccurrences(c *gin.Context) {
	start := c.Query("start_date")
	end := c.Query("end_date")
	if start == "" || end == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required (YYYY-MM-DD)"})
		return
	}
	if _, err := time.Parse(feedDateLayout, start); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}
	if _, err := time.Parse(feedDateLayout, end); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	word := c.Query("word")
	if word == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "word parameter is required"})
		return
	}

	// sources
	srcIDs := []int{}
	srcIDs = utils.ParseIntList(c.Query("sources"))

	rows, err := repositories.WordCoOccurrences(
		c.Request.Context(),
		start, end,
		word,
		srcIDs,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to compute co-occurrences"})
		return
	}

	c.JSON(http.StatusOK, rows)
}

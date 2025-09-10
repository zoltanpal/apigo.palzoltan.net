package routes

import (
	"golang-restapi/handlers"
	"time"

	//"golang-restapi/middlewares"

	"golang-restapi/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRoutes # Function to handle the API routes
func SetupRoutes(r *gin.Engine, cfg config.Config) {
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5000",
			"http://127.0.0.1:5000",
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"https://palzoltan.net",
			"https://pow.palzoltan.net",
			"http://devpow.palzoltan.net",
		},
		AllowMethods:     []string{"GET"}, // "PUT", "DELETE", "PATCH", "POST",  "OPTIONS"
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	protected := r.Group("/")
	//protected.Use(middlewares.Authentication())

	protected.GET("/pow/feeds", handlers.GetFeeds)
	protected.GET("/pow/most_common_words", handlers.MostCommonWordsHandler)
	protected.GET("/pow/get_sentiment_grouped", handlers.GetSentimentGrouped)
	protected.GET("/pow/count_sentiments", handlers.CountSentiments)
	protected.GET("/pow/top_feeds", handlers.TopFeeds)
	protected.GET("/pow/bias_detection", handlers.BiasDetection)
	protected.GET("/pow/correlation_between_sources_avg_compound", handlers.CorrelationBetweenSourcesAvgCompound)
	protected.GET("/pow/word_co_occurences", handlers.WordCoOccurrences)
	protected.GET("/pow/phrase_frequency_trends", handlers.PhraseFrequencyTrends)
}

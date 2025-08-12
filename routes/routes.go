package routes

import (
	"golang-restapi/handlers"
	//"golang-restapi/middlewares"

	"github.com/gin-gonic/gin"
	"golang-restapi/config"
)

// Function to handle the API routes
func SetupRoutes(r *gin.Engine, cfg config.Config) {
	protected := r.Group("/")
	// protected.Use(middlewares.AuthMiddleware())

	protected.GET("/pow/feed_words", handlers.GetFeedsWords)
	protected.GET("/pow/get_sentiment_grouped", handlers.GetSentimentGrouped)
	protected.GET("/pow/count_sentiments", handlers.CountSentiments)
}

package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"golang-restapi/config"
	"golang-restapi/db"
	"golang-restapi/middlewares"
	"golang-restapi/routes"
)

func main() {
	// Load config & set Gin mode
	cfg := config.LoadConfig()
	//gin.SetMode(gin.DebugMode)

	// Init DB
	db.InitDB(cfg)
	defer db.DB.Close()

	// router init
	router := gin.New()
	router.Use(
		gin.Logger(),
		gin.Recovery(),
		middlewares.ErrorLogger(),
	)

	// Register your routes
	routes.SetupRoutes(router, cfg)

	if err := router.Run(":1985"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

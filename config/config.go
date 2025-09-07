package config

import (
	//"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	AuthIssuers        []string
	AuthEmails         []string
	USGSApiHost        string
	IMDBApiHost        string
	IMDBApiKey         string
	CORSAllowedOrigins string
	APP_PORT           string
	SentimentURL       string
	SentimentToken     string
}

// LoadConfig loads environment variables from .env
func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "1985" // default
	}

	var issuers, emails []string

	if envIssuers := os.Getenv("ALLOWED_ISSUERS"); envIssuers != "" {
		issuers = strings.SplitAfter(envIssuers, ",")
	}

	if envEmails := os.Getenv("ALLOWED_EMAILS"); envEmails != "" {
		emails = strings.SplitAfter(envEmails, ",")
	}

	return Config{
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             os.Getenv("DB_PORT"),
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		CORSAllowedOrigins: os.Getenv("CORS_ALLOWED_ORIGINS"),
		AuthIssuers:        issuers,
		AuthEmails:         emails,
		USGSApiHost:        os.Getenv("USGS_API_HOST"),
		IMDBApiHost:        os.Getenv("IMDB_BASE_URL"),
		IMDBApiKey:         os.Getenv("IMDB_API_KEY"),
		SentimentURL:       os.Getenv("SENTIMENT_URL"),
		SentimentToken:     os.Getenv("SENTIMENT_TOKEN"),
		APP_PORT:           appPort,
	}
}

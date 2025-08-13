package config

import (
	//"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

// Config holds application configuration
type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	AuthIssuers   []string
	AuthEmails    []string
	USGS_API_HOST string
	IMDB_BASE_URL string
	IMDB_API_KEY  string
}

// LoadConfig loads environment variables from .env
func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var issuers, emails []string

	if envIssuers := os.Getenv("ALLOWED_ISSUERS"); envIssuers != "" {
		issuers = strings.SplitAfter(envIssuers, ",")
	}

	if envEmails := os.Getenv("ALLOWED_EMAILS"); envEmails != "" {
		emails = strings.SplitAfter(envEmails, ",")
	}

	return Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		AuthIssuers:   issuers,
		AuthEmails:    emails,
		USGS_API_HOST: os.Getenv("USGS_API_HOST"),
		IMDB_BASE_URL: os.Getenv("IMDB_BASE_URL"),
		IMDB_API_KEY:  os.Getenv("IMDB_API_KEY"),
	}
}

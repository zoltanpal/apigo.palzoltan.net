package db

import (
	"database/sql"

	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"golang-restapi/config"
)

var DB *sql.DB

func InitDB(cfg config.Config) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var err error

	DB, err = sql.Open("postgres", dsn)

	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Set connection pool parameters
	DB.SetMaxOpenConns(10)                  // Max open connections
	DB.SetMaxIdleConns(5)                   // Max idle connections
	DB.SetConnMaxLifetime(10 * time.Minute) // lifetime to 10 minutes

	err = DB.Ping()
	if err != nil {
		log.Fatal("Database connection is not alive:", err)
	}

	fmt.Println("Connected to the database successfully!")
}

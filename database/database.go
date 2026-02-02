package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"kasir-api/config"
)

var DB *sql.DB

func InitDB() {
	var err error

	if config.AppConfig.DBConn == "" {
		log.Println("DB_CONN not set, running without database (in-memory mode)")
		return
	}

	DB, err = sql.Open("pgx", config.AppConfig.DBConn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	// Connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connected successfully")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

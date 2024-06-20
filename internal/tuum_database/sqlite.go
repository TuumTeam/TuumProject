package tuum_database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"tuum.com/internal/config"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)

	connMaxLifetime, err := time.ParseDuration(cfg.Database.ConnMaxLifetime)
	if err != nil {
		log.Fatalf("Failed to parse conn_max_lifetime: %v", err)
	}
	db.SetConnMaxLifetime(connMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connection established successfully.")
	return db, nil
}

package database

import (
	"database/sql"
	"fmt"
	"time"

	"task-engine/config"

	_ "github.com/lib/pq" // anonymous import: registers the PostgreSQL driver in database/sql via init(), even though it's not used directly
)

var DB *sql.DB

func Connect(cfg *config.Config) error {
	db, err := sql.Open("postgres", cfg.GetDatabaseURL())
	if err != nil {
		return err
	}

	// Connections pool config
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Testing connection
	if err := db.Ping(); err != nil {
		return err
	}

	DB = db
	return nil
}

func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("DB not connected")
	}
	return DB.Ping()
}

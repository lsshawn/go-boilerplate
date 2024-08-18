package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tursodatabase/go-libsql"
)

var DB *sql.DB

func Init() error {
	dbName := "local.db"
	primaryURL := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_DATABASE_AUTH_TOKEN")

	if primaryURL == "" {
		return fmt.Errorf("TURSO_DATABASE_URL environment variable is not set")
	}

	if authToken == "" {
		return fmt.Errorf("TURSO_AUTH_TOKEN environment variable is not set")
	}

	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		return fmt.Errorf("error creating temporary directory: %w", err)
	}

	dbPath := filepath.Join(dir, dbName)

	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryURL,
		libsql.WithAuthToken(authToken),
	)
	if err != nil {
		return fmt.Errorf("error creating connector: %w", err)
	}

	DB = sql.OpenDB(connector)

	// Test the connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	fmt.Println("Successfully connected to Turso database")
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

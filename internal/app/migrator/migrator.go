package migrator

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/Shemistan/manager/internal/config"
	"github.com/joho/godotenv"
)

// Run initializes and runs the database migrator.
func Run() error {
	// Load .env file for local development (ignore if not exists).
	_ = godotenv.Load()

	// Load configuration.
	cfg, err := config.Load("app.toml")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger.
	logger := log.New(os.Stdout, "[migrator] ", log.LstdFlags|log.Lshortfile)

	// Connect to database.
	dsn := config.BuildDSN(cfg)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close() // nolint:errcheck

	// Verify database connection.
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	logger.Println("connected to database")

	// Read migration files.
	migrationDir := "migration"
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %w", err)
	}

	// Filter and sort SQL files.
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	// Execute migrations.
	for _, fileName := range sqlFiles {
		filePath := filepath.Join(migrationDir, fileName)
		logger.Printf("executing migration: %s", fileName)

		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", fileName, err)
		}

		logger.Printf("completed migration: %s", fileName)
	}

	logger.Println("all migrations completed successfully")
	return nil
}

package manager

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Shemistan/manager/internal/api/manager"
	"github.com/Shemistan/manager/internal/config"
	servicemanager "github.com/Shemistan/manager/internal/service/manager"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Run initializes and runs the manager application.
func Run() error {
	// Load .env file for local development (ignore if not exists).
	_ = godotenv.Load()

	// Load configuration.
	cfg, err := config.Load("app.toml")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize logger.
	logger := log.New(os.Stdout, "[manager] ", log.LstdFlags|log.Lshortfile)

	//// Connect to database.
	//dsn := config.BuildDSN(cfg)
	//db, err := sql.Open("postgres", dsn)
	//if err != nil {
	//	return fmt.Errorf("failed to open database: %w", err)
	//}
	//defer db.Close() // nolint:errcheck
	//
	//// Verify database connection.
	//if err := db.Ping(); err != nil {
	//	return fmt.Errorf("failed to ping database: %w", err)
	//}
	//logger.Println("connected to database")
	//
	//// Initialize storage.
	//healthStorage := storagemanager.NewHealthStorage(db)

	// Initialize service.
	healthService := servicemanager.NewHealthService(nil)

	// Initialize API handler.
	handler := manager.NewHandler(healthService, logger)
	router := handler.RegisterRoutes()

	// Start HTTP server.
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Printf("starting HTTP server on port %d", cfg.HTTPPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

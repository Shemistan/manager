package manager

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Shemistan/manager/internal/service"
)

// Handler manages HTTP handlers for the manager API.
type Handler struct {
	healthService service.HealthService
	logger        *log.Logger
}

// NewHandler creates a new Handler instance.
func NewHandler(healthService service.HealthService, logger *log.Logger) *Handler {
	return &Handler{
		healthService: healthService,
		logger:        logger,
	}
}

// RegisterRoutes registers all API routes.
func (h *Handler) RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	return mux
}

// Health handles the GET /health endpoint.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	// Call the service layer.
	err := h.healthService.HandleHealth(ctx)
	if err != nil {
		h.logger.Printf("error handling health check: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status: "error",
			Error:  "failed to record health check",
		})
		return
	}

	// Return success response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HealthResponse{
		Status: "success",
	})
}

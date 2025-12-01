package service

import "context"

// HealthService defines the interface for health check business logic.
type HealthService interface {
	// HandleHealth processes a health check request.
	HandleHealth(ctx context.Context) error
}

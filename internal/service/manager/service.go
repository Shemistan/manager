package manager

import (
	"context"

	"github.com/Shemistan/manager/internal/storage"
)

// HealthService is the business logic implementation for health checks.
type HealthService struct {
	storage storage.HealthStorage
}

// NewHealthService creates a new HealthService instance.
func NewHealthService(storage storage.HealthStorage) *HealthService {
	return &HealthService{
		storage: storage,
	}
}

// HandleHealth processes a health check request and saves it to storage.
func (s *HealthService) HandleHealth(ctx context.Context) error {
	return s.storage.SaveHealthCall(ctx)
}

package storage

import "context"

// HealthStorage defines the interface for health call logging.
type HealthStorage interface {
	// SaveHealthCall records a health check call in the database.
	SaveHealthCall(ctx context.Context) error
}

package manager

import (
	"context"
	"database/sql"
)

// HealthStorage is a PostgreSQL implementation of the HealthStorage interface.
type HealthStorage struct {
	db *sql.DB
}

// NewHealthStorage creates a new HealthStorage instance.
func NewHealthStorage(db *sql.DB) *HealthStorage {
	return &HealthStorage{
		db: db,
	}
}

// SaveHealthCall records a health check call in the health_calls table.
func (s *HealthStorage) SaveHealthCall(ctx context.Context) error {
	query := `INSERT INTO health_calls (called_at) VALUES (NOW())`
	_, err := s.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

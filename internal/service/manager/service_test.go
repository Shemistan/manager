package manager

import (
	"context"
	"errors"
	"testing"
)

// MockHealthStorage is a mock implementation of HealthStorage for testing.
type MockHealthStorage struct {
	SaveHealthCallFunc func(ctx context.Context) error
}

func (m *MockHealthStorage) SaveHealthCall(ctx context.Context) error {
	if m.SaveHealthCallFunc != nil {
		return m.SaveHealthCallFunc(ctx)
	}
	return nil
}

func TestHandleHealth_Success(t *testing.T) {
	mock := &MockHealthStorage{
		SaveHealthCallFunc: func(ctx context.Context) error {
			return nil
		},
	}

	service := NewHealthService(mock)
	err := service.HandleHealth(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestHandleHealth_StorageError(t *testing.T) {
	expectedErr := errors.New("storage error")
	mock := &MockHealthStorage{
		SaveHealthCallFunc: func(ctx context.Context) error {
			return expectedErr
		},
	}

	service := NewHealthService(mock)
	err := service.HandleHealth(context.Background())

	if err != expectedErr {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
}

func TestHandleHealth_ContextCancel(t *testing.T) {
	mock := &MockHealthStorage{
		SaveHealthCallFunc: func(ctx context.Context) error {
			return ctx.Err()
		},
	}

	service := NewHealthService(mock)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := service.HandleHealth(ctx)

	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}

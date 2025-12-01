package manager

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// MockHealthService is a mock implementation of HealthService for testing.
type MockHealthService struct {
	HandleHealthFunc func(ctx context.Context) error
}

func (m *MockHealthService) HandleHealth(ctx context.Context) error {
	if m.HandleHealthFunc != nil {
		return m.HandleHealthFunc(ctx)
	}
	return nil
}

func TestHealthHandler_Success(t *testing.T) {
	mock := &MockHealthService{
		HandleHealthFunc: func(ctx context.Context) error {
			return nil
		},
	}

	logger := log.New(os.Stdout, "[test] ", log.LstdFlags)
	handler := NewHandler(mock, logger)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp HealthResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "success" {
		t.Fatalf("expected status 'success', got '%s'", resp.Status)
	}
}

func TestHealthHandler_ServiceError(t *testing.T) {
	mock := &MockHealthService{
		HandleHealthFunc: func(ctx context.Context) error {
			return context.DeadlineExceeded
		},
	}

	logger := log.New(os.Stdout, "[test] ", log.LstdFlags)
	handler := NewHandler(mock, logger)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}

	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "error" {
		t.Fatalf("expected status 'error', got '%s'", resp.Status)
	}
}

func TestHealthHandler_ContentType(t *testing.T) {
	mock := &MockHealthService{
		HandleHealthFunc: func(ctx context.Context) error {
			return nil
		},
	}

	logger := log.New(os.Stdout, "[test] ", log.LstdFlags)
	handler := NewHandler(mock, logger)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.Health(w, req)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected Content-Type 'application/json', got '%s'", contentType)
	}
}

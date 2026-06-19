package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jory/urlwatch/internal/domain"
)

func TestMemoryStore_SaveAndGet(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	batch := domain.Batch{
		ID:        "b_abc123",
		CreatedAt: time.Now().UTC(),
		Summary:   domain.Summary{Total: 1, Up: 1, Down: 0, DurationMs: 42},
		Results:   []domain.CheckResult{{URL: "https://go.dev", StatusCode: 200, OK: true, LatencyMs: 42}},
	}

	if err := s.Save(ctx, batch); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	got, err := s.Get(ctx, "b_abc123")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if got.ID != batch.ID {
		t.Errorf("ID = %q, want %q", got.ID, batch.ID)
	}
	if got.Summary.Total != 1 {
		t.Errorf("Total = %d, want 1", got.Summary.Total)
	}
}

func TestMemoryStore_GetNotFound(t *testing.T) {
	s := NewMemoryStore()
	ctx := context.Background()

	_, err := s.Get(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrBatchNotFound) {
		t.Errorf("expected ErrBatchNotFound, got: %v", err)
	}
}

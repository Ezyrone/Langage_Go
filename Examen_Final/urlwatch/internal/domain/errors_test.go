package domain

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrBatchNotFound_Sentinel(t *testing.T) {
	wrapped := fmt.Errorf("batch %q: %w", "b_abc", ErrBatchNotFound)

	if !errors.Is(wrapped, ErrBatchNotFound) {
		t.Error("errors.Is should match ErrBatchNotFound through wrapping")
	}
}

func TestValidationError_ErrorsAs(t *testing.T) {
	err := fmt.Errorf("request failed: %w", &ValidationError{Field: "urls", Message: "too many"})

	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatal("errors.As should unwrap ValidationError")
	}
	if ve.Field != "urls" {
		t.Errorf("Field = %q, want %q", ve.Field, "urls")
	}
}

func TestValidationError_Message(t *testing.T) {
	err := &ValidationError{Field: "concurrency", Message: "must be between 1 and 50"}
	expected := `validation error on field "concurrency": must be between 1 and 50`
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}

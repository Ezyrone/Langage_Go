package api

import (
	"errors"
	"testing"

	"github.com/jory/urlwatch/internal/domain"
)

func intPtr(i int) *int { return &i }

func TestValidateAndNormalize(t *testing.T) {
	tests := []struct {
		name    string
		req     CheckRequest
		wantErr bool
		errField string
	}{
		{
			name:    "valid minimal",
			req:     CheckRequest{URLs: []string{"https://go.dev"}},
			wantErr: false,
		},
		{
			name:    "empty urls",
			req:     CheckRequest{URLs: []string{}},
			wantErr: true,
			errField: "urls",
		},
		{
			name:    "invalid url scheme",
			req:     CheckRequest{URLs: []string{"ftp://example.com"}},
			wantErr: true,
			errField: "urls",
		},
		{
			name: "concurrency too high",
			req: CheckRequest{
				URLs:    []string{"https://go.dev"},
				Options: &CheckOptions{Concurrency: intPtr(100)},
			},
			wantErr:  true,
			errField: "options.concurrency",
		},
		{
			name: "concurrency too low",
			req: CheckRequest{
				URLs:    []string{"https://go.dev"},
				Options: &CheckOptions{Concurrency: intPtr(0)},
			},
			wantErr:  true,
			errField: "options.concurrency",
		},
		{
			name: "timeout too low",
			req: CheckRequest{
				URLs:    []string{"https://go.dev"},
				Options: &CheckOptions{TimeoutMs: intPtr(10)},
			},
			wantErr:  true,
			errField: "options.timeout_ms",
		},
		{
			name: "timeout too high",
			req: CheckRequest{
				URLs:    []string{"https://go.dev"},
				Options: &CheckOptions{TimeoutMs: intPtr(50000)},
			},
			wantErr:  true,
			errField: "options.timeout_ms",
		},
		{
			name: "valid with options",
			req: CheckRequest{
				URLs:    []string{"https://go.dev", "https://example.com"},
				Options: &CheckOptions{Concurrency: intPtr(4), TimeoutMs: intPtr(2000)},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := validateAndNormalize(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil {
				var ve *domain.ValidationError
				if !errors.As(err, &ve) {
					t.Errorf("expected ValidationError, got %T", err)
				} else if ve.Field != tt.errField {
					t.Errorf("field = %q, want %q", ve.Field, tt.errField)
				}
			}
		})
	}
}

func TestValidateAndNormalize_Defaults(t *testing.T) {
	_, opts, err := validateAndNormalize(CheckRequest{URLs: []string{"https://go.dev"}})
	if err != nil {
		t.Fatal(err)
	}
	if opts.Concurrency != 8 {
		t.Errorf("default concurrency = %d, want 8", opts.Concurrency)
	}
	if opts.TimeoutMs != 5000 {
		t.Errorf("default timeout = %d, want 5000", opts.TimeoutMs)
	}
}

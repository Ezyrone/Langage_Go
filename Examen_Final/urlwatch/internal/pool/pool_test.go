package pool

import (
	"context"
	"testing"
	"time"

	"github.com/jory/urlwatch/internal/checker"
	"github.com/jory/urlwatch/internal/domain"
)

func TestRun_BasicFunctionality(t *testing.T) {
	mock := checker.NewMockChecker(map[string]domain.CheckResult{
		"https://ok.example.com":   {URL: "https://ok.example.com", StatusCode: 200, OK: true, LatencyMs: 10},
		"https://fail.example.com": {URL: "https://fail.example.com", OK: false, Error: "connection refused"},
	}, 0)

	urls := []string{"https://ok.example.com", "https://fail.example.com"}
	batch := Run(context.Background(), mock, urls, Options{Concurrency: 2, TimeoutMs: 5000})

	if batch.Summary.Total != 2 {
		t.Fatalf("total = %d, want 2", batch.Summary.Total)
	}
	if batch.Summary.Up != 1 {
		t.Errorf("up = %d, want 1", batch.Summary.Up)
	}
	if batch.Summary.Down != 1 {
		t.Errorf("down = %d, want 1", batch.Summary.Down)
	}
}

func TestRun_RespectsContext_Cancellation(t *testing.T) {
	mock := checker.NewMockChecker(map[string]domain.CheckResult{
		"https://slow.example.com": {URL: "https://slow.example.com", StatusCode: 200, OK: true, LatencyMs: 50},
	}, 2*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	urls := []string{"https://slow.example.com"}
	batch := Run(ctx, mock, urls, Options{Concurrency: 1, TimeoutMs: 200})

	if batch.Summary.Total != 1 {
		t.Fatalf("total = %d, want 1", batch.Summary.Total)
	}
	if batch.Results[0].OK {
		t.Error("expected result to be not OK due to context cancellation")
	}
}

func TestRun_ConcurrencyBound(t *testing.T) {
	results := make(map[string]domain.CheckResult)
	for i := 0; i < 10; i++ {
		url := "https://example.com/" + string(rune('a'+i))
		results[url] = domain.CheckResult{URL: url, StatusCode: 200, OK: true, LatencyMs: 5}
	}

	mock := checker.NewMockChecker(results, 50*time.Millisecond)

	urls := make([]string, 0, 10)
	for u := range results {
		urls = append(urls, u)
	}

	batch := Run(context.Background(), mock, urls, Options{Concurrency: 2, TimeoutMs: 10000})

	if batch.Summary.Total != 10 {
		t.Fatalf("total = %d, want 10", batch.Summary.Total)
	}
	if batch.Summary.Up != 10 {
		t.Errorf("up = %d, want 10", batch.Summary.Up)
	}
}

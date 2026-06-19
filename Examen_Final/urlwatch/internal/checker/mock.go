package checker

import (
	"context"
	"time"

	"github.com/jory/urlwatch/internal/domain"
)

type MockChecker struct {
	Results map[string]domain.CheckResult
	Delay   time.Duration
}

func NewMockChecker(results map[string]domain.CheckResult, delay time.Duration) *MockChecker {
	return &MockChecker{Results: results, Delay: delay}
}

func (m *MockChecker) Check(ctx context.Context, url string) domain.CheckResult {
	if m.Delay > 0 {
		select {
		case <-time.After(m.Delay):
		case <-ctx.Done():
			return domain.CheckResult{
				URL:       url,
				OK:        false,
				LatencyMs: m.Delay.Milliseconds(),
				Error:     ctx.Err().Error(),
			}
		}
	}

	if r, ok := m.Results[url]; ok {
		return r
	}
	return domain.CheckResult{
		URL:       url,
		OK:        false,
		LatencyMs: 0,
		Error:     "unknown url in mock",
	}
}

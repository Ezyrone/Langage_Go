package checker

import (
	"context"
	"net/http"
	"time"

	"github.com/jory/urlwatch/internal/domain"
)

type HTTPChecker struct {
	Client *http.Client
}

func NewHTTPChecker() *HTTPChecker {
	return &HTTPChecker{
		Client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (c *HTTPChecker) Check(ctx context.Context, url string) domain.CheckResult {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.CheckResult{
			URL:       url,
			OK:        false,
			LatencyMs: time.Since(start).Milliseconds(),
			Error:     err.Error(),
		}
	}

	resp, err := c.Client.Do(req)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return domain.CheckResult{
			URL:       url,
			OK:        false,
			LatencyMs: latency,
			Error:     err.Error(),
		}
	}
	defer resp.Body.Close()

	return domain.CheckResult{
		URL:        url,
		StatusCode: resp.StatusCode,
		OK:         resp.StatusCode >= 200 && resp.StatusCode < 400,
		LatencyMs:  latency,
	}
}

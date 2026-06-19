package api

import (
	"net/url"
	"strconv"

	"github.com/jory/urlwatch/internal/domain"
	"github.com/jory/urlwatch/internal/pool"
)

const (
	defaultConcurrency = 8
	maxConcurrency     = 50
	defaultTimeoutMs   = 5000
	minTimeoutMs       = 100
	maxTimeoutMs       = 30000
	maxURLs            = 100
)

func validateAndNormalize(req CheckRequest) ([]string, pool.Options, error) {
	if len(req.URLs) == 0 {
		return nil, pool.Options{}, &domain.ValidationError{Field: "urls", Message: "must contain at least 1 URL"}
	}
	if len(req.URLs) > maxURLs {
		return nil, pool.Options{}, &domain.ValidationError{Field: "urls", Message: "must contain at most 100 URLs"}
	}

	for i, u := range req.URLs {
		parsed, err := url.ParseRequestURI(u)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			return nil, pool.Options{}, &domain.ValidationError{
				Field:   "urls",
				Message: "invalid URL at index " + strconv.Itoa(i) + ": must be http or https",
			}
		}
	}

	concurrency := defaultConcurrency
	timeoutMs := defaultTimeoutMs

	if req.Options != nil {
		if req.Options.Concurrency != nil {
			concurrency = *req.Options.Concurrency
			if concurrency < 1 || concurrency > maxConcurrency {
				return nil, pool.Options{}, &domain.ValidationError{
					Field:   "options.concurrency",
					Message: "must be between 1 and 50",
				}
			}
		}
		if req.Options.TimeoutMs != nil {
			timeoutMs = *req.Options.TimeoutMs
			if timeoutMs < minTimeoutMs || timeoutMs > maxTimeoutMs {
				return nil, pool.Options{}, &domain.ValidationError{
					Field:   "options.timeout_ms",
					Message: "must be between 100 and 30000",
				}
			}
		}
	}

	return req.URLs, pool.Options{Concurrency: concurrency, TimeoutMs: timeoutMs}, nil
}

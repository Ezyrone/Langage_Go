package pool

import (
	"context"
	"sync"
	"time"

	"github.com/jory/urlwatch/internal/domain"
)

type Options struct {
	Concurrency int
	TimeoutMs   int
}

func Run(ctx context.Context, checker domain.Checker, urls []string, opts Options) domain.Batch {
	start := time.Now()

	batchCtx, batchCancel := context.WithTimeout(ctx, time.Duration(opts.TimeoutMs)*time.Millisecond)
	defer batchCancel()

	jobs := make(chan string, len(urls))
	results := make(chan domain.CheckResult, len(urls))

	var wg sync.WaitGroup
	for i := 0; i < opts.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				perURLCtx, perURLCancel := context.WithTimeout(batchCtx, time.Duration(opts.TimeoutMs)*time.Millisecond)
				r := checker.Check(perURLCtx, url)
				perURLCancel()
				results <- r
			}
		}()
	}

	for _, u := range urls {
		jobs <- u
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	collected := make([]domain.CheckResult, 0, len(urls))
	for r := range results {
		collected = append(collected, r)
	}

	durationMs := time.Since(start).Milliseconds()

	return domain.Batch{
		CreatedAt: time.Now().UTC(),
		Summary:   domain.AggregateSummary(collected, durationMs),
		Results:   collected,
	}
}

package domain

import "testing"

func TestAggregateSummary(t *testing.T) {
	tests := []struct {
		name       string
		results    []CheckResult
		durationMs int64
		wantUp     int
		wantDown   int
		wantTotal  int
	}{
		{
			name:       "all up",
			results:    []CheckResult{{OK: true}, {OK: true}, {OK: true}},
			durationMs: 100,
			wantUp:     3, wantDown: 0, wantTotal: 3,
		},
		{
			name:       "all down",
			results:    []CheckResult{{OK: false}, {OK: false}},
			durationMs: 200,
			wantUp:     0, wantDown: 2, wantTotal: 2,
		},
		{
			name:       "mixed",
			results:    []CheckResult{{OK: true}, {OK: false}, {OK: true}},
			durationMs: 300,
			wantUp:     2, wantDown: 1, wantTotal: 3,
		},
		{
			name:       "empty",
			results:    []CheckResult{},
			durationMs: 0,
			wantUp:     0, wantDown: 0, wantTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := AggregateSummary(tt.results, tt.durationMs)
			if s.Total != tt.wantTotal {
				t.Errorf("Total = %d, want %d", s.Total, tt.wantTotal)
			}
			if s.Up != tt.wantUp {
				t.Errorf("Up = %d, want %d", s.Up, tt.wantUp)
			}
			if s.Down != tt.wantDown {
				t.Errorf("Down = %d, want %d", s.Down, tt.wantDown)
			}
			if s.DurationMs != tt.durationMs {
				t.Errorf("DurationMs = %d, want %d", s.DurationMs, tt.durationMs)
			}
		})
	}
}

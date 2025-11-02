package ui

import (
	"testing"
	"time"

	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

func TestCalculateStats(t *testing.T) {
	tests := []struct {
		name    string
		results []models.TestResult
		want    TestStats
	}{
		{
			name:    "empty results",
			results: []models.TestResult{},
			want: TestStats{
				Total:       0,
				Passed:      0,
				Failed:      0,
				FastestTime: 0,
			},
		},
		{
			name: "all passed",
			results: []models.TestResult{
				{Method: "GET", Endpoint: "/users", Status: "200", Duration: 100 * time.Millisecond},
				{Method: "POST", Endpoint: "/users", Status: "201", Duration: 150 * time.Millisecond},
				{Method: "GET", Endpoint: "/posts", Status: "200", Duration: 80 * time.Millisecond},
			},
			want: TestStats{
				Total:           3,
				Passed:          3,
				Failed:          0,
				AverageTime:     110 * time.Millisecond,
				TotalTime:       330 * time.Millisecond,
				FastestTime:     80 * time.Millisecond,
				SlowestTime:     150 * time.Millisecond,
				FastestEndpoint: "GET /posts",
				SlowestEndpoint: "POST /users",
			},
		},
		{
			name: "mixed results",
			results: []models.TestResult{
				{Method: "GET", Endpoint: "/users", Status: "200", Duration: 100 * time.Millisecond},
				{Method: "POST", Endpoint: "/users", Status: "500", Duration: 200 * time.Millisecond},
				{Method: "GET", Endpoint: "/posts", Status: "404", Duration: 50 * time.Millisecond},
			},
			want: TestStats{
				Total:           3,
				Passed:          1,
				Failed:          2,
				AverageTime:     116666666 * time.Nanosecond, // ~116.67ms
				TotalTime:       350 * time.Millisecond,
				FastestTime:     50 * time.Millisecond,
				SlowestTime:     200 * time.Millisecond,
				FastestEndpoint: "GET /posts",
				SlowestEndpoint: "POST /users",
			},
		},
		{
			name: "all failed",
			results: []models.TestResult{
				{Method: "GET", Endpoint: "/users", Status: "500", Duration: 100 * time.Millisecond},
				{Method: "POST", Endpoint: "/users", Status: "ERR", Duration: 50 * time.Millisecond},
			},
			want: TestStats{
				Total:           2,
				Passed:          0,
				Failed:          2,
				AverageTime:     75 * time.Millisecond,
				TotalTime:       150 * time.Millisecond,
				FastestTime:     50 * time.Millisecond,
				SlowestTime:     100 * time.Millisecond,
				FastestEndpoint: "POST /users",
				SlowestEndpoint: "GET /users",
			},
		},
		{
			name: "zero durations",
			results: []models.TestResult{
				{Method: "GET", Endpoint: "/users", Status: "200", Duration: 0},
				{Method: "POST", Endpoint: "/posts", Status: "201", Duration: 0},
			},
			want: TestStats{
				Total:       2,
				Passed:      2,
				Failed:      0,
				AverageTime: 0,
				TotalTime:   0,
				FastestTime: 0,
				SlowestTime: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateStats(tt.results)

			if got.Total != tt.want.Total {
				t.Errorf("Total = %d, want %d", got.Total, tt.want.Total)
			}
			if got.Passed != tt.want.Passed {
				t.Errorf("Passed = %d, want %d", got.Passed, tt.want.Passed)
			}
			if got.Failed != tt.want.Failed {
				t.Errorf("Failed = %d, want %d", got.Failed, tt.want.Failed)
			}
			if got.AverageTime != tt.want.AverageTime {
				t.Errorf("AverageTime = %v, want %v", got.AverageTime, tt.want.AverageTime)
			}
			if got.TotalTime != tt.want.TotalTime {
				t.Errorf("TotalTime = %v, want %v", got.TotalTime, tt.want.TotalTime)
			}
			if got.FastestTime != tt.want.FastestTime {
				t.Errorf("FastestTime = %v, want %v", got.FastestTime, tt.want.FastestTime)
			}
			if got.SlowestTime != tt.want.SlowestTime {
				t.Errorf("SlowestTime = %v, want %v", got.SlowestTime, tt.want.SlowestTime)
			}
			if got.FastestEndpoint != tt.want.FastestEndpoint {
				t.Errorf("FastestEndpoint = %v, want %v", got.FastestEndpoint, tt.want.FastestEndpoint)
			}
			if got.SlowestEndpoint != tt.want.SlowestEndpoint {
				t.Errorf("SlowestEndpoint = %v, want %v", got.SlowestEndpoint, tt.want.SlowestEndpoint)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{
			name:     "zero duration",
			duration: 0,
			want:     "0ms",
		},
		{
			name:     "microseconds",
			duration: 500 * time.Microsecond,
			want:     "500µs",
		},
		{
			name:     "milliseconds",
			duration: 150 * time.Millisecond,
			want:     "150ms",
		},
		{
			name:     "seconds",
			duration: 2500 * time.Millisecond,
			want:     "2.50s",
		},
		{
			name:     "exact second",
			duration: 1 * time.Second,
			want:     "1.00s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.duration)
			if got != tt.want {
				t.Errorf("formatDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatStats(t *testing.T) {
	// Test that FormatStats doesn't panic and returns a non-empty string
	stats := TestStats{
		Total:           10,
		Passed:          8,
		Failed:          2,
		AverageTime:     100 * time.Millisecond,
		TotalTime:       1000 * time.Millisecond,
		FastestTime:     50 * time.Millisecond,
		SlowestTime:     200 * time.Millisecond,
		FastestEndpoint: "GET /fast",
		SlowestEndpoint: "POST /slow",
	}

	result := FormatStats(stats)

	if result == "" {
		t.Error("FormatStats() returned empty string")
	}

	// Basic sanity check - should contain some key text
	if len(result) < 50 {
		t.Errorf("FormatStats() result seems too short: %d characters", len(result))
	}
}

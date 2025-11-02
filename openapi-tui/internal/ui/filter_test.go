package ui

import (
	"testing"
	"time"

	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

func TestFilterResults(t *testing.T) {
	// Sample test results
	results := []models.TestResult{
		{Method: "GET", Endpoint: "/users", Status: "200", Message: "OK", Duration: 100 * time.Millisecond},
		{Method: "POST", Endpoint: "/users", Status: "201", Message: "Created", Duration: 150 * time.Millisecond},
		{Method: "GET", Endpoint: "/posts", Status: "200", Message: "OK", Duration: 80 * time.Millisecond},
		{Method: "DELETE", Endpoint: "/users/123", Status: "404", Message: "Not Found", Duration: 50 * time.Millisecond},
		{Method: "GET", Endpoint: "/api/v1/products", Status: "500", Message: "Internal Server Error", Duration: 200 * time.Millisecond},
		{Method: "PUT", Endpoint: "/posts/1", Status: "ERR", Message: "Connection refused", Duration: 0},
	}

	tests := []struct {
		name          string
		query         string
		expectedCount int
		description   string
	}{
		{
			name:          "empty query returns all",
			query:         "",
			expectedCount: 6,
			description:   "No filter should return all results",
		},
		{
			name:          "filter by status 200",
			query:         "200",
			expectedCount: 2,
			description:   "Should return only 200 status results",
		},
		{
			name:          "filter by status 404",
			query:         "404",
			expectedCount: 1,
			description:   "Should return only 404 status results",
		},
		{
			name:          "filter by method GET",
			query:         "GET",
			expectedCount: 3,
			description:   "Should return all GET requests",
		},
		{
			name:          "filter by method POST",
			query:         "post",
			expectedCount: 3,
			description:   "Should match POST method AND /posts endpoints (flexible matching)",
		},
		{
			name:          "filter by endpoint users",
			query:         "users",
			expectedCount: 3,
			description:   "Should match /users and /users/123",
		},
		{
			name:          "filter by endpoint products",
			query:         "products",
			expectedCount: 1,
			description:   "Should match /api/v1/products",
		},
		{
			name:          "filter by endpoint path segment",
			query:         "/api/",
			expectedCount: 1,
			description:   "Should match partial paths",
		},
		{
			name:          "filter by message OK",
			query:         "OK",
			expectedCount: 2,
			description:   "Should match 'OK' in messages",
		},
		{
			name:          "filter by message error",
			query:         "error",
			expectedCount: 1,
			description:   "Should match 'Internal Server Error' in message",
		},
		{
			name:          "special keyword pass",
			query:         "pass",
			expectedCount: 3,
			description:   "Should match all 2xx status codes",
		},
		{
			name:          "special keyword passed",
			query:         "passed",
			expectedCount: 3,
			description:   "Should match all 2xx status codes",
		},
		{
			name:          "special keyword success",
			query:         "success",
			expectedCount: 3,
			description:   "Should match all 2xx status codes",
		},
		{
			name:          "special keyword fail",
			query:         "fail",
			expectedCount: 3,
			description:   "Should match non-2xx status codes",
		},
		{
			name:          "special keyword error",
			query:         "err",
			expectedCount: 3,
			description:   "Should match non-2xx status codes and ERR",
		},
		{
			name:          "no matches",
			query:         "xyz123",
			expectedCount: 0,
			description:   "Should return empty when no matches",
		},
		{
			name:          "whitespace in query",
			query:         "  GET  ",
			expectedCount: 3,
			description:   "Should trim whitespace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := FilterResults(results, tt.query)
			if len(filtered) != tt.expectedCount {
				t.Errorf("%s: got %d results, want %d. Query: '%s'",
					tt.description, len(filtered), tt.expectedCount, tt.query)
			}
		})
	}
}

func TestFilterResultsEdgeCases(t *testing.T) {
	t.Run("empty results slice", func(t *testing.T) {
		results := []models.TestResult{}
		filtered := FilterResults(results, "GET")
		if len(filtered) != 0 {
			t.Errorf("Expected 0 results for empty slice, got %d", len(filtered))
		}
	})

	t.Run("nil results slice", func(t *testing.T) {
		var results []models.TestResult
		filtered := FilterResults(results, "GET")
		if filtered != nil && len(filtered) != 0 {
			t.Errorf("Expected nil or empty for nil slice, got %d results", len(filtered))
		}
	})

	t.Run("case insensitive matching", func(t *testing.T) {
		results := []models.TestResult{
			{Method: "GET", Endpoint: "/Users", Status: "200", Message: "OK"},
		}
		
		queries := []string{"users", "USERS", "Users", "uSeRs"}
		for _, query := range queries {
			filtered := FilterResults(results, query)
			if len(filtered) != 1 {
				t.Errorf("Case insensitive match failed for query '%s'", query)
			}
		}
	})
}

func TestMatchesFilter(t *testing.T) {
	tests := []struct {
		name     string
		result   models.TestResult
		query    string
		expected bool
	}{
		{
			name:     "matches status",
			result:   models.TestResult{Status: "200"},
			query:    "200",
			expected: true,
		},
		{
			name:     "matches method",
			result:   models.TestResult{Method: "GET"},
			query:    "get",
			expected: true,
		},
		{
			name:     "matches endpoint",
			result:   models.TestResult{Endpoint: "/users"},
			query:    "users",
			expected: true,
		},
		{
			name:     "matches message",
			result:   models.TestResult{Message: "Not Found"},
			query:    "not found",
			expected: true,
		},
		{
			name:     "matches pass keyword with 2xx",
			result:   models.TestResult{Status: "201"},
			query:    "pass",
			expected: true,
		},
		{
			name:     "matches fail keyword with 4xx",
			result:   models.TestResult{Status: "404"},
			query:    "fail",
			expected: true,
		},
		{
			name:     "does not match unrelated query",
			result:   models.TestResult{Method: "GET", Endpoint: "/users", Status: "200"},
			query:    "posts",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesFilter(tt.result, tt.query)
			if got != tt.expected {
				t.Errorf("matchesFilter() = %v, want %v for query '%s'", got, tt.expected, tt.query)
			}
		})
	}
}

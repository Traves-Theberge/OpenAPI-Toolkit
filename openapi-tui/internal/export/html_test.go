package export

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

func TestExportResultsToHTML(t *testing.T) {
	tests := []struct {
		name     string
		results  []models.TestResult
		specPath string
		baseURL  string
		wantErr  bool
	}{
		{
			name: "export with mixed results",
			results: []models.TestResult{
				{
					Method:   "GET",
					Endpoint: "/users",
					Status:   "200",
					Message:  "OK",
					Duration: 100 * time.Millisecond,
				},
				{
					Method:   "POST",
					Endpoint: "/users",
					Status:   "201",
					Message:  "Created",
					Duration: 150 * time.Millisecond,
				},
				{
					Method:   "GET",
					Endpoint: "/users/999",
					Status:   "404",
					Message:  "Not Found",
					Duration: 50 * time.Millisecond,
				},
			},
			specPath: "/path/to/spec.yaml",
			baseURL:  "https://api.example.com",
			wantErr:  false,
		},
		{
			name: "export with all passing results",
			results: []models.TestResult{
				{
					Method:   "GET",
					Endpoint: "/health",
					Status:   "200",
					Message:  "OK",
					Duration: 10 * time.Millisecond,
				},
				{
					Method:   "GET",
					Endpoint: "/version",
					Status:   "200",
					Message:  "OK",
					Duration: 15 * time.Millisecond,
				},
			},
			specPath: "spec.yaml",
			baseURL:  "http://localhost:8080",
			wantErr:  false,
		},
		{
			name: "export with all failing results",
			results: []models.TestResult{
				{
					Method:   "GET",
					Endpoint: "/api/fail1",
					Status:   "500",
					Message:  "Internal Server Error",
					Duration: 200 * time.Millisecond,
				},
				{
					Method:   "POST",
					Endpoint: "/api/fail2",
					Status:   "503",
					Message:  "Service Unavailable",
					Duration: 100 * time.Millisecond,
				},
			},
			specPath: "test-spec.yaml",
			baseURL:  "https://test.api.com",
			wantErr:  false,
		},
		{
			name:     "export with empty results",
			results:  []models.TestResult{},
			specPath: "empty-spec.yaml",
			baseURL:  "http://localhost:3000",
			wantErr:  false,
		},
		{
			name: "export with verbose log data",
			results: []models.TestResult{
				{
					Method:   "POST",
					Endpoint: "/api/users",
					Status:   "201",
					Message:  "Created",
					Duration: 125 * time.Millisecond,
					LogEntry: &models.LogEntry{
						RequestURL:      "https://api.example.com/api/users",
						RequestHeaders:  map[string]string{"Content-Type": "application/json"},
						RequestBody:     `{"name":"John Doe","email":"john@example.com"}`,
						ResponseHeaders: map[string]string{"Content-Type": "application/json"},
						ResponseBody:    `{"id":123,"name":"John Doe","email":"john@example.com"}`,
						Duration:        125 * time.Millisecond,
						Timestamp:       time.Now(),
					},
				},
			},
			specPath: "api-spec.yaml",
			baseURL:  "https://api.example.com",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Export to HTML
			filename, err := ExportResultsToHTML(tt.results, tt.specPath, tt.baseURL)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ExportResultsToHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return // Expected error, test passed
			}

			// Verify file was created
			if filename == "" {
				t.Error("ExportResultsToHTML() returned empty filename")
				return
			}

			// Clean up
			defer func() {
				if err := os.Remove(filename); err != nil {
					t.Logf("Warning: failed to remove test file %s: %v", filename, err)
				}
			}()

			// Read and validate file content
			content, err := os.ReadFile(filename)
			if err != nil {
				t.Fatalf("Failed to read exported HTML file: %v", err)
			}

			htmlContent := string(content)

			// Verify HTML structure
			if !strings.Contains(htmlContent, "<!DOCTYPE html>") {
				t.Error("HTML output missing DOCTYPE declaration")
			}
			if !strings.Contains(htmlContent, "<title>OpenAPI Test Results") {
				t.Error("HTML output missing title")
			}
			if !strings.Contains(htmlContent, "OpenAPI Test Results") {
				t.Error("HTML output missing header")
			}

			// Verify metadata
			if tt.specPath != "" && !strings.Contains(htmlContent, tt.specPath) {
				t.Errorf("HTML output missing spec path: %s", tt.specPath)
			}
			if tt.baseURL != "" && !strings.Contains(htmlContent, tt.baseURL) {
				t.Errorf("HTML output missing base URL: %s", tt.baseURL)
			}

			// Verify statistics section
			if !strings.Contains(htmlContent, "Total Tests") {
				t.Error("HTML output missing 'Total Tests' statistic")
			}
			if !strings.Contains(htmlContent, "Passed") {
				t.Error("HTML output missing 'Passed' statistic")
			}
			if !strings.Contains(htmlContent, "Failed") {
				t.Error("HTML output missing 'Failed' statistic")
			}
			if !strings.Contains(htmlContent, "Pass Rate") {
				t.Error("HTML output missing 'Pass Rate' statistic")
			}

			// Verify results table
			if !strings.Contains(htmlContent, "<table") {
				t.Error("HTML output missing results table")
			}
			if !strings.Contains(htmlContent, "<thead>") {
				t.Error("HTML output missing table header")
			}
			if !strings.Contains(htmlContent, "<tbody>") {
				t.Error("HTML output missing table body")
			}

			// Verify each result is present
			for _, result := range tt.results {
				if !strings.Contains(htmlContent, result.Method) {
					t.Errorf("HTML output missing method: %s", result.Method)
				}
				if !strings.Contains(htmlContent, result.Endpoint) {
					t.Errorf("HTML output missing endpoint: %s", result.Endpoint)
				}
				if !strings.Contains(htmlContent, result.Status) {
					t.Errorf("HTML output missing status: %s", result.Status)
				}
				if !strings.Contains(htmlContent, result.Message) {
					t.Errorf("HTML output missing message: %s", result.Message)
				}
			}

			// Verify styling
			if !strings.Contains(htmlContent, "<style>") {
				t.Error("HTML output missing CSS styling")
			}
			if !strings.Contains(htmlContent, ".results-table") {
				t.Error("HTML output missing results table CSS")
			}
			if !strings.Contains(htmlContent, ".stat-box") {
				t.Error("HTML output missing stat box CSS")
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
			name:     "zero duration",
			duration: 0,
			want:     "0µs",
		},
		{
			name:     "1 microsecond",
			duration: 1 * time.Microsecond,
			want:     "1µs",
		},
		{
			name:     "999 microseconds",
			duration: 999 * time.Microsecond,
			want:     "999µs",
		},
		{
			name:     "1 millisecond",
			duration: 1 * time.Millisecond,
			want:     "1ms",
		},
		{
			name:     "999 milliseconds",
			duration: 999 * time.Millisecond,
			want:     "999ms",
		},
		{
			name:     "1 second",
			duration: 1 * time.Second,
			want:     "1.00s",
		},
		{
			name:     "large duration",
			duration: 45 * time.Second,
			want:     "45.00s",
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

func TestHTMLTemplateData(t *testing.T) {
	// Test that HTMLTemplateData structure is correctly populated
	results := []models.TestResult{
		{
			Method:   "GET",
			Endpoint: "/test",
			Status:   "200",
			Message:  "OK",
			Duration: 100 * time.Millisecond,
		},
		{
			Method:   "POST",
			Endpoint: "/test",
			Status:   "500",
			Message:  "Error",
			Duration: 200 * time.Millisecond,
		},
	}

	filename, err := ExportResultsToHTML(results, "test.yaml", "http://localhost")
	if err != nil {
		t.Fatalf("ExportResultsToHTML() failed: %v", err)
	}
	defer os.Remove(filename)

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	htmlContent := string(content)

	// Verify pass rate calculation (1 passed out of 2 = 50%)
	if !strings.Contains(htmlContent, "50.0%") {
		t.Error("HTML should contain 50.0% pass rate")
	}

	// Verify total time is displayed
	if !strings.Contains(htmlContent, "Total Time") {
		t.Error("HTML should contain total time")
	}

	// Verify average time is displayed
	if !strings.Contains(htmlContent, "Average Time") {
		t.Error("HTML should contain average time")
	}

	// Verify row classes are applied
	if !strings.Contains(htmlContent, `class="success"`) {
		t.Error("HTML should contain success row class")
	}
	if !strings.Contains(htmlContent, `class="failure"`) {
		t.Error("HTML should contain failure row class")
	}
}

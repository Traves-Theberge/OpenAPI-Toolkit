package export

import (
	"encoding/xml"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

func TestExportResultsToJUnit(t *testing.T) {
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
					Status:   "404",
					Message:  "Not Found",
					Duration: 50 * time.Millisecond,
				},
				{
					Method:   "DELETE",
					Endpoint: "/users/123",
					Status:   "ERR",
					Message:  "Connection refused",
					Duration: 0,
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
					Status:   "201",
					Message:  "Created",
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
						RequestHeaders:  map[string]string{"Content-Type": "application/json", "Authorization": "Bearer token123"},
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
			// Export to JUnit XML
			filename, err := ExportResultsToJUnit(tt.results, tt.specPath, tt.baseURL)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExportResultsToJUnit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return // Expected error, test passed
			}

			// Verify file was created
			if filename == "" {
				t.Error("ExportResultsToJUnit() returned empty filename")
				return
			}

			// Clean up
			defer func() {
				if err := os.Remove(filename); err != nil {
					t.Logf("Warning: failed to remove test file %s: %v", filename, err)
				}
			}()

			// Read and parse XML file
			content, err := os.ReadFile(filename)
			if err != nil {
				t.Fatalf("Failed to read exported JUnit XML file: %v", err)
			}

			// Verify XML is valid by unmarshaling
			var suites JUnitTestSuites
			if err := xml.Unmarshal(content, &suites); err != nil {
				t.Fatalf("Failed to parse JUnit XML: %v", err)
			}

			xmlContent := string(content)

			// Verify XML declaration
			if !strings.HasPrefix(xmlContent, "<?xml") {
				t.Error("JUnit XML missing XML declaration")
			}

			// Verify root element
			if !strings.Contains(xmlContent, "<testsuites>") {
				t.Error("JUnit XML missing testsuites root element")
			}

			// Verify test suite exists
			if len(suites.Suites) == 0 {
				t.Fatal("JUnit XML should contain at least one test suite")
			}

			suite := suites.Suites[0]

			// Verify test suite name
			if suite.Name != "OpenAPI Tests" {
				t.Errorf("Test suite name = %s, want 'OpenAPI Tests'", suite.Name)
			}

			// Verify test count
			if suite.Tests != len(tt.results) {
				t.Errorf("Test suite tests count = %d, want %d", suite.Tests, len(tt.results))
			}

			// Verify properties
			foundSpecPath := false
			foundBaseURL := false
			for _, prop := range suite.Properties {
				if prop.Name == "spec_path" && prop.Value == tt.specPath {
					foundSpecPath = true
				}
				if prop.Name == "base_url" && prop.Value == tt.baseURL {
					foundBaseURL = true
				}
			}
			if tt.specPath != "" && !foundSpecPath {
				t.Error("JUnit XML missing spec_path property")
			}
			if tt.baseURL != "" && !foundBaseURL {
				t.Error("JUnit XML missing base_url property")
			}

			// Verify test cases
			if len(suite.TestCases) != len(tt.results) {
				t.Errorf("Test case count = %d, want %d", len(suite.TestCases), len(tt.results))
			}

			// Verify each test case
			for i, result := range tt.results {
				if i >= len(suite.TestCases) {
					break
				}
				testCase := suite.TestCases[i]

				// Verify test case name contains method and endpoint
				if !strings.Contains(testCase.Name, result.Method) {
					t.Errorf("Test case %d name missing method: %s", i, result.Method)
				}
				if !strings.Contains(testCase.Name, result.Endpoint) {
					t.Errorf("Test case %d name missing endpoint: %s", i, result.Endpoint)
				}

				// Verify failure/error is set correctly
				if result.Status == "ERR" {
					if testCase.Error == nil {
						t.Errorf("Test case %d should have error for ERR status", i)
					}
				} else if !strings.HasPrefix(result.Status, "2") {
					if testCase.Failure == nil {
						t.Errorf("Test case %d should have failure for non-2xx status", i)
					}
				} else {
					// Passing test should have no failure or error
					if testCase.Failure != nil || testCase.Error != nil {
						t.Errorf("Test case %d should not have failure/error for 2xx status", i)
					}
				}

				// Verify system-out for verbose logs
				if result.LogEntry != nil {
					if testCase.SystemOut == "" {
						t.Errorf("Test case %d should have system-out for verbose log", i)
					}
					if !strings.Contains(testCase.SystemOut, result.LogEntry.RequestURL) {
						t.Errorf("Test case %d system-out missing request URL", i)
					}
				}
			}

			// Verify failure and error counts
			expectedFailures := 0
			expectedErrors := 0
			for _, r := range tt.results {
				if r.Status == "ERR" {
					expectedErrors++
				} else if !strings.HasPrefix(r.Status, "2") {
					expectedFailures++
				}
			}

			if suite.Failures != expectedFailures {
				t.Errorf("Test suite failures = %d, want %d", suite.Failures, expectedFailures)
			}
			if suite.Errors != expectedErrors {
				t.Errorf("Test suite errors = %d, want %d", suite.Errors, expectedErrors)
			}
		})
	}
}

func TestSanitizeClassName(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "https URL",
			url:  "https://api.example.com",
			want: "api.example.com",
		},
		{
			name: "http URL",
			url:  "http://localhost:8080",
			want: "localhost.8080",
		},
		{
			name: "URL with path",
			url:  "https://api.example.com/v1/users",
			want: "api.example.com.v1.users",
		},
		{
			name: "URL with port and dashes",
			url:  "http://my-api.test.com:3000",
			want: "my_api.test.com.3000",
		},
		{
			name: "empty URL",
			url:  "",
			want: "openapi.tests",
		},
		{
			name: "URL with trailing slash",
			url:  "https://api.example.com/",
			want: "api.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeClassName(tt.url)
			if got != tt.want {
				t.Errorf("sanitizeClassName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDurationSeconds(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{
			name:     "milliseconds",
			duration: 125 * time.Millisecond,
			want:     "0.125",
		},
		{
			name:     "seconds",
			duration: 2500 * time.Millisecond,
			want:     "2.500",
		},
		{
			name:     "zero duration",
			duration: 0,
			want:     "0.000",
		},
		{
			name:     "microseconds",
			duration: 500 * time.Microsecond,
			want:     "0.001",
		},
		{
			name:     "1 second",
			duration: 1 * time.Second,
			want:     "1.000",
		},
		{
			name:     "large duration",
			duration: 45 * time.Second,
			want:     "45.000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDurationSeconds(tt.duration)
			if got != tt.want {
				t.Errorf("formatDurationSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJUnitXMLStructure(t *testing.T) {
	// Test that JUnit XML structure is correctly formed
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

	filename, err := ExportResultsToJUnit(results, "test.yaml", "http://localhost")
	if err != nil {
		t.Fatalf("ExportResultsToJUnit() failed: %v", err)
	}
	defer os.Remove(filename)

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	xmlContent := string(content)

	// Verify essential JUnit XML elements
	requiredElements := []string{
		"<testsuites>",
		"<testsuite",
		"tests=",
		"failures=",
		"errors=",
		"time=",
		"timestamp=",
		"<properties>",
		"<property",
		"<testcase",
		"name=",
		"classname=",
	}

	for _, elem := range requiredElements {
		if !strings.Contains(xmlContent, elem) {
			t.Errorf("JUnit XML missing required element: %s", elem)
		}
	}

	// Verify failure element for non-2xx status
	if !strings.Contains(xmlContent, "<failure") {
		t.Error("JUnit XML should contain failure element for 500 status")
	}

	// Parse and validate XML structure
	var suites JUnitTestSuites
	if err := xml.Unmarshal(content, &suites); err != nil {
		t.Fatalf("Failed to unmarshal JUnit XML: %v", err)
	}

	// Verify structure
	if len(suites.Suites) != 1 {
		t.Errorf("Expected 1 test suite, got %d", len(suites.Suites))
	}

	suite := suites.Suites[0]
	if suite.Tests != 2 {
		t.Errorf("Expected 2 tests, got %d", suite.Tests)
	}
	if suite.Failures != 1 {
		t.Errorf("Expected 1 failure, got %d", suite.Failures)
	}
	if suite.Errors != 0 {
		t.Errorf("Expected 0 errors, got %d", suite.Errors)
	}
}

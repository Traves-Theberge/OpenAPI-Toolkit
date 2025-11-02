package export

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

// TestExportResults tests successful export with all fields
func TestExportResults(t *testing.T) {
	results := []models.TestResult{
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
	}

	filename, err := ExportResults(results, "/path/to/spec.yaml")
	if err != nil {
		t.Fatalf("ExportResults failed: %v", err)
	}

	// Cleanup
	defer os.Remove(filename)

	// Verify filename format
	if !strings.HasPrefix(filename, "openapi-test-results_") {
		t.Errorf("Expected filename to start with 'openapi-test-results_', got: %s", filename)
	}
	if !strings.HasSuffix(filename, ".json") {
		t.Errorf("Expected filename to end with '.json', got: %s", filename)
	}

	// Verify file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("Expected file to exist: %s", filename)
	}

	// Read and verify file content
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read exported file: %v", err)
	}

	var exportData models.ExportData
	if err := json.Unmarshal(data, &exportData); err != nil {
		t.Fatalf("Failed to unmarshal exported data: %v", err)
	}

	// Verify metadata
	if exportData.SpecPath != "/path/to/spec.yaml" {
		t.Errorf("Expected SpecPath '/path/to/spec.yaml', got: %s", exportData.SpecPath)
	}
	if exportData.TotalTests != 2 {
		t.Errorf("Expected TotalTests 2, got: %d", exportData.TotalTests)
	}
	if exportData.Passed != 2 {
		t.Errorf("Expected Passed 2, got: %d", exportData.Passed)
	}
	if exportData.Failed != 0 {
		t.Errorf("Expected Failed 0, got: %d", exportData.Failed)
	}

	// Verify results
	if len(exportData.Results) != 2 {
		t.Fatalf("Expected 2 results, got: %d", len(exportData.Results))
	}

	// Verify first result
	if exportData.Results[0].Method != "GET" {
		t.Errorf("Expected Method 'GET', got: %s", exportData.Results[0].Method)
	}
	if exportData.Results[0].Endpoint != "/users" {
		t.Errorf("Expected Endpoint '/users', got: %s", exportData.Results[0].Endpoint)
	}
	if exportData.Results[0].Status != "200" {
		t.Errorf("Expected Status '200', got: %s", exportData.Results[0].Status)
	}
}

// TestExportResults_EmptyResults tests export with no results
func TestExportResults_EmptyResults(t *testing.T) {
	results := []models.TestResult{}

	filename, err := ExportResults(results, "/path/to/spec.yaml")
	if err != nil {
		t.Fatalf("ExportResults failed: %v", err)
	}

	defer os.Remove(filename)

	// Read and verify
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read exported file: %v", err)
	}

	var exportData models.ExportData
	if err := json.Unmarshal(data, &exportData); err != nil {
		t.Fatalf("Failed to unmarshal exported data: %v", err)
	}

	if exportData.TotalTests != 0 {
		t.Errorf("Expected TotalTests 0, got: %d", exportData.TotalTests)
	}
	if exportData.Passed != 0 {
		t.Errorf("Expected Passed 0, got: %d", exportData.Passed)
	}
	if exportData.Failed != 0 {
		t.Errorf("Expected Failed 0, got: %d", exportData.Failed)
	}
	if len(exportData.Results) != 0 {
		t.Errorf("Expected 0 results, got: %d", len(exportData.Results))
	}
}

// TestExportResults_FailedTests tests export with failed tests
func TestExportResults_FailedTests(t *testing.T) {
	results := []models.TestResult{
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
			Status:   "ERR",
			Message:  "Connection failed",
			Duration: 50 * time.Millisecond,
		},
		{
			Method:   "DELETE",
			Endpoint: "/users/1",
			Status:   "404",
			Message:  "validation failed",
			Duration: 75 * time.Millisecond,
		},
	}

	filename, err := ExportResults(results, "/spec.yaml")
	if err != nil {
		t.Fatalf("ExportResults failed: %v", err)
	}

	defer os.Remove(filename)

	// Read and verify
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read exported file: %v", err)
	}

	var exportData models.ExportData
	if err := json.Unmarshal(data, &exportData); err != nil {
		t.Fatalf("Failed to unmarshal exported data: %v", err)
	}

	// Verify statistics
	if exportData.TotalTests != 3 {
		t.Errorf("Expected TotalTests 3, got: %d", exportData.TotalTests)
	}
	if exportData.Passed != 1 {
		t.Errorf("Expected Passed 1, got: %d", exportData.Passed)
	}
	if exportData.Failed != 2 {
		t.Errorf("Expected Failed 2, got: %d", exportData.Failed)
	}
}

// TestExportResults_MixedStatuses tests export with various status codes
func TestExportResults_MixedStatuses(t *testing.T) {
	results := []models.TestResult{
		{Method: "GET", Endpoint: "/api/v1", Status: "200", Message: "OK"},
		{Method: "POST", Endpoint: "/api/v2", Status: "201", Message: "Created"},
		{Method: "PUT", Endpoint: "/api/v3", Status: "404", Message: "Not Found"},
		{Method: "DELETE", Endpoint: "/api/v4", Status: "500", Message: "Server Error"},
		{Method: "PATCH", Endpoint: "/api/v5", Status: "ERR", Message: "Network timeout"},
	}

	filename, err := ExportResults(results, "/spec.yaml")
	if err != nil {
		t.Fatalf("ExportResults failed: %v", err)
	}

	defer os.Remove(filename)

	// Read and verify
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read exported file: %v", err)
	}

	var exportData models.ExportData
	if err := json.Unmarshal(data, &exportData); err != nil {
		t.Fatalf("Failed to unmarshal exported data: %v", err)
	}

	if exportData.TotalTests != 5 {
		t.Errorf("Expected TotalTests 5, got: %d", exportData.TotalTests)
	}
	// Logic: passed if Status != "ERR" AND Message doesn't contain "failed"
	// Status "200", "201", "404", "500" without "failed" in message = 4 passed
	// Status "ERR" = 1 failed
	if exportData.Passed != 4 {
		t.Errorf("Expected Passed 4, got: %d", exportData.Passed)
	}
	if exportData.Failed != 1 {
		t.Errorf("Expected Failed 1, got: %d", exportData.Failed)
	}
}

// TestExportResultsToFile tests export with custom filename
func TestExportResultsToFile(t *testing.T) {
	results := []models.TestResult{
		{
			Method:   "GET",
			Endpoint: "/test",
			Status:   "200",
			Message:  "OK",
			Duration: 100 * time.Millisecond,
		},
	}

	customFilename := "custom_test_results.json"
	err := ExportResultsToFile(results, "/spec.yaml", customFilename)
	if err != nil {
		t.Fatalf("ExportResultsToFile failed: %v", err)
	}

	defer os.Remove(customFilename)

	// Verify file exists
	if _, err := os.Stat(customFilename); os.IsNotExist(err) {
		t.Errorf("Expected file to exist: %s", customFilename)
	}

	// Read and verify
	data, err := os.ReadFile(customFilename)
	if err != nil {
		t.Fatalf("Failed to read exported file: %v", err)
	}

	var exportData models.ExportData
	if err := json.Unmarshal(data, &exportData); err != nil {
		t.Fatalf("Failed to unmarshal exported data: %v", err)
	}

	if exportData.TotalTests != 1 {
		t.Errorf("Expected TotalTests 1, got: %d", exportData.TotalTests)
	}
}

// TestExportResultsToFile_InvalidPath tests export to invalid path
func TestExportResultsToFile_InvalidPath(t *testing.T) {
	results := []models.TestResult{
		{Method: "GET", Endpoint: "/test", Status: "200", Message: "OK"},
	}

	// Try to write to an invalid directory
	err := ExportResultsToFile(results, "/spec.yaml", "/nonexistent/directory/file.json")
	if err == nil {
		t.Fatal("Expected error when writing to invalid path, got nil")
	}

	if !strings.Contains(err.Error(), "failed to write file") {
		t.Errorf("Expected error to contain 'failed to write file', got: %s", err.Error())
	}
}

// TestExportResultsToFile_EmptyFilename tests export with empty filename
func TestExportResultsToFile_EmptyFilename(t *testing.T) {
	results := []models.TestResult{
		{Method: "GET", Endpoint: "/test", Status: "200", Message: "OK"},
	}

	err := ExportResultsToFile(results, "/spec.yaml", "")
	if err == nil {
		t.Fatal("Expected error when using empty filename, got nil")
	}
}

// TestFormatExportSummary tests summary formatting
func TestFormatExportSummary(t *testing.T) {
	summary := FormatExportSummary("test_results.json", 5)

	// Verify summary contains expected elements
	if !strings.Contains(summary, "test_results.json") {
		t.Errorf("Expected summary to contain filename, got: %s", summary)
	}
	if !strings.Contains(summary, "5 result(s)") {
		t.Errorf("Expected summary to contain result count, got: %s", summary)
	}
	if !strings.Contains(summary, "✓") {
		t.Errorf("Expected summary to contain checkmark, got: %s", summary)
	}
	if !strings.Contains(summary, "Timestamp:") {
		t.Errorf("Expected summary to contain timestamp, got: %s", summary)
	}
}

// TestFormatExportSummary_ZeroResults tests summary with zero results
func TestFormatExportSummary_ZeroResults(t *testing.T) {
	summary := FormatExportSummary("empty_results.json", 0)

	if !strings.Contains(summary, "0 result(s)") {
		t.Errorf("Expected summary to contain '0 result(s)', got: %s", summary)
	}
}

// TestFormatExportSummary_LargeCount tests summary with large result count
func TestFormatExportSummary_LargeCount(t *testing.T) {
	summary := FormatExportSummary("large_results.json", 1000)

	if !strings.Contains(summary, "1000 result(s)") {
		t.Errorf("Expected summary to contain '1000 result(s)', got: %s", summary)
	}
}

// TestExportResults_DurationFormatting tests that durations are properly converted to strings
func TestExportResults_DurationFormatting(t *testing.T) {
	results := []models.TestResult{
		{
			Method:   "GET",
			Endpoint: "/test",
			Status:   "200",
			Message:  "OK",
			Duration: 123 * time.Millisecond,
		},
	}

	filename, err := ExportResults(results, "/spec.yaml")
	if err != nil {
		t.Fatalf("ExportResults failed: %v", err)
	}

	defer os.Remove(filename)

	// Read and verify duration is a string
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read exported file: %v", err)
	}

	var exportData models.ExportData
	if err := json.Unmarshal(data, &exportData); err != nil {
		t.Fatalf("Failed to unmarshal exported data: %v", err)
	}

	// Verify duration is formatted as string
	if exportData.Results[0].Duration == "" {
		t.Error("Expected duration to be non-empty string")
	}
	if !strings.Contains(exportData.Results[0].Duration, "ms") {
		t.Errorf("Expected duration to contain 'ms', got: %s", exportData.Results[0].Duration)
	}
}

// TestExportResults_TimestampFormat tests timestamp format
func TestExportResults_TimestampFormat(t *testing.T) {
	results := []models.TestResult{
		{Method: "GET", Endpoint: "/test", Status: "200", Message: "OK"},
	}

	filename, err := ExportResults(results, "/spec.yaml")
	if err != nil {
		t.Fatalf("ExportResults failed: %v", err)
	}

	defer os.Remove(filename)

	// Read and verify timestamp format
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read exported file: %v", err)
	}

	var exportData models.ExportData
	if err := json.Unmarshal(data, &exportData); err != nil {
		t.Fatalf("Failed to unmarshal exported data: %v", err)
	}

	// Verify timestamp is in RFC3339 format
	_, err = time.Parse(time.RFC3339, exportData.Timestamp)
	if err != nil {
		t.Errorf("Expected timestamp in RFC3339 format, got parse error: %v", err)
	}
}

// TestExportResults_SpecPathPreserved tests that spec path is preserved correctly
func TestExportResults_SpecPathPreserved(t *testing.T) {
	testCases := []struct {
		name     string
		specPath string
	}{
		{"Relative path", "./spec.yaml"},
		{"Absolute path", "/home/user/spec.yaml"},
		{"With spaces", "/path with spaces/spec.yaml"},
		{"Empty path", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results := []models.TestResult{
				{Method: "GET", Endpoint: "/test", Status: "200", Message: "OK"},
			}

			filename, err := ExportResults(results, tc.specPath)
			if err != nil {
				t.Fatalf("ExportResults failed: %v", err)
			}

			defer os.Remove(filename)

			// Read and verify
			data, err := os.ReadFile(filename)
			if err != nil {
				t.Fatalf("Failed to read exported file: %v", err)
			}

			var exportData models.ExportData
			if err := json.Unmarshal(data, &exportData); err != nil {
				t.Fatalf("Failed to unmarshal exported data: %v", err)
			}

			if exportData.SpecPath != tc.specPath {
				t.Errorf("Expected SpecPath '%s', got: '%s'", tc.specPath, exportData.SpecPath)
			}
		})
	}
}

// TestExportResults_JSONFormat tests that output is valid JSON
func TestExportResults_JSONFormat(t *testing.T) {
	results := []models.TestResult{
		{Method: "GET", Endpoint: "/test", Status: "200", Message: "OK"},
	}

	filename, err := ExportResults(results, "/spec.yaml")
	if err != nil {
		t.Fatalf("ExportResults failed: %v", err)
	}

	defer os.Remove(filename)

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read exported file: %v", err)
	}

	// Verify it's valid JSON by unmarshaling to generic interface
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		t.Errorf("Exported file is not valid JSON: %v", err)
	}

	// Verify it's properly formatted (has indentation)
	if !strings.Contains(string(data), "\n") {
		t.Error("Expected JSON to be formatted with indentation")
	}
}

package export

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

// ExportResults exports test results to a JSON file
// Returns the filename and any error
func ExportResults(results []models.TestResult, specPath string) (string, error) {
	// Calculate statistics
	passed := 0
	failed := 0
	for _, r := range results {
		if r.Status != "ERR" && !strings.Contains(r.Message, "failed") {
			passed++
		} else {
			failed++
		}
	}

	// Create export data structure
	data := models.ExportData{
		Timestamp:  time.Now().Format(time.RFC3339),
		SpecPath:   specPath,
		BaseURL:    "",
		TotalTests: len(results),
		Passed:     passed,
		Failed:     failed,
		Results:    make([]models.ExportResult, len(results)),
	}

	// Convert test results to export format
	for i, r := range results {
		data.Results[i] = models.ExportResult{
			Method:     r.Method,
			Endpoint:   r.Endpoint,
			Status:     r.Status,
			Message:    r.Message,
			Duration:   r.Duration.String(), // Convert duration to string
			RetryCount: r.RetryCount,        // Include retry count
		}
	}

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("openapi-test-results_%s.json", timestamp)

	// Write to file
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filename, nil
}

// ExportResultsToFile exports results with a custom filename
func ExportResultsToFile(results []models.TestResult, specPath, filename string) error {
	// Calculate statistics
	passed := 0
	failed := 0
	for _, r := range results {
		if r.Status != "ERR" && !strings.Contains(r.Message, "failed") {
			passed++
		} else {
			failed++
		}
	}

	// Create export data structure
	data := models.ExportData{
		Timestamp:  time.Now().Format(time.RFC3339),
		SpecPath:   specPath,
		BaseURL:    "",
		TotalTests: len(results),
		Passed:     passed,
		Failed:     failed,
		Results:    make([]models.ExportResult, len(results)),
	}

	// Convert test results to export format
	for i, r := range results {
		data.Results[i] = models.ExportResult{
			Method:     r.Method,
			Endpoint:   r.Endpoint,
			Status:     r.Status,
			Message:    r.Message,
			Duration:   r.Duration.String(), // Convert duration to string
			RetryCount: r.RetryCount,        // Include retry count
		}
	}

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// FormatExportSummary creates a human-readable summary of export
func FormatExportSummary(filename string, resultCount int) string {
	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("✓ Exported %d result(s) to %s\n", resultCount, filename))
	summary.WriteString(fmt.Sprintf("  Timestamp: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	return summary.String()
}

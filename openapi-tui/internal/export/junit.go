package export

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

// JUnitTestSuites represents the root element of JUnit XML format
type JUnitTestSuites struct {
	XMLName xml.Name         `xml:"testsuites"`
	Suites  []JUnitTestSuite `xml:"testsuite"`
}

// JUnitTestSuite represents a single test suite
type JUnitTestSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Name       string          `xml:"name,attr"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Errors     int             `xml:"errors,attr"`
	Skipped    int             `xml:"skipped,attr"`
	Time       string          `xml:"time,attr"`
	Timestamp  string          `xml:"timestamp,attr"`
	Hostname   string          `xml:"hostname,attr,omitempty"`
	Properties []JUnitProperty `xml:"properties>property,omitempty"`
	TestCases  []JUnitTestCase `xml:"testcase"`
}

// JUnitProperty represents a property in the test suite
type JUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// JUnitTestCase represents a single test case
type JUnitTestCase struct {
	Name      string         `xml:"name,attr"`
	Classname string         `xml:"classname,attr"`
	Time      string         `xml:"time,attr"`
	Failure   *JUnitFailure  `xml:"failure,omitempty"`
	Error     *JUnitError    `xml:"error,omitempty"`
	SystemOut string         `xml:"system-out,omitempty"`
	SystemErr string         `xml:"system-err,omitempty"`
}

// JUnitFailure represents a test failure
type JUnitFailure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

// JUnitError represents a test error
type JUnitError struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

// ExportResultsToJUnit exports test results to JUnit XML format
// Returns the filename and any error
func ExportResultsToJUnit(results []models.TestResult, specPath, baseURL string) (string, error) {
	// Calculate statistics
	failures := 0
	errors := 0
	var totalDuration time.Duration

	for _, r := range results {
		totalDuration += r.Duration
		
		// ERR status is an error, non-2xx is a failure
		if r.Status == "ERR" {
			errors++
		} else if !strings.HasPrefix(r.Status, "2") {
			failures++
		}
	}

	// Build test cases
	testCases := make([]JUnitTestCase, len(results))
	for i, r := range results {
		// Create test case name and classname
		testName := fmt.Sprintf("%s %s", r.Method, r.Endpoint)
		className := sanitizeClassName(baseURL)

		testCase := JUnitTestCase{
			Name:      testName,
			Classname: className,
			Time:      formatDurationSeconds(r.Duration),
		}

		// Add failure or error if test didn't pass
		if r.Status == "ERR" {
			testCase.Error = &JUnitError{
				Message: r.Message,
				Type:    "Error",
				Content: fmt.Sprintf("Test failed with error: %s\nEndpoint: %s\nMethod: %s",
					r.Message, r.Endpoint, r.Method),
			}
		} else if !strings.HasPrefix(r.Status, "2") {
			testCase.Failure = &JUnitFailure{
				Message: fmt.Sprintf("HTTP %s: %s", r.Status, r.Message),
				Type:    "AssertionFailure",
				Content: fmt.Sprintf("Expected 2xx status code, got %s\nEndpoint: %s\nMethod: %s\nMessage: %s",
					r.Status, r.Endpoint, r.Method, r.Message),
			}
		}

		// Add verbose log data to system-out if available
		if r.LogEntry != nil {
			var sysOut strings.Builder
			sysOut.WriteString(fmt.Sprintf("Request: %s %s\n", r.Method, r.LogEntry.RequestURL))
			sysOut.WriteString(fmt.Sprintf("Duration: %s\n", r.Duration))
			if r.RetryCount > 0 {
				sysOut.WriteString(fmt.Sprintf("Retries: %d\n", r.RetryCount))
			}
			
			if len(r.LogEntry.RequestHeaders) > 0 {
				sysOut.WriteString("Request Headers:\n")
				for k, v := range r.LogEntry.RequestHeaders {
					sysOut.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
				}
			}
			
			if r.LogEntry.RequestBody != "" {
				sysOut.WriteString(fmt.Sprintf("Request Body:\n%s\n", r.LogEntry.RequestBody))
			}
			
			if len(r.LogEntry.ResponseHeaders) > 0 {
				sysOut.WriteString("Response Headers:\n")
				for k, v := range r.LogEntry.ResponseHeaders {
					sysOut.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
				}
			}
			
			if r.LogEntry.ResponseBody != "" {
				sysOut.WriteString(fmt.Sprintf("Response Body:\n%s\n", r.LogEntry.ResponseBody))
			}
			
			testCase.SystemOut = sysOut.String()
		}

		testCases[i] = testCase
	}

	// Create test suite
	suite := JUnitTestSuite{
		Name:      "OpenAPI Tests",
		Tests:     len(results),
		Failures:  failures,
		Errors:    errors,
		Skipped:   0,
		Time:      formatDurationSeconds(totalDuration),
		Timestamp: time.Now().Format(time.RFC3339),
		Properties: []JUnitProperty{
			{Name: "spec_path", Value: specPath},
			{Name: "base_url", Value: baseURL},
			{Name: "test_framework", Value: "openapi-tui"},
		},
		TestCases: testCases,
	}

	// Create root element with single suite
	suites := JUnitTestSuites{
		Suites: []JUnitTestSuite{suite},
	}

	// Marshal to XML with proper formatting
	xmlData, err := xml.MarshalIndent(suites, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JUnit XML: %w", err)
	}

	// Add XML declaration
	xmlOutput := []byte(xml.Header + string(xmlData))

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("openapi-test-results_%s.xml", timestamp)

	// Write to file
	if err := os.WriteFile(filename, xmlOutput, 0644); err != nil {
		return "", fmt.Errorf("failed to write JUnit XML file: %w", err)
	}

	return filename, nil
}

// sanitizeClassName converts a URL to a valid Java-style class name
func sanitizeClassName(url string) string {
	// Remove protocol
	className := strings.TrimPrefix(url, "https://")
	className = strings.TrimPrefix(className, "http://")
	
	// Replace special characters with dots
	className = strings.ReplaceAll(className, "/", ".")
	className = strings.ReplaceAll(className, ":", ".")
	className = strings.ReplaceAll(className, "-", "_")
	
	// Remove trailing dots
	className = strings.TrimSuffix(className, ".")
	
	// Default if empty
	if className == "" {
		className = "openapi.tests"
	}
	
	return className
}

// formatDurationSeconds formats a duration as seconds with decimal precision
func formatDurationSeconds(d time.Duration) string {
	return fmt.Sprintf("%.3f", d.Seconds())
}

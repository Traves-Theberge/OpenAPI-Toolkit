package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

// ExecuteCustomRequest executes a manually created API request
func ExecuteCustomRequest(method, endpoint string, headers map[string]string, body string, auth *models.AuthConfig, verbose bool) (models.TestResult, error) {
	startTime := time.Now()

	// Validate method
	method = strings.ToUpper(method)
	validMethods := map[string]bool{"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true, "HEAD": true, "OPTIONS": true}
	if !validMethods[method] {
		return models.TestResult{
			Method:   method,
			Endpoint: endpoint,
			Status:   "ERR",
			Message:  fmt.Sprintf("Invalid HTTP method: %s", method),
			Duration: time.Since(startTime),
		}, fmt.Errorf("invalid HTTP method: %s", method)
	}

	// Validate and parse body if present
	var bodyReader io.Reader
	if body != "" {
		// Try to parse as JSON to validate
		var jsonData interface{}
		if err := json.Unmarshal([]byte(body), &jsonData); err != nil {
			return models.TestResult{
				Method:   method,
				Endpoint: endpoint,
				Status:   "ERR",
				Message:  fmt.Sprintf("Invalid JSON body: %v", err),
				Duration: time.Since(startTime),
			}, fmt.Errorf("invalid JSON body: %w", err)
		}
		bodyReader = bytes.NewBufferString(body)
	}

	// Create HTTP request
	req, err := http.NewRequest(method, endpoint, bodyReader)
	if err != nil {
		return models.TestResult{
			Method:   method,
			Endpoint: endpoint,
			Status:   "ERR",
			Message:  fmt.Sprintf("Failed to create request: %v", err),
			Duration: time.Since(startTime),
		}, fmt.Errorf("failed to create request: %w", err)
	}

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set Content-Type if body is present and not already set
	if body != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Apply authentication
	if auth != nil {
		switch auth.AuthType {
		case "Bearer":
			req.Header.Set("Authorization", "Bearer "+auth.Token)
		case "API Key":
			if auth.APIKeyIn == "header" {
				req.Header.Set(auth.APIKeyName, auth.Token)
			} else if auth.APIKeyIn == "query" {
				q := req.URL.Query()
				q.Add(auth.APIKeyName, auth.Token)
				req.URL.RawQuery = q.Encode()
			}
		case "Basic":
			req.SetBasicAuth(auth.Username, auth.Password)
		}
	}

	// Execute request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	duration := time.Since(startTime)

	// Capture request details for verbose logging
	var logEntry *models.LogEntry
	if verbose {
		requestHeaders := make(map[string]string)
		for k, v := range req.Header {
			requestHeaders[k] = strings.Join(v, ", ")
		}

		logEntry = &models.LogEntry{
			RequestURL:     endpoint,
			RequestHeaders: requestHeaders,
			RequestBody:    body,
			Duration:       duration,
			Timestamp:      startTime,
		}
	}

	if err != nil {
		message := fmt.Sprintf("Request failed: %v", err)
		if logEntry != nil {
			logEntry.ResponseBody = message
		}
		return models.TestResult{
			Method:   method,
			Endpoint: endpoint,
			Status:   "ERR",
			Message:  message,
			Duration: duration,
			LogEntry: logEntry,
		}, err
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		responseBody = []byte(fmt.Sprintf("Failed to read response: %v", err))
	}

	// Capture response details for verbose logging
	if verbose && logEntry != nil {
		responseHeaders := make(map[string]string)
		for k, v := range resp.Header {
			responseHeaders[k] = strings.Join(v, ", ")
		}
		logEntry.ResponseHeaders = responseHeaders
		logEntry.ResponseBody = string(responseBody)
	}

	// Format status
	statusCode := resp.StatusCode
	statusText := http.StatusText(statusCode)
	statusStr := fmt.Sprintf("%d", statusCode)

	// Create message
	message := fmt.Sprintf("%d %s", statusCode, statusText)
	if len(responseBody) > 0 {
		// Try to pretty-print JSON
		var jsonData interface{}
		if json.Unmarshal(responseBody, &jsonData) == nil {
			if prettyJSON, err := json.MarshalIndent(jsonData, "", "  "); err == nil {
				message = fmt.Sprintf("%d %s\n%s", statusCode, statusText, string(prettyJSON))
			}
		}
	}

	return models.TestResult{
		Method:   method,
		Endpoint: endpoint,
		Status:   statusStr,
		Message:  message,
		Duration: duration,
		LogEntry: logEntry,
	}, nil
}

// ExecuteCustomRequestCmd wraps ExecuteCustomRequest as a Bubble Tea command
func ExecuteCustomRequestCmd(method, endpoint string, headers map[string]string, body string, auth *models.AuthConfig, verbose bool) tea.Cmd {
	return func() tea.Msg {
		result, err := ExecuteCustomRequest(method, endpoint, headers, body, auth, verbose)
		if err != nil {
			return TestErrorMsg{Err: err}
		}
		return TestCompleteMsg{Results: []models.TestResult{result}}
	}
}

// ValidateJSONBody validates that a string is valid JSON
func ValidateJSONBody(body string) error {
	if body == "" {
		return nil // Empty body is valid
	}
	var jsonData interface{}
	return json.Unmarshal([]byte(body), &jsonData)
}

// FormatJSONBody formats a JSON string with indentation
func FormatJSONBody(body string) (string, error) {
	if body == "" {
		return "", nil
	}
	var jsonData interface{}
	if err := json.Unmarshal([]byte(body), &jsonData); err != nil {
		return "", err
	}
	formatted, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return "", err
	}
	return string(formatted), nil
}

package testing

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

// isRetryableError determines if an error should trigger a retry
// Returns true for network errors, timeouts, and server errors (5xx)
// Returns false for validation errors, client errors (4xx), and successful responses
func isRetryableError(err error, statusCode int) bool {
	// If no error and status is not a server error, don't retry
	if err == nil && statusCode < 500 {
		return false
	}

	// Server errors (5xx) are retryable
	if statusCode >= 500 && statusCode < 600 {
		return true
	}

	// If there's an error, check if it's a network/timeout error
	if err != nil {
		errStr := strings.ToLower(err.Error())
		
		// Network connectivity issues
		if strings.Contains(errStr, "connection refused") ||
			strings.Contains(errStr, "connection reset") ||
			strings.Contains(errStr, "no such host") ||
			strings.Contains(errStr, "network is unreachable") ||
			strings.Contains(errStr, "broken pipe") {
			return true
		}

		// Timeout errors
		if strings.Contains(errStr, "timeout") ||
			strings.Contains(errStr, "deadline exceeded") ||
			strings.Contains(errStr, "context canceled") {
			return true
		}

		// TLS/SSL errors that might be temporary
		if strings.Contains(errStr, "tls handshake") ||
			strings.Contains(errStr, "eof") {
			return true
		}
	}

	return false
}

// executeWithRetry executes an HTTP request with automatic retry logic
// Implements exponential backoff: delay doubles after each retry
// Only retries on network errors and server errors (5xx), not on client errors (4xx)
func executeWithRetry(
	method, url string,
	body []byte,
	auth *models.AuthConfig,
	verbose bool,
	maxRetries int,
	initialDelay int,
) (int, *http.Response, *models.LogEntry, int, error) {
	var lastErr error
	var statusCode int
	var resp *http.Response
	var log *models.LogEntry
	retryCount := 0

	// Ensure reasonable retry limits
	if maxRetries < 0 {
		maxRetries = 0
	}
	if maxRetries > 10 {
		maxRetries = 10 // Cap at 10 retries
	}
	if initialDelay < 100 {
		initialDelay = 100 // Minimum 100ms
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Execute the request
		statusCode, resp, log, lastErr = TestEndpoint(method, url, body, auth, verbose)

		// Check if we should retry
		shouldRetry := isRetryableError(lastErr, statusCode)

		// If success or non-retryable error, return immediately
		if !shouldRetry {
			return statusCode, resp, log, retryCount, lastErr
		}

		// If this was not the last attempt, increment retry count and wait
		if attempt < maxRetries {
			retryCount++
			
			// Calculate exponential backoff delay
			delay := time.Duration(initialDelay) * time.Millisecond
			for i := 0; i < attempt; i++ {
				delay *= 2 // Double the delay for each retry
			}

			// Cap maximum delay at 30 seconds
			if delay > 30*time.Second {
				delay = 30 * time.Second
			}

			time.Sleep(delay)
		}
	}

	// All retries exhausted
	if lastErr != nil {
		return statusCode, resp, log, retryCount, fmt.Errorf("max retries (%d) exceeded: %w", maxRetries, lastErr)
	}
	return statusCode, resp, log, retryCount, fmt.Errorf("max retries (%d) exceeded: server returned %d", maxRetries, statusCode)
}

// TestEndpointWithRetry is a wrapper around executeWithRetry that returns the standard TestEndpoint signature
// plus the retry count for tracking
func TestEndpointWithRetry(
	method, url string,
	body []byte,
	auth *models.AuthConfig,
	verbose bool,
	maxRetries int,
	retryDelay int,
) (int, *http.Response, *models.LogEntry, int, error) {
	return executeWithRetry(method, url, body, auth, verbose, maxRetries, retryDelay)
}

package testing

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestIsRetryableError_ServerErrors verifies 5xx errors are retryable
func TestIsRetryableError_ServerErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		err        error
		want       bool
	}{
		{"500 Internal Server Error", 500, nil, true},
		{"502 Bad Gateway", 502, nil, true},
		{"503 Service Unavailable", 503, nil, true},
		{"504 Gateway Timeout", 504, nil, true},
		{"599 Custom Server Error", 599, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRetryableError(tt.err, tt.statusCode)
			if got != tt.want {
				t.Errorf("isRetryableError() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsRetryableError_ClientErrors verifies 4xx errors are NOT retryable
func TestIsRetryableError_ClientErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		err        error
		want       bool
	}{
		{"400 Bad Request", 400, nil, false},
		{"401 Unauthorized", 401, nil, false},
		{"403 Forbidden", 403, nil, false},
		{"404 Not Found", 404, nil, false},
		{"422 Unprocessable Entity", 422, nil, false},
		{"429 Too Many Requests", 429, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRetryableError(tt.err, tt.statusCode)
			if got != tt.want {
				t.Errorf("isRetryableError() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsRetryableError_SuccessfulResponses verifies 2xx/3xx are NOT retryable
func TestIsRetryableError_SuccessfulResponses(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		err        error
		want       bool
	}{
		{"200 OK", 200, nil, false},
		{"201 Created", 201, nil, false},
		{"204 No Content", 204, nil, false},
		{"301 Moved Permanently", 301, nil, false},
		{"302 Found", 302, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRetryableError(tt.err, tt.statusCode)
			if got != tt.want {
				t.Errorf("isRetryableError() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsRetryableError_NetworkErrors verifies network errors are retryable
func TestIsRetryableError_NetworkErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"Connection refused", errors.New("connection refused"), true},
		{"Connection reset", errors.New("connection reset by peer"), true},
		{"No such host", errors.New("no such host"), true},
		{"Network unreachable", errors.New("network is unreachable"), true},
		{"Broken pipe", errors.New("broken pipe"), true},
		{"Timeout", errors.New("i/o timeout"), true},
		{"Deadline exceeded", errors.New("context deadline exceeded"), true},
		{"TLS handshake", errors.New("tls handshake error"), true},
		{"EOF", errors.New("EOF"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRetryableError(tt.err, 0)
			if got != tt.want {
				t.Errorf("isRetryableError() = %v, want %v for error: %v", got, tt.want, tt.err)
			}
		})
	}
}

// TestIsRetryableError_NonRetryableErrors verifies validation/application errors are NOT retryable
func TestIsRetryableError_NonRetryableErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"Validation error", errors.New("validation failed"), false},
		{"Parse error", errors.New("failed to parse response"), false},
		{"Generic error", errors.New("something went wrong"), false},
		{"Nil error with 200", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRetryableError(tt.err, 200)
			if got != tt.want {
				t.Errorf("isRetryableError() = %v, want %v for error: %v", got, tt.want, tt.err)
			}
		})
	}
}

// TestExecuteWithRetry_SuccessFirstAttempt verifies no retry on immediate success
func TestExecuteWithRetry_SuccessFirstAttempt(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	statusCode, resp, _, retryCount, err := executeWithRetry(
		"GET",
		server.URL+"/test",
		nil,
		nil,
		false,
		3,
		100,
	)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if resp != nil {
		resp.Body.Close()
	}
	if statusCode != 200 {
		t.Errorf("Expected status 200, got: %d", statusCode)
	}
	if retryCount != 0 {
		t.Errorf("Expected 0 retries, got: %d", retryCount)
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got: %d", attempts)
	}
}

// TestExecuteWithRetry_SuccessAfterRetries verifies successful retry
func TestExecuteWithRetry_SuccessAfterRetries(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			// Fail first 2 attempts
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"error": "service unavailable"}`))
		} else {
			// Success on 3rd attempt
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		}
	}))
	defer server.Close()

	statusCode, resp, _, retryCount, err := executeWithRetry(
		"GET",
		server.URL+"/test",
		nil,
		nil,
		false,
		3,
		50, // Short delay for fast testing
	)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if resp != nil {
		resp.Body.Close()
	}
	if statusCode != 200 {
		t.Errorf("Expected status 200, got: %d", statusCode)
	}
	if retryCount != 2 {
		t.Errorf("Expected 2 retries, got: %d", retryCount)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got: %d", attempts)
	}
}

// TestExecuteWithRetry_MaxRetriesExceeded verifies retry exhaustion
func TestExecuteWithRetry_MaxRetriesExceeded(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		// Always fail
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server error"}`))
	}))
	defer server.Close()

	statusCode, resp, _, retryCount, err := executeWithRetry(
		"GET",
		server.URL+"/test",
		nil,
		nil,
		false,
		3,
		50, // Short delay for fast testing
	)

	// Should return an error after exhausting retries
	if err == nil {
		t.Error("Expected error after max retries exceeded, got nil")
	}
	if resp != nil {
		resp.Body.Close()
	}
	if statusCode != 500 {
		t.Errorf("Expected status 500, got: %d", statusCode)
	}
	if retryCount != 3 {
		t.Errorf("Expected 3 retries, got: %d", retryCount)
	}
	// 1 initial attempt + 3 retries = 4 total attempts
	if attempts != 4 {
		t.Errorf("Expected 4 attempts, got: %d", attempts)
	}
}

// TestExecuteWithRetry_NoRetryOn4xx verifies client errors don't retry
func TestExecuteWithRetry_NoRetryOn4xx(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	}))
	defer server.Close()

	statusCode, resp, _, retryCount, err := executeWithRetry(
		"GET",
		server.URL+"/test",
		nil,
		nil,
		false,
		3,
		100,
	)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if resp != nil {
		resp.Body.Close()
	}
	if statusCode != 404 {
		t.Errorf("Expected status 404, got: %d", statusCode)
	}
	if retryCount != 0 {
		t.Errorf("Expected 0 retries (4xx should not retry), got: %d", retryCount)
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt (no retries), got: %d", attempts)
	}
}

// TestExecuteWithRetry_ExponentialBackoff verifies delay increases exponentially
func TestExecuteWithRetry_ExponentialBackoff(t *testing.T) {
	attempts := 0
	attemptTimes := []time.Time{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		attemptTimes = append(attemptTimes, time.Now())
		// Always fail to force retries
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	initialDelay := 100 // 100ms
	_, resp, _, retryCount, _ := executeWithRetry(
		"GET",
		server.URL+"/test",
		nil,
		nil,
		false,
		3,
		initialDelay,
	)
	if resp != nil {
		resp.Body.Close()
	}

	if retryCount != 3 {
		t.Errorf("Expected 3 retries, got: %d", retryCount)
	}

	// Verify exponential backoff delays
	// Expected: 100ms, 200ms, 400ms (approximately)
	if len(attemptTimes) == 4 {
		delay1 := attemptTimes[1].Sub(attemptTimes[0]).Milliseconds()
		delay2 := attemptTimes[2].Sub(attemptTimes[1]).Milliseconds()
		delay3 := attemptTimes[3].Sub(attemptTimes[2]).Milliseconds()

		// Allow some tolerance (±20ms) due to execution time
		tolerance := int64(30)
		
		if delay1 < int64(initialDelay)-tolerance || delay1 > int64(initialDelay)+tolerance {
			t.Errorf("First delay expected ~%dms, got %dms", initialDelay, delay1)
		}
		if delay2 < int64(initialDelay*2)-tolerance || delay2 > int64(initialDelay*2)+tolerance {
			t.Errorf("Second delay expected ~%dms, got %dms", initialDelay*2, delay2)
		}
		if delay3 < int64(initialDelay*4)-tolerance || delay3 > int64(initialDelay*4)+tolerance {
			t.Errorf("Third delay expected ~%dms, got %dms", initialDelay*4, delay3)
		}
	}
}

// TestExecuteWithRetry_RetryCaps verifies max retries and delay caps
func TestExecuteWithRetry_RetryCaps(t *testing.T) {
	t.Run("Max retries capped at 10", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		_, resp, _, retryCount, _ := executeWithRetry(
			"GET",
			server.URL+"/test",
			nil,
			nil,
			false,
			100, // Request 100 retries, should cap at 10
			10,
		)
		if resp != nil {
			resp.Body.Close()
		}

		if retryCount > 10 {
			t.Errorf("Expected max 10 retries, got: %d", retryCount)
		}
	})

	t.Run("Negative retries becomes 0", func(t *testing.T) {
		attempts := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		_, resp, _, retryCount, _ := executeWithRetry(
			"GET",
			server.URL+"/test",
			nil,
			nil,
			false,
			-5, // Negative retries
			100,
		)
		if resp != nil {
			resp.Body.Close()
		}

		if retryCount != 0 {
			t.Errorf("Expected 0 retries (negative clamped to 0), got: %d", retryCount)
		}
		if attempts != 1 {
			t.Errorf("Expected 1 attempt (no retries), got: %d", attempts)
		}
	})

	t.Run("Min delay is 100ms", func(t *testing.T) {
		attempts := 0
		attemptTimes := []time.Time{}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attempts++
			attemptTimes = append(attemptTimes, time.Now())
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		_, resp, _, _, _ := executeWithRetry(
			"GET",
			server.URL+"/test",
			nil,
			nil,
			false,
			1,
			10, // Try to set very short delay, should be clamped to 100ms
		)
		if resp != nil {
			resp.Body.Close()
		}

		if len(attemptTimes) == 2 {
			delay := attemptTimes[1].Sub(attemptTimes[0]).Milliseconds()
			if delay < 100 {
				t.Errorf("Expected min delay 100ms, got %dms", delay)
			}
		}
	})
}

// TestTestEndpointWithRetry_Integration verifies the wrapper function
func TestTestEndpointWithRetry_Integration(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusBadGateway)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data": "success"}`))
		}
	}))
	defer server.Close()

	statusCode, resp, _, retryCount, err := TestEndpointWithRetry(
		"GET",
		server.URL+"/api/test",
		nil,
		nil,
		false,
		3,
		50,
	)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if resp != nil {
		resp.Body.Close()
	}
	if statusCode != 200 {
		t.Errorf("Expected status 200, got: %d", statusCode)
	}
	if retryCount != 1 {
		t.Errorf("Expected 1 retry, got: %d", retryCount)
	}
	if attempts != 2 {
		t.Errorf("Expected 2 total attempts, got: %d", attempts)
	}
}

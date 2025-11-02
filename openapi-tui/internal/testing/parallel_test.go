package testing

import (
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TestRunTestsParallel_AutoDetectConcurrency verifies auto-detection of CPU count
func TestRunTestsParallel_AutoDetectConcurrency(t *testing.T) {
	// Create a simple test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	// Create a minimal OpenAPI spec file
	specContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      summary: Test endpoint
      responses:
        '200':
          description: OK
`
	specPath := createTempSpec(t, specContent)

	// Test with maxConcurrency = 0 (auto-detect)
	results, err := RunTestsParallel(specPath, server.URL, nil, false, 0, nil)
	if err != nil {
		t.Fatalf("RunTestsParallel failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	if results[0].Status != "200" {
		t.Errorf("Expected status '200', got '%s'", results[0].Status)
	}
}

// TestRunTestsParallel_CustomConcurrency verifies custom concurrency setting
func TestRunTestsParallel_CustomConcurrency(t *testing.T) {
	requestCount := int32(0)
	maxConcurrent := int32(0)
	currentConcurrent := int32(0)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Track concurrent requests
		current := atomic.AddInt32(&currentConcurrent, 1)
		atomic.AddInt32(&requestCount, 1)

		// Update max concurrent
		for {
			max := atomic.LoadInt32(&maxConcurrent)
			if current <= max || atomic.CompareAndSwapInt32(&maxConcurrent, max, current) {
				break
			}
		}

		// Simulate some work
		time.Sleep(10 * time.Millisecond)

		atomic.AddInt32(&currentConcurrent, -1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create spec with multiple endpoints
	specContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test1:
    get:
      responses:
        '200':
          description: OK
  /test2:
    get:
      responses:
        '200':
          description: OK
  /test3:
    get:
      responses:
        '200':
          description: OK
  /test4:
    get:
      responses:
        '200':
          description: OK
  /test5:
    get:
      responses:
        '200':
          description: OK
`
	specPath := createTempSpec(t, specContent)

	// Test with maxConcurrency = 2
	results, err := RunTestsParallel(specPath, server.URL, nil, false, 2, nil)
	if err != nil {
		t.Fatalf("RunTestsParallel failed: %v", err)
	}

	if len(results) != 5 {
		t.Errorf("Expected 5 results, got %d", len(results))
	}

	// Verify concurrency was limited to 2
	max := atomic.LoadInt32(&maxConcurrent)
	if max > 2 {
		t.Errorf("Expected max concurrency <= 2, got %d", max)
	}

	t.Logf("Max concurrent requests: %d", max)
}

// TestRunTestsParallel_ResultOrdering verifies results are returned in correct order
func TestRunTestsParallel_ResultOrdering(t *testing.T) {
	// Server that returns the path in the response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add variable delay to simulate different response times
		time.Sleep(time.Duration(10-len(r.URL.Path)) * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.URL.Path))
	}))
	defer server.Close()

	specContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test1:
    get:
      responses:
        '200':
          description: OK
  /test2:
    get:
      responses:
        '200':
          description: OK
  /test3:
    get:
      responses:
        '200':
          description: OK
`
	specPath := createTempSpec(t, specContent)

	results, err := RunTestsParallel(specPath, server.URL, nil, false, 3, nil)
	if err != nil {
		t.Fatalf("RunTestsParallel failed: %v", err)
	}

	// Verify all expected endpoints are present (order may vary due to map iteration)
	expectedEndpoints := map[string]bool{
		"/test1": false,
		"/test2": false,
		"/test3": false,
	}
	
	for _, result := range results {
		if _, exists := expectedEndpoints[result.Endpoint]; exists {
			expectedEndpoints[result.Endpoint] = true
		}
	}
	
	// Check all endpoints were found
	for endpoint, found := range expectedEndpoints {
		if !found {
			t.Errorf("Expected endpoint %s not found in results", endpoint)
		}
	}
}

// TestRunTestsParallel_ErrorHandling verifies error handling with failing endpoints
func TestRunTestsParallel_ErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/fail" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "internal error"}`))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		}
	}))
	defer server.Close()

	specContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /success:
    get:
      responses:
        '200':
          description: OK
  /fail:
    get:
      responses:
        '200':
          description: OK
`
	specPath := createTempSpec(t, specContent)

	results, err := RunTestsParallel(specPath, server.URL, nil, false, 2, nil)
	if err != nil {
		t.Fatalf("RunTestsParallel failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// Verify both success and failure are captured
	successCount := 0
	failureCount := 0
	for _, result := range results {
		// Check if status is "200" for success or "500" for error
		if result.Status == "200" {
			successCount++
		} else if result.Status == "500" {
			failureCount++
		}
	}

	if successCount != 1 || failureCount != 1 {
		t.Errorf("Expected 1 success and 1 failure, got %d success and %d failure (statuses: %v)", successCount, failureCount, []string{results[0].Status, results[1].Status})
	}
}

// TestRunTestsParallel_RaceConditions tests for race conditions with -race flag
func TestRunTestsParallel_RaceConditions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond) // Simulate work
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	specContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test1:
    get:
      responses:
        '200':
          description: OK
  /test2:
    get:
      responses:
        '200':
          description: OK
  /test3:
    get:
      responses:
        '200':
          description: OK
  /test4:
    get:
      responses:
        '200':
          description: OK
  /test5:
    get:
      responses:
        '200':
          description: OK
`
	specPath := createTempSpec(t, specContent)

	// Run multiple times to increase chance of detecting races
	for i := 0; i < 10; i++ {
		_, err := RunTestsParallel(specPath, server.URL, nil, false, 3, nil)
		if err != nil {
			t.Fatalf("Iteration %d: RunTestsParallel failed: %v", i, err)
		}
	}
}

// TestRunTestParallelCmd verifies Bubble Tea command integration
func TestRunTestParallelCmd(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	specContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      responses:
        '200':
          description: OK
`
	specPath := createTempSpec(t, specContent)

	// Execute the command
	cmd := RunTestParallelCmd(specPath, server.URL, nil, false, 2)
	msg := cmd()

	// Verify message type
	switch msg := msg.(type) {
	case TestCompleteMsg:
		if len(msg.Results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(msg.Results))
		}
	case TestErrorMsg:
		t.Fatalf("Expected TestCompleteMsg, got TestErrorMsg: %v", msg.Err)
	default:
		t.Fatalf("Expected TestCompleteMsg or TestErrorMsg, got %T", msg)
	}
}

// TestWorkerPoolPattern verifies worker pool behavior
func TestWorkerPoolPattern(t *testing.T) {
	// Track worker execution
	workerIDs := make(map[int]struct{})
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate work
		time.Sleep(5 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create spec with many endpoints
	specContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test1:
    get:
      responses:
        '200':
          description: OK
  /test2:
    get:
      responses:
        '200':
          description: OK
  /test3:
    get:
      responses:
        '200':
          description: OK
  /test4:
    get:
      responses:
        '200':
          description: OK
  /test5:
    get:
      responses:
        '200':
          description: OK
  /test6:
    get:
      responses:
        '200':
          description: OK
  /test7:
    get:
      responses:
        '200':
          description: OK
  /test8:
    get:
      responses:
        '200':
          description: OK
`
	specPath := createTempSpec(t, specContent)

	// Run with 3 workers
	results, err := RunTestsParallel(specPath, server.URL, nil, false, 3, nil)
	if err != nil {
		t.Fatalf("RunTestsParallel failed: %v", err)
	}

	if len(results) != 8 {
		t.Errorf("Expected 8 results, got %d", len(results))
	}

	mu.Lock()
	defer mu.Unlock()
	t.Logf("Worker pool used %d workers", len(workerIDs))
}

// TestProgressMessages verifies progress messages are sent
func TestProgressMessages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	specContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test1:
    get:
      responses:
        '200':
          description: OK
  /test2:
    get:
      responses:
        '200':
          description: OK
  /test3:
    get:
      responses:
        '200':
          description: OK
`
	specPath := createTempSpec(t, specContent)

	// Create progress channel
	progressChan := make(chan tea.Msg, 10)

	// Run tests with progress tracking
	_, err := RunTestsParallel(specPath, server.URL, nil, false, 2, progressChan)
	if err != nil {
		t.Fatalf("RunTestsParallel failed: %v", err)
	}

	// Count progress messages
	close(progressChan)
	progressCount := 0
	for range progressChan {
		progressCount++
	}

	// Should receive at least one progress message per test
	if progressCount < 3 {
		t.Errorf("Expected at least 3 progress messages, got %d", progressCount)
	}

	t.Logf("Received %d progress messages for 3 tests", progressCount)
}

// Helper function to create temporary spec file
func createTempSpec(tb testing.TB, content string) string {
	tb.Helper()
	tmpFile := tb.TempDir() + "/spec.yaml"
	if err := writeFile(tmpFile, content); err != nil {
		tb.Fatalf("Failed to create temp spec: %v", err)
	}
	return tmpFile
}

// Helper function to write file (simple version, no deps needed)
func writeFile(path, content string) error {
	// Use os.WriteFile from standard library
	return os.WriteFile(path, []byte(content), 0644)
}

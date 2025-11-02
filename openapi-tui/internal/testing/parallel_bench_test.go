package testing

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// BenchmarkSequentialExecution benchmarks the original sequential test execution
func BenchmarkSequentialExecution(b *testing.B) {
	// Create test server with simulated delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate realistic API response time (10ms)
		// time.Sleep(10 * time.Millisecond) // Commented to avoid slowing CI
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	// Create spec with multiple endpoints
	specContent := generateBenchmarkSpec(10) // 10 endpoints
	specPath := createTempSpec(b, specContent)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := RunTests(specPath, server.URL, nil, false)
		if err != nil {
			b.Fatalf("RunTests failed: %v", err)
		}
	}
}

// BenchmarkParallelExecution benchmarks the new parallel test execution
func BenchmarkParallelExecution(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate realistic API response time
		// time.Sleep(10 * time.Millisecond) // Commented to avoid slowing CI
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	specContent := generateBenchmarkSpec(10) // 10 endpoints
	specPath := createTempSpec(b, specContent)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := RunTestsParallel(specPath, server.URL, nil, false, 0, nil)
		if err != nil {
			b.Fatalf("RunTestsParallel failed: %v", err)
		}
	}
}

// BenchmarkParallelExecution_CustomConcurrency tests different concurrency levels
func BenchmarkParallelExecution_CustomConcurrency(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	specContent := generateBenchmarkSpec(20) // 20 endpoints
	specPath := createTempSpec(b, specContent)

	// Test different concurrency levels
	for _, concurrency := range []int{1, 2, 4, 8} {
		b.Run(fmt.Sprintf("Concurrency%d", concurrency), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := RunTestsParallel(specPath, server.URL, nil, false, concurrency, nil)
				if err != nil {
					b.Fatalf("RunTestsParallel failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkScaling tests performance with varying endpoint counts
func BenchmarkScaling(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	for _, endpointCount := range []int{5, 10, 25, 50} {
		specContent := generateBenchmarkSpec(endpointCount)
		specPath := createTempSpec(b, specContent)

		b.Run(fmt.Sprintf("Sequential_%dEndpoints", endpointCount), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := RunTests(specPath, server.URL, nil, false)
				if err != nil {
					b.Fatalf("RunTests failed: %v", err)
				}
			}
		})

		b.Run(fmt.Sprintf("Parallel_%dEndpoints", endpointCount), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := RunTestsParallel(specPath, server.URL, nil, false, 0, nil)
				if err != nil {
					b.Fatalf("RunTestsParallel failed: %v", err)
				}
			}
		})
	}
}

// generateBenchmarkSpec creates an OpenAPI spec with N endpoints for benchmarking
func generateBenchmarkSpec(endpointCount int) string {
	spec := `
openapi: 3.0.0
info:
  title: Benchmark API
  version: 1.0.0
paths:
`
	for i := 1; i <= endpointCount; i++ {
		spec += fmt.Sprintf(`  /test%d:
    get:
      summary: Test endpoint %d
      responses:
        '200':
          description: OK
`, i, i)
	}
	return spec
}

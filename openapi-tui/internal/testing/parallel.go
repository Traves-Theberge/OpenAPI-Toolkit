package testing

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/errors"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/validation"
	"github.com/getkin/kin-openapi/openapi3"
)

// TestJob represents a single test to execute
type TestJob struct {
	Method      string
	Path        string
	Endpoint    string
	RequestBody []byte
	Operation   *openapi3.Operation
}

// TestProgressMsg is sent during parallel execution to update progress
type TestProgressMsg struct {
	Completed int
	Total     int
	Latest    *models.TestResult // Most recent result
}

// RunTestsParallel executes API tests concurrently with a worker pool
// maxConcurrency: maximum number of concurrent requests (0 = auto-detect)
func RunTestsParallel(specPath, baseURL string, auth *models.AuthConfig, verbose bool, maxConcurrency int, maxRetries int, retryDelay int, progressChan chan<- tea.Msg) ([]models.TestResult, error) {
	// Load and validate the OpenAPI spec
	loader := &openapi3.Loader{IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		return nil, errors.EnhanceFileError(err, specPath)
	}

	// Auto-detect concurrency if not specified
	if maxConcurrency <= 0 {
		maxConcurrency = runtime.NumCPU()
		// Cap at reasonable maximum to avoid overwhelming servers
		if maxConcurrency > 10 {
			maxConcurrency = 10
		}
	}

	// Collect all test jobs
	var jobs []TestJob
	if doc.Paths != nil {
		for path, pathItem := range doc.Paths.Map() {
			for method, operation := range pathItem.Operations() {
				// Construct full endpoint URL
				endpoint := baseURL + ReplacePlaceholders(path)
				endpoint += BuildQueryParams(operation)

				// Generate request body if needed
				var requestBody []byte
				if strings.ToUpper(method) == "POST" || strings.ToUpper(method) == "PUT" || strings.ToUpper(method) == "PATCH" {
					requestBody, err = GenerateRequestBody(operation)
					if err != nil {
						// Add error result and continue
						jobs = append(jobs, TestJob{
							Method:      method,
							Path:        path,
							Endpoint:    endpoint,
							RequestBody: nil,
							Operation:   nil, // Signal error
						})
						continue
					}
				}

				jobs = append(jobs, TestJob{
					Method:      method,
					Path:        path,
					Endpoint:    endpoint,
					RequestBody: requestBody,
					Operation:   operation,
				})
			}
		}
	}

	totalJobs := len(jobs)
	if totalJobs == 0 {
		return []models.TestResult{}, nil
	}

	// Create worker pool with indexed jobs for maintaining order
	type IndexedJob struct {
		Index int
		Job   TestJob
	}
	type IndexedResult struct {
		Index  int
		Result models.TestResult
	}
	
	jobChan := make(chan IndexedJob, totalJobs)
	resultChan := make(chan IndexedResult, totalJobs)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for indexedJob := range jobChan {
				result := executeTestJob(indexedJob.Job, auth, verbose, maxRetries, retryDelay)
				resultChan <- IndexedResult{Index: indexedJob.Index, Result: result}
				
				// Send progress update if channel provided
				if progressChan != nil {
					select {
					case progressChan <- TestProgressMsg{
						Completed: indexedJob.Index + 1,
						Total:     totalJobs,
						Latest:    &result,
					}:
					default:
						// Don't block if UI isn't ready
					}
				}
			}
		}()
	}

	// Send all jobs to workers with indices
	for i, job := range jobs {
		jobChan <- IndexedJob{Index: i, Job: job}
	}
	close(jobChan)

	// Wait for all workers to finish
	wg.Wait()
	close(resultChan)

	// Collect results maintaining original order
	resultMap := make(map[int]models.TestResult)
	for indexedResult := range resultChan {
		resultMap[indexedResult.Index] = indexedResult.Result
	}
	
	// Reconstruct ordered results
	results := make([]models.TestResult, totalJobs)
	for i := 0; i < totalJobs; i++ {
		results[i] = resultMap[i]
	}

	return results, nil
}

// executeTestJob runs a single test job with retry support
func executeTestJob(job TestJob, auth *models.AuthConfig, verbose bool, maxRetries int, retryDelay int) models.TestResult {
	// Handle jobs that failed during body generation
	if job.Operation == nil && job.RequestBody == nil {
		return models.TestResult{
			Method:     job.Method,
			Endpoint:   job.Path,
			Status:     "ERR",
			Message:    "Failed to generate request body",
			RetryCount: 0,
		}
	}

	// Execute the test with retry logic
	startTime := time.Now()
	status, resp, logEntry, retryCount, err := TestEndpointWithRetry(job.Method, job.Endpoint, job.RequestBody, auth, verbose, maxRetries, retryDelay)
	duration := time.Since(startTime)

	message := "OK"
	if err != nil {
		message = err.Error()
	} else if resp != nil {
		// Validate response against spec
		validationResult := validation.ValidateResponse(resp, job.Operation, status)
		
		// Close response body after validation
		if resp.Body != nil {
			io.Copy(io.Discard, resp.Body) // Drain body
			resp.Body.Close()
		}
		
		if !validationResult.Valid {
			message = "Response validation failed"
			if len(validationResult.SchemaErrors) > 0 {
				message = validationResult.SchemaErrors[0]
			}
		} else if validationResult.StatusValid {
			if retryCount > 0 {
				message = fmt.Sprintf("OK (validated, %d retries)", retryCount)
			} else {
				message = "OK (validated)"
			}
		}
	}

	// Format status for display
	statusStr := fmt.Sprintf("%d", status)
	if err != nil {
		statusStr = "ERR"
	}

	return models.TestResult{
		Method:     job.Method,
		Endpoint:   job.Path,
		Status:     statusStr,
		Message:    message,
		Duration:   duration,
		LogEntry:   logEntry,
		RetryCount: retryCount,
	}
}

// RunTestParallelCmd wraps RunTestsParallel in a Bubble Tea command
// Uses same interface as RunTestCmd for easy migration
func RunTestParallelCmd(specPath, baseURL string, auth *models.AuthConfig, verbose bool, maxConcurrency int, maxRetries int, retryDelay int) tea.Cmd {
	return func() tea.Msg {
		results, err := RunTestsParallel(specPath, baseURL, auth, verbose, maxConcurrency, maxRetries, retryDelay, nil)
		if err != nil {
			return TestErrorMsg{Err: err}
		}
		return TestCompleteMsg{Results: results}
	}
}

// RunTestParallelCmdWithSelection executes tests for only selected endpoints
func RunTestParallelCmdWithSelection(specPath, baseURL string, auth *models.AuthConfig, verbose bool, maxConcurrency int, maxRetries int, retryDelay int, selectedEndpoints []models.EndpointInfo) tea.Cmd {
	return func() tea.Msg {
		results, err := RunTestsParallelWithSelection(specPath, baseURL, auth, verbose, maxConcurrency, maxRetries, retryDelay, nil, selectedEndpoints)
		if err != nil {
			return TestErrorMsg{Err: err}
		}
		return TestCompleteMsg{Results: results}
	}
}

// RunTestsParallelWithSelection runs tests for only the selected endpoints
func RunTestsParallelWithSelection(specPath, baseURL string, auth *models.AuthConfig, verbose bool, maxConcurrency int, maxRetries int, retryDelay int, progressChan chan<- tea.Msg, selectedEndpoints []models.EndpointInfo) ([]models.TestResult, error) {
	// Load and validate the OpenAPI spec
	loader := &openapi3.Loader{IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		return nil, errors.EnhanceFileError(err, specPath)
	}

	// Auto-detect concurrency if not specified
	if maxConcurrency <= 0 {
		maxConcurrency = runtime.NumCPU()
		if maxConcurrency > 10 {
			maxConcurrency = 10
		}
	}

	// Create a map of selected endpoints for quick lookup
	selectedMap := make(map[string]map[string]bool)
	for _, ep := range selectedEndpoints {
		if selectedMap[ep.Path] == nil {
			selectedMap[ep.Path] = make(map[string]bool)
		}
		selectedMap[ep.Path][ep.Method] = true
	}

	// Build job queue with only selected endpoints
	var jobs []TestJob
	for path, pathItem := range doc.Paths.Map() {
		if pathMethods, exists := selectedMap[path]; exists {
			operations := map[string]*openapi3.Operation{
				"GET":     pathItem.Get,
				"POST":    pathItem.Post,
				"PUT":     pathItem.Put,
				"PATCH":   pathItem.Patch,
				"DELETE":  pathItem.Delete,
				"HEAD":    pathItem.Head,
				"OPTIONS": pathItem.Options,
			}

			for method, operation := range operations {
				if operation != nil && pathMethods[method] {
					// Generate request body if needed
					requestBody, _ := GenerateRequestBody(operation)

					// Build full endpoint URL
					endpoint := baseURL + ReplacePlaceholders(path)
					queryParams := BuildQueryParams(operation)
					if queryParams != "" {
						endpoint += queryParams
					}

					jobs = append(jobs, TestJob{
						Method:      method,
						Path:        path,
						Endpoint:    endpoint,
						RequestBody: requestBody,
						Operation:   operation,
					})
				}
			}
		}
	}

	// Execute jobs with worker pool
	results := make([]models.TestResult, len(jobs))
	var wg sync.WaitGroup
	jobChan := make(chan struct {
		job   TestJob
		index int
	}, len(jobs))

	// Send all jobs to channel
	for i, job := range jobs {
		jobChan <- struct {
			job   TestJob
			index int
		}{job, i}
	}
	close(jobChan)

	// Start workers
	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for jobWithIndex := range jobChan {
				result := executeTestJob(jobWithIndex.job, auth, verbose, maxRetries, retryDelay)
				results[jobWithIndex.index] = result

				// Send progress update if channel provided
				if progressChan != nil {
					progressChan <- TestProgressMsg{
						Completed: jobWithIndex.index + 1,
						Total:     len(jobs),
						Latest:    &result,
					}
				}
			}
		}()
	}

	wg.Wait()
	return results, nil
}

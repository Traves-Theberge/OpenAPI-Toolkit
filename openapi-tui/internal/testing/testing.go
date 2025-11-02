package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/errors"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/validation"
	"github.com/getkin/kin-openapi/openapi3"
)

// replacePlaceholders replaces path parameters like {id} with sample values
func ReplacePlaceholders(path string) string {
	// Replace {param} with sensible defaults
	re := regexp.MustCompile(`\{[^}]+\}`)
	return re.ReplaceAllString(path, "1")
}

// buildQueryParams constructs query parameters from operation parameters
func BuildQueryParams(operation *openapi3.Operation) string {
	if operation == nil || operation.Parameters == nil {
		return ""
	}

	var params []string
	for _, paramRef := range operation.Parameters {
		param := paramRef.Value
		if param == nil || param.In != "query" {
			continue
		}

		// Generate sample value based on schema
		value := "1" // Default
		if param.Schema != nil && param.Schema.Value != nil {
			schema := param.Schema.Value
			// Check if type contains specific values
			if schema.Type.Is("string") {
				if len(schema.Enum) > 0 {
					value = fmt.Sprintf("%v", schema.Enum[0])
				} else if schema.Example != nil {
					value = fmt.Sprintf("%v", schema.Example)
				} else {
					value = "test"
				}
			} else if schema.Type.Is("integer") || schema.Type.Is("number") {
				if schema.Example != nil {
					value = fmt.Sprintf("%v", schema.Example)
				} else {
					value = "1"
				}
			} else if schema.Type.Is("boolean") {
				value = "true"
			} else if schema.Type.Is("array") {
				value = "1,2,3" // Simple array representation
			}
		}

		params = append(params, fmt.Sprintf("%s=%s", param.Name, value))
	}

	if len(params) == 0 {
		return ""
	}
	return "?" + strings.Join(params, "&")
}

// generateRequestBody creates a sample JSON request body from an OpenAPI schema
// Generates realistic sample data based on schema properties, types, and examples
func GenerateRequestBody(operation *openapi3.Operation) ([]byte, error) {
	if operation == nil || operation.RequestBody == nil {
		return nil, nil
	}

	// Get the request body content for JSON
	requestBody := operation.RequestBody.Value
	if requestBody == nil || requestBody.Content == nil {
		return nil, nil
	}

	// Look for JSON content type
	jsonContent := requestBody.Content.Get("application/json")
	if jsonContent == nil || jsonContent.Schema == nil || jsonContent.Schema.Value == nil {
		return nil, nil
	}

	schema := jsonContent.Schema.Value

	// Generate sample data from schema
	sample := GenerateSampleFromSchema(schema)
	
	// Marshal to JSON
	jsonData, err := json.Marshal(sample)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	return jsonData, nil
}

// generateSampleFromSchema recursively generates sample data from an OpenAPI schema
func GenerateSampleFromSchema(schema *openapi3.Schema) interface{} {
	if schema == nil {
		return nil
	}

	// Use example if available
	if schema.Example != nil {
		return schema.Example
	}

	// Use default if available
	if schema.Default != nil {
		return schema.Default
	}

	// Generate based on type
	if schema.Type.Is("object") {
		obj := make(map[string]interface{})
		for propName, propRef := range schema.Properties {
			if propRef != nil && propRef.Value != nil {
				obj[propName] = GenerateSampleFromSchema(propRef.Value)
			}
		}
		return obj
	}

	if schema.Type.Is("array") {
		if schema.Items != nil && schema.Items.Value != nil {
			// Generate a single-item array
			return []interface{}{GenerateSampleFromSchema(schema.Items.Value)}
		}
		return []interface{}{}
	}

	if schema.Type.Is("string") {
		if len(schema.Enum) > 0 {
			return schema.Enum[0]
		}
		if schema.Format == "email" {
			return "user@example.com"
		}
		if schema.Format == "uri" || schema.Format == "url" {
			return "https://example.com"
		}
		if schema.Format == "date" {
			return "2024-01-01"
		}
		if schema.Format == "date-time" {
			return "2024-01-01T00:00:00Z"
		}
		return "sample"
	}

	if schema.Type.Is("integer") {
		if schema.Min != nil {
			return int(*schema.Min)
		}
		return 1
	}

	if schema.Type.Is("number") {
		if schema.Min != nil {
			return *schema.Min
		}
		return 1.0
	}

	if schema.Type.Is("boolean") {
		return true
	}

	// Default fallback
	return nil
}

// applyAuth applies authentication configuration to an HTTP request
func ApplyAuth(req *http.Request, auth *models.AuthConfig) {
	if auth == nil || auth.AuthType == "none" || auth.AuthType == "" {
		return
	}

	switch auth.AuthType {
	case "bearer":
		if auth.Token != "" {
			req.Header.Set("Authorization", "Bearer "+auth.Token)
		}
	case "apiKey":
		if auth.APIKeyName != "" && auth.Token != "" {
			if auth.APIKeyIn == "header" {
				req.Header.Set(auth.APIKeyName, auth.Token)
			} else if auth.APIKeyIn == "query" {
				// Add to query parameters
				q := req.URL.Query()
				q.Add(auth.APIKeyName, auth.Token)
				req.URL.RawQuery = q.Encode()
			}
		}
	case "basic":
		if auth.Username != "" {
			req.SetBasicAuth(auth.Username, auth.Password)
		}
	}
}

// testEndpoint performs an HTTP request to test an API endpoint
// Supports GET, POST, PUT, PATCH, DELETE methods with optional request bodies
// Returns status code, response object, log entry, and error
func TestEndpoint(method, url string, body []byte, auth *models.AuthConfig, verbose bool) (int, *http.Response, *models.LogEntry, error) {
	var req *http.Request
	var err error

	// Create request based on HTTP method
	method = strings.ToUpper(method)
	
	if body != nil && len(body) > 0 {
		// Create request with body
		req, err = http.NewRequest(method, url, bytes.NewReader(body))
		if err != nil {
			return 0, nil, nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		// Create request without body
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return 0, nil, nil, err
		}
	}

	// Apply authentication if configured
	ApplyAuth(req, auth)

	// Capture start time for duration measurement
	startTime := time.Now()

	// Execute request with timeout to prevent hanging
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		return 0, nil, nil, errors.EnhanceNetworkError(err, url)
	}

	// Create log entry if verbose mode is enabled
	var log *models.LogEntry
	if verbose {
		log = &models.LogEntry{
			RequestURL:  url,
			Duration:    duration,
			Timestamp:   startTime,
			RequestHeaders: make(map[string]string),
			ResponseHeaders: make(map[string]string),
		}

		// Capture request headers
		for k, v := range req.Header {
			if len(v) > 0 {
				log.RequestHeaders[k] = v[0]
			}
		}

		// Capture request body
		if len(body) > 0 {
			log.RequestBody = string(body)
			// Truncate if too large
			if len(log.RequestBody) > 500 {
				log.RequestBody = log.RequestBody[:500] + "... (truncated)"
			}
		}

		// Capture response headers
		for k, v := range resp.Header {
			if len(v) > 0 {
				log.ResponseHeaders[k] = v[0]
			}
		}

		// Capture response body (read and restore)
		if resp.Body != nil {
			bodyBytes, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil {
				log.ResponseBody = string(bodyBytes)
				// Truncate if too large
				if len(log.ResponseBody) > 500 {
					log.ResponseBody = log.ResponseBody[:500] + "... (truncated)"
				}
				// Restore body for further processing
				resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}
	}

	return resp.StatusCode, resp, log, nil
}

// runTests executes API tests against endpoints defined in OpenAPI spec
// Tests each endpoint with a simple request and records results
// Accepts optional auth configuration and verbose flag for detailed logging
func RunTests(specPath, baseURL string, auth *models.AuthConfig, verbose bool) ([]models.TestResult, error) {
	// Load and validate the OpenAPI spec
	loader := &openapi3.Loader{IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		return nil, errors.EnhanceFileError(err, specPath)
	}

	var results []models.TestResult

	// Iterate through all paths and operations in the spec
	if doc.Paths != nil {
		for path, pathItem := range doc.Paths.Map() {
			// Iterate through all paths and operations in the spec
			for method, operation := range pathItem.Operations() {
				// Construct full endpoint URL with placeholder replacement
				endpoint := baseURL + ReplacePlaceholders(path)
				
				// Add query parameters if defined
				queryParams := BuildQueryParams(operation)
				endpoint += queryParams

				// Generate request body if needed
				var requestBody []byte
				if strings.ToUpper(method) == "POST" || strings.ToUpper(method) == "PUT" || strings.ToUpper(method) == "PATCH" {
					requestBody, err = GenerateRequestBody(operation)
					if err != nil {
						// Log error but continue testing
						results = append(results, models.TestResult{
							Method:   method,
							Endpoint: path,
							Status:   "ERR",
							Message:  fmt.Sprintf("Failed to generate request body: %v", err),
						})
						continue
					}
				}

				// Test the endpoint and record result
				startTime := time.Now()
				status, resp, logEntry, err := TestEndpoint(method, endpoint, requestBody, auth, verbose)
				duration := time.Since(startTime)
				message := "OK"
				if err != nil {
					message = err.Error()
				} else if resp != nil {
					// Validate response against spec
					validationResult := validation.ValidateResponse(resp, operation, status)
					
					// Close response body after validation
					if resp.Body != nil {
						io.Copy(io.Discard, resp.Body) // Drain body
						resp.Body.Close()
					}
					
					if !validationResult.Valid {
						message = "Response validation failed"
						if len(validationResult.SchemaErrors) > 0 {
							message = validationResult.SchemaErrors[0] // Show first error
						}
					} else if validationResult.StatusValid {
						message = "OK (validated)"
					}
				}

				// Format status for display
				statusStr := fmt.Sprintf("%d", status)
				if err != nil {
					statusStr = "ERR"
				}

				// Add result to collection
				results = append(results, models.TestResult{
					Method:   method,
					Endpoint: path,
					Status:   statusStr,
					Message:  message,
					Duration: duration,
					LogEntry: logEntry,
				})
			}
		}
	}

	return results, nil
}

// runTestCmd wraps runTests in a Bubble Tea command
// Returns a message containing test results or error
func RunTestCmd(specPath, baseURL string, auth *models.AuthConfig, verbose bool) tea.Cmd {
	return func() tea.Msg {
		results, err := RunTests(specPath, baseURL, auth, verbose)
		if err != nil {
			return TestErrorMsg{Err: err}
		}
		return TestCompleteMsg{Results: results}
	}
}

// TestErrorMsg is sent when testing encounters an error
type TestErrorMsg struct {
	Err error
}

// TestCompleteMsg is sent when testing completes successfully
type TestCompleteMsg struct {
	Results []models.TestResult
}

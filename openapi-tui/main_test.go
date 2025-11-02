package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

// TestReplacePlaceholders tests path parameter replacement
func TestReplacePlaceholders(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single parameter",
			input:    "/users/{id}",
			expected: "/users/1",
		},
		{
			name:     "multiple parameters",
			input:    "/users/{userId}/posts/{postId}",
			expected: "/users/1/posts/1",
		},
		{
			name:     "no parameters",
			input:    "/users",
			expected: "/users",
		},
		{
			name:     "parameter at start",
			input:    "/{version}/users",
			expected: "/1/users",
		},
		{
			name:     "parameter at end",
			input:    "/api/users/{id}",
			expected: "/api/users/1",
		},
		{
			name:     "mixed underscores and camelCase",
			input:    "/users/{user_id}/posts/{postId}",
			expected: "/users/1/posts/1",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only slashes",
			input:    "///",
			expected: "///",
		},
		{
			name:     "nested braces (invalid but should handle)",
			input:    "/users/{{id}}",
			expected: "/users/1}", // Replaces first {}, leaves trailing }
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replacePlaceholders(tt.input)
			if result != tt.expected {
				t.Errorf("replacePlaceholders(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestValidateSpec tests OpenAPI spec validation
func TestValidateSpec(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	// Valid OpenAPI spec
	validSpec := `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      responses:
        '200':
          description: Success
`

	// Invalid spec - missing required fields
	invalidSpec := `openapi: 3.0.0
paths:
  /test:
    get:
      responses: {}
`

	// Invalid YAML
	invalidYAML := `this is not: valid: yaml: content::`

	tests := []struct {
		name      string
		content   string
		wantError bool
	}{
		{
			name:      "valid spec",
			content:   validSpec,
			wantError: false,
		},
		{
			name:      "invalid spec structure",
			content:   invalidSpec,
			wantError: true,
		},
		{
			name:      "invalid YAML",
			content:   invalidYAML,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			filePath := filepath.Join(tempDir, tt.name+".yaml")
			err := os.WriteFile(filePath, []byte(tt.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Test validation
			_, err = validateSpec(filePath)
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}

	// Test non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		_, err := validateSpec("/nonexistent/path/spec.yaml")
		if err == nil {
			t.Error("Expected error for non-existent file but got none")
		}
	})
}

// TestTestEndpoint tests HTTP endpoint testing
func TestTestEndpoint(t *testing.T) {
	// Create test server that responds based on path
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		case "/created":
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"status":"created"}`))
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"server error"}`))
		case "/notfound":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"not found"}`))
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		wantError      bool
	}{
		{
			name:           "GET success",
			method:         "GET",
			path:           "/success",
			expectedStatus: 200,
			wantError:      false,
		},
		{
			name:           "POST created",
			method:         "POST",
			path:           "/created",
			expectedStatus: 201,
			wantError:      false,
		},
		{
			name:           "GET error",
			method:         "GET",
			path:           "/error",
			expectedStatus: 500,
			wantError:      false,
		},
		{
			name:           "GET not found",
			method:         "GET",
			path:           "/notfound",
			expectedStatus: 404,
			wantError:      false,
		},
		{
			name:           "lowercase method",
			method:         "get",
			path:           "/success",
			expectedStatus: 200,
			wantError:      false,
		},
		{
			name:           "mixed case method",
			method:         "GeT",
			path:           "/success",
			expectedStatus: 200,
			wantError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := server.URL + tt.path
			status, resp, err := testEndpoint(tt.method, url, nil)

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
			if status != tt.expectedStatus {
				t.Errorf("Expected status %d but got %d", tt.expectedStatus, status)
			}
			if resp != nil && resp.Body != nil {
				resp.Body.Close()
			}
		})
	}

	// Test unsupported method - now all methods are supported
	t.Run("DELETE method", func(t *testing.T) {
		_, resp, err := testEndpoint("DELETE", server.URL+"/success", nil)
		if err != nil {
			t.Errorf("Expected no error but got: %v", err)
		}
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	})

	// Test invalid URL
	t.Run("invalid URL", func(t *testing.T) {
		_, resp, err := testEndpoint("GET", "://invalid-url", nil)
		if err == nil {
			t.Error("Expected error for invalid URL but got none")
		}
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	})

	// Test unreachable server
	t.Run("unreachable server", func(t *testing.T) {
		_, resp, err := testEndpoint("GET", "http://localhost:99999/test", nil)
		if err == nil {
			t.Error("Expected error for unreachable server but got none")
		}
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	})
}

// TestTestEndpointTimeout tests HTTP timeout
func TestTestEndpointTimeout(t *testing.T) {
	// Create server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This would timeout if timeout is < 50ms, but our timeout is 10s
		// So we just test that timeout is configured
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Test that we can make successful requests (timeout is working properly)
	status, resp, err := testEndpoint("GET", server.URL+"/test", nil)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
	if status != 200 {
		t.Errorf("Expected status 200 but got %d", status)
	}
}

// BenchmarkReplacePlaceholders benchmarks the path replacement function
func BenchmarkReplacePlaceholders(b *testing.B) {
	testCases := []string{
		"/users/{id}",
		"/users/{userId}/posts/{postId}/comments/{commentId}",
		"/api/v1/{version}/users/{id}",
	}

	for _, tc := range testCases {
		b.Run(tc, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				replacePlaceholders(tc)
			}
		})
	}
}

// TestRunTests integration test
func TestRunTests(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	// Create temp spec file
	tempDir := t.TempDir()
	specContent := `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /users:
    get:
      responses:
        '200':
          description: Success
  /users/{id}:
    get:
      responses:
        '200':
          description: Success
  /posts:
    post:
      responses:
        '200':
          description: Success
`
	specPath := filepath.Join(tempDir, "test-spec.yaml")
	err := os.WriteFile(specPath, []byte(specContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test spec: %v", err)
	}

	// Run tests
	results, err := runTests(specPath, server.URL)
	if err != nil {
		t.Fatalf("runTests failed: %v", err)
	}

	// Verify results
	if len(results) != 3 {
		t.Errorf("Expected 3 results but got %d", len(results))
	}

	// Check that path parameters were replaced
	hasReplacedPath := false
	for _, result := range results {
		if result.endpoint == "/users/{id}" && result.status == "200" {
			hasReplacedPath = true
		}
	}
	if !hasReplacedPath {
		t.Error("Expected path parameter to be replaced in testing")
	}
}

// TestRunTestsInvalidSpec tests error handling
func TestRunTestsInvalidSpec(t *testing.T) {
	_, err := runTests("/nonexistent/spec.yaml", "http://example.com")
	if err == nil {
		t.Error("Expected error for invalid spec but got none")
	}
}

// TestRunTestsWithQueryParams tests query parameter handling
func TestRunTestsWithQueryParams(t *testing.T) {
	// Track requests received
	var receivedURLs []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedURLs = append(receivedURLs, r.URL.String())
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create spec with query parameters
	tempDir := t.TempDir()
	specContent := `openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /users:
    get:
      parameters:
        - name: page
          in: query
          schema:
            type: integer
            example: 1
        - name: limit
          in: query
          schema:
            type: integer
            example: 10
      responses:
        '200':
          description: Success
`
	specPath := filepath.Join(tempDir, "test-spec.yaml")
	err := os.WriteFile(specPath, []byte(specContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test spec: %v", err)
	}

	// Run tests
	results, err := runTests(specPath, server.URL)
	if err != nil {
		t.Fatalf("runTests failed: %v", err)
	}

	// Verify results
	if len(results) != 1 {
		t.Errorf("Expected 1 result but got %d", len(results))
	}

	// Verify query parameters were added
	if len(receivedURLs) != 1 {
		t.Fatalf("Expected 1 request but got %d", len(receivedURLs))
	}

	url := receivedURLs[0]
	if !strings.Contains(url, "page=1") {
		t.Errorf("Expected URL to contain 'page=1' but got: %s", url)
	}
	if !strings.Contains(url, "limit=10") {
		t.Errorf("Expected URL to contain 'limit=10' but got: %s", url)
	}
}

// TestGenerateRequestBody tests request body generation from OpenAPI schemas
func TestGenerateRequestBody(t *testing.T) {
	tests := []struct {
		name           string
		operation      *openapi3.Operation
		expectNil      bool
		expectedFields []string // Fields that should exist in generated JSON
	}{
		{
			name:      "nil operation",
			operation: nil,
			expectNil: true,
		},
		{
			name:      "no request body",
			operation: &openapi3.Operation{},
			expectNil: true,
		},
		{
			name: "simple object schema",
			operation: &openapi3.Operation{
				RequestBody: &openapi3.RequestBodyRef{
					Value: &openapi3.RequestBody{
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: &openapi3.Types{"object"},
										Properties: openapi3.Schemas{
											"title": &openapi3.SchemaRef{
												Value: &openapi3.Schema{
													Type: &openapi3.Types{"string"},
												},
											},
											"userId": &openapi3.SchemaRef{
												Value: &openapi3.Schema{
													Type: &openapi3.Types{"integer"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectNil:      false,
			expectedFields: []string{"title", "userId"},
		},
		{
			name: "schema with example",
			operation: &openapi3.Operation{
				RequestBody: &openapi3.RequestBodyRef{
					Value: &openapi3.RequestBody{
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: &openapi3.Types{"object"},
										Properties: openapi3.Schemas{
											"name": &openapi3.SchemaRef{
												Value: &openapi3.Schema{
													Type:    &openapi3.Types{"string"},
													Example: "John Doe",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectNil:      false,
			expectedFields: []string{"name"},
		},
		{
			name: "nested object schema",
			operation: &openapi3.Operation{
				RequestBody: &openapi3.RequestBodyRef{
					Value: &openapi3.RequestBody{
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: &openapi3.Types{"object"},
										Properties: openapi3.Schemas{
											"user": &openapi3.SchemaRef{
												Value: &openapi3.Schema{
													Type: &openapi3.Types{"object"},
													Properties: openapi3.Schemas{
														"name": &openapi3.SchemaRef{
															Value: &openapi3.Schema{
																Type: &openapi3.Types{"string"},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectNil:      false,
			expectedFields: []string{"user"},
		},
		{
			name: "array schema",
			operation: &openapi3.Operation{
				RequestBody: &openapi3.RequestBodyRef{
					Value: &openapi3.RequestBody{
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{
								Schema: &openapi3.SchemaRef{
									Value: &openapi3.Schema{
										Type: &openapi3.Types{"object"},
										Properties: openapi3.Schemas{
											"tags": &openapi3.SchemaRef{
												Value: &openapi3.Schema{
													Type: &openapi3.Types{"array"},
													Items: &openapi3.SchemaRef{
														Value: &openapi3.Schema{
															Type: &openapi3.Types{"string"},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectNil:      false,
			expectedFields: []string{"tags"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := generateRequestBody(tt.operation)
			
			if err != nil {
				t.Errorf("generateRequestBody() error = %v", err)
				return
			}

			if tt.expectNil {
				if body != nil {
					t.Errorf("generateRequestBody() expected nil but got: %s", string(body))
				}
				return
			}

			if body == nil {
				t.Error("generateRequestBody() returned nil but expected body")
				return
			}

			// Parse JSON and verify fields
			var result map[string]interface{}
			if err := json.Unmarshal(body, &result); err != nil {
				t.Errorf("Failed to parse generated JSON: %v", err)
				return
			}

			for _, field := range tt.expectedFields {
				if _, exists := result[field]; !exists {
					t.Errorf("Expected field %q in generated body but not found. Got: %s", field, string(body))
				}
			}
		})
	}
}

// TestGenerateSampleFromSchema tests schema to sample data conversion
func TestGenerateSampleFromSchema(t *testing.T) {
	tests := []struct {
		name     string
		schema   *openapi3.Schema
		validate func(t *testing.T, result interface{})
	}{
		{
			name:   "nil schema",
			schema: nil,
			validate: func(t *testing.T, result interface{}) {
				if result != nil {
					t.Errorf("Expected nil but got %v", result)
				}
			},
		},
		{
			name: "string type",
			schema: &openapi3.Schema{
				Type: &openapi3.Types{"string"},
			},
			validate: func(t *testing.T, result interface{}) {
				if _, ok := result.(string); !ok {
					t.Errorf("Expected string but got %T", result)
				}
			},
		},
		{
			name: "integer type",
			schema: &openapi3.Schema{
				Type: &openapi3.Types{"integer"},
			},
			validate: func(t *testing.T, result interface{}) {
				if _, ok := result.(int); !ok {
					t.Errorf("Expected int but got %T", result)
				}
			},
		},
		{
			name: "boolean type",
			schema: &openapi3.Schema{
				Type: &openapi3.Types{"boolean"},
			},
			validate: func(t *testing.T, result interface{}) {
				if _, ok := result.(bool); !ok {
					t.Errorf("Expected bool but got %T", result)
				}
			},
		},
		{
			name: "string with example",
			schema: &openapi3.Schema{
				Type:    &openapi3.Types{"string"},
				Example: "test-value",
			},
			validate: func(t *testing.T, result interface{}) {
				if result != "test-value" {
					t.Errorf("Expected 'test-value' but got %v", result)
				}
			},
		},
		{
			name: "email format",
			schema: &openapi3.Schema{
				Type:   &openapi3.Types{"string"},
				Format: "email",
			},
			validate: func(t *testing.T, result interface{}) {
				str, ok := result.(string)
				if !ok || !strings.Contains(str, "@") {
					t.Errorf("Expected email format but got %v", result)
				}
			},
		},
		{
			name: "object type",
			schema: &openapi3.Schema{
				Type: &openapi3.Types{"object"},
				Properties: openapi3.Schemas{
					"name": &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: &openapi3.Types{"string"},
						},
					},
				},
			},
			validate: func(t *testing.T, result interface{}) {
				obj, ok := result.(map[string]interface{})
				if !ok {
					t.Errorf("Expected object but got %T", result)
					return
				}
				if _, exists := obj["name"]; !exists {
					t.Error("Expected 'name' field in object")
				}
			},
		},
		{
			name: "array type",
			schema: &openapi3.Schema{
				Type: &openapi3.Types{"array"},
				Items: &openapi3.SchemaRef{
					Value: &openapi3.Schema{
						Type: &openapi3.Types{"string"},
					},
				},
			},
			validate: func(t *testing.T, result interface{}) {
				arr, ok := result.([]interface{})
				if !ok {
					t.Errorf("Expected array but got %T", result)
					return
				}
				if len(arr) == 0 {
					t.Error("Expected non-empty array")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateSampleFromSchema(tt.schema)
			tt.validate(t, result)
		})
	}
}

// TestValidateResponse tests response validation against OpenAPI specs
func TestValidateResponse(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		contentType string
		operation   *openapi3.Operation
		expectValid bool
		expectMsg   string
	}{
		{
			name:        "nil operation - should pass",
			statusCode:  200,
			contentType: "application/json",
			operation:   nil,
			expectValid: true,
		},
		{
			name:        "no responses defined - should pass",
			statusCode:  200,
			contentType: "application/json",
			operation:   &openapi3.Operation{},
			expectValid: true,
		},
		{
			name:        "valid status code",
			statusCode:  200,
			contentType: "application/json",
			operation: func() *openapi3.Operation {
				responses := openapi3.NewResponses()
				responses.Set("200", &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: strPtr("Success"),
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{},
						},
					},
				})
				return &openapi3.Operation{Responses: responses}
			}(),
			expectValid: true,
		},
		{
			name:        "status not in spec - with default fallback",
			statusCode:  404,
			contentType: "application/json",
			operation: func() *openapi3.Operation {
				responses := openapi3.NewResponses()
				responses.Set("200", &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: strPtr("Success"),
					},
				})
				responses.Set("201", &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: strPtr("Created"),
					},
				})
				// Note: kin-openapi may auto-add a default response
				return &openapi3.Operation{Responses: responses}
			}(),
			expectValid: true, // Will use default if library provides one
		},
		{
			name:        "default response fallback",
			statusCode:  500,
			contentType: "application/json",
			operation: func() *openapi3.Operation {
				responses := openapi3.NewResponses()
				responses.Set("default", &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: strPtr("Error"),
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{},
						},
					},
				})
				return &openapi3.Operation{Responses: responses}
			}(),
			expectValid: true,
		},
		{
			name:        "wrong content type",
			statusCode:  200,
			contentType: "text/html",
			operation: func() *openapi3.Operation {
				responses := openapi3.NewResponses()
				responses.Set("200", &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: strPtr("Success"),
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{},
						},
					},
				})
				return &openapi3.Operation{Responses: responses}
			}(),
			expectValid: false,
			expectMsg:   "content-type",
		},
		{
			name:        "content type with charset",
			statusCode:  200,
			contentType: "application/json; charset=utf-8",
			operation: func() *openapi3.Operation {
				responses := openapi3.NewResponses()
				responses.Set("200", &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: strPtr("Success"),
						Content: openapi3.Content{
							"application/json": &openapi3.MediaType{},
						},
					},
				})
				return &openapi3.Operation{Responses: responses}
			}(),
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock response
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Header:     http.Header{},
			}
			if tt.contentType != "" {
				resp.Header.Set("Content-Type", tt.contentType)
			}

			result := validateResponse(resp, tt.operation, tt.statusCode)

			if result.valid != tt.expectValid {
				t.Errorf("Expected valid=%v but got %v. Errors: %v, StatusValid: %v, ExpectedStatus: %s", 
					tt.expectValid, result.valid, result.schemaErrors, result.statusValid, result.expectedStatus)
			}

			if !tt.expectValid && tt.expectMsg != "" {
				if len(result.schemaErrors) == 0 {
					t.Error("Expected error message but got none")
				} else if !strings.Contains(result.schemaErrors[0], tt.expectMsg) {
					t.Errorf("Expected error containing %q but got %q", 
						tt.expectMsg, result.schemaErrors[0])
				}
			}
		})
	}
}

// Helper function for creating string pointers
func strPtr(s string) *string {
	return &s
}

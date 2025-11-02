package testing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
	"github.com/getkin/kin-openapi/openapi3"
)

// TestReplacePlaceholders tests path parameter replacement
func TestReplacePlaceholders(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"Single parameter", "/users/{id}", "/users/1"},
		{"Multiple parameters", "/users/{userId}/posts/{postId}", "/users/1/posts/1"},
		{"No parameters", "/users", "/users"},
		{"Parameter at start", "{id}/users", "1/users"},
		{"Parameter at end", "/users/{id}", "/users/1"},
		{"Mixed with query", "/users/{id}?page=1", "/users/1?page=1"},
		{"Underscore param", "/items/{item_id}", "/items/1"},
		{"Camel case param", "/items/{itemId}", "/items/1"},
		{"Multiple same name", "/users/{id}/friends/{id}", "/users/1/friends/1"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ReplacePlaceholders(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// TestBuildQueryParams tests query parameter generation
func TestBuildQueryParams(t *testing.T) {
	t.Run("Nil operation", func(t *testing.T) {
		result := BuildQueryParams(nil)
		if result != "" {
			t.Errorf("Expected empty string for nil operation, got: %s", result)
		}
	})

	t.Run("No parameters", func(t *testing.T) {
		operation := &openapi3.Operation{}
		result := BuildQueryParams(operation)
		if result != "" {
			t.Errorf("Expected empty string for no parameters, got: %s", result)
		}
	})

	t.Run("Single string query param", func(t *testing.T) {
		strSchema := openapi3.NewStringSchema()
		operation := &openapi3.Operation{
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:   "name",
						In:     "query",
						Schema: &openapi3.SchemaRef{Value: strSchema},
					},
				},
			},
		}

		result := BuildQueryParams(operation)
		if !strings.Contains(result, "name=test") {
			t.Errorf("Expected result to contain 'name=test', got: %s", result)
		}
		if !strings.HasPrefix(result, "?") {
			t.Errorf("Expected result to start with '?', got: %s", result)
		}
	})

	t.Run("Integer query param", func(t *testing.T) {
		intSchema := openapi3.NewIntegerSchema()
		operation := &openapi3.Operation{
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:   "page",
						In:     "query",
						Schema: &openapi3.SchemaRef{Value: intSchema},
					},
				},
			},
		}

		result := BuildQueryParams(operation)
		if !strings.Contains(result, "page=1") {
			t.Errorf("Expected result to contain 'page=1', got: %s", result)
		}
	})

	t.Run("Boolean query param", func(t *testing.T) {
		boolSchema := openapi3.NewBoolSchema()
		operation := &openapi3.Operation{
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:   "active",
						In:     "query",
						Schema: &openapi3.SchemaRef{Value: boolSchema},
					},
				},
			},
		}

		result := BuildQueryParams(operation)
		if !strings.Contains(result, "active=true") {
			t.Errorf("Expected result to contain 'active=true', got: %s", result)
		}
	})

	t.Run("Multiple query params", func(t *testing.T) {
		strSchema := openapi3.NewStringSchema()
		intSchema := openapi3.NewIntegerSchema()
		
		operation := &openapi3.Operation{
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:   "name",
						In:     "query",
						Schema: &openapi3.SchemaRef{Value: strSchema},
					},
				},
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:   "page",
						In:     "query",
						Schema: &openapi3.SchemaRef{Value: intSchema},
					},
				},
			},
		}

		result := BuildQueryParams(operation)
		if !strings.Contains(result, "name=test") {
			t.Errorf("Expected result to contain 'name=test', got: %s", result)
		}
		if !strings.Contains(result, "page=1") {
			t.Errorf("Expected result to contain 'page=1', got: %s", result)
		}
		if !strings.Contains(result, "&") {
			t.Errorf("Expected result to contain '&', got: %s", result)
		}
	})

	t.Run("Path parameters ignored", func(t *testing.T) {
		strSchema := openapi3.NewStringSchema()
		operation := &openapi3.Operation{
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:   "id",
						In:     "path",
						Schema: &openapi3.SchemaRef{Value: strSchema},
					},
				},
			},
		}

		result := BuildQueryParams(operation)
		if result != "" {
			t.Errorf("Expected empty string for path parameters, got: %s", result)
		}
	})

	t.Run("Enum query param", func(t *testing.T) {
		strSchema := openapi3.NewStringSchema()
		strSchema.Enum = []interface{}{"active", "inactive"}
		
		operation := &openapi3.Operation{
			Parameters: openapi3.Parameters{
				&openapi3.ParameterRef{
					Value: &openapi3.Parameter{
						Name:   "status",
						In:     "query",
						Schema: &openapi3.SchemaRef{Value: strSchema},
					},
				},
			},
		}

		result := BuildQueryParams(operation)
		if !strings.Contains(result, "status=active") {
			t.Errorf("Expected result to contain 'status=active', got: %s", result)
		}
	})
}

// TestGenerateRequestBody tests request body generation
func TestGenerateRequestBody(t *testing.T) {
	t.Run("Nil operation", func(t *testing.T) {
		body, err := GenerateRequestBody(nil)
		if err != nil {
			t.Errorf("Expected no error for nil operation, got: %v", err)
		}
		if body != nil {
			t.Errorf("Expected nil body for nil operation, got: %v", body)
		}
	})

	t.Run("No request body", func(t *testing.T) {
		operation := &openapi3.Operation{}
		body, err := GenerateRequestBody(operation)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if body != nil {
			t.Errorf("Expected nil body, got: %v", body)
		}
	})

	t.Run("Simple object request body", func(t *testing.T) {
		schema := openapi3.NewObjectSchema()
		schema.Properties = openapi3.Schemas{
			"name": &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
			"age":  &openapi3.SchemaRef{Value: openapi3.NewIntegerSchema()},
		}

		operation := &openapi3.Operation{
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{Value: schema},
						},
					},
				},
			},
		}

		body, err := GenerateRequestBody(operation)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if body == nil {
			t.Fatal("Expected body to be generated")
		}

		// Verify it's valid JSON
		var jsonData map[string]interface{}
		if err := json.Unmarshal(body, &jsonData); err != nil {
			t.Errorf("Expected valid JSON, got: %v", err)
		}

		// Verify fields exist
		if _, ok := jsonData["name"]; !ok {
			t.Error("Expected 'name' field in generated body")
		}
		if _, ok := jsonData["age"]; !ok {
			t.Error("Expected 'age' field in generated body")
		}
	})

	t.Run("Request body with example", func(t *testing.T) {
		schema := openapi3.NewStringSchema()
		schema.Example = "test@example.com"

		operation := &openapi3.Operation{
			RequestBody: &openapi3.RequestBodyRef{
				Value: &openapi3.RequestBody{
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{Value: schema},
						},
					},
				},
			},
		}

		body, err := GenerateRequestBody(operation)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Verify the example is used
		if !strings.Contains(string(body), "test@example.com") {
			t.Errorf("Expected body to contain example value, got: %s", string(body))
		}
	})
}

// TestGenerateSampleFromSchema tests schema-based sample generation
func TestGenerateSampleFromSchema(t *testing.T) {
	t.Run("Nil schema", func(t *testing.T) {
		result := GenerateSampleFromSchema(nil)
		if result != nil {
			t.Errorf("Expected nil for nil schema, got: %v", result)
		}
	})

	t.Run("String schema", func(t *testing.T) {
		schema := openapi3.NewStringSchema()
		result := GenerateSampleFromSchema(schema)
		if result != "sample" {
			t.Errorf("Expected 'sample', got: %v", result)
		}
	})

	t.Run("String with enum", func(t *testing.T) {
		schema := openapi3.NewStringSchema()
		schema.Enum = []interface{}{"red", "green", "blue"}
		result := GenerateSampleFromSchema(schema)
		if result != "red" {
			t.Errorf("Expected 'red', got: %v", result)
		}
	})

	t.Run("Email format", func(t *testing.T) {
		schema := openapi3.NewStringSchema()
		schema.Format = "email"
		result := GenerateSampleFromSchema(schema)
		if result != "user@example.com" {
			t.Errorf("Expected 'user@example.com', got: %v", result)
		}
	})

	t.Run("URL format", func(t *testing.T) {
		schema := openapi3.NewStringSchema()
		schema.Format = "url"
		result := GenerateSampleFromSchema(schema)
		if result != "https://example.com" {
			t.Errorf("Expected 'https://example.com', got: %v", result)
		}
	})

	t.Run("Date format", func(t *testing.T) {
		schema := openapi3.NewStringSchema()
		schema.Format = "date"
		result := GenerateSampleFromSchema(schema)
		if result != "2024-01-01" {
			t.Errorf("Expected '2024-01-01', got: %v", result)
		}
	})

	t.Run("DateTime format", func(t *testing.T) {
		schema := openapi3.NewStringSchema()
		schema.Format = "date-time"
		result := GenerateSampleFromSchema(schema)
		if result != "2024-01-01T00:00:00Z" {
			t.Errorf("Expected '2024-01-01T00:00:00Z', got: %v", result)
		}
	})

	t.Run("Integer schema", func(t *testing.T) {
		schema := openapi3.NewIntegerSchema()
		result := GenerateSampleFromSchema(schema)
		if result != 1 {
			t.Errorf("Expected 1, got: %v", result)
		}
	})

	t.Run("Number schema", func(t *testing.T) {
		schema := openapi3.NewFloat64Schema()
		result := GenerateSampleFromSchema(schema)
		if result != 1.0 {
			t.Errorf("Expected 1.0, got: %v", result)
		}
	})

	t.Run("Boolean schema", func(t *testing.T) {
		schema := openapi3.NewBoolSchema()
		result := GenerateSampleFromSchema(schema)
		if result != true {
			t.Errorf("Expected true, got: %v", result)
		}
	})

	t.Run("Object schema", func(t *testing.T) {
		schema := openapi3.NewObjectSchema()
		schema.Properties = openapi3.Schemas{
			"name": &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
			"age":  &openapi3.SchemaRef{Value: openapi3.NewIntegerSchema()},
		}

		result := GenerateSampleFromSchema(schema)
		obj, ok := result.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected map[string]interface{}, got: %T", result)
		}

		if _, ok := obj["name"]; !ok {
			t.Error("Expected 'name' field in object")
		}
		if _, ok := obj["age"]; !ok {
			t.Error("Expected 'age' field in object")
		}
	})

	t.Run("Array schema", func(t *testing.T) {
		schema := openapi3.NewArraySchema()
		schema.Items = &openapi3.SchemaRef{Value: openapi3.NewStringSchema()}

		result := GenerateSampleFromSchema(schema)
		arr, ok := result.([]interface{})
		if !ok {
			t.Fatalf("Expected []interface{}, got: %T", result)
		}

		if len(arr) != 1 {
			t.Errorf("Expected array with 1 item, got: %d", len(arr))
		}
	})

	t.Run("Schema with example", func(t *testing.T) {
		schema := openapi3.NewStringSchema()
		schema.Example = "custom-example"
		result := GenerateSampleFromSchema(schema)
		if result != "custom-example" {
			t.Errorf("Expected 'custom-example', got: %v", result)
		}
	})

	t.Run("Schema with default", func(t *testing.T) {
		schema := openapi3.NewStringSchema()
		schema.Default = "default-value"
		result := GenerateSampleFromSchema(schema)
		if result != "default-value" {
			t.Errorf("Expected 'default-value', got: %v", result)
		}
	})
}

// TestApplyAuth tests authentication application
func TestApplyAuth(t *testing.T) {
	t.Run("Nil auth", func(t *testing.T) {
		req := httptest.NewRequest("GET", "https://api.example.com/users", nil)
		ApplyAuth(req, nil)
		// Should not panic or modify request
		if req.Header.Get("Authorization") != "" {
			t.Error("Expected no Authorization header")
		}
	})

	t.Run("No auth type", func(t *testing.T) {
		req := httptest.NewRequest("GET", "https://api.example.com/users", nil)
		auth := &models.AuthConfig{AuthType: "none"}
		ApplyAuth(req, auth)
		if req.Header.Get("Authorization") != "" {
			t.Error("Expected no Authorization header")
		}
	})

	t.Run("Bearer token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "https://api.example.com/users", nil)
		auth := &models.AuthConfig{
			AuthType: "bearer",
			Token:    "test-token-123",
		}
		ApplyAuth(req, auth)
		
		authHeader := req.Header.Get("Authorization")
		if authHeader != "Bearer test-token-123" {
			t.Errorf("Expected 'Bearer test-token-123', got: %s", authHeader)
		}
	})

	t.Run("API key in header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "https://api.example.com/users", nil)
		auth := &models.AuthConfig{
			AuthType:   "apiKey",
			APIKeyName: "X-API-Key",
			APIKeyIn:   "header",
			Token:      "api-key-456",
		}
		ApplyAuth(req, auth)
		
		apiKey := req.Header.Get("X-API-Key")
		if apiKey != "api-key-456" {
			t.Errorf("Expected 'api-key-456', got: %s", apiKey)
		}
	})

	t.Run("API key in query", func(t *testing.T) {
		req := httptest.NewRequest("GET", "https://api.example.com/users", nil)
		auth := &models.AuthConfig{
			AuthType:   "apiKey",
			APIKeyName: "api_key",
			APIKeyIn:   "query",
			Token:      "query-key-789",
		}
		ApplyAuth(req, auth)
		
		queryParam := req.URL.Query().Get("api_key")
		if queryParam != "query-key-789" {
			t.Errorf("Expected 'query-key-789', got: %s", queryParam)
		}
	})

	t.Run("Basic auth", func(t *testing.T) {
		req := httptest.NewRequest("GET", "https://api.example.com/users", nil)
		auth := &models.AuthConfig{
			AuthType: "basic",
			Username: "testuser",
			Password: "testpass",
		}
		ApplyAuth(req, auth)
		
		username, password, ok := req.BasicAuth()
		if !ok {
			t.Error("Expected basic auth to be set")
		}
		if username != "testuser" {
			t.Errorf("Expected username 'testuser', got: %s", username)
		}
		if password != "testpass" {
			t.Errorf("Expected password 'testpass', got: %s", password)
		}
	})
}

// TestTestEndpoint tests endpoint testing with mock server
func TestTestEndpoint(t *testing.T) {
	t.Run("Successful GET request", func(t *testing.T) {
		// Create mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"message":"success"}`))
		}))
		defer server.Close()

		statusCode, resp, logEntry, err := TestEndpoint("GET", server.URL, nil, nil, false)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if statusCode != http.StatusOK {
			t.Errorf("Expected status 200, got: %d", statusCode)
		}
		if resp == nil {
			t.Error("Expected response to be non-nil")
		}
		if logEntry != nil {
			t.Error("Expected no log entry when verbose=false")
		}
	})

	t.Run("POST request with body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				t.Errorf("Expected POST, got: %s", r.Method)
			}
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type application/json, got: %s", r.Header.Get("Content-Type"))
			}
			w.WriteHeader(http.StatusCreated)
		}))
		defer server.Close()

		body := []byte(`{"name":"test"}`)
		statusCode, _, _, err := TestEndpoint("POST", server.URL, body, nil, false)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if statusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got: %d", statusCode)
		}
	})

	t.Run("Verbose mode creates log entry", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Response-ID", "12345")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"data":"test"}`))
		}))
		defer server.Close()

		_, _, logEntry, err := TestEndpoint("GET", server.URL, nil, nil, true)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if logEntry == nil {
			t.Fatal("Expected log entry when verbose=true")
		}
		if logEntry.RequestURL != server.URL {
			t.Errorf("Expected URL %s, got: %s", server.URL, logEntry.RequestURL)
		}
		if logEntry.Duration == 0 {
			t.Error("Expected duration to be recorded")
		}
		if len(logEntry.ResponseHeaders) == 0 {
			t.Error("Expected response headers to be captured")
		}
	})

	t.Run("With bearer auth", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "Bearer test-token" {
				t.Errorf("Expected 'Bearer test-token', got: %s", authHeader)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		auth := &models.AuthConfig{
			AuthType: "bearer",
			Token:    "test-token",
		}

		statusCode, _, _, err := TestEndpoint("GET", server.URL, nil, auth, false)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if statusCode != http.StatusOK {
			t.Errorf("Expected status 200, got: %d", statusCode)
		}
	})

	t.Run("Invalid URL returns error", func(t *testing.T) {
		_, _, _, err := TestEndpoint("GET", "http://invalid-domain-that-does-not-exist-12345.com", nil, nil, false)
		if err == nil {
			t.Error("Expected error for invalid URL")
		}
	})
}

// TestRunTestCmd tests the Bubble Tea command wrapper
func TestRunTestCmd(t *testing.T) {
	t.Run("Invalid spec path returns error message", func(t *testing.T) {
		cmd := RunTestCmd("/nonexistent/spec.yaml", "https://api.example.com", nil, false)
		msg := cmd()

		errMsg, ok := msg.(TestErrorMsg)
		if !ok {
			t.Fatalf("Expected TestErrorMsg, got: %T", msg)
		}
		if errMsg.Err == nil {
			t.Error("Expected error in TestErrorMsg")
		}
	})
}

// TestBuildQueryParams_WithExample tests example values in query params
func TestBuildQueryParams_WithExample(t *testing.T) {
	strSchema := openapi3.NewStringSchema()
	strSchema.Example = "custom-value"
	
	operation := &openapi3.Operation{
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:   "query",
					In:     "query",
					Schema: &openapi3.SchemaRef{Value: strSchema},
				},
			},
		},
	}

	result := BuildQueryParams(operation)
	if !strings.Contains(result, "query=custom-value") {
		t.Errorf("Expected result to contain 'query=custom-value', got: %s", result)
	}
}

// TestGenerateSampleFromSchema_NestedObject tests nested object generation
func TestGenerateSampleFromSchema_NestedObject(t *testing.T) {
	addressSchema := openapi3.NewObjectSchema()
	addressSchema.Properties = openapi3.Schemas{
		"street": &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
		"city":   &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
	}

	userSchema := openapi3.NewObjectSchema()
	userSchema.Properties = openapi3.Schemas{
		"name":    &openapi3.SchemaRef{Value: openapi3.NewStringSchema()},
		"address": &openapi3.SchemaRef{Value: addressSchema},
	}

	result := GenerateSampleFromSchema(userSchema)
	obj, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got: %T", result)
	}

	address, ok := obj["address"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected nested address object")
	}

	if _, ok := address["street"]; !ok {
		t.Error("Expected 'street' field in nested address")
	}
	if _, ok := address["city"]; !ok {
		t.Error("Expected 'city' field in nested address")
	}
}

// TestTestEndpoint_Timeout tests that timeout is enforced
func TestTestEndpoint_Timeout(t *testing.T) {
	// Create server that delays longer than timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(15 * time.Second) // Longer than 10s timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	startTime := time.Now()
	_, _, _, err := TestEndpoint("GET", server.URL, nil, nil, false)
	duration := time.Since(startTime)

	if err == nil {
		t.Error("Expected timeout error")
	}

	// Should timeout around 10 seconds, not wait 15
	if duration > 12*time.Second {
		t.Errorf("Expected timeout around 10s, took: %v", duration)
	}
}

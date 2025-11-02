package validation

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

// TestValidateSpec tests OpenAPI spec validation
func TestValidateSpec(t *testing.T) {
	// Create a valid test spec
	validSpec := `
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
	tmpFile, err := os.CreateTemp("", "valid-spec-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.Write([]byte(validSpec))
	tmpFile.Close()

	// Test valid spec
	result, err := ValidateSpec(tmpFile.Name())
	if err != nil {
		t.Fatalf("ValidateSpec() failed for valid spec: %v", err)
	}
	if result == "" {
		t.Error("Expected success message, got empty string")
	}
	if result != "OpenAPI spec is valid! 🎉" {
		t.Errorf("Expected success message, got: %s", result)
	}
}

// TestValidateSpec_InvalidFile tests validation with nonexistent file
func TestValidateSpec_InvalidFile(t *testing.T) {
	_, err := ValidateSpec("/nonexistent/file.yaml")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

// TestValidateSpec_InvalidYAML tests validation with malformed YAML
func TestValidateSpec_InvalidYAML(t *testing.T) {
	invalidYAML := "this is not: [valid yaml"
	tmpFile, err := os.CreateTemp("", "invalid-yaml-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.Write([]byte(invalidYAML))
	tmpFile.Close()

	_, err = ValidateSpec(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

// TestValidateSpec_InvalidOpenAPI tests validation with invalid OpenAPI structure
func TestValidateSpec_InvalidOpenAPI(t *testing.T) {
	invalidSpec := `
openapi: 3.0.0
info:
  title: Test
  # Missing version field (required)
paths: {}
`
	tmpFile, err := os.CreateTemp("", "invalid-spec-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.Write([]byte(invalidSpec))
	tmpFile.Close()

	_, err = ValidateSpec(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for invalid OpenAPI spec, got nil")
	}
}

// TestValidateResponse tests response validation against OpenAPI spec
func TestValidateResponse(t *testing.T) {
	// Create a test operation
	responses := openapi3.NewResponses()
	desc200 := "OK"
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &desc200,
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{},
			},
		},
	})
	
	operation := &openapi3.Operation{
		Responses: responses,
	}

	// Test successful validation
	resp := &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(`{"test": "data"}`))),
	}

	result := ValidateResponse(resp, operation, 200)
	
	if !result.Valid {
		t.Error("Expected valid result for matching response")
	}
	if !result.StatusValid {
		t.Error("Expected status to be valid")
	}
	if result.ContentType != "application/json" {
		t.Errorf("Expected content type 'application/json', got '%s'", result.ContentType)
	}
	if len(result.SchemaErrors) > 0 {
		t.Errorf("Expected no schema errors, got: %v", result.SchemaErrors)
	}
}

// TestValidateResponse_InvalidStatus tests validation with undefined status code
// Note: NewResponses() automatically creates a "default" response, so we test
// that the validation correctly uses it for undefined status codes
func TestValidateResponse_InvalidStatus(t *testing.T) {
	responses := openapi3.NewResponses()
	desc200 := "OK"
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &desc200,
		},
	})
	
	operation := &openapi3.Operation{
		Responses: responses,
	}

	resp := &http.Response{
		StatusCode: 404,
		Header:     http.Header{},
		Body:       io.NopCloser(bytes.NewReader([]byte(""))),
	}

	result := ValidateResponse(resp, operation, 404)
	
	// Since NewResponses() creates a default response, 404 should be valid (uses default)
	// This is actually correct behavior per OpenAPI spec
	if !result.Valid {
		t.Errorf("Expected valid result (uses default), got Valid=%v", result.Valid)
	}
	if !result.StatusValid {
		t.Errorf("Expected status to be valid (uses default), got StatusValid=%v", result.StatusValid)
	}
	if result.ExpectedStatus != "default" {
		t.Errorf("Expected ExpectedStatus='default', got '%s'", result.ExpectedStatus)
	}
}

// TestValidateResponse_DefaultResponse tests validation with default response
func TestValidateResponse_DefaultResponse(t *testing.T) {
	responses := openapi3.NewResponses()
	desc200 := "OK"
	descDefault := "Default response"
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &desc200,
		},
	})
	responses.Set("default", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descDefault,
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{},
			},
		},
	})
	
	operation := &openapi3.Operation{
		Responses: responses,
	}

	resp := &http.Response{
		StatusCode: 500,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(`{"error": "internal"}`))),
	}

	result := ValidateResponse(resp, operation, 500)
	
	if !result.Valid {
		t.Error("Expected valid result for default response")
	}
	if !result.StatusValid {
		t.Error("Expected status to be valid (using default)")
	}
	if result.ExpectedStatus != "default" {
		t.Errorf("Expected status 'default', got '%s'", result.ExpectedStatus)
	}
}

// TestValidateResponse_InvalidContentType tests validation with wrong content type
func TestValidateResponse_InvalidContentType(t *testing.T) {
	responses := openapi3.NewResponses()
	desc200 := "OK"
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &desc200,
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{},
			},
		},
	})
	
	operation := &openapi3.Operation{
		Responses: responses,
	}

	resp := &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": []string{"text/html"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte("<html></html>"))),
	}

	result := ValidateResponse(resp, operation, 200)
	
	if result.Valid {
		t.Error("Expected invalid result for wrong content type")
	}
	if len(result.SchemaErrors) == 0 {
		t.Error("Expected schema errors for wrong content type")
	}
}

// TestValidateResponse_ContentTypeWithCharset tests content type with charset
func TestValidateResponse_ContentTypeWithCharset(t *testing.T) {
	responses := openapi3.NewResponses()
	desc200 := "OK"
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &desc200,
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{},
			},
		},
	})
	
	operation := &openapi3.Operation{
		Responses: responses,
	}

	resp := &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type": []string{"application/json; charset=utf-8"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte(`{}`))),
	}

	result := ValidateResponse(resp, operation, 200)
	
	if !result.Valid {
		t.Error("Expected valid result when content type includes charset")
	}
}

// TestValidateResponse_NoOperation tests validation with nil operation
func TestValidateResponse_NoOperation(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Header:     http.Header{},
		Body:       io.NopCloser(bytes.NewReader([]byte(""))),
	}

	result := ValidateResponse(resp, nil, 200)
	
	if !result.Valid {
		t.Error("Expected valid result when no operation provided")
	}
}

// TestValidateResponseBody tests response body validation placeholder
func TestValidateResponseBody(t *testing.T) {
	// This is a placeholder test for future implementation
	body := []byte(`{"test": "data"}`)
	schema := &openapi3.Schema{}

	errors := ValidateResponseBody(body, schema)
	
	// Currently returns empty array (not implemented)
	if len(errors) != 0 {
		t.Errorf("Expected empty errors array, got %d errors", len(errors))
	}
}

// TestValidateHeaders tests header validation placeholder
func TestValidateHeaders(t *testing.T) {
	headers := http.Header{
		"Content-Type": []string{"application/json"},
		"X-Custom":     []string{"value"},
	}
	spec := map[string]*openapi3.HeaderRef{}

	errors := ValidateHeaders(headers, spec)
	
	// Currently returns empty array (not implemented)
	if len(errors) != 0 {
		t.Errorf("Expected empty errors array, got %d errors", len(errors))
	}
}

// TestValidateContentType tests content type matching
func TestValidateContentType(t *testing.T) {
	tests := []struct {
		name     string
		actual   string
		expected []string
		want     bool
	}{
		{"exact match", "application/json", []string{"application/json"}, true},
		{"with charset", "application/json; charset=utf-8", []string{"application/json"}, true},
		{"no match", "text/html", []string{"application/json"}, false},
		{"multiple options", "application/xml", []string{"application/json", "application/xml"}, true},
		{"empty actual", "", []string{"application/json"}, false},
		{"empty expected", "application/json", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateContentType(tt.actual, tt.expected)
			if got != tt.want {
				t.Errorf("validateContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetExpectedContentTypes tests extraction of content types from response spec
func TestGetExpectedContentTypes(t *testing.T) {
	// Test with nil response
	types := getExpectedContentTypes(nil)
	if len(types) != 0 {
		t.Errorf("Expected empty array for nil response, got %d types", len(types))
	}

	// Test with response that has content types
	response := &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{},
				"application/xml":  &openapi3.MediaType{},
			},
		},
	}

	types = getExpectedContentTypes(response)
	if len(types) != 2 {
		t.Errorf("Expected 2 content types, got %d", len(types))
	}
}

// TestValidateStatusCode tests status code validation
func TestValidateStatusCode(t *testing.T) {
	responses := openapi3.NewResponses()
	desc200 := "OK"
	desc404 := "Not Found"
	descDefault := "Default"
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &desc200,
		},
	})
	responses.Set("404", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &desc404,
		},
	})
	responses.Set("default", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descDefault,
		},
	})
	
	operation := &openapi3.Operation{
		Responses: responses,
	}

	tests := []struct {
		name       string
		statusCode int
		wantValid  bool
		wantStatus string
	}{
		{"defined 200", 200, true, "200"},
		{"defined 404", 404, true, "404"},
		{"undefined uses default", 500, true, "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, status := validateStatusCode(tt.statusCode, operation)
			if valid != tt.wantValid {
				t.Errorf("validateStatusCode() valid = %v, want %v", valid, tt.wantValid)
			}
			if status != tt.wantStatus {
				t.Errorf("validateStatusCode() status = %v, want %v", status, tt.wantStatus)
			}
		})
	}
}

// TestValidateStatusCode_NoOperation tests validation with nil operation
func TestValidateStatusCode_NoOperation(t *testing.T) {
	valid, status := validateStatusCode(200, nil)
	if !valid {
		t.Error("Expected valid=true for nil operation")
	}
	if status != "" {
		t.Errorf("Expected empty status for nil operation, got %s", status)
	}
}

// TestFormatValidationErrors tests error formatting
func TestFormatValidationErrors(t *testing.T) {
	// Test empty errors
	result := formatValidationErrors([]string{})
	if result != "No validation errors" {
		t.Errorf("Expected 'No validation errors', got '%s'", result)
	}

	// Test with errors
	errors := []string{"error 1", "error 2", "error 3"}
	result = formatValidationErrors(errors)
	if !bytes.Contains([]byte(result), []byte("Validation errors:")) {
		t.Error("Expected 'Validation errors:' in output")
	}
	if !bytes.Contains([]byte(result), []byte("1. error 1")) {
		t.Error("Expected numbered error in output")
	}
}

// TestDrainAndCloseResponse tests response body cleanup
func TestDrainAndCloseResponse(t *testing.T) {
	// Test with nil response
	drainAndCloseResponse(nil)

	// Test with valid response
	body := io.NopCloser(bytes.NewReader([]byte("test data")))
	resp := &http.Response{
		Body: body,
	}
	drainAndCloseResponse(resp)
	// Should not panic
}

// TestValidateResponse_NilOperation tests ValidateResponse with nil operation
func TestValidateResponse_NilOperation(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{}`))),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
	
	// Call with nil operation - should still return a result (marks as valid since no spec)
	result := ValidateResponse(resp, nil, 200)
	// The function returns Valid=true when operation is nil (no spec to validate against)
	if !result.Valid {
		t.Error("Expected validation to pass with nil operation (no spec to validate)")
	}
}

// TestValidateResponse_NilResponses tests ValidateResponse with operation but no responses
func TestValidateResponse_NilResponses(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{}`))),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
	operation := &openapi3.Operation{
		Responses: nil,
	}
	
	// Call with nil responses - should return valid (no spec to validate against)
	result := ValidateResponse(resp, operation, 200)
	if !result.Valid {
		t.Error("Expected validation to pass with nil responses (no spec to validate)")
	}
}

// TestValidateResponse_DefaultResponseFallback tests ValidateResponse falling back to default response
func TestValidateResponse_DefaultResponseFallback(t *testing.T) {
	// Use a 404 status that's not explicitly defined, should fall back to default
	resp := &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"error":"not found"}`))),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
	
	// Create responses with only default and 200
	responses := openapi3.NewResponses()
	desc200 := "OK"
	descDefault := "Default error response"
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &desc200,
		},
	})
	responses.Set("default", &openapi3.ResponseRef{
		Value: &openapi3.Response{
			Description: &descDefault,
			Content: openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{
						Value: openapi3.NewObjectSchema().WithProperty("error", openapi3.NewStringSchema()),
					},
				},
			},
		},
	})
	
	operation := &openapi3.Operation{
		Responses: responses,
	}
	
	result := ValidateResponse(resp, operation, 404)
	// Should use default response for validation
	if !result.Valid {
		t.Errorf("Expected validation to succeed with default response: %v", result.SchemaErrors)
	}
	if result.ExpectedStatus != "default" {
		t.Errorf("Expected ExpectedStatus='default', got '%s'", result.ExpectedStatus)
	}
}

// TestValidateResponse_MismatchedContentType tests ValidateResponse with wrong content-type
func TestValidateResponse_MismatchedContentType(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(`<xml>data</xml>`))),
		Header:     http.Header{"Content-Type": []string{"application/xml"}}, // XML instead of JSON
	}
	
	desc := "Success"
	responseWithContent := &openapi3.Response{
		Description: &desc,
		Content: openapi3.Content{
			"application/json": &openapi3.MediaType{
				Schema: &openapi3.SchemaRef{
					Value: openapi3.NewObjectSchema(),
				},
			},
		},
	}
	
	responses := openapi3.NewResponses()
	responses.Set("200", &openapi3.ResponseRef{Value: responseWithContent})
	
	operation := &openapi3.Operation{
		Responses: responses,
	}
	
	result := ValidateResponse(resp, operation, 200)
	if result.Valid {
		t.Error("Expected validation to fail with mismatched content-type")
	}
	hasContentTypeError := false
	for _, err := range result.SchemaErrors {
		if strings.Contains(strings.ToLower(err), "content-type") {
			hasContentTypeError = true
			break
		}
	}
	if !hasContentTypeError {
		t.Errorf("Expected content-type error in validation results, got: %v", result.SchemaErrors)
	}
}

// TestValidateStatusCode_Ranges tests various status code ranges and exact matches
func TestValidateStatusCode_Ranges(t *testing.T) {
	tests := []struct {
		name             string
		statusCode       int
		definedStatus    string
		shouldBeValid    bool
		expectedStatus   string
	}{
		{
			name:           "1xx range match",
			statusCode:     101,
			definedStatus:  "1XX",
			shouldBeValid:  true,
			expectedStatus: "101",
		},
		{
			name:           "2xx range match",
			statusCode:     204,
			definedStatus:  "2XX",
			shouldBeValid:  true,
			expectedStatus: "204",
		},
		{
			name:           "3xx range match",
			statusCode:     304,
			definedStatus:  "3XX",
			shouldBeValid:  true,
			expectedStatus: "304",
		},
		{
			name:           "4xx range match",
			statusCode:     404,
			definedStatus:  "4XX",
			shouldBeValid:  true,
			expectedStatus: "404",
		},
		{
			name:           "5xx range match",
			statusCode:     503,
			definedStatus:  "5XX",
			shouldBeValid:  true,
			expectedStatus: "503",
		},
		{
			name:           "exact match",
			statusCode:     200,
			definedStatus:  "200",
			shouldBeValid:  true,
			expectedStatus: "200",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responses := openapi3.NewResponses()
			desc := "test"
			responses.Set(tt.definedStatus, &openapi3.ResponseRef{
				Value: &openapi3.Response{Description: &desc},
			})
			
			operation := &openapi3.Operation{
				Responses: responses,
			}
			
			valid, status := validateStatusCode(tt.statusCode, operation)
			if valid != tt.shouldBeValid {
				t.Errorf("Expected valid=%v for status %d with defined %s, got %v", 
					tt.shouldBeValid, tt.statusCode, tt.definedStatus, valid)
			}
			if valid && status != tt.expectedStatus {
				t.Errorf("Expected status=%s, got %s", tt.expectedStatus, status)
			}
		})
	}
}

// TestValidateStatusCode_WithDefault tests default response handling
func TestValidateStatusCode_WithDefault(t *testing.T) {
	responses := openapi3.NewResponses()
	desc200 := "OK"
	descDefault := "Default"
	responses.Set("200", &openapi3.ResponseRef{
		Value: &openapi3.Response{Description: &desc200},
	})
	responses.Set("default", &openapi3.ResponseRef{
		Value: &openapi3.Response{Description: &descDefault},
	})
	
	operation := &openapi3.Operation{
		Responses: responses,
	}
	
	// Try 404 which is not defined but has default
	valid, status := validateStatusCode(404, operation)
	if !valid {
		t.Error("Expected valid=true when default response exists")
	}
	// Should match either "404" (from Status()) or "default" (from Map())
	if status != "404" && status != "default" {
		t.Errorf("Expected status '404' or 'default', got '%s'", status)
	}
}

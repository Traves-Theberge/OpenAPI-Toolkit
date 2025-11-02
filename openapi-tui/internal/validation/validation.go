package validation

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/errors"
)

// validateSpec validates an OpenAPI specification file
// Returns success message or error with helpful suggestions
func ValidateSpec(filePath string) (string, error) {
	// Load OpenAPI document with external references allowed
	loader := &openapi3.Loader{IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(filePath)
	if err != nil {
		return "", errors.EnhanceFileError(err, filePath)
	}

	// Validate the loaded document
	err = doc.Validate(loader.Context)
	if err != nil {
		return "", errors.EnhanceValidationError(err)
	}

	return "OpenAPI spec is valid! 🎉", nil
}

// validateResponse validates an HTTP response against OpenAPI spec
// Returns validation result with detailed error information
func ValidateResponse(resp *http.Response, operation *openapi3.Operation, statusCode int) models.ValidationResult {
	result := models.ValidationResult{
		Valid:       true,
		StatusValid: false,
		ContentType: resp.Header.Get("Content-Type"),
	}

	if operation == nil || operation.Responses == nil {
		// No spec to validate against - mark as valid
		return result
	}

	// Check if status code is defined in spec
	statusStr := fmt.Sprintf("%d", statusCode)
	response := operation.Responses.Status(statusCode)
	
	if response != nil {
		// Found exact status match
		result.StatusValid = true
		result.ExpectedStatus = statusStr
	} else {
		// Check if there's an explicit "default" response
		respMap := operation.Responses.Map()
		if defaultResp, hasDefault := respMap["default"]; hasDefault && defaultResp != nil {
			response = defaultResp
			result.StatusValid = true
			result.ExpectedStatus = "default"
		} else {
			// Status code not defined in spec and no explicit default
			result.StatusValid = false
			result.Valid = false
			result.SchemaErrors = append(result.SchemaErrors, 
				fmt.Sprintf("status %d not defined in spec", statusCode))
			return result
		}
	}

	// Validate content type if response has content
	if response.Value != nil && response.Value.Content != nil {
		// Extract base content type (ignore charset, etc.)
		contentType := strings.Split(result.ContentType, ";")[0]
		contentType = strings.TrimSpace(contentType)
		
		// Check if content type is defined in spec
		mediaType := response.Value.Content.Get(contentType)
		if mediaType == nil {
			// Try common alternatives
			if contentType == "" {
				contentType = "application/json" // Default assumption
				mediaType = response.Value.Content.Get(contentType)
			}
		}

		if mediaType == nil && len(response.Value.Content) > 0 {
			result.Valid = false
			result.SchemaErrors = append(result.SchemaErrors,
				fmt.Sprintf("content-type '%s' not defined in spec", result.ContentType))
		}

		// TODO: Add JSON schema validation against response body
		// This would require parsing the response body and validating against mediaType.Schema
		// For now, we just validate status code and content type
	}

	return result
}

// validateResponseBody validates response body against OpenAPI schema
// This is a placeholder for future implementation
func ValidateResponseBody(body []byte, schema *openapi3.Schema) []string {
	// TODO: Implement full JSON schema validation
	// For now, just return empty array (no errors)
	return []string{}
}

// validateHeaders validates response headers against OpenAPI spec
// This is a placeholder for future implementation
func ValidateHeaders(headers http.Header, spec map[string]*openapi3.HeaderRef) []string {
	var errors []string
	// TODO: Implement header validation
	// Check required headers, types, patterns, etc.
	return errors
}

// validateContentType checks if response content type matches spec
func validateContentType(actual string, expected []string) bool {
	// Normalize actual content type (remove charset, etc.)
	actual = strings.Split(actual, ";")[0]
	actual = strings.TrimSpace(actual)
	
	for _, exp := range expected {
		exp = strings.TrimSpace(exp)
		if actual == exp {
			return true
		}
	}
	return false
}

// getExpectedContentTypes extracts expected content types from response spec
func getExpectedContentTypes(response *openapi3.ResponseRef) []string {
	var contentTypes []string
	if response != nil && response.Value != nil && response.Value.Content != nil {
		for ct := range response.Value.Content {
			contentTypes = append(contentTypes, ct)
		}
	}
	return contentTypes
}

// validateStatusCode checks if status code is defined in operation responses
func validateStatusCode(statusCode int, operation *openapi3.Operation) (bool, string) {
	if operation == nil || operation.Responses == nil {
		return true, "" // No spec to validate against
	}

	statusStr := fmt.Sprintf("%d", statusCode)
	response := operation.Responses.Status(statusCode)
	
	if response != nil {
		return true, statusStr
	}

	// Check for default response
	respMap := operation.Responses.Map()
	if defaultResp, hasDefault := respMap["default"]; hasDefault && defaultResp != nil {
		return true, "default"
	}

	return false, ""
}

// formatValidationErrors formats validation errors for display
func formatValidationErrors(errors []string) string {
	if len(errors) == 0 {
		return "No validation errors"
	}
	
	var sb strings.Builder
	sb.WriteString("Validation errors:\n")
	for i, err := range errors {
		sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, err))
	}
	return sb.String()
}

// drainAndCloseResponse drains and closes an HTTP response body
// This prevents connection reuse issues
func drainAndCloseResponse(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

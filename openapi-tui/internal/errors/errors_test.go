package errors

import (
	"errors"
	"strings"
	"testing"
)

// TestEnhancedError_Error tests the Error() method implementation
func TestEnhancedError_Error(t *testing.T) {
	t.Run("Error with suggestions", func(t *testing.T) {
		err := &EnhancedError{
			Title:       "Test Error",
			Description: "This is a test error",
			Suggestions: []string{"Try this", "Or try that"},
			Original:    errors.New("original error"),
		}

		errStr := err.Error()
		if !strings.Contains(errStr, "Test Error") {
			t.Errorf("Expected error to contain title, got: %s", errStr)
		}
		if !strings.Contains(errStr, "This is a test error") {
			t.Errorf("Expected error to contain description, got: %s", errStr)
		}
		if !strings.Contains(errStr, "Try this") {
			t.Errorf("Expected error to contain suggestions, got: %s", errStr)
		}
		if !strings.Contains(errStr, "Suggestions:") {
			t.Errorf("Expected error to contain 'Suggestions:', got: %s", errStr)
		}
	})

	t.Run("Error without suggestions", func(t *testing.T) {
		err := &EnhancedError{
			Title:       "Simple Error",
			Description: "Just a description",
			Suggestions: []string{},
		}

		errStr := err.Error()
		if strings.Contains(errStr, "Suggestions:") {
			t.Error("Expected no 'Suggestions:' section when suggestions are empty")
		}
	})
}

// TestFormatEnhancedError tests styled error formatting
func TestFormatEnhancedError(t *testing.T) {
	t.Run("Format EnhancedError", func(t *testing.T) {
		err := &EnhancedError{
			Title:       "Format Test",
			Description: "Testing formatting",
			Suggestions: []string{"Suggestion 1", "Suggestion 2"},
		}

		formatted := FormatEnhancedError(err)
		if !strings.Contains(formatted, "Format Test") {
			t.Errorf("Expected formatted error to contain title, got: %s", formatted)
		}
		if !strings.Contains(formatted, "Testing formatting") {
			t.Errorf("Expected formatted error to contain description, got: %s", formatted)
		}
		if !strings.Contains(formatted, "💡") {
			t.Errorf("Expected formatted error to contain suggestion icon, got: %s", formatted)
		}
	})

	t.Run("Format regular error", func(t *testing.T) {
		err := errors.New("regular error message")
		formatted := FormatEnhancedError(err)
		
		if !strings.Contains(formatted, "regular error message") {
			t.Errorf("Expected formatted error to contain message, got: %s", formatted)
		}
		if !strings.Contains(formatted, "❌") {
			t.Errorf("Expected formatted error to contain error icon, got: %s", formatted)
		}
	})
}

// TestEnhanceFileError tests file error enhancement
func TestEnhanceFileError(t *testing.T) {
	t.Run("Nil error returns nil", func(t *testing.T) {
		result := EnhanceFileError(nil, "/path/to/file")
		if result != nil {
			t.Error("Expected nil for nil input")
		}
	})

	t.Run("File not found error", func(t *testing.T) {
		originalErr := errors.New("no such file or directory")
		enhanced := EnhanceFileError(originalErr, "/path/to/spec.yaml")

		enhancedErr, ok := enhanced.(*EnhancedError)
		if !ok {
			t.Fatal("Expected EnhancedError")
		}

		if enhancedErr.Title != "File Not Found" {
			t.Errorf("Expected title 'File Not Found', got: %s", enhancedErr.Title)
		}
		if !strings.Contains(enhancedErr.Description, "/path/to/spec.yaml") {
			t.Errorf("Expected description to contain file path, got: %s", enhancedErr.Description)
		}
		if len(enhancedErr.Suggestions) == 0 {
			t.Error("Expected suggestions for file not found error")
		}
		if enhancedErr.Original != originalErr {
			t.Error("Expected original error to be preserved")
		}
	})

	t.Run("Cannot find variant", func(t *testing.T) {
		originalErr := errors.New("cannot find the specified file")
		enhanced := EnhanceFileError(originalErr, "/missing.yaml")

		enhancedErr, ok := enhanced.(*EnhancedError)
		if !ok {
			t.Fatal("Expected EnhancedError for 'cannot find' error")
		}

		if enhancedErr.Title != "File Not Found" {
			t.Errorf("Expected title 'File Not Found', got: %s", enhancedErr.Title)
		}
	})

	t.Run("Permission denied error", func(t *testing.T) {
		originalErr := errors.New("permission denied")
		enhanced := EnhanceFileError(originalErr, "/etc/secret.yaml")

		enhancedErr, ok := enhanced.(*EnhancedError)
		if !ok {
			t.Fatal("Expected EnhancedError")
		}

		if enhancedErr.Title != "Permission Denied" {
			t.Errorf("Expected title 'Permission Denied', got: %s", enhancedErr.Title)
		}
		if !strings.Contains(enhancedErr.Description, "/etc/secret.yaml") {
			t.Errorf("Expected description to contain file path, got: %s", enhancedErr.Description)
		}
		if len(enhancedErr.Suggestions) == 0 {
			t.Error("Expected suggestions for permission denied error")
		}
	})

	t.Run("YAML parse error", func(t *testing.T) {
		originalErr := errors.New("yaml: unmarshal error")
		enhanced := EnhanceFileError(originalErr, "/spec.yaml")

		enhancedErr, ok := enhanced.(*EnhancedError)
		if !ok {
			t.Fatal("Expected EnhancedError")
		}

		if enhancedErr.Title != "Invalid File Format" {
			t.Errorf("Expected title 'Invalid File Format', got: %s", enhancedErr.Title)
		}
		if len(enhancedErr.Suggestions) == 0 {
			t.Error("Expected suggestions for parse error")
		}
	})

	t.Run("Unmarshal error", func(t *testing.T) {
		originalErr := errors.New("failed to unmarshal spec")
		enhanced := EnhanceFileError(originalErr, "/spec.yaml")

		enhancedErr, ok := enhanced.(*EnhancedError)
		if !ok {
			t.Fatal("Expected EnhancedError for unmarshal error")
		}

		if enhancedErr.Title != "Invalid File Format" {
			t.Errorf("Expected title 'Invalid File Format', got: %s", enhancedErr.Title)
		}
	})

	t.Run("Unrecognized error passes through", func(t *testing.T) {
		originalErr := errors.New("some other error")
		result := EnhanceFileError(originalErr, "/spec.yaml")

		if result != originalErr {
			t.Error("Expected unrecognized error to pass through unchanged")
		}
	})
}

// TestEnhanceNetworkError tests network error enhancement
func TestEnhanceNetworkError(t *testing.T) {
	t.Run("Nil error returns nil", func(t *testing.T) {
		result := EnhanceNetworkError(nil, "https://api.example.com")
		if result != nil {
			t.Error("Expected nil for nil input")
		}
	})

	t.Run("Connection refused error", func(t *testing.T) {
		originalErr := errors.New("connection refused")
		enhanced := EnhanceNetworkError(originalErr, "https://api.example.com")

		enhancedErr, ok := enhanced.(*EnhancedError)
		if !ok {
			t.Fatal("Expected EnhancedError")
		}

		if enhancedErr.Title != "Connection Refused" {
			t.Errorf("Expected title 'Connection Refused', got: %s", enhancedErr.Title)
		}
		if !strings.Contains(enhancedErr.Description, "https://api.example.com") {
			t.Errorf("Expected description to contain URL, got: %s", enhancedErr.Description)
		}
		if len(enhancedErr.Suggestions) == 0 {
			t.Error("Expected suggestions for connection refused error")
		}
	})

	t.Run("Timeout error", func(t *testing.T) {
		originalErr := errors.New("request timeout")
		enhanced := EnhanceNetworkError(originalErr, "https://slow-api.com")

		enhancedErr, ok := enhanced.(*EnhancedError)
		if !ok {
			t.Fatal("Expected EnhancedError")
		}

		if enhancedErr.Title != "Request Timeout" {
			t.Errorf("Expected title 'Request Timeout', got: %s", enhancedErr.Title)
		}
		if len(enhancedErr.Suggestions) == 0 {
			t.Error("Expected suggestions for timeout error")
		}
	})

	t.Run("Deadline exceeded error", func(t *testing.T) {
		originalErr := errors.New("context deadline exceeded")
		enhanced := EnhanceNetworkError(originalErr, "https://api.example.com")

		enhancedErr, ok := enhanced.(*EnhancedError)
		if !ok {
			t.Fatal("Expected EnhancedError for deadline exceeded")
		}

		if enhancedErr.Title != "Request Timeout" {
			t.Errorf("Expected title 'Request Timeout', got: %s", enhancedErr.Title)
		}
	})

	t.Run("DNS resolution error", func(t *testing.T) {
		originalErr := errors.New("no such host")
		enhanced := EnhanceNetworkError(originalErr, "https://nonexistent.example.com")

		enhancedErr, ok := enhanced.(*EnhancedError)
		if !ok {
			t.Fatal("Expected EnhancedError")
		}

		if enhancedErr.Title != "DNS Resolution Failed" {
			t.Errorf("Expected title 'DNS Resolution Failed', got: %s", enhancedErr.Title)
		}
		if !strings.Contains(enhancedErr.Description, "https://nonexistent.example.com") {
			t.Errorf("Expected description to contain URL, got: %s", enhancedErr.Description)
		}
		if len(enhancedErr.Suggestions) == 0 {
			t.Error("Expected suggestions for DNS error")
		}
	})

	t.Run("DNS error variant", func(t *testing.T) {
		originalErr := errors.New("dns lookup failed")
		enhanced := EnhanceNetworkError(originalErr, "https://example.com")

		enhancedErr, ok := enhanced.(*EnhancedError)
		if !ok {
			t.Fatal("Expected EnhancedError for dns error")
		}

		if enhancedErr.Title != "DNS Resolution Failed" {
			t.Errorf("Expected title 'DNS Resolution Failed', got: %s", enhancedErr.Title)
		}
	})

	t.Run("TLS error", func(t *testing.T) {
		originalErr := errors.New("tls handshake failed")
		enhanced := EnhanceNetworkError(originalErr, "https://api.example.com")

		enhancedErr, ok := enhanced.(*EnhancedError)
		if !ok {
			t.Fatal("Expected EnhancedError")
		}

		if enhancedErr.Title != "TLS/SSL Error" {
			t.Errorf("Expected title 'TLS/SSL Error', got: %s", enhancedErr.Title)
		}
		if len(enhancedErr.Suggestions) == 0 {
			t.Error("Expected suggestions for TLS error")
		}
	})

	t.Run("Certificate error", func(t *testing.T) {
		originalErr := errors.New("certificate is invalid")
		enhanced := EnhanceNetworkError(originalErr, "https://api.example.com")

		enhancedErr, ok := enhanced.(*EnhancedError)
		if !ok {
			t.Fatal("Expected EnhancedError for certificate error")
		}

		if enhancedErr.Title != "TLS/SSL Error" {
			t.Errorf("Expected title 'TLS/SSL Error', got: %s", enhancedErr.Title)
		}
	})

	t.Run("Unrecognized error passes through", func(t *testing.T) {
		originalErr := errors.New("unknown network error")
		result := EnhanceNetworkError(originalErr, "https://api.example.com")

		if result != originalErr {
			t.Error("Expected unrecognized error to pass through unchanged")
		}
	})
}

// TestEnhanceValidationError tests validation error enhancement
func TestEnhanceValidationError(t *testing.T) {
	t.Run("Nil error returns nil", func(t *testing.T) {
		result := EnhanceValidationError(nil)
		if result != nil {
			t.Error("Expected nil for nil input")
		}
	})

	t.Run("Required field error", func(t *testing.T) {
		originalErr := errors.New("field is required")
		enhanced := EnhanceValidationError(originalErr)

		enhancedErr, ok := enhanced.(*EnhancedError)
		if !ok {
			t.Fatal("Expected EnhancedError")
		}

		if enhancedErr.Title != "Validation Failed" {
			t.Errorf("Expected title 'Validation Failed', got: %s", enhancedErr.Title)
		}
		if len(enhancedErr.Suggestions) == 0 {
			t.Error("Expected suggestions for validation error")
		}
	})

	t.Run("Version error", func(t *testing.T) {
		originalErr := errors.New("unsupported openapi version")
		enhanced := EnhanceValidationError(originalErr)

		enhancedErr, ok := enhanced.(*EnhancedError)
		if !ok {
			t.Fatal("Expected EnhancedError")
		}

		if enhancedErr.Title != "Unsupported Version" {
			t.Errorf("Expected title 'Unsupported Version', got: %s", enhancedErr.Title)
		}
		if len(enhancedErr.Suggestions) == 0 {
			t.Error("Expected suggestions for version error")
		}
		// Check for specific suggestions about OpenAPI 3.x
		found := false
		for _, s := range enhancedErr.Suggestions {
			if strings.Contains(s, "OpenAPI 3.x") || strings.Contains(s, "3.") {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected suggestions to mention OpenAPI 3.x")
		}
	})

	t.Run("Version field error", func(t *testing.T) {
		originalErr := errors.New("invalid version field")
		enhanced := EnhanceValidationError(originalErr)

		enhancedErr, ok := enhanced.(*EnhancedError)
		if !ok {
			t.Fatal("Expected EnhancedError for version field error")
		}

		if enhancedErr.Title != "Unsupported Version" {
			t.Errorf("Expected title 'Unsupported Version', got: %s", enhancedErr.Title)
		}
	})

	t.Run("Unrecognized error passes through", func(t *testing.T) {
		originalErr := errors.New("some other validation issue")
		result := EnhanceValidationError(originalErr)

		if result != originalErr {
			t.Error("Expected unrecognized error to pass through unchanged")
		}
	})
}

// TestEnhancedError_Suggestions tests various suggestion scenarios
func TestEnhancedError_Suggestions(t *testing.T) {
	t.Run("Multiple suggestions", func(t *testing.T) {
		err := &EnhancedError{
			Title:       "Test",
			Description: "Test error",
			Suggestions: []string{"First", "Second", "Third"},
		}

		errStr := err.Error()
		if !strings.Contains(errStr, "First") {
			t.Error("Expected first suggestion")
		}
		if !strings.Contains(errStr, "Second") {
			t.Error("Expected second suggestion")
		}
		if !strings.Contains(errStr, "Third") {
			t.Error("Expected third suggestion")
		}
	})

	t.Run("Single suggestion", func(t *testing.T) {
		err := &EnhancedError{
			Title:       "Test",
			Description: "Test error",
			Suggestions: []string{"Only one"},
		}

		errStr := err.Error()
		if !strings.Contains(errStr, "Suggestions:") {
			t.Error("Expected 'Suggestions:' header even for single suggestion")
		}
		if !strings.Contains(errStr, "Only one") {
			t.Error("Expected suggestion to be present")
		}
	})
}

// TestEnhancedError_OriginalPreserved tests that original error is preserved
func TestEnhancedError_OriginalPreserved(t *testing.T) {
	originalErr := errors.New("original error message")
	
	testCases := []struct {
		name     string
		enhancer func(error, ...string) error
		args     []string
	}{
		{"File error", func(err error, args ...string) error { return EnhanceFileError(err, args[0]) }, []string{"/test"}},
		{"Network error", func(err error, args ...string) error { return EnhanceNetworkError(err, args[0]) }, []string{"https://test.com"}},
		{"Validation error", func(err error, args ...string) error { return EnhanceValidationError(err) }, []string{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			enhanced := tc.enhancer(originalErr, tc.args...)
			
			// If it's enhanced, check Original field
			if enhancedErr, ok := enhanced.(*EnhancedError); ok {
				if enhancedErr.Original != originalErr {
					t.Error("Expected original error to be preserved")
				}
			}
		})
	}
}

// TestErrorChaining tests that errors maintain error chain
func TestErrorChaining(t *testing.T) {
	// Use an error message that will trigger enhancement
	originalErr := errors.New("no such file or directory: root cause")
	
	// Enhance file error (this will be enhanced because of "no such file")
	fileErr := EnhanceFileError(originalErr, "/test.yaml")
	
	if enhancedErr, ok := fileErr.(*EnhancedError); ok {
		if enhancedErr.Original == nil {
			t.Error("Expected original error to be set")
		}
		if enhancedErr.Original != originalErr {
			t.Error("Expected original error to match root cause")
		}
	} else {
		t.Error("Expected EnhancedError type for 'no such file' error")
	}
}

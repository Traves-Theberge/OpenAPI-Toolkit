package errors

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// EnhancedError provides detailed error information with actionable suggestions
type EnhancedError struct {
	Title       string   // Short error title
	Description string   // Detailed error description
	Suggestions []string // Actionable suggestions for fixing
	Original    error    // Original error for reference
}

// Error implements the error interface
func (e *EnhancedError) Error() string {
	msg := fmt.Sprintf("%s: %s", e.Title, e.Description)
	if len(e.Suggestions) > 0 {
		msg += "\n\nSuggestions:"
		for _, s := range e.Suggestions {
			msg += "\n  • " + s
		}
	}
	return msg
}

// FormatEnhancedError creates a styled error message for display
func FormatEnhancedError(err error) string {
	if enhanced, ok := err.(*EnhancedError); ok {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")). // Red
			Bold(true)

		suggestionStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")) // Yellow

		msg := errorStyle.Render("❌ " + enhanced.Title)
		msg += "\n\n" + enhanced.Description

		if len(enhanced.Suggestions) > 0 {
			msg += "\n\n" + suggestionStyle.Render("💡 Suggestions:")
			for _, s := range enhanced.Suggestions {
				msg += "\n  • " + s
			}
		}
		return msg
	}

	// Fallback for regular errors
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Render("❌ Error: " + err.Error())
}

// enhanceFileError wraps file-related errors with helpful suggestions
func EnhanceFileError(err error, filePath string) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// File not found
	if strings.Contains(errStr, "no such file") || strings.Contains(errStr, "cannot find") {
		return &EnhancedError{
			Title:       "File Not Found",
			Description: fmt.Sprintf("Could not find the file: %s", filePath),
			Suggestions: []string{
				"Check if the file path is correct",
				"Use an absolute path (e.g., /home/user/spec.yaml)",
				"Verify the file exists using 'ls' command",
				"Make sure you have read permissions for the file",
			},
			Original: err,
		}
	}

	// Permission denied
	if strings.Contains(errStr, "permission denied") {
		return &EnhancedError{
			Title:       "Permission Denied",
			Description: fmt.Sprintf("Cannot read file: %s", filePath),
			Suggestions: []string{
				"Check file permissions with 'ls -l'",
				"Try running with appropriate permissions",
				"Make sure you own the file or have read access",
			},
			Original: err,
		}
	}

	// Parse errors
	if strings.Contains(errStr, "yaml") || strings.Contains(errStr, "unmarshal") {
		return &EnhancedError{
			Title:       "Invalid File Format",
			Description: "The file is not a valid OpenAPI specification",
			Suggestions: []string{
				"Ensure the file is valid YAML or JSON",
				"Check for syntax errors (quotes, indentation)",
				"Validate YAML at https://www.yamllint.com/",
				"Make sure it's an OpenAPI 3.x specification",
			},
			Original: err,
		}
	}

	return err
}

// enhanceNetworkError wraps network-related errors with helpful suggestions
func EnhanceNetworkError(err error, url string) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// Connection refused
	if strings.Contains(errStr, "connection refused") {
		return &EnhancedError{
			Title:       "Connection Refused",
			Description: fmt.Sprintf("Cannot connect to: %s", url),
			Suggestions: []string{
				"Check if the server is running",
				"Verify the URL and port are correct",
				"Check firewall settings",
				"Try pinging the host to verify connectivity",
			},
			Original: err,
		}
	}

	// Timeout
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded") {
		return &EnhancedError{
			Title:       "Request Timeout",
			Description: "The server took too long to respond",
			Suggestions: []string{
				"Check your internet connection",
				"The server might be overloaded - try again later",
				"Verify the URL points to the correct endpoint",
				"Check if the server is experiencing issues",
			},
			Original: err,
		}
	}

	// DNS resolution failure
	if strings.Contains(errStr, "no such host") || strings.Contains(errStr, "dns") {
		return &EnhancedError{
			Title:       "DNS Resolution Failed",
			Description: fmt.Sprintf("Cannot resolve hostname: %s", url),
			Suggestions: []string{
				"Check if the URL is spelled correctly",
				"Verify your DNS settings",
				"Try using the IP address directly",
				"Check your internet connection",
			},
			Original: err,
		}
	}

	// TLS/SSL errors
	if strings.Contains(errStr, "tls") || strings.Contains(errStr, "certificate") {
		return &EnhancedError{
			Title:       "TLS/SSL Error",
			Description: "Cannot establish secure connection",
			Suggestions: []string{
				"The server's SSL certificate might be invalid",
				"Check if the URL should use 'http' instead of 'https'",
				"Verify the server's certificate is up to date",
				"Try accessing the URL in a browser first",
			},
			Original: err,
		}
	}

	return err
}

// enhanceValidationError wraps validation errors with helpful suggestions
func EnhanceValidationError(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// Missing required fields
	if strings.Contains(errStr, "required") {
		return &EnhancedError{
			Title:       "Validation Failed",
			Description: "The OpenAPI specification is missing required fields",
			Suggestions: []string{
				"Check that all required fields are present (openapi, info, paths)",
				"Verify the spec follows OpenAPI 3.x format",
				"Use a linter like Spectral to validate your spec",
				"See OpenAPI specification at https://spec.openapis.org/",
			},
			Original: err,
		}
	}

	// Version mismatch
	if strings.Contains(errStr, "version") || strings.Contains(errStr, "openapi") {
		return &EnhancedError{
			Title:       "Unsupported Version",
			Description: "This tool only supports OpenAPI 3.x specifications",
			Suggestions: []string{
				"Check the 'openapi' field in your spec",
				"If using Swagger 2.0, convert it to OpenAPI 3.x",
				"Use https://converter.swagger.io/ for conversion",
				"Ensure 'openapi' field starts with '3.' (e.g., '3.0.3')",
			},
			Original: err,
		}
	}

	return err
}

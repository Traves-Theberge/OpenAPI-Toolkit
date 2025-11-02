package ui

import (
	"strings"

	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

// FilterResults filters test results based on a query string
// Supports filtering by:
// - Status code (e.g., "200", "404", "500")
// - HTTP method (e.g., "GET", "POST")
// - Endpoint path (partial match, e.g., "users", "/api/")
// - Special keywords: "pass", "fail", "error", "success"
func FilterResults(results []models.TestResult, query string) []models.TestResult {
	if query == "" {
		return results
	}

	query = strings.ToLower(strings.TrimSpace(query))
	var filtered []models.TestResult

	for _, result := range results {
		if matchesFilter(result, query) {
			filtered = append(filtered, result)
		}
	}

	return filtered
}

// matchesFilter checks if a result matches the filter query
func matchesFilter(result models.TestResult, query string) bool {
	// First check if query matches special keywords (full word match only)
	switch query {
	case "pass", "passed", "success", "successful":
		// Match 2xx status codes
		return len(result.Status) > 0 && result.Status[0] == '2'
	case "fail", "failed", "err":
		// Match non-2xx status codes or "ERR" status
		return len(result.Status) > 0 && (result.Status[0] != '2' || result.Status == "ERR")
	}

	// Then check regular substring matches
	// Check status code
	if strings.Contains(strings.ToLower(result.Status), query) {
		return true
	}

	// Check HTTP method
	if strings.Contains(strings.ToLower(result.Method), query) {
		return true
	}

	// Check endpoint path
	if strings.Contains(strings.ToLower(result.Endpoint), query) {
		return true
	}

	// Check message (this allows "ok" and "error" to match message text)
	if strings.Contains(strings.ToLower(result.Message), query) {
		return true
	}

	return false
}

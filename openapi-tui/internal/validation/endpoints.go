package validation

import (
	"fmt"
	"os"
	"strings"

	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
	"github.com/getkin/kin-openapi/openapi3"
)

// ExtractEndpoints parses an OpenAPI spec and extracts all endpoints
func ExtractEndpoints(specPath string) ([]models.EndpointInfo, error) {
	// Read the spec file
	data, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec file: %w", err)
	}

	// Parse the OpenAPI spec
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	// Validate the spec
	if err := doc.Validate(loader.Context); err != nil {
		return nil, fmt.Errorf("invalid OpenAPI spec: %w", err)
	}

	// Extract endpoints
	var endpoints []models.EndpointInfo

	for path, pathItem := range doc.Paths.Map() {
		// Process each HTTP method
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
			if operation == nil {
				continue
			}

			endpoint := models.EndpointInfo{
				Path:        path,
				Method:      method,
				OperationID: operation.OperationID,
				Tags:        operation.Tags,
				Summary:     operation.Summary,
				Description: operation.Description,
				Selected:    true, // Default to selected
			}

			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints, nil
}

// FilterEndpoints filters endpoints based on a search query
func FilterEndpoints(endpoints []models.EndpointInfo, query string) []models.EndpointInfo {
	if query == "" {
		return endpoints
	}

	query = strings.ToLower(strings.TrimSpace(query))
	var filtered []models.EndpointInfo

	// Check for special filters (method:GET, tag:users, etc.)
	if strings.Contains(query, ":") {
		parts := strings.SplitN(query, ":", 2)
		if len(parts) == 2 {
			filterType := strings.TrimSpace(parts[0])
			filterValue := strings.ToLower(strings.TrimSpace(parts[1]))

			for _, ep := range endpoints {
				switch filterType {
				case "method", "m":
					if strings.ToLower(ep.Method) == filterValue {
						filtered = append(filtered, ep)
					}
				case "tag", "t":
					for _, tag := range ep.Tags {
						if strings.Contains(strings.ToLower(tag), filterValue) {
							filtered = append(filtered, ep)
							break
						}
					}
				case "path", "p":
					if strings.Contains(strings.ToLower(ep.Path), filterValue) {
						filtered = append(filtered, ep)
					}
				}
			}

			if len(filtered) > 0 {
				return filtered
			}
		}
	}

	// General fuzzy search across all fields
	for _, ep := range endpoints {
		if matchesQuery(ep, query) {
			filtered = append(filtered, ep)
		}
	}

	return filtered
}

// matchesQuery checks if an endpoint matches the search query
func matchesQuery(ep models.EndpointInfo, query string) bool {
	query = strings.ToLower(query)

	// Search in path
	if strings.Contains(strings.ToLower(ep.Path), query) {
		return true
	}

	// Search in method
	if strings.Contains(strings.ToLower(ep.Method), query) {
		return true
	}

	// Search in operation ID
	if strings.Contains(strings.ToLower(ep.OperationID), query) {
		return true
	}

	// Search in tags
	for _, tag := range ep.Tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}

	// Search in summary
	if strings.Contains(strings.ToLower(ep.Summary), query) {
		return true
	}

	// Search in description
	if strings.Contains(strings.ToLower(ep.Description), query) {
		return true
	}

	return false
}

// GetSelectedEndpoints returns only the selected endpoints
func GetSelectedEndpoints(endpoints []models.EndpointInfo) []models.EndpointInfo {
	var selected []models.EndpointInfo
	for _, ep := range endpoints {
		if ep.Selected {
			selected = append(selected, ep)
		}
	}
	return selected
}

// SelectAllEndpoints marks all endpoints as selected
func SelectAllEndpoints(endpoints []models.EndpointInfo) []models.EndpointInfo {
	for i := range endpoints {
		endpoints[i].Selected = true
	}
	return endpoints
}

// DeselectAllEndpoints marks all endpoints as deselected
func DeselectAllEndpoints(endpoints []models.EndpointInfo) []models.EndpointInfo {
	for i := range endpoints {
		endpoints[i].Selected = false
	}
	return endpoints
}

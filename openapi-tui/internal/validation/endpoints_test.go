package validation

import (
	"os"
	"testing"

	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

// TestExtractEndpoints tests endpoint extraction from OpenAPI spec
func TestExtractEndpoints(t *testing.T) {
	// Create a test OpenAPI spec
	specContent := `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /users:
    get:
      summary: Get all users
      operationId: getUsers
      tags:
        - users
      responses:
        '200':
          description: OK
    post:
      summary: Create a user
      operationId: createUser
      tags:
        - users
      responses:
        '201':
          description: Created
  /users/{id}:
    get:
      summary: Get user by ID
      operationId: getUserById
      tags:
        - users
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: OK
    delete:
      summary: Delete user
      operationId: deleteUser
      tags:
        - users
        - admin
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        '204':
          description: No Content
  /products:
    get:
      summary: List products
      operationId: listProducts
      tags:
        - products
      responses:
        '200':
          description: OK
`
	// Write spec to temp file
	tmpFile, err := os.CreateTemp("", "test-spec-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(specContent)); err != nil {
		t.Fatalf("Failed to write spec: %v", err)
	}
	tmpFile.Close()

	// Extract endpoints
	endpoints, err := ExtractEndpoints(tmpFile.Name())
	if err != nil {
		t.Fatalf("ExtractEndpoints failed: %v", err)
	}

	// Verify count
	if len(endpoints) != 5 {
		t.Errorf("Expected 5 endpoints, got %d", len(endpoints))
	}

	// Verify endpoint details
	foundGetUsers := false
	foundDeleteUser := false
	for _, ep := range endpoints {
		if ep.Path == "/users" && ep.Method == "GET" {
			foundGetUsers = true
			if ep.Summary != "Get all users" {
				t.Errorf("Expected summary 'Get all users', got '%s'", ep.Summary)
			}
			if ep.OperationID != "getUsers" {
				t.Errorf("Expected operationID 'getUsers', got '%s'", ep.OperationID)
			}
			if len(ep.Tags) != 1 || ep.Tags[0] != "users" {
				t.Errorf("Expected tags [users], got %v", ep.Tags)
			}
			if !ep.Selected {
				t.Error("Expected endpoint to be selected by default")
			}
		}
		if ep.Path == "/users/{id}" && ep.Method == "DELETE" {
			foundDeleteUser = true
			if len(ep.Tags) != 2 {
				t.Errorf("Expected 2 tags, got %d", len(ep.Tags))
			}
		}
	}

	if !foundGetUsers {
		t.Error("GET /users endpoint not found")
	}
	if !foundDeleteUser {
		t.Error("DELETE /users/{id} endpoint not found")
	}
}

// TestExtractEndpoints_InvalidFile tests error handling for invalid files
func TestExtractEndpoints_InvalidFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"nonexistent file", "/nonexistent/file.yaml", true},
		{"empty path", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ExtractEndpoints(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractEndpoints() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestExtractEndpoints_InvalidSpec tests handling of invalid OpenAPI specs
func TestExtractEndpoints_InvalidSpec(t *testing.T) {
	invalidSpec := `
this is not valid YAML
  or OpenAPI spec
`
	tmpFile, err := os.CreateTemp("", "invalid-spec-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.Write([]byte(invalidSpec))
	tmpFile.Close()

	_, err = ExtractEndpoints(tmpFile.Name())
	if err == nil {
		t.Error("Expected error for invalid spec, got nil")
	}
}

// TestExtractEndpoints_EmptySpec tests handling of spec with no paths
func TestExtractEndpoints_EmptySpec(t *testing.T) {
	emptySpec := `
openapi: 3.0.0
info:
  title: Empty API
  version: 1.0.0
paths: {}
`
	tmpFile, err := os.CreateTemp("", "empty-spec-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.Write([]byte(emptySpec))
	tmpFile.Close()

	endpoints, err := ExtractEndpoints(tmpFile.Name())
	if err != nil {
		t.Fatalf("ExtractEndpoints failed: %v", err)
	}

	if len(endpoints) != 0 {
		t.Errorf("Expected 0 endpoints, got %d", len(endpoints))
	}
}

// TestFilterEndpoints tests endpoint filtering
func TestFilterEndpoints(t *testing.T) {
	endpoints := []models.EndpointInfo{
		{Path: "/users", Method: "GET", Tags: []string{"users"}, Summary: "Get all users", Selected: true},
		{Path: "/users", Method: "POST", Tags: []string{"users"}, Summary: "Create user", Selected: true},
		{Path: "/users/{id}", Method: "GET", Tags: []string{"users"}, Summary: "Get user by ID", Selected: true},
		{Path: "/users/{id}", Method: "DELETE", Tags: []string{"users", "admin"}, Summary: "Delete user", Selected: true},
		{Path: "/products", Method: "GET", Tags: []string{"products"}, Summary: "List products", Selected: true},
		{Path: "/orders", Method: "POST", Tags: []string{"orders"}, Summary: "Create order", Selected: true},
	}

	tests := []struct {
		name     string
		query    string
		wantLen  int
		checkFn  func([]models.EndpointInfo) error
	}{
		{
			name:    "empty query returns all",
			query:   "",
			wantLen: 6,
		},
		{
			name:    "filter by path",
			query:   "users",
			wantLen: 4,
		},
		{
			name:    "filter by method",
			query:   "GET",
			wantLen: 3,
		},
		{
			name:    "filter by tag",
			query:   "admin",
			wantLen: 1,
		},
		{
			name:    "filter by summary",
			query:   "create",
			wantLen: 2,
		},
		{
			name:    "special filter method:GET",
			query:   "method:GET",
			wantLen: 3,
		},
		{
			name:    "special filter tag:users",
			query:   "tag:users",
			wantLen: 4,
		},
		{
			name:    "special filter path:/products",
			query:   "path:/products",
			wantLen: 1,
		},
		{
			name:    "case insensitive",
			query:   "USERS",
			wantLen: 4,
		},
		{
			name:    "no matches",
			query:   "nonexistent",
			wantLen: 0,
		},
		{
			name:    "partial path match",
			query:   "/users/",
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterEndpoints(endpoints, tt.query)
			if len(result) != tt.wantLen {
				t.Errorf("FilterEndpoints() returned %d endpoints, want %d", len(result), tt.wantLen)
			}
			if tt.checkFn != nil {
				if err := tt.checkFn(result); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

// TestMatchesQuery tests the matching logic
func TestMatchesQuery(t *testing.T) {
	endpoint := models.EndpointInfo{
		Path:        "/users/{id}",
		Method:      "GET",
		OperationID: "getUserById",
		Tags:        []string{"users", "public"},
		Summary:     "Get user by ID",
		Description: "Retrieves a single user by their unique identifier",
	}

	tests := []struct {
		name  string
		query string
		want  bool
	}{
		{"match path", "users", true},
		{"match method", "get", true},
		{"match operation ID", "getUserById", true},
		{"match tag", "public", true},
		{"match summary", "user by", true},
		{"match description", "identifier", true},
		{"no match", "products", false},
		{"case insensitive", "USERS", true},
		{"partial match", "user", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchesQuery(endpoint, tt.query)
			if got != tt.want {
				t.Errorf("matchesQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetSelectedEndpoints tests filtering for selected endpoints
func TestGetSelectedEndpoints(t *testing.T) {
	endpoints := []models.EndpointInfo{
		{Path: "/users", Method: "GET", Selected: true},
		{Path: "/users", Method: "POST", Selected: false},
		{Path: "/products", Method: "GET", Selected: true},
		{Path: "/orders", Method: "POST", Selected: false},
	}

	selected := GetSelectedEndpoints(endpoints)

	if len(selected) != 2 {
		t.Errorf("Expected 2 selected endpoints, got %d", len(selected))
	}

	for _, ep := range selected {
		if !ep.Selected {
			t.Errorf("GetSelectedEndpoints returned unselected endpoint: %s %s", ep.Method, ep.Path)
		}
	}
}

// TestGetSelectedEndpoints_NoneSelected tests when no endpoints are selected
func TestGetSelectedEndpoints_NoneSelected(t *testing.T) {
	endpoints := []models.EndpointInfo{
		{Path: "/users", Method: "GET", Selected: false},
		{Path: "/products", Method: "POST", Selected: false},
	}

	selected := GetSelectedEndpoints(endpoints)

	if len(selected) != 0 {
		t.Errorf("Expected 0 selected endpoints, got %d", len(selected))
	}
}

// TestGetSelectedEndpoints_AllSelected tests when all endpoints are selected
func TestGetSelectedEndpoints_AllSelected(t *testing.T) {
	endpoints := []models.EndpointInfo{
		{Path: "/users", Method: "GET", Selected: true},
		{Path: "/products", Method: "POST", Selected: true},
	}

	selected := GetSelectedEndpoints(endpoints)

	if len(selected) != 2 {
		t.Errorf("Expected 2 selected endpoints, got %d", len(selected))
	}
}

// TestSelectAllEndpoints tests selecting all endpoints
func TestSelectAllEndpoints(t *testing.T) {
	endpoints := []models.EndpointInfo{
		{Path: "/users", Method: "GET", Selected: false},
		{Path: "/users", Method: "POST", Selected: true},
		{Path: "/products", Method: "GET", Selected: false},
	}

	result := SelectAllEndpoints(endpoints)

	for i, ep := range result {
		if !ep.Selected {
			t.Errorf("Endpoint %d not selected after SelectAllEndpoints", i)
		}
	}
}

// TestDeselectAllEndpoints tests deselecting all endpoints
func TestDeselectAllEndpoints(t *testing.T) {
	endpoints := []models.EndpointInfo{
		{Path: "/users", Method: "GET", Selected: true},
		{Path: "/users", Method: "POST", Selected: true},
		{Path: "/products", Method: "GET", Selected: false},
	}

	result := DeselectAllEndpoints(endpoints)

	for i, ep := range result {
		if ep.Selected {
			t.Errorf("Endpoint %d still selected after DeselectAllEndpoints", i)
		}
	}
}

// TestFilterEndpoints_SpecialFilters tests special filter syntax
func TestFilterEndpoints_SpecialFilters(t *testing.T) {
	endpoints := []models.EndpointInfo{
		{Path: "/users", Method: "GET", Tags: []string{"users", "public"}, Selected: true},
		{Path: "/users", Method: "POST", Tags: []string{"users"}, Selected: true},
		{Path: "/admin/users", Method: "DELETE", Tags: []string{"admin"}, Selected: true},
	}

	tests := []struct {
		name    string
		query   string
		wantLen int
	}{
		{"method filter lowercase", "method:get", 1},
		{"method filter uppercase", "method:GET", 1},
		{"method shorthand", "m:post", 1},
		{"tag filter", "tag:admin", 1},
		{"tag shorthand", "t:users", 2},
		{"path filter", "path:/admin", 1},
		{"path shorthand", "p:admin", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterEndpoints(endpoints, tt.query)
			if len(result) != tt.wantLen {
				t.Errorf("FilterEndpoints(%q) returned %d endpoints, want %d", tt.query, len(result), tt.wantLen)
			}
		})
	}
}

// TestFilterEndpoints_EdgeCases tests edge cases
func TestFilterEndpoints_EdgeCases(t *testing.T) {
	t.Run("nil endpoints", func(t *testing.T) {
		result := FilterEndpoints(nil, "test")
		if result != nil {
			t.Error("Expected nil result for nil input")
		}
	})

	t.Run("empty endpoints", func(t *testing.T) {
		result := FilterEndpoints([]models.EndpointInfo{}, "test")
		if len(result) != 0 {
			t.Errorf("Expected 0 results for empty input, got %d", len(result))
		}
	})

	t.Run("whitespace query", func(t *testing.T) {
		endpoints := []models.EndpointInfo{
			{Path: "/users", Method: "GET", Selected: true},
		}
		result := FilterEndpoints(endpoints, "   ")
		if len(result) != 1 {
			t.Errorf("Expected 1 result for whitespace query, got %d", len(result))
		}
	})
}

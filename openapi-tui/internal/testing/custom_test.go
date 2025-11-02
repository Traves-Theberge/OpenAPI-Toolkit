package testing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

// TestExecuteCustomRequest tests basic custom request execution
func TestExecuteCustomRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	result, err := ExecuteCustomRequest("GET", server.URL+"/test", nil, "", nil, false)
	if err != nil {
		t.Fatalf("ExecuteCustomRequest failed: %v", err)
	}

	if result.Method != "GET" {
		t.Errorf("Expected method GET, got %s", result.Method)
	}
	if result.Status != "200" {
		t.Errorf("Expected status 200, got %s", result.Status)
	}
}

// TestExecuteCustomRequest_Methods tests different HTTP methods
func TestExecuteCustomRequest_Methods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != method {
					t.Errorf("Expected method %s, got %s", method, r.Method)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			result, err := ExecuteCustomRequest(method, server.URL, nil, "", nil, false)
			if err != nil {
				t.Fatalf("ExecuteCustomRequest failed for %s: %v", method, err)
			}

			if result.Status != "200" {
				t.Errorf("Expected status 200, got %s", result.Status)
			}
		})
	}
}

// TestExecuteCustomRequest_InvalidMethod tests invalid HTTP method
func TestExecuteCustomRequest_InvalidMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result, err := ExecuteCustomRequest("INVALID", server.URL, nil, "", nil, false)
	if err == nil {
		t.Fatal("Expected error for invalid method, got nil")
	}

	if result.Status != "ERR" {
		t.Errorf("Expected status ERR, got %s", result.Status)
	}
}

// TestExecuteCustomRequest_WithHeaders tests custom headers
func TestExecuteCustomRequest_WithHeaders(t *testing.T) {
	expectedHeaders := map[string]string{
		"X-Custom-Header": "test-value",
		"Authorization":   "Bearer token123",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, expectedValue := range expectedHeaders {
			if r.Header.Get(key) != expectedValue {
				t.Errorf("Expected header %s: %s, got %s", key, expectedValue, r.Header.Get(key))
			}
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result, err := ExecuteCustomRequest("GET", server.URL, expectedHeaders, "", nil, false)
	if err != nil {
		t.Fatalf("ExecuteCustomRequest failed: %v", err)
	}

	if result.Status != "200" {
		t.Errorf("Expected status 200, got %s", result.Status)
	}
}

// TestExecuteCustomRequest_WithBody tests request with JSON body
func TestExecuteCustomRequest_WithBody(t *testing.T) {
	expectedBody := `{"name": "test", "value": 123}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1}`))
	}))
	defer server.Close()

	result, err := ExecuteCustomRequest("POST", server.URL, nil, expectedBody, nil, false)
	if err != nil {
		t.Fatalf("ExecuteCustomRequest failed: %v", err)
	}

	if result.Status != "201" {
		t.Errorf("Expected status 201, got %s", result.Status)
	}
}

// TestExecuteCustomRequest_InvalidJSON tests invalid JSON body
func TestExecuteCustomRequest_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	invalidJSON := `{"invalid": json}`

	result, err := ExecuteCustomRequest("POST", server.URL, nil, invalidJSON, nil, false)
	if err == nil {
		t.Fatal("Expected error for invalid JSON, got nil")
	}

	if result.Status != "ERR" {
		t.Errorf("Expected status ERR, got %s", result.Status)
	}
}

// TestExecuteCustomRequest_WithAuth tests authentication
func TestExecuteCustomRequest_WithAuth(t *testing.T) {
	t.Run("Bearer Token", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth != "Bearer test-token" {
				t.Errorf("Expected Authorization header 'Bearer test-token', got '%s'", auth)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		auth := &models.AuthConfig{
			AuthType: "Bearer",
			Token:    "test-token",
		}

		result, err := ExecuteCustomRequest("GET", server.URL, nil, "", auth, false)
		if err != nil {
			t.Fatalf("ExecuteCustomRequest failed: %v", err)
		}

		if result.Status != "200" {
			t.Errorf("Expected status 200, got %s", result.Status)
		}
	})

	t.Run("API Key Header", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			if apiKey != "test-api-key" {
				t.Errorf("Expected X-API-Key header 'test-api-key', got '%s'", apiKey)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		auth := &models.AuthConfig{
			AuthType:   "API Key",
			Token:      "test-api-key",
			APIKeyIn:   "header",
			APIKeyName: "X-API-Key",
		}

		result, err := ExecuteCustomRequest("GET", server.URL, nil, "", auth, false)
		if err != nil {
			t.Fatalf("ExecuteCustomRequest failed: %v", err)
		}

		if result.Status != "200" {
			t.Errorf("Expected status 200, got %s", result.Status)
		}
	})

	t.Run("Basic Auth", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok || username != "testuser" || password != "testpass" {
				t.Errorf("Expected Basic Auth testuser:testpass, got %s:%s (ok=%v)", username, password, ok)
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		auth := &models.AuthConfig{
			AuthType: "Basic",
			Username: "testuser",
			Password: "testpass",
		}

		result, err := ExecuteCustomRequest("GET", server.URL, nil, "", auth, false)
		if err != nil {
			t.Fatalf("ExecuteCustomRequest failed: %v", err)
		}

		if result.Status != "200" {
			t.Errorf("Expected status 200, got %s", result.Status)
		}
	})
}

// TestExecuteCustomRequest_VerboseLogging tests verbose mode logging
func TestExecuteCustomRequest_VerboseLogging(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Response-Header", "test-value")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "logged"}`))
	}))
	defer server.Close()

	headers := map[string]string{"X-Request-Header": "test"}
	body := `{"request": "data"}`

	result, err := ExecuteCustomRequest("POST", server.URL, headers, body, nil, true)
	if err != nil {
		t.Fatalf("ExecuteCustomRequest failed: %v", err)
	}

	if result.LogEntry == nil {
		t.Fatal("Expected LogEntry to be populated in verbose mode")
	}

	if result.LogEntry.RequestURL != server.URL {
		t.Errorf("Expected RequestURL %s, got %s", server.URL, result.LogEntry.RequestURL)
	}

	if result.LogEntry.RequestBody != body {
		t.Errorf("Expected RequestBody %s, got %s", body, result.LogEntry.RequestBody)
	}

	if result.LogEntry.ResponseBody == "" {
		t.Error("Expected ResponseBody to be populated")
	}
}

// TestValidateJSONBody tests JSON validation
func TestValidateJSONBody(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{"empty body", "", false},
		{"valid JSON object", `{"key": "value"}`, false},
		{"valid JSON array", `[1, 2, 3]`, false},
		{"valid nested JSON", `{"user": {"name": "test", "age": 25}}`, false},
		{"invalid JSON", `{"invalid": }`, true},
		{"invalid JSON missing quote", `{"key": value}`, true},
		{"plain text", `this is not json`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateJSONBody(tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJSONBody() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestFormatJSONBody tests JSON formatting
func TestFormatJSONBody(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    string
		wantErr bool
	}{
		{
			"empty body",
			"",
			"",
			false,
		},
		{
			"compact JSON",
			`{"key":"value","number":123}`,
			"{\n  \"key\": \"value\",\n  \"number\": 123\n}",
			false,
		},
		{
			"already formatted",
			"{\n  \"key\": \"value\"\n}",
			"{\n  \"key\": \"value\"\n}",
			false,
		},
		{
			"invalid JSON",
			`{"invalid": }`,
			"",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatJSONBody(tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatJSONBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FormatJSONBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestExecuteCustomRequestCmd tests Bubble Tea command wrapper
func TestExecuteCustomRequestCmd(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
	defer server.Close()

	cmd := ExecuteCustomRequestCmd("GET", server.URL, nil, "", nil, false)
	msg := cmd()

	switch msg := msg.(type) {
	case TestCompleteMsg:
		if len(msg.Results) != 1 {
			t.Errorf("Expected 1 result, got %d", len(msg.Results))
		}
		if msg.Results[0].Status != "200" {
			t.Errorf("Expected status 200, got %s", msg.Results[0].Status)
		}
	case TestErrorMsg:
		t.Fatalf("Expected TestCompleteMsg, got TestErrorMsg: %v", msg.Err)
	default:
		t.Fatalf("Expected TestCompleteMsg, got %T", msg)
	}
}

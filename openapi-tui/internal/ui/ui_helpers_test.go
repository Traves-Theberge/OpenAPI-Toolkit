package ui

import (
	"testing"
	"time"

	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

func TestInitialConfigEditorModel_WithConfig(t *testing.T) {
	cfg := models.Config{
		SpecPath:       "openapi.yaml",
		BaseURL:        "http://example.com",
		VerboseMode:    true,
		MaxConcurrency: 5,
		MaxRetries:     3,
		RetryDelay:     1000,
		Auth: &models.AuthConfig{
			AuthType:   "bearer",
			Token:      "test-token",
			APIKeyName: "X-API-Key",
			APIKeyIn:   "header",
			Username:   "user",
			Password:   "pass",
		},
	}

	cem := InitialConfigEditorModel(cfg)

	// Verify all fields are populated from config
	if cem.SpecPathInput.Value() != "openapi.yaml" {
		t.Errorf("Expected SpecPath openapi.yaml, got %s", cem.SpecPathInput.Value())
	}
	if cem.BaseURLInput.Value() != "http://example.com" {
		t.Errorf("Expected BaseURL http://example.com, got %s", cem.BaseURLInput.Value())
	}
	if cem.AuthTypeInput.Value() != "bearer" {
		t.Errorf("Expected AuthType bearer, got %s", cem.AuthTypeInput.Value())
	}
	if cem.TokenInput.Value() != "test-token" {
		t.Errorf("Expected Token test-token, got %s", cem.TokenInput.Value())
	}
	if cem.APIKeyNameInput.Value() != "X-API-Key" {
		t.Errorf("Expected APIKeyName X-API-Key, got %s", cem.APIKeyNameInput.Value())
	}
	if cem.APIKeyInInput.Value() != "header" {
		t.Errorf("Expected APIKeyIn header, got %s", cem.APIKeyInInput.Value())
	}
	if cem.UsernameInput.Value() != "user" {
		t.Errorf("Expected Username user, got %s", cem.UsernameInput.Value())
	}
	if cem.PasswordInput.Value() != "pass" {
		t.Errorf("Expected Password pass, got %s", cem.PasswordInput.Value())
	}
	if cem.VerboseInput.Value() != "true" {
		t.Errorf("Expected VerboseInput true, got %s", cem.VerboseInput.Value())
	}
}

func TestInitialConfigEditorModel_VerboseFalse(t *testing.T) {
	cfg := models.Config{
		VerboseMode: false,
	}

	cem := InitialConfigEditorModel(cfg)

	if cem.VerboseInput.Value() != "false" {
		t.Errorf("Expected VerboseInput false, got %s", cem.VerboseInput.Value())
	}
}

func TestFormatStats_AllBranches(t *testing.T) {
	tests := []struct {
		name  string
		stats TestStats
	}{
		{
			name: "with failed tests",
			stats: TestStats{
				Total:           10,
				Passed:          7,
				Failed:          3,
				AverageTime:     150 * time.Millisecond,
				TotalTime:       1505 * time.Millisecond,
				FastestTime:     50 * time.Millisecond,
				SlowestTime:     300 * time.Millisecond,
				FastestEndpoint: "/fast",
				SlowestEndpoint: "/slow",
			},
		},
		{
			name: "all passed",
			stats: TestStats{
				Total:           5,
				Passed:          5,
				Failed:          0,
				AverageTime:     100 * time.Millisecond,
				TotalTime:       500 * time.Millisecond,
				FastestTime:     80 * time.Millisecond,
				SlowestTime:     120 * time.Millisecond,
				FastestEndpoint: "/users",
				SlowestEndpoint: "/api",
			},
		},
		{
			name: "all failed",
			stats: TestStats{
				Total:       5,
				Passed:      0,
				Failed:      5,
				AverageTime: 50 * time.Millisecond,
				TotalTime:   250 * time.Millisecond,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := FormatStats(tt.stats)
			if output == "" {
				t.Errorf("FormatStats returned empty string for %s", tt.name)
			}
		})
	}
}

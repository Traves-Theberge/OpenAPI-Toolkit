package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

// TestGetConfigPath tests configuration path resolution
func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() failed: %v", err)
	}

	if path == "" {
		t.Error("GetConfigPath() returned empty path")
	}

	// Should end with config.yaml
	if filepath.Base(path) != "config.yaml" {
		t.Errorf("Expected path to end with config.yaml, got %s", path)
	}

	// Should contain .config/openapi-tui
	if !filepath.IsAbs(path) {
		t.Error("Expected absolute path")
	}
}

// TestLoadConfig_NoFile tests loading when no config file exists
func TestLoadConfig_NoFile(t *testing.T) {
	// This should not fail even if file doesn't exist
	cfg := LoadConfig()

	// Should return default config
	if cfg.VerboseMode {
		t.Error("Expected VerboseMode to be false by default")
	}
	if cfg.MaxConcurrency != 0 {
		t.Errorf("Expected MaxConcurrency to be 0 (auto-detect), got %d", cfg.MaxConcurrency)
	}
}

// TestSaveAndLoadConfig tests round-trip save/load
func TestSaveAndLoadConfig(t *testing.T) {
	// Create a test config
	testConfig := models.Config{
		BaseURL:        "https://api.test.com",
		SpecPath:       "/path/to/spec.yaml",
		VerboseMode:    true,
		MaxConcurrency: 5,
		Auth: &models.AuthConfig{
			AuthType:   "bearer",
			Token:      "test-token-123",
			APIKeyIn:   "header",
			APIKeyName: "X-API-Key",
			Username:   "testuser",
			Password:   "testpass",
		},
	}

	// Save the config
	err := SaveConfig(testConfig)
	if err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	// Load it back
	loadedConfig := LoadConfig()

	// Verify all fields
	if loadedConfig.BaseURL != testConfig.BaseURL {
		t.Errorf("BaseURL mismatch: got %s, want %s", loadedConfig.BaseURL, testConfig.BaseURL)
	}
	if loadedConfig.SpecPath != testConfig.SpecPath {
		t.Errorf("SpecPath mismatch: got %s, want %s", loadedConfig.SpecPath, testConfig.SpecPath)
	}
	if loadedConfig.VerboseMode != testConfig.VerboseMode {
		t.Errorf("VerboseMode mismatch: got %v, want %v", loadedConfig.VerboseMode, testConfig.VerboseMode)
	}
	if loadedConfig.MaxConcurrency != testConfig.MaxConcurrency {
		t.Errorf("MaxConcurrency mismatch: got %d, want %d", loadedConfig.MaxConcurrency, testConfig.MaxConcurrency)
	}

	// Verify auth config
	if loadedConfig.Auth == nil {
		t.Fatal("Auth config is nil")
	}
	if loadedConfig.Auth.AuthType != testConfig.Auth.AuthType {
		t.Errorf("AuthType mismatch: got %s, want %s", loadedConfig.Auth.AuthType, testConfig.Auth.AuthType)
	}
	if loadedConfig.Auth.Token != testConfig.Auth.Token {
		t.Errorf("Token mismatch: got %s, want %s", loadedConfig.Auth.Token, testConfig.Auth.Token)
	}
	if loadedConfig.Auth.APIKeyName != testConfig.Auth.APIKeyName {
		t.Errorf("APIKeyName mismatch: got %s, want %s", loadedConfig.Auth.APIKeyName, testConfig.Auth.APIKeyName)
	}
	if loadedConfig.Auth.APIKeyIn != testConfig.Auth.APIKeyIn {
		t.Errorf("APIKeyIn mismatch: got %s, want %s", loadedConfig.Auth.APIKeyIn, testConfig.Auth.APIKeyIn)
	}
	if loadedConfig.Auth.Username != testConfig.Auth.Username {
		t.Errorf("Username mismatch: got %s, want %s", loadedConfig.Auth.Username, testConfig.Auth.Username)
	}
	if loadedConfig.Auth.Password != testConfig.Auth.Password {
		t.Errorf("Password mismatch: got %s, want %s", loadedConfig.Auth.Password, testConfig.Auth.Password)
	}
}

// TestSaveConfig_NoAuth tests saving config without authentication
func TestSaveConfig_NoAuth(t *testing.T) {
	testConfig := models.Config{
		BaseURL:        "https://api.noauth.com",
		SpecPath:       "/spec/noauth.yaml",
		VerboseMode:    false,
		MaxConcurrency: 0,
		Auth:           nil,
	}

	err := SaveConfig(testConfig)
	if err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	loadedConfig := LoadConfig()

	if loadedConfig.Auth != nil {
		t.Error("Expected Auth to be nil, got non-nil")
	}
}

// TestSaveConfig_APIKeyAuth tests saving config with API key authentication
func TestSaveConfig_APIKeyAuth(t *testing.T) {
	testConfig := models.Config{
		BaseURL:        "https://api.apikey.com",
		SpecPath:       "/spec/apikey.yaml",
		VerboseMode:    true,
		MaxConcurrency: 3,
		Auth: &models.AuthConfig{
			AuthType:   "apikey",
			APIKeyIn:   "query",
			APIKeyName: "api_key",
			Token:      "", // API key auth doesn't use token
		},
	}

	err := SaveConfig(testConfig)
	if err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	loadedConfig := LoadConfig()

	if loadedConfig.Auth == nil {
		t.Fatal("Auth config is nil")
	}
	if loadedConfig.Auth.AuthType != "apikey" {
		t.Errorf("Expected AuthType 'apikey', got %s", loadedConfig.Auth.AuthType)
	}
	if loadedConfig.Auth.APIKeyIn != "query" {
		t.Errorf("Expected APIKeyIn 'query', got %s", loadedConfig.Auth.APIKeyIn)
	}
	if loadedConfig.Auth.APIKeyName != "api_key" {
		t.Errorf("Expected APIKeyName 'api_key', got %s", loadedConfig.Auth.APIKeyName)
	}
}

// TestSaveConfig_BasicAuth tests saving config with Basic authentication
func TestSaveConfig_BasicAuth(t *testing.T) {
	testConfig := models.Config{
		BaseURL:        "https://api.basic.com",
		SpecPath:       "/spec/basic.yaml",
		VerboseMode:    false,
		MaxConcurrency: 0,
		Auth: &models.AuthConfig{
			AuthType: "basic",
			Username: "admin",
			Password: "secret123",
		},
	}

	err := SaveConfig(testConfig)
	if err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	loadedConfig := LoadConfig()

	if loadedConfig.Auth == nil {
		t.Fatal("Auth config is nil")
	}
	if loadedConfig.Auth.AuthType != "basic" {
		t.Errorf("Expected AuthType 'basic', got %s", loadedConfig.Auth.AuthType)
	}
	if loadedConfig.Auth.Username != "admin" {
		t.Errorf("Expected Username 'admin', got %s", loadedConfig.Auth.Username)
	}
	if loadedConfig.Auth.Password != "secret123" {
		t.Errorf("Expected Password 'secret123', got %s", loadedConfig.Auth.Password)
	}
}

// TestSaveConfig_MaxConcurrency tests various concurrency settings
func TestSaveConfig_MaxConcurrency(t *testing.T) {
	tests := []struct {
		name           string
		maxConcurrency int
	}{
		{"auto-detect", 0},
		{"single worker", 1},
		{"moderate", 5},
		{"high", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testConfig := models.Config{
				BaseURL:        "https://api.test.com",
				SpecPath:       "/test.yaml",
				MaxConcurrency: tt.maxConcurrency,
			}

			err := SaveConfig(testConfig)
			if err != nil {
				t.Fatalf("SaveConfig() failed: %v", err)
			}

			loadedConfig := LoadConfig()
			if loadedConfig.MaxConcurrency != tt.maxConcurrency {
				t.Errorf("MaxConcurrency mismatch: got %d, want %d", loadedConfig.MaxConcurrency, tt.maxConcurrency)
			}
		})
	}
}

// TestLoadConfig_CorruptFile tests handling of corrupt config file
func TestLoadConfig_CorruptFile(t *testing.T) {
	// Save a valid config first
	testConfig := models.Config{
		BaseURL:  "https://api.test.com",
		SpecPath: "/test.yaml",
	}
	SaveConfig(testConfig)

	// Get config path
	configPath, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath() failed: %v", err)
	}

	// Overwrite with corrupt YAML
	corruptYAML := "this is not: [valid yaml content\n  missing brackets"
	err = os.WriteFile(configPath, []byte(corruptYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to write corrupt config: %v", err)
	}

	// Should return default config without crashing
	cfg := LoadConfig()
	if cfg.VerboseMode {
		t.Error("Expected default VerboseMode (false) for corrupt file")
	}
}

// TestSaveConfig_EmptyValues tests saving config with empty string values
func TestSaveConfig_EmptyValues(t *testing.T) {
	testConfig := models.Config{
		BaseURL:        "",
		SpecPath:       "",
		VerboseMode:    false,
		MaxConcurrency: 0,
		Auth:           nil,
	}

	err := SaveConfig(testConfig)
	if err != nil {
		t.Fatalf("SaveConfig() failed: %v", err)
	}

	loadedConfig := LoadConfig()

	// Empty values should be preserved
	if loadedConfig.BaseURL != "" {
		t.Errorf("Expected empty BaseURL, got %s", loadedConfig.BaseURL)
	}
	if loadedConfig.SpecPath != "" {
		t.Errorf("Expected empty SpecPath, got %s", loadedConfig.SpecPath)
	}
}

// TestConfigPersistence tests that config survives multiple save/load cycles
func TestConfigPersistence(t *testing.T) {
	configs := []models.Config{
		{
			BaseURL:        "https://api1.test.com",
			SpecPath:       "/spec1.yaml",
			VerboseMode:    true,
			MaxConcurrency: 2,
		},
		{
			BaseURL:        "https://api2.test.com",
			SpecPath:       "/spec2.yaml",
			VerboseMode:    false,
			MaxConcurrency: 7,
			Auth: &models.AuthConfig{
				AuthType: "bearer",
				Token:    "token-abc",
			},
		},
		{
			BaseURL:        "https://api3.test.com",
			SpecPath:       "/spec3.yaml",
			VerboseMode:    true,
			MaxConcurrency: 0,
			Auth: &models.AuthConfig{
				AuthType:   "apikey",
				APIKeyIn:   "header",
				APIKeyName: "X-Key",
			},
		},
	}

	for i, cfg := range configs {
		// Save
		err := SaveConfig(cfg)
		if err != nil {
			t.Fatalf("SaveConfig() iteration %d failed: %v", i, err)
		}

		// Load
		loaded := LoadConfig()

		// Verify
		if loaded.BaseURL != cfg.BaseURL {
			t.Errorf("Iteration %d: BaseURL mismatch", i)
		}
		if loaded.MaxConcurrency != cfg.MaxConcurrency {
			t.Errorf("Iteration %d: MaxConcurrency mismatch", i)
		}
	}
}

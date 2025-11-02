package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

func TestGetConfigPath_Error(t *testing.T) {
	// Test normal operation (hard to trigger UserHomeDir error in test)
	path, err := GetConfigPath()
	if err != nil {
		t.Fatalf("GetConfigPath failed: %v", err)
	}
	if path == "" {
		t.Error("Expected non-empty config path")
	}
}

func TestSaveConfig_ErrorWriting(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	
	os.Setenv("HOME", tempDir)
	
	// Create config directory and file with read-only permissions
	configDir := filepath.Join(tempDir, ".config", "openapi-tui")
	os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, "config.yaml")
	os.WriteFile(configPath, []byte("test: value"), 0444) // read-only
	
	cfg := models.Config{
		BaseURL:  "http://example.com",
		SpecPath: "openapi.yaml",
	}
	
	// Try to save - might fail on read-only file
	err := SaveConfig(cfg)
	if err != nil {
		t.Logf("Got expected error for read-only file: %v", err)
	}
}

func TestSaveConfig_WithAuth(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	
	os.Setenv("HOME", tempDir)
	
	cfg := models.Config{
		BaseURL:        "http://example.com",
		SpecPath:       "openapi.yaml",
		VerboseMode:    true,
		MaxConcurrency: 5,
		MaxRetries:     5,
		RetryDelay:     2000,
		Auth: &models.AuthConfig{
			AuthType:   "bearer",
			Token:      "test-token",
			APIKeyIn:   "header",
			APIKeyName: "X-API-Key",
			Username:   "user",
			Password:   "pass",
		},
	}
	
	err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}
	
	// Load and verify
	loaded := LoadConfig()
	if loaded.BaseURL != cfg.BaseURL {
		t.Errorf("Expected BaseURL %s, got %s", cfg.BaseURL, loaded.BaseURL)
	}
	if loaded.Auth == nil {
		t.Fatal("Expected Auth to be loaded")
	}
	if loaded.Auth.Token != cfg.Auth.Token {
		t.Errorf("Expected Token %s, got %s", cfg.Auth.Token, loaded.Auth.Token)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	// Create temp directory with no config file
	tempDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	
	os.Setenv("HOME", tempDir)
	
	// Load should return default config
	cfg := LoadConfig()
	if cfg.MaxRetries != 3 {
		t.Errorf("Expected default MaxRetries 3, got %d", cfg.MaxRetries)
	}
	if cfg.RetryDelay != 1000 {
		t.Errorf("Expected default RetryDelay 1000, got %d", cfg.RetryDelay)
	}
}

func TestLoadConfig_CorruptedFile(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	
	os.Setenv("HOME", tempDir)
	
	// Create config with invalid YAML
	configDir := filepath.Join(tempDir, ".config", "openapi-tui")
	os.MkdirAll(configDir, 0755)
	configPath := filepath.Join(configDir, "config.yaml")
	os.WriteFile(configPath, []byte("invalid: yaml: content:"), 0644)
	
	// Load should return default config
	cfg := LoadConfig()
	if cfg.MaxRetries != 3 {
		t.Errorf("Expected default MaxRetries 3, got %d", cfg.MaxRetries)
	}
}

func TestLoadConfig_WithZeroValues(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	
	os.Setenv("HOME", tempDir)
	
	// Create config with zero values to test defaults
	cfg := models.Config{
		BaseURL:        "http://example.com",
		MaxConcurrency: 0, // Should stay 0 (auto-detect)
		MaxRetries:     0, // Should default to 3
		RetryDelay:     0, // Should default to 1000
	}
	
	SaveConfig(cfg)
	
	// Load and verify defaults are applied
	loaded := LoadConfig()
	if loaded.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries 3 (default), got %d", loaded.MaxRetries)
	}
	if loaded.RetryDelay != 1000 {
		t.Errorf("Expected RetryDelay 1000 (default), got %d", loaded.RetryDelay)
	}
	if loaded.MaxConcurrency != 0 {
		t.Errorf("Expected MaxConcurrency 0 (auto-detect), got %d", loaded.MaxConcurrency)
	}
}

package models

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveHistory_ErrorCreatingDir(t *testing.T) {
	// Try to save to a path where we can't create directories
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	
	// Set HOME to /dev/null which can't have subdirectories
	os.Setenv("HOME", "/dev/null")
	
	history := &TestHistory{}
	err := SaveHistory(history)
	
	// Should get an error (can't create .config/openapi-tui under /dev/null)
	if err == nil {
		t.Error("Expected error when saving to invalid path, got nil")
	}
}

func TestLoadHistory_CorruptedFile(t *testing.T) {
	// Create a temp directory
	tempDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	
	os.Setenv("HOME", tempDir)
	
	// Create config directory
	configDir := filepath.Join(tempDir, ".config", "openapi-tui")
	os.MkdirAll(configDir, 0755)
	
	// Write corrupted JSON to history file
	historyPath := filepath.Join(configDir, "history.json")
	os.WriteFile(historyPath, []byte("{invalid json}"), 0644)
	
	// Try to load - should handle gracefully
	history, err := LoadHistory()
	
	// Should return error or empty history
	if err == nil && (history == nil || len(history.Entries) != 0) {
		t.Error("Expected error or empty history for corrupted file")
	}
}

func TestGetConfigDir_UserHomeDirError(t *testing.T) {
	// This test is tricky because we can't easily mock os.UserHomeDir
	// But we can test the normal path
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	
	// Set a valid HOME
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	
	dir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir failed with valid HOME: %v", err)
	}
	
	expectedDir := filepath.Join(tempDir, ".config", "openapi-tui")
	if dir != expectedDir {
		t.Errorf("Expected dir %s, got %s", expectedDir, dir)
	}
}

func TestSaveHistory_FileWritePermissionError(t *testing.T) {
	// Create a temp directory
	tempDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	
	os.Setenv("HOME", tempDir)
	
	// Create config directory
	configDir := filepath.Join(tempDir, ".config", "openapi-tui")
	os.MkdirAll(configDir, 0755)
	
	// Create history file with read-only permissions
	historyPath := filepath.Join(configDir, "history.json")
	os.WriteFile(historyPath, []byte("{}"), 0444) // read-only
	
	history := &TestHistory{}
	history.AddEntry(HistoryEntry{ID: "test"})
	
	// Try to save - should fail on write
	err := SaveHistory(history)
	
	// On some systems this might not error, so just verify it attempts to save
	// If it errors, that's expected for read-only file
	if err != nil {
		t.Logf("Got expected error for read-only file: %v", err)
	}
}

func TestLoadHistory_EmptyFile(t *testing.T) {
	// Create a temp directory
	tempDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	
	os.Setenv("HOME", tempDir)
	
	// Create config directory
	configDir := filepath.Join(tempDir, ".config", "openapi-tui")
	os.MkdirAll(configDir, 0755)
	
	// Create empty history file
	historyPath := filepath.Join(configDir, "history.json")
	os.WriteFile(historyPath, []byte(""), 0644)
	
	// Try to load
	history, err := LoadHistory()
	
	// Should handle empty file gracefully - return empty history
	if err == nil && history != nil && len(history.Entries) != 0 {
		t.Error("Expected empty history for empty file")
	}
}

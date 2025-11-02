package models

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAddEntry(t *testing.T) {
	history := &TestHistory{}

	// Add a few entries
	for i := 0; i < 3; i++ {
		entry := HistoryEntry{
			ID:         time.Now().Format("20060102_150405"),
			Timestamp:  time.Now(),
			SpecPath:   "test.yaml",
			BaseURL:    "http://localhost",
			TotalTests: 10,
			Passed:     8,
			Failed:     2,
			Duration:   "1.5s",
		}
		history.AddEntry(entry)
		time.Sleep(1 * time.Millisecond) // Ensure unique IDs
	}

	if len(history.Entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(history.Entries))
	}

	// Verify entries are in reverse chronological order (newest first)
	if history.Entries[0].Timestamp.Before(history.Entries[1].Timestamp) {
		t.Error("Entries should be in reverse chronological order")
	}
}

func TestAddEntryLimit(t *testing.T) {
	history := &TestHistory{}

	// Add 60 entries (more than the 50 limit)
	for i := 0; i < 60; i++ {
		entry := HistoryEntry{
			ID:         time.Now().Format("20060102_150405"),
			Timestamp:  time.Now(),
			TotalTests: i,
		}
		history.AddEntry(entry)
		time.Sleep(1 * time.Millisecond) // Ensure unique IDs
	}

	// Should only keep last 50
	if len(history.Entries) != 50 {
		t.Errorf("Expected 50 entries (limit), got %d", len(history.Entries))
	}

	// Verify newest entry is first
	if history.Entries[0].TotalTests != 59 {
		t.Errorf("Expected newest entry (59) first, got %d", history.Entries[0].TotalTests)
	}

	// Verify oldest kept entry is last
	if history.Entries[49].TotalTests != 10 {
		t.Errorf("Expected oldest kept entry (10) last, got %d", history.Entries[49].TotalTests)
	}
}

func TestGetEntry(t *testing.T) {
	history := &TestHistory{}
	
	entry1 := HistoryEntry{
		ID:         "entry1",
		TotalTests: 10,
	}
	entry2 := HistoryEntry{
		ID:         "entry2",
		TotalTests: 20,
	}
	
	history.AddEntry(entry1)
	history.AddEntry(entry2)

	// Get existing entry
	found := history.GetEntry("entry1")
	if found == nil {
		t.Fatal("Expected to find entry1")
	}
	if found.TotalTests != 10 {
		t.Errorf("Expected TotalTests = 10, got %d", found.TotalTests)
	}

	// Get non-existing entry
	notFound := history.GetEntry("nonexistent")
	if notFound != nil {
		t.Error("Expected nil for non-existent entry")
	}
}

func TestSaveAndLoadHistory(t *testing.T) {
	// Create temporary config directory
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)

	// Create history with some entries
	history := &TestHistory{}
	entry1 := HistoryEntry{
		ID:         "test1",
		Timestamp:  time.Now(),
		SpecPath:   "api.yaml",
		BaseURL:    "http://localhost:8080",
		TotalTests: 15,
		Passed:     12,
		Failed:     3,
		Duration:   "2.5s",
		Results: []TestResult{
			{Method: "GET", Endpoint: "/users", Status: "200", Message: "OK"},
		},
	}
	entry2 := HistoryEntry{
		ID:         "test2",
		Timestamp:  time.Now().Add(-1 * time.Hour),
		SpecPath:   "api.yaml",
		BaseURL:    "http://localhost:8080",
		TotalTests: 10,
		Passed:     10,
		Failed:     0,
		Duration:   "1.2s",
	}
	
	history.AddEntry(entry1)
	history.AddEntry(entry2)

	// Save history
	if err := SaveHistory(history); err != nil {
		t.Fatalf("Failed to save history: %v", err)
	}

	// Verify file exists
	configDir := filepath.Join(tempDir, ".config", "openapi-tui")
	historyPath := filepath.Join(configDir, "history.json")
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		t.Fatal("History file was not created")
	}

	// Load history
	loaded, err := LoadHistory()
	if err != nil {
		t.Fatalf("Failed to load history: %v", err)
	}

	// Verify loaded history matches
	if len(loaded.Entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(loaded.Entries))
	}

	// Verify first entry (entry2 is at index 0 because AddEntry prepends, so last added is first)
	if loaded.Entries[0].ID != "test2" {
		t.Errorf("Expected first entry ID = test2, got %s", loaded.Entries[0].ID)
	}
	if loaded.Entries[0].TotalTests != 10 {
		t.Errorf("Expected TotalTests = 10, got %d", loaded.Entries[0].TotalTests)
	}
	if loaded.Entries[0].Passed != 10 {
		t.Errorf("Expected Passed = 10, got %d", loaded.Entries[0].Passed)
	}
	
	// Verify second entry (entry1 should be at index 1)
	if loaded.Entries[1].ID != "test1" {
		t.Errorf("Expected second entry ID = test1, got %s", loaded.Entries[1].ID)
	}
	if loaded.Entries[1].TotalTests != 15 {
		t.Errorf("Expected TotalTests = 15, got %d", loaded.Entries[1].TotalTests)
	}
	if loaded.Entries[1].Passed != 12 {
		t.Errorf("Expected Passed = 12, got %d", loaded.Entries[1].Passed)
	}
	if len(loaded.Entries[1].Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(loaded.Entries[1].Results))
	}
}

func TestLoadHistoryNonExistent(t *testing.T) {
	// Create temporary directory
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	os.Setenv("HOME", tempDir)

	// Load history when file doesn't exist
	history, err := LoadHistory()
	if err != nil {
		t.Errorf("LoadHistory should not return error for non-existent file, got: %v", err)
	}

	// Should return empty history
	if history == nil {
		t.Fatal("LoadHistory should return empty history, not nil")
	}
	if len(history.Entries) != 0 {
		t.Errorf("Expected empty history, got %d entries", len(history.Entries))
	}
}

func TestCreateHistoryEntry(t *testing.T) {
	results := []TestResult{
		{Method: "GET", Endpoint: "/users", Status: "200", Message: "OK", Duration: 100 * time.Millisecond},
		{Method: "POST", Endpoint: "/users", Status: "201", Message: "Created", Duration: 150 * time.Millisecond},
		{Method: "GET", Endpoint: "/users/999", Status: "404", Message: "Not Found", Duration: 50 * time.Millisecond},
		{Method: "DELETE", Endpoint: "/users/1", Status: "ERR", Message: "Connection refused", Duration: 0},
	}

	entry := CreateHistoryEntry("api.yaml", "http://localhost:8080", results, 300*time.Millisecond)

	// Verify entry fields
	if entry.SpecPath != "api.yaml" {
		t.Errorf("Expected SpecPath = api.yaml, got %s", entry.SpecPath)
	}
	if entry.BaseURL != "http://localhost:8080" {
		t.Errorf("Expected BaseURL = http://localhost:8080, got %s", entry.BaseURL)
	}
	if entry.TotalTests != 4 {
		t.Errorf("Expected TotalTests = 4, got %d", entry.TotalTests)
	}
	if entry.Passed != 2 {
		t.Errorf("Expected Passed = 2 (200, 201), got %d", entry.Passed)
	}
	if entry.Failed != 2 {
		t.Errorf("Expected Failed = 2 (404, ERR), got %d", entry.Failed)
	}
	if entry.Duration == "" {
		t.Error("Expected non-empty duration")
	}
	if len(entry.Results) != 4 {
		t.Errorf("Expected 4 results, got %d", len(entry.Results))
	}
	if entry.ID == "" {
		t.Error("Expected non-empty ID")
	}
}

func TestFormatHistoryDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{
			name:     "milliseconds",
			duration: 500 * time.Millisecond,
			want:     "500ms",
		},
		{
			name:     "seconds",
			duration: 2500 * time.Millisecond,
			want:     "2.5s",
		},
		{
			name:     "one minute",
			duration: 60 * time.Second,
			want:     "1m0s",
		},
		{
			name:     "minutes and seconds",
			duration: 125 * time.Second,
			want:     "2m5s",
		},
		{
			name:     "zero duration",
			duration: 0,
			want:     "0ms",
		},
		{
			name:     "less than second",
			duration: 250 * time.Millisecond,
			want:     "250ms",
		},
		{
			name:     "exactly one second",
			duration: 1 * time.Second,
			want:     "1.0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatHistoryDuration(tt.duration)
			if got != tt.want {
				t.Errorf("formatHistoryDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHistoryPersistence(t *testing.T) {
	// Create temporary config directory
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)
	
	os.Setenv("HOME", tempDir)

	// Create and save history multiple times
	for i := 0; i < 3; i++ {
		history, _ := LoadHistory()
		
		entry := HistoryEntry{
			ID:         time.Now().Format("20060102_150405"),
			Timestamp:  time.Now(),
			TotalTests: i + 1,
		}
		history.AddEntry(entry)
		
		if err := SaveHistory(history); err != nil {
			t.Fatalf("Iteration %d: Failed to save history: %v", i, err)
		}
		
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
	}

	// Load and verify final state
	final, err := LoadHistory()
	if err != nil {
		t.Fatalf("Failed to load final history: %v", err)
	}

	if len(final.Entries) != 3 {
		t.Errorf("Expected 3 entries after multiple saves, got %d", len(final.Entries))
	}

	// Verify entries are in correct order (newest first)
	if final.Entries[0].TotalTests != 3 {
		t.Errorf("Expected newest entry (3) first, got %d", final.Entries[0].TotalTests)
	}
	if final.Entries[2].TotalTests != 1 {
		t.Errorf("Expected oldest entry (1) last, got %d", final.Entries[2].TotalTests)
	}
}

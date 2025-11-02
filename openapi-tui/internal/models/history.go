package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// HistoryEntry represents a single test run in history
type HistoryEntry struct {
	ID          string       `json:"id"`
	Timestamp   time.Time    `json:"timestamp"`
	SpecPath    string       `json:"specPath"`
	BaseURL     string       `json:"baseUrl"`
	TotalTests  int          `json:"totalTests"`
	Passed      int          `json:"passed"`
	Failed      int          `json:"failed"`
	Duration    string       `json:"duration"`
	Results     []TestResult `json:"results"`
}

// TestHistory manages the history of test runs
type TestHistory struct {
	Entries []HistoryEntry `json:"entries"`
}

// AddEntry adds a new entry to the history
func (h *TestHistory) AddEntry(entry HistoryEntry) {
	// Limit history to last 50 entries
	h.Entries = append([]HistoryEntry{entry}, h.Entries...)
	if len(h.Entries) > 50 {
		h.Entries = h.Entries[:50]
	}
}

// GetEntry retrieves an entry by ID
func (h *TestHistory) GetEntry(id string) *HistoryEntry {
	for i := range h.Entries {
		if h.Entries[i].ID == id {
			return &h.Entries[i]
		}
	}
	return nil
}

// SaveHistory saves test history to a file
func SaveHistory(history *TestHistory) error {
	// Get config directory
	configDir, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// History file path
	historyPath := filepath.Join(configDir, "history.json")

	// Marshal to JSON
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}

	// Write to file
	if err := os.WriteFile(historyPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write history file: %w", err)
	}

	return nil
}

// LoadHistory loads test history from file
func LoadHistory() (*TestHistory, error) {
	// Get config directory
	configDir, err := getConfigDir()
	if err != nil {
		return &TestHistory{}, nil // Return empty history on error
	}

	// History file path
	historyPath := filepath.Join(configDir, "history.json")

	// Check if file exists
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		return &TestHistory{}, nil // Return empty history if file doesn't exist
	}

	// Read file
	data, err := os.ReadFile(historyPath)
	if err != nil {
		return &TestHistory{}, nil // Return empty history on error
	}

	// Unmarshal JSON
	var history TestHistory
	if err := json.Unmarshal(data, &history); err != nil {
		return &TestHistory{}, nil // Return empty history on parse error
	}

	return &history, nil
}

// CreateHistoryEntry creates a new history entry from test results
func CreateHistoryEntry(specPath, baseURL string, results []TestResult, duration time.Duration) HistoryEntry {
	// Calculate statistics
	passed := 0
	failed := 0
	for _, r := range results {
		if r.Status != "ERR" && r.Status[0] == '2' {
			passed++
		} else {
			failed++
		}
	}

	// Generate ID from timestamp
	id := time.Now().Format("20060102_150405")

	return HistoryEntry{
		ID:         id,
		Timestamp:  time.Now(),
		SpecPath:   specPath,
		BaseURL:    baseURL,
		TotalTests: len(results),
		Passed:     passed,
		Failed:     failed,
		Duration:   formatHistoryDuration(duration),
		Results:    results,
	}
}

// formatHistoryDuration formats a duration for display
func formatHistoryDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%dm%ds", minutes, seconds)
}

// getConfigDir returns the configuration directory path
func getConfigDir() (string, error) {
	// Get user home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Return config directory path
	return filepath.Join(home, ".config", "openapi-tui"), nil
}

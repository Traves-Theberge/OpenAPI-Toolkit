package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

// TestStats holds calculated statistics for test results
type TestStats struct {
	Total          int
	Passed         int
	Failed         int
	AverageTime    time.Duration
	TotalTime      time.Duration
	FastestTime    time.Duration
	SlowestTime    time.Duration
	FastestEndpoint string
	SlowestEndpoint string
}

// CalculateStats computes statistics from test results
func CalculateStats(results []models.TestResult) TestStats {
	stats := TestStats{
		Total: len(results),
	}

	if len(results) == 0 {
		return stats
	}

	// Initialize fastest to a high value only if we have results
	stats.FastestTime = time.Hour * 999

	var totalDuration time.Duration

	for _, result := range results {
		// Count passed/failed based on status
		// Passed: 2xx status codes (200, 201, 204, etc.) or "OK" status
		if result.Status[0] == '2' || result.Status == "OK" {
			stats.Passed++
		} else {
			stats.Failed++
		}

		// Track timing
		totalDuration += result.Duration

		// Find fastest
		if result.Duration < stats.FastestTime && result.Duration > 0 {
			stats.FastestTime = result.Duration
			stats.FastestEndpoint = fmt.Sprintf("%s %s", result.Method, result.Endpoint)
		}

		// Find slowest
		if result.Duration > stats.SlowestTime {
			stats.SlowestTime = result.Duration
			stats.SlowestEndpoint = fmt.Sprintf("%s %s", result.Method, result.Endpoint)
		}
	}

	stats.TotalTime = totalDuration
	if stats.Total > 0 {
		stats.AverageTime = totalDuration / time.Duration(stats.Total)
	}

	// Reset fastest if no valid durations found
	if stats.FastestTime == time.Hour*999 {
		stats.FastestTime = 0
	}

	return stats
}

// FormatStats renders statistics in a styled format
func FormatStats(stats TestStats) string {
	// Define colors
	successColor := lipgloss.Color("#4ECDC4")
	errorColor := lipgloss.Color("#FF6B6B")
	neutralColor := lipgloss.Color("#888")
	
	// Calculate pass rate
	passRate := 0.0
	if stats.Total > 0 {
		passRate = float64(stats.Passed) / float64(stats.Total) * 100
	}

	// Choose color based on pass rate
	passRateColor := successColor
	if passRate < 50 {
		passRateColor = errorColor
	} else if passRate < 100 {
		passRateColor = lipgloss.Color("#F9CA24") // Yellow for partial success
	}

	// Build the stats display
	statsLines := []string{
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).Render("📊 Test Summary"),
		"",
		fmt.Sprintf("Total Tests:  %s",
			lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%d", stats.Total))),
		fmt.Sprintf("Passed:       %s",
			lipgloss.NewStyle().Foreground(successColor).Bold(true).Render(fmt.Sprintf("%d", stats.Passed))),
		fmt.Sprintf("Failed:       %s",
			lipgloss.NewStyle().Foreground(errorColor).Bold(true).Render(fmt.Sprintf("%d", stats.Failed))),
		fmt.Sprintf("Pass Rate:    %s",
			lipgloss.NewStyle().Foreground(passRateColor).Bold(true).Render(fmt.Sprintf("%.1f%%", passRate))),
		"",
		lipgloss.NewStyle().Foreground(neutralColor).Render("⏱️  Timing:"),
		fmt.Sprintf("  Total:      %s", formatDuration(stats.TotalTime)),
		fmt.Sprintf("  Average:    %s", formatDuration(stats.AverageTime)),
	}

	// Add fastest/slowest if available
	if stats.FastestTime > 0 {
		statsLines = append(statsLines,
			fmt.Sprintf("  Fastest:    %s %s",
				formatDuration(stats.FastestTime),
				lipgloss.NewStyle().Foreground(neutralColor).Render("("+stats.FastestEndpoint+")")))
	}
	if stats.SlowestTime > 0 {
		statsLines = append(statsLines,
			fmt.Sprintf("  Slowest:    %s %s",
				formatDuration(stats.SlowestTime),
				lipgloss.NewStyle().Foreground(neutralColor).Render("("+stats.SlowestEndpoint+")")))
	}

	// Join all lines
	content := lipgloss.JoinVertical(lipgloss.Left, statsLines...)

	// Add border and padding
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4ECDC4")).
		Padding(1, 2).
		Render(content)
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d == 0 {
		return "0ms"
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

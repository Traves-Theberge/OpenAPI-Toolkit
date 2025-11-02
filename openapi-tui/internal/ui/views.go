package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/errors"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

// ViewMenu renders the main menu screen
func ViewMenu(m models.Model) string {
	// Title with styling
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		MarginBottom(1).
		Render("🚀 OpenAPI CLI TUI")

	// Menu options with styling - highlight selected item
	var menuItems []string
	options := []string{"📋 Validate OpenAPI Spec", "🧪 Test API", "❓ Help", "👋 Quit"}

	for i, option := range options {
		var cursor string
		var style lipgloss.Style

		if m.Cursor == i {
			// Selected item: red cursor and bold text
			cursor = "▶ "
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B")).
				Bold(true)
		} else {
			// Unselected item: cyan text
			cursor = "  "
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4ECDC4"))
		}

		menuItems = append(menuItems, style.Render(cursor+option))
	}

	menu := lipgloss.JoinVertical(lipgloss.Left, menuItems...)

	// Status indicators
	var statusIndicators []string
	if m.VerboseMode {
		statusIndicators = append(statusIndicators, "🔍 Verbose: ON")
	}
	if m.Config.BaseURL != "" || m.Config.SpecPath != "" {
		statusIndicators = append(statusIndicators, "💾 Config loaded")
	}
	
	verboseStatus := ""
	if len(statusIndicators) > 0 {
		verboseStatus = " • " + strings.Join(statusIndicators, " • ")
	}

	// Status bar with navigation hints
	statusBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		Padding(0, 1).
		MarginTop(1).
		Render("Press h or ? for help • v to toggle verbose • ↑↓/jk to navigate • Enter to select" + verboseStatus)

	// Combine sections vertically centered
	content := lipgloss.JoinVertical(lipgloss.Center, title, menu, statusBar)

	// Center the entire content with dynamic sizing and border
	borderedContent := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4ECDC4")).
		Padding(1, 2).
		Render(content)

	return lipgloss.Place(
		m.Width, m.Height,
		lipgloss.Center, lipgloss.Center,
		borderedContent,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#333")),
	)
}

// ViewHelp renders the help screen with keyboard shortcuts and usage information
func ViewHelp(m models.Model) string {
	// Help screen title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4")).
		Render("❓ Help & Shortcuts")

	// Help content - compact version
	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4ECDC4")).
		Width(54).
		Render(`Navigation:
  ↑/↓ j/k  - Navigate    Enter - Select
  h or ?   - Help        q/Esc - Back/Quit

Features:
  📋 Validate - Check spec validity
  🧪 Test     - Test all endpoints
  🎨 Styled   - Beautiful terminal UI

Tips:
  • Use relative or absolute paths
  • URLs need http:// or https://`)

	// Footer with exit instruction
	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Width(54).
		Render("\nPress q or Esc to return to menu")

	content := lipgloss.JoinVertical(lipgloss.Left, title, "", helpText, footer)

	// Center the entire content with dynamic sizing and border
	borderedContent := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4ECDC4")).
		Padding(1, 2).
		Render(content)

	return lipgloss.Place(
		m.Width, m.Height,
		lipgloss.Center, lipgloss.Center,
		borderedContent,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#333")),
	)
}

// ViewValidate renders the validation screen
func ViewValidate(m models.Model) string {
	var content string

	// Show validation results if completed
	if m.ValidateModel.Done {
		if m.ValidateModel.Err != nil {
			// Display enhanced validation error with suggestions
			content = errors.FormatEnhancedError(m.ValidateModel.Err)
		} else {
			// Display success message in green
			content = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4ECDC4")).
				Bold(true).
				Render(m.ValidateModel.Result)
		}
		// Add exit instruction
		content += "\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render("Press Enter or Esc to return to menu")
	} else {
		// Show input field for spec file path
		input := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4ECDC4")).
			Padding(1, 2).
			Render("> " + m.ValidateModel.TextInput.View())

		if m.ValidateModel.Err != nil {
			// Show enhanced input error with suggestions
			content = input + "\n\n" + errors.FormatEnhancedError(m.ValidateModel.Err)
		} else {
			// Show input instructions
			content = input + "\n\n" + lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888")).
				Render("Enter path to OpenAPI spec file and press Enter")
		}
	}

	// Center the entire content with dynamic sizing and border
	borderedContent := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4ECDC4")).
		Padding(1, 2).
		Render(content)

	return lipgloss.Place(
		m.Width, m.Height,
		lipgloss.Center, lipgloss.Center,
		borderedContent,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#333")),
	)
}

// ViewTest renders the testing screen
func ViewTest(m models.Model) string {
	var content string

	// Multi-step testing workflow UI
	switch m.TestModel.Step {
	case 0: // Spec file input
		input := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4ECDC4")).
			Padding(1, 2).
			Render("> " + m.TestModel.SpecInput.View())

		if m.TestModel.Err != nil {
			// Show enhanced input error for spec file with suggestions
			content = input + "\n\n" + errors.FormatEnhancedError(m.TestModel.Err)
		} else {
			// Show spec file input instructions
			content = input + "\n\n" + lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888")).
				Render("Enter path to OpenAPI spec file and press Enter")
		}
	case 1: // Base URL input
		input := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4ECDC4")).
			Padding(1, 2).
			Render("> " + m.TestModel.UrlInput.View())

		if m.TestModel.Err != nil {
			// Show enhanced input error for base URL with suggestions
			content = input + "\n\n" + errors.FormatEnhancedError(m.TestModel.Err)
		} else {
			// Show base URL input instructions
			content = input + "\n\n" + lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888")).
				Render("Enter base URL (e.g., https://api.example.com) and press Enter")
		}
	case 2: // Testing in progress
		// Show animated spinner with testing message
		spinnerView := m.TestModel.Spinner.View() + " Testing API endpoints..."
		content = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true).
			Render(spinnerView) + "\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render("Press Ctrl+C to cancel")
	case 3: // Results display
		if m.TestModel.Err != nil {
			// Show enhanced testing error with actionable suggestions
			content = errors.FormatEnhancedError(m.TestModel.Err)
		} else {
			// Determine which results to display
			resultsToShow := m.TestModel.Results
			filterQuery := m.TestModel.FilterInput.Value()
			
			// Apply filter if active and query is not empty
			if m.TestModel.FilterActive && filterQuery != "" {
				resultsToShow = FilterResults(m.TestModel.Results, filterQuery)
			}
			
			// Calculate and display summary statistics
			stats := CalculateStats(resultsToShow)
			statsView := FormatStats(stats)

			// Populate table with results (filtered or all)
			var rows []table.Row
			for _, r := range resultsToShow {
				rows = append(rows, table.Row{r.Method, r.Endpoint, r.Status, r.Message})
			}
			m.TestModel.Table.SetRows(rows)

			// Show filter input if active
			filterView := ""
			if m.TestModel.FilterActive {
				filterStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#4ECDC4")).
					Bold(true)
				filterView = filterStyle.Render("🔍 Filter: ") + m.TestModel.FilterInput.View() + 
					"\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#888")).
					Render(fmt.Sprintf("(Showing %d of %d results. Keywords: pass, fail, err)", 
						len(resultsToShow), len(m.TestModel.Results))) + "\n\n"
			}

			// Show success message, filter, stats, and results table
			content = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4ECDC4")).
				Bold(true).
				Render("✅ Testing Complete!") + "\n\n" + 
				filterView +
				statsView + "\n\n" +
				m.TestModel.Table.View()
			
			// Show export success message if results were exported
			if m.TestModel.ExportSuccess != "" {
				content += "\n\n" + lipgloss.NewStyle().
					Foreground(lipgloss.Color("#4ECDC4")).
					Render(m.TestModel.ExportSuccess)
			}
		}
		// Add instructions
		instructions := "Press 'f' to filter | 'e' JSON | 'h' HTML | 'j' JUnit XML"
		if m.VerboseMode {
			instructions += " | 'l' logs"
		}
		instructions += " | Enter to return"
		content += "\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render(instructions)
	case 4: // Log detail view
		if m.TestModel.SelectedLog >= 0 && m.TestModel.SelectedLog < len(m.TestModel.Results) {
			result := m.TestModel.Results[m.TestModel.SelectedLog]
			if result.LogEntry != nil {
				content = ViewLogDetail(m, result, result.LogEntry)
			}
		}
	}

	// Center the entire content with dynamic sizing and border
	borderedContent := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4ECDC4")).
		Padding(1, 2).
		Render(content)

	return lipgloss.Place(
		m.Width, m.Height,
		lipgloss.Center, lipgloss.Center,
		borderedContent,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#333")),
	)
}

// viewLogDetail renders the detailed log view for a test result
func ViewLogDetail(m models.Model, result models.TestResult, log *models.LogEntry) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))
	
	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#4ECDC4"))
	
	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA"))
	
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888"))
	
	// Title
	title := titleStyle.Render(fmt.Sprintf("📊 Log Details: %s %s", result.Method, result.Endpoint))
	
	// Request section
	requestSection := labelStyle.Render("🔼 REQUEST") + "\n"
	requestSection += labelStyle.Render("URL: ") + valueStyle.Render(log.RequestURL) + "\n"
	requestSection += labelStyle.Render("Timestamp: ") + valueStyle.Render(log.Timestamp.Format(time.RFC3339)) + "\n"
	requestSection += labelStyle.Render("Duration: ") + valueStyle.Render(log.Duration.String()) + "\n"
	
	// Request headers
	if len(log.RequestHeaders) > 0 {
		requestSection += "\n" + headerStyle.Render("Headers:")
		for k, v := range log.RequestHeaders {
			requestSection += "\n  " + labelStyle.Render(k+": ") + valueStyle.Render(v)
		}
	}
	
	// Request body
	if log.RequestBody != "" {
		requestSection += "\n\n" + headerStyle.Render("Body:")
		requestSection += "\n" + valueStyle.Render(log.RequestBody)
	}
	
	// Response section
	responseSection := "\n\n" + labelStyle.Render("🔽 RESPONSE") + "\n"
	responseSection += labelStyle.Render("Status: ") + valueStyle.Render(result.Status) + " - " + valueStyle.Render(result.Message) + "\n"
	
	// Response headers
	if len(log.ResponseHeaders) > 0 {
		responseSection += "\n" + headerStyle.Render("Headers:")
		for k, v := range log.ResponseHeaders {
			responseSection += "\n  " + labelStyle.Render(k+": ") + valueStyle.Render(v)
		}
	}
	
	// Response body
	if log.ResponseBody != "" {
		responseSection += "\n\n" + headerStyle.Render("Body:")
		responseSection += "\n" + valueStyle.Render(log.ResponseBody)
	}
	
	// Footer
	footer := "\n\n" + lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Render("Press Esc or Enter to return to results")
	
	return title + "\n\n" + requestSection + responseSection + footer
}

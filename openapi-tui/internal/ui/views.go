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
	options := []string{"📋 Validate OpenAPI Spec", "🧪 Test All Endpoints", "🎯 Select & Test Endpoints", "✏️  Custom Request", "📜 History", "❓ Help", "👋 Quit"}

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
		instructions := "Press 'f' to filter | 'e' JSON | 'h' HTML | 'j' JUnit XML | 'r' history"
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

// ViewHistory renders the test run history screen
func ViewHistory(m models.Model) string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		MarginBottom(1)

	title := titleStyle.Render("📜 Test Run History")

	if len(m.History.Entries) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Italic(true).
			Padding(2)
		
		empty := emptyStyle.Render("No test history yet. Run some tests to see them here!")
		
		footer := "\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render("Press Esc to return")
		
		return title + "\n\n" + empty + footer
	}

	// Create table for history entries
	columns := []table.Column{
		{Title: "Date & Time", Width: 20},
		{Title: "Spec", Width: 25},
		{Title: "Tests", Width: 12},
		{Title: "Passed", Width: 8},
		{Title: "Failed", Width: 8},
		{Title: "Duration", Width: 12},
	}

	rows := []table.Row{}
	for _, entry := range m.History.Entries {
		timestamp := entry.Timestamp.Format("2006-01-02 15:04:05")
		
		// Truncate spec path if too long
		spec := entry.SpecPath
		if len(spec) > 23 {
			spec = "..." + spec[len(spec)-20:]
		}
		
		tests := fmt.Sprintf("%d", entry.TotalTests)
		passed := fmt.Sprintf("%d", entry.Passed)
		failed := fmt.Sprintf("%d", entry.Failed)
		
		rows = append(rows, table.Row{
			timestamp,
			spec,
			tests,
			passed,
			failed,
			entry.Duration,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Bold(false)
	t.SetStyles(s)

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		MarginTop(1).
		Render("↑/↓: Navigate | Enter: Replay selected test | Esc: Return to results")

	return title + "\n\n" + t.View() + "\n\n" + instructions
}

// ViewCustomRequest renders the custom request screen
func ViewCustomRequest(m models.Model) string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		MarginBottom(1).
		Render("✏️  Custom API Request")

	crm := m.CustomRequestModel

	var content string

	switch crm.Step {
	case 0: // Method input
		stepTitle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#4ECDC4")).
			Render("Step 1/4: HTTP Method")
		
		input := crm.MethodInput.View()
		hint := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render("Enter HTTP method (GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS)")
		
		content = fmt.Sprintf("%s\n\n%s\n%s", stepTitle, input, hint)

	case 1: // Endpoint input
		stepTitle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#4ECDC4")).
			Render("Step 2/4: Endpoint URL")
		
		methodInfo := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render(fmt.Sprintf("Method: %s", crm.Request.Method))
		
		input := crm.EndpointInput.View()
		hint := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render("Enter full URL (e.g., https://api.example.com/users)")
		
		content = fmt.Sprintf("%s\n%s\n\n%s\n%s", stepTitle, methodInfo, input, hint)

	case 2: // Headers input
		stepTitle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#4ECDC4")).
			Render("Step 3/4: Request Headers (Optional)")
		
		methodInfo := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render(fmt.Sprintf("Method: %s | Endpoint: %s", crm.Request.Method, crm.Request.Endpoint))
		
		var headersList string
		if len(crm.Request.Headers) > 0 {
			headersList = "\nCurrent Headers:\n"
			for k, v := range crm.Request.Headers {
				headersList += fmt.Sprintf("  %s: %s\n", k, v)
			}
		}
		
		var inputs string
		if crm.HeaderKeyInput.Focused() {
			inputs = "Header Key:\n" + crm.HeaderKeyInput.View()
		} else {
			inputs = fmt.Sprintf("Header Key: %s\n\nHeader Value:\n%s", crm.HeaderKeyInput.Value(), crm.HeaderValueInput.View())
		}
		
		hint := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render("\nEnter header name, then value. Press Enter on empty key to continue to body.")
		
		content = fmt.Sprintf("%s\n%s%s\n\n%s%s", stepTitle, methodInfo, headersList, inputs, hint)

	case 3: // Body input
		stepTitle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#4ECDC4")).
			Render("Step 4/4: Request Body (Optional)")
		
		methodInfo := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render(fmt.Sprintf("Method: %s | Endpoint: %s", crm.Request.Method, crm.Request.Endpoint))
		
		var headersList string
		if len(crm.Request.Headers) > 0 {
			headersList = "\nHeaders:\n"
			for k, v := range crm.Request.Headers {
				headersList += fmt.Sprintf("  %s: %s\n", k, v)
			}
		}
		
		input := crm.BodyInput.View()
		hint := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render("\nEnter JSON body (or leave empty). Press Enter to execute request.")
		
		content = fmt.Sprintf("%s\n%s%s\n\nBody:\n%s%s", stepTitle, methodInfo, headersList, input, hint)

	case 4: // Executing
		statusMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Render(fmt.Sprintf("\n%s Sending %s request to %s...\n", crm.Spinner.View(), crm.Request.Method, crm.Request.Endpoint))
		content = statusMsg

	case 5: // Results
		if crm.Result != nil {
			resultStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#00FF00"))
			if crm.Result.Status != "200" && crm.Result.Status != "201" && crm.Result.Status != "204" {
				resultStyle = resultStyle.Foreground(lipgloss.Color("#FF6B6B"))
			}
			
			statusLine := resultStyle.Render(fmt.Sprintf("✓ %s %s", crm.Result.Status, crm.Result.Method))
			endpoint := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888")).
				Render(fmt.Sprintf("Endpoint: %s", crm.Result.Endpoint))
			duration := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888")).
				Render(fmt.Sprintf("Duration: %s", crm.Result.Duration))
			
			message := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4ECDC4")).
				MarginTop(1).
				Render(crm.Result.Message)
			
			content = fmt.Sprintf("%s\n%s\n%s\n%s", statusLine, endpoint, duration, message)
		}
	}

	// Error display
	var errorMsg string
	if crm.Err != nil {
		errorMsg = "\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true).
			Render(fmt.Sprintf("❌ Error: %v", crm.Err))
	}

	// Instructions
	var instructions string
	if crm.Step < 4 {
		instructions = "\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render("Enter: Next | Esc: Cancel and return to menu")
	} else if crm.Step == 5 {
		instructions = "\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render("Enter/Esc: Return to menu")
	}

	return title + "\n\n" + content + errorMsg + instructions
}

// ViewEndpointSelector renders the endpoint selector screen
func ViewEndpointSelector(m models.Model) string {
	esm := m.EndpointSelectorModel

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		MarginBottom(1).
		Render("🔍 Select Endpoints to Test")

	// Error display
	if esm.Err != nil {
		errorMsg := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true).
			Render(fmt.Sprintf("❌ Error: %v", esm.Err))
		instructions := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render("\nEsc: Return to menu")
		return title + "\n\n" + errorMsg + instructions
	}

	// Not ready (loading)
	if !esm.Ready {
		loading := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Render("Loading endpoints from spec...")
		return title + "\n\n" + loading
	}

	// No endpoints found
	if len(esm.AllEndpoints) == 0 {
		noEndpoints := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Render("No endpoints found in OpenAPI spec")
		instructions := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render("\nEsc: Return to menu")
		return title + "\n\n" + noEndpoints + instructions
	}

	// Search box
	searchLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4ECDC4")).
		Bold(true).
		Render("Search: ")
	searchBox := searchLabel + esm.SearchInput.View()

	// Selected count
	selectedCount := 0
	for _, ep := range esm.AllEndpoints {
		if ep.Selected {
			selectedCount++
		}
	}
	countStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		MarginTop(1)
	countText := countStyle.Render(fmt.Sprintf("Selected: %d/%d endpoints", selectedCount, len(esm.AllEndpoints)))

	// Filtered count (if different)
	var filterInfo string
	if len(esm.FilteredEndpoints) < len(esm.AllEndpoints) {
		filterInfo = countStyle.Render(fmt.Sprintf(" | Showing: %d", len(esm.FilteredEndpoints)))
	}

	// Endpoint list
	var listItems []string
	visibleHeight := 15 // Number of visible items
	endpoints := esm.FilteredEndpoints
	if len(endpoints) == 0 {
		endpoints = esm.AllEndpoints
	}

	// Calculate visible range with scrolling
	start := esm.Offset
	end := start + visibleHeight
	if end > len(endpoints) {
		end = len(endpoints)
	}

	for i := start; i < end; i++ {
		ep := endpoints[i]
		
		// Checkbox
		checkbox := "[ ]"
		if ep.Selected {
			checkbox = "[✓]"
		}

		// Cursor
		cursor := "  "
		if i == esm.Cursor {
			cursor = "▶ "
		}

		// Method with color
		var methodStyle lipgloss.Style
		switch ep.Method {
		case "GET":
			methodStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ECDC4"))
		case "POST":
			methodStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#95E1D3"))
		case "PUT":
			methodStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD93D"))
		case "PATCH":
			methodStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FAB95B"))
		case "DELETE":
			methodStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
		default:
			methodStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888"))
		}
		method := methodStyle.Render(fmt.Sprintf("%-7s", ep.Method))

		// Path
		path := ep.Path

		// Summary (optional)
		summary := ""
		if ep.Summary != "" {
			summary = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888")).
				Render(fmt.Sprintf(" - %s", ep.Summary))
			// Truncate if too long
			if len(summary) > 60 {
				summary = summary[:57] + "..."
			}
		}

		// Tags (optional)
		tags := ""
		if len(ep.Tags) > 0 {
			tags = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9B59B6")).
				Render(fmt.Sprintf(" [%s]", strings.Join(ep.Tags, ", ")))
		}

		// Build line
		line := fmt.Sprintf("%s%s %s %s%s%s", cursor, checkbox, method, path, summary, tags)
		
		// Highlight selected line
		if i == esm.Cursor {
			line = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FF6B6B")).
				Render(line)
		}

		listItems = append(listItems, line)
	}

	list := strings.Join(listItems, "\n")

	// Scroll indicators
	var scrollIndicator string
	if start > 0 {
		scrollIndicator += lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render("▲ More above\n")
	}
	if end < len(endpoints) {
		scrollIndicator += lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render("\n▼ More below")
	}

	// Instructions
	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		MarginTop(1).
		Render("↑/↓: Navigate | Space: Toggle | a: Select All | d: Deselect All | Enter: Test Selected | Esc: Cancel")

	return title + "\n\n" + searchBox + "\n" + countText + filterInfo + "\n\n" + scrollIndicator + list + scrollIndicator + "\n" + instructions
}

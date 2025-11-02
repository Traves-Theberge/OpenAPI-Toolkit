// Package main implements a professional Terminal User Interface (TUI) for OpenAPI specification
// validation and API testing. It uses the Charm Bracelet Bubble Tea framework for reactive
// terminal applications and Lip Gloss for beautiful styling.
//
// Architecture:
// - Single Bubble Tea program with multiple screen states
// - Screen-based navigation (menu → help/validate/test)
// - Embedded models for each feature (validation, testing)
// - Dynamic borders that adapt to terminal size
// - Async operations with spinners and progress indicators
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/getkin/kin-openapi/openapi3"
)

// screen represents the different UI states/screens in the application
type screen int

const (
	menuScreen screen = iota // Main menu with navigation options
	helpScreen               // Help/documentation screen
	validateScreen           // OpenAPI spec validation screen
	testScreen               // API endpoint testing screen
)

// errMsg is a type alias for error messages used in Bubble Tea messaging
type errMsg error

// validateModel handles the state for OpenAPI specification validation
type validateModel struct {
	textInput textinput.Model // Text input for file path
	err       error           // Validation error (if any)
	result    string          // Validation result message
	done      bool            // Whether validation is complete
}

// testResultMsg is a Bubble Tea message containing test results or errors
type testResultMsg struct {
	results []testResult // Array of individual test results
	err     error        // Testing error (if any)
}

// testModel manages the multi-step API testing workflow
type testModel struct {
	step      int                // Current step: 0=file input, 1=URL input, 2=testing, 3=results
	specInput textinput.Model    // Text input for OpenAPI spec file path
	urlInput  textinput.Model    // Text input for base API URL
	spinner   spinner.Model      // Loading spinner for async operations
	table     table.Model        // Results table display
	err       error              // Testing error (if any)
	result    string             // Result message
	done      bool               // Whether testing is complete
	testing   bool               // Whether testing is currently in progress
	results   []testResult       // Array of test results
}

// testResult represents a single API endpoint test result
type testResult struct {
	method   string // HTTP method (GET, POST, etc.)
	endpoint string // API endpoint path
	status   string // HTTP status code or "ERR"
	message  string // Success message or error details
}

// validationResult contains detailed validation information for a response
type validationResult struct {
	valid          bool     // Whether the response is valid
	statusValid    bool     // Whether status code matches spec
	contentType    string   // Actual content type
	schemaErrors   []string // Schema validation errors
	expectedStatus string   // Expected status code(s)
}

// model is the main application state, containing all sub-models and UI state
type model struct {
	cursor       int           // Currently selected menu item (0-3)
	choice       string        // Selected menu choice (legacy, may be removed)
	screen       screen        // Current screen being displayed
	width        int           // Terminal width (updated via WindowSizeMsg)
	height       int           // Terminal height (updated via WindowSizeMsg)
	validateModel validateModel // Embedded validation model
	testModel     testModel     // Embedded testing model
}

func initialModel() model {
	// Initialize the main application model with default values
	return model{
		cursor:         0,                    // Start cursor at first menu item
		screen:         menuScreen,           // Begin with main menu
		width:          80,                   // Default terminal width
		height:         24,                   // Default terminal height
		validateModel:  initialValidateModel(), // Initialize validation sub-model
		testModel:      initialTestModel(),     // Initialize testing sub-model
	}
}

// Init returns the initial command to run when the program starts
// Starts the spinner animation for the testing screen
func (m model) Init() tea.Cmd {
	return m.testModel.spinner.Tick
}

// Update handles all incoming messages and updates the model accordingly
// Routes messages to appropriate screen-specific update handlers
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle terminal resize - update dimensions for responsive layout
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		// Route key messages to current screen's update handler
		switch m.screen {
		case menuScreen:
			return m.updateMenu(msg)
		case helpScreen:
			return m.updateHelp(msg)
		case validateScreen:
			return m.updateValidate(msg)
		case testScreen:
			return m.updateTest(msg)
		}
	case testResultMsg:
		// Handle async test results in test screen
		if m.screen == testScreen {
			return m.updateTest(msg)
		}
	}
	return m, nil
}

func (m model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		// Quit the application
		return m, tea.Quit
	case "up", "k":
		// Navigate up in menu (vim-style 'k' key)
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		// Navigate down in menu (vim-style 'j' key)
		if m.cursor < 3 {
			m.cursor++
		}
	case "h", "?":
		// Show help screen
		m.screen = helpScreen
		return m, nil
	case "enter":
		// Select current menu item
		switch m.cursor {
		case 0:
			// Navigate to validation screen
			m.screen = validateScreen
			return m, nil
		case 1:
			// Navigate to testing screen
			m.screen = testScreen
			return m, nil
		case 2:
			// Navigate to help screen
			m.screen = helpScreen
			return m, nil
		case 3:
			// Quit application
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) updateHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "h", "?":
		// Return to menu from help screen
		m.screen = menuScreen
		return m, nil
	}
	return m, nil
}

func (m model) updateValidate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// Execute validation when user presses Enter
			filePath := m.validateModel.textInput.Value()
			if filePath == "" {
				m.validateModel.err = fmt.Errorf("file path cannot be empty")
				return m, nil
			}
			result, err := validateSpec(filePath)
			if err != nil {
				m.validateModel.err = err
				return m, nil
			}
			m.validateModel.result = result
			m.validateModel.done = true
			return m, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			// Return to menu and reset validation state
			m.screen = menuScreen
			m.validateModel = initialValidateModel()
			return m, nil
		}
	}

	// Handle text input for all messages (typing, cursor movement, etc.)
	m.validateModel.textInput, cmd = m.validateModel.textInput.Update(msg)
	return m, cmd
}

func (m model) updateTest(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Multi-step testing workflow: spec file -> base URL -> testing -> results
	switch m.testModel.step {
	case 0: // Spec file input
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				// Validate spec file path and advance to URL input
				if m.testModel.specInput.Value() == "" {
					m.testModel.err = fmt.Errorf("spec file path cannot be empty")
					return m, nil
				}
				m.testModel.step = 1
				m.testModel.urlInput.Focus()
				return m, nil
			case tea.KeyCtrlC, tea.KeyEsc:
				// Return to menu and reset test state
				m.screen = menuScreen
				m.testModel = initialTestModel()
				return m, nil
			}
			// Handle text input for spec file path
			m.testModel.specInput, cmd = m.testModel.specInput.Update(msg)
		case testResultMsg:
			// Handle async test results (shouldn't happen in step 0, but defensive)
			m.testModel.results = msg.results
			m.testModel.err = msg.err
			m.testModel.step = 3
			m.testModel.testing = false
			return m, nil
		}
	case 1: // Base URL input
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				// Validate base URL and start testing
				if m.testModel.urlInput.Value() == "" {
					m.testModel.err = fmt.Errorf("base URL cannot be empty")
					return m, nil
				}
				m.testModel.step = 2
				m.testModel.testing = true
				// Start async testing command
				return m, runTestCmd(m.testModel.specInput.Value(), m.testModel.urlInput.Value())
			case tea.KeyCtrlC, tea.KeyEsc:
				// Return to menu and reset test state
				m.screen = menuScreen
				m.testModel = initialTestModel()
				return m, nil
			}
			// Handle text input for base URL
			m.testModel.urlInput, cmd = m.testModel.urlInput.Update(msg)
		}
	case 2: // Testing in progress
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				// Cancel testing and return to menu
				m.screen = menuScreen
				m.testModel = initialTestModel()
				return m, nil
			}
		case testResultMsg:
			// Handle completed test results
			m.testModel.results = msg.results
			m.testModel.err = msg.err
			m.testModel.step = 3
			m.testModel.testing = false
			return m, nil
		}
		// Update spinner animation
		m.testModel.spinner, cmd = m.testModel.spinner.Update(msg)
	case 3: // Results display
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
				// Return to menu and reset test state
				m.screen = menuScreen
				m.testModel = initialTestModel()
				return m, nil
			}
		}
		// Handle table navigation and selection
		m.testModel.table, cmd = m.testModel.table.Update(msg)
	}

	return m, cmd
}

// View renders the current screen based on the application state
// Routes to appropriate view method based on current screen
func (m model) View() string {
	switch m.screen {
	case menuScreen:
		return m.viewMenu()
	case helpScreen:
		return m.viewHelp()
	case validateScreen:
		return m.viewValidate()
	case testScreen:
		return m.viewTest()
	default:
		return ""
	}
}

func (m model) viewMenu() string {
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

		if m.cursor == i {
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

	// Status bar with navigation hints
	statusBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		Padding(0, 1).
		MarginTop(1).
		Render("Press h or ? for help • ↑↓/jk to navigate • Enter to select")

	// Combine sections vertically centered
	content := lipgloss.JoinVertical(lipgloss.Center, title, menu, statusBar)

	// Center the entire content with dynamic sizing and border
	borderedContent := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4ECDC4")).
		Padding(1, 2).
		Render(content)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		borderedContent,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#333")),
	)
}

func (m model) viewHelp() string {
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
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		borderedContent,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#333")),
	)
}

func (m model) viewValidate() string {
	var content string

	// Show validation results if completed
	if m.validateModel.done {
		if m.validateModel.err != nil {
			// Display validation error in red
			content = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B")).
				Bold(true).
				Render(fmt.Sprintf("❌ Validation Failed:\n%s", m.validateModel.err.Error()))
		} else {
			// Display success message in green
			content = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4ECDC4")).
				Bold(true).
				Render(m.validateModel.result)
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
			Render("> " + m.validateModel.textInput.View())

		if m.validateModel.err != nil {
			// Show input error
			content = input + "\n\n" + lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B")).
				Render(fmt.Sprintf("Error: %s", m.validateModel.err.Error()))
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
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		borderedContent,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#333")),
	)
}

func (m model) viewTest() string {
	var content string

	// Multi-step testing workflow UI
	switch m.testModel.step {
	case 0: // Spec file input
		input := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4ECDC4")).
			Padding(1, 2).
			Render("> " + m.testModel.specInput.View())

		if m.testModel.err != nil {
			// Show input error for spec file
			content = input + "\n\n" + lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B")).
				Render(fmt.Sprintf("Error: %s", m.testModel.err.Error()))
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
			Render("> " + m.testModel.urlInput.View())

		if m.testModel.err != nil {
			// Show input error for base URL
			content = input + "\n\n" + lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B")).
				Render(fmt.Sprintf("Error: %s", m.testModel.err.Error()))
		} else {
			// Show base URL input instructions
			content = input + "\n\n" + lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888")).
				Render("Enter base URL (e.g., https://api.example.com) and press Enter")
		}
	case 2: // Testing in progress
		// Show animated spinner with testing message
		spinnerView := m.testModel.spinner.View() + " Testing API endpoints..."
		content = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true).
			Render(spinnerView) + "\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render("Press Ctrl+C to cancel")
	case 3: // Results display
		if m.testModel.err != nil {
			// Show testing error
			content = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B")).
				Bold(true).
				Render(fmt.Sprintf("❌ Testing Failed:\n%s", m.testModel.err.Error()))
		} else {
			// Populate table with test results
			var rows []table.Row
			for _, r := range m.testModel.results {
				rows = append(rows, table.Row{r.method, r.endpoint, r.status, r.message})
			}
			m.testModel.table.SetRows(rows)

			// Show success message and results table
			content = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4ECDC4")).
				Bold(true).
				Render("✅ Testing Complete!") + "\n\n" + m.testModel.table.View()
		}
		// Add exit instruction
		content += "\n\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888")).
			Render("Press Enter to return to menu")
	}

	// Center the entire content with dynamic sizing and border
	borderedContent := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4ECDC4")).
		Padding(1, 2).
		Render(content)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		borderedContent,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#333")),
	)
}

func initialValidateModel() validateModel {
	// Create and configure text input for spec file path
	ti := textinput.New()
	ti.Placeholder = "Path to OpenAPI spec file (e.g., openapi.yaml)"
	ti.Focus() // Start with focus for immediate typing
	ti.CharLimit = 156
	ti.Width = 60

	return validateModel{
		textInput: ti,
		err:       nil,
	}
}

func initialTestModel() testModel {
	// Configure spec file input
	specTi := textinput.New()
	specTi.Placeholder = "Path to OpenAPI spec file (e.g., openapi.yaml)"
	specTi.Focus() // Start with focus for immediate typing
	specTi.CharLimit = 156
	specTi.Width = 60

	// Configure base URL input
	urlTi := textinput.New()
	urlTi.Placeholder = "Base URL (e.g., https://api.example.com)"
	urlTi.CharLimit = 156
	urlTi.Width = 60

	// Configure spinner for testing progress
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ECDC4"))

	// Configure results table
	columns := []table.Column{
		{Title: "Method", Width: 8},
		{Title: "Endpoint", Width: 40},
		{Title: "Status", Width: 10},
		{Title: "Message", Width: 30},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(false), // Read-only table
		table.WithHeight(10),
	)

	// Configure table styles for professional appearance
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ECDC4"))
	t.SetStyles(table.Styles{
		Header: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Bold(true),
		Cell: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")),
		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true),
	})

	return testModel{
		specInput: specTi,
		urlInput:  urlTi,
		spinner:   s,
		table:     t,
	}
}

// runTestCmd creates a Bubble Tea command to run tests asynchronously
// Returns a command that executes runTests in a goroutine and sends results via message
func runTestCmd(specPath, baseURL string) tea.Cmd {
	return func() tea.Msg {
		results, err := runTests(specPath, baseURL)
		return testResultMsg{results: results, err: err}
	}
}

// validateSpec validates an OpenAPI specification file
// Returns success message or detailed error
func validateSpec(filePath string) (string, error) {
	// Load OpenAPI document with external references allowed
	loader := &openapi3.Loader{IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to load spec: %v", err)
	}

	// Validate the loaded document
	err = doc.Validate(loader.Context)
	if err != nil {
		return "", fmt.Errorf("validation failed: %v", err)
	}

	return "OpenAPI spec is valid! 🎉", nil
}

// replacePlaceholders replaces path parameters like {id} with sample values
func replacePlaceholders(path string) string {
	// Replace {param} with sensible defaults
	re := regexp.MustCompile(`\{[^}]+\}`)
	return re.ReplaceAllString(path, "1")
}

// buildQueryParams constructs query parameters from operation parameters
func buildQueryParams(operation *openapi3.Operation) string {
	if operation == nil || operation.Parameters == nil {
		return ""
	}

	var params []string
	for _, paramRef := range operation.Parameters {
		param := paramRef.Value
		if param == nil || param.In != "query" {
			continue
		}

		// Generate sample value based on schema
		value := "1" // Default
		if param.Schema != nil && param.Schema.Value != nil {
			schema := param.Schema.Value
			// Check if type contains specific values
			if schema.Type.Is("string") {
				if len(schema.Enum) > 0 {
					value = fmt.Sprintf("%v", schema.Enum[0])
				} else if schema.Example != nil {
					value = fmt.Sprintf("%v", schema.Example)
				} else {
					value = "test"
				}
			} else if schema.Type.Is("integer") || schema.Type.Is("number") {
				if schema.Example != nil {
					value = fmt.Sprintf("%v", schema.Example)
				} else {
					value = "1"
				}
			} else if schema.Type.Is("boolean") {
				value = "true"
			} else if schema.Type.Is("array") {
				value = "1,2,3" // Simple array representation
			}
		}

		params = append(params, fmt.Sprintf("%s=%s", param.Name, value))
	}

	if len(params) == 0 {
		return ""
	}
	return "?" + strings.Join(params, "&")
}

// generateRequestBody creates a sample JSON request body from an OpenAPI schema
// Generates realistic sample data based on schema properties, types, and examples
func generateRequestBody(operation *openapi3.Operation) ([]byte, error) {
	if operation == nil || operation.RequestBody == nil {
		return nil, nil
	}

	// Get the request body content for JSON
	requestBody := operation.RequestBody.Value
	if requestBody == nil || requestBody.Content == nil {
		return nil, nil
	}

	// Look for JSON content type
	jsonContent := requestBody.Content.Get("application/json")
	if jsonContent == nil || jsonContent.Schema == nil || jsonContent.Schema.Value == nil {
		return nil, nil
	}

	schema := jsonContent.Schema.Value

	// Generate sample data from schema
	sample := generateSampleFromSchema(schema)
	
	// Marshal to JSON
	jsonData, err := json.Marshal(sample)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	return jsonData, nil
}

// generateSampleFromSchema recursively generates sample data from an OpenAPI schema
func generateSampleFromSchema(schema *openapi3.Schema) interface{} {
	if schema == nil {
		return nil
	}

	// Use example if available
	if schema.Example != nil {
		return schema.Example
	}

	// Use default if available
	if schema.Default != nil {
		return schema.Default
	}

	// Generate based on type
	if schema.Type.Is("object") {
		obj := make(map[string]interface{})
		for propName, propRef := range schema.Properties {
			if propRef != nil && propRef.Value != nil {
				obj[propName] = generateSampleFromSchema(propRef.Value)
			}
		}
		return obj
	}

	if schema.Type.Is("array") {
		if schema.Items != nil && schema.Items.Value != nil {
			// Generate a single-item array
			return []interface{}{generateSampleFromSchema(schema.Items.Value)}
		}
		return []interface{}{}
	}

	if schema.Type.Is("string") {
		if len(schema.Enum) > 0 {
			return schema.Enum[0]
		}
		if schema.Format == "email" {
			return "user@example.com"
		}
		if schema.Format == "uri" || schema.Format == "url" {
			return "https://example.com"
		}
		if schema.Format == "date" {
			return "2024-01-01"
		}
		if schema.Format == "date-time" {
			return "2024-01-01T00:00:00Z"
		}
		return "sample"
	}

	if schema.Type.Is("integer") {
		if schema.Min != nil {
			return int(*schema.Min)
		}
		return 1
	}

	if schema.Type.Is("number") {
		if schema.Min != nil {
			return *schema.Min
		}
		return 1.0
	}

	if schema.Type.Is("boolean") {
		return true
	}

	// Default fallback
	return nil
}

// runTests executes API tests against endpoints defined in OpenAPI spec
// Tests each endpoint with a simple request and records results
func runTests(specPath, baseURL string) ([]testResult, error) {
	// Load and validate the OpenAPI spec
	loader := &openapi3.Loader{IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load spec: %v", err)
	}

	var results []testResult

	// Iterate through all paths and operations in the spec
	if doc.Paths != nil {
		for path, pathItem := range doc.Paths.Map() {
			// Iterate through all paths and operations in the spec
			for method, operation := range pathItem.Operations() {
				// Construct full endpoint URL with placeholder replacement
				endpoint := baseURL + replacePlaceholders(path)
				
				// Add query parameters if defined
				queryParams := buildQueryParams(operation)
				endpoint += queryParams

				// Generate request body if needed
				var requestBody []byte
				if strings.ToUpper(method) == "POST" || strings.ToUpper(method) == "PUT" || strings.ToUpper(method) == "PATCH" {
					requestBody, err = generateRequestBody(operation)
					if err != nil {
						// Log error but continue testing
						results = append(results, testResult{
							method:   method,
							endpoint: path,
							status:   "ERR",
							message:  fmt.Sprintf("Failed to generate request body: %v", err),
						})
						continue
					}
				}

				// Test the endpoint and record result
				status, resp, err := testEndpoint(method, endpoint, requestBody)
				message := "OK"
				if err != nil {
					message = err.Error()
				} else if resp != nil {
					// Validate response against spec
					validation := validateResponse(resp, operation, status)
					
					// Close response body after validation
					if resp.Body != nil {
						io.Copy(io.Discard, resp.Body) // Drain body
						resp.Body.Close()
					}
					
					if !validation.valid {
						message = "Response validation failed"
						if len(validation.schemaErrors) > 0 {
							message = validation.schemaErrors[0] // Show first error
						}
					} else if validation.statusValid {
						message = "OK (validated)"
					}
				}

				// Format status for display
				statusStr := fmt.Sprintf("%d", status)
				if err != nil {
					statusStr = "ERR"
				}

				// Add result to collection
				results = append(results, testResult{
					method:   method,
					endpoint: path,
					status:   statusStr,
					message:  message,
				})
			}
		}
	}

	return results, nil
}

// validateResponse validates an HTTP response against OpenAPI spec
// Returns validation result with detailed error information
func validateResponse(resp *http.Response, operation *openapi3.Operation, statusCode int) validationResult {
	result := validationResult{
		valid:       true,
		statusValid: false,
		contentType: resp.Header.Get("Content-Type"),
	}

	if operation == nil || operation.Responses == nil {
		// No spec to validate against - mark as valid
		return result
	}

	// Check if status code is defined in spec
	statusStr := fmt.Sprintf("%d", statusCode)
	response := operation.Responses.Status(statusCode)
	
	if response != nil {
		// Found exact status match
		result.statusValid = true
		result.expectedStatus = statusStr
	} else {
		// Check if there's an explicit "default" response
		respMap := operation.Responses.Map()
		if defaultResp, hasDefault := respMap["default"]; hasDefault && defaultResp != nil {
			response = defaultResp
			result.statusValid = true
			result.expectedStatus = "default"
		} else {
			// Status code not defined in spec and no explicit default
			result.statusValid = false
			result.valid = false
			result.schemaErrors = append(result.schemaErrors, 
				fmt.Sprintf("status %d not defined in spec", statusCode))
			return result
		}
	}

	// Validate content type if response has content
	if response.Value != nil && response.Value.Content != nil {
		// Extract base content type (ignore charset, etc.)
		contentType := strings.Split(result.contentType, ";")[0]
		contentType = strings.TrimSpace(contentType)
		
		// Check if content type is defined in spec
		mediaType := response.Value.Content.Get(contentType)
		if mediaType == nil {
			// Try common alternatives
			if contentType == "" {
				contentType = "application/json" // Default assumption
				mediaType = response.Value.Content.Get(contentType)
			}
		}

		if mediaType == nil && len(response.Value.Content) > 0 {
			result.valid = false
			result.schemaErrors = append(result.schemaErrors,
				fmt.Sprintf("content-type '%s' not defined in spec", result.contentType))
		}

		// TODO: Add JSON schema validation against response body
		// This would require parsing the response body and validating against mediaType.Schema
		// For now, we just validate status code and content type
	}

	return result
}

// testEndpoint performs an HTTP request to test an API endpoint
// Supports GET, POST, PUT, PATCH, DELETE methods with optional request bodies
// Returns status code, response object, and error
func testEndpoint(method, url string, body []byte) (int, *http.Response, error) {
	var req *http.Request
	var err error

	// Create request based on HTTP method
	method = strings.ToUpper(method)
	
	if body != nil && len(body) > 0 {
		// Create request with body
		req, err = http.NewRequest(method, url, bytes.NewReader(body))
		if err != nil {
			return 0, nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		// Create request without body
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return 0, nil, err
		}
	}

	// Execute request with timeout to prevent hanging
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, resp, nil
}

// main initializes and runs the Bubble Tea TUI program
// Uses alt screen mode for full terminal control
func main() {
	// Create and run the Bubble Tea program with initial model
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Program handles all user interaction internally - no external processing needed
}
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
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v3"
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
	method   string    // HTTP method (GET, POST, etc.)
	endpoint string    // API endpoint path
	status   string    // HTTP status code or "ERR"
	message  string    // Success message or error details
	duration time.Duration // Request duration
	logEntry *logEntry // Detailed log information (optional)
}

// logEntry contains detailed request/response information for debugging
type logEntry struct {
	requestURL     string            // Full request URL
	requestHeaders map[string]string // Request headers sent
	requestBody    string            // Request body (if any)
	responseHeaders map[string]string // Response headers received
	responseBody    string            // Response body (truncated if large)
	duration       time.Duration     // Total request duration
	timestamp      time.Time         // When the request was made
}

// validationResult contains detailed validation information for a response
type validationResult struct {
	valid          bool     // Whether the response is valid
	statusValid    bool     // Whether status code matches spec
	contentType    string   // Actual content type
	schemaErrors   []string // Schema validation errors
	expectedStatus string   // Expected status code(s)
}

// authConfig contains authentication configuration for API requests
type authConfig struct {
	authType string // Type: "bearer", "apiKey", "basic", "none"
	token    string // For bearer tokens or API keys
	apiKeyIn string // For API keys: "header" or "query"
	apiKeyName string // Header/query parameter name for API key
	username string // For basic auth
	password string // For basic auth
}

// config represents the application configuration that persists between sessions
type config struct {
	BaseURL     string      `yaml:"baseURL"`     // Default base URL for API testing
	SpecPath    string      `yaml:"specPath"`    // Default OpenAPI spec file path
	VerboseMode bool        `yaml:"verboseMode"` // Default verbose logging state
	Auth        *authConfig `yaml:"auth,omitempty"` // Optional authentication configuration
}

// configFile represents the on-disk configuration structure
type configFile struct {
	BaseURL     string `yaml:"baseURL"`
	SpecPath    string `yaml:"specPath"`
	VerboseMode bool   `yaml:"verboseMode"`
	Auth        *struct {
		Type       string `yaml:"type"`
		Token      string `yaml:"token,omitempty"`
		APIKeyIn   string `yaml:"apiKeyIn,omitempty"`
		APIKeyName string `yaml:"apiKeyName,omitempty"`
		Username   string `yaml:"username,omitempty"`
		Password   string `yaml:"password,omitempty"`
	} `yaml:"auth,omitempty"`
}

// enhancedError provides detailed error information with actionable suggestions
type enhancedError struct {
	title       string   // Short error title
	description string   // Detailed error description
	suggestions []string // Actionable suggestions for fixing
	original    error    // Original error for reference
}

// Error implements the error interface
func (e *enhancedError) Error() string {
	msg := fmt.Sprintf("%s: %s", e.title, e.description)
	if len(e.suggestions) > 0 {
		msg += "\n\nSuggestions:"
		for _, s := range e.suggestions {
			msg += "\n  • " + s
		}
	}
	return msg
}

// formatEnhancedError creates a styled error message for display
func formatEnhancedError(err error) string {
	if enhanced, ok := err.(*enhancedError); ok {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")). // Red
			Bold(true)
		
		suggestionStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")) // Yellow
		
		msg := errorStyle.Render("❌ " + enhanced.title)
		msg += "\n\n" + enhanced.description
		
		if len(enhanced.suggestions) > 0 {
			msg += "\n\n" + suggestionStyle.Render("💡 Suggestions:")
			for _, s := range enhanced.suggestions {
				msg += "\n  • " + s
			}
		}
		return msg
	}
	
	// Fallback for regular errors
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Render("❌ Error: " + err.Error())
}

// enhanceFileError wraps file-related errors with helpful suggestions
func enhanceFileError(err error, filePath string) error {
	if err == nil {
		return nil
	}
	
	errStr := err.Error()
	
	// File not found
	if strings.Contains(errStr, "no such file") || strings.Contains(errStr, "cannot find") {
		return &enhancedError{
			title:       "File Not Found",
			description: fmt.Sprintf("Could not find the file: %s", filePath),
			suggestions: []string{
				"Check if the file path is correct",
				"Use an absolute path (e.g., /home/user/spec.yaml)",
				"Verify the file exists using 'ls' command",
				"Make sure you have read permissions for the file",
			},
			original: err,
		}
	}
	
	// Permission denied
	if strings.Contains(errStr, "permission denied") {
		return &enhancedError{
			title:       "Permission Denied",
			description: fmt.Sprintf("Cannot read file: %s", filePath),
			suggestions: []string{
				"Check file permissions with 'ls -l'",
				"Try running with appropriate permissions",
				"Make sure you own the file or have read access",
			},
			original: err,
		}
	}
	
	// Parse errors
	if strings.Contains(errStr, "yaml") || strings.Contains(errStr, "unmarshal") {
		return &enhancedError{
			title:       "Invalid File Format",
			description: "The file is not a valid OpenAPI specification",
			suggestions: []string{
				"Ensure the file is valid YAML or JSON",
				"Check for syntax errors (quotes, indentation)",
				"Validate YAML at https://www.yamllint.com/",
				"Make sure it's an OpenAPI 3.x specification",
			},
			original: err,
		}
	}
	
	return err
}

// enhanceNetworkError wraps network-related errors with helpful suggestions
func enhanceNetworkError(err error, url string) error {
	if err == nil {
		return nil
	}
	
	errStr := err.Error()
	
	// Connection refused
	if strings.Contains(errStr, "connection refused") {
		return &enhancedError{
			title:       "Connection Refused",
			description: fmt.Sprintf("Cannot connect to: %s", url),
			suggestions: []string{
				"Check if the server is running",
				"Verify the URL and port are correct",
				"Check firewall settings",
				"Try pinging the host to verify connectivity",
			},
			original: err,
		}
	}
	
	// Timeout
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded") {
		return &enhancedError{
			title:       "Request Timeout",
			description: "The server took too long to respond",
			suggestions: []string{
				"Check your internet connection",
				"The server might be overloaded - try again later",
				"Verify the URL points to the correct endpoint",
				"Check if the server is experiencing issues",
			},
			original: err,
		}
	}
	
	// DNS resolution failure
	if strings.Contains(errStr, "no such host") || strings.Contains(errStr, "dns") {
		return &enhancedError{
			title:       "DNS Resolution Failed",
			description: fmt.Sprintf("Cannot resolve hostname: %s", url),
			suggestions: []string{
				"Check if the URL is spelled correctly",
				"Verify your DNS settings",
				"Try using the IP address directly",
				"Check your internet connection",
			},
			original: err,
		}
	}
	
	// TLS/SSL errors
	if strings.Contains(errStr, "tls") || strings.Contains(errStr, "certificate") {
		return &enhancedError{
			title:       "TLS/SSL Error",
			description: "Cannot establish secure connection",
			suggestions: []string{
				"The server's SSL certificate might be invalid",
				"Check if the URL should use 'http' instead of 'https'",
				"Verify the server's certificate is up to date",
				"Try accessing the URL in a browser first",
			},
			original: err,
		}
	}
	
	return err
}

// enhanceValidationError wraps validation errors with helpful suggestions
func enhanceValidationError(err error) error {
	if err == nil {
		return nil
	}
	
	errStr := err.Error()
	
	// Missing required fields
	if strings.Contains(errStr, "required") {
		return &enhancedError{
			title:       "Validation Failed",
			description: "The OpenAPI specification is missing required fields",
			suggestions: []string{
				"Check that all required fields are present (openapi, info, paths)",
				"Verify the spec follows OpenAPI 3.x format",
				"Use a linter like Spectral to validate your spec",
				"See OpenAPI specification at https://spec.openapis.org/",
			},
			original: err,
		}
	}
	
	// Version mismatch
	if strings.Contains(errStr, "version") || strings.Contains(errStr, "openapi") {
		return &enhancedError{
			title:       "Unsupported Version",
			description: "This tool only supports OpenAPI 3.x specifications",
			suggestions: []string{
				"Check the 'openapi' field in your spec",
				"If using Swagger 2.0, convert it to OpenAPI 3.x",
				"Use https://converter.swagger.io/ for conversion",
				"Ensure 'openapi' field starts with '3.' (e.g., '3.0.3')",
			},
			original: err,
		}
	}
	
	return err
}

// getConfigPath returns the path to the configuration file
// Uses ~/.config/openapi-tui/config.yaml on Unix-like systems
// Uses %APPDATA%\openapi-tui\config.yaml on Windows
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	
	configDir := filepath.Join(homeDir, ".config", "openapi-tui")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	
	return filepath.Join(configDir, "config.yaml"), nil
}

// loadConfig loads configuration from the config file
// Returns default config if file doesn't exist
func loadConfig() config {
	cfg := config{
		VerboseMode: false,
	}
	
	configPath, err := getConfigPath()
	if err != nil {
		return cfg
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		// File doesn't exist or can't be read - return defaults
		return cfg
	}
	
	var fileConfig configFile
	if err := yaml.Unmarshal(data, &fileConfig); err != nil {
		// Invalid YAML - return defaults
		return cfg
	}
	
	// Convert fileConfig to config
	cfg.BaseURL = fileConfig.BaseURL
	cfg.SpecPath = fileConfig.SpecPath
	cfg.VerboseMode = fileConfig.VerboseMode
	
	if fileConfig.Auth != nil {
		cfg.Auth = &authConfig{
			authType:   fileConfig.Auth.Type,
			token:      fileConfig.Auth.Token,
			apiKeyIn:   fileConfig.Auth.APIKeyIn,
			apiKeyName: fileConfig.Auth.APIKeyName,
			username:   fileConfig.Auth.Username,
			password:   fileConfig.Auth.Password,
		}
	}
	
	return cfg
}

// saveConfig saves the current configuration to the config file
func saveConfig(cfg config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	
	// Convert config to fileConfig
	fileConfig := configFile{
		BaseURL:     cfg.BaseURL,
		SpecPath:    cfg.SpecPath,
		VerboseMode: cfg.VerboseMode,
	}
	
	if cfg.Auth != nil {
		fileConfig.Auth = &struct {
			Type       string `yaml:"type"`
			Token      string `yaml:"token,omitempty"`
			APIKeyIn   string `yaml:"apiKeyIn,omitempty"`
			APIKeyName string `yaml:"apiKeyName,omitempty"`
			Username   string `yaml:"username,omitempty"`
			Password   string `yaml:"password,omitempty"`
		}{
			Type:       cfg.Auth.authType,
			Token:      cfg.Auth.token,
			APIKeyIn:   cfg.Auth.apiKeyIn,
			APIKeyName: cfg.Auth.apiKeyName,
			Username:   cfg.Auth.username,
			Password:   cfg.Auth.password,
		}
	}
	
	data, err := yaml.Marshal(fileConfig)
	if err != nil {
		return err
	}
	
	return os.WriteFile(configPath, data, 0644)
}

// model is the main application state, containing all sub-models and UI state
type model struct {
	cursor       int           // Currently selected menu item (0-3)
	choice       string        // Selected menu choice (legacy, may be removed)
	screen       screen        // Current screen being displayed
	width        int           // Terminal width (updated via WindowSizeMsg)
	height       int           // Terminal height (updated via WindowSizeMsg)
	verboseMode  bool          // Whether verbose logging is enabled
	config       config        // Application configuration
	validateModel validateModel // Embedded validation model
	testModel     testModel     // Embedded testing model
}

func initialModel() model {
	// Load configuration from file
	cfg := loadConfig()
	
	// Initialize the main application model with default values and loaded config
	m := model{
		cursor:         0,                    // Start cursor at first menu item
		screen:         menuScreen,           // Begin with main menu
		width:          80,                   // Default terminal width
		height:         24,                   // Default terminal height
		verboseMode:    cfg.VerboseMode,      // Load verbose mode from config
		config:         cfg,                  // Store config
		validateModel:  initialValidateModel(), // Initialize validation sub-model
		testModel:      initialTestModel(),     // Initialize testing sub-model
	}
	
	// Pre-fill spec path and base URL if saved in config
	if cfg.SpecPath != "" {
		m.testModel.specInput.SetValue(cfg.SpecPath)
	}
	if cfg.BaseURL != "" {
		m.testModel.urlInput.SetValue(cfg.BaseURL)
	}
	
	return m
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
	case "v":
		// Toggle verbose mode and save to config
		m.verboseMode = !m.verboseMode
		m.config.VerboseMode = m.verboseMode
		saveConfig(m.config) // Save configuration change
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
				
				// Save spec path and base URL to config for next session
				m.config.SpecPath = m.testModel.specInput.Value()
				m.config.BaseURL = m.testModel.urlInput.Value()
				saveConfig(m.config)
				
				m.testModel.step = 2
				m.testModel.testing = true
				// Start async testing command
				// TODO: Collect auth configuration from user in future enhancement
				return m, runTestCmd(m.testModel.specInput.Value(), m.testModel.urlInput.Value(), nil, m.verboseMode)
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

	// Status indicators
	var statusIndicators []string
	if m.verboseMode {
		statusIndicators = append(statusIndicators, "🔍 Verbose: ON")
	}
	if m.config.BaseURL != "" || m.config.SpecPath != "" {
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
			// Display enhanced validation error with suggestions
			content = formatEnhancedError(m.validateModel.err)
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
			// Show enhanced input error with suggestions
			content = input + "\n\n" + formatEnhancedError(m.validateModel.err)
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
			// Show enhanced input error for spec file with suggestions
			content = input + "\n\n" + formatEnhancedError(m.testModel.err)
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
			// Show enhanced input error for base URL with suggestions
			content = input + "\n\n" + formatEnhancedError(m.testModel.err)
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
			// Show enhanced testing error with actionable suggestions
			content = formatEnhancedError(m.testModel.err)
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
func runTestCmd(specPath, baseURL string, auth *authConfig, verbose bool) tea.Cmd {
	return func() tea.Msg {
		results, err := runTests(specPath, baseURL, auth, verbose)
		return testResultMsg{results: results, err: err}
	}
}

// validateSpec validates an OpenAPI specification file
// Returns success message or detailed error with actionable suggestions
func validateSpec(filePath string) (string, error) {
	// Load OpenAPI document with external references allowed
	loader := &openapi3.Loader{IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(filePath)
	if err != nil {
		return "", enhanceFileError(err, filePath)
	}

	// Validate the loaded document
	err = doc.Validate(loader.Context)
	if err != nil {
		return "", enhanceValidationError(err)
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
// Accepts optional auth configuration and verbose flag for detailed logging
func runTests(specPath, baseURL string, auth *authConfig, verbose bool) ([]testResult, error) {
	// Load and validate the OpenAPI spec
	loader := &openapi3.Loader{IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile(specPath)
	if err != nil {
		return nil, enhanceFileError(err, specPath)
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
				startTime := time.Now()
				status, resp, logEntry, err := testEndpoint(method, endpoint, requestBody, auth, verbose)
				duration := time.Since(startTime)
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
					duration: duration,
					logEntry: logEntry,
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

// applyAuth applies authentication configuration to an HTTP request
func applyAuth(req *http.Request, auth *authConfig) {
	if auth == nil || auth.authType == "none" || auth.authType == "" {
		return
	}

	switch auth.authType {
	case "bearer":
		if auth.token != "" {
			req.Header.Set("Authorization", "Bearer "+auth.token)
		}
	case "apiKey":
		if auth.apiKeyName != "" && auth.token != "" {
			if auth.apiKeyIn == "header" {
				req.Header.Set(auth.apiKeyName, auth.token)
			} else if auth.apiKeyIn == "query" {
				// Add to query parameters
				q := req.URL.Query()
				q.Add(auth.apiKeyName, auth.token)
				req.URL.RawQuery = q.Encode()
			}
		}
	case "basic":
		if auth.username != "" {
			req.SetBasicAuth(auth.username, auth.password)
		}
	}
}

// testEndpoint performs an HTTP request to test an API endpoint
// Supports GET, POST, PUT, PATCH, DELETE methods with optional request bodies
// Returns status code, response object, log entry, and error
func testEndpoint(method, url string, body []byte, auth *authConfig, verbose bool) (int, *http.Response, *logEntry, error) {
	var req *http.Request
	var err error

	// Create request based on HTTP method
	method = strings.ToUpper(method)
	
	if body != nil && len(body) > 0 {
		// Create request with body
		req, err = http.NewRequest(method, url, bytes.NewReader(body))
		if err != nil {
			return 0, nil, nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		// Create request without body
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return 0, nil, nil, err
		}
	}

	// Apply authentication if configured
	applyAuth(req, auth)

	// Capture start time for duration measurement
	startTime := time.Now()

	// Execute request with timeout to prevent hanging
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		return 0, nil, nil, enhanceNetworkError(err, url)
	}

	// Create log entry if verbose mode is enabled
	var log *logEntry
	if verbose {
		log = &logEntry{
			requestURL:  url,
			duration:    duration,
			timestamp:   startTime,
			requestHeaders: make(map[string]string),
			responseHeaders: make(map[string]string),
		}

		// Capture request headers
		for k, v := range req.Header {
			if len(v) > 0 {
				log.requestHeaders[k] = v[0]
			}
		}

		// Capture request body
		if len(body) > 0 {
			log.requestBody = string(body)
			// Truncate if too large
			if len(log.requestBody) > 500 {
				log.requestBody = log.requestBody[:500] + "... (truncated)"
			}
		}

		// Capture response headers
		for k, v := range resp.Header {
			if len(v) > 0 {
				log.responseHeaders[k] = v[0]
			}
		}

		// Capture response body (read and restore)
		if resp.Body != nil {
			bodyBytes, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil {
				log.responseBody = string(bodyBytes)
				// Truncate if too large
				if len(log.responseBody) > 500 {
					log.responseBody = log.responseBody[:500] + "... (truncated)"
				}
				// Restore body for further processing
				resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
		}
	}

	return resp.StatusCode, resp, log, nil
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
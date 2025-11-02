// Package main implements a professional Terminal User Interface (TUI) for OpenAPI specification
// validation and API testing. It uses the Charm Bracelet Bubble Tea framework for reactive
// terminal applications and Lip Gloss for beautiful styling.
//
// Architecture:
// - Modular design with separate packages for models, views, testing, validation, etc.
// - Single Bubble Tea program with multiple screen states
// - Screen-based navigation (menu → help/validate/test)
// - Embedded models for each feature (validation, testing)
// - Dynamic borders that adapt to terminal size
// - Async operations with spinners and progress indicators
package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/config"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/errors"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/export"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/testing"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/ui"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/validation"
)

// model is a local wrapper around models.Model to implement tea.Model interface
type model struct {
	models.Model
}

// initialModel creates and initializes the main application model
// Loads configuration and prepares all sub-models
func initialModel() model {
	// Load configuration from file
	cfg := config.LoadConfig()

	// Load test run history
	history, err := models.LoadHistory()
	if err != nil {
		// If history can't be loaded, start with empty history
		// Error is logged but doesn't prevent app from starting
		history = &models.TestHistory{}
	}

	// Initialize the main application model with default values and loaded config
	m := model{
		Model: models.Model{
			Cursor:                0,
			Screen:                models.MenuScreen,
			Width:                 80,
			Height:                24,
			VerboseMode:           cfg.VerboseMode,
			Config:                cfg,
			ValidateModel:         ui.InitialValidateModel(),
			TestModel:             ui.InitialTestModel(),
			CustomRequestModel:    ui.InitialCustomRequestModel(),
			EndpointSelectorModel: ui.InitialEndpointSelectorModel(),
			History:               history,
			HistoryIndex:          0,
		},
	}

	// Pre-fill spec path and base URL if saved in config
	if cfg.SpecPath != "" {
		m.TestModel.SpecInput.SetValue(cfg.SpecPath)
	}
	if cfg.BaseURL != "" {
		m.TestModel.UrlInput.SetValue(cfg.BaseURL)
	}

	return m
}

// Init returns the initial command to run when the program starts
func (m model) Init() tea.Cmd {
	return m.TestModel.Spinner.Tick
}

// Update handles all incoming messages and updates the model accordingly
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch m.Screen {
		case models.MenuScreen:
			return m.updateMenu(msg)
		case models.HelpScreen:
			return m.updateHelp(msg)
		case models.ValidateScreen:
			return m.updateValidate(msg)
		case models.TestScreen:
			return m.updateTest(msg)
		case models.CustomRequestScreen:
			return m.updateCustomRequest(msg)
		case models.EndpointSelectorScreen:
			return m.updateEndpointSelector(msg)
		case models.HistoryScreen:
			return m.updateHistory(msg)
		}
	case testing.TestCompleteMsg:
		if m.Screen == models.TestScreen {
			return m.updateTest(msg)
		}
		if m.Screen == models.CustomRequestScreen {
			return m.updateCustomRequest(msg)
		}
	case testing.TestErrorMsg:
		if m.Screen == models.TestScreen {
			return m.updateTest(msg)
		}
		if m.Screen == models.CustomRequestScreen {
			return m.updateCustomRequest(msg)
		}
	}
	return m, nil
}

// updateMenu handles key events in the main menu screen
func (m model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < 6 {
			m.Cursor++
		}
	case "h", "?":
		m.Screen = models.HelpScreen
		return m, nil
	case "v":
		m.VerboseMode = !m.VerboseMode
		m.Config.VerboseMode = m.VerboseMode
		config.SaveConfig(m.Config)
		return m, nil
	case "enter":
		switch m.Cursor {
		case 0:
			m.Screen = models.ValidateScreen
			return m, nil
		case 1:
			// Test All Endpoints
			m.Screen = models.TestScreen
			return m, nil
		case 2:
			// Select & Test Endpoints - need spec path and base URL first
			m.Screen = models.TestScreen
			m.TestModel.Step = 0
			m.TestModel.SelectEndpoints = true  // Flag to show endpoint selector after step 1
			return m, nil
		case 3:
			m.Screen = models.CustomRequestScreen
			m.CustomRequestModel = ui.InitialCustomRequestModel()
			return m, nil
		case 4:
			m.Screen = models.HistoryScreen
			m.HistoryIndex = 0
			return m, nil
		case 5:
			m.Screen = models.HelpScreen
			return m, nil
		case 6:
			return m, tea.Quit
		}
	}
	return m, nil
}

// updateHelp handles key events in the help screen
func (m model) updateHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "h", "?":
		m.Screen = models.MenuScreen
		return m, nil
	}
	return m, nil
}

// updateValidate handles events in the validation screen
func (m model) updateValidate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.ValidateModel.Done {
				m.Screen = models.MenuScreen
				m.ValidateModel = ui.InitialValidateModel()
				return m, nil
			}

			filePath := m.ValidateModel.TextInput.Value()
			if filePath == "" {
				m.ValidateModel.Err = fmt.Errorf("file path cannot be empty")
				return m, nil
			}
			result, err := validation.ValidateSpec(filePath)
			if err != nil {
				m.ValidateModel.Err = err
				return m, nil
			}
			m.ValidateModel.Result = result
			m.ValidateModel.Done = true
			return m, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			m.Screen = models.MenuScreen
			m.ValidateModel = ui.InitialValidateModel()
			return m, nil
		}
	}

	m.ValidateModel.TextInput, cmd = m.ValidateModel.TextInput.Update(msg)
	return m, cmd
}

// updateTest handles events in the testing screen
func (m model) updateTest(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.TestModel.Step {
	case 0:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				if m.TestModel.SpecInput.Value() == "" {
					m.TestModel.Err = fmt.Errorf("spec file path cannot be empty")
					return m, nil
				}
				m.TestModel.Step = 1
				m.TestModel.UrlInput.Focus()
				return m, nil
			case tea.KeyCtrlC, tea.KeyEsc:
				m.Screen = models.MenuScreen
				m.TestModel = ui.InitialTestModel()
				return m, nil
			}
			m.TestModel.SpecInput, cmd = m.TestModel.SpecInput.Update(msg)
		case testing.TestCompleteMsg:
			m.TestModel.Results = msg.Results
			m.TestModel.Err = nil
			m.TestModel.Step = 3
			m.TestModel.Testing = false
			return m, nil
		case testing.TestErrorMsg:
			m.TestModel.Err = msg.Err
			m.TestModel.Step = 3
			m.TestModel.Testing = false
			return m, nil
		}
	case 1:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				if m.TestModel.UrlInput.Value() == "" {
					m.TestModel.Err = fmt.Errorf("base URL cannot be empty")
					return m, nil
				}

				m.Config.SpecPath = m.TestModel.SpecInput.Value()
				m.Config.BaseURL = m.TestModel.UrlInput.Value()
				config.SaveConfig(m.Config)

				// Check if we should show endpoint selector
				if m.TestModel.SelectEndpoints {
					// Load endpoints from spec
					endpoints, err := validation.ExtractEndpoints(m.Config.SpecPath)
					if err != nil {
						m.TestModel.Err = fmt.Errorf("failed to load endpoints: %w", err)
						return m, nil
					}

					// Initialize endpoint selector
					m.EndpointSelectorModel = ui.InitialEndpointSelectorModel()
					m.EndpointSelectorModel.AllEndpoints = endpoints
					m.EndpointSelectorModel.FilteredEndpoints = endpoints
					m.EndpointSelectorModel.Ready = true

					// Switch to endpoint selector screen
					m.Screen = models.EndpointSelectorScreen
					return m, nil
				}

				// Normal flow: test all endpoints
				m.TestModel.Step = 2
				m.TestModel.Testing = true
				m.TestModel.TestStartTime = time.Now()
				return m, testing.RunTestParallelCmd(m.TestModel.SpecInput.Value(), m.TestModel.UrlInput.Value(), nil, m.VerboseMode, m.Config.MaxConcurrency)
			case tea.KeyCtrlC, tea.KeyEsc:
				m.Screen = models.MenuScreen
				m.TestModel = ui.InitialTestModel()
				return m, nil
			}
			m.TestModel.UrlInput, cmd = m.TestModel.UrlInput.Update(msg)
		}
	case 2:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				m.Screen = models.MenuScreen
				m.TestModel = ui.InitialTestModel()
				return m, nil
			}
		case testing.TestCompleteMsg:
			m.TestModel.Results = msg.Results
			m.TestModel.Err = nil
			m.TestModel.Step = 3
			m.TestModel.Testing = false
			
			// Save to history
			duration := time.Since(m.TestModel.TestStartTime)
			entry := models.CreateHistoryEntry(
				m.TestModel.SpecInput.Value(),
				m.TestModel.UrlInput.Value(),
				msg.Results,
				duration,
			)
			m.History.AddEntry(entry)
			
			// Persist history to disk (ignore errors to not disrupt user flow)
			_ = models.SaveHistory(m.History)
			
			return m, nil
		case testing.TestErrorMsg:
			m.TestModel.Err = msg.Err
			m.TestModel.Step = 3
			m.TestModel.Testing = false
			return m, nil
		}
		m.TestModel.Spinner, cmd = m.TestModel.Spinner.Update(msg)
	case 3:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			// If filter is active, handle filter input first
			if m.TestModel.FilterActive {
				switch msg.Type {
				case tea.KeyEsc:
					// Esc while filtering: exit filter mode
					m.TestModel.FilterActive = false
					m.TestModel.FilterInput.Blur()
					m.TestModel.FilterInput.SetValue("")
					return m, nil
				case tea.KeyEnter:
					// Enter while filtering: return to menu
					m.Screen = models.MenuScreen
					m.TestModel = ui.InitialTestModel()
					return m, nil
				default:
					// Route all other keys to filter input
					m.TestModel.FilterInput, cmd = m.TestModel.FilterInput.Update(msg)
					return m, cmd
				}
			}
			
			// Normal key handling when filter is not active
			switch msg.String() {
			case "f":
				// Toggle filter mode
				m.TestModel.FilterActive = !m.TestModel.FilterActive
				if m.TestModel.FilterActive {
					m.TestModel.FilterInput.Focus()
				} else {
					m.TestModel.FilterInput.Blur()
					m.TestModel.FilterInput.SetValue("")
				}
				return m, nil
			case "e":
				if len(m.TestModel.Results) > 0 {
					specPath := m.TestModel.SpecInput.Value()
					filename, err := export.ExportResults(m.TestModel.Results, specPath)
					if err != nil {
						m.TestModel.Err = errors.EnhanceFileError(err, "export file")
					} else {
						m.TestModel.ExportSuccess = fmt.Sprintf("✅ Exported JSON to %s", filename)
					}
				}
				return m, nil
			case "h":
				if len(m.TestModel.Results) > 0 {
					specPath := m.TestModel.SpecInput.Value()
					baseURL := m.TestModel.UrlInput.Value()
					filename, err := export.ExportResultsToHTML(m.TestModel.Results, specPath, baseURL)
					if err != nil {
						m.TestModel.Err = errors.EnhanceFileError(err, "HTML export file")
					} else {
						m.TestModel.ExportSuccess = fmt.Sprintf("✅ Exported HTML to %s", filename)
					}
				}
				return m, nil
			case "j":
				if len(m.TestModel.Results) > 0 {
					specPath := m.TestModel.SpecInput.Value()
					baseURL := m.TestModel.UrlInput.Value()
					filename, err := export.ExportResultsToJUnit(m.TestModel.Results, specPath, baseURL)
					if err != nil {
						m.TestModel.Err = errors.EnhanceFileError(err, "JUnit XML export file")
					} else {
						m.TestModel.ExportSuccess = fmt.Sprintf("✅ Exported JUnit XML to %s", filename)
					}
				}
				return m, nil
			case "r":
				// View test run history
				m.Screen = models.HistoryScreen
				m.HistoryIndex = 0
				return m, nil
			case "l":
				if m.VerboseMode && len(m.TestModel.Results) > 0 {
					selectedIdx := m.TestModel.Table.Cursor()
					if selectedIdx >= 0 && selectedIdx < len(m.TestModel.Results) {
						result := m.TestModel.Results[selectedIdx]
						if result.LogEntry != nil {
							m.TestModel.ShowingLog = true
							m.TestModel.SelectedLog = selectedIdx
							m.TestModel.Step = 4
							return m, nil
						}
					}
				}
				return m, nil
			}
			switch msg.Type {
			case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
				m.Screen = models.MenuScreen
				m.TestModel = ui.InitialTestModel()
				return m, nil
			}
		}
		m.TestModel.Table, cmd = m.TestModel.Table.Update(msg)
	case 4:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEsc, tea.KeyEnter:
				m.TestModel.ShowingLog = false
				m.TestModel.Step = 3
				return m, nil
			case tea.KeyCtrlC:
				m.Screen = models.MenuScreen
				m.TestModel = ui.InitialTestModel()
				return m, nil
			}
		}
	}

	return m, cmd
}

// updateHistory handles key events in the history screen
func (m model) updateHistory(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Return to results screen
		m.Screen = models.TestScreen
		return m, nil
	case "up", "k":
		if m.HistoryIndex > 0 {
			m.HistoryIndex--
		}
		return m, nil
	case "down", "j":
		if m.HistoryIndex < len(m.History.Entries)-1 {
			m.HistoryIndex++
		}
		return m, nil
	case "enter":
		// Replay selected test
		if m.HistoryIndex >= 0 && m.HistoryIndex < len(m.History.Entries) {
			entry := m.History.Entries[m.HistoryIndex]
			
			// Set spec and URL from history
			m.TestModel.SpecInput.SetValue(entry.SpecPath)
			m.TestModel.UrlInput.SetValue(entry.BaseURL)
			
			// Save to config
			m.Config.SpecPath = entry.SpecPath
			m.Config.BaseURL = entry.BaseURL
			config.SaveConfig(m.Config)
			
			// Start testing
			m.Screen = models.TestScreen
			m.TestModel.Step = 2
			m.TestModel.Testing = true
			m.TestModel.Results = nil
			m.TestModel.Err = nil
			m.TestModel.ExportSuccess = ""
			m.TestModel.TestStartTime = time.Now()
			
			return m, testing.RunTestCmd(entry.SpecPath, entry.BaseURL, nil, m.VerboseMode)
		}
		return m, nil
	case "ctrl+c", "q":
		return m, tea.Quit
	}
	return m, nil
}

// updateCustomRequest handles events in the custom request screen
func (m model) updateCustomRequest(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.CustomRequestModel.Step {
	case 0: // Method input
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				method := strings.ToUpper(strings.TrimSpace(m.CustomRequestModel.MethodInput.Value()))
				if method == "" {
					m.CustomRequestModel.Err = fmt.Errorf("HTTP method cannot be empty")
					return m, nil
				}
				validMethods := map[string]bool{"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true, "HEAD": true, "OPTIONS": true}
				if !validMethods[method] {
					m.CustomRequestModel.Err = fmt.Errorf("invalid HTTP method: %s (use GET, POST, PUT, PATCH, DELETE, HEAD, or OPTIONS)", method)
					return m, nil
				}
				m.CustomRequestModel.Request.Method = method
				m.CustomRequestModel.Step = 1
				m.CustomRequestModel.MethodInput.Blur()
				m.CustomRequestModel.EndpointInput.Focus()
				m.CustomRequestModel.Err = nil
				return m, nil
			case tea.KeyCtrlC, tea.KeyEsc:
				m.Screen = models.MenuScreen
				m.CustomRequestModel = ui.InitialCustomRequestModel()
				return m, nil
			}
			m.CustomRequestModel.MethodInput, cmd = m.CustomRequestModel.MethodInput.Update(msg)
		}
	
	case 1: // Endpoint input
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				endpoint := strings.TrimSpace(m.CustomRequestModel.EndpointInput.Value())
				if endpoint == "" {
					m.CustomRequestModel.Err = fmt.Errorf("endpoint URL cannot be empty")
					return m, nil
				}
				m.CustomRequestModel.Request.Endpoint = endpoint
				m.CustomRequestModel.Step = 2
				m.CustomRequestModel.EndpointInput.Blur()
				m.CustomRequestModel.HeaderKeyInput.Focus()
				m.CustomRequestModel.Err = nil
				return m, nil
			case tea.KeyCtrlC, tea.KeyEsc:
				m.Screen = models.MenuScreen
				m.CustomRequestModel = ui.InitialCustomRequestModel()
				return m, nil
			}
			m.CustomRequestModel.EndpointInput, cmd = m.CustomRequestModel.EndpointInput.Update(msg)
		}
	
	case 2: // Headers input
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				headerKey := strings.TrimSpace(m.CustomRequestModel.HeaderKeyInput.Value())
				if headerKey == "" {
					// Skip headers, go to body
					m.CustomRequestModel.Step = 3
					m.CustomRequestModel.HeaderKeyInput.Blur()
					m.CustomRequestModel.BodyInput.Focus()
					m.CustomRequestModel.Err = nil
					return m, nil
				}
				// Need header value
				if m.CustomRequestModel.HeaderValueInput.Value() == "" {
					m.CustomRequestModel.HeaderKeyInput.Blur()
					m.CustomRequestModel.HeaderValueInput.Focus()
					return m, nil
				}
				// Save header
				headerValue := strings.TrimSpace(m.CustomRequestModel.HeaderValueInput.Value())
				m.CustomRequestModel.Request.Headers[headerKey] = headerValue
				// Reset for next header
				m.CustomRequestModel.HeaderKeyInput.SetValue("")
				m.CustomRequestModel.HeaderValueInput.SetValue("")
				m.CustomRequestModel.HeaderValueInput.Blur()
				m.CustomRequestModel.HeaderKeyInput.Focus()
				m.CustomRequestModel.Err = nil
				return m, nil
			case tea.KeyCtrlC, tea.KeyEsc:
				m.Screen = models.MenuScreen
				m.CustomRequestModel = ui.InitialCustomRequestModel()
				return m, nil
			}
			if m.CustomRequestModel.HeaderKeyInput.Focused() {
				m.CustomRequestModel.HeaderKeyInput, cmd = m.CustomRequestModel.HeaderKeyInput.Update(msg)
			} else {
				m.CustomRequestModel.HeaderValueInput, cmd = m.CustomRequestModel.HeaderValueInput.Update(msg)
			}
		}
	
	case 3: // Body input
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				body := strings.TrimSpace(m.CustomRequestModel.BodyInput.Value())
				if body != "" {
					// Validate JSON
					if err := testing.ValidateJSONBody(body); err != nil {
						m.CustomRequestModel.Err = fmt.Errorf("invalid JSON: %v", err)
						return m, nil
					}
				}
				m.CustomRequestModel.Request.Body = body
				m.CustomRequestModel.Step = 4
				m.CustomRequestModel.BodyInput.Blur()
				m.CustomRequestModel.Testing = true
				// Execute the request
				return m, testing.ExecuteCustomRequestCmd(
					m.CustomRequestModel.Request.Method,
					m.CustomRequestModel.Request.Endpoint,
					m.CustomRequestModel.Request.Headers,
					m.CustomRequestModel.Request.Body,
					nil, // TODO: Add auth support
					m.VerboseMode,
				)
			case tea.KeyCtrlC, tea.KeyEsc:
				m.Screen = models.MenuScreen
				m.CustomRequestModel = ui.InitialCustomRequestModel()
				return m, nil
			}
			m.CustomRequestModel.BodyInput, cmd = m.CustomRequestModel.BodyInput.Update(msg)
		}
	
	case 4: // Executing request (showing spinner)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				m.Screen = models.MenuScreen
				m.CustomRequestModel = ui.InitialCustomRequestModel()
				return m, nil
			}
		case testing.TestCompleteMsg:
			if len(msg.Results) > 0 {
				result := msg.Results[0]
				m.CustomRequestModel.Result = &result
				m.CustomRequestModel.Err = nil
				m.CustomRequestModel.Step = 5
				m.CustomRequestModel.Testing = false
				
				// Save to history
				entry := models.CreateHistoryEntry(
					"Custom Request",
					m.CustomRequestModel.Request.Endpoint,
					msg.Results,
					result.Duration,
				)
				m.History.AddEntry(entry)
				_ = models.SaveHistory(m.History)
			}
			return m, nil
		case testing.TestErrorMsg:
			m.CustomRequestModel.Err = msg.Err
			m.CustomRequestModel.Step = 5
			m.CustomRequestModel.Testing = false
			return m, nil
		}
		m.CustomRequestModel.Spinner, cmd = m.CustomRequestModel.Spinner.Update(msg)
	
	case 5: // Show results
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
				m.Screen = models.MenuScreen
				m.CustomRequestModel = ui.InitialCustomRequestModel()
				return m, nil
			}
		}
	}
	
	return m, cmd
}

// updateEndpointSelector handles events in the endpoint selector screen
func (m model) updateEndpointSelector(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			// Cancel and return to menu
			m.Screen = models.MenuScreen
			m.EndpointSelectorModel = ui.InitialEndpointSelectorModel()
			return m, nil

		case tea.KeyEnter:
			// Confirm selection and start testing
			if !m.EndpointSelectorModel.Ready {
				return m, nil
			}

			selected := validation.GetSelectedEndpoints(m.EndpointSelectorModel.AllEndpoints)
			if len(selected) == 0 {
				m.EndpointSelectorModel.Err = fmt.Errorf("no endpoints selected")
				return m, nil
			}

			// Update config with spec path
			m.Config.SpecPath = m.TestModel.SpecInput.Value()
			m.Config.BaseURL = m.TestModel.UrlInput.Value()
			config.SaveConfig(m.Config)

			// Move to test screen with spinner
			m.Screen = models.TestScreen
			m.TestModel.Step = 2  // Spinner step
			m.TestModel.Testing = true
			m.TestModel.Err = nil
			m.TestModel.TestStartTime = time.Now()

			// Start parallel test execution with selected endpoints
			return m, testing.RunTestParallelCmdWithSelection(
				m.Config.SpecPath,
				m.Config.BaseURL,
				m.Config.Auth,
				m.VerboseMode,
				m.Config.MaxConcurrency,
				selected,
			)

		case tea.KeyUp, tea.KeyCtrlP:
			// Move cursor up
			if m.EndpointSelectorModel.Cursor > 0 {
				m.EndpointSelectorModel.Cursor--
				// Scroll up if needed
				if m.EndpointSelectorModel.Cursor < m.EndpointSelectorModel.Offset {
					m.EndpointSelectorModel.Offset = m.EndpointSelectorModel.Cursor
				}
			}
			return m, nil

		case tea.KeyDown, tea.KeyCtrlN:
			// Move cursor down
			endpoints := m.EndpointSelectorModel.FilteredEndpoints
			if len(endpoints) == 0 {
				endpoints = m.EndpointSelectorModel.AllEndpoints
			}
			if m.EndpointSelectorModel.Cursor < len(endpoints)-1 {
				m.EndpointSelectorModel.Cursor++
				// Scroll down if needed
				visibleHeight := 15
				if m.EndpointSelectorModel.Cursor >= m.EndpointSelectorModel.Offset+visibleHeight {
					m.EndpointSelectorModel.Offset = m.EndpointSelectorModel.Cursor - visibleHeight + 1
				}
			}
			return m, nil

		case tea.KeyRunes:
			switch string(msg.Runes) {
			case " ":
				// Toggle selection for current endpoint
				if m.EndpointSelectorModel.Ready {
					endpoints := &m.EndpointSelectorModel.AllEndpoints
					if m.EndpointSelectorModel.Cursor < len(*endpoints) {
						(*endpoints)[m.EndpointSelectorModel.Cursor].Selected = !(*endpoints)[m.EndpointSelectorModel.Cursor].Selected
					}
				}
				return m, nil

			case "a", "A":
				// Select all
				m.EndpointSelectorModel.AllEndpoints = validation.SelectAllEndpoints(m.EndpointSelectorModel.AllEndpoints)
				return m, nil

			case "d", "D":
				// Deselect all
				m.EndpointSelectorModel.AllEndpoints = validation.DeselectAllEndpoints(m.EndpointSelectorModel.AllEndpoints)
				return m, nil
			}
		}

		// Update search input
		m.EndpointSelectorModel.SearchInput, cmd = m.EndpointSelectorModel.SearchInput.Update(msg)
		
		// Filter endpoints based on search
		query := m.EndpointSelectorModel.SearchInput.Value()
		m.EndpointSelectorModel.FilteredEndpoints = validation.FilterEndpoints(m.EndpointSelectorModel.AllEndpoints, query)
		
		// Reset cursor if out of bounds
		if m.EndpointSelectorModel.Cursor >= len(m.EndpointSelectorModel.FilteredEndpoints) {
			m.EndpointSelectorModel.Cursor = 0
			m.EndpointSelectorModel.Offset = 0
		}
	}

	return m, cmd
}

// View renders the current screen based on the application state
func (m model) View() string {
	switch m.Screen {
	case models.MenuScreen:
		return ui.ViewMenu(m.Model)
	case models.HelpScreen:
		return ui.ViewHelp(m.Model)
	case models.ValidateScreen:
		return ui.ViewValidate(m.Model)
	case models.TestScreen:
		return ui.ViewTest(m.Model)
	case models.CustomRequestScreen:
		return ui.ViewCustomRequest(m.Model)
	case models.HistoryScreen:
		return ui.ViewHistory(m.Model)
	case models.EndpointSelectorScreen:
		return ui.ViewEndpointSelector(m.Model)
	default:
		return "Unknown screen"
	}
}

// main initializes and runs the Bubble Tea TUI program
func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
}

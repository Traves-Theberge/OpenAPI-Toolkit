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

	// Initialize the main application model with default values and loaded config
	m := model{
		Model: models.Model{
			Cursor:        0,
			Screen:        models.MenuScreen,
			Width:         80,
			Height:        24,
			VerboseMode:   cfg.VerboseMode,
			Config:        cfg,
			ValidateModel: ui.InitialValidateModel(),
			TestModel:     ui.InitialTestModel(),
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
		}
	case testing.TestCompleteMsg:
		if m.Screen == models.TestScreen {
			return m.updateTest(msg)
		}
	case testing.TestErrorMsg:
		if m.Screen == models.TestScreen {
			return m.updateTest(msg)
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
		if m.Cursor < 3 {
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
			m.Screen = models.TestScreen
			return m, nil
		case 2:
			m.Screen = models.HelpScreen
			return m, nil
		case 3:
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

				m.TestModel.Step = 2
				m.TestModel.Testing = true
				return m, testing.RunTestCmd(m.TestModel.SpecInput.Value(), m.TestModel.UrlInput.Value(), nil, m.VerboseMode)
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

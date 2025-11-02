// Package main implements a professional Terminal User Interface (TUI) for OpenAPI specification
// validation and API testing. It uses the Charm Bracelet Bubble Tea framework for reactive
// terminal applications and Lip Gloss for beautiful styling.
//
// Architecture:
// - Modular design with separate files for models, views, testing, validation, etc.
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
)

// initialModel creates and initializes the main application model
// Loads configuration and prepares all sub-models
func initialModel() model {
// Load configuration from file
cfg := loadConfig()

// Initialize the main application model with default values and loaded config
m := model{
cursor:         0,
screen:         menuScreen,
width:          80,
height:         24,
verboseMode:    cfg.VerboseMode,
config:         cfg,
validateModel:  initialValidateModel(),
testModel:      initialTestModel(),
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
func (m model) Init() tea.Cmd {
return m.testModel.spinner.Tick
}

// Update handles all incoming messages and updates the model accordingly
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
switch msg := msg.(type) {
case tea.WindowSizeMsg:
m.width = msg.Width
m.height = msg.Height
return m, nil
case tea.KeyMsg:
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
case testCompleteMsg:
if m.screen == testScreen {
return m.updateTest(msg)
}
case testErrorMsg:
if m.screen == testScreen {
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
if m.cursor > 0 {
m.cursor--
}
case "down", "j":
if m.cursor < 3 {
m.cursor++
}
case "h", "?":
m.screen = helpScreen
return m, nil
case "v":
m.verboseMode = !m.verboseMode
m.config.VerboseMode = m.verboseMode
saveConfig(m.config)
return m, nil
case "enter":
switch m.cursor {
case 0:
m.screen = validateScreen
return m, nil
case 1:
m.screen = testScreen
return m, nil
case 2:
m.screen = helpScreen
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
m.screen = menuScreen
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
if m.validateModel.done {
m.screen = menuScreen
m.validateModel = initialValidateModel()
return m, nil
}

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
m.screen = menuScreen
m.validateModel = initialValidateModel()
return m, nil
}
}

m.validateModel.textInput, cmd = m.validateModel.textInput.Update(msg)
return m, cmd
}

// updateTest handles events in the testing screen
func (m model) updateTest(msg tea.Msg) (tea.Model, tea.Cmd) {
var cmd tea.Cmd

switch m.testModel.step {
case 0:
switch msg := msg.(type) {
case tea.KeyMsg:
switch msg.Type {
case tea.KeyEnter:
if m.testModel.specInput.Value() == "" {
m.testModel.err = fmt.Errorf("spec file path cannot be empty")
return m, nil
}
m.testModel.step = 1
m.testModel.urlInput.Focus()
return m, nil
case tea.KeyCtrlC, tea.KeyEsc:
m.screen = menuScreen
m.testModel = initialTestModel()
return m, nil
}
m.testModel.specInput, cmd = m.testModel.specInput.Update(msg)
case testCompleteMsg:
m.testModel.results = msg.results
m.testModel.err = nil
m.testModel.step = 3
m.testModel.testing = false
return m, nil
case testErrorMsg:
m.testModel.err = msg.err
m.testModel.step = 3
m.testModel.testing = false
return m, nil
}
case 1:
switch msg := msg.(type) {
case tea.KeyMsg:
switch msg.Type {
case tea.KeyEnter:
if m.testModel.urlInput.Value() == "" {
m.testModel.err = fmt.Errorf("base URL cannot be empty")
return m, nil
}

m.config.SpecPath = m.testModel.specInput.Value()
m.config.BaseURL = m.testModel.urlInput.Value()
saveConfig(m.config)

m.testModel.step = 2
m.testModel.testing = true
return m, runTestCmd(m.testModel.specInput.Value(), m.testModel.urlInput.Value(), nil, m.verboseMode)
case tea.KeyCtrlC, tea.KeyEsc:
m.screen = menuScreen
m.testModel = initialTestModel()
return m, nil
}
m.testModel.urlInput, cmd = m.testModel.urlInput.Update(msg)
}
case 2:
switch msg := msg.(type) {
case tea.KeyMsg:
switch msg.Type {
case tea.KeyCtrlC, tea.KeyEsc:
m.screen = menuScreen
m.testModel = initialTestModel()
return m, nil
}
case testCompleteMsg:
m.testModel.results = msg.results
m.testModel.err = nil
m.testModel.step = 3
m.testModel.testing = false
return m, nil
case testErrorMsg:
m.testModel.err = msg.err
m.testModel.step = 3
m.testModel.testing = false
return m, nil
}
m.testModel.spinner, cmd = m.testModel.spinner.Update(msg)
case 3:
switch msg := msg.(type) {
case tea.KeyMsg:
switch msg.String() {
case "e":
if len(m.testModel.results) > 0 {
specPath := m.testModel.specInput.Value()
filename, err := exportResults(m.testModel.results, specPath)
if err != nil {
m.testModel.err = enhanceFileError(err, "export file")
} else {
m.testModel.exportSuccess = fmt.Sprintf("✅ Exported to %s", filename)
}
}
return m, nil
case "l":
if m.verboseMode && len(m.testModel.results) > 0 {
selectedIdx := m.testModel.table.Cursor()
if selectedIdx >= 0 && selectedIdx < len(m.testModel.results) {
result := m.testModel.results[selectedIdx]
if result.logEntry != nil {
m.testModel.showingLog = true
m.testModel.selectedLog = selectedIdx
m.testModel.step = 4
return m, nil
}
}
}
return m, nil
}
switch msg.Type {
case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
m.screen = menuScreen
m.testModel = initialTestModel()
return m, nil
}
}
m.testModel.table, cmd = m.testModel.table.Update(msg)
case 4:
switch msg := msg.(type) {
case tea.KeyMsg:
switch msg.Type {
case tea.KeyEsc, tea.KeyEnter:
m.testModel.showingLog = false
m.testModel.step = 3
return m, nil
case tea.KeyCtrlC:
m.screen = menuScreen
m.testModel = initialTestModel()
return m, nil
}
}
}

return m, cmd
}

// View renders the current screen based on the application state
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

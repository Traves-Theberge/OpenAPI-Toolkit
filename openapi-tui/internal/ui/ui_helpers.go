package ui

import (
"fmt"
"github.com/charmbracelet/bubbles/spinner"
"github.com/charmbracelet/bubbles/table"
"github.com/charmbracelet/bubbles/textinput"
"github.com/charmbracelet/lipgloss"
"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

// InitialValidateModel creates and configures the validation model
func InitialValidateModel() models.ValidateModel {
ti := textinput.New()
ti.Placeholder = "Path to OpenAPI spec file (e.g., openapi.yaml)"
ti.Focus()
ti.CharLimit = 156
ti.Width = 60

return models.ValidateModel{
TextInput: ti,
Err:       nil,
}
}

// InitialTestModel creates and configures the test model
func InitialTestModel() models.TestModel {
specTi := textinput.New()
specTi.Placeholder = "Path to OpenAPI spec file (e.g., openapi.yaml)"
specTi.Focus()
specTi.CharLimit = 156
specTi.Width = 60

urlTi := textinput.New()
urlTi.Placeholder = "Base URL (e.g., https://api.example.com)"
urlTi.CharLimit = 156
urlTi.Width = 60

s := spinner.New()
s.Spinner = spinner.Dot
s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ECDC4"))

columns := []table.Column{
{Title: "Method", Width: 8},
{Title: "Endpoint", Width: 40},
{Title: "Status", Width: 10},
{Title: "Message", Width: 30},
}

t := table.New(
table.WithColumns(columns),
table.WithFocused(false),
table.WithHeight(10),
)

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

	// Custom request inputs
	methodTi := textinput.New()
	methodTi.Placeholder = "HTTP Method (GET, POST, PUT, PATCH, DELETE)"
	methodTi.CharLimit = 10
	methodTi.Width = 60

	endpointTi := textinput.New()
	endpointTi.Placeholder = "Full URL (e.g., https://api.example.com/users)"
	endpointTi.CharLimit = 256
	endpointTi.Width = 60

	headerKeyTi := textinput.New()
	headerKeyTi.Placeholder = "Header name (e.g., Content-Type)"
	headerKeyTi.CharLimit = 100
	headerKeyTi.Width = 60

	headerValueTi := textinput.New()
	headerValueTi.Placeholder = "Header value (e.g., application/json)"
	headerValueTi.CharLimit = 256
	headerValueTi.Width = 60

	bodyTi := textinput.New()
	bodyTi.Placeholder = "Request body (JSON)"
	bodyTi.CharLimit = 0 // No limit for body
	bodyTi.Width = 60

	filterTi := textinput.New()
	filterTi.Placeholder = "Filter results (status, method, endpoint, message)"
	filterTi.CharLimit = 100
	filterTi.Width = 60

	return models.TestModel{
		SpecInput:   specTi,
		UrlInput:    urlTi,
		Spinner:     s,
		Table:       t,
		FilterInput: filterTi,
	}
}

// InitialCustomRequestModel creates and configures the custom request model
func InitialCustomRequestModel() models.CustomRequestModel {
	methodTi := textinput.New()
	methodTi.Placeholder = "HTTP Method (GET, POST, PUT, PATCH, DELETE)"
	methodTi.CharLimit = 10
	methodTi.Width = 60
	methodTi.Focus()

	endpointTi := textinput.New()
	endpointTi.Placeholder = "Full URL (e.g., https://api.example.com/users)"
	endpointTi.CharLimit = 256
	endpointTi.Width = 60

	headerKeyTi := textinput.New()
	headerKeyTi.Placeholder = "Header name (e.g., Content-Type) or leave empty to skip"
	headerKeyTi.CharLimit = 100
	headerKeyTi.Width = 60

	headerValueTi := textinput.New()
	headerValueTi.Placeholder = "Header value (e.g., application/json)"
	headerValueTi.CharLimit = 256
	headerValueTi.Width = 60

	bodyTi := textinput.New()
	bodyTi.Placeholder = `Request body (JSON, e.g., {"key": "value"}) or leave empty`
	bodyTi.CharLimit = 0 // No limit for body
	bodyTi.Width = 60

	filterTi := textinput.New()
	filterTi.Placeholder = "Filter results"
	filterTi.CharLimit = 100
	filterTi.Width = 60

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#4ECDC4"))

	columns := []table.Column{
		{Title: "Method", Width: 8},
		{Title: "Endpoint", Width: 50},
		{Title: "Status", Width: 10},
		{Title: "Duration", Width: 15},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(false),
		table.WithHeight(5),
	)

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

	return models.CustomRequestModel{
		MethodInput:      methodTi,
		EndpointInput:    endpointTi,
		HeaderKeyInput:   headerKeyTi,
		HeaderValueInput: headerValueTi,
		BodyInput:        bodyTi,
		FilterInput:      filterTi,
		Spinner:          s,
		Table:            t,
		Request:          models.CustomRequest{Headers: make(map[string]string), QueryParams: make(map[string]string), IsCustom: true},
	}
}

// InitialEndpointSelectorModel creates and configures the endpoint selector model
func InitialEndpointSelectorModel() models.EndpointSelectorModel {
	searchTi := textinput.New()
	searchTi.Placeholder = "Search endpoints (path, method, tag) or use filters (method:GET, tag:users)"
	searchTi.CharLimit = 100
	searchTi.Width = 80
	searchTi.Focus()

	return models.EndpointSelectorModel{
		SearchInput:       searchTi,
		AllEndpoints:      []models.EndpointInfo{},
		FilteredEndpoints: []models.EndpointInfo{},
		Cursor:            0,
		Offset:            0,
		Ready:             false,
	}
}

// InitialConfigEditorModel creates and configures the configuration editor model
func InitialConfigEditorModel(cfg models.Config) models.ConfigEditorModel {
	// Spec Path input
	specPathTi := textinput.New()
	specPathTi.Placeholder = "Path to OpenAPI spec (e.g., openapi.yaml)"
	specPathTi.CharLimit = 200
	specPathTi.Width = 60
	if cfg.SpecPath != "" {
		specPathTi.SetValue(cfg.SpecPath)
	}

	// Base URL input
	baseURLTi := textinput.New()
	baseURLTi.Placeholder = "Base URL (e.g., https://api.example.com)"
	baseURLTi.CharLimit = 200
	baseURLTi.Width = 60
	if cfg.BaseURL != "" {
		baseURLTi.SetValue(cfg.BaseURL)
	}

	// Auth Type input
	authTypeTi := textinput.New()
	authTypeTi.Placeholder = "none, bearer, apikey, basic"
	authTypeTi.CharLimit = 20
	authTypeTi.Width = 30
	if cfg.Auth != nil && cfg.Auth.AuthType != "" {
		authTypeTi.SetValue(cfg.Auth.AuthType)
	}

	// Token input
	tokenTi := textinput.New()
	tokenTi.Placeholder = "Bearer token value"
	tokenTi.CharLimit = 200
	tokenTi.Width = 60
	tokenTi.EchoMode = textinput.EchoPassword
	tokenTi.EchoCharacter = '•'
	if cfg.Auth != nil && cfg.Auth.Token != "" {
		tokenTi.SetValue(cfg.Auth.Token)
	}

	// API Key Name input
	apiKeyNameTi := textinput.New()
	apiKeyNameTi.Placeholder = "API key header/query name (e.g., X-API-Key)"
	apiKeyNameTi.CharLimit = 100
	apiKeyNameTi.Width = 50
	if cfg.Auth != nil && cfg.Auth.APIKeyName != "" {
		apiKeyNameTi.SetValue(cfg.Auth.APIKeyName)
	}

	// API Key Location input
	apiKeyInTi := textinput.New()
	apiKeyInTi.Placeholder = "header or query"
	apiKeyInTi.CharLimit = 10
	apiKeyInTi.Width = 20
	if cfg.Auth != nil && cfg.Auth.APIKeyIn != "" {
		apiKeyInTi.SetValue(cfg.Auth.APIKeyIn)
	}

	// Username input
	usernameTi := textinput.New()
	usernameTi.Placeholder = "Basic auth username"
	usernameTi.CharLimit = 100
	usernameTi.Width = 40
	if cfg.Auth != nil && cfg.Auth.Username != "" {
		usernameTi.SetValue(cfg.Auth.Username)
	}

	// Password input
	passwordTi := textinput.New()
	passwordTi.Placeholder = "Basic auth password"
	passwordTi.CharLimit = 200
	passwordTi.Width = 40
	passwordTi.EchoMode = textinput.EchoPassword
	passwordTi.EchoCharacter = '•'
	if cfg.Auth != nil && cfg.Auth.Password != "" {
		passwordTi.SetValue(cfg.Auth.Password)
	}

	// Max Concurrency input
	maxConcurrTi := textinput.New()
	maxConcurrTi.Placeholder = "0 (auto-detect)"
	maxConcurrTi.CharLimit = 3
	maxConcurrTi.Width = 10
	if cfg.MaxConcurrency > 0 {
		maxConcurrTi.SetValue(string(rune(cfg.MaxConcurrency + '0')))
	}

	// Verbose Mode input
	verboseTi := textinput.New()
	verboseTi.Placeholder = "true or false"
	verboseTi.CharLimit = 5
	verboseTi.Width = 10
	if cfg.VerboseMode {
		verboseTi.SetValue("true")
	} else {
		verboseTi.SetValue("false")
	}

	// Max Retries input
	maxRetriesTi := textinput.New()
	maxRetriesTi.Placeholder = "3 (default: 3 retries)"
	maxRetriesTi.CharLimit = 2
	maxRetriesTi.Width = 10
	if cfg.MaxRetries > 0 {
		maxRetriesTi.SetValue(fmt.Sprintf("%d", cfg.MaxRetries))
	}

	// Retry Delay input
	retryDelayTi := textinput.New()
	retryDelayTi.Placeholder = "1000 (milliseconds)"
	retryDelayTi.CharLimit = 5
	retryDelayTi.Width = 15
	if cfg.RetryDelay > 0 {
		retryDelayTi.SetValue(fmt.Sprintf("%d", cfg.RetryDelay))
	}

	// Focus first field
	specPathTi.Focus()

	return models.ConfigEditorModel{
		FocusedField:    0,
		SpecPathInput:   specPathTi,
		BaseURLInput:    baseURLTi,
		AuthTypeInput:   authTypeTi,
		TokenInput:      tokenTi,
		APIKeyNameInput: apiKeyNameTi,
		APIKeyInInput:   apiKeyInTi,
		UsernameInput:   usernameTi,
		PasswordInput:   passwordTi,
		MaxConcurrInput: maxConcurrTi,
		VerboseInput:    verboseTi,
		MaxRetriesInput: maxRetriesTi,
		RetryDelayInput: retryDelayTi,
		OriginalConfig:  cfg,
	}
}
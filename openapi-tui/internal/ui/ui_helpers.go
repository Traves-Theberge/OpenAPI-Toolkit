package ui

import (
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

return models.TestModel{
SpecInput: specTi,
UrlInput:  urlTi,
Spinner:   s,
Table:     t,
}
}

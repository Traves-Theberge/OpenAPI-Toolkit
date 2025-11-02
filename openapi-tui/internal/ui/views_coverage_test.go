package ui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui/internal/models"
)

func TestViewLogDetail(t *testing.T) {
	// Test with a log entry
	log := &models.LogEntry{
		RequestURL: "http://example.com/users",
		RequestHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		RequestBody: `{"test": "data"}`,
		ResponseHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		ResponseBody: `{"result": "success"}`,
		Duration:     100 * time.Millisecond,
		Timestamp:    time.Now(),
	}

	result := models.TestResult{
		Method:   "GET",
		Endpoint: "/users",
		Status:   "200",
		Message:  "OK",
		Duration: 100 * time.Millisecond,
		LogEntry: log,
	}

	m := models.Model{
		Width:  100,
		Height: 50,
	}

	output := ViewLogDetail(m, result, log)
	if output == "" {
		t.Fatal("ViewLogDetail returned empty string")
	}
	if !strings.Contains(output, "GET") {
		t.Error("ViewLogDetail should contain method GET")
	}
	if !strings.Contains(output, "/users") {
		t.Error("ViewLogDetail should contain endpoint /users")
	}
	if !strings.Contains(output, "http://example.com/users") {
		t.Error("ViewLogDetail should contain request URL")
	}
}

func TestViewCustomRequest(t *testing.T) {
	methodTi := textinput.New()
	methodTi.SetValue("GET")
	
	endpointTi := textinput.New()
	endpointTi.SetValue("http://example.com/api")
	
	headerKeyTi := textinput.New()
	headerValueTi := textinput.New()
	bodyTi := textinput.New()
	filterTi := textinput.New()

	columns := []table.Column{
		{Title: "Method", Width: 8},
		{Title: "Endpoint", Width: 50},
		{Title: "Status", Width: 10},
	}
	tbl := table.New(table.WithColumns(columns))

	tests := []struct {
		name  string
		model models.Model
	}{
		{
			name: "step 0 - method input",
			model: models.Model{
				CustomRequestModel: models.CustomRequestModel{
					Step:             0,
					MethodInput:      methodTi,
					EndpointInput:    endpointTi,
					HeaderKeyInput:   headerKeyTi,
					HeaderValueInput: headerValueTi,
					BodyInput:        bodyTi,
					FilterInput:      filterTi,
					Table:            tbl,
					Request:          models.CustomRequest{Headers: make(map[string]string)},
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "step 1 - endpoint input",
			model: models.Model{
				CustomRequestModel: models.CustomRequestModel{
					Step:             1,
					MethodInput:      methodTi,
					EndpointInput:    endpointTi,
					HeaderKeyInput:   headerKeyTi,
					HeaderValueInput: headerValueTi,
					BodyInput:        bodyTi,
					FilterInput:      filterTi,
					Table:            tbl,
					Request:          models.CustomRequest{Headers: make(map[string]string)},
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "step 2 - header input",
			model: models.Model{
				CustomRequestModel: models.CustomRequestModel{
					Step:             2,
					MethodInput:      methodTi,
					EndpointInput:    endpointTi,
					HeaderKeyInput:   headerKeyTi,
					HeaderValueInput: headerValueTi,
					BodyInput:        bodyTi,
					FilterInput:      filterTi,
					Table:            tbl,
					Request:          models.CustomRequest{Headers: make(map[string]string)},
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "step 3 - body input",
			model: models.Model{
				CustomRequestModel: models.CustomRequestModel{
					Step:             3,
					MethodInput:      methodTi,
					EndpointInput:    endpointTi,
					HeaderKeyInput:   headerKeyTi,
					HeaderValueInput: headerValueTi,
					BodyInput:        bodyTi,
					FilterInput:      filterTi,
					Table:            tbl,
					Request:          models.CustomRequest{Headers: make(map[string]string)},
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "testing state",
			model: models.Model{
				CustomRequestModel: models.CustomRequestModel{
					Step:             0,
					MethodInput:      methodTi,
					EndpointInput:    endpointTi,
					HeaderKeyInput:   headerKeyTi,
					HeaderValueInput: headerValueTi,
					BodyInput:        bodyTi,
					FilterInput:      filterTi,
					Table:            tbl,
					Testing:          true,
					Request:          models.CustomRequest{Headers: make(map[string]string)},
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "with result",
			model: models.Model{
				CustomRequestModel: models.CustomRequestModel{
					Step:             0,
					MethodInput:      methodTi,
					EndpointInput:    endpointTi,
					HeaderKeyInput:   headerKeyTi,
					HeaderValueInput: headerValueTi,
					BodyInput:        bodyTi,
					FilterInput:      filterTi,
					Table:            tbl,
					Result: &models.TestResult{
						Method:   "GET",
						Endpoint: "http://example.com",
						Status:   "200",
						Message:  "OK",
					},
					Request: models.CustomRequest{Headers: make(map[string]string)},
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "with error",
			model: models.Model{
				CustomRequestModel: models.CustomRequestModel{
					Step:             0,
					MethodInput:      methodTi,
					EndpointInput:    endpointTi,
					HeaderKeyInput:   headerKeyTi,
					HeaderValueInput: headerValueTi,
					BodyInput:        bodyTi,
					FilterInput:      filterTi,
					Table:            tbl,
					Err:              fmt.Errorf("test error"),
					Request:          models.CustomRequest{Headers: make(map[string]string)},
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "showing log",
			model: models.Model{
				CustomRequestModel: models.CustomRequestModel{
					Step:             0,
					MethodInput:      methodTi,
					EndpointInput:    endpointTi,
					HeaderKeyInput:   headerKeyTi,
					HeaderValueInput: headerValueTi,
					BodyInput:        bodyTi,
					FilterInput:      filterTi,
					Table:            tbl,
					ShowingLog:       true,
					Result: &models.TestResult{
						LogEntry: &models.LogEntry{
							RequestURL:  "http://example.com",
							RequestBody: "test",
						},
					},
					Request: models.CustomRequest{Headers: make(map[string]string)},
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "step 4 - executing",
			model: models.Model{
				CustomRequestModel: models.CustomRequestModel{
					Step:             4,
					MethodInput:      methodTi,
					EndpointInput:    endpointTi,
					HeaderKeyInput:   headerKeyTi,
					HeaderValueInput: headerValueTi,
					BodyInput:        bodyTi,
					FilterInput:      filterTi,
					Table:            tbl,
					Request: models.CustomRequest{
						Method:   "GET",
						Endpoint: "http://example.com",
						Headers:  make(map[string]string),
					},
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "step 5 - results success",
			model: models.Model{
				CustomRequestModel: models.CustomRequestModel{
					Step:             5,
					MethodInput:      methodTi,
					EndpointInput:    endpointTi,
					HeaderKeyInput:   headerKeyTi,
					HeaderValueInput: headerValueTi,
					BodyInput:        bodyTi,
					FilterInput:      filterTi,
					Table:            tbl,
					Result: &models.TestResult{
						Method:   "GET",
						Endpoint: "http://example.com",
						Status:   "200",
						Message:  "OK",
						Duration: 100 * time.Millisecond,
					},
					Request: models.CustomRequest{Headers: make(map[string]string)},
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "step 5 - results error status",
			model: models.Model{
				CustomRequestModel: models.CustomRequestModel{
					Step:             5,
					MethodInput:      methodTi,
					EndpointInput:    endpointTi,
					HeaderKeyInput:   headerKeyTi,
					HeaderValueInput: headerValueTi,
					BodyInput:        bodyTi,
					FilterInput:      filterTi,
					Table:            tbl,
					Result: &models.TestResult{
						Method:   "GET",
						Endpoint: "http://example.com",
						Status:   "404",
						Message:  "Not Found",
						Duration: 50 * time.Millisecond,
					},
					Request: models.CustomRequest{Headers: make(map[string]string)},
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "step 2 with headers populated",
			model: models.Model{
				CustomRequestModel: models.CustomRequestModel{
					Step:             2,
					MethodInput:      methodTi,
					EndpointInput:    endpointTi,
					HeaderKeyInput:   headerKeyTi,
					HeaderValueInput: headerValueTi,
					BodyInput:        bodyTi,
					FilterInput:      filterTi,
					Table:            tbl,
					Request: models.CustomRequest{
						Method:   "POST",
						Endpoint: "http://example.com/api",
						Headers: map[string]string{
							"Content-Type": "application/json",
							"Authorization": "Bearer token",
						},
					},
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "step 3 with headers",
			model: models.Model{
				CustomRequestModel: models.CustomRequestModel{
					Step:             3,
					MethodInput:      methodTi,
					EndpointInput:    endpointTi,
					HeaderKeyInput:   headerKeyTi,
					HeaderValueInput: headerValueTi,
					BodyInput:        bodyTi,
					FilterInput:      filterTi,
					Table:            tbl,
					Request: models.CustomRequest{
						Method:   "POST",
						Endpoint: "http://example.com/api",
						Headers: map[string]string{
							"Content-Type": "application/json",
						},
					},
				},
				Width:  100,
				Height: 50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := ViewCustomRequest(tt.model)
			if output == "" {
				t.Errorf("ViewCustomRequest returned empty string for %s", tt.name)
			}
		})
	}
}

func TestViewEndpointSelector(t *testing.T) {
	searchTi := textinput.New()
	searchTi.SetValue("users")

	endpoints := []models.EndpointInfo{
		{
			Path:        "/users",
			Method:      "GET",
			OperationID: "getUsers",
			Tags:        []string{"users"},
			Summary:     "Get all users",
			Selected:    true,
		},
		{
			Path:        "/users/{id}",
			Method:      "GET",
			OperationID: "getUserById",
			Tags:        []string{"users"},
			Summary:     "Get user by ID",
			Selected:    false,
		},
	}

	tests := []struct {
		name  string
		model models.Model
	}{
		{
			name: "with endpoints",
			model: models.Model{
				EndpointSelectorModel: models.EndpointSelectorModel{
					SearchInput:       searchTi,
					AllEndpoints:      endpoints,
					FilteredEndpoints: endpoints,
					Cursor:            0,
					Ready:             true,
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "not ready",
			model: models.Model{
				EndpointSelectorModel: models.EndpointSelectorModel{
					SearchInput:       searchTi,
					AllEndpoints:      []models.EndpointInfo{},
					FilteredEndpoints: []models.EndpointInfo{},
					Cursor:            0,
					Ready:             false,
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "with error",
			model: models.Model{
				EndpointSelectorModel: models.EndpointSelectorModel{
					SearchInput:       searchTi,
					AllEndpoints:      []models.EndpointInfo{},
					FilteredEndpoints: []models.EndpointInfo{},
					Cursor:            0,
					Ready:             true,
					Err:               fmt.Errorf("test error"),
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "empty results",
			model: models.Model{
				EndpointSelectorModel: models.EndpointSelectorModel{
					SearchInput:       searchTi,
					AllEndpoints:      endpoints,
					FilteredEndpoints: []models.EndpointInfo{},
					Cursor:            0,
					Ready:             true,
				},
				Width:  100,
				Height: 50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := ViewEndpointSelector(tt.model)
			if output == "" {
				t.Errorf("ViewEndpointSelector returned empty string for %s", tt.name)
			}
		})
	}
}

func TestViewConfigEditor(t *testing.T) {
	specPathTi := textinput.New()
	specPathTi.SetValue("openapi.yaml")
	
	baseURLTi := textinput.New()
	baseURLTi.SetValue("http://example.com")
	
	authTypeTi := textinput.New()
	authTypeTi.SetValue("bearer")
	
	tokenTi := textinput.New()
	apiKeyNameTi := textinput.New()
	apiKeyInTi := textinput.New()
	usernameTi := textinput.New()
	passwordTi := textinput.New()
	maxConcurrTi := textinput.New()
	verboseTi := textinput.New()
	maxRetriesTi := textinput.New()
	retryDelayTi := textinput.New()

	tests := []struct {
		name  string
		model models.Model
	}{
		{
			name: "default view",
			model: models.Model{
				ConfigEditorModel: models.ConfigEditorModel{
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
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "with validation error",
			model: models.Model{
				ConfigEditorModel: models.ConfigEditorModel{
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
					ValidationError: "Invalid URL format",
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "focus on different field",
			model: models.Model{
				ConfigEditorModel: models.ConfigEditorModel{
					FocusedField:    5,
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
				},
				Width:  100,
				Height: 50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := ViewConfigEditor(tt.model)
			if output == "" {
				t.Errorf("ViewConfigEditor returned empty string for %s", tt.name)
			}
		})
	}
}

func TestViewValidateExtended(t *testing.T) {
	ti := textinput.New()
	ti.SetValue("openapi.yaml")

	tests := []struct {
		name  string
		model models.Model
	}{
		{
			name: "with result",
			model: models.Model{
				ValidateModel: models.ValidateModel{
					TextInput: ti,
					Result:    "Validation successful",
					Done:      true,
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "with error",
			model: models.Model{
				ValidateModel: models.ValidateModel{
					TextInput: ti,
					Err:       fmt.Errorf("validation error"),
					Done:      true,
				},
				Width:  100,
				Height: 50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := ViewValidate(tt.model)
			if output == "" {
				t.Errorf("ViewValidate returned empty string for %s", tt.name)
			}
		})
	}
}

func TestViewTestExtended(t *testing.T) {
	specTi := textinput.New()
	specTi.SetValue("openapi.yaml")
	
	urlTi := textinput.New()
	urlTi.SetValue("http://example.com")
	
	filterTi := textinput.New()

	columns := []table.Column{
		{Title: "Method", Width: 8},
		{Title: "Endpoint", Width: 40},
		{Title: "Status", Width: 10},
		{Title: "Message", Width: 30},
	}
	tbl := table.New(table.WithColumns(columns))

	results := []models.TestResult{
		{
			Method:   "GET",
			Endpoint: "/users",
			Status:   "200",
			Message:  "OK",
			Duration: 100 * time.Millisecond,
		},
	}

	tests := []struct {
		name  string
		model models.Model
	}{
		{
			name: "step 0 - spec input",
			model: models.Model{
				TestModel: models.TestModel{
					Step:        0,
					SpecInput:   specTi,
					UrlInput:    urlTi,
					FilterInput: filterTi,
					Table:       tbl,
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "step 1 - url input",
			model: models.Model{
				TestModel: models.TestModel{
					Step:        1,
					SpecInput:   specTi,
					UrlInput:    urlTi,
					FilterInput: filterTi,
					Table:       tbl,
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "testing state",
			model: models.Model{
				TestModel: models.TestModel{
					Step:        0,
					SpecInput:   specTi,
					UrlInput:    urlTi,
					FilterInput: filterTi,
					Table:       tbl,
					Testing:     true,
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "with results",
			model: models.Model{
				TestModel: models.TestModel{
					Step:        0,
					SpecInput:   specTi,
					UrlInput:    urlTi,
					FilterInput: filterTi,
					Table:       tbl,
					Results:     results,
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "with error",
			model: models.Model{
				TestModel: models.TestModel{
					Step:        0,
					SpecInput:   specTi,
					UrlInput:    urlTi,
					FilterInput: filterTi,
					Table:       tbl,
					Err:         fmt.Errorf("test error"),
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "with export success",
			model: models.Model{
				TestModel: models.TestModel{
					Step:          0,
					SpecInput:     specTi,
					UrlInput:      urlTi,
					FilterInput:   filterTi,
					Table:         tbl,
					Results:       results,
					ExportSuccess: "Exported to file.json",
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "showing log",
			model: models.Model{
				TestModel: models.TestModel{
					Step:        0,
					SpecInput:   specTi,
					UrlInput:    urlTi,
					FilterInput: filterTi,
					Table:       tbl,
					Results: []models.TestResult{
						{
							LogEntry: &models.LogEntry{
								RequestURL:  "http://example.com",
								RequestBody: "test",
							},
						},
					},
					ShowingLog:  true,
					SelectedLog: 0,
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "filter active",
			model: models.Model{
				TestModel: models.TestModel{
					Step:          0,
					SpecInput:     specTi,
					UrlInput:      urlTi,
					FilterInput:   filterTi,
					Table:         tbl,
					Results:       results,
					FilterActive:  true,
					FilteredResults: results,
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "select endpoints mode",
			model: models.Model{
				TestModel: models.TestModel{
					Step:            1,
					SpecInput:       specTi,
					UrlInput:        urlTi,
					FilterInput:     filterTi,
					Table:           tbl,
					SelectEndpoints: true,
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "step 2 - testing in progress",
			model: models.Model{
				TestModel: models.TestModel{
					Step:        2,
					SpecInput:   specTi,
					UrlInput:    urlTi,
					FilterInput: filterTi,
					Table:       tbl,
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "step 3 - results with stats",
			model: models.Model{
				TestModel: models.TestModel{
					Step:        3,
					SpecInput:   specTi,
					UrlInput:    urlTi,
					FilterInput: filterTi,
					Table:       tbl,
					Results:     results,
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "step 3 - with export message",
			model: models.Model{
				TestModel: models.TestModel{
					Step:          3,
					SpecInput:     specTi,
					UrlInput:      urlTi,
					FilterInput:   filterTi,
					Table:         tbl,
					Results:       results,
					ExportSuccess: "Exported successfully",
				},
				Width:       100,
				Height:      50,
				VerboseMode: true,
			},
		},
		{
			name: "step 3 - filter active with query",
			model: models.Model{
				TestModel: models.TestModel{
					Step:            3,
					SpecInput:       specTi,
					UrlInput:        urlTi,
					FilterInput:     filterTi,
					Table:           tbl,
					Results:         results,
					FilterActive:    true,
					FilteredResults: results,
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "step 4 - log detail",
			model: models.Model{
				TestModel: models.TestModel{
					Step:        4,
					SpecInput:   specTi,
					UrlInput:    urlTi,
					FilterInput: filterTi,
					Table:       tbl,
					Results: []models.TestResult{
						{
							Method:   "GET",
							Endpoint: "/test",
							LogEntry: &models.LogEntry{
								RequestURL: "http://test.com",
							},
						},
					},
					SelectedLog: 0,
				},
				Width:  100,
				Height: 50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := ViewTest(tt.model)
			if output == "" {
				t.Errorf("ViewTest returned empty string for %s", tt.name)
			}
		})
	}
}

func TestViewHistoryExtended(t *testing.T) {
	entries := []models.HistoryEntry{
		{
			ID:         "test1",
			Timestamp:  time.Now(),
			SpecPath:   "openapi.yaml",
			BaseURL:    "http://example.com",
			TotalTests: 10,
			Passed:     8,
			Failed:     2,
			Duration:   "1.5s",
			Results: []models.TestResult{
				{Method: "GET", Endpoint: "/users", Status: "200", Message: "OK"},
			},
		},
	}

	tests := []struct {
		name  string
		model models.Model
	}{
		{
			name: "with entries",
			model: models.Model{
				History: &models.TestHistory{
					Entries: entries,
				},
				HistoryIndex: 0,
				Width:        100,
				Height:       50,
			},
		},
		{
			name: "selected entry",
			model: models.Model{
				History: &models.TestHistory{
					Entries: entries,
				},
				HistoryIndex: 0,
				Width:        100,
				Height:       50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := ViewHistory(tt.model)
			if output == "" {
				t.Errorf("ViewHistory returned empty string for %s", tt.name)
			}
		})
	}
}

func TestViewMenuExtended(t *testing.T) {
	tests := []struct {
		name  string
		model models.Model
	}{
		{
			name: "verbose mode on",
			model: models.Model{
				Cursor:      0,
				VerboseMode: true,
				Width:       100,
				Height:      50,
			},
		},
		{
			name: "with config",
			model: models.Model{
				Cursor: 1,
				Config: models.Config{
					BaseURL:  "http://example.com",
					SpecPath: "openapi.yaml",
				},
				Width:  100,
				Height: 50,
			},
		},
		{
			name: "different cursor positions",
			model: models.Model{
				Cursor: 5,
				Width:  100,
				Height: 50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := ViewMenu(tt.model)
			if output == "" {
				t.Errorf("ViewMenu returned empty string for %s", tt.name)
			}
		})
	}
}

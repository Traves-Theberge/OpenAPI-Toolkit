package models

import (
"time"

"github.com/charmbracelet/bubbles/spinner"
"github.com/charmbracelet/bubbles/table"
"github.com/charmbracelet/bubbles/textinput"
)

// Screen represents different views in the application
type Screen int

const (
MenuScreen Screen = iota
HelpScreen
ValidateScreen
TestScreen
CustomRequestScreen
HistoryScreen
)

// Model is the main application state
type Model struct {
Cursor             int
Screen             Screen
Width              int
Height             int
VerboseMode        bool
Config             Config
ValidateModel      ValidateModel
TestModel          TestModel
CustomRequestModel CustomRequestModel
History            *TestHistory
HistoryIndex       int  // Selected index in history view
}

// ValidateModel holds state for the validation screen
type ValidateModel struct {
TextInput textinput.Model
Err       error
Result    string
Done      bool
}

// TestModel holds state for the testing screen
type TestModel struct {
Step          int
SpecInput     textinput.Model
UrlInput      textinput.Model
Spinner       spinner.Model
Table         table.Model
Results       []TestResult
Err           error
Testing       bool
ExportSuccess string
ShowingLog    bool
SelectedLog   int
FilterActive  bool
FilterInput   textinput.Model
FilteredResults []TestResult
TestStartTime time.Time  // Track when test run started for history
}

// CustomRequestModel holds state for the custom request screen
type CustomRequestModel struct {
Step             int
MethodInput      textinput.Model
EndpointInput    textinput.Model
HeaderKeyInput   textinput.Model
HeaderValueInput textinput.Model
BodyInput        textinput.Model
Spinner          spinner.Model
Table            table.Model
Request          CustomRequest
Result           *TestResult
Err              error
Testing          bool
ExportSuccess    string
ShowingLog       bool
FilterActive     bool
FilterInput      textinput.Model
}

// CustomRequest holds a manually created API request
type CustomRequest struct {
Method      string
Endpoint    string
Headers     map[string]string
Body        string
QueryParams map[string]string
IsCustom    bool // Flag for history tracking
}

// TestResult represents the result of testing an API endpoint
type TestResult struct {
Method   string
Endpoint string
Status   string
Message  string
Duration time.Duration
LogEntry *LogEntry
}

// LogEntry captures detailed request/response information
type LogEntry struct {
RequestURL      string
RequestHeaders  map[string]string
RequestBody     string
ResponseHeaders map[string]string
ResponseBody    string
Duration        time.Duration
Timestamp       time.Time
}

// ValidationResult contains OpenAPI validation results
type ValidationResult struct {
	Valid          bool
	StatusValid    bool
	ContentType    string
	SchemaErrors   []string
	ExpectedStatus string
}// AuthConfig holds authentication configuration
type AuthConfig struct {
AuthType   string
Token      string
APIKeyIn   string
APIKeyName string
Username   string
Password   string
}

// Config holds application configuration
type Config struct {
BaseURL        string
SpecPath       string
VerboseMode    bool
Auth           *AuthConfig
MaxConcurrency int  // Maximum number of concurrent test requests (0 = auto-detect)
}

// ConfigFile represents the YAML configuration file structure
type ConfigFile struct {
BaseURL        string `yaml:"baseUrl"`
SpecPath       string `yaml:"specPath"`
VerboseMode    bool   `yaml:"verboseMode"`
MaxConcurrency int    `yaml:"maxConcurrency,omitempty"`
Auth           *struct {
Type       string `yaml:"type"`
Token      string `yaml:"token,omitempty"`
APIKeyIn   string `yaml:"apiKeyIn,omitempty"`
APIKeyName string `yaml:"apiKeyName,omitempty"`
Username   string `yaml:"username,omitempty"`
Password   string `yaml:"password,omitempty"`
} `yaml:"auth,omitempty"`
}

// ExportResult represents a single test result in the export format
type ExportResult struct {
Method   string `json:"method"`
Endpoint string `json:"endpoint"`
Status   string `json:"status"`
Message  string `json:"message"`
Duration string `json:"duration"`
}

// ExportData represents the complete export structure
type ExportData struct {
Timestamp  string         `json:"timestamp"`
SpecPath   string         `json:"specPath"`
BaseURL    string         `json:"baseUrl"`
TotalTests int            `json:"totalTests"`
Passed     int            `json:"passed"`
Failed     int            `json:"failed"`
Results    []ExportResult `json:"results"`
}

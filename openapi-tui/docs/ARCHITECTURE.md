# OpenAPI CLI TUI Architecture

## Overview

This is a **single-program Bubble Tea TUI** that implements a professional OpenAPI testing tool with screen-based navigation. The architecture follows the **Model-View-Update (MVU)** pattern with embedded sub-models for feature isolation.

## Core Components

```mermaid
graph TD
    A[main.go] --> B[Bubble Tea Framework]
    A --> C[Charmbracelet Libraries]
    A --> D[Kin OpenAPI]

    B --> E[Init/Update/View Methods]
    B --> F[Message Passing]
    B --> G[Commands]

    C --> H[Lip Gloss - Styling]
    C --> I[Bubbles - Components]
    C --> J[Kin OpenAPI - Spec Parsing]
```

## User Flow Diagrams

### Main Navigation Flow

```mermaid
stateDiagram-v2
    [*] --> MenuScreen
    MenuScreen --> ValidateScreen: Enter (cursor=0)
    MenuScreen --> TestScreen: Enter (cursor=1)
    MenuScreen --> HelpScreen: h or ?
    MenuScreen --> [*]: Enter (cursor=3) or q/Ctrl+C

    ValidateScreen --> MenuScreen: Esc or Enter (after validation)
    TestScreen --> MenuScreen: Esc or Enter (after testing)
    TestScreen --> HistoryScreen: r (from results view)
    HistoryScreen --> TestScreen: Esc or Enter (replay test)
    HelpScreen --> MenuScreen: q/Esc/h/?
```

### Validation Screen Flow

```mermaid
flowchart TD
    A[User Input] --> B{Path Empty?}
    B -->|Yes| C[Show Error]
    B -->|No| D[Load Spec]
    D --> E{Load Success?}
    E -->|No| F[Show Load Error]
    E -->|Yes| G[Validate Spec]
    G --> H{Valid?}
    H -->|No| I[Show Validation Error]
    H -->|Yes| J[Show Success]
    C --> K[Enter/Esc to continue]
    F --> K
    I --> K
    J --> K
```

### Testing Screen Flow

```mermaid
stateDiagram-v2
    [*] --> Step0: Spec File Input
    Step0 --> Step1: Enter (valid path)
    Step0 --> [*]: Esc

    Step1 --> Step2: Enter (valid URL)
    Step1 --> [*]: Esc

    Step2 --> Step3: Testing Complete
    Step2 --> [*]: Esc/Ctrl+C

    Step3 --> [*]: Enter/Esc
```

### Message Flow Architecture

```mermaid
sequenceDiagram
    participant U as User
    participant BT as Bubble Tea
    participant U as Update()
    participant S as Screen Handler
    participant M as Model

    U->>BT: Key Input
    BT->>U: tea.KeyMsg
    U->>S: Route to Screen
    S->>M: Update Model
    M->>BT: New Model + Cmd
    BT->>U: Re-render View
```

## Screen State Management

The application uses a **finite state machine** with four screens:

```go
type screen int
const (
    menuScreen screen = iota // Main menu with navigation options
    helpScreen               // Help/documentation screen
    validateScreen           // OpenAPI spec validation screen
    testScreen               // API endpoint testing screen
)
```

### Screen Transitions

```mermaid
flowchart TD
    A[Menu Screen] --> B[Validate Screen]
    A --> C[Test Screen]
    A --> D[Quit]
    A --> E[Help Screen]

    B --> A
    C --> A
    E --> A

    D
```

## Embedded Model Architecture

The main model embeds specialized sub-models for each feature:

```go
type model struct {
    screen         screen
    cursor         int
    width, height  int
    validateModel  validateModel  // Embedded validation state
    testModel      testModel       // Embedded testing state
}
```

### Sub-Model Responsibilities

**validateModel:**
- Text input for spec file path
- Validation result/error state
- Completion flag

**testModel:**
- Multi-step workflow state (0-3)
- Spec file input
- Base URL input
- Spinner for progress
- Results table
- Testing flag

## Async Operations with Commands

The testing workflow uses Bubble Tea commands for non-blocking operations:

```go
func runTestCmd(specPath, baseURL string) tea.Cmd {
    return func() tea.Msg {
        results, err := runTests(specPath, baseURL)
        return testResultMsg{results: results, err: err}
    }
}
```

### Command Flow

```mermaid
sequenceDiagram
    participant UI as User Input
    participant CMD as runTestCmd()
    participant GO as Goroutine
    participant HTTP as HTTP Requests
    participant MSG as testResultMsg

    UI->>CMD: Execute Command
    CMD->>GO: Start Goroutine
    GO->>HTTP: Make API Calls
    HTTP-->>GO: Response Data
    GO->>MSG: Create Message
    MSG->>UI: Update UI
```

## Dynamic Responsive Layout

The application adapts to terminal size changes:

```go
case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
    return m, nil
```

All screens use `lipgloss.Place()` for centering:

```go
return lipgloss.Place(
    m.width, m.height,
    lipgloss.Center, lipgloss.Center,
    borderedContent,
    lipgloss.WithWhitespaceChars(" "),
    lipgloss.WithWhitespaceForeground(lipgloss.Color("#333")),
)
```

## Business Logic Layer

### Validation Logic

```mermaid
flowchart TD
    A["validateSpec()"] --> B["Load OpenAPI Document"]
    B --> C{"Load Success?"}
    C -->|No| D["Return Load Error"]
    C -->|Yes| E["Validate Document"]
    E --> F{"Valid?"}
    F -->|No| G["Return Validation Error"]
    F -->|Yes| H["Return Success Message"]
```

### Testing Logic

```mermaid
flowchart TD
    A["runTests()"] --> B["Load OpenAPI Spec"]
    B --> C{"Load Success?"}
    C -->|No| D["Return Load Error"]
    C -->|Yes| E["Iterate Paths & Operations"]
    E --> F["replacePlaceholders()"]
    F --> G["buildQueryParams()"]
    G --> H{"POST/PUT/PATCH?"}
    H -->|Yes| I["generateRequestBody()"]
    H -->|No| J["Test Endpoint"]
    I --> K["generateSampleFromSchema()"]
    K --> J
    J --> L{"Test Success?"}
    L -->|No| M["Record Error"]
    L -->|Yes| N["Record Success"]
    M --> O{"More Endpoints?"}
    N --> O
    O -->|Yes| E
    O -->|No| P["Return Results"]
```

## Request Body Generation

The application intelligently generates request bodies for POST/PUT/PATCH operations:

```mermaid
flowchart TD
    A["generateRequestBody()"] --> B{"RequestBody exists?"}
    B -->|No| C["Return nil"]
    B -->|Yes| D["Get JSON Schema"]
    D --> E{"Schema exists?"}
    E -->|No| C
    E -->|Yes| F["generateSampleFromSchema()"]
    F --> G{"Schema Type?"}
    G -->|Object| H["Iterate Properties"]
    G -->|Array| I["Generate Single Item"]
    G -->|String| J["Check Format/Example"]
    G -->|Integer| K["Use Min or 1"]
    G -->|Number| L["Use Min or 1.0"]
    G -->|Boolean| M["Return true"]
    H --> N["Recursive Call for Each Property"]
    I --> N
    J --> O["Return Sample Value"]
    K --> O
    L --> O
    M --> O
    N --> O
    O --> P["Marshal to JSON"]
    P --> Q["Return Bytes"]
```

### Schema Type Handling

The request body generator supports:

- **Primitives**: string, integer, number, boolean
- **Complex Types**: object (nested), array
- **String Formats**: email, uri/url, date, date-time
- **Schema Features**: example, default, enum, min/max
- **Nested Structures**: Recursive generation for deep schemas

## Query Parameter Handling

Query parameters are automatically extracted from the OpenAPI spec and appended to URLs:

```mermaid
flowchart TD
    A["buildQueryParams()"] --> B{"Parameters exist?"}
    B -->|No| C["Return empty string"]
    B -->|Yes| D["Iterate Parameters"]
    D --> E{"In: query?"}
    E -->|No| D
    E -->|Yes| F{"Schema Type?"}
    F -->|String| G["Use example or 'sample'"]
    F -->|Integer/Number| H["Use example or 1"]
    F -->|Boolean| I["Use example or true"]
    F -->|Array| J["Use example or ['item']"]
    G --> K["Append param=value"]
    H --> K
    I --> K
    J --> K
    K --> L{"More Params?"}
    L -->|Yes| D
    L -->|No| M["Join with &"]
    M --> N["Return ?param1=val1&param2=val2"]
```

## Path Parameter Handling

Path parameters (e.g., `/users/{id}`) are automatically replaced with sample values:

```go
// replacePlaceholders() uses regex to replace {param} with "1"
"/users/{id}" → "/users/1"
"/users/{userId}/posts/{postId}" → "/users/1/posts/1"
```

## Styling System

The application uses a consistent color scheme:

- **Primary**: `#4ECDC4` (Teal) - Success, borders, cells
- **Secondary**: `#7D56F4` (Purple) - Headers, titles
- **Accent**: `#FF6B6B` (Red) - Errors, selected items
- **Background**: `#FAFAFA` (Light Gray) - Headers
- **Text**: `#888` (Gray) - Instructions, status

## Error Handling

Comprehensive error handling at all levels:

- **Input Validation**: Empty fields, invalid paths
- **File Operations**: Spec loading failures
- **Network Requests**: Connection timeouts, invalid responses
- **Spec Validation**: OpenAPI compliance issues
- **UI State**: Graceful fallbacks for edge cases

## Key Design Decisions

1. **Single Program Architecture**: Eliminates terminal window breaking by keeping everything in one Bubble Tea program
2. **Screen State Management**: Clean separation of concerns with dedicated handlers per screen
3. **Embedded Models**: Feature isolation while maintaining shared state
4. **Async Commands**: Non-blocking operations for smooth UX during testing
5. **Responsive Design**: Adapts to any terminal size with dynamic centering
6. **Professional Styling**: Consistent color scheme and modern UI elements
7. **Intelligent Request Generation**: Automatic body and parameter generation from OpenAPI schemas
8. **HTTP Timeout Protection**: 10-second timeout prevents hanging on unresponsive endpoints
9. **Comprehensive Test Coverage**: 30.3% coverage with table-driven tests and edge cases

## Component Architecture

```mermaid
classDiagram
    class model {
        +screen: screen
        +cursor: int
        +width: int
        +height: int
        +validateModel: validateModel
        +testModel: testModel
        +Init(): tea.Cmd
        +Update(msg): (tea.Model, tea.Cmd)
        +View(): string
    }

    class validateModel {
        +textInput: textinput.Model
        +err: error
        +result: string
        +done: bool
    }

    class testModel {
        +step: int
        +specInput: textinput.Model
        +urlInput: textinput.Model
        +spinner: spinner.Model
        +table: table.Model
        +err: error
        +result: string
        +done: bool
        +testing: bool
        +results: []testResult
    }

    class testResult {
        +method: string
        +endpoint: string
        +status: string
        +message: string
    }

    model *-- validateModel
    model *-- testModel
    testModel *-- testResult
```

## Testing Functions

Key business logic functions with test coverage:

### Core Testing Functions
- `runTests(specPath, baseURL, auth)` - Main testing orchestrator with auth support
- `testEndpoint(method, url, body, auth)` - HTTP client with 10s timeout and authentication
- `applyAuth(req, auth)` - Apply authentication to HTTP requests
- `replacePlaceholders(path)` - Replace `{id}` with `"1"` using regex
- `buildQueryParams(operation)` - Generate query strings from spec
- `generateRequestBody(operation)` - Create JSON from schema
- `generateSampleFromSchema(schema)` - Recursive schema-to-sample converter
- `validateResponse(resp, operation, status)` - Validate response against spec

### Validation Functions
- `validateSpec(specPath)` - OpenAPI spec validation
- `validateResponse(...)` - Response status code and content-type validation

### Test Coverage
- **Current Coverage**: 37.8% of statements
- **Test Files**: 11 test functions, 70+ test cases
- **Test Types**: Unit tests, integration tests, table-driven tests
- **Edge Cases**: Nil handling, invalid inputs, nested structures, response validation, authentication scenarios

## Authentication Support

The application supports multiple authentication methods for protected APIs:

```mermaid
flowchart TD
    A["applyAuth(req, auth)"] --> B{"auth is nil?"}
    B -->|Yes| C["Return without modification"]
    B -->|No| D{"authType?"}
    D -->|"none" or ""| C
    D -->|"bearer"| E["Set Authorization: Bearer token"]
    D -->|"apiKey"| F{"apiKeyIn?"}
    D -->|"basic"| G["Set Basic Auth header"]
    F -->|"header"| H["Set custom header"]
    F -->|"query"| I["Add query parameter"]
    E --> J["Continue with request"]
    H --> J
    I --> J
    G --> J
```

### Authentication Types

1. **Bearer Token**
   - Sets `Authorization: Bearer <token>` header
   - Common for OAuth 2.0 and JWT authentication
   - Example: `Authorization: Bearer eyJhbGciOiJIUzI1...`

2. **API Key**
   - Supports both header and query parameter placement
   - Header: Sets custom header name (e.g., `X-API-Key: abc123`)
   - Query: Adds to URL (e.g., `?api_key=abc123`)
   - Flexible configuration for different API providers

3. **Basic Authentication**
   - Standard HTTP Basic Auth with username/password
   - Uses Go's `SetBasicAuth()` for proper encoding
   - Example: `Authorization: Basic dXNlcjpwYXNz`

4. **None**
   - Default for public/unauthenticated APIs
   - No headers or parameters added

### Implementation Details

**authConfig Structure:**
```go
type authConfig struct {
    authType   string // "bearer", "apiKey", "basic", "none"
    token      string // For bearer and apiKey
    apiKeyIn   string // "header" or "query"
    apiKeyName string // Custom header/param name
    username   string // For basic auth
    password   string // For basic auth
}
```

**Integration:**
- `applyAuth()` is called in `testEndpoint()` before HTTP execution
- All test functions accept optional `auth *authConfig` parameter
- Nil auth defaults to no authentication
- Function signatures updated across 12 call sites

### Future Enhancements
- UI for collecting authentication details
- Support for OAuth 2.0 flows
- Certificate-based authentication
- Auth credential persistence

## Response Validation

The application validates API responses against OpenAPI specifications to detect mismatches:

```mermaid
flowchart TD
    A["validateResponse()"] --> B{"Operation has responses?"}
    B -->|No| C["Mark as valid - no spec"]
    B -->|Yes| D["Look for status code"]
    D --> E{"Exact match found?"}
    E -->|Yes| F["Use matched response"]
    E -->|No| G{"Default response exists?"}
    G -->|Yes| H["Use default response"]
    G -->|No| I["Mark invalid - status not in spec"]
    F --> J["Validate Content-Type"]
    H --> J
    J --> K{"Content-Type defined?"}
    K -->|No| L["Skip content validation"]
    K -->|Yes| M{"Content-Type matches?"}
    M -->|Yes| N["Mark as valid"]
    M -->|No| O["Mark invalid - wrong content-type"]
    L --> N
    I --> P["Return validation errors"]
    O --> P
    N --> Q["Return success"]
```

### Validation Checks

1. **Status Code Validation**
   - Checks if response status matches spec
   - Falls back to "default" response if defined
   - Reports "status X not defined in spec" errors

2. **Content-Type Validation**
   - Extracts base content type (ignores charset)
   - Compares against spec's response content types
   - Handles common cases like `application/json; charset=utf-8`

3. **Error Reporting**
   - Collects all validation errors
   - Shows first error in test results
   - Marks successful validations as "OK (validated)"

### Future Enhancements
- JSON schema validation against response body
- Header validation
- Response time tracking

## Recent Enhancements

### Phase 1 Implementation (5 of 5 Complete) ✅✅✅

1. **Unit Tests (#1)** ✅
   - Created comprehensive test suite
   - Table-driven tests for all core functions
   - Coverage progression: 0% → 21.9% → 30.3% → 35.2% → 37.8%

2. **Request Body Generation (#2)** ✅
   - Intelligent JSON generation from OpenAPI schemas
   - Support for nested objects, arrays, primitives
   - Format-aware generation (email, date, uri)
   - Example and default value usage
   - 14 test cases covering all schema types

3. **Response Schema Validation (#3)** ✅
   - Validates status codes against spec
   - Content-type verification
   - Default response fallback support
   - 7 test cases with edge case coverage
   - Integrated into testing pipeline

4. **Query Parameter Handling (#4)** ✅
   - Automatic extraction from operation specs
   - Type-aware sample value generation
   - Proper URL encoding with `?` and `&`
   - Integration with testing pipeline

5. **Authentication Support (#5)** ✅
   - Bearer token, API key, and Basic auth support
   - Flexible header/query parameter placement for API keys
   - 11 comprehensive auth tests (7 unit + 4 integration)
   - Function signatures updated across entire codebase
   - Ready for UI integration

### Phase 1 Complete! 🎉

**Achievements:**
- Test coverage: 0% → 37.8% (nearly 40%!)
- Test-to-code ratio: 0.99:1 (1128 test lines for 1140 code lines)
- All 5 critical foundation features delivered
- 70+ test cases, all passing
- Robust, production-ready testing infrastructure

### Phase 2 Implementation (10 of 15 Complete) 🚀

**Completed Features:**

1. **Enhanced Error Messages** ✅ - Actionable suggestions with context-aware guidance
2. **Verbose Logging** ✅ - Full HTTP transaction details with 'l' key
3. **Configuration Persistence** ✅ - Auto-save/load from `~/.config/openapi-tui/config.yaml`
4. **JSON Export** ✅ - CI/CD integration with metadata and statistics
5. **Standard Go Layout** ✅ - Modular cmd/ and internal/ package structure
6. **Summary Statistics** ✅ - Pass rates, timing analysis, performance metrics
7. **Response Filtering** ✅ - Multi-field filtering with special keywords
8. **HTML Export** ✅ - Professional web reports with embedded CSS
9. **JUnit XML Export** ✅ - CI/CD pipeline integration (Jenkins, GitLab, GitHub Actions)
10. **Request History** ✅ - Track, replay, and analyze test runs with persistence

**Architecture Updates:**
- **Package Structure**: Migrated from monolithic 1,794-line file to modular packages
  - `cmd/openapi-tui/` - Application entry point
  - `internal/models/` - Shared type definitions (Screen, Model, TestResult, HistoryEntry)
  - `internal/config/` - Configuration management
  - `internal/errors/` - Enhanced error handling
  - `internal/export/` - JSON, HTML, JUnit XML export
  - `internal/validation/` - OpenAPI validation
  - `internal/testing/` - API testing logic
  - `internal/ui/` - View rendering (stats, filter, history views)

**New Screen: HistoryScreen**
- Displays table of past test runs
- Navigation with arrow keys or vim keys (j/k)
- Select and replay tests with Enter
- Persistent storage in `~/.config/openapi-tui/history.json`
- Automatic 50-entry limit (keeps most recent)

**Key Bindings:**
- `v` - Toggle verbose mode
- `f` - Filter results
- `e` - Export JSON
- `h` - Export HTML
- `j` - Export JUnit XML
- `r` - View history
- `l` - View detailed logs

**Test Coverage:** 94 tests passing (170+ test runs including subtests)

### Remaining Phase 2 Features (5 of 15)

- Parallel test execution
- Custom request editing
- Configuration UI
- Endpoint search
- Test retries

This architecture provides a robust, maintainable, and production-ready TUI for OpenAPI testing. The modular structure enables easy feature additions and maintenance. Phase 2 enhancements deliver professional-grade developer experience with comprehensive export formats, historical analysis, and CI/CD integration.

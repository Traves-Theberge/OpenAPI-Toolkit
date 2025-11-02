# OpenAPI Toolkit - CLI Architecture

## Overview

The OpenAPI CLI is a command-line tool built with TypeScript and Node.js for validating OpenAPI specifications and testing API endpoints. It's designed for automation, CI/CD pipelines, and scripting workflows.

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      CLI Entry Point                         │
│                       (cli.ts)                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  validate    │  │     test     │  │     help     │      │
│  │   command    │  │   command    │  │   command    │      │
│  └──────┬───────┘  └──────┬───────┘  └──────────────┘      │
└─────────┼──────────────────┼──────────────────────────────────┘
          │                  │
          ▼                  ▼
┌─────────────────┐  ┌─────────────────────────────────────┐
│   validate.ts   │  │          test.ts                    │
│                 │  │                                     │
│ - Load spec     │  │ - Load spec                         │
│ - Parse YAML    │  │ - Build URLs                        │
│ - Validate      │  │ - Test endpoints                    │
│ - Show errors   │  │ - Collect results                   │
│   with          │  │ - Export JSON                       │
│   suggestions   │  │ - Exit codes                        │
└─────────────────┘  └─────────────────────────────────────┘
          │                  │
          ▼                  ▼
    ┌─────────────────────────────────┐
    │       File System (fs)          │
    │   - Read OpenAPI specs          │
    │   - Write JSON exports          │
    └─────────────────────────────────┘
                     │
                     ▼
              ┌─────────────┐
              │    Axios    │
              │ HTTP Client │
              └─────────────┘
```

## Component Architecture

### 1. CLI Layer (`src/cli.ts`)

**Responsibility:** Command parsing and routing

**Key Components:**
- Commander.js for argument parsing
- Command definitions (validate, test)
- Option flags (`-v`, `-e`)
- Error handling and exit codes

**Flow:**
```
User Input → Commander.js → Command Handler → Exit Code
```

### 2. Validation Layer (`src/commands/validate.ts`)

**Responsibility:** OpenAPI spec validation

**Key Components:**
- `validateSpec()` - Main validation function
- `loadSpec()` - File loading (JSON/YAML)
- Error collection with suggestions
- Validation rules

**Data Structures:**
```typescript
interface ValidationError {
  path: string;           // JSON path (e.g., "info.version")
  message: string;        // Error description
  suggestion?: string;    // Actionable fix suggestion
}
```

**Validation Rules:**
1. OpenAPI version must be 3.x
2. `info` object required with `title` and `version`
3. `paths` object required with at least one path
4. Each operation must have `responses`

### 3. Testing Layer (`src/commands/test.ts`)

**Responsibility:** API endpoint testing

**Key Components:**
- `runTests()` - Main test orchestrator
- `testEndpoint()` - Individual endpoint tester
- `replacePlaceholders()` - Path parameter replacement
- `buildQueryParams()` - Query string builder
- `loadSpec()` - Spec file loader

**Data Structures:**
```typescript
interface TestOptions {
  export?: string;    // JSON export file path
  verbose?: boolean;  // Enable verbose logging
}

interface TestResult {
  method: string;              // HTTP method
  endpoint: string;            // API endpoint
  status: number | null;       // HTTP status code
  success: boolean;            // Pass/fail
  message: string;             // Result message
  duration?: number;           // Request timing (ms)
  timestamp?: string;          // ISO timestamp
  requestHeaders?: Record;     // Request headers (verbose)
  responseHeaders?: Record;    // Response headers (verbose)
}
```

**Test Flow:**
```
Load Spec → Parse Paths → For Each Endpoint:
  1. Replace path parameters ({id} → 1)
  2. Build query parameters
  3. Send HTTP request (Axios)
  4. Capture timing & headers
  5. Record result
→ Display Results → Export JSON → Exit Code
```

### 4. HTTP Methods Supported

| Method  | Body Support | Usage |
|---------|-------------|-------|
| GET     | No          | Retrieve resources |
| POST    | Yes         | Create resources |
| PUT     | Yes         | Update resources |
| PATCH   | Yes         | Partial updates |
| DELETE  | No          | Delete resources |
| HEAD    | No          | Headers only |
| OPTIONS | No          | CORS preflight |

## Data Flow

### Validate Command Flow

```
┌─────────────┐
│ User Input  │
│  spec.yaml  │
└──────┬──────┘
       │
       ▼
┌─────────────────┐
│  Read File      │
│  (fs.readFile)  │
└──────┬──────────┘
       │
       ▼
┌──────────────────┐
│  Parse YAML/JSON │
│  (js-yaml)       │
└──────┬───────────┘
       │
       ▼
┌──────────────────────┐
│  Validate Structure  │
│  - openapi version   │
│  - info object       │
│  - paths object      │
│  - operations        │
└──────┬───────────────┘
       │
       ├─(errors)─→ ┌────────────────────┐
       │            │ Display Errors +   │
       │            │ Suggestions        │
       │            │ Exit Code: 1       │
       │            └────────────────────┘
       │
       └─(valid)──→ ┌────────────────────┐
                    │ Display Success    │
                    │ Exit Code: 0       │
                    └────────────────────┘
```

### Test Command Flow

```
┌─────────────────────┐
│  User Input         │
│  spec.yaml + URL    │
│  + flags (-v, -e)   │
└──────┬──────────────┘
       │
       ▼
┌─────────────────────┐
│  Load & Parse Spec  │
└──────┬──────────────┘
       │
       ▼
┌────────────────────────┐
│  Iterate Paths         │
│  ┌──────────────────┐  │
│  │ For each method  │  │
│  └──────┬───────────┘  │
└─────────┼──────────────┘
          │
          ▼
    ┌─────────────────────┐
    │  Build Request      │
    │  - Replace {params} │
    │  - Add query params │
    │  - Add body         │
    └──────┬──────────────┘
           │
           ▼
    ┌────────────────────┐
    │  Start Timer       │
    │  Send HTTP Request │
    │  (Axios)           │
    └──────┬─────────────┘
           │
           ▼
    ┌────────────────────┐
    │  Record Duration   │
    │  Capture Headers   │
    │  (if verbose)      │
    └──────┬─────────────┘
           │
           ▼
    ┌────────────────────┐
    │  Store Result      │
    │  - Status code     │
    │  - Success/fail    │
    │  - Timing          │
    │  - Headers         │
    └──────┬─────────────┘
           │
           ▼
    ┌────────────────────┐
    │  Display Result    │
    │  ✓ or ✗            │
    │  (+ verbose info)  │
    └──────┬─────────────┘
           │
     ┌─────┴─────┐
     │  More?    │
     └─────┬─────┘
           │ No
           ▼
    ┌──────────────────────┐
    │  Display Summary     │
    │  - Pass/fail counts  │
    │  - Total tests       │
    └──────┬───────────────┘
           │
           ▼
    ┌──────────────────────┐
    │  Export JSON?        │
    │  (if -e flag)        │
    └──────┬───────────────┘
           │
           ▼
    ┌──────────────────────┐
    │  Exit Code           │
    │  0 = all passed      │
    │  1 = some failed     │
    └──────────────────────┘
```

## Design Patterns

### 1. Command Pattern
- Each command (validate, test) is encapsulated
- Commander.js handles routing
- Commands are independently testable

### 2. Strategy Pattern
- Different HTTP methods handled via switch statement
- Request body generation varies by content type
- Error handling adapts to error type

### 3. Builder Pattern
- URL building: base + path + query
- Query parameters built from spec definition
- Request configuration assembled incrementally

### 4. Observer Pattern (Results)
- Results collected in array
- Statistics calculated on-the-fly
- Final reporting after all tests

## Error Handling

### File Errors
```typescript
if (!fs.existsSync(filePath)) {
  throw new Error('File not found')
  // Suggestion: Check the file path
}
```

### Network Errors
```typescript
catch (error) {
  if (error.code === 'ECONNREFUSED') {
    message = 'Connection refused'
    // API server not running
  } else if (error.code === 'ETIMEDOUT') {
    message = 'Request timeout'
    // Default: 10 seconds
  }
}
```

### Validation Errors
```typescript
errors.push({
  path: 'openapi',
  message: 'Missing required field',
  suggestion: 'Add: openapi: "3.0.0"'
})
```

## Performance Considerations

### 1. Sequential Testing
- Endpoints tested one at a time
- Prevents overwhelming API server
- More predictable timing results

**Future Enhancement:** Parallel testing with concurrency limit

### 2. Timeout Management
- Default: 10 seconds per request
- Prevents hanging on slow APIs
- Configurable in code (not CLI yet)

### 3. Memory Usage
- Results stored in memory array
- Acceptable for typical APIs (<1000 endpoints)
- JSON export happens at end

## Security Considerations

### 1. File Access
- Only reads specified files
- No directory traversal
- File paths validated by Node.js

### 2. HTTP Requests
- Uses Axios with default security
- No authentication credentials stored
- HTTPS supported

### 3. Error Messages
- Sensitive data not logged
- File paths shown (user-controlled)
- No credential exposure

## Dependencies

### Production Dependencies
```json
{
  "axios": "1.6.0",        // HTTP client
  "commander": "12.0.0",   // CLI framework
  "js-yaml": "4.1.0"       // YAML parser
}
```

### Development Dependencies
```json
{
  "typescript": "5.9.2",   // Type safety
  "jest": "29.7.0",        // Unit testing
  "@types/*": "*"          // TypeScript types
}
```

## Extension Points

### Adding New Commands
```typescript
program
  .command('new-command')
  .description('Description')
  .action(async () => {
    // Implementation
  });
```

### Adding New Validators
```typescript
// In validate.ts
if (!condition) {
  errors.push({
    path: 'field.path',
    message: 'Error description',
    suggestion: 'How to fix'
  });
}
```

### Adding New HTTP Methods
```typescript
// In test.ts
case 'TRACE':
  response = await axios.trace(url, config);
  break;
```

## Testing Strategy

See [TESTING-GUIDE.md](TESTING-GUIDE.md) for comprehensive testing documentation.

## Future Enhancements

### Planned Features
1. **Parallel Testing** - Test multiple endpoints concurrently
2. **Custom Timeouts** - CLI flag for timeout configuration
3. **Request Bodies** - Generate from schema (not just examples)
4. **Authentication** - Support Bearer, API Key, Basic auth
5. **Response Validation** - Schema validation against spec
6. **HTML Reports** - In addition to JSON export
7. **Configuration File** - Store default options

### Architecture Changes Needed
- Parallel testing: Worker pool or Promise.all with limit
- Auth: Credential management system
- Schema validation: OpenAPI parser library (kin-openapi alternative)
- Config file: YAML config loader in home directory

## Comparison with TUI

| Aspect | CLI | TUI |
|--------|-----|-----|
| **Architecture** | Procedural | Event-driven (Elm) |
| **State** | Stateless | Stateful (model) |
| **Output** | Streaming | Screen updates |
| **Interaction** | One-shot | Interactive |
| **Use Case** | Automation | Development |

## Conclusion

The OpenAPI CLI follows a simple, linear architecture optimized for automation and scripting. Its stateless design makes it ideal for CI/CD pipelines, while its comprehensive error messages and JSON export enable integration with other tools.

For interactive development and debugging, consider using the companion **openapi-tui** tool.

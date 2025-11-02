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
  export?: string;           // JSON export file path
  exportHtml?: string;       // HTML export file path
  exportJunit?: string;      // JUnit XML export file path
  verbose?: boolean;         // Enable verbose logging
  quiet?: boolean;           // Quiet mode (only errors)
  timeout?: string;          // Request timeout in ms
  parallel?: string;         // Concurrency limit
  retry?: string;            // Max retry attempts
  validateSchema?: boolean;  // Validate response schemas
  watch?: boolean;           // Watch mode for file changes

  // Authentication
  authBearer?: string;       // Bearer token
  authApiKey?: string;       // API key
  authHeader?: string;       // API key header name
  authQuery?: string;        // API key query param name
  authBasic?: string;        // Basic auth (user:pass)

  // Filtering
  methods?: string;          // Filter by HTTP methods
  paths?: string;            // Filter by path pattern
  header?: string[];         // Custom headers
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
  schemaErrors?: string[];     // Schema validation errors
  retryCount?: number;         // Number of retries
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

### 1. Parallel Execution ✅
- Configurable concurrency with `--parallel <limit>` flag
- Default: 5 concurrent requests
- Promise-based concurrency control prevents overwhelming servers
- Performance gain: 8% faster (parallel 5), 23% faster (parallel 10) for 8 endpoints
- Scales well with larger APIs (50+ endpoints)

### 2. Timeout Management ✅
- Configurable via `--timeout <ms>` flag
- Default: 10000ms (10 seconds)
- Prevents hanging on slow APIs
- Works with retry logic for better reliability

### 3. Memory Usage
- Results stored in memory array
- Acceptable for typical APIs (<1000 endpoints)
- Export formats (JSON/HTML/JUnit) written to disk at completion
- Watch mode persists results between runs

### 4. Schema Validation Overhead
- AJV validation adds ~10-50ms per response
- Disabled by default, enable with `--validate-schema`
- Minimal impact for most use cases

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
  "axios": "^1.6.0",            // HTTP client
  "commander": "^12.0.0",       // CLI framework
  "js-yaml": "^4.1.0",          // YAML parser
  "ajv": "^8.12.0",             // JSON Schema validator
  "chokidar": "^4.0.3",         // File watcher (watch mode)
  "openapi-types": "^12.1.3"    // TypeScript types for OpenAPI
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

## Implemented Features (Phase 3 Complete)

### Core Features ✅
1. **Parallel Testing** ✅ - Promise-based concurrency with configurable limits (`--parallel`)
2. **Custom Timeouts** ✅ - CLI flag for timeout configuration (`--timeout`)
3. **Request Bodies** ✅ - Schema-based generation with fallback to examples
4. **Authentication** ✅ - Bearer, API Key (header/query), Basic auth
5. **Response Validation** ✅ - AJV-based schema validation (`--validate-schema`)
6. **HTML Reports** ✅ - Beautiful styled HTML export (`--export-html`)
7. **JUnit XML Export** ✅ - CI/CD integration format (`--export-junit`)
8. **Configuration File** ✅ - YAML/JSON with auto-discovery (`--config`)
9. **Retry Logic** ✅ - Exponential backoff for network errors (`--retry`)
10. **Watch Mode** ✅ - Auto re-run on file changes (`--watch`)
11. **Progress Indicator** ✅ - Test count and progress counter
12. **Method Filtering** ✅ - Test specific HTTP methods (`--methods`)
13. **Path Filtering** ✅ - Wildcard pattern matching (`--paths`)
14. **Quiet Mode** ✅ - Suppress output for CI/CD (`--quiet`)
15. **Custom Headers** ✅ - Repeatable header flag (`-H`)

### Architecture Implementation
- **Parallel testing**: Promise.all with concurrency limiter using async iteration
- **Auth**: Credential injection via Axios config (Authorization header, query params)
- **Schema validation**: AJV library with OpenAPI 3.x schema support
- **Config file**: YAML/JSON loader with directory tree search and precedence
- **Watch mode**: Chokidar file watching with persistent watcher
- **Retry logic**: Exponential backoff with network error detection

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

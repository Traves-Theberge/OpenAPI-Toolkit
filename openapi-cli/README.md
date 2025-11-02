# OpenAPI CLI Tester

A professional command-line tool for validating OpenAPI specifications and testing APIs against them.

## Features

- **Comprehensive Validation** - Validate OpenAPI 3.x specifications with detailed error reporting
- **Full HTTP Method Support** - Test GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS endpoints
- **Path Parameter Handling** - Automatically replaces `{id}` placeholders with test values
- **Query Parameter Support** - Builds query strings from OpenAPI parameter definitions
- **Colored Output** - Clear visual feedback with success (green) and error (red) indicators
- **Request Timeout Handling** - Configurable 10-second timeout with graceful error handling
- **JSON & YAML Support** - Parse and validate both `.json` and `.yaml`/`.yml` files
- **Detailed Error Messages** - Connection errors, timeouts, and HTTP status codes clearly reported
- **Summary Statistics** - Test results summary with pass/fail counts
- **JSON Export** - Export test results to JSON file for CI/CD integration
- **Verbose Mode** - Show detailed request/response headers and timing info
- **Enhanced Error Messages** - Actionable suggestions for fixing validation errors
- **Authentication Support** - Bearer tokens, API keys (header/query), and Basic authentication
- **Custom Timeouts** - Configurable request timeouts for slow APIs or fast-fail scenarios
- **Custom Headers** - Add custom HTTP headers to all requests (repeatable -H flag)
- **Method Filtering** - Test only specific HTTP methods (--methods GET,POST)
- **Path Filtering** - Test only paths matching a pattern with * wildcard (--paths /users/*)
- **Quiet Mode** - Suppress output except errors and exit codes for CI/CD (--quiet)
- **Parallel Execution** - Run tests concurrently with configurable concurrency limit (--parallel 5)
- **Schema-Based Body Generation** - Automatically generate request bodies from JSON Schema definitions
- **HTML Export** - Generate beautiful, styled HTML reports for test results (--export-html)
- **JUnit XML Export** - Generate JUnit XML reports for CI/CD integration (--export-junit)

## Installation

### Global Installation

```bash
npm install
npm run build
npm link
```

Now `openapi-test` will be available globally.

### Local Development

```bash
npm install
npm run build
```

## Usage

### Validate an OpenAPI Specification

Validates the structure and required fields of an OpenAPI spec:

```bash
openapi-test validate path/to/spec.yaml
```

**Output Example:**
```
📄 Validating OpenAPI specification: openapi.yaml
ℹ Found 2 paths with 3 operations
✓ Validation successful!
  OpenAPI Version: 3.0.3
  Title: My API
  Version: 1.0.0
```

### Test API Endpoints

Tests all endpoints defined in the OpenAPI spec against a live API:

```bash
# Basic testing
openapi-test test path/to/spec.yaml http://api.example.com

# Verbose mode (show request/response details)
openapi-test test path/to/spec.yaml http://api.example.com --verbose

# Export results to JSON
openapi-test test path/to/spec.yaml http://api.example.com --export results.json

# Custom timeout (30 seconds)
openapi-test test path/to/spec.yaml http://api.example.com --timeout 30000

# Combine flags
openapi-test test path/to/spec.yaml http://api.example.com -v -e results.json -t 30000

# Bearer token authentication
openapi-test test path/to/spec.yaml http://api.example.com --auth-bearer YOUR_TOKEN

# API key in header (default: X-API-Key)
openapi-test test path/to/spec.yaml http://api.example.com --auth-api-key YOUR_API_KEY

# API key in custom header
openapi-test test path/to/spec.yaml http://api.example.com --auth-api-key YOUR_KEY --auth-header X-Custom-API-Key

# API key in query parameter
openapi-test test path/to/spec.yaml http://api.example.com --auth-api-key YOUR_KEY --auth-query api_key

# Basic authentication
openapi-test test path/to/spec.yaml http://api.example.com --auth-basic username:password

# Custom headers (repeatable)
openapi-test test path/to/spec.yaml http://api.example.com -H "X-Custom-Header: Value" -H "X-Another: Value2"

# Filter by HTTP methods
openapi-test test path/to/spec.yaml http://api.example.com --methods GET,POST

# Quiet mode (only errors and exit code)
openapi-test test path/to/spec.yaml http://api.example.com --quiet

# Filter by path pattern (supports * wildcard)
openapi-test test path/to/spec.yaml http://api.example.com --paths "/users/*"
```

**Output Example:**
```
🧪 Testing API: JSONPlaceholder API
📍 Base URL: https://jsonplaceholder.typicode.com

✓ GET     /posts                                   - 200 OK
  Duration: 156ms
  Response Headers: {"content-type":"application/json; charset=utf-8"}
✓ POST    /posts                                   - 201 OK
  Duration: 142ms
✓ GET     /posts/1                                 - 200 OK
✗ DELETE  /posts/999                               - HTTP 404 Not Found

================================================================================
📊 Summary: 3 passed, 1 failed, 4 total
✓ Results exported to results.json
✗ Some tests failed
```

### Testing with Sample Server

Start the included sample server:

```bash
npm run server
```

Then in another terminal:

```bash
openapi-test test openapi.yaml http://localhost:3000
```

## Development

### Available Scripts

```bash
npm run build   # Compile TypeScript to JavaScript
npm run dev     # Run with ts-node (no build needed)
npm run server  # Start sample Express server
npm test        # Run Jest unit tests
npm run lint    # Lint TypeScript code
npm run format  # Format code with Prettier
```

### Running Tests

```bash
npm test
```

The test suite includes:
- Validation tests for valid/invalid OpenAPI specs
- File handling tests
- Mock filesystem tests

## Advanced Features

### Verbose Mode

Show detailed request/response information including headers and timing:

```bash
openapi-test test spec.yaml http://api.example.com --verbose
# or short form:
openapi-test test spec.yaml http://api.example.com -v
```

**Output includes:**
- Request duration in milliseconds
- Response headers
- Detailed timing for each endpoint test

### JSON Export

Export test results to a JSON file for integration with CI/CD pipelines:

```bash
openapi-test test spec.yaml http://api.example.com --export results.json
# or short form:
openapi-test test spec.yaml http://api.example.com -e results.json
```

**JSON Export Format:**
```json
{
  "timestamp": "2025-01-15T10:30:00.000Z",
  "specPath": "spec.yaml",
  "baseUrl": "http://api.example.com",
  "totalTests": 10,
  "passed": 8,
  "failed": 2,
  "results": [
    {
      "method": "GET",
      "endpoint": "/users",
      "status": 200,
      "success": true,
      "message": "OK",
      "duration": 156,
      "timestamp": "2025-01-15T10:30:00.123Z"
    }
  ]
}
```

### Custom Timeout

Configure request timeout for slow APIs or testing timeout behavior:

```bash
openapi-test test spec.yaml http://api.example.com --timeout 30000
# or short form:
openapi-test test spec.yaml http://api.example.com -t 30000
```

**Parameters:**
- Timeout value in milliseconds
- Default: 10000ms (10 seconds)
- Minimum: 1000ms (1 second) recommended
- Maximum: No hard limit (use responsibly)

**Use Cases:**
- **Slow APIs**: Increase timeout for APIs with long response times
- **Fast APIs**: Decrease timeout to fail fast
- **Timeout Testing**: Test how your monitoring handles timeouts
- **Network Issues**: Adjust based on network reliability

**Example:**
```bash
# Test with 60-second timeout for slow endpoint
openapi-test test slow-api.yaml https://slow-api.com -t 60000 -v

# Fast fail with 2-second timeout
openapi-test test fast-api.yaml https://fast-api.com -t 2000
```

### Authentication

Test protected APIs with multiple authentication methods:

#### Bearer Token Authentication

```bash
openapi-test test spec.yaml https://api.example.com --auth-bearer YOUR_TOKEN
```

Sends an `Authorization: Bearer YOUR_TOKEN` header with all requests.

**Use Cases:**
- OAuth 2.0 APIs
- JWT-based authentication
- Modern REST APIs

#### API Key Authentication

**In Header (default: X-API-Key):**
```bash
openapi-test test spec.yaml https://api.example.com --auth-api-key YOUR_KEY
```

**In Custom Header:**
```bash
openapi-test test spec.yaml https://api.example.com --auth-api-key YOUR_KEY --auth-header X-Custom-API-Key
```

**In Query Parameter:**
```bash
openapi-test test spec.yaml https://api.example.com --auth-api-key YOUR_KEY --auth-query api_key
```

Adds the API key to either:
- Request headers (default or custom header name)
- Query string parameter

**Use Cases:**
- API key-based services
- Third-party integrations
- Legacy APIs

#### Basic Authentication

```bash
openapi-test test spec.yaml https://api.example.com --auth-basic username:password
```

Sends an `Authorization: Basic <base64-encoded-credentials>` header.

**Use Cases:**
- HTTP Basic Auth APIs
- Simple authentication scenarios
- Internal APIs

**Example with Multiple Auth Options:**
```bash
# Test authenticated endpoints with verbose output and export
openapi-test test protected-api.yaml https://api.example.com \
  --auth-bearer $API_TOKEN \
  -v -e auth-test-results.json

# Test API with key in custom header
openapi-test test spec.yaml https://api.example.com \
  --auth-api-key $API_KEY \
  --auth-header X-API-Secret \
  -t 30000
```

**Security Notes:**
- Never commit credentials to version control
- Use environment variables for sensitive data
- Credentials are only sent to the specified base URL

### Custom Headers

Add custom HTTP headers to all requests:

```bash
# Single custom header
openapi-test test spec.yaml https://api.example.com -H "X-Request-ID: 12345"

# Multiple custom headers (repeatable flag)
openapi-test test spec.yaml https://api.example.com \
  -H "X-Request-ID: 12345" \
  -H "X-Client-Version: 1.0.0" \
  -H "X-Environment: staging"

# Override default headers
openapi-test test spec.yaml https://api.example.com -H "Content-Type: application/xml"
```

**Header Format:**
- Format: `"Name: Value"`
- Separate name and value with a colon
- Whitespace around colon is trimmed
- Case-sensitive header names

**Use Cases:**
- Request tracking (X-Request-ID, X-Correlation-ID)
- Client identification (User-Agent, X-Client-Version)
- Feature flags (X-Feature-Enabled)
- A/B testing (X-Variant)
- Content negotiation (Accept, Content-Type)

**Combined with Authentication:**
```bash
# Custom headers + Bearer token
openapi-test test spec.yaml https://api.example.com \
  --auth-bearer $TOKEN \
  -H "X-Request-ID: req-123" \
  -H "X-Client: CLI-Tester"
```

**Notes:**
- Custom headers are applied to all requests
- Later headers override earlier ones if the same name is used
- Authentication headers take precedence over custom headers

### Filter by HTTP Method

Test only specific HTTP methods to speed up targeted testing:

```bash
# Test only GET requests
openapi-test test spec.yaml https://api.example.com --methods GET

# Test only POST and PUT requests
openapi-test test spec.yaml https://api.example.com -m POST,PUT

# Case insensitive (get,post works too)
openapi-test test spec.yaml https://api.example.com --methods get,post
```

**Supported Methods:**
- GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS

**Use Cases:**
- **Read-only testing**: Filter to only GET requests (`--methods GET`)
- **Write operations**: Test POST, PUT, PATCH, DELETE (`--methods POST,PUT,PATCH,DELETE`)
- **Quick smoke tests**: Test only critical endpoints
- **Debugging**: Isolate specific method types
- **CI/CD stages**: Different test stages for different methods

**Format:**
- Comma-separated list of HTTP methods
- Case-insensitive (GET, get, Get all work)
- No spaces around commas (or quote the value)
- Unknown methods are silently ignored

**Examples:**
```bash
# Test only safe methods
openapi-test test spec.yaml https://api.example.com --methods GET,HEAD,OPTIONS

# Test data modification endpoints
openapi-test test spec.yaml https://api.example.com --methods POST,PUT,DELETE

# Combined with other flags
openapi-test test spec.yaml https://api.example.com \
  --methods GET \
  --auth-bearer $TOKEN \
  -v -e get-results.json
```

### Quiet Mode

Suppress all output except errors and the final exit code - ideal for CI/CD:

```bash
# Quiet mode - only errors shown
openapi-test test spec.yaml https://api.example.com --quiet

# Short form
openapi-test test spec.yaml https://api.example.com -q

# With export (errors + export shown)
openapi-test test spec.yaml https://api.example.com -q -e results.json
```

**Behavior:**
- **Hides**: Test progress, success messages, summaries
- **Shows**: Error messages, export failures
- **Exit codes**: 0 for success, 1 for failures (same as normal mode)

**Use Cases:**
- **CI/CD pipelines**: Clean logs, only failures visible
- **Cron jobs**: Quiet unless there's a problem
- **Scripting**: Check exit code without noise
- **Automated testing**: Parse JSON export instead of console output

**What You See:**

Success (quiet):
```bash
$ openapi-test test spec.yaml https://api.example.com -q
All tests passed.
$ echo $?
0
```

Failure (quiet):
```bash
$ openapi-test test spec.yaml https://api.example.com -q
✗ GET     /invalid                                 - HTTP 404 Not Found
$ echo $?
1
```

**CI/CD Example:**
```bash
#!/bin/bash
if openapi-test test spec.yaml https://staging.api.com -q -e results.json; then
  echo "API tests passed"
else
  echo "API tests failed - check results.json"
  exit 1
fi
```

### Filter by Path Pattern

Test only paths matching a specific pattern with wildcard support:

```bash
# Test exact path
openapi-test test spec.yaml https://api.example.com --paths "/users"

# Test all paths starting with /users/
openapi-test test spec.yaml https://api.example.com -p "/users/*"

# Test all /posts paths (including /posts/123)
openapi-test test spec.yaml https://api.example.com --paths "/posts*"
```

**Pattern Syntax:**
- Exact match: `/users` matches only `/users`
- Wildcard: `*` matches any characters
- `/users/*` matches `/users/123`, `/users/abc`, etc.
- `/posts*` matches `/posts`, `/posts/1`, `/posts/comments`

**Use Cases:**
- **Endpoint isolation**: Test specific API sections
- **Incremental testing**: Test one resource at a time
- **Debugging**: Focus on problematic endpoints
- **Staged rollouts**: Test new paths before release
- **Resource-based CI/CD**: Different pipelines for different resources

**Examples:**
```bash
# Test only user endpoints
openapi-test test spec.yaml https://api.example.com --paths "/users*"

# Test specific nested path
openapi-test test spec.yaml https://api.example.com --paths "/api/v1/users/*"

# Combine with method filter
openapi-test test spec.yaml https://api.example.com \
  --paths "/admin/*" \
  --methods GET \
  -q
```

**Notes:**
- Pattern matching uses regex internally
- Special regex characters (except *) are automatically escaped
- Case-sensitive matching
- If no paths match, zero tests will run

### Parallel Execution

Run tests concurrently to significantly improve performance for large APIs:

```bash
# Default parallel execution (concurrency limit: 5)
openapi-test test spec.yaml https://api.example.com --parallel 5

# Higher concurrency for fast APIs
openapi-test test spec.yaml https://api.example.com --parallel 10

# Lower concurrency for rate-limited APIs
openapi-test test spec.yaml https://api.example.com --parallel 2

# Sequential execution (default behavior if --parallel not specified)
openapi-test test spec.yaml https://api.example.com
```

**How It Works:**
- Tests are executed concurrently up to the specified limit
- New tests start as soon as a slot becomes available
- Results are displayed as they complete (order may vary)
- All other features work with parallel execution (filters, auth, export, etc.)

**Performance Impact:**
- **Sequential (no --parallel)**: ~478ms for 8 endpoints
- **Parallel --parallel 5**: ~439ms for 8 endpoints (8% faster)
- **Parallel --parallel 10**: ~366ms for 8 endpoints (23% faster)
- **Speedup increases with more endpoints and higher latency**

**Use Cases:**
- **Large APIs**: Test 100+ endpoints faster
- **CI/CD optimization**: Reduce pipeline execution time
- **Development**: Quick feedback during API development
- **Performance testing**: Simulate concurrent load

**Choosing Concurrency Limit:**
- **5 (default)**: Good balance for most APIs
- **10-20**: Fast, reliable APIs without rate limiting
- **2-3**: Rate-limited APIs or slow backends
- **1**: Sequential execution (same as no --parallel flag)

**Examples:**
```bash
# Fast test run with all features
openapi-test test spec.yaml https://api.example.com \
  --parallel 10 \
  -m GET,POST \
  --auth-bearer "$TOKEN" \
  -q -e results.json

# Rate-limited API with authentication
openapi-test test spec.yaml https://api.example.com \
  --parallel 2 \
  --auth-api-key "$API_KEY" \
  --timeout 30000

# Maximum performance for internal APIs
openapi-test test spec.yaml https://internal.api.com \
  --parallel 20 \
  -v
```

**Notes:**
- Parallel execution uses Promise-based concurrency control
- Results order may differ between runs
- All tests complete before summary is shown
- Exit codes work the same as sequential mode
- Export includes all results regardless of execution order

### Schema-Based Request Body Generation

Automatically generate request bodies from JSON Schema definitions when examples are not provided:

**How It Works:**
- If an `example` is provided in the spec, it's used (backward compatible)
- If no example exists, the CLI generates a body from the `schema` definition
- Supports all JSON Schema types: string, number, integer, boolean, object, array, null
- Honors schema constraints: min/max length, min/max values, enums, required fields
- Handles special string formats: email, uri, date, date-time, uuid

**Supported Schema Features:**
```yaml
requestBody:
  content:
    application/json:
      schema:
        type: object
        required:
          - name
          - email
        properties:
          name:
            type: string
            minLength: 3
            maxLength: 50
          email:
            type: string
            format: email        # → test@example.com
          age:
            type: integer
            minimum: 18
            maximum: 100         # → 59 (midpoint)
          active:
            type: boolean        # → true
          role:
            type: string
            enum: [admin, user]  # → admin (first value)
          tags:
            type: array
            items:
              type: string       # → ["testx", "testx"]
          address:
            type: object
            properties:
              city:
                type: string     # → "testx"
```

**Generated Values:**
- **String**: `"testx"` (respects minLength/maxLength)
- **Email**: `test@example.com`
- **URI/URL**: `https://example.com`
- **Date**: `2025-01-01`
- **DateTime**: `2025-01-01T00:00:00Z`
- **UUID**: `123e4567-e89b-12d3-a456-426614174000`
- **Number/Integer**: Midpoint between min and max (default: 50)
- **Boolean**: `true`
- **Enum**: First value in the enum array
- **Array**: 2 generated items
- **Object**: All required fields + 50% of optional fields

**Example Spec:**
```yaml
paths:
  /users:
    post:
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required:
                - name
                - email
              properties:
                name:
                  type: string
                email:
                  type: string
                  format: email
                age:
                  type: integer
                  minimum: 18
```

**Generated Body:**
```json
{
  "name": "testx",
  "email": "test@example.com",
  "age": 59
}
```

**Use Cases:**
- **Specs without examples**: Test APIs that only define schemas
- **Rapid testing**: No need to manually craft request bodies
- **Schema validation**: Verify your API accepts valid schema-compliant data
- **CI/CD**: Automated testing without maintaining example data

**Notes:**
- Examples always take precedence over schema generation
- Optional fields have a 50% chance of being included
- Arrays generate 2 items by default
- Supports `oneOf`, `anyOf`, `allOf` schema composition
- Works with all other CLI features (parallel, filters, auth, etc.)

### HTML Export

Generate beautiful, styled HTML reports for test results - perfect for sharing with teams or archiving:

**Basic Usage:**
```bash
# Generate HTML report
openapi-test test spec.yaml https://api.example.com --export-html report.html

# Combine with JSON export
openapi-test test spec.yaml https://api.example.com \
  --export results.json \
  --export-html report.html
```

**Features:**
- 📊 **Summary Cards**: Total tests, passed, failed, success rate
- 🎨 **Color-Coded Results**: Green for success, red for failures
- 📱 **Responsive Design**: Works on desktop, tablet, and mobile
- 🖨️ **Print-Friendly**: Optimized for PDF export and printing
- ⚡ **Self-Contained**: All CSS embedded, no external dependencies
- 🔍 **Detailed Table**: Method, endpoint, status, message, duration for each test

**HTML Report Includes:**
- API title and metadata (base URL, spec path, timestamp)
- Visual summary with pass/fail statistics
- Color-coded HTTP methods (GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS)
- Success rate percentage
- Sortable results table
- Duration in milliseconds for each request
- Professional gradient header design

**Sample Report Structure:**
```
┌─────────────────────────────────────────────┐
│ 🧪 API Test Results                         │
│ Integration Test API                        │
│ https://api.example.com                     │
├─────────────────────────────────────────────┤
│ [Total: 10] [Passed: 8] [Failed: 2] [80%]  │
├─────────────────────────────────────────────┤
│ ✓ GET  /users      200 OK        156ms     │
│ ✓ POST /users      201 OK        234ms     │
│ ✗ GET  /invalid    404 Not Found 45ms      │
└─────────────────────────────────────────────┘
```

**Use Cases:**
- **Team Sharing**: Email HTML reports to stakeholders
- **CI/CD Artifacts**: Archive test results as build artifacts
- **Documentation**: Embed in wikis or documentation sites
- **Debugging**: Visual overview of API health
- **Historical Records**: Track API behavior over time

**Example with All Features:**
```bash
# Complete test run with HTML report
openapi-test test spec.yaml https://api.example.com \
  --parallel 10 \
  -m GET,POST \
  --auth-bearer "$TOKEN" \
  -v \
  --export-html detailed-report.html
```

**Notes:**
- HTML file is self-contained (no external CSS/JS)
- Compatible with all browsers (Chrome, Firefox, Safari, Edge)
- Can be opened directly or served via web server
- Embedded CSS uses modern flexbox/grid layouts
- Print media queries optimize for PDF export

### JUnit XML Export

Generate JUnit XML format test reports for seamless CI/CD integration with Jenkins, GitLab CI, GitHub Actions, and other tools:

**Basic Usage:**
```bash
# Generate JUnit XML report
openapi-test test spec.yaml https://api.example.com --export-junit results.xml

# Combine with other export formats
openapi-test test spec.yaml https://api.example.com \
  --export results.json \
  --export-html report.html \
  --export-junit junit.xml
```

**XML Format:**
- Standard JUnit XML schema compliant
- Compatible with all major CI/CD platforms
- Includes test duration, failure messages, and metadata
- Groups tests by API title as test suite
- Classname format: `{API Title}.{HTTP Method}`

**Sample XML Structure:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuites name="My API" tests="10" failures="2" errors="0" time="1.234">
  <testsuite name="My API" tests="10" failures="2" errors="0" time="1.234">
    <properties>
      <property name="spec" value="spec.yaml"/>
      <property name="baseUrl" value="https://api.example.com"/>
    </properties>
    <testcase name="GET /users" classname="My API.GET" time="0.156"/>
    <testcase name="POST /users" classname="My API.POST" time="0.234">
      <failure message="HTTP 404 Not Found" type="HTTP 404">
        Test: POST /users
        Status: 404
        Message: HTTP 404 Not Found
        Endpoint: https://api.example.com/users
      </failure>
    </testcase>
  </testsuite>
</testsuites>
```

**CI/CD Integration Examples:**

**Jenkins:**
```groovy
pipeline {
  stages {
    stage('API Tests') {
      steps {
        sh 'openapi-test test spec.yaml https://api.example.com --export-junit results.xml'
      }
      post {
        always {
          junit 'results.xml'
        }
      }
    }
  }
}
```

**GitHub Actions:**
```yaml
- name: Run API Tests
  run: openapi-test test spec.yaml https://api.example.com --export-junit results.xml

- name: Publish Test Results
  uses: EnricoMi/publish-unit-test-result-action@v2
  if: always()
  with:
    files: results.xml
```

**GitLab CI:**
```yaml
api-tests:
  script:
    - openapi-test test spec.yaml https://api.example.com --export-junit results.xml
  artifacts:
    reports:
      junit: results.xml
```

**Use Cases:**
- **CI/CD Integration**: Automatic test reporting in pipelines
- **Trend Analysis**: Track API reliability over time
- **Failure Tracking**: Detailed failure information for debugging
- **Team Dashboards**: Visual test results in CI tools
- **Quality Gates**: Block deployments on test failures

**Features:**
- Test duration in seconds (3 decimal places)
- Failure type and message for failed tests
- System output with test summary
- Properties for spec path, base URL, timestamp
- Hostname extracted from base URL
- XML character escaping for special characters

**Notes:**
- Time values are in seconds (JUnit standard)
- Failed tests include `<failure>` element with details
- All tests have unique names: `{METHOD} {ENDPOINT}`
- Compatible with JUnit 4 and 5 report parsers
- Works with all CLI features (filters, auth, parallel, etc.)

### Enhanced Error Messages

Validation errors now include actionable suggestions:

```bash
openapi-test validate invalid-spec.yaml
```

**Output:**
```
✗ Validation failed with 2 error(s):

  1. openapi: Missing required field "openapi"
     💡 Add: openapi: "3.0.0" or openapi: "3.1.0" at the root level

  2. info.version: Missing required field "info.version"
     💡 Add: version: "1.0.0" under the info object
```

### Path Parameter Replacement

The CLI automatically replaces path parameters with test values:

```yaml
paths:
  /users/{id}:
    get: ...
```

Becomes: `/users/1` when tested

### Query Parameter Generation

Query parameters from the OpenAPI spec are automatically added:

```yaml
parameters:
  - name: userId
    in: query
    schema:
      type: integer
```

Generates: `?userId=1` in the request

### Request Body Handling

For POST, PUT, and PATCH requests, the CLI uses example request bodies if defined:

```yaml
requestBody:
  content:
    application/json:
      example:
        title: "Test Post"
        body: "Test content"
```

### Error Handling

The CLI gracefully handles various error conditions:
- **ECONNREFUSED** - Connection refused
- **ETIMEDOUT** - Request timeout (10s default)
- **HTTP Errors** - Non-2xx status codes with detailed messages
- **Parse Errors** - Invalid JSON/YAML with specific error messages

## Examples

### Example 1: Validate a Spec

```bash
$ openapi-test validate examples/petstore.yaml

📄 Validating OpenAPI specification: examples/petstore.yaml
ℹ Found 14 paths with 20 operations
✓ Validation successful!
  OpenAPI Version: 3.0.0
  Title: Petstore API
  Version: 1.0.0
```

### Example 2: Test JSONPlaceholder API

```bash
$ openapi-test test ../openapi-tui/tests/jsonplaceholder-spec.yaml https://jsonplaceholder.typicode.com

🧪 Testing API: JSONPlaceholder API
📍 Base URL: https://jsonplaceholder.typicode.com

✓ GET     /posts                                   - 200 OK
✓ GET     /posts/1                                 - 200 OK
✓ POST    /posts                                   - 201 OK
✓ PUT     /posts/1                                 - 200 OK
✓ PATCH   /posts/1                                 - 200 OK
✓ DELETE  /posts/1                                 - 200 OK
✓ GET     /users                                   - 200 OK
✓ GET     /users/1                                 - 200 OK

================================================================================
📊 Summary: 8 passed, 0 failed, 8 total
✓ All tests passed!
```

## Technical Details

### Supported HTTP Methods

- GET
- POST
- PUT
- PATCH
- DELETE
- HEAD
- OPTIONS

(TRACE is defined in spec but not commonly supported by Axios)

### OpenAPI Version Support

- OpenAPI 3.0.x
- OpenAPI 3.1.x

### Dependencies

- **commander** - CLI framework
- **axios** - HTTP client for testing
- **js-yaml** - YAML parsing
- **TypeScript** - Type safety

## Troubleshooting

### "File not found" error

Ensure the path to your OpenAPI spec is correct (absolute or relative path).

### "Connection refused" error

Check that:
1. The API server is running
2. The base URL is correct
3. There are no firewall/network issues

### "Request timeout" error

The default timeout is 10 seconds. If your API is slow, endpoints may timeout. This is expected behavior for performance monitoring.

## License

MIT

## Related Projects

- **[openapi-tui](../openapi-tui/)** - Interactive terminal UI version with Bubble Tea
- Built with the same core functionality, different interfaces
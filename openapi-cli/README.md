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
- **Configuration File Support** - Store default options in YAML or JSON config files (--config)
- **Response Schema Validation** - Validate API responses against OpenAPI schemas (--validate-schema)
- **Retry Logic** - Automatically retry failed requests with exponential backoff (--retry)
- **Watch Mode** - Automatically re-run tests when spec file changes (--watch)
- **Progress Indicator** - Shows test progress during execution

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

## 📚 Documentation

For comprehensive guides and documentation:

- **[Complete User Guide](docs/USER-GUIDE.md)** - End-to-end guide with real-world examples
- **[Architecture Guide](docs/ARCHITECTURE.md)** - System design and components
- **[Progress Tracking](docs/PROGRESS.md)** - Feature roadmap and status
- **[Testing Guide](docs/TESTING-GUIDE.md)** - Testing procedures

---

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

### Configuration File Support

Store default options in a configuration file instead of passing them via command-line every time. Supports both YAML and JSON formats.

**Automatic Discovery:**

The CLI automatically searches for config files in the current directory and parent directories:
- `.openapi-cli.yaml`
- `.openapi-cli.yml`
- `.openapi-cli.json`
- `openapi-cli.yaml`
- `openapi-cli.yml`
- `openapi-cli.json`

**Explicit Config File:**
```bash
# Use specific config file
openapi-test test spec.yaml https://api.example.com --config my-config.yaml
```

**Example YAML Config:**
```yaml
# .openapi-cli.yaml
# Authentication
auth-bearer: "your-jwt-token-here"

# Custom headers
headers:
  - "User-Agent: My-App/1.0"
  - "Accept: application/json"

# Request options
timeout: 15000
verbose: true
quiet: false

# Filtering
methods: "GET,POST,PUT"
paths: "/api/v1/*"

# Execution
parallel: 10

# Export
export: "results.json"
export-html: "report.html"
export-junit: "junit.xml"
```

**Example JSON Config:**
```json
{
  "auth-bearer": "your-jwt-token",
  "headers": [
    "User-Agent: My-App/1.0",
    "Accept: application/json"
  ],
  "timeout": 15000,
  "verbose": true,
  "parallel": 10,
  "export-html": "report.html"
}
```

**Supported Options:**

All command-line options can be configured in the config file:
- Authentication: `auth-bearer`, `auth-api-key`, `auth-header`, `auth-query`, `auth-basic`
- Headers: `headers` (array of strings)
- Request: `timeout`, `verbose`, `quiet`
- Filtering: `methods`, `paths`
- Execution: `parallel`
- Export: `export`, `export-html`, `export-junit`

**Option Precedence:**

CLI options always take precedence over config file options:

```bash
# Config file has verbose: false
# This command will use verbose: true
openapi-test test spec.yaml https://api.example.com --verbose
```

**Header Merging:**

Headers from both config file and CLI are merged:

```yaml
# Config file
headers:
  - "User-Agent: My-App"
```

```bash
# Both headers will be sent
openapi-test test spec.yaml https://api.example.com -H "X-Custom: value"
```

**Use Cases:**
- **Project Defaults**: Store project-specific settings in `.openapi-cli.yaml`
- **Environment Configs**: Different configs for dev/staging/prod
- **Team Sharing**: Commit config to git for consistent team settings
- **CI/CD**: Separate configs for local dev vs CI/CD pipeline
- **Multiple APIs**: Different config files for different API endpoints

**Notes:**
- Config files are optional - all options can still be passed via CLI
- YAML format supports comments for documentation
- JSON format is more strict but easier to generate programmatically
- Config files are searched up the directory tree (like .gitignore)
- Use `.openapi-cli.yaml` prefix to hide from directory listings

### Response Schema Validation

Validate API response bodies against the schemas defined in your OpenAPI specification to ensure your API returns correctly structured data.

**Basic Usage:**
```bash
# Enable schema validation
openapi-test test spec.yaml https://api.example.com --validate-schema

# Combine with other options
openapi-test test spec.yaml https://api.example.com \
  --validate-schema \
  --verbose \
  --export-junit results.xml
```

**How It Works:**

The CLI validates responses against the schema defined in the OpenAPI spec:

```yaml
paths:
  /users/{id}:
    get:
      responses:
        '200':
          description: User object
          content:
            application/json:
              schema:
                type: object
                required:
                  - id
                  - name
                  - email
                properties:
                  id:
                    type: integer
                  name:
                    type: string
                  email:
                    type: string
                    format: email
```

**Validation Features:**

- **Type Checking**: Validates data types (string, number, integer, boolean, object, array)
- **Required Fields**: Ensures required properties are present
- **Constraints**: Validates min/max values, min/max length, patterns
- **Enums**: Checks values are in allowed list
- **Formats**: Validates email, uri, date, date-time, uuid formats
- **Additional Properties**: Detects unexpected properties
- **Nested Objects**: Recursively validates nested structures
- **Arrays**: Validates array items against schemas

**Error Reporting:**

When validation fails, detailed error messages are shown:

```bash
✗ GET     /users/1                                 - Schema validation failed: 3 error(s)
  ⚠  /name: expected type string, got number
  ⚠  /email: must be valid email format
  ⚠  root: missing required property 'age'
```

**Error Types:**

- `expected type X, got Y` - Type mismatch
- `missing required property 'X'` - Missing required field
- `value must be one of [...]` - Enum violation
- `must match pattern X` - Regex pattern mismatch
- `must be valid X format` - Format validation failed (email, uri, etc.)
- `additional property 'X' not allowed` - Unexpected property
- `must be <= X` / `must be >= X` - Min/max constraint violation

**Use Cases:**

- **API Contract Testing**: Ensure backend follows OpenAPI contract
- **Regression Testing**: Detect schema changes in API responses
- **Integration Testing**: Validate third-party API responses
- **Development**: Catch schema mismatches early in development
- **CI/CD**: Automate API contract verification in pipelines

**Configuration:**

Enable in config file:

```yaml
# .openapi-cli.yaml
validate-schema: true
verbose: true
export-junit: "schema-validation-results.xml"
```

**Notes:**

- Validation only runs on successful HTTP responses (2xx status codes)
- Uses JSON Schema validation via AJV library
- Follows OpenAPI 3.x schema specification
- Supports $ref references within schemas
- Partial validation: tests continue even if some fail
- Schema errors count as test failures in summary

**Performance:**

Schema validation adds minimal overhead (~10-50ms per response depending on complexity). Use sparingly for large test suites or disable for performance testing.

### Retry Logic

Automatically retry failed requests with exponential backoff to handle transient network errors and improve test reliability.

**Basic Usage:**
```bash
# Retry failed requests up to 3 times
openapi-test test spec.yaml https://api.example.com --retry 3

# Short form
openapi-test test spec.yaml https://api.example.com -r 3

# Combine with verbose to see retry attempts
openapi-test test spec.yaml https://api.example.com --retry 3 --verbose
```

**How It Works:**

The CLI automatically retries requests that fail with network errors using exponential backoff:

1. **First retry**: Wait 1 second
2. **Second retry**: Wait 2 seconds
3. **Third retry**: Wait 4 seconds
4. **Fourth retry**: Wait 8 seconds
5. **Maximum**: 10 seconds between retries

**Retryable Errors:**

Only network and connection errors are retried (not HTTP errors):

- `ECONNREFUSED` - Connection refused
- `ETIMEDOUT` - Request timeout
- `ENOTFOUND` - DNS lookup failed
- `ECONNRESET` - Connection reset
- `ENETUNREACH` - Network unreachable
- Network timeout errors

**Non-Retryable:**

HTTP errors (4xx, 5xx) are **NOT** retried:
- `401 Unauthorized` - Authentication issue, won't resolve with retry
- `404 Not Found` - Resource doesn't exist
- `500 Internal Server Error` - Server-side issue
- Other 4xx/5xx errors

**Retry Output:**

When retries occur, you'll see progress indicators:

```bash
🧪 Testing API: My API
📍 Base URL: https://api.example.com

  ↻  Retry attempt 1/3 after 1000ms...
  ↻  Retry attempt 2/3 after 2000ms...
✓ GET     /users                                   - 200 OK
```

**Configuration:**

Enable in config file:

```yaml
# .openapi-cli.yaml
retry: 3
verbose: true
timeout: 15000
```

**Use Cases:**

- **Unstable Networks**: Retry on flaky connections
- **Rate Limiting**: Handle temporary rate limit errors
- **Service Warmup**: Retry when services are starting up
- **CI/CD**: Improve test reliability in CI/CD pipelines
- **Development**: Handle local dev server restarts

**Performance Impact:**

Retries add time when failures occur:
- **No failures**: No impact (0ms overhead)
- **1 retry**: Adds 1 second
- **2 retries**: Adds 3 seconds total (1s + 2s)
- **3 retries**: Adds 7 seconds total (1s + 2s + 4s)

**Best Practices:**

1. **Use sparingly**: Only enable when needed (flaky networks, dev environments)
2. **Combine with timeout**: Set appropriate `--timeout` to fail fast
3. **Verbose mode**: Use `--verbose` to see retry attempts
4. **CI/CD**: Consider `--retry 2` for CI/CD to handle transient failures
5. **Don't over-retry**: More than 3 retries is usually excessive

**Notes:**

- Retry count of 0 (default) disables retries
- Exponential backoff prevents overwhelming failed services
- Maximum backoff capped at 10 seconds
- Retries are per-request, not per-test
- Works with parallel execution (`--parallel`)
- Compatible with all other options

### Watch Mode

Automatically re-run tests when the OpenAPI spec file changes, perfect for development workflows.

**Basic Usage:**
```bash
# Watch mode
openapi-test test spec.yaml https://api.example.com --watch

# Short form
openapi-test test spec.yaml https://api.example.com -w

# Combine with other options
openapi-test test spec.yaml https://api.example.com \
  --watch \
  --validate-schema \
  --verbose
```

**How It Works:**

The CLI monitors the spec file for changes and automatically re-runs tests when changes are detected:

1. Runs tests initially
2. Watches spec file for changes
3. Re-runs tests on any file change
4. Continues watching until you press Ctrl+C
5. Shows change notification before re-running

**Output:**
```bash
👁  Watching /path/to/spec.yaml for changes...

Press Ctrl+C to stop

🧪 Testing API: My API
📍 Base URL: https://api.example.com

ℹ Running 10 tests...

✓ GET     /users                                   - 200 OK
...

🔄 File changed, re-running tests...

🧪 Testing API: My API
...
```

**Use Cases:**

- **Development**: Test API changes in real-time during development
- **TDD Workflow**: Test-driven development with immediate feedback
- **Spec Editing**: Validate spec changes as you write them
- **Debugging**: Quickly iterate on API fixes
- **Documentation**: Keep tests running while updating OpenAPI docs

**Notes:**

- Watch mode runs indefinitely until stopped (Ctrl+C)
- File changes are debounced (no duplicate runs)
- Works with all other options (filters, auth, validation, etc.)
- Shows clear indicators for file changes
- Exits gracefully on Ctrl+C
- Not recommended for CI/CD (use normal mode)

**Best Practices:**

1. **Local Development Only**: Don't use in production or CI/CD
2. **Combine with Verbose**: Use `--verbose` to see detailed changes
3. **Filter Tests**: Use `--methods` or `--paths` to run subset of tests
4. **Fast Iteration**: Great for rapid development cycles
5. **Clear Terminal**: Watch mode output can be long, clear terminal periodically

### Progress Indicator

Shows test progress during execution for better visibility with large test suites.

**Features:**

**1. Test Count Display:**
Shows total number of tests before execution:
```bash
ℹ Running 25 tests...
```

**2. Progress Counter (Sequential Mode):**
Shows current test number during sequential execution:
```bash
[1/25] ✓ GET     /users                               - 200 OK
[2/25] ✓ POST    /users                               - 201 OK
[3/25] ✓ GET     /posts                               - 200 OK
...
```

**When Progress Counter Appears:**
- Sequential mode (`--parallel 1` or default)
- Only for test suites with more than 3 tests
- Disabled in quiet mode (`--quiet`)

**Output Modes:**

**Normal Mode (< 3 tests):**
```bash
ℹ Running 2 tests...

✓ GET     /users                                   - 200 OK
✓ POST    /users                                   - 201 OK
```

**Sequential Mode (> 3 tests):**
```bash
ℹ Running 10 tests...

[1/10] ✓ GET     /users                               - 200 OK
[2/10] ✓ POST    /users                               - 201 OK
[3/10] ✓ GET     /posts                               - 200 OK
...
```

**Parallel Mode:**
```bash
ℹ Running 50 tests...

✓ GET     /users                                   - 200 OK
✓ GET     /posts                                   - 200 OK
✓ GET     /comments                                - 200 OK
...
```

**Benefits:**

- **Visibility**: See how many tests are running
- **Progress Tracking**: Know how far along the test run is
- **Large Suites**: Especially useful for 50+ endpoint specs
- **ETA Estimation**: Mentally estimate completion time
- **Better UX**: More professional appearance

**Notes:**

- Progress counter only shows in sequential mode
- Parallel mode shows count but not individual progress
- Quiet mode disables all progress indicators
- No performance overhead

### Error Handling

The CLI gracefully handles various error conditions:
- **ECONNREFUSED** - Connection refused
- **ETIMEDOUT** - Request timeout (10s default)
- **HTTP Errors** - Non-2xx status codes with detailed messages
- **Parse Errors** - Invalid JSON/YAML with specific error messages
- **Config Errors** - Invalid or missing config files with helpful messages

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
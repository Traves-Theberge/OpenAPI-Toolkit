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
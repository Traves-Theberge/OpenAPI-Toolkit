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
openapi-test test path/to/spec.yaml http://api.example.com
```

**Output Example:**
```
🧪 Testing API: JSONPlaceholder API
📍 Base URL: https://jsonplaceholder.typicode.com

✓ GET     /posts                                   - 200 OK
✓ POST    /posts                                   - 201 OK
✓ GET     /posts/1                                 - 200 OK
✗ DELETE  /posts/999                               - HTTP 404 Not Found

================================================================================
📊 Summary: 3 passed, 1 failed, 4 total
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
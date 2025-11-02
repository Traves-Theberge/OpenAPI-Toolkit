# OpenAPI CLI - Complete User Guide

A comprehensive end-to-end guide for using the OpenAPI CLI tool.

## Table of Contents

- [Getting Started](#getting-started)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Basic Usage](#basic-usage)
- [Advanced Features](#advanced-features)
- [Configuration](#configuration)
- [CI/CD Integration](#cicd-integration)
- [Real-World Examples](#real-world-examples)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)

---

## Getting Started

### What is OpenAPI CLI?

OpenAPI CLI is a command-line tool that:
- ✅ Validates OpenAPI 3.x specifications
- ✅ Tests API endpoints automatically
- ✅ Generates request bodies from schemas
- ✅ Validates API responses against schemas
- ✅ Exports results in multiple formats (JSON, HTML, JUnit XML)
- ✅ Supports authentication (Bearer, API Key, Basic)
- ✅ Runs tests in parallel for better performance
- ✅ Integrates seamlessly with CI/CD pipelines

### When to Use CLI

**Perfect for:**
- ✅ CI/CD pipelines (GitHub Actions, GitLab CI, Jenkins)
- ✅ Automated API testing
- ✅ Pre-commit hooks
- ✅ API monitoring and health checks
- ✅ Regression testing
- ✅ Contract testing

**Not ideal for:**
- ❌ Interactive API exploration (use TUI instead)
- ❌ Manual debugging sessions (use TUI instead)
- ❌ Visual feedback requirements (use TUI instead)

---

## Installation

### Prerequisites

- **Node.js**: Version 16 or higher
- **npm**: Comes with Node.js

### Install from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/OpenAPI-Toolkit.git
cd OpenAPI-Toolkit/openapi-cli

# Install dependencies
npm install

# Build TypeScript to JavaScript
npm run build

# Link globally (makes 'openapi-test' available everywhere)
npm link

# Verify installation
openapi-test --version
```

### Verify Installation

```bash
# Should display version number
openapi-test --version

# Should display help information
openapi-test --help
```

---

## Quick Start

### 1. Validate Your First Spec

```bash
# Create a simple OpenAPI spec
cat > api-spec.yaml <<EOF
openapi: 3.0.0
info:
  title: My API
  version: 1.0.0
paths:
  /users:
    get:
      summary: Get all users
      responses:
        '200':
          description: Success
EOF

# Validate it
openapi-test validate api-spec.yaml
```

**Expected Output:**
```
📄 Validating OpenAPI specification: api-spec.yaml
ℹ Found 1 paths with 1 operations
✓ Validation successful!
  OpenAPI Version: 3.0.0
  Title: My API
  Version: 1.0.0
```

### 2. Test Your First API

```bash
# Test a public API (JSONPlaceholder)
openapi-test test \
  https://raw.githubusercontent.com/typicode/jsonplaceholder/master/openapi.yaml \
  https://jsonplaceholder.typicode.com
```

**Expected Output:**
```
🧪 Testing API: JSONPlaceholder API
📍 Base URL: https://jsonplaceholder.typicode.com

ℹ Running 8 tests...

✓ GET     /posts                                   - 200 OK
✓ POST    /posts                                   - 201 OK
✓ PUT     /posts/1                                 - 200 OK
✓ DELETE  /posts/1                                 - 200 OK

================================================================================
📊 Summary: 8 passed, 0 failed, 8 total
✓ All tests passed!
```

---

## Basic Usage

### Command Structure

```bash
openapi-test <command> [arguments] [options]
```

### Available Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `validate` | Validate OpenAPI spec | `openapi-test validate spec.yaml` |
| `test` | Test API endpoints | `openapi-test test spec.yaml https://api.example.com` |
| `--help` | Show help | `openapi-test --help` |
| `--version` | Show version | `openapi-test --version` |

### Common Options

| Option | Short | Description | Example |
|--------|-------|-------------|---------|
| `--verbose` | `-v` | Show detailed output | `openapi-test test spec.yaml URL -v` |
| `--quiet` | `-q` | Suppress output | `openapi-test test spec.yaml URL -q` |
| `--export` | `-e` | Export to JSON | `openapi-test test spec.yaml URL -e results.json` |
| `--timeout` | `-t` | Set timeout (ms) | `openapi-test test spec.yaml URL -t 30000` |
| `--parallel` | | Parallel execution | `openapi-test test spec.yaml URL --parallel 10` |

---

## Advanced Features

### 1. Authentication

#### Bearer Token

```bash
# Using environment variable (recommended)
export API_TOKEN="your-jwt-token"
openapi-test test spec.yaml https://api.example.com \
  --auth-bearer "$API_TOKEN"

# Direct (not recommended for production)
openapi-test test spec.yaml https://api.example.com \
  --auth-bearer "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

#### API Key in Header

```bash
# Default header: X-API-Key
openapi-test test spec.yaml https://api.example.com \
  --auth-api-key "your-api-key"

# Custom header name
openapi-test test spec.yaml https://api.example.com \
  --auth-api-key "your-api-key" \
  --auth-header "X-Custom-API-Key"
```

#### API Key in Query Parameter

```bash
openapi-test test spec.yaml https://api.example.com \
  --auth-api-key "your-api-key" \
  --auth-query "api_key"
```

#### Basic Authentication

```bash
export USERNAME="admin"
export PASSWORD="secret123"
openapi-test test spec.yaml https://api.example.com \
  --auth-basic "$USERNAME:$PASSWORD"
```

### 2. Custom Headers

```bash
# Single header
openapi-test test spec.yaml https://api.example.com \
  -H "X-Request-ID: abc123"

# Multiple headers
openapi-test test spec.yaml https://api.example.com \
  -H "X-Request-ID: abc123" \
  -H "X-Client-Version: 1.0.0" \
  -H "Accept-Language: en-US"
```

### 3. Filtering Tests

#### Filter by HTTP Method

```bash
# Test only GET requests
openapi-test test spec.yaml https://api.example.com \
  --methods GET

# Test multiple methods
openapi-test test spec.yaml https://api.example.com \
  --methods GET,POST,PUT
```

#### Filter by Path Pattern

```bash
# Test exact path
openapi-test test spec.yaml https://api.example.com \
  --paths "/users"

# Test with wildcard
openapi-test test spec.yaml https://api.example.com \
  --paths "/users/*"

# Test all admin endpoints
openapi-test test spec.yaml https://api.example.com \
  --paths "/admin/*"
```

#### Combine Filters

```bash
# Only GET requests on user endpoints
openapi-test test spec.yaml https://api.example.com \
  --methods GET \
  --paths "/users/*"
```

### 4. Export Formats

#### JSON Export

```bash
openapi-test test spec.yaml https://api.example.com \
  --export results.json

# Output: results.json
```

**JSON Format:**
```json
{
  "timestamp": "2025-11-02T18:00:00.000Z",
  "specPath": "spec.yaml",
  "baseUrl": "https://api.example.com",
  "totalTests": 10,
  "passed": 8,
  "failed": 2,
  "results": [...]
}
```

#### HTML Export

```bash
openapi-test test spec.yaml https://api.example.com \
  --export-html report.html

# Opens in browser
open report.html  # macOS
xdg-open report.html  # Linux
start report.html  # Windows
```

#### JUnit XML Export

```bash
openapi-test test spec.yaml https://api.example.com \
  --export-junit results.xml

# Perfect for CI/CD integration
```

#### Export All Formats

```bash
openapi-test test spec.yaml https://api.example.com \
  --export results.json \
  --export-html report.html \
  --export-junit results.xml
```

### 5. Schema Validation

```bash
# Enable response schema validation
openapi-test test spec.yaml https://api.example.com \
  --validate-schema

# With verbose output to see validation details
openapi-test test spec.yaml https://api.example.com \
  --validate-schema \
  --verbose
```

**Example Output:**
```
✓ GET     /users/1                                 - 200 OK (schema valid)
✗ GET     /users/2                                 - Schema validation failed: 2 error(s)
  ⚠  /email: must be valid email format
  ⚠  root: missing required property 'age'
```

### 6. Retry Logic

```bash
# Retry failed requests up to 3 times
openapi-test test spec.yaml https://api.example.com \
  --retry 3

# With verbose to see retry attempts
openapi-test test spec.yaml https://api.example.com \
  --retry 3 \
  --verbose
```

**Example Output:**
```
  ↻  Retry attempt 1/3 after 1000ms...
  ↻  Retry attempt 2/3 after 2000ms...
✓ GET     /users                                   - 200 OK
```

### 7. Parallel Execution

```bash
# Run with 10 concurrent requests
openapi-test test spec.yaml https://api.example.com \
  --parallel 10

# Auto-detect concurrency (default: 5)
openapi-test test spec.yaml https://api.example.com \
  --parallel 5
```

**Performance Impact:**
- Sequential: 100 endpoints ≈ 30 seconds
- Parallel (5): 100 endpoints ≈ 25 seconds (20% faster)
- Parallel (10): 100 endpoints ≈ 20 seconds (33% faster)

### 8. Watch Mode

```bash
# Watch for file changes and re-run tests
openapi-test test spec.yaml https://api.example.com \
  --watch

# Combine with other options
openapi-test test spec.yaml https://api.example.com \
  --watch \
  --validate-schema \
  --verbose
```

**Use Cases:**
- Development workflow
- TDD (Test-Driven Development)
- Real-time API validation
- Spec editing

**Press Ctrl+C to stop watch mode**

---

## Configuration

### Configuration Files

Create a configuration file to avoid repeating options:

#### YAML Format

```yaml
# .openapi-cli.yaml
auth-bearer: "your-jwt-token"
timeout: 15000
verbose: true
parallel: 10
validate-schema: true

headers:
  - "User-Agent: OpenAPI-CLI/1.0"
  - "Accept: application/json"

methods: "GET,POST,PUT"
export-html: "report.html"
export-junit: "results.xml"
```

#### JSON Format

```json
{
  "auth-bearer": "your-jwt-token",
  "timeout": 15000,
  "verbose": true,
  "parallel": 10,
  "validate-schema": true,
  "headers": [
    "User-Agent: OpenAPI-CLI/1.0",
    "Accept: application/json"
  ],
  "methods": "GET,POST,PUT",
  "export-html": "report.html"
}
```

### Auto-Discovery

The CLI automatically searches for config files:

1. `.openapi-cli.yaml` (current directory)
2. `.openapi-cli.yml` (current directory)
3. `.openapi-cli.json` (current directory)
4. `openapi-cli.yaml` (current directory)
5. Parent directories (up to root)

### Explicit Config File

```bash
# Use specific config file
openapi-test test spec.yaml https://api.example.com \
  --config my-config.yaml
```

### Option Precedence

**Priority (highest to lowest):**
1. Command-line options
2. Config file options
3. Default values

**Example:**
```bash
# Config file has: timeout: 10000
# This command uses timeout: 30000 (CLI takes precedence)
openapi-test test spec.yaml URL --config .openapi-cli.yaml --timeout 30000
```

---

## CI/CD Integration

### GitHub Actions

```yaml
name: API Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - uses: actions/setup-node@v3
        with:
          node-version: '18'

      - name: Install OpenAPI CLI
        run: |
          cd openapi-cli
          npm install
          npm run build
          npm link

      - name: Run API Tests
        env:
          API_TOKEN: ${{ secrets.API_TOKEN }}
        run: |
          openapi-test test api-spec.yaml https://api.example.com \
            --auth-bearer "$API_TOKEN" \
            --export-junit results.xml \
            --quiet

      - name: Publish Test Results
        uses: EnricoMi/publish-unit-test-result-action@v2
        if: always()
        with:
          files: results.xml
```

### GitLab CI

```yaml
api-tests:
  stage: test
  image: node:18
  script:
    - cd openapi-cli
    - npm install
    - npm run build
    - npm link
    - openapi-test test ../api-spec.yaml https://api.example.com
        --auth-bearer "$API_TOKEN"
        --export-junit results.xml
        --quiet
  artifacts:
    reports:
      junit: results.xml
    when: always
```

### Jenkins

```groovy
pipeline {
    agent any

    environment {
        API_TOKEN = credentials('api-token')
    }

    stages {
        stage('Install') {
            steps {
                sh '''
                    cd openapi-cli
                    npm install
                    npm run build
                    npm link
                '''
            }
        }

        stage('Test API') {
            steps {
                sh '''
                    openapi-test test api-spec.yaml https://api.example.com \
                        --auth-bearer "$API_TOKEN" \
                        --export-junit results.xml \
                        --quiet
                '''
            }
        }
    }

    post {
        always {
            junit 'results.xml'
        }
    }
}
```

### Docker

**Dockerfile:**
```dockerfile
FROM node:18-alpine

WORKDIR /app

COPY openapi-cli/package*.json ./
RUN npm install

COPY openapi-cli/ ./
RUN npm run build && npm link

ENTRYPOINT ["openapi-test"]
```

**Build and Run:**
```bash
# Build image
docker build -t openapi-cli .

# Run tests
docker run --rm \
  -v $(pwd)/api-spec.yaml:/spec.yaml \
  openapi-cli test /spec.yaml https://api.example.com
```

---

## Real-World Examples

### Example 1: Complete Test Suite

```bash
#!/bin/bash
# complete-test.sh - Full API test with all features

SPEC="api-spec.yaml"
BASE_URL="https://api.production.com"
TOKEN="${API_TOKEN}"

echo "🧪 Running comprehensive API tests..."

openapi-test test "$SPEC" "$BASE_URL" \
  --auth-bearer "$TOKEN" \
  --validate-schema \
  --retry 3 \
  --parallel 10 \
  --timeout 30000 \
  --verbose \
  --export results.json \
  --export-html report.html \
  --export-junit junit.xml \
  -H "X-Client-Version: 1.0.0" \
  -H "X-Environment: production"

if [ $? -eq 0 ]; then
  echo "✅ All tests passed!"
  exit 0
else
  echo "❌ Some tests failed!"
  exit 1
fi
```

### Example 2: Smoke Tests

```bash
#!/bin/bash
# smoke-test.sh - Quick health check

openapi-test test api-spec.yaml https://api.example.com \
  --methods GET \
  --paths "/health,/status,/version" \
  --timeout 5000 \
  --quiet

exit $?
```

### Example 3: Environment-Specific Tests

```bash
#!/bin/bash
# test-environment.sh - Test different environments

ENV=${1:-staging}

case $ENV in
  staging)
    URL="https://staging.api.example.com"
    CONFIG=".openapi-cli.staging.yaml"
    ;;
  production)
    URL="https://api.example.com"
    CONFIG=".openapi-cli.prod.yaml"
    ;;
  *)
    echo "Unknown environment: $ENV"
    exit 1
    ;;
esac

echo "Testing $ENV environment..."

openapi-test test api-spec.yaml "$URL" \
  --config "$CONFIG" \
  --export-html "report-$ENV.html"
```

### Example 4: Pre-Commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Validating OpenAPI specification..."

openapi-test validate api/openapi.yaml

if [ $? -ne 0 ]; then
  echo "❌ OpenAPI spec validation failed!"
  echo "Fix the spec before committing."
  exit 1
fi

echo "✅ OpenAPI spec is valid"
exit 0
```

### Example 5: Monitoring Script

```bash
#!/bin/bash
# monitor.sh - API health monitoring

SPEC="api-spec.yaml"
URL="https://api.production.com"
SLACK_WEBHOOK="${SLACK_WEBHOOK_URL}"

# Run tests
openapi-test test "$SPEC" "$URL" \
  --quiet \
  --retry 2 \
  --export results.json

EXIT_CODE=$?

# Parse results
TOTAL=$(jq '.totalTests' results.json)
FAILED=$(jq '.failed' results.json)
PASSED=$(jq '.passed' results.json)

if [ $EXIT_CODE -ne 0 ]; then
  # Send alert to Slack
  curl -X POST "$SLACK_WEBHOOK" \
    -H 'Content-Type: application/json' \
    -d "{\"text\":\"🚨 API Health Check Failed\n$FAILED/$TOTAL tests failed\"}"

  exit 1
fi

echo "✅ API health check passed ($PASSED/$TOTAL)"
exit 0
```

---

## Troubleshooting

### Common Issues

#### 1. "Command not found: openapi-test"

**Problem:** CLI not installed globally

**Solution:**
```bash
cd openapi-cli
npm link
```

#### 2. Connection Refused

**Problem:** API server not running or URL incorrect

**Solutions:**
```bash
# Check if server is running
curl https://api.example.com/health

# Verify URL in spec matches base URL
openapi-test test spec.yaml https://api.example.com --verbose

# Check firewall/network
ping api.example.com
```

#### 3. Timeout Errors

**Problem:** Requests taking too long

**Solutions:**
```bash
# Increase timeout
openapi-test test spec.yaml URL --timeout 60000

# Check API performance
curl -w "@curl-format.txt" -o /dev/null -s https://api.example.com/endpoint
```

#### 4. Authentication Failures

**Problem:** 401 Unauthorized errors

**Solutions:**
```bash
# Verify token
echo $API_TOKEN

# Check token format
openapi-test test spec.yaml URL --auth-bearer "$API_TOKEN" --verbose

# Test authentication manually
curl -H "Authorization: Bearer $API_TOKEN" https://api.example.com/endpoint
```

#### 5. Schema Validation Errors

**Problem:** Response doesn't match schema

**Solutions:**
```bash
# View detailed errors
openapi-test test spec.yaml URL --validate-schema --verbose

# Check specific endpoint response
curl https://api.example.com/endpoint | jq .

# Verify schema in OpenAPI spec
yq eval '.paths."/endpoint".get.responses.200.content."application/json".schema' spec.yaml
```

### Debug Mode

```bash
# Maximum verbosity
openapi-test test spec.yaml URL \
  --verbose \
  --validate-schema \
  --retry 3 \
  2>&1 | tee debug.log
```

---

## Best Practices

### 1. Use Configuration Files

**Don't:**
```bash
openapi-test test spec.yaml URL --auth-bearer "$TOKEN" --timeout 30000 --parallel 10 --validate-schema --export report.html
```

**Do:**
```yaml
# .openapi-cli.yaml
auth-bearer: ${API_TOKEN}
timeout: 30000
parallel: 10
validate-schema: true
export-html: report.html
```

```bash
openapi-test test spec.yaml URL
```

### 2. Use Environment Variables for Secrets

**Don't:**
```bash
openapi-test test spec.yaml URL --auth-bearer "hardcoded-token-123"
```

**Do:**
```bash
export API_TOKEN="your-secret-token"
openapi-test test spec.yaml URL --auth-bearer "$API_TOKEN"
```

### 3. Use Quiet Mode in CI/CD

**Don't:**
```bash
# Too much output in CI logs
openapi-test test spec.yaml URL --verbose
```

**Do:**
```bash
# Clean CI logs, only errors shown
openapi-test test spec.yaml URL --quiet --export-junit results.xml
```

### 4. Use Parallel Execution for Large APIs

**Don't:**
```bash
# Sequential testing of 100 endpoints (slow)
openapi-test test large-spec.yaml URL
```

**Do:**
```bash
# Parallel testing (much faster)
openapi-test test large-spec.yaml URL --parallel 10
```

### 5. Use Filtering for Focused Testing

**Don't:**
```bash
# Test all 100 endpoints when only testing user endpoints
openapi-test test spec.yaml URL
```

**Do:**
```bash
# Filter to relevant endpoints
openapi-test test spec.yaml URL --paths "/users/*"
```

### 6. Use Retry Logic for Flaky Networks

**Don't:**
```bash
# Fail immediately on network hiccups
openapi-test test spec.yaml URL
```

**Do:**
```bash
# Retry on network errors
openapi-test test spec.yaml URL --retry 3
```

### 7. Version Your Config Files

**Do:**
```bash
# Commit config files to git
git add .openapi-cli.yaml
git commit -m "Add API test configuration"
```

**Structure:**
```
project/
├── .openapi-cli.yaml          # Default config
├── .openapi-cli.staging.yaml  # Staging env
├── .openapi-cli.prod.yaml     # Production env
└── api/
    └── openapi.yaml           # API spec
```

### 8. Use Schema Validation in CI/CD

**Do:**
```bash
# Catch contract violations early
openapi-test test spec.yaml URL \
  --validate-schema \
  --export-junit results.xml \
  --quiet
```

### 9. Export Multiple Formats

**Do:**
```bash
# JSON for automation, HTML for humans, JUnit for CI
openapi-test test spec.yaml URL \
  --export results.json \
  --export-html report.html \
  --export-junit results.xml
```

### 10. Document Your Test Scripts

**Do:**
```bash
#!/bin/bash
# test-api.sh - Production API health check
#
# Usage: ./test-api.sh [environment]
#
# Environments:
#   - staging (default)
#   - production
#
# Environment variables:
#   - API_TOKEN: Authentication token (required)
#
# Exit codes:
#   0: All tests passed
#   1: Some tests failed

ENV=${1:-staging}

# ... script continues
```

---

## Summary

The OpenAPI CLI provides a comprehensive toolkit for automated API testing. Key takeaways:

- ✅ **Easy to install** - npm install and link
- ✅ **Flexible authentication** - Bearer, API Key, Basic
- ✅ **Multiple export formats** - JSON, HTML, JUnit XML
- ✅ **CI/CD ready** - Exit codes, quiet mode, exports
- ✅ **High performance** - Parallel execution, configurable concurrency
- ✅ **Developer friendly** - Config files, watch mode, verbose logging
- ✅ **Production ready** - Retry logic, schema validation, timeout control

### Quick Reference

```bash
# Basic usage
openapi-test validate spec.yaml
openapi-test test spec.yaml https://api.example.com

# With authentication
openapi-test test spec.yaml URL --auth-bearer "$TOKEN"

# Full featured
openapi-test test spec.yaml URL \
  --auth-bearer "$TOKEN" \
  --validate-schema \
  --retry 3 \
  --parallel 10 \
  --export-html report.html \
  --quiet
```

For more details, see:
- [README.md](../README.md) - Feature overview
- [ARCHITECTURE.md](ARCHITECTURE.md) - System design
- [PROGRESS.md](PROGRESS.md) - Development status

---

**Made with ❤️ for the OpenAPI community**

💡 **Pro Tip:** Combine with the OpenAPI TUI for the complete testing experience!

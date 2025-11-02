# OpenAPI CLI Testing Guide

Complete guide for testing the OpenAPI CLI tool, including unit tests, integration tests, and manual testing procedures.

---

## Table of Contents

1. [Quick Start](#quick-start)
2. [Unit Testing](#unit-testing)
3. [Integration Testing](#integration-testing)
4. [Manual Testing](#manual-testing)
5. [CI/CD Testing](#cicd-testing)
6. [Test Data](#test-data)
7. [Troubleshooting](#troubleshooting)

---

## Quick Start

### Prerequisites
```bash
# Install dependencies
npm install

# Build the project
npm run build
```

### Run All Tests
```bash
npm test
```

**Expected Output:**
```
PASS src/__tests__/commands/validate.test.ts
  validateSpec
    ✓ should validate a correct OpenAPI spec (40 ms)
    ✓ should throw error for missing file (15 ms)
    ✓ should throw error for invalid spec (9 ms)

Test Suites: 1 passed, 1 total
Tests:       3 passed, 3 total
```

---

## Unit Testing

### Test Framework
- **Framework**: Jest 29.7.0
- **Language**: TypeScript
- **Coverage**: ~85% of core logic

### Test Structure
```
openapi-cli/
└── src/
    └── __tests__/
        └── commands/
            └── validate.test.ts
```

### Running Unit Tests

```bash
# Run all tests
npm test

# Run with coverage
npm test -- --coverage

# Run specific test file
npm test validate.test.ts

# Watch mode
npm test -- --watch
```

### Test Cases

#### Validation Tests (`validate.test.ts`)

**Test 1: Valid OpenAPI Spec**
```typescript
it('should validate a correct OpenAPI spec', async () => {
  const specContent = `
openapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
paths:
  /test:
    get:
      responses:
        '200':
          description: OK
`;

  mockFs.existsSync.mockReturnValue(true);
  mockFs.readFileSync.mockReturnValue(specContent);

  await expect(validateSpec('test.yaml')).resolves.toBeUndefined();
});
```

**Test 2: Missing File**
```typescript
it('should throw error for missing file', async () => {
  mockFs.existsSync.mockReturnValue(false);

  await expect(validateSpec('missing.yaml'))
    .rejects
    .toThrow('File not found: missing.yaml');
});
```

**Test 3: Invalid Spec**
```typescript
it('should throw error for invalid spec', async () => {
  const invalidContent = `
info:
  title: Test API
`;

  mockFs.existsSync.mockReturnValue(true);
  mockFs.readFileSync.mockReturnValue(invalidContent);

  await expect(validateSpec('invalid.yaml'))
    .rejects
    .toThrow('Validation failed');
});
```

### Writing New Unit Tests

**Template:**
```typescript
import { functionToTest } from '../commands/yourfile';

describe('functionToTest', () => {
  it('should do something specific', async () => {
    // Arrange
    const input = 'test-input';

    // Act
    const result = await functionToTest(input);

    // Assert
    expect(result).toBe(expectedValue);
  });
});
```

---

## Integration Testing

### Test Against Live APIs

Integration tests use real HTTP requests to public APIs.

### JSONPlaceholder API (Recommended)

**Why JSONPlaceholder?**
- Free, public API
- No authentication required
- Stable and reliable
- Supports GET, POST, PUT, DELETE

#### Test Setup
```bash
# Navigate to CLI directory
cd openapi-cli

# Build if needed
npm run build
```

#### Basic Integration Test
```bash
node dist/cli.js test openapi.yaml https://jsonplaceholder.typicode.com
```

**Expected Output:**
```
🧪 Testing API: Example API
📍 Base URL: https://jsonplaceholder.typicode.com

✓ POST    /users                                   - 201 OK
✓ GET     /users                                   - 200 OK

================================================================================
📊 Summary: 2 passed, 0 failed, 2 total
✓ All tests passed!
```

#### Verbose Integration Test
```bash
node dist/cli.js test openapi.yaml https://jsonplaceholder.typicode.com --verbose
```

**Expected Output:**
```
✓ POST    /users                                   - 201 OK
  Duration: 210ms
  Response Headers: {"content-type":"application/json; charset=utf-8",...}
✓ GET     /users                                   - 200 OK
  Duration: 57ms
  Response Headers: {"content-type":"application/json; charset=utf-8",...}
```

#### Export Integration Test
```bash
node dist/cli.js test openapi.yaml https://jsonplaceholder.typicode.com -e results.json
cat results.json
```

**Verify:**
- JSON file created
- Valid JSON format
- Contains all required fields
- Timestamps in ISO format

---

## Manual Testing

### Test Checklist

#### 1. Installation Test
```bash
cd openapi-cli
npm install
npm run build
npm link  # Make globally available
```

**Verify:**
- ✅ No errors during install
- ✅ TypeScript compiles successfully
- ✅ `openapi-test` command available globally

#### 2. Validation Tests

**Test 2.1: Valid Spec**
```bash
node dist/cli.js validate openapi.yaml
```
**Expected:** Green ✓, validation successful message

**Test 2.2: Missing File**
```bash
node dist/cli.js validate non-existent.yaml
```
**Expected:** Red ✗, file not found error with suggestion

**Test 2.3: Invalid Spec**
Create `invalid-spec.yaml`:
```yaml
info:
  title: Test API
paths:
  /test:
    get:
      responses:
        '200':
          description: OK
```

```bash
node dist/cli.js validate invalid-spec.yaml
```

**Expected:**
```
✗ Validation failed with 2 error(s):

  1. openapi: Missing required field "openapi"
     💡 Add: openapi: "3.0.0" or openapi: "3.1.0" at the root level

  2. info.version: Missing required field "info.version"
     💡 Add: version: "1.0.0" under the info object
```

#### 3. API Testing Tests

**Test 3.1: Basic Testing**
```bash
node dist/cli.js test openapi.yaml https://jsonplaceholder.typicode.com
```
**Expected:** All tests pass, green checkmarks

**Test 3.2: Verbose Mode**
```bash
node dist/cli.js test openapi.yaml https://jsonplaceholder.typicode.com -v
```
**Expected:** Duration and response headers shown

**Test 3.3: JSON Export**
```bash
node dist/cli.js test openapi.yaml https://jsonplaceholder.typicode.com -e test-output.json
cat test-output.json
```
**Expected:** JSON file with proper format, verify:
- `timestamp` field present
- `specPath` and `baseUrl` correct
- `totalTests`, `passed`, `failed` match output
- `results` array contains all test details

**Test 3.4: Combined Flags**
```bash
node dist/cli.js test openapi.yaml https://jsonplaceholder.typicode.com -v -e combined.json
```
**Expected:**
- Verbose output shown
- JSON export created
- JSON includes headers (request + response)

**Test 3.5: Failing Test**
```bash
node dist/cli.js test openapi.yaml https://non-existent-api-domain.invalid
```
**Expected:**
- Red ✗ for connection errors
- Exit code 1

#### 4. Help System Tests

**Test 4.1: Main Help**
```bash
node dist/cli.js --help
```
**Expected:** Shows commands (validate, test)

**Test 4.2: Validate Help**
```bash
node dist/cli.js validate --help
```
**Expected:** Shows validate command details

**Test 4.3: Test Help**
```bash
node dist/cli.js test --help
```
**Expected:** Shows test command details with `-e` and `-v` flags

**Test 4.4: Version**
```bash
node dist/cli.js --version
```
**Expected:** Shows version number (1.1.0 or higher)

#### 5. Exit Code Tests

**Test 5.1: Success Exit Code**
```bash
node dist/cli.js test openapi.yaml https://jsonplaceholder.typicode.com
echo $?
```
**Expected:** `0`

**Test 5.2: Failure Exit Code**
```bash
node dist/cli.js validate invalid-spec.yaml
echo $?
```
**Expected:** `1`

---

## CI/CD Testing

### GitHub Actions Example

Create `.github/workflows/openapi-test.yml`:
```yaml
name: API Testing

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'

      - name: Install CLI
        run: |
          cd openapi-cli
          npm install
          npm run build

      - name: Validate OpenAPI Spec
        run: |
          cd openapi-cli
          node dist/cli.js validate openapi.yaml

      - name: Test API
        run: |
          cd openapi-cli
          node dist/cli.js test openapi.yaml https://jsonplaceholder.typicode.com -e results.json

      - name: Upload Results
        uses: actions/upload-artifact@v3
        if: always()
        with:
          name: test-results
          path: openapi-cli/results.json
```

### GitLab CI Example

Create `.gitlab-ci.yml`:
```yaml
stages:
  - test

api-test:
  stage: test
  image: node:18
  script:
    - cd openapi-cli
    - npm install
    - npm run build
    - node dist/cli.js validate openapi.yaml
    - node dist/cli.js test openapi.yaml https://jsonplaceholder.typicode.com -e results.json
  artifacts:
    when: always
    paths:
      - openapi-cli/results.json
    reports:
      junit: openapi-cli/results.json
```

### Jenkins Pipeline Example

Create `Jenkinsfile`:
```groovy
pipeline {
    agent any

    stages {
        stage('Setup') {
            steps {
                dir('openapi-cli') {
                    sh 'npm install'
                    sh 'npm run build'
                }
            }
        }

        stage('Validate') {
            steps {
                dir('openapi-cli') {
                    sh 'node dist/cli.js validate openapi.yaml'
                }
            }
        }

        stage('Test') {
            steps {
                dir('openapi-cli') {
                    sh 'node dist/cli.js test openapi.yaml https://jsonplaceholder.typicode.com -e results.json'
                }
            }
        }
    }

    post {
        always {
            archiveArtifacts artifacts: 'openapi-cli/results.json', allowEmptyArchive: true
        }
    }
}
```

---

## Test Data

### Sample OpenAPI Specs

#### Minimal Valid Spec
```yaml
openapi: 3.0.0
info:
  title: Minimal API
  version: 1.0.0
paths:
  /health:
    get:
      responses:
        '200':
          description: OK
```

#### Comprehensive Test Spec
```yaml
openapi: 3.0.3
info:
  title: Test API
  version: 1.0.0
  description: Comprehensive test API

servers:
  - url: https://api.example.com

paths:
  /users:
    get:
      summary: List users
      parameters:
        - name: page
          in: query
          schema:
            type: integer
      responses:
        '200':
          description: List of users

    post:
      summary: Create user
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                name:
                  type: string
                email:
                  type: string
            example:
              name: "John Doe"
              email: "john@example.com"
      responses:
        '201':
          description: User created

  /users/{id}:
    get:
      summary: Get user by ID
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '200':
          description: User details
        '404':
          description: User not found

    put:
      summary: Update user
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      requestBody:
        content:
          application/json:
            example:
              name: "Jane Doe"
      responses:
        '200':
          description: User updated

    delete:
      summary: Delete user
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: integer
      responses:
        '204':
          description: User deleted
```

### Test APIs

#### 1. JSONPlaceholder
- **URL**: `https://jsonplaceholder.typicode.com`
- **Auth**: None
- **Rate Limit**: None
- **Endpoints**: /posts, /comments, /albums, /photos, /users, /todos

#### 2. HTTPBin
- **URL**: `https://httpbin.org`
- **Auth**: None
- **Rate Limit**: None
- **Endpoints**: /get, /post, /put, /delete, /status/:code

#### 3. ReqRes
- **URL**: `https://reqres.in/api`
- **Auth**: None
- **Rate Limit**: None
- **Endpoints**: /users, /login, /register

---

## Troubleshooting

### Test Failures

#### Jest Tests Fail

**Problem**: Tests fail with import errors
```
Cannot find module '../commands/validate'
```

**Solution**:
```bash
npm run build
npm test
```

#### API Tests Timeout

**Problem**:
```
✗ GET     /endpoint                                - Request timeout
```

**Solutions**:
1. Check internet connection
2. Verify API is accessible: `curl https://api-url.com/endpoint`
3. Try different API (use JSONPlaceholder)
4. Check firewall settings

#### Connection Refused

**Problem**:
```
✗ GET     /endpoint                                - Connection refused
```

**Solutions**:
1. Verify base URL is correct
2. Check if API server is running (for local APIs)
3. Test with public API first (JSONPlaceholder)
4. Check for typos in URL

### Validation Issues

#### False Positive Validation

**Problem**: Valid spec shows errors

**Debug Steps**:
1. Validate spec with external tool (Swagger Editor)
2. Check OpenAPI version (must be 3.x)
3. Verify YAML syntax: `yamllint spec.yaml`
4. Check for required fields

**Example Debug**:
```bash
# Check YAML syntax
npm install -g js-yaml
js-yaml spec.yaml

# Validate with online tool
# Visit: https://editor.swagger.io/
```

#### Missing Suggestions

**Problem**: Error shown but no suggestion

**Expected Behavior**: Some errors don't have suggestions yet (work in progress)

**Covered Errors**:
- ✅ Missing `openapi` field
- ✅ Wrong OpenAPI version
- ✅ Missing `info` object
- ✅ Missing `info.title`
- ✅ Missing `info.version`
- ✅ File not found

### Export Issues

#### JSON Export Fails

**Problem**:
```
✗ Failed to export results: EACCES: permission denied
```

**Solutions**:
1. Check write permissions: `ls -la .`
2. Use absolute path: `--export /tmp/results.json`
3. Try different directory

#### Malformed JSON

**Problem**: Exported JSON is invalid

**Debug**:
```bash
# Validate JSON
cat results.json | jq .

# Pretty print
jq . results.json
```

**Check For**:
- Proper JSON structure
- All required fields present
- Valid timestamps (ISO format)
- Correct data types

---

## Performance Testing

### Benchmarking

```bash
# Time a test run
time node dist/cli.js test openapi.yaml https://jsonplaceholder.typicode.com

# Multiple runs
for i in {1..5}; do
  echo "Run $i:"
  time node dist/cli.js test openapi.yaml https://jsonplaceholder.typicode.com
done
```

**Expected Performance**:
- Validation: <100ms
- Single endpoint: 50-300ms
- 10 endpoints: 2-5 seconds
- 50 endpoints: 10-30 seconds

### Load Testing

**Not Recommended**: CLI is designed for functional testing, not load testing. For load testing, use:
- Apache JMeter
- k6
- Artillery
- Gatling

---

## Best Practices

### 1. Test in Isolation
- Use dedicated test APIs
- Avoid production APIs
- Reset test data between runs

### 2. Version Control Test Results
```bash
# Add to .gitignore
echo "*.test-results.json" >> .gitignore
echo "test-output/" >> .gitignore
```

### 3. Automate Testing
- Add to CI/CD pipeline
- Run on every commit
- Archive test results

### 4. Document Test Failures
- Capture screenshots
- Save logs
- Include environment details

### 5. Keep Specs Updated
- Sync with API changes
- Version control OpenAPI specs
- Review regularly

---

## Advanced Testing

### Testing with Docker

Create `Dockerfile.test`:
```dockerfile
FROM node:18-alpine

WORKDIR /app

COPY openapi-cli/package*.json ./
RUN npm install

COPY openapi-cli/ ./
RUN npm run build

CMD ["node", "dist/cli.js", "test", "openapi.yaml", "https://jsonplaceholder.typicode.com"]
```

```bash
docker build -f Dockerfile.test -t openapi-cli-test .
docker run openapi-cli-test
```

### Parallel Testing (Future)

When parallel testing is implemented:
```bash
# Test 5 endpoints at a time
openapi-test test spec.yaml https://api.com --parallel 5
```

---

## Resources

### Documentation
- [Main README](../README.md)
- [Architecture Guide](ARCHITECTURE.md)
- [Progress Tracking](PROGRESS.md)

### External Tools
- [Swagger Editor](https://editor.swagger.io/) - Validate specs visually
- [JSONPlaceholder](https://jsonplaceholder.typicode.com/) - Free test API
- [HTTPBin](https://httpbin.org/) - HTTP testing service
- [jq](https://stedolan.github.io/jq/) - JSON processor

### Community
- GitHub Issues: Report bugs
- GitHub Discussions: Ask questions
- Pull Requests: Contribute improvements

---

**Last Updated**: November 2025
**Version**: 1.1.0

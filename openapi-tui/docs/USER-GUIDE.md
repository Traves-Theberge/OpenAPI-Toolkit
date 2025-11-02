# OpenAPI TUI - Complete User Guide

A comprehensive end-to-end guide for using the OpenAPI TUI (Terminal User Interface).

## Table of Contents

- [Getting Started](#getting-started)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [User Interface](#user-interface)
- [Core Features](#core-features)
- [Advanced Features](#advanced-features)
- [Configuration](#configuration)
- [Export & Reporting](#export--reporting)
- [Real-World Workflows](#real-world-workflows)
- [Troubleshooting](#troubleshooting)
- [Best Practices](#best-practices)

---

## Getting Started

### What is OpenAPI TUI?

OpenAPI TUI is an interactive terminal application that:
- ✅ Validates OpenAPI 3.x specifications with visual feedback
- ✅ Tests API endpoints interactively
- ✅ Generates request bodies from schemas automatically
- ✅ Validates API responses against schemas
- ✅ Provides beautiful, styled terminal UI
- ✅ Exports results in multiple formats (JSON, HTML, JUnit XML)
- ✅ Supports authentication (Bearer, API Key, Basic)
- ✅ Tracks test history with replay functionality
- ✅ Runs tests in parallel for better performance

### When to Use TUI

**Perfect for:**
- ✅ Interactive API development and testing
- ✅ Real-time API exploration
- ✅ Manual debugging and troubleshooting
- ✅ Learning API behavior
- ✅ Quick spec validation during development
- ✅ Visual feedback requirements
- ✅ Test run history and replay

**Not ideal for:**
- ❌ CI/CD automation (use CLI instead)
- ❌ Headless environments (use CLI instead)
- ❌ Scripting and batch processing (use CLI instead)

---

## Installation

### Prerequisites

- **Go**: Version 1.21 or higher
- **Terminal**: Modern terminal with color support (iTerm2, Alacritty, Windows Terminal, etc.)

### Install from Source

```bash
# Clone the repository
git clone https://github.com/Traves-Theberge/OpenAPI-Toolkit.git
cd OpenAPI-Toolkit/openapi-tui

# Install dependencies
go mod tidy

# Build the binary
go build -o openapi-tui .

# Optionally install globally
go install .
```

### Install Directly

```bash
go install github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui@latest
```

### Verify Installation

```bash
# Run the TUI
openapi-tui

# Should display the main menu
```

---

## Quick Start

### 1. Launch the TUI

```bash
openapi-tui
```

You'll see the main menu with 8 options:

```
┌─────────────────────────────────────┐
│     📋 OpenAPI Specification        │
│         Validator & Tester          │
├─────────────────────────────────────┤
│ 1. 📋 Validate OpenAPI Spec         │
│ 2. 🧪 Test API Endpoints            │
│ 3. 📄 Custom Request                │
│ 4. 🎯 Select Endpoints              │
│ 5. 📜 View Test History             │
│ 6. ⚙️  Settings                     │
│ 7. 📚 Help                          │
│ 8. 🚪 Quit                          │
└─────────────────────────────────────┘
```

### 2. Validate Your First Spec

1. Select option **1** (Validate OpenAPI Spec)
2. Enter the path to your spec file (e.g., `api-spec.yaml`)
3. View validation results instantly

**Example:**
```
✓ Validation Successful!
  OpenAPI Version: 3.0.0
  Title: My API
  Version: 1.0.0
  Paths: 5
  Operations: 12
```

### 3. Test Your First API

1. Select option **2** (Test API Endpoints)
2. Enter spec file path
3. Enter base URL (e.g., `https://jsonplaceholder.typicode.com`)
4. Watch tests run with real-time updates
5. View results in styled table

**Interactive features:**
- Press **'v'** to toggle verbose mode
- Press **'f'** to filter results
- Press **'e'** to export JSON
- Press **'h'** to export HTML
- Press **'r'** to view history

---

## User Interface

### Navigation Keys

#### Global Keys (Work Everywhere)
| Key | Action |
|-----|--------|
| **↑ / ↓** | Navigate up/down |
| **j / k** | Vim-style navigation (down/up) |
| **Enter** | Select / Confirm |
| **Esc** | Go back / Cancel |
| **q** | Quit (from menu) |
| **h / ?** | Show help screen |

#### Main Menu
| Key | Action |
|-----|--------|
| **1-8** | Select menu option (direct) |
| **↑ / ↓** | Navigate menu |
| **Enter** | Confirm selection |
| **v** | Toggle verbose mode (shows in status bar) |
| **q** | Quit application |

#### Test Results Screen
| Key | Action |
|-----|--------|
| **v** | Toggle verbose logging |
| **f** | Enter filter mode |
| **e** | Export results to JSON |
| **h** | Export results to HTML |
| **j** | Export results to JUnit XML |
| **r** | View test run history |
| **l** | View detailed logs (verbose mode only) |
| **↑ / ↓** | Scroll through results |
| **Enter** | Return to menu |

#### Filter Mode
| Key | Action |
|-----|--------|
| **Type** | Enter filter text |
| **Backspace** | Delete character |
| **Esc** | Exit filter mode |
| **Enter** | Apply filter and return to menu |

Filters match: status codes, HTTP methods, endpoints, or keywords

#### History Screen
| Key | Action |
|-----|--------|
| **↑ / ↓** or **j / k** | Navigate history entries |
| **Enter** | Replay selected test run |
| **Esc** | Return to results screen |

#### Configuration Editor
| Key | Action |
|-----|--------|
| **Tab** | Move to next field |
| **Shift+Tab** | Move to previous field |
| **↑ / ↓** | Also navigate fields |
| **Type** | Edit current field |
| **Enter** | Save configuration |
| **Esc** | Cancel and return to menu |

**12 Configurable Fields:**
1. Base URL
2. Spec File Path
3. Timeout (ms)
4. Auth Type (none/bearer/apikey/basic)
5. Auth Token/Key
6. API Key Header Name
7. API Key Query Param
8. Basic Auth Username
9. Basic Auth Password
10. Worker Count (parallel testing)
11. Max Retries (0-10)
12. Verbose Mode (true/false)

#### Custom Request Screen
| Key | Action |
|-----|--------|
| **Tab** | Cycle through fields (Method → URL → Headers → Body) |
| **Type** | Edit current field |
| **Enter** | Execute request (when focused on body) |
| **Esc** | Return to menu |

#### Endpoint Selector
| Key | Action |
|-----|--------|
| **Space** | Toggle selection for current endpoint |
| **a** | Select all endpoints |
| **n** | Deselect all endpoints |
| **↑ / ↓** | Navigate endpoint list |
| **Enter** | Run tests on selected endpoints |
| **Esc** | Cancel and return to menu |

---

## Core Features

### 1. Spec Validation

**Purpose**: Verify OpenAPI specification is valid before testing

**How to Use:**
1. Main menu → **1. Validate OpenAPI Spec**
2. Enter spec file path (supports `.yaml`, `.yml`, `.json`)
3. View instant validation results

**What Gets Validated:**
- ✅ OpenAPI version (3.0.x or 3.1.x)
- ✅ Required fields (info, paths, operations)
- ✅ Path format (must start with `/`)
- ✅ Response definitions
- ✅ Schema structures
- ✅ Reference resolution ($ref)

**Example Output:**
```
┌───────────────────────────────────────┐
│      ✓ Validation Successful!        │
├───────────────────────────────────────┤
│ OpenAPI Version: 3.0.0                │
│ Title: JSONPlaceholder API            │
│ Version: 1.0.0                        │
│ Paths: 8                              │
│ Operations: 12                        │
└───────────────────────────────────────┘
```

### 2. Endpoint Testing

**Purpose**: Automatically test all endpoints defined in your spec

**How to Use:**
1. Main menu → **2. Test API Endpoints**
2. Enter spec file path
3. Enter base URL (e.g., `https://api.example.com`)
4. Watch tests run with progress indicators

**What Happens:**
- Loads and parses OpenAPI spec
- Generates requests for each operation
- Automatically generates request bodies from schemas
- Substitutes path parameters (`{id}` → `1`)
- Builds query parameters from spec
- Executes HTTP requests
- Validates responses against schemas
- Displays results in real-time

**Test Results Table:**
```
┌──────────────────────────────────────────────────────┐
│           📊 Test Results - 12 endpoints             │
├──────────────────────────────────────────────────────┤
│ Method │ Endpoint        │ Status │ Duration │ Msg  │
├────────┼─────────────────┼────────┼──────────┼──────┤
│ GET    │ /posts          │ ✅ 200 │ 125ms    │ OK   │
│ POST   │ /posts          │ ✅ 201 │ 156ms    │ OK   │
│ PUT    │ /posts/1        │ ✅ 200 │ 98ms     │ OK   │
│ DELETE │ /posts/1        │ ✅ 200 │ 87ms     │ OK   │
└────────┴─────────────────┴────────┴──────────┴──────┘

Statistics:
  Total: 12 | Passed: 12 | Failed: 0 | Pass Rate: 100%
  Timing: Total 1.2s | Avg 100ms | Fastest 87ms | Slowest 156ms
```

### 3. Custom Requests

**Purpose**: Manually craft and execute HTTP requests

**How to Use:**
1. Main menu → **3. Custom Request**
2. Use **Tab** to navigate fields:
   - **Method**: GET, POST, PUT, PATCH, DELETE
   - **URL**: Full endpoint URL
   - **Headers**: JSON object (e.g., `{"X-Custom": "value"}`)
   - **Body**: JSON request body
3. Press **Enter** to execute

**Example:**
```
Method: POST
URL: https://api.example.com/users
Headers: {"Content-Type": "application/json", "X-API-Key": "secret"}
Body: {"name": "John Doe", "email": "john@example.com"}

Press Enter to execute...
```

**Response Display:**
- Status code with color coding
- Response headers
- Response body (formatted JSON)
- Duration
- Timestamp

### 4. Endpoint Selection

**Purpose**: Test only specific endpoints instead of all

**How to Use:**
1. Main menu → **4. Select Endpoints**
2. Use **Space** to toggle selection
3. Use **'a'** to select all, **'n'** to deselect all
4. Press **Enter** to run tests on selected endpoints

**Example:**
```
┌─────────────────────────────────────┐
│    Select Endpoints to Test         │
├─────────────────────────────────────┤
│ [✓] GET /posts                      │
│ [✓] POST /posts                     │
│ [ ] PUT /posts/{id}                 │
│ [ ] DELETE /posts/{id}              │
│ [✓] GET /users                      │
└─────────────────────────────────────┘

Space: toggle | a: all | n: none | Enter: run
```

**Use Cases:**
- Test specific features
- Debug individual endpoints
- Faster iteration during development

---

## Advanced Features

### 1. Authentication

The TUI supports 3 authentication methods, configurable via Settings (option 6).

#### Bearer Token Authentication

**When to Use**: JWT tokens, OAuth 2.0

**Configuration:**
1. Settings → Auth Type: `bearer`
2. Auth Token: `your-jwt-token`
3. Save

**Result**: Adds `Authorization: Bearer <token>` header to all requests

#### API Key Authentication

**When to Use**: API key-based services

**Header-Based:**
1. Settings → Auth Type: `apikey`
2. Auth Token/Key: `your-api-key`
3. API Key Header Name: `X-API-Key` (or custom)
4. Save

**Query Parameter-Based:**
1. Settings → Auth Type: `apikey`
2. Auth Token/Key: `your-api-key`
3. API Key Query Param: `api_key` (or custom)
4. Save

#### Basic Authentication

**When to Use**: HTTP Basic Auth

**Configuration:**
1. Settings → Auth Type: `basic`
2. Basic Auth Username: `your-username`
3. Basic Auth Password: `your-password`
4. Save

**Result**: Adds `Authorization: Basic <base64-credentials>` header

### 2. Verbose Mode

**Purpose**: See full HTTP request/response details

**How to Enable:**
- Main menu: Press **'v'** (shows in status bar)
- Test results: Press **'v'** during tests

**What You Get:**
- Full request headers
- Request body
- Response headers
- Response body
- Timing breakdown
- Detailed error messages

**Viewing Logs:**
1. Enable verbose mode (**'v'**)
2. Run tests
3. Press **'l'** on any result to view full details

**Log Display:**
```
┌────────────────────────────────────────┐
│     Detailed Log: GET /posts/1         │
├────────────────────────────────────────┤
│ Request:                               │
│   Method: GET                          │
│   URL: https://api.example.com/posts/1 │
│   Headers:                             │
│     Accept: application/json           │
│     User-Agent: OpenAPI-TUI/1.0        │
│                                        │
│ Response:                              │
│   Status: 200 OK                       │
│   Duration: 125ms                      │
│   Headers:                             │
│     Content-Type: application/json     │
│     Cache-Control: max-age=3600        │
│   Body:                                │
│     {"id": 1, "title": "Post 1"...}    │
└────────────────────────────────────────┘
```

### 3. Response Filtering

**Purpose**: Find specific results quickly in large test runs

**How to Use:**
1. Run tests (option 2)
2. Press **'f'** to enter filter mode
3. Type your filter text
4. Results update in real-time

**Filter Matches:**
- Status codes (e.g., `200`, `404`, `5xx`)
- HTTP methods (e.g., `GET`, `POST`)
- Endpoints (e.g., `/users`, `/posts/1`)
- Keywords in messages

**Examples:**
```
Filter: "200"     → Shows only 200 OK responses
Filter: "GET"     → Shows only GET requests
Filter: "/users"  → Shows only /users endpoints
Filter: "error"   → Shows results with "error" in message
Filter: "4"       → Shows 4xx status codes
```

### 4. Parallel Testing

**Purpose**: Run tests faster by testing multiple endpoints concurrently

**Configuration:**
1. Settings → Worker Count: `1-50` (default: auto-detect CPU cores)
2. Save

**How It Works:**
- Worker pool executes tests concurrently
- Auto-detects CPU cores (default)
- Limits concurrent requests to prevent overwhelming servers
- Results display as they complete

**Performance Impact:**
- Sequential (1 worker): 50 endpoints ≈ 15 seconds
- Parallel (4 workers): 50 endpoints ≈ 5 seconds (3x faster)
- Parallel (8 workers): 50 endpoints ≈ 3 seconds (5x faster)

**Best Practices:**
- Use auto-detect (default) for most cases
- Reduce workers for rate-limited APIs
- Increase workers for fast, reliable APIs

### 5. Retry Logic

**Purpose**: Handle transient network errors automatically

**Configuration:**
1. Settings → Max Retries: `0-10` (default: 3)
2. Save

**How It Works:**
- Retries only on network errors (not HTTP errors)
- Exponential backoff: 1s, 2s, 4s, 8s...
- Max backoff: 10 seconds
- Shows retry attempts in verbose mode

**Retryable Errors:**
- Connection refused
- Timeout
- DNS lookup failed
- Connection reset
- Network unreachable

**Non-Retryable:**
- 4xx HTTP errors (client errors)
- 5xx HTTP errors (server errors)
- Successful responses (2xx, 3xx)

**Visual Feedback:**
```
  ↻ Retry attempt 1/3 after 1000ms...
  ↻ Retry attempt 2/3 after 2000ms...
✓ GET /posts - 200 OK (succeeded on retry 2)
```

### 6. Test History

**Purpose**: Track test runs over time and replay previous tests

**How to Use:**
1. Run some tests (option 2)
2. Press **'r'** to view history
3. Navigate with **↑ / ↓**
4. Press **Enter** to replay a test

**History Display:**
```
┌──────────────────────────────────────────┐
│         📜 Test Run History              │
├──────────────────────────────────────────┤
│ 2025-11-02 18:30:45                      │
│   ✅ 12/12 passed (100%)                 │
│   🕒 Duration: 1.2s                      │
│   📍 https://api.example.com             │
│                                          │
│ 2025-11-02 17:15:22                      │
│   ⚠️  10/12 passed (83%)                 │
│   🕒 Duration: 1.5s                      │
│   📍 https://staging.api.example.com     │
└──────────────────────────────────────────┘
```

**Features:**
- Stores last 50 test runs
- Persistent storage in `~/.config/openapi-tui/history.json`
- Full test results saved
- One-click replay
- Track API health trends

---

## Configuration

### Configuration File

Location: `~/.config/openapi-tui/config.yaml`

**Example Configuration:**
```yaml
# General Settings
base_url: https://api.example.com
spec_file: api-spec.yaml
timeout: 15000
verbose: false

# Authentication
auth_type: bearer
auth_token: your-jwt-token
api_key_header: X-API-Key
api_key_query: ""
basic_username: ""
basic_password: ""

# Performance
worker_count: 4
max_retries: 3
```

### Configuration Editor (Recommended)

**Access:**
1. Main menu → **6. Settings**
2. Navigate with **Tab** / **Shift+Tab**
3. Edit fields
4. Press **Enter** to save

**Fields:**

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| Base URL | String | Default API base URL | `https://api.example.com` |
| Spec File Path | String | Default spec file | `./api-spec.yaml` |
| Timeout (ms) | Number | Request timeout | `15000` (15 seconds) |
| Auth Type | Enum | Authentication method | `none`, `bearer`, `apikey`, `basic` |
| Auth Token/Key | String | Token or API key | `your-token-123` |
| API Key Header | String | Header name for API key | `X-API-Key` |
| API Key Query | String | Query param for API key | `api_key` |
| Basic Username | String | Basic auth username | `admin` |
| Basic Password | String | Basic auth password | `secret123` |
| Worker Count | Number | Parallel workers (1-50) | `4` |
| Max Retries | Number | Retry attempts (0-10) | `3` |
| Verbose Mode | Boolean | Enable verbose logging | `true` / `false` |

**Auto-Save:**
- Configuration automatically saved to `~/.config/openapi-tui/config.yaml`
- Loaded on next TUI launch
- Overridden by manual input during testing

---

## Export & Reporting

### JSON Export

**How to Use:**
1. Run tests
2. Press **'e'** from results screen

**Filename**: `openapi-test-results-YYYY-MM-DD-HHMMSS.json`

**Format:**
```json
{
  "timestamp": "2025-11-02T18:30:00Z",
  "spec_path": "api-spec.yaml",
  "base_url": "https://api.example.com",
  "total_tests": 12,
  "passed": 12,
  "failed": 0,
  "results": [
    {
      "method": "GET",
      "endpoint": "/posts",
      "status": 200,
      "success": true,
      "message": "OK",
      "duration": 125,
      "timestamp": "2025-11-02T18:30:01Z"
    }
  ],
  "statistics": {
    "pass_rate": 100.0,
    "avg_duration": 100.5,
    "total_duration": 1206,
    "fastest": 87,
    "slowest": 156
  }
}
```

**Use Cases:**
- Programmatic analysis
- CI/CD integration
- Custom reporting scripts
- Trend analysis

### HTML Export

**How to Use:**
1. Run tests
2. Press **'h'** from results screen

**Filename**: `openapi-test-results-YYYY-MM-DD-HHMMSS.html`

**Features:**
- Professional web report
- Embedded CSS (self-contained)
- Statistics dashboard
- Color-coded results
- Responsive design
- Print-friendly

**Sections:**
1. Header with API title and timestamp
2. Summary cards (total, passed, failed, success rate)
3. Statistics (timing data)
4. Results table (sortable, color-coded)

**Use Cases:**
- Share with stakeholders
- Email reports
- Documentation
- Archive test results

### JUnit XML Export

**How to Use:**
1. Run tests
2. Press **'j'** from results screen

**Filename**: `openapi-test-results-YYYY-MM-DD-HHMMSS.xml`

**Format:**
```xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuites name="API Tests" tests="12" failures="0" errors="0" time="1.206">
  <testsuite name="API Tests" tests="12" failures="0" errors="0" time="1.206">
    <testcase name="GET /posts" classname="APITest.GET" time="0.125"/>
    <testcase name="POST /posts" classname="APITest.POST" time="0.156"/>
  </testsuite>
</testsuites>
```

**CI/CD Integration:**

**Jenkins:**
```groovy
post {
  always {
    junit '**/openapi-test-results-*.xml'
  }
}
```

**GitLab CI:**
```yaml
artifacts:
  reports:
    junit: openapi-test-results-*.xml
```

**GitHub Actions:**
```yaml
- name: Publish Test Results
  uses: EnricoMi/publish-unit-test-result-action@v2
  with:
    files: openapi-test-results-*.xml
```

---

## Real-World Workflows

### Workflow 1: API Development

**Scenario**: Developing new API endpoints

**Steps:**
1. Write OpenAPI spec for new endpoint
2. **Validate** spec (option 1)
3. Fix any validation errors
4. Implement the endpoint in your code
5. **Test** endpoints (option 2) with verbose mode
6. Review detailed logs (**'l'** key)
7. Fix any issues
8. **Export** results (**'h'** for HTML report)
9. Commit spec and code

**Benefits:**
- Catch spec errors early
- Visual feedback on API behavior
- Documentation alongside code

### Workflow 2: API Debugging

**Scenario**: Investigating failing endpoint

**Steps:**
1. **Custom Request** (option 3)
2. Enter failing endpoint details
3. Execute request
4. Review response (status, headers, body)
5. Modify request (headers, body)
6. Re-execute until issue found
7. Fix backend code
8. **Test** all endpoints (option 2) to verify fix

**Benefits:**
- Quick iteration
- Full request/response visibility
- No need for curl or Postman

### Workflow 3: Regression Testing

**Scenario**: Ensure API hasn't broken after changes

**Steps:**
1. Configure **Settings** (option 6) with production URL
2. **Test** all endpoints (option 2)
3. **Filter** results (**'f'**) to find failures
4. Review detailed logs for failed tests
5. **Export** JUnit XML (**'j'**) for CI/CD
6. Fix regressions
7. **Replay** test from history (**'r'**) to verify fix

**Benefits:**
- Quick regression detection
- Historical comparison
- CI/CD integration

### Workflow 4: Performance Monitoring

**Scenario**: Track API performance over time

**Steps:**
1. Enable **Parallel Testing** (Settings → Worker Count: 8)
2. **Test** endpoints (option 2)
3. Note timing statistics (avg, P50, P95)
4. **Export** JSON (**'e'**) with timestamp
5. Run regularly (daily/weekly)
6. Compare statistics across runs
7. View **History** (**'r'**) for trends

**Benefits:**
- Performance baseline
- Detect degradation early
- Historical data for analysis

### Workflow 5: Multi-Environment Testing

**Scenario**: Test staging vs production

**Staging Test:**
1. Settings → Base URL: `https://staging.api.example.com`
2. Test endpoints
3. Export HTML: `staging-report.html`

**Production Test:**
1. Settings → Base URL: `https://api.example.com`
2. Test endpoints
3. Export HTML: `production-report.html`

**Compare:**
- Open both HTML reports side-by-side
- Compare pass rates
- Compare response times
- Identify environment-specific issues

**Benefits:**
- Environment comparison
- Pre-production validation
- Performance comparison

---

## Troubleshooting

### Issue 1: Spec Validation Fails

**Symptoms**: Validation errors on valid-looking spec

**Solutions:**
1. Check OpenAPI version (must be 3.0.x or 3.1.x)
2. Ensure all required fields present (`info`, `paths`)
3. Validate JSON/YAML syntax externally
4. Check $ref references resolve correctly
5. Use verbose mode for detailed error messages

**Common Errors:**
```
✗ Missing required field: info.version
💡 Add: version: "1.0.0" under info object

✗ Invalid path: users (must start with /)
💡 Change to: /users
```

### Issue 2: Connection Errors

**Symptoms**: "Connection refused" or "Timeout" errors

**Solutions:**
1. Verify API server is running: `curl https://api.example.com/health`
2. Check base URL is correct (no trailing slash)
3. Increase timeout in Settings (default: 10000ms)
4. Check firewall/network settings
5. Enable retry logic (Settings → Max Retries: 3)

**Debug Steps:**
1. Settings → Verbose: true
2. Run test
3. Press 'l' to view full error details
4. Check network connectivity
5. Verify API server logs

### Issue 3: Authentication Failures

**Symptoms**: 401 Unauthorized errors

**Solutions:**
1. Verify auth type is correct (bearer/apikey/basic)
2. Check token/credentials are valid
3. For API key: verify header name or query param matches API expectations
4. Use Custom Request (option 3) to test auth manually
5. Check token expiration

**Testing Auth:**
```
Settings:
  Auth Type: bearer
  Auth Token: your-token

Custom Request:
  Method: GET
  URL: https://api.example.com/protected
  Execute and verify Authorization header is sent
```

### Issue 4: Slow Performance

**Symptoms**: Tests take very long to complete

**Solutions:**
1. Enable parallel testing (Settings → Worker Count: 4-8)
2. Filter to specific endpoints (option 4: Select Endpoints)
3. Increase timeout if API is genuinely slow
4. Check network latency: `ping api.example.com`
5. Use retry sparingly (reduces retry delay)

**Performance Tuning:**
```
Fast APIs: Worker Count = 8-10
Normal APIs: Worker Count = 4-5 (default)
Slow/Rate-Limited: Worker Count = 1-2
```

### Issue 5: Export Failures

**Symptoms**: Export files not created

**Solutions:**
1. Check write permissions in current directory
2. Ensure disk space available
3. Try exporting to specific directory: `/tmp/results.json`
4. Check for special characters in filenames
5. Verify export actually completed (success message)

**Debug:**
```bash
# Check permissions
ls -la .

# Export to known-good directory
cd /tmp
openapi-tui
# Run tests and export
```

---

## Best Practices

### 1. Use Configuration Files

**Don't:** Re-enter settings every time

**Do:** Configure once in Settings (option 6)
- Saves to `~/.config/openapi-tui/config.yaml`
- Auto-loaded on next launch
- Override as needed during testing

### 2. Enable Verbose Mode for Development

**Don't:** Debug with minimal output

**Do:** Enable verbose mode (press 'v')
- See full request/response details
- Press 'l' to view logs
- Understand API behavior
- Catch issues early

### 3. Use Endpoint Selection for Large APIs

**Don't:** Test all 100 endpoints every time

**Do:** Select specific endpoints (option 4)
- Faster iteration
- Focus on changes
- Reduce API load
- Save time

### 4. Export Results Regularly

**Don't:** Lose test results after closing TUI

**Do:** Export in multiple formats
- JSON for analysis
- HTML for sharing
- JUnit for CI/CD
- Keep historical records

### 5. Use Test History for Comparisons

**Don't:** Manually track changes

**Do:** Use history (press 'r')
- Compare current vs previous runs
- Track API health trends
- Replay specific test runs
- Identify regressions

### 6. Configure Parallel Testing Wisely

**Don't:** Use maximum workers always

**Do:** Tune based on API characteristics
- Fast APIs: 8-10 workers
- Normal APIs: 4-5 workers (default)
- Rate-limited: 1-2 workers
- Respect API limits

### 7. Use Retry for Flaky Networks

**Don't:** Manually re-run failed tests

**Do:** Configure retries (Settings → Max Retries: 3)
- Handles transient errors
- Exponential backoff
- Improves reliability
- Don't over-retry (max 3-5)

### 8. Filter Results for Large Test Suites

**Don't:** Scroll through 100 results manually

**Do:** Use filter mode (press 'f')
- Find failures: filter "error" or "4"
- Find specific methods: filter "POST"
- Find endpoints: filter "/users"
- Quick analysis

### 9. Use Custom Requests for Debugging

**Don't:** Use external tools for manual testing

**Do:** Use Custom Request (option 3)
- Test edge cases
- Debug specific requests
- Try different headers/bodies
- All in one place

### 10. Validate Before Testing

**Don't:** Test against invalid specs

**Do:** Always validate first (option 1)
- Catch spec errors early
- Save time on broken tests
- Ensure spec quality
- Better API documentation

---

## Summary

The OpenAPI TUI provides a comprehensive, interactive toolkit for API development and testing. Key takeaways:

- ✅ **Easy to use** - Intuitive keyboard navigation
- ✅ **Visual feedback** - Styled tables, colors, real-time updates
- ✅ **Flexible authentication** - Bearer, API Key, Basic
- ✅ **Multiple exports** - JSON, HTML, JUnit XML
- ✅ **Developer friendly** - Verbose mode, filtering, history
- ✅ **High performance** - Parallel testing, retry logic
- ✅ **Production ready** - 453 tests, 100% core coverage

### Quick Reference

```
Main Menu:
  1. Validate Spec
  2. Test Endpoints  ← Most used
  3. Custom Request  ← Debugging
  4. Select Endpoints
  5. View History    ← Track trends
  6. Settings        ← Configure once
  7. Help
  8. Quit

During Tests:
  v - Verbose mode
  f - Filter results
  l - View logs
  r - View history
  e - Export JSON
  h - Export HTML
  j - Export JUnit XML
```

For more details, see:
- [README.md](../README.md) - Feature overview
- [ARCHITECTURE.md](ARCHITECTURE.md) - System design
- [PROGRESS.md](PROGRESS.md) - Development status
- [TESTING-GUIDE.md](TESTING-GUIDE.md) - Testing procedures

---

**Made with ❤️ for the OpenAPI community**

💡 **Pro Tip:** Combine with the OpenAPI CLI for the complete testing experience!

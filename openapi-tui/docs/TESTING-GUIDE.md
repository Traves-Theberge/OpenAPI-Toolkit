# Step-by-Step Testing Guide for OpenAPI Toolkit

This guide walks you through testing both the **TypeScript CLI** and **Go TUI** applications using the JSONPlaceholder API.

---

## Prerequisites

### Required
- Node.js (v18+) and npm
- Go (v1.21+)
- A terminal with color support
- Internet connection

### Test API Details
- **Base URL:** `https://jsonplaceholder.typicode.com`
- **OpenAPI Spec:** `jsonplaceholder-spec.yaml` (included in this repo)
- **Available Endpoints:**
  - `/posts` - 100 blog posts
  - `/users` - 10 users
  - `/comments` - 500 comments
  - `/todos` - 200 todos
  - `/albums` - 100 albums
  - `/photos` - 5000 photos

---

# Part 1: Testing openapi-cli (TypeScript CLI)

## Step 1: Install Dependencies

```bash
cd /home/traves/Development/OpenAPI-Toolkit/openapi-cli
npm install
```

**Expected output:**
```
added XX packages in Xs
```

## Step 2: Build the Project

```bash
npm run build
```

**Expected output:**
```
> openapi-cli@1.0.0 build
> tsc

✓ Compiled successfully
```

## Step 3: Link for Global Use (Optional)

```bash
npm link
```

**Expected output:**
```
added 1 package
```

Now you can use `openapi-test` globally.

## Step 4: Validate the OpenAPI Spec

```bash
# From the openapi-cli directory
npx ts-node src/cli.ts validate ../jsonplaceholder-spec.yaml

# OR if you ran npm link:
openapi-test validate ../jsonplaceholder-spec.yaml
```

**Expected output:**
```
✅ Validation successful!
OpenAPI specification is valid.
```

**What this tests:**
- YAML parsing works
- OpenAPI structure is correct
- All `$ref` references resolve
- Schema definitions are valid

## Step 5: Test API Endpoints

```bash
# Basic test against JSONPlaceholder
npx ts-node src/cli.ts test ../jsonplaceholder-spec.yaml https://jsonplaceholder.typicode.com

# OR with npm link:
openapi-test test ../jsonplaceholder-spec.yaml https://jsonplaceholder.typicode.com
```

**Expected output:**
```
Testing API endpoints...
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ GET /posts → 200 OK
✅ GET /posts/1 → 200 OK
✅ GET /users → 200 OK
✅ GET /users/1 → 200 OK
✅ GET /comments → 200 OK
✅ GET /comments/1 → 200 OK
✅ GET /todos → 200 OK
✅ GET /todos/1 → 200 OK
✅ GET /albums → 200 OK
✅ GET /albums/1 → 200 OK
✅ GET /photos → 200 OK
✅ GET /photos/1 → 200 OK

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Results: 12/12 passed ✅
```

**What this tests:**
- HTTP client works
- URL construction is correct
- All GET endpoints are accessible
- Response parsing works
- Status code checking works

## Step 6: Test Individual Endpoints

```bash
# Test only posts endpoints
openapi-test test ../jsonplaceholder-spec.yaml https://jsonplaceholder.typicode.com --filter=/posts

# Test with verbose output
openapi-test test ../jsonplaceholder-spec.yaml https://jsonplaceholder.typicode.com --verbose
```

## Step 7: Run Unit Tests

```bash
npm test
```

**Expected output:**
```
PASS src/__tests__/commands/validate.test.ts
  ✓ validates correct OpenAPI spec (XXms)
  ✓ detects invalid spec (XXms)

Test Suites: 1 passed, 1 total
Tests:       2 passed, 2 total
```

## Troubleshooting openapi-cli

### Issue: "Cannot find module"
```bash
npm install
npm run build
```

### Issue: "Validation failed"
Check the spec file exists:
```bash
ls -la ../jsonplaceholder-spec.yaml
```

### Issue: "Network error"
Test direct API access:
```bash
curl https://jsonplaceholder.typicode.com/posts
```

---

# Part 2: Testing openapi-cli-tui (Go TUI)

## Step 1: Install Dependencies

```bash
cd /home/traves/Development/OpenAPI-Toolkit/openapi-cli-tui
go mod tidy
```

**Expected output:**
```
go: downloading github.com/charmbracelet/bubbletea vX.XX.X
go: downloading github.com/charmbracelet/lipgloss vX.XX.X
...
```

## Step 2: Build the Application

```bash
go build -o openapi-tui .
```

**Expected output:**
- No output = success
- Creates `openapi-tui` binary in current directory

Verify:
```bash
ls -lh openapi-tui
```

**Should see:**
```
-rwxr-xr-x 1 user user 8.5M Nov 1 12:00 openapi-tui
```

## Step 3: Run the TUI

```bash
./openapi-tui
```

**Expected screen:**
```
┌───────────────────────────────────────┐
│                                       │
│         🚀 OpenAPI CLI TUI            │
│                                       │
│   ┌─────────────────────────────┐    │
│   │ 📋 Validate OpenAPI Spec    │    │
│   │ 🧪 Test API Endpoints       │    │
│   │ ❓ Help                     │    │
│   │ 🚪 Quit                     │    │
│   └─────────────────────────────┘    │
│                                       │
│   ↑/↓ or j/k to navigate             │
│   Enter to select • q to quit        │
│                                       │
└───────────────────────────────────────┘
```

## Step 4: Test Validation Feature

### Actions:
1. **Press ↓ or j** to highlight "Validate OpenAPI Spec"
2. **Press Enter**
3. **Type the path:** `../jsonplaceholder-spec.yaml`
4. **Press Enter**

### Expected Output:
```
┌────────────────────────────────────────┐
│      Validate OpenAPI Spec             │
├────────────────────────────────────────┤
│                                        │
│  Enter path to OpenAPI spec file:     │
│  > ../jsonplaceholder-spec.yaml        │
│                                        │
│  ✅ Validation successful!             │
│                                        │
│  Spec is valid OpenAPI 3.0.3          │
│  Found 12 paths                        │
│  Found 6 schemas                       │
│                                        │
│  Press Enter or Esc to continue        │
│                                        │
└────────────────────────────────────────┘
```

5. **Press Enter or Esc** to return to menu

## Step 5: Test API Testing Feature

### Actions:
1. From main menu, select "Test API Endpoints"
2. **Step 1 - Spec File:**
   - Type: `../jsonplaceholder-spec.yaml`
   - Press Enter

3. **Step 2 - Base URL:**
   - Type: `https://jsonplaceholder.typicode.com`
   - Press Enter

4. **Step 3 - Testing (automatic):**
   - Watch the spinner animate
   - Wait 5-10 seconds

### Expected Output:
```
┌──────────────────────────────────────────────────────────┐
│              Test API Endpoints - Results                │
├──────────────────────────────────────────────────────────┤
│                                                          │
│  ┌────────┬─────────────────┬─────────┬──────────────┐ │
│  │ Method │ Endpoint        │ Status  │ Message      │ │
│  ├────────┼─────────────────┼─────────┼──────────────┤ │
│  │ GET    │ /posts          │ ✅ 200  │ Success      │ │
│  │ GET    │ /posts/{id}     │ ✅ 200  │ Success      │ │
│  │ POST   │ /posts          │ ✅ 201  │ Created      │ │
│  │ GET    │ /users          │ ✅ 200  │ Success      │ │
│  │ GET    │ /users/{id}     │ ✅ 200  │ Success      │ │
│  │ GET    │ /comments       │ ✅ 200  │ Success      │ │
│  │ GET    │ /comments/{id}  │ ✅ 200  │ Success      │ │
│  │ GET    │ /todos          │ ✅ 200  │ Success      │ │
│  │ GET    │ /todos/{id}     │ ✅ 200  │ Success      │ │
│  │ GET    │ /albums         │ ✅ 200  │ Success      │ │
│  │ GET    │ /albums/{id}    │ ✅ 200  │ Success      │ │
│  │ GET    │ /photos         │ ✅ 200  │ Success      │ │
│  │ GET    │ /photos/{id}    │ ✅ 200  │ Success      │ │
│  └────────┴─────────────────┴─────────┴──────────────┘ │
│                                                          │
│  Summary: 13/13 endpoints passed ✅                      │
│                                                          │
│  Press Enter or Esc to return to menu                   │
│                                                          │
└──────────────────────────────────────────────────────────┘
```

5. **Press Enter or Esc** to return to menu

## Step 6: Test Help Screen

1. From main menu, **press h or ?**
2. View help documentation
3. **Press q, Esc, h, or ?** to return to menu

## Step 7: Exit Application

1. From main menu, select "Quit"
2. **OR press q or Ctrl+C** at any time

## Troubleshooting openapi-cli-tui

### Issue: "go: command not found"
Install Go:
```bash
# Ubuntu/Debian
sudo apt install golang-go

# Or download from: https://go.dev/dl/
```

### Issue: Build fails with missing dependencies
```bash
go mod download
go mod tidy
go build -o openapi-tui .
```

### Issue: Terminal display is garbled
- Ensure your terminal supports colors (most do)
- Try a different terminal (iTerm2, Windows Terminal, Alacritty)
- Resize terminal to at least 80x24

### Issue: Can't find spec file
Use absolute path:
```bash
# From TUI input
/home/traves/Development/OpenAPI-Toolkit/jsonplaceholder-spec.yaml
```

### Issue: Network timeout
Check internet connection:
```bash
curl https://jsonplaceholder.typicode.com/posts
```

---

# Part 3: Advanced Testing Scenarios

## Scenario 1: Test Against Local API

If you have a local API running on `http://localhost:3000`:

```bash
# CLI
openapi-test test ./your-local-spec.yaml http://localhost:3000

# TUI
./openapi-tui
# Then enter: http://localhost:3000 when prompted
```

## Scenario 2: Test Invalid Spec

Create an invalid spec to test error handling:

```bash
echo "invalid: yaml: content" > invalid-spec.yaml

# CLI
openapi-test validate ./invalid-spec.yaml
# Expected: ❌ Validation failed with error message

# TUI
./openapi-tui → Validate → ./invalid-spec.yaml
# Expected: Red error message displayed
```

## Scenario 3: Test Unreachable API

```bash
# CLI
openapi-test test ../jsonplaceholder-spec.yaml http://localhost:9999
# Expected: Connection refused errors

# TUI
./openapi-tui → Test → http://localhost:9999
# Expected: Error messages in results table
```

## Scenario 4: Performance Testing

Time how long it takes to test all endpoints:

```bash
time openapi-test test ../jsonplaceholder-spec.yaml https://jsonplaceholder.typicode.com
```

Expected: ~5-10 seconds for 12 endpoints

---

# Part 4: Verification Checklist

## ✅ openapi-cli Tests

- [ ] Installs without errors
- [ ] Builds successfully
- [ ] Validates correct spec
- [ ] Detects invalid spec
- [ ] Tests all endpoints
- [ ] Displays results clearly
- [ ] Handles network errors gracefully
- [ ] Unit tests pass

## ✅ openapi-cli-tui Tests

- [ ] Builds successfully
- [ ] Launches TUI without errors
- [ ] Menu navigation works (arrows/jk)
- [ ] Validation feature works
- [ ] Shows validation results
- [ ] Testing feature works
- [ ] Shows spinner during testing
- [ ] Displays results table
- [ ] Help screen accessible
- [ ] Can exit cleanly

## ✅ Integration Tests

- [ ] Both tools validate same spec successfully
- [ ] Both tools test same API successfully
- [ ] Results are consistent between tools
- [ ] Both handle errors appropriately

---

# Quick Reference Commands

## TypeScript CLI

```bash
# Setup
cd openapi-cli && npm install && npm run build

# Validate
openapi-test validate ../jsonplaceholder-spec.yaml

# Test
openapi-test test ../jsonplaceholder-spec.yaml https://jsonplaceholder.typicode.com

# Run tests
npm test
```

## Go TUI

```bash
# Setup
cd openapi-cli-tui && go mod tidy && go build -o openapi-tui .

# Run
./openapi-tui

# Navigation
↑/↓ or j/k - Navigate
Enter - Select
q/Esc - Back/Quit
h/? - Help
```

---

# Success Criteria

Both applications should:
1. ✅ **Validate** the JSONPlaceholder spec without errors
2. ✅ **Test** all 12+ endpoints successfully
3. ✅ **Display** clear, formatted output
4. ✅ **Handle errors** gracefully (invalid paths, network issues)
5. ✅ **Complete testing** in under 15 seconds

---

# Next Steps

Once both tools are working:
1. Test against your own OpenAPI specs
2. Add more endpoints to the test spec
3. Implement authentication testing
4. Add response validation (schema checking)
5. Create CI/CD pipeline for automated testing
6. Add support for POST/PUT/DELETE operations

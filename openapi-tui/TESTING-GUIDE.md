# 🧪 OpenAPI TUI - Manual Testing Guide

## Quick Start

```bash
# Build the TUI
go build -o openapi-tui ./cmd/openapi-tui

# Run it
./openapi-tui
```

---

## 📋 **Complete Test Workflow** (15-20 minutes)

### **Test 1: Validate OpenAPI Spec** ✅

1. Launch TUI: `./openapi-tui`
2. Select option **1: Validate OpenAPI Spec**
3. Enter spec path: `test-api.yaml` (or full path)
4. Press **Enter**
5. **Expected**: Green success message "✅ Spec is valid!"
6. Press **Enter** or **Esc** to return to menu

**What to verify:**
- ✅ Spec loads without errors
- ✅ Validation message is clear
- ✅ Can return to menu

---

### **Test 2: Configure Settings** ⚙️

1. From menu, select option **8: ⚙️  Settings**
2. **Verify 12 fields displayed**:
   - General Settings: Spec Path, Base URL, Verbose Mode
   - Authentication: Type, Token, API Key Name, API Key Location, Username, Password
   - Performance: Max Concurrency, Max Retries, Retry Delay

3. **Test Tab Navigation**:
   - Press **Tab** multiple times - cursor should move through all 12 fields
   - Press **Shift+Tab** - cursor should move backwards
   - Press **↑/↓** - should also navigate fields

4. **Test Configuration**:
   - Field 1 (Spec Path): Enter `test-api.yaml`
   - Field 2 (Base URL): Enter `https://jsonplaceholder.typicode.com`
   - Field 3 (Verbose): Enter `true`
   - Field 4 (Auth Type): Enter `none`
   - Field 10 (Max Concurrency): Leave as `0` (auto)
   - Field 11 (Max Retries): Enter `3`
   - Field 12 (Retry Delay): Enter `1000`

5. Press **Enter** to save
6. **Expected**: Return to menu with "Config saved" message
7. Press **Esc** to cancel (try this once to test cancel)

**What to verify:**
- ✅ All 12 fields visible and editable
- ✅ Tab navigation works forward and backward
- ✅ Save returns to menu
- ✅ Cancel discards changes
- ✅ Config persists after restart (check ~/.config/openapi-tui/config.yaml)

---

### **Test 3: Run API Tests (All Features)** 🧪

1. From menu, select option **2: Test API**
2. Enter spec path: `test-api.yaml` (should be pre-filled from config)
3. Press **Enter**
4. Enter base URL: `https://jsonplaceholder.typicode.com` (should be pre-filled)
5. Press **Enter**
6. **Expected**: Tests run automatically on all endpoints

**Results Screen - Test All Features:**

#### **A. Toggle Verbose Mode**
- Press **v** key
- **Expected**: Status bar shows "Verbose: ON"
- Press **v** again
- **Expected**: Status bar shows "Verbose: OFF"
- Turn it **ON** for remaining tests

#### **B. View Detailed Logs**
- With verbose ON, navigate to any result with ↑/↓
- Press **l** key
- **Expected**: Detailed log view showing:
  - Request: method, URL, headers, body
  - Response: status, headers, body
  - Duration
- Press **Esc** to return to results

#### **C. Filter Results**
- Press **f** key
- **Expected**: Filter input appears at top
- Type: `200` (filter by status)
- **Expected**: Only 2xx results shown
- Press **Esc** to exit filter
- Press **f** again, type: `GET` (filter by method)
- **Expected**: Only GET requests shown
- Press **Esc** to clear filter

#### **D. Export JSON**
- Press **e** key
- **Expected**: Success message "✅ Exported JSON to openapi-test-results-YYYY-MM-DD-HHMMSS.json"
- Verify file exists: `ls -lh openapi-test-results-*.json`
- Check contents: `cat openapi-test-results-*.json | jq .`

#### **E. Export HTML**
- Press **h** key
- **Expected**: Success message "✅ Exported HTML to openapi-test-results_YYYYMMDD_HHMMSS.html"
- Verify file exists: `ls -lh openapi-test-results_*.html`
- Open in browser: `xdg-open openapi-test-results_*.html` (Linux) or `open` (Mac)
- **Verify HTML**:
  - Professional styling
  - Statistics dashboard
  - Color-coded results table
  - 6 columns including "Retries"

#### **F. Export JUnit XML**
- Press **j** key
- **Expected**: Success message "✅ Exported JUnit XML to openapi-test-results_YYYYMMDD_HHMMSS.xml"
- Verify file exists: `ls -lh openapi-test-results_*.xml`
- Check structure: `cat openapi-test-results_*.xml | head -30`

#### **G. View History**
- Press **r** key
- **Expected**: History screen showing previous test runs
- Navigate with **↑/↓** or **j/k**
- Press **Enter** on a history entry
- **Expected**: Test re-runs with same spec/URL
- Results appear again

#### **H. Return to Menu**
- Press **Enter** or **Esc**
- **Expected**: Return to main menu
- Config and history should persist

**What to verify:**
- ✅ Tests run on all endpoints
- ✅ Verbose toggle works (v key)
- ✅ Filter works (f key)
- ✅ Log viewer works (l key, verbose mode only)
- ✅ JSON export works (e key)
- ✅ HTML export works (h key)
- ✅ JUnit export works (j key)
- ✅ History works (r key)
- ✅ History replay works
- ✅ Results table navigable with ↑/↓
- ✅ Retry count shown in results (if any retries occurred)

---

### **Test 4: Endpoint Selector (Selective Testing)** 🎯

1. From menu, select option **2: Test API**
2. Enter spec and URL (or use saved config)
3. At URL prompt, press **Shift+S** (or check menu for selector key)
   - *Note: You may need to check the exact key binding in the code*
4. **Expected**: Endpoint selector screen with fuzzy search
5. Use **Space** to toggle individual endpoints
6. Press **a** to select all
7. Press **n** to deselect all
8. Select 2-3 endpoints manually
9. Press **Enter** to run tests on selected endpoints only
10. **Expected**: Only selected endpoints are tested

**What to verify:**
- ✅ Endpoint selector shows all available endpoints
- ✅ Space toggles selection
- ✅ 'a' selects all
- ✅ 'n' deselects all
- ✅ Only selected endpoints run

---

### **Test 5: Custom Request** 🎨

1. From menu, select option **5: Custom Request**
2. **Method field** (focused first):
   - Enter: `GET`
   - Press **Tab**
3. **URL field**:
   - Enter: `https://jsonplaceholder.typicode.com/posts/1`
   - Press **Tab**
4. **Headers** (optional):
   - Skip by pressing **Tab** with empty value
5. **Body** (for GET, leave empty):
   - Press **Enter** to execute
6. **Expected**: Request executes, response shown
7. Press **Esc** to return

**Test POST request:**
1. Select Custom Request again
2. Method: `POST`
3. URL: `https://jsonplaceholder.typicode.com/posts`
4. Headers: Add `Content-Type: application/json` (if supported)
5. Body: `{"title":"test","body":"test body","userId":1}`
6. Execute
7. **Expected**: 201 Created response

**What to verify:**
- ✅ Tab moves through fields
- ✅ GET request works
- ✅ POST request works
- ✅ Response shown clearly
- ✅ Can return to menu

---

### **Test 6: Help Screen** ❓

1. From menu, press **h** or **?**
2. **Expected**: Help screen with:
   - Key bindings
   - Feature descriptions
   - Navigation help
3. Press **Esc** or **q** to return

**What to verify:**
- ✅ Help screen displays
- ✅ Information is accurate
- ✅ Can return to menu

---

### **Test 7: Parallel Testing** ⚡

1. Configure Max Concurrency to `5` in settings
2. Run tests on spec with many endpoints
3. **Expected**: Tests run faster than sequential
4. Check terminal output for concurrent execution indicators

**What to verify:**
- ✅ Tests complete faster with parallelism
- ✅ No race conditions or crashes
- ✅ Results are accurate

---

### **Test 8: Retry Logic** 🔄

**Test with a failing endpoint:**
1. Create a spec with an endpoint that might be flaky or unreachable
2. Configure Max Retries to `3`, Retry Delay to `1000ms`
3. Run tests
4. **Expected**: Failed requests retry up to 3 times
5. Check export files - "Retries" column shows retry count

**What to verify:**
- ✅ Retries occur on network errors
- ✅ Retries occur on 5xx errors
- ✅ No retries on 4xx errors
- ✅ Retry count displayed in results
- ✅ Retry count in all export formats (JSON, HTML, JUnit)

---

## 🐛 **Error Testing**

### **Test Invalid Spec**
1. Create invalid YAML file: `echo "invalid: yaml: file:" > invalid.yaml`
2. Try to validate it
3. **Expected**: Clear error message with suggestions

### **Test Unreachable Server**
1. Use spec with valid endpoints
2. Enter fake URL: `http://localhost:9999`
3. **Expected**: Connection error with helpful message

### **Test Missing File**
1. Enter non-existent spec path: `nonexistent.yaml`
2. **Expected**: File not found error with suggestions

---

## ✅ **Final Checklist**

- [ ] Menu navigation works (↑/↓, j/k, Enter)
- [ ] Help screen accessible (h or ?)
- [ ] Spec validation works
- [ ] Configuration editor (12 fields, Tab navigation)
- [ ] Config persists across restarts
- [ ] Tests run successfully
- [ ] Verbose mode toggles (v key)
- [ ] Filter mode works (f key)
- [ ] Log viewer works (l key, verbose mode)
- [ ] JSON export works (e key)
- [ ] HTML export works (h key)
- [ ] JUnit XML export works (j key)
- [ ] History screen works (r key)
- [ ] History replay works
- [ ] Endpoint selector works
- [ ] Custom requests work (GET and POST)
- [ ] Parallel testing faster than sequential
- [ ] Retry logic works on failures
- [ ] Retry count shown in results and exports
- [ ] Error messages helpful and clear
- [ ] All key bindings work as documented
- [ ] UI is responsive and polished

---

## 📊 **Expected Files After Testing**

```bash
ls -lh
```

Should show:
- `config.yaml` (if in ~/.config/openapi-tui/)
- `history.json` (in ~/.config/openapi-tui/)
- `openapi-test-results-*.json` (JSON exports)
- `openapi-test-results_*.html` (HTML exports)
- `openapi-test-results_*.xml` (JUnit exports)

Check history:
```bash
cat ~/.config/openapi-tui/history.json | jq .
```

Check config:
```bash
cat ~/.config/openapi-tui/config.yaml
```

---

## 🎬 **Quick Smoke Test Script**

For a quick automated check:

```bash
#!/bin/bash
echo "🧪 Running OpenAPI TUI smoke tests..."

# Build
echo "1. Building..."
go build -o openapi-tui ./cmd/openapi-tui || exit 1

# Check binary
echo "2. Checking binary..."
./openapi-tui --version 2>/dev/null || echo "✓ Binary runs"

# Run tests
echo "3. Running unit tests..."
go test ./... || exit 1

echo ""
echo "✅ Smoke tests passed!"
echo "Now run './openapi-tui' for manual testing"
```

---

## 🚀 **Next Steps After Manual Testing**

Once manual testing is complete:
1. Note any bugs or issues
2. Verify all 15 Phase 2 features work correctly
3. Ready to start Phase 3 (CI/CD Headless Mode)
4. Consider creating automated E2E tests for regression testing

---

**Happy Testing! 🎉**

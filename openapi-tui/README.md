# 🚀 OpenAPI CLI TUI

A modern, polished terminal user interface for validating and testing APIs against OpenAPI specifications. Built with Go and the Charmbracelet ecosystem for a fast, friendly CLI experience.

![Demo](https://img.shields.io/badge/Demo-Coming%20Soon-FF6B6B?style=for-the-badge&logo=terminal)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-v0.25.0-7D56F4?style=for-the-badge)

## ✨ Features

### Core Capabilities
- 📋 **Validate OpenAPI Specs** — Ensure your API docs are valid and compliant
- 🧪 **Auto-Test Endpoints** — Automatically test endpoints defined in your spec
- 🎨 **Polished TUI** — Clear, styled terminal interface with useful feedback
- ⚡ **Fast & Lightweight** — Written in Go for performance and simple distribution
- 🔄 **Real-time Feedback** — See validation and test results as they run

### Phase 1 — Critical Foundation (Complete ✅)
- 🔐 **Authentication Support** — Bearer tokens, API keys (header/query), Basic auth
- 📦 **Smart Request Bodies** — Auto-generate realistic JSON from OpenAPI schemas
- ✅ **Response Validation** — Validate status codes and content-types against spec
- 🔗 **Path Parameters** — Automatic substitution of `{id}` placeholders
- 🔍 **Query Parameters** — Auto-generated query strings
- 🧪 **Excellent Test Coverage** — 90.7% average across all packages, 453+ tests

### Phase 2 — Developer Experience (100% Complete ✅🎉)
- 💡 **Enhanced Error Messages** — Actionable suggestions for common issues ✅
- 📊 **Verbose Logging & Display** — Full HTTP details with 'v' toggle, 'l' to view logs ✅
- 💾 **Configuration Support** — Auto-save/load settings from `~/.config/openapi-tui/config.yaml` ✅
- ⚙️ **Configuration Editor** — Form-based UI to edit all settings with validation ✅
- 📤 **Export Results** — JSON, HTML, and JUnit XML export for CI/CD integration ✅
- 🏗️ **Standard Go Layout** — Modular architecture with cmd/ and internal/ packages ✅
- 📊 **Summary Statistics** — Pass rates, timing analysis, performance metrics ✅
- 🔍 **Response Filtering** — Filter results by status, method, endpoint, or keywords ✅
- 📄 **HTML Reports** — Professional web reports with embedded CSS and statistics ✅
- 🤖 **JUnit XML** — CI/CD integration with Jenkins, GitLab CI, GitHub Actions ✅
- 📜 **Request History** — Track, replay, and analyze test runs over time ✅
- ⚡ **Parallel Testing** — Concurrent endpoint testing with worker pools ✅
- 🎨 **Custom Requests** — Execute custom HTTP requests with full control ✅
- 🎯 **Selective Testing** — Interactive multi-select UI for choosing specific endpoints ✅
- 🔄 **Test Retry Logic** — Exponential backoff with configurable retries (0-10, default 3) ✅
- 🧪 **453+ Tests** — Comprehensive test suite with 90.7% average coverage

## 🛠️ Installation

### Prerequisites
- Go 1.21 or later
- A terminal that supports colors (most modern terminals do)

### Install from source
```bash
git clone https://github.com/Traves-Theberge/OpenAPI-Toolkit.git
cd OpenAPI-Toolkit/openapi-tui
go mod tidy
go install .
```

### Or install directly
```bash
go install github.com/Traves-Theberge/OpenAPI-Toolkit/openapi-tui@latest
```

## 🎮 Usage

Run the installed binary (or build and run locally):

```bash
openapi-tui
```

### Navigation & Key Bindings

#### Global Keys
- **↑/↓ or j/k** — Navigate menu/lists
- **Enter** — Select option / Confirm
- **Esc** — Go back / Cancel
- **q** — Quit (from menu or help)
- **h or ?** — Show help screen

#### Menu Screen
- **v** — Toggle verbose mode (shows in status bar)
- **Enter** — Select menu option (0-7)

#### Test Results Screen
- **v** — Toggle verbose logging (enables 'l' key)
- **f** — Toggle filter mode (filter by status/method/endpoint)
- **e** — Export results to JSON
- **h** — Export results to HTML
- **j** — Export results to JUnit XML
- **r** — View test run history
- **l** — View detailed logs (only when verbose mode enabled)
- **↑/↓** — Navigate results table
- **Enter** — Return to menu

#### Filter Mode (when active)
- **Type** — Enter filter text (matches status, method, endpoint)
- **Esc** — Exit filter mode
- **Enter** — Return to menu with filter active

#### History Screen
- **↑/↓ or j/k** — Navigate history entries
- **Enter** — Replay selected test
- **Esc** — Return to results

#### Configuration Editor
- **Tab/Shift+Tab** — Navigate fields (12 fields total)
- **↑/↓** — Also navigate fields
- **Enter** — Save configuration
- **Esc** — Cancel and return to menu

#### Custom Request Screen
- **Tab** — Move through: Method → URL → Headers → Body
- **Enter** — Execute request (when on body field)
- **Esc** — Return to menu

#### Endpoint Selector
- **Space** — Toggle endpoint selection
- **a** — Select all endpoints
- **n** — Deselect all
- **Enter** — Run tests on selected endpoints
- **Esc** — Cancel selection

### Typical Workflow
1. **Validate Spec** 📋
   - Enter path to your OpenAPI YAML/JSON file
   - Get instant validation results and errors

2. **Configure Settings** ⚙️
   - Access from main menu (option 6: Settings)
   - Form-based editor with 10 configurable fields
   - Sections: General Settings / Authentication / Performance
   - Navigate with Tab/Shift+Tab, Enter to save, Esc to cancel
   - Auto-saves to `~/.config/openapi-tui/config.yaml`

3. **Test API** 🧪
   - Provide spec file path
   - Enter base URL (e.g., `https://api.example.com`)
   - Watch automated endpoint testing with live statistics
   - Press **'v'** to toggle verbose logging (full HTTP details)
   - Press **'f'** to filter results by status, method, or keywords
   - Press **'l'** on a result to view detailed logs (request/response headers, bodies, timing)
   - Press **'r'** to view test run history and replay previous tests
   - Press **'e'**, **'h'**, or **'j'** to export results
   - Use parallel testing for faster execution on large specs

### Export & Analysis
After running tests, export results in multiple formats:

**JSON Export** (press **'e'**):
- Filename: `openapi-test-results-YYYY-MM-DD-HHMMSS.json`
- Machine-readable format for CI/CD integration
- Includes metadata, statistics, and full test details
- Contains verbose log data when enabled

**HTML Report** (press **'h'**):
- Professional web report with embedded CSS
- Statistics dashboard with visual indicators
- Color-coded results table
- Perfect for sharing with stakeholders

**JUnit XML** (press **'j'**):
- Standard CI/CD format for Jenkins, GitLab CI, GitHub Actions
- Test suite with proper failure/error distinction
- Timing data and metadata properties
- Automated pipeline integration

**Request History** (press **'r'**):
- View past test runs with timestamps and statistics
- Replay any previous test with one keystroke
- Track API health trends over time
- Persistent storage in `~/.config/openapi-tui/history.json`

## 📚 Documentation

For comprehensive guides and documentation:

- **[Complete User Guide](docs/USER-GUIDE.md)** - End-to-end guide with real-world workflows
- **[Architecture Guide](docs/ARCHITECTURE.md)** - System design and components
- **[Progress Tracking](docs/PROGRESS.md)** - Feature roadmap and status
- **[Testing Guide](docs/TESTING-GUIDE.md)** - Testing procedures

---

## 🏗️ Architecture

Built with industry-standard libraries:

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** — Elm-inspired TUI framework
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** — Styling for terminal apps
- **[Kin OpenAPI](https://github.com/getkin/kin-openapi)** — OpenAPI parsing & validation

## 🚀 Development

```bash
# Clone and setup
git clone https://github.com/Traves-Theberge/OpenAPI-Toolkit.git
cd OpenAPI-Toolkit/openapi-tui

# Install dependencies
go mod tidy

# Run locally
go run .

# Build for production
go build -o openapi-tui .
```

### Tests
Run the comprehensive test suite (453 tests across 8 packages):

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run with verbose output
go test ./... -v

# Run with race detection
go test ./... -race

# Run specific package tests
go test ./internal/validation -v
go test ./internal/export -v
go test ./internal/testing -v
```

**Test Coverage:**
- ✅ **errors** — 100.0% coverage - Enhanced error messages with suggestions
- ✅ **ui** — 96.7% coverage - View rendering, filters, forms (45.3% improvement)
- ✅ **validation** — 94.3% coverage - OpenAPI spec & response validation
- ✅ **export** — 93.0% coverage - JSON/HTML/JUnit export formats
- ✅ **models** — 90.2% coverage - Data structures and history (5.9% improvement)
- ✅ **config** — 88.1% coverage - Configuration save/load, all auth types
- ✅ **testing** — 80.3% coverage - Request generation, auth, parallel, retry logic (10% improvement)
- ✅ **Total:** 453+ tests, 90.7% average coverage, all passing ✅

See [COVERAGE-STATUS.md](COVERAGE-STATUS.md) for detailed coverage analysis and roadmap to 100%.

## 🤝 Contributing

PRs welcome! Please keep contributions professional and well-tested.

1. Fork it
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -am 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request and reference the issue/feature

## 📄 License

MIT - See [LICENSE](LICENSE) file

---

**Made with ❤️ and lots of ☕ by the OpenAPI CLI TUI team**
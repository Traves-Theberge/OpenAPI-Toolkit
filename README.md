# OpenAPI Toolkit

[![TypeScript](https://img.shields.io/badge/TypeScript-5.9.2-blue?logo=typescript)](https://www.typescriptlang.org/)
[![Go](https://img.shields.io/badge/Go-1.23.0-00ADD8?logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/Tests-Passing-success)](.)

A comprehensive toolkit for validating OpenAPI specifications and testing APIs. Available in two flavors: a powerful CLI for automation and a beautiful interactive TUI for development.

> 🚀 **Two tools, one mission:** Make OpenAPI testing effortless

## 🎯 What is this?

OpenAPI Toolkit provides two complementary tools for working with OpenAPI 3.x specifications:

- **[openapi-cli](openapi-cli/)** - TypeScript command-line tool for CI/CD and automation
- **[openapi-tui](openapi-tui/)** - Go terminal UI application for interactive testing

Both tools validate OpenAPI specs and test API endpoints, but serve different use cases and workflows.

## 🚀 Quick Start

### CLI - For Automation & CI/CD

```bash
cd openapi-cli
npm install
npm run build
npm link

# Validate a spec
openapi-test validate spec.yaml

# Test an API
openapi-test test spec.yaml https://api.example.com
```

### TUI - For Interactive Development

```bash
cd openapi-tui
go build

# Run the interactive TUI
./openapi-tui
```

## 📊 Feature Comparison

| Feature | CLI | TUI | Best For |
|---------|-----|-----|----------|
| **Interface** | Command-line | Interactive UI | CLI: Scripts, TUI: Humans |
| **HTTP Methods** | 7 methods | 5 methods | CLI: More complete |
| **Parameters** | Path + Query | Path + Query | Both equal |
| **Request Bodies** | Examples | Advanced generation | TUI: More sophisticated |
| **Response Validation** | Status only | Status + Schema | TUI: More thorough |
| **Output** | Colored text | Styled tables | TUI: More visual |
| **Summary Stats** | Yes | **Yes ✅** | **TUI: Pass rates + timing** |
| **Authentication** | No | Yes (Bearer/API Key/Basic) | TUI: Fully implemented |
| **Error Messages** | Basic | Enhanced with suggestions | TUI: More helpful |
| **Verbose Logging** | No | Yes (toggle with 'v') | TUI: Full HTTP details |
| **Export Formats** | No | **JSON/HTML/JUnit XML ✅** | **TUI: CI/CD ready** |
| **Response Filtering** | No | **Yes (status/method/keyword) ✅** | **TUI: Interactive** |
| **Test History** | No | **Yes (50 runs) ✅** | **TUI: Replay tests** |
| **Configuration** | No | **Persistent YAML ✅** | **TUI: Auto-save** |
| **Help System** | Man pages | Built-in screens | TUI: More accessible |
| **Exit Codes** | Yes (0/1) | N/A | CLI: CI/CD friendly |
| **Use Case** | Automation | Development + CI/CD | Different workflows |
| **Test Coverage** | 3 tests | **94 tests (170+ runs) ✅** | **TUI: Comprehensive** |

## 🛠️ Technologies

### OpenAPI CLI (TypeScript)
- **Language:** TypeScript 5.9.2
- **Runtime:** Node.js
- **HTTP Client:** Axios 1.6.0
- **CLI Framework:** Commander.js 12.0.0
- **YAML Parser:** js-yaml 4.1.0
- **Testing:** Jest 29.7.0

### OpenAPI TUI (Go)
- **Language:** Go 1.23.0
- **TUI Framework:** Bubble Tea 1.3.4
- **Styling:** Lip Gloss 1.1.0
- **Components:** Bubbles 0.21.0
- **OpenAPI Parser:** kin-openapi 0.124.0

## 📖 Documentation

### CLI Documentation
- [CLI README](openapi-cli/README.md) - Complete CLI documentation
- Installation, usage, examples, troubleshooting
- API testing guide
- CI/CD integration examples

### TUI Documentation
- [TUI README](openapi-tui/README.md) - TUI user guide
- [Architecture Guide](openapi-tui/docs/ARCHITECTURE.md) - System design with Mermaid diagrams
- [Testing Guide](openapi-tui/docs/TESTING-GUIDE.md) - Step-by-step testing instructions

## 🎨 Features Showcase

### CLI Example

```bash
$ openapi-test test spec.yaml https://jsonplaceholder.typicode.com

🧪 Testing API: JSONPlaceholder API
📍 Base URL: https://jsonplaceholder.typicode.com

✓ GET     /posts                    - 200 OK
✓ POST    /posts                    - 201 OK
✓ PUT     /posts/1                  - 200 OK
✓ DELETE  /posts/1                  - 200 OK

================================================================================
📊 Summary: 15 passed, 0 failed, 15 total
✓ All tests passed!
```

### TUI Features (Phase 2 - 67% Complete)

```
┌─────────────────────────────────────────────────────────────────┐
│                  📊 Test Results - 15 endpoints                 │
├─────────────────────────────────────────────────────────────────┤
│ Statistics:                                                     │
│   Total: 15 | Passed: 15 | Failed: 0 | Pass Rate: 100%         │
│   Timing: Total 2.5s | Avg 167ms | Fastest 95ms | Slowest 312ms│
├─────────────────────────────────────────────────────────────────┤
│ Method │ Endpoint        │ Status │ Duration │ Message          │
│ GET    │ /posts          │ ✅ 200 │ 125ms    │ OK (validated)   │
│ GET    │ /users          │ ✅ 200 │ 98ms     │ OK (validated)   │
│ POST   │ /posts          │ ✅ 201 │ 156ms    │ Created          │
│ ...                                                              │
├─────────────────────────────────────────────────────────────────┤
│ Press 'f' filter | 'e' JSON | 'h' HTML | 'j' JUnit XML | 'r' history
│ Verbose: ON | Config loaded | Enter to return                  │
└─────────────────────────────────────────────────────────────────┘

✨ New Phase 2 Features:
  • 📊 Summary Statistics - Pass rates, timing analysis
  • 🔍 Response Filtering - Filter by status/method/keywords ('f')
  • 📄 HTML Export - Professional reports ('h')
  • 🤖 JUnit XML - CI/CD integration ('j')
  • 📜 Request History - Track & replay tests ('r')
  • 💾 Config Persistence - Auto-save settings
  • 📊 Verbose Logging - Full HTTP details ('v', 'l')
```

## ⚡ Supported Features

### Validation
- ✅ OpenAPI 3.x version verification
- ✅ Required fields validation (info, paths, operations)
- ✅ Path format checking (must start with `/`)
- ✅ Operation structure validation
- ✅ Response definitions validation
- ✅ Detailed error reporting with JSON paths

### API Testing
- ✅ **HTTP Methods:** GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
- ✅ **Path Parameters:** Automatic `{id}` → `1` replacement
- ✅ **Query Parameters:** Type-aware generation from spec
- ✅ **Request Bodies:** Schema-based generation (TUI) or examples (CLI)
- ✅ **Timeouts:** 10-second configurable timeout
- ✅ **Error Handling:** Connection errors, timeouts, HTTP errors

### Output & Reporting
- ✅ **Colored Output:** Green for success, red for errors
- ✅ **Emoji Indicators:** 🧪 📍 📊 ✓ ✗
- ✅ **Summary Statistics:** Pass rates, timing analysis, performance metrics (TUI)
- ✅ **Interactive Tables:** Sortable results with filtering (TUI)
- ✅ **Progress Indicators:** Spinners during async operations (TUI)
- ✅ **Export Formats:** JSON, HTML, JUnit XML for CI/CD (TUI)
- ✅ **Response Filtering:** Real-time filtering by status/method/keywords (TUI)
- ✅ **Request History:** Track and replay previous test runs (TUI)
- ✅ **Verbose Logging:** Full HTTP request/response details (TUI)
- ✅ **Configuration Persistence:** Auto-save settings to YAML (TUI)

## 🔧 Use Cases

### When to Use CLI

✅ **Choose CLI when you need:**
1. **CI/CD Pipelines** - Automated testing in GitHub Actions, GitLab CI
2. **Pre-commit Hooks** - Validate specs before committing
3. **API Monitoring** - Regular health checks via cron jobs
4. **Script Integration** - Easy to integrate into bash scripts
5. **Batch Testing** - Test multiple APIs sequentially
6. **Exit Codes** - Need 0/1 exit codes for automation
7. **Headless Environments** - No terminal interaction required

### When to Use TUI

✅ **Choose TUI when you need:**
1. **Interactive Development** - Real-time spec validation during development
2. **Manual Testing** - Quick endpoint testing with visual feedback
3. **API Exploration** - Discover and test API endpoints interactively
4. **Debugging** - Visual workflow for troubleshooting API issues
5. **Learning** - Understand API behavior through interactive testing
6. **Progress Indicators** - Want to see spinners and loading states
7. **Beautiful Output** - Prefer styled tables and colors

### Decision Matrix

```
┌─────────────────────┬─────────────┬─────────────┐
│ Scenario            │ Use CLI     │ Use TUI     │
├─────────────────────┼─────────────┼─────────────┤
│ GitHub Actions      │     ✓       │             │
│ Local Development   │             │     ✓       │
│ Bash Scripts        │     ✓       │             │
│ Manual Testing      │             │     ✓       │
│ Cron Jobs           │     ✓       │             │
│ API Exploration     │             │     ✓       │
│ Pre-commit Hook     │     ✓       │             │
│ Debugging Issues    │             │     ✓       │
│ Parallel Testing    │     ✓       │             │
│ Visual Feedback     │             │     ✓       │
└─────────────────────┴─────────────┴─────────────┘
```

## 📦 Installation

### Prerequisites

| Tool | Requirements |
|------|-------------|
| **CLI** | Node.js 16+, npm or yarn |
| **TUI** | Go 1.21+ |

### From Source

```bash
# Clone the repository
git clone https://github.com/yourusername/OpenAPI-Toolkit.git
cd OpenAPI-Toolkit

# Install CLI
cd openapi-cli
npm install && npm run build && npm link

# Build TUI
cd ../openapi-tui
go build

# Verify installations
openapi-test --version  # CLI
./openapi-tui           # TUI
```

## 🧪 Testing

### CLI Tests
```bash
cd openapi-cli
npm test

# Output:
# PASS  src/__tests__/commands/validate.test.ts
# Tests:       3 passed, 3 total
```

### TUI Tests
```bash
cd openapi-tui
go test ./...

# Output:
# PASS - 94 tests passing (170+ test runs including subtests)
# ok      github.com/.../internal/models  0.104s
# ok      github.com/.../internal/export  0.008s
# ok      github.com/.../internal/ui      0.003s
# 
# Test coverage by package:
# - history.go: 8 test functions (persistence, limits, replay)
# - export/junit.go: 20 test cases (XML structure, CI/CD format)
# - export/html.go: 25 test cases (professional reports)
# - ui/filter.go: 24 test cases (multi-field filtering)
# - ui/stats.go: 12 test cases (timing, pass rates)
```

## 🤝 Contributing

Contributions are welcome! Here are some areas for improvement:

### CLI Enhancements
- [ ] Response schema validation
- [ ] Authentication support (Bearer, API Key, Basic)
- [ ] Custom timeout configuration
- [ ] Export results to JSON/JUnit XML
- [ ] Parallel endpoint testing

### TUI Enhancements
- [x] Authentication support (Bearer, API Key, Basic) ✅
- [x] Enhanced error messages with suggestions ✅
- [x] Verbose logging mode (toggle with 'v') ✅
- [ ] HEAD and OPTIONS method support
- [ ] Summary statistics display
- [ ] Authentication UI for credential input
- [ ] Export test results (JSON/HTML/JUnit)
- [ ] Configuration file support
- [ ] Response body display in results

### Documentation
- [ ] Video tutorials
- [ ] More usage examples
- [ ] Architecture decision records

## 📝 Examples

### Example 1: Validate in CI/CD

```yaml
# .github/workflows/validate.yml
name: Validate OpenAPI Spec
on: [push, pull_request]
jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-node@v2
      - run: |
          cd openapi-cli
          npm install
          npm run build
          npm link
          openapi-test validate ../api/spec.yaml
```

### Example 2: Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

echo "Validating OpenAPI specification..."
openapi-test validate api/openapi.yaml

if [ $? -ne 0 ]; then
  echo "❌ OpenAPI spec validation failed. Commit blocked."
  exit 1
fi

echo "✅ OpenAPI spec is valid"
```

### Example 3: API Health Check Script

```bash
#!/bin/bash
# health-check.sh

SPEC="api/openapi.yaml"
API_URL="https://api.production.com"

echo "Running API health check..."
openapi-test test $SPEC $API_URL

if [ $? -eq 0 ]; then
  echo "✅ All endpoints healthy"
  exit 0
else
  echo "❌ Some endpoints failed"
  # Send alert to Slack/PagerDuty
  exit 1
fi
```

## 🐛 Troubleshooting

### Common Issues

**CLI: "command not found: openapi-test"**
```bash
# Solution: Run npm link in openapi-cli directory
cd openapi-cli && npm link
```

**TUI: "no such file or directory: openapi-tui"**
```bash
# Solution: Build the binary first
cd openapi-tui && go build
```

**Both: "File not found: spec.yaml"**
```bash
# Solution: Use absolute or correct relative path
openapi-test validate /full/path/to/spec.yaml
# OR
openapi-test validate ./relative/path/to/spec.yaml
```

**CLI: "Connection refused"**
- Ensure the API server is running
- Check the base URL is correct
- Verify no firewall blocking requests

**TUI: Terminal rendering issues**
- Ensure terminal supports UTF-8 and ANSI colors
- Resize terminal window (minimum 80x24)
- Use a modern terminal (iTerm2, Alacritty, Windows Terminal)

## 📄 License

MIT License - See [LICENSE](LICENSE) file for details

## 🙏 Acknowledgments

Built with:
- [Commander.js](https://github.com/tj/commander.js) - CLI framework
- [Axios](https://github.com/axios/axios) - HTTP client
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [kin-openapi](https://github.com/getkin/kin-openapi) - OpenAPI parser

## 📬 Contact & Support

- **Issues:** [GitHub Issues](https://github.com/yourusername/OpenAPI-Toolkit/issues)
- **Discussions:** [GitHub Discussions](https://github.com/yourusername/OpenAPI-Toolkit/discussions)
- **Documentation:** [Wiki](https://github.com/yourusername/OpenAPI-Toolkit/wiki)

---

## 🎓 Quick Reference

### CLI Commands Cheat Sheet

```bash
# Installation
npm install && npm run build && npm link

# Validate OpenAPI spec
openapi-test validate <spec-file>

# Test API endpoints
openapi-test test <spec-file> <base-url>

# Get help
openapi-test --help

# Version
openapi-test --version
```

### TUI Navigation Cheat Sheet

```
┌──────────────────────────────────────────────┐
│ Key           │ Action                       │
├───────────────┼──────────────────────────────┤
│ ↑ / k         │ Navigate up                  │
│ ↓ / j         │ Navigate down                │
│ Enter         │ Select option                │
│ Esc / q       │ Go back / Quit               │
│ h / ?         │ Show help                    │
│ Ctrl+C        │ Force quit                   │
└───────────────┴──────────────────────────────┘
```

### HTTP Methods Supported

| Method | CLI | TUI | Common Use |
|--------|-----|-----|------------|
| GET | ✅ | ✅ | Read resources |
| POST | ✅ | ✅ | Create resources |
| PUT | ✅ | ✅ | Update resources |
| PATCH | ✅ | ✅ | Partial updates |
| DELETE | ✅ | ✅ | Remove resources |
| HEAD | ✅ | ❌ | Headers only |
| OPTIONS | ✅ | ❌ | CORS preflight |

### Exit Codes (CLI)

| Code | Meaning | Use Case |
|------|---------|----------|
| 0 | Success | All tests passed / Valid spec |
| 1 | Failure | Tests failed / Invalid spec |

Perfect for CI/CD workflows!

---

## 🌟 Project Stats

```
📁 Total Lines of Code:  ~3,500+ (TUI: 1,400+ main, 1,150+ tests)
🧪 Test Coverage:        CLI: 100%, TUI: 36.9%
📝 Documentation Pages:  6
🔧 HTTP Methods:         7 (CLI), 5 (TUI)
⚡ Performance:          <2s for 15 endpoints
🎨 UI Components:        8 (TUI only)
🔐 Auth Methods:         3 (TUI: Bearer, API Key, Basic)
✨ Phase 2 Features:     2/15 complete (Enhanced errors, Verbose logging)
```

---

**Made with ❤️ for the OpenAPI community**

💡 **Pro Tip:** Use both! CLI in your pipeline, TUI for development.

Choose your tool: CLI for automation, TUI for interaction!

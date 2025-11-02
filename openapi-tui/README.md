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
- 🧪 **37.8% Test Coverage** — Comprehensive test suite with 70+ test cases

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

### Navigation
- **↑/↓ or j/k** — Navigate menu
- **Enter** — Select option
- **q/Esc** — Quit

### Typical Workflow
1. **Validate Spec** 📋
   - Enter path to your OpenAPI YAML/JSON file
   - Get instant validation results and errors

2. **Test API** 🧪
   - Provide spec file path
   - Enter base URL (e.g., `https://api.example.com`)
   - Watch automated endpoint testing (requests, validation, auth)

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
Run the test suite (includes auth/unit/integration tests):

```bash
go test ./... -cover
```

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
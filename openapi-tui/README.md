# 🚀 OpenAPI CLI TUI

A **fucking sick** terminal user interface for validating and testing APIs against OpenAPI specifications. Built with Go and the Charmbracelet ecosystem for maximum terminal awesomeness.

![Demo](https://img.shields.io/badge/Demo-Coming%20Soon-FF6B6B?style=for-the-badge&logo=terminal)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-v0.25.0-7D56F4?style=for-the-badge)

## ✨ Features

- 📋 **Validate OpenAPI Specs** - Ensure your API docs are valid and compliant
- 🧪 **Auto-Test Endpoints** - Automatically test all endpoints in your spec
- 🎨 **Beautiful TUI** - Gorgeous terminal interface with colors and borders
- ⚡ **Fast & Lightweight** - Built with Go for speed and efficiency
- 🔄 **Real-time Feedback** - See results instantly with styled output

## 🛠️ Installation

### Prerequisites
- Go 1.21 or later
- A terminal that supports colors (most modern terminals do)

### Install from source
```bash
git clone https://github.com/yourusername/openapi-cli-tui.git
cd openapi-cli-tui
go mod tidy
go install .
```

### Or install directly
```bash
go install github.com/yourusername/openapi-cli-tui@latest
```

## 🎮 Usage

```bash
openapi-cli-tui
```

### Navigation
- **↑/↓ or j/k** - Navigate menu
- **Enter** - Select option
- **q/Esc** - Quit

### Workflow
1. **Validate Spec** 📋
   - Enter path to your OpenAPI YAML/JSON file
   - Get instant validation results

2. **Test API** 🧪
   - Provide spec file path
   - Enter base URL (e.g., `https://api.example.com`)
   - Watch automated endpoint testing

## 🏗️ Architecture

Built with industry-standard libraries:

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** - Elm-inspired TUI framework
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** - CSS-like styling for terminals
- **[Kin OpenAPI](https://github.com/getkin/kin-openapi)** - Comprehensive OpenAPI toolkit

## 🚀 Development

```bash
# Clone and setup
git clone https://github.com/yourusername/openapi-cli-tui.git
cd openapi-cli-tui

# Install dependencies
go mod tidy

# Run locally
go run .

# Build for production
go build -o openapi-tui .
```

## 🤝 Contributing

PRs welcome! Make sure your code looks fucking sick too.

1. Fork it
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -am 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

MIT - See [LICENSE](LICENSE) file

---

**Made with ❤️ and lots of ☕ by the OpenAPI CLI TUI team**
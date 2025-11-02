# OpenAPI CLI Tester

A command-line tool for validating OpenAPI specifications and testing APIs against them.

## Installation

```bash
npm install -g .
```

Or run locally:

```bash
npm install
npm run build
npm link
```

## Usage

### Validate an OpenAPI spec

```bash
openapi-test validate path/to/spec.yaml
```

### Run tests against an API

```bash
openapi-test test path/to/spec.yaml http://api.example.com
```

## Development

```bash
npm run dev  # Run with ts-node
npm test     # Run tests
npm run lint # Lint code
```

## Features

- Validate OpenAPI 3.x specifications
- Test API endpoints automatically based on the spec
- Support for JSON and YAML formats
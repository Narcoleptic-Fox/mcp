# Contributing to MCP Go SDK

Thank you for your interest in contributing to the Model Context Protocol (MCP) Go SDK! This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to contact@narcolepticfox.com.

## Getting Started

### Prerequisites

- Go 1.16 or higher
- Git

### Development Setup

1. Fork the repository on GitHub
2. Clone your fork locally
   ```bash
   git clone https://github.com/YOUR-USERNAME/mcp.git
   cd mcp
   ```
3. Add the original repository as an upstream remote
   ```bash
   git remote add upstream https://github.com/narcolepticfox/mcp.git
   ```
4. Create a branch for your feature or bugfix
   ```bash
   git checkout -b feature-or-bugfix-name
   ```

## Development Workflow

### Running Tests

Run all tests with:

```bash
go test ./...
```

Or use the provided script:

```bash
./run_tests.sh
```

### Code Coverage

Generate code coverage reports with:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Code Style

- Follow standard Go coding conventions
- Use `gofmt` to format your code
- Run `golint` and `go vet` before submitting changes

## Pull Request Process

1. Update documentation as needed
2. Add tests for new functionality
3. Ensure all tests pass
4. Update the CHANGELOG.md with details of changes
5. Submit a pull request with a clear description of the changes

### Pull Request Guidelines

- Keep PRs focused on a single topic
- Write clear commit messages
- Reference issues and PRs in commit messages
- Update documentation as needed
- Rebase your branch before submitting a PR

## Releasing

Release process is managed by project maintainers. Version numbers follow [Semantic Versioning](https://semver.org/).

## Questions?

Feel free to open an issue with your question or contact the project maintainers directly.

Thank you for your contributions!

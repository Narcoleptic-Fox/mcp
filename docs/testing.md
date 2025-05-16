# MCP Testing Guide

This document outlines the testing strategy and implementation for the Model Context Protocol (MCP) Go SDK.

## Testing Philosophy

The MCP testing suite is built on the following principles:

1. **Test-driven development** - Write tests before or alongside code to ensure proper design and complete test coverage
2. **Comprehensive testing** - Cover every component with at least unit tests and integration tests where appropriate
3. **Performance awareness** - Use benchmarks to track performance and prevent regressions
4. **Maintainable tests** - Keep tests simple, descriptive, and focused on single responsibilities

## Test Coverage

The MCP SDK has complete test coverage across all major components:

| Component   | Unit Tests | Integration Tests | Benchmarks |
| ----------- | ---------- | ----------------- | ---------- |
| Core Models | ✅        | N/A                | ✅        |
| Status      | ✅        | N/A                | N/A        |
| Validation  | ✅        | N/A                | ✅        |
| Client      | ✅        | ✅                | ✅         |
| Server      | ✅        | ✅                | ✅         |
| Handlers    | ✅        | ✅                | ✅         |

## Types of Tests

### Unit Tests

Unit tests focus on testing individual components in isolation:

- Test a single function or method
- Mock external dependencies
- Fast execution
- Focused on specific behaviors

### Integration Tests

Integration tests verify interactions between components:

- Test multiple components working together
- Test real network communication
- Verify end-to-end workflows
- Test protocol compatibility

### Benchmark Tests

Benchmark tests measure performance characteristics:

- Response time
- Throughput
- Resource usage
- Scaling behavior

## Test Tags

Tests are organized using build tags to allow selective execution:

- `unit` - Unit tests
- `integration` - Integration tests
- `benchmark` - Performance benchmarks

## Test Structure

The test suite follows a consistent pattern across all packages:

### Unit Tests
Each package has its own `*_test.go` files that test individual components:

- `models_test.go` - Tests for core data models
- `status_test.go` - Tests for status tracking and status events
- `options_test.go` - Tests for configurable client and server options
- `handler_test.go` - Tests for request handlers
- `server_test.go` - Tests for server functionality
- `client_test.go` - Tests for client functionality

### Integration Tests
Integration tests are implemented in the same test files but are controlled via build tags.

### Benchmarks
Benchmark tests are collected in the root `benchmark_test.go` file and measure:
- Request/response performance
- Handler processing speed
- Client-server round trip time
- Concurrent request handling

## Test Utilities

The `testutil` package provides helpers for writing consistent and maintainable tests:

### Assertion Utilities

```go
// AssertEqual compares expected and actual values
testutil.AssertEqual(t, expected, actual, "Values should match")

// AssertNil checks for nil values
testutil.AssertNil(t, err, "Error should be nil")

// AssertNotNil checks for non-nil values
testutil.AssertNotNil(t, response, "Response should not be nil")

// AssertTrue checks for true conditions
testutil.AssertTrue(t, client.IsConnected(), "Client should be connected")

// AssertFalse checks for false conditions
testutil.AssertFalse(t, server.IsStopped(), "Server should not be stopped")
```

### Network Utilities

```go
// Get a free port for testing
port, err := testutil.GetFreePort()

// Create a test request
req := testutil.CreateTestModelRequest()

// Wait for a condition with timeout
success := testutil.WaitForCondition(5*time.Second, 100*time.Millisecond, func() bool {
    return client.Status() == core.StatusRunning
})
```

### Mock Server

The `testutil` package includes a mock server implementation for testing client code:

```go
// Create a mock server for testing
mockServer := testutil.NewMockServer()

// Configure mock responses
mockServer.SetResponseHandler(func(req *core.ModelRequest) *core.ModelResponse {
    resp := core.NewModelResponse(req)
    resp.Results["status"] = "processed"
    return resp
})

// Start the mock server
port, err := mockServer.Start()

// Create a client connecting to the mock server
client := client.New(client.WithServerPort(port))
```

## Running Tests

The MCP SDK includes a comprehensive test runner script that makes it easy to execute different types of tests:

```bash
# Run all tests and checks
./run_tests.sh --all

# Run only unit tests
./run_tests.sh --unit-only

# Run only integration tests
./run_tests.sh --integration-only

# Run benchmark tests
./run_tests.sh --benchmarks

# Generate test coverage report
./run_tests.sh --coverage

# Run code quality checks
./run_tests.sh --code-quality
```

### Running Tests Manually

You can also run tests directly with Go's test command:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./core/...

# Run tests with coverage report
go test -cover ./...

# Run benchmarks
go test -bench=. -benchmem ./...
```

## Continuous Integration

The test suite is designed to integrate with CI/CD pipelines. For CI environments:

1. Run `go mod download` to ensure dependencies are available
2. Run `go test ./...` to execute all tests
3. Run `go test -bench=. -benchmem ./...` to execute benchmarks
4. Use `go test -coverprofile=coverage.out ./...` to generate coverage data
5. Process coverage data with `go tool cover -html=coverage.out -o coverage.html`

### File Naming

- Unit tests: `package/*_test.go`
- Integration tests: `integration_test.go` or `*_integration_test.go`
- Benchmark tests: `benchmark_test.go` or `*_benchmark_test.go`

### Test Functions

- Unit tests: `TestXxx` (where Xxx describes what is being tested)
- Table-driven tests: Use `t.Run` with descriptive subtests
- Integration tests: `TestIntegrationXxx`
- Benchmark tests: `BenchmarkXxx`

## Test Utilities

### Helper Packages

- `testutil` - Common utilities for testing:
  - `mock_server.go` - Mock server implementation
  - `testutil.go` - Generic testing utilities

### Assertion Patterns

We use both the standard testing package and testify for assertions:

```go
// Standard library
if got != want {
    t.Errorf("Result = %v, want %v", got, want)
}

// Using testify
assert.Equal(t, want, got, "Result should match expected value")
require.NoError(t, err, "Operation should not return an error")
```

## Running Tests

### Command Line

```bash
# Run all tests
go test ./...

# Run only unit tests
go test -tags=unit ./...

# Run only integration tests
go test -tags=integration ./...

# Run benchmarks
go test -bench=. -benchmem ./...

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Using the Test Script

A convenience script is provided for running tests with various configurations:

```bash
# Run all tests
./run_tests.sh --all

# Run only unit tests
./run_tests.sh --unit-only

# Run with coverage
./run_tests.sh --coverage
```

## Continuous Integration

Tests are automatically run on CI for:

- Every push to main branch
- Every pull request
- Scheduled runs (daily)

The CI pipeline includes:

- Unit and integration tests across multiple Go versions
- Benchmark tests to detect performance regressions
- Code coverage reporting
- Static code analysis and linting

## Writing Good Tests

### Do's

- Write tests before or alongside implementation
- Keep tests simple and focused
- Use clear, descriptive names
- Test edge cases and error conditions
- Use table-driven tests for testing multiple similar cases
- Include comments explaining complex test setups

### Don'ts

- Don't test implementation details that may change
- Don't write brittle tests that break with minor changes
- Don't use time.Sleep() for synchronization (use channels or test utilities)
- Don't skip error checks
- Don't make tests dependent on each other

## Test Data

Test data should be:

- Representative of real-world usage
- Minimal (just enough to test the functionality)
- Clearly labeled with its purpose
- Generated programmatically when possible

## Mocking

Two approaches to mocking are used:

1. **Interface-based mocking** - Create mock implementations of interfaces
2. **Testify mocks** - Use the testify/mock package for more complex behaviors

Example:

```go
// Custom mock
type MockHandler struct {
    methods []string
}

func (m *MockHandler) Methods() []string {
    return m.methods
}

// Testify mock
type MockHandler struct {
    mock.Mock
}

func (m *MockHandler) Methods() []string {
    args := m.Called()
    return args.Get(0).([]string)
}
```

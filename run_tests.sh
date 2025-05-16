#!/bin/bash
# filepath: c:\Users\z004ecdm\source\repos\Personal\MCP SDKs\mcp\run_tests.sh

# MCP Test Suite Runner

set -e  # Exit on any error

# ANSI Color codes for formatting output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print header
echo -e "${BLUE}======================================${NC}"
echo -e "${BLUE}     MCP Test Suite Runner           ${NC}"
echo -e "${BLUE}======================================${NC}"

# Function to print section headers
print_header() {
    echo -e "\n${YELLOW}$1${NC}"
    echo -e "${YELLOW}$(printf '=%.0s' $(seq 1 ${#1}))${NC}\n"
}

# Parse command line arguments
run_unit=true
run_integration=true
run_benchmarks=false
run_code_quality=false
coverage=false

for arg in "$@"
do
    case $arg in
        --unit-only)
        run_integration=false
        run_benchmarks=false
        run_code_quality=false
        shift
        ;;
        --integration-only)
        run_unit=false
        run_benchmarks=false
        run_code_quality=false
        shift
        ;;
        --benchmarks)
        run_benchmarks=true
        shift
        ;;
        --code-quality)
        run_code_quality=true
        shift
        ;;
        --coverage)
        coverage=true
        shift
        ;;
        --all)
        run_unit=true
        run_integration=true
        run_benchmarks=true
        run_code_quality=true
        shift
        ;;
        --help)
        echo "Usage: run_tests.sh [options]"
        echo "Options:"
        echo "  --unit-only       Run only unit tests"
        echo "  --integration-only Run only integration tests"
        echo "  --benchmarks      Run benchmark tests"
        echo "  --code-quality    Run code quality checks"
        echo "  --coverage        Generate test coverage report"
        echo "  --all             Run all tests and checks"
        echo "  --help            Show this help message"
        exit 0
        ;;
    esac
done

# Ensure dependencies are installed
print_header "Checking dependencies"
go mod tidy
echo -e "${GREEN}Dependencies are up to date${NC}"

# Unit Tests
if $run_unit; then
    print_header "Running Unit Tests"
    
    if $coverage; then
        # Run tests with coverage
        go test -tags=unit ./... -cover -coverprofile=coverage.out
        go tool cover -html=coverage.out -o coverage.html
        echo -e "${GREEN}Coverage report generated: coverage.html${NC}"
    else
        # Run tests without coverage
        go test -tags=unit ./... -v
    fi
    
    echo -e "${GREEN}Unit tests completed${NC}"
fi

# Integration Tests
if $run_integration; then
    print_header "Running Integration Tests"
    go test -tags=integration ./... -v
    echo -e "${GREEN}Integration tests completed${NC}"
fi

# Benchmark Tests
if $run_benchmarks; then
    print_header "Running Benchmark Tests"
    go test -tags=benchmark -bench=. -benchmem ./...
    echo -e "${GREEN}Benchmark tests completed${NC}"
fi

# Code Quality Checks
if $run_code_quality; then
    print_header "Running Code Quality Checks"
    
    echo -e "${BLUE}Running go fmt...${NC}"
    go fmt ./...
    
    echo -e "${BLUE}Running go vet...${NC}"
    go vet ./...
    
    # Check if golangci-lint is installed
    if command -v golangci-lint &> /dev/null; then
        echo -e "${BLUE}Running golangci-lint...${NC}"
        golangci-lint run ./...
    else
        echo -e "${YELLOW}golangci-lint not found, skipping linting${NC}"
        echo -e "${YELLOW}To install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest${NC}"
    fi
    
    # Check if staticcheck is installed
    if command -v staticcheck &> /dev/null; then
        echo -e "${BLUE}Running staticcheck...${NC}"
        staticcheck ./...
    else
        echo -e "${YELLOW}staticcheck not found, skipping static analysis${NC}"
        echo -e "${YELLOW}To install: go install honnef.co/go/tools/cmd/staticcheck@latest${NC}"
    fi
    
    echo -e "${GREEN}Code quality checks completed${NC}"
fi

echo -e "\n${GREEN}All tests and checks completed successfully!${NC}"

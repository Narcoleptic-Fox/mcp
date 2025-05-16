// Package testutil provides utilities for testing MCP components.
package testutil

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/yourorg/mcp/core"
)

// AssertEqual compares expected and actual values and fails the test if they don't match.
func AssertEqual(t *testing.T, expected, actual interface{}, message string) {
	t.Helper()
	if expected != actual {
		t.Errorf("%s: expected %v but got %v", message, expected, actual)
	}
}

// AssertNil fails the test if the value is not nil.
func AssertNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value != nil {
		t.Errorf("%s: expected nil but got %v", message, value)
	}
}

// AssertNotNil fails the test if the value is nil.
func AssertNotNil(t *testing.T, value interface{}, message string) {
	t.Helper()
	if value == nil {
		t.Errorf("%s: unexpectedly got nil", message)
	}
}

// AssertTrue fails the test if the condition is not true.
func AssertTrue(t *testing.T, condition bool, message string) {
	t.Helper()
	if !condition {
		t.Errorf("%s: condition not true", message)
	}
}

// AssertFalse fails the test if the condition is not false.
func AssertFalse(t *testing.T, condition bool, message string) {
	t.Helper()
	if condition {
		t.Errorf("%s: condition not false", message)
	}
}

// GetFreePort returns an available TCP port.
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// CreateTestModelRequest creates a model request for testing.
func CreateTestModelRequest() *core.ModelRequest {
	req := core.NewModelRequest()
	req.ModelData["name"] = "Test Model"
	req.ModelData["value"] = 42
	req.Parameters = append(req.Parameters, core.Parameter{
		Name:  "param1",
		Value: "value1",
		Type:  "string",
	})
	return req
}

// WaitForCondition waits for a condition to be true or times out.
func WaitForCondition(timeout time.Duration, interval time.Duration, condition func() bool) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(interval)
	}
	return false
}

// ContextWithTimeout creates a context with timeout and handles defer cancel
func ContextWithTimeout(t *testing.T, duration time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	t.Cleanup(func() { cancel() })
	return ctx, cancel
}

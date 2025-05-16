// Package testutil provides utilities for testing MCP components.
package testutil

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/narcolepticfox/mcp/core"
)

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

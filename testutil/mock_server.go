// Package testutil provides utilities for testing MCP components.
package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/narcolepticfox/mcp/core"
	"github.com/sourcegraph/jsonrpc2"
)

// MockServer provides a test implementation of an MCP server.
type MockServer struct {
	t           *testing.T
	listener    net.Listener
	port        int
	conn        *jsonrpc2.Conn
	mutex       sync.Mutex
	handler     func(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error)
	shouldError bool
}

// NewMockServer creates a new mock server for testing.
func NewMockServer(t *testing.T) (*MockServer, error) {
	port, err := GetFreePort()
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return nil, err
	}

	mockServer := &MockServer{
		t:        t,
		listener: listener,
		port:     port,
		handler:  func(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error) { return nil, nil },
	}

	go mockServer.serve()

	return mockServer, nil
}

// serve handles connections to the mock server.
func (m *MockServer) serve() {
	for {
		conn, err := m.listener.Accept()
		if err != nil {
			// Expected when closing
			return
		}

		m.mutex.Lock()
		m.conn = jsonrpc2.NewConn(
			context.Background(),
			jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{}),
			jsonrpc2.HandlerWithError(m.handle),
		)
		m.mutex.Unlock()
	}
}

// handle processes JSON-RPC requests.
func (m *MockServer) handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	switch req.Method {
	case "mcp.processModel":
		if m.shouldError {
			return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeInvalidParams, Message: "Test error"}
		}

		var modelReq core.ModelRequest
		if err := json.Unmarshal(*req.Params, &modelReq); err != nil {
			return nil, fmt.Errorf("failed to unmarshal request params: %w", err)
		}

		return m.handler(ctx, &modelReq)
	default:
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: "Method not found"}
	}
}

// Check if this is correctly initialized
func (m *MockServer) Port() int {
	// Add nil check
	if m == nil {
		return 0 // or another appropriate default
	}
	return m.port
}

// SetShouldError configures the mock server to return errors.
func (m *MockServer) SetShouldError(shouldError bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.shouldError = shouldError
}

// SetupModelHandler configures a custom handler function for model processing requests.
func (m *MockServer) SetupModelHandler(handler func(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.handler = handler
}

// Close shuts down the mock server.
func (m *MockServer) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.conn != nil {
		m.conn.Close()
		m.conn = nil
	}

	return m.listener.Close()
}

// Start starts the mock server.
func (m *MockServer) Start() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	// Server is already started by NewMockServer, no need to do anything here
	return nil
}

// Stop stops the mock server.
func (m *MockServer) Stop() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.conn != nil {
		m.conn.Close()
		m.conn = nil
	}

	if m.listener != nil {
		return m.listener.Close()
	}

	return nil
}

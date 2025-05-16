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
	"github.com/narcolepticfox/mcp/server"
	"github.com/sourcegraph/jsonrpc2"
)

// MockServer provides a test implementation of an MCP server.
type MockServer struct {
	t           *testing.T
	listener    net.Listener
	port        int
	conn        *jsonrpc2.Conn
	mutex       sync.Mutex
	handler     server.ModelHandler
	shouldError bool
}

// NewMockServer creates a new mock server for testing.
func NewMockServer(t *testing.T) (*MockServer, error) {
	port, err := GetFreePort()
	if err != nil {
		return nil, err
	}

	listener, err := net.Listen("tcp", "localhost:"+string(port))
	if err != nil {
		return nil, err
	}

	mockServer := &MockServer{
		t:        t,
		listener: listener,
		port:     port,
		handler:  server.NewDefaultModelHandler(),
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

		return m.handler.ProcessModel(ctx, &modelReq)
	default:
		return nil, &jsonrpc2.Error{Code: jsonrpc2.CodeMethodNotFound, Message: "Method not found"}
	}
}

// Port returns the port number the mock server is listening on.
func (m *MockServer) Port() int {
	return m.port
}

// SetShouldError configures the mock server to return errors.
func (m *MockServer) SetShouldError(shouldError bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.shouldError = shouldError
}

// SetModelHandler configures the mock server to use a custom model handler.
func (m *MockServer) SetModelHandler(handler server.ModelHandler) {
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

package server

import (
	"context"
	"testing"
	"time"

	"github.com/narcolepticfox/mcp/client"
	"github.com/narcolepticfox/mcp/core"
	"github.com/narcolepticfox/mcp/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerLifecycle(t *testing.T) {
	// Get a free port for testing
	port, err := testutil.GetFreePort()
	require.NoError(t, err, "Failed to get free port")

	// Create a server
	srv := New(WithPort(port))

	// Verify initial state
	assert.Equal(t, core.StatusStopped, srv.Status(), "Server should start in stopped state")

	// Register for status change events
	var statusEvents []core.StatusChangeEvent
	srv.OnStatusChange(func(event core.StatusChangeEvent) {
		statusEvents = append(statusEvents, event)
	})

	// Start the server
	err = srv.Start()
	assert.NoError(t, err, "Start should succeed")
	defer srv.Stop() // Clean up after test

	// Server should be in running state
	assert.Equal(t, core.StatusRunning, srv.Status(), "Server should be in running state after start")

	// Stop the server
	err = srv.Stop()
	assert.NoError(t, err, "Stop should succeed")

	// Server should return to stopped state
	assert.Equal(t, core.StatusStopped, srv.Status(), "Server should return to stopped state after stop")

	// Check that at least two status events were recorded (idle->running, running->idle)
	assert.GreaterOrEqual(t, len(statusEvents), 2, "At least two status events should have been emitted")
}

func TestHandlerRegistration(t *testing.T) {
	// Create a server
	srv := New()

	// Create a custom handler
	handler := &MockModelHandler{
		methods: []string{"mcp.processModel"},
	}

	// Register the handler
	err := srv.RegisterHandler(handler)
	assert.NoError(t, err, "Handler registration should succeed")

	// Try to register the same method again
	duplicateHandler := &MockModelHandler{
		methods: []string{"mcp.processModel"},
	}
	err = srv.RegisterHandler(duplicateHandler)
	assert.Error(t, err, "Registering a duplicate method should fail")
}

func TestServerWithClient(t *testing.T) {
	// Get a free port for testing
	port, err := testutil.GetFreePort()
	require.NoError(t, err, "Failed to get free port")

	// Create a server
	srv := New(WithPort(port))

	// Register a handler
	handler := NewDefaultModelHandler()
	err = srv.RegisterHandler(handler)
	require.NoError(t, err, "Handler registration should succeed")

	// Start the server
	err = srv.Start()
	require.NoError(t, err, "Server should start successfully")
	defer srv.Stop()

	// Create a client that connects to our server
	c := client.New(
		client.WithServerPort(port),
		client.WithConnectionTimeout(2*time.Second),
	)

	// Start the client
	err = c.Start()
	require.NoError(t, err, "Client should connect to server")
	defer c.Stop()

	// Wait for the client to fully connect
	assert.True(t, testutil.WaitForCondition(2*time.Second, 100*time.Millisecond, func() bool {
		return c.Status() == core.StatusRunning
	}), "Client should enter running state")

	// Create a request
	req := testutil.CreateTestModelRequest()

	// Process a model request
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := c.ProcessModel(ctx, req)
	assert.NoError(t, err, "ProcessModel should not return an error")
	assert.NotNil(t, resp, "Response should not be nil")

	// Verify response
	assert.Equal(t, req.ID, resp.ID, "Response ID should match request ID")
	assert.True(t, resp.Success, "Response should indicate success")
	assert.Equal(t, "processed", resp.Results["status"], "Status should be set to 'processed'")
}

func TestServerRejectedRequest(t *testing.T) {
	// Get a free port for testing
	port, err := testutil.GetFreePort()
	require.NoError(t, err, "Failed to get free port")

	// Create a server
	srv := New(WithPort(port))

	// Register a custom handler that rejects requests
	handler := &MockModelHandler{
		methods: []string{"mcp.processModel"},
		processResponse: &core.ModelResponse{
			Success:      false,
			ErrorMessage: "rejected request",
			Results:      map[string]interface{}{},
		},
	}
	err = srv.RegisterHandler(handler)
	require.NoError(t, err, "Handler registration should succeed")

	// Start the server
	err = srv.Start()
	require.NoError(t, err, "Server should start successfully")
	defer srv.Stop()

	// Create a client that connects to our server
	c := client.New(
		client.WithServerPort(port),
		client.WithConnectionTimeout(2*time.Second),
	)

	// Start the client
	err = c.Start()
	require.NoError(t, err, "Client should connect to server")
	defer c.Stop()

	// Create a request
	req := testutil.CreateTestModelRequest()

	// Process a model request
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := c.ProcessModel(ctx, req)
	assert.NoError(t, err, "ProcessModel should not return an error even for rejected requests")
	assert.NotNil(t, resp, "Response should not be nil")

	// Verify the response reflects the rejection
	assert.False(t, resp.Success, "Response should indicate failure")
	assert.Equal(t, "rejected request", resp.ErrorMessage, "Error message should be set")
}

func TestServerRequestTimeout(t *testing.T) {
	// Get a free port for testing
	port, err := testutil.GetFreePort()
	require.NoError(t, err, "Failed to get free port")

	// Create a server with the specific port
	srv := New(WithPort(port))

	// Register a handler that sleeps for a period
	handler := &SlowModelHandler{delay: 500 * time.Millisecond}
	err = srv.RegisterHandler(handler)
	require.NoError(t, err, "Handler registration should succeed")

	// Start the server
	err = srv.Start()
	require.NoError(t, err, "Server should start successfully")
	defer srv.Stop()

	// Create a client that connects to our server with the correct port
	c := client.New(
		client.WithServerPort(port),
		client.WithConnectionTimeout(2*time.Second),
	)

	// Start the client
	err = c.Start()
	require.NoError(t, err, "Client should connect to server")
	defer c.Stop()

	// Create a request
	req := testutil.CreateTestModelRequest()

	// Process a model request with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	resp, err := c.ProcessModel(ctx, req)
	assert.Error(t, err, "ProcessModel should return an error when context times out")
	assert.Nil(t, resp, "Response should be nil when request times out")

	// Process a model request with sufficient timeout
	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel2()

	resp2, err2 := c.ProcessModel(ctx2, req)
	assert.NoError(t, err2, "ProcessModel should succeed with sufficient timeout")
	assert.NotNil(t, resp2, "Response should not be nil")
}

// SlowModelHandler implements a handler that sleeps before responding
type SlowModelHandler struct {
	delay time.Duration
}

func (h *SlowModelHandler) Methods() []string {
	return []string{"mcp.processModel"}
}

func (h *SlowModelHandler) ProcessModel(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(h.delay):
		resp := core.NewModelResponse(req)
		resp.Results["status"] = "processed after delay"
		return resp, nil
	}
}

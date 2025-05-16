package client

import (
	"context"
	"testing"
	"time"

	"github.com/narcolepticfox/mcp/core"
	"github.com/narcolepticfox/mcp/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientLifecycle(t *testing.T) {
	// Create a client
	client := New()

	// Verify initial state
	assert.Equal(t, core.StatusStopped, client.Status(), "Client should start in idle state")

	// Register for status change events
	var statusEvents []core.StatusChangeEvent
	client.OnStatusChange(func(event core.StatusChangeEvent) {
		statusEvents = append(statusEvents, event)
	})

	// Start the client (this should fail since there's no server running)
	assert.Error(t, client.Start(), "Start should fail when server is not available")

	// Client should be in error state after failed start
	assert.Equal(t, core.StatusFailed, client.Status(), "Client should be in failed state after failed start")

	// Check that at least one status event was recorded
	assert.NotEmpty(t, statusEvents, "At least one status event should have been emitted")
}

func TestClientWithMockServer(t *testing.T) {
	// Create a mock server
	mockServer, err := testutil.NewMockServer(t)

	// Create a client that connects to our mock server
	client := New(
		WithServerHost("localhost"),
		WithServerPort(mockServer.Port()),
		WithConnectionTimeout(2*time.Second),
		WithAutoReconnect(false),
	)

	// Set up the mock server to handle our test request
	testReq := testutil.CreateTestModelRequest()
	testResp := &core.ModelResponse{
		ID:           testReq.ID,
		Success:      true,
		ErrorMessage: "",
		Results: map[string]interface{}{
			"result": "test success",
		},
		Timestamp: time.Now(),
	}

	mockServer.SetupModelHandler(func(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error) {
		// Verify the request matches what we expect
		assert.Equal(t, testReq.ID, req.ID, "Request ID should match")
		return testResp, nil
	})

	// Start the client
	err = client.Start()
	assert.NoError(t, err, "Client should start successfully with mock server running")

	// Wait for the client to fully connect
	assert.True(t, testutil.WaitForCondition(2*time.Second, 100*time.Millisecond, func() bool {
		return client.Status() == core.StatusRunning
	}), "Client should enter running state")

	// Process a model request
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resp, err := client.ProcessModel(ctx, testReq)
	assert.NoError(t, err, "ProcessModel should not return an error")
	assert.NotNil(t, resp, "Response should not be nil")

	// Verify response matches what the mock server returned
	assert.Equal(t, testResp.ID, resp.ID, "Response ID should match")
	assert.Equal(t, testResp.Success, resp.Success, "Success flag should match")
	assert.Equal(t, testResp.Results["result"], resp.Results["result"], "Result value should match")

	err = client.Stop()
	assert.NoError(t, err, "Client should stop successfully")
}

func TestClientReconnect(t *testing.T) {
	// Only run this test if reconnect feature is implemented
	t.Skip("Reconnect test requires implementation of auto-reconnect feature")

	// Create a mock server
	mockServer, err := testutil.NewMockServer(t)
	require.NoError(t, err, "Failed to start mock server")

	// Create a client with auto-reconnect
	client := New(
		WithServerHost("localhost"),
		WithServerPort(mockServer.Port()),
		WithConnectionTimeout(2*time.Second),
		WithAutoReconnect(true),
		WithMaxReconnectAttempts(5),
		WithReconnectDelay(500*time.Millisecond),
	)

	// Start the client
	err = client.Start()
	assert.NoError(t, err, "Client should start successfully")
	err = client.Stop()
	assert.NoError(t, err, "Client should stop successfully")

	// Wait for the client to connect
	assert.True(t, testutil.WaitForCondition(2*time.Second, 100*time.Millisecond, func() bool {
		return client.Status() == core.StatusRunning
	}), "Client should enter running state")

	// Stop the server to simulate disconnection
	err = mockServer.Stop()
	require.NoError(t, err, "Failed to stop mock server")

	// Wait a bit for the client to detect disconnection
	time.Sleep(100 * time.Millisecond)

	// Restart the server
	err = mockServer.Start()
	require.NoError(t, err, "Failed to restart mock server")
	err = mockServer.Stop()
	require.NoError(t, err, "Failed to stop mock server")

	// Client should auto-reconnect
	assert.True(t, testutil.WaitForCondition(5*time.Second, 500*time.Millisecond, func() bool {
		return client.Status() == core.StatusRunning
	}), "Client should reconnect and return to running state")
}

func TestClientContextCancellation(t *testing.T) {
	// Create a mock server
	mockServer, err := testutil.NewMockServer(t)
	require.NoError(t, err, "Failed to create mock server")
	err = mockServer.Start()
	require.NoError(t, err, "Failed to start mock server")
	err = mockServer.Stop()
	require.NoError(t, err, "Failed to stop mock server")

	// Create a client
	client := New(
		WithServerHost("localhost"),
		WithServerPort(mockServer.Port()),
	)

	// Start the client
	err = client.Start()
	require.NoError(t, err, "Client should start successfully")
	err = client.Stop()
	require.NoError(t, err, "Client should stop successfully")

	// Configure mock server to delay response
	testReq := testutil.CreateTestModelRequest()
	mockServer.SetupModelHandler(func(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error) {
		// Check if context is canceled
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
			// This should not happen due to context cancellation
			return core.NewModelResponse(req), nil
		}
	})

	// Create a context with immediate cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Process model should fail due to canceled context
	resp, err := client.ProcessModel(ctx, testReq)
	assert.Error(t, err, "ProcessModel should return an error when context is canceled")
	assert.Nil(t, resp, "Response should be nil when context is canceled")
}

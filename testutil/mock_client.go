// Package testutil provides utilities for testing MCP components.
package testutil

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/narcolepticfox/mcp/core"
)

// MockClient provides a mock implementation of a client for testing.
type MockClient struct {
	status           core.Status
	statusMu         sync.RWMutex
	isConnected      bool
	callbacks        []func(core.StatusChangeEvent)
	processResponse  *core.ModelResponse
	processError     error
	startError       error
	stopError        error
	connectDelay     time.Duration
	processDelay     time.Duration
	requestsReceived []*core.ModelRequest
	mu               sync.Mutex
}

// NewMockClient creates a new mock client for testing.
func NewMockClient() *MockClient {
	return &MockClient{
		status:           core.StatusStopped,
		callbacks:        make([]func(core.StatusChangeEvent), 0),
		requestsReceived: make([]*core.ModelRequest, 0),
	}
}

// Start simulates starting the client.
func (c *MockClient) Start() error {
	c.statusMu.Lock()
	defer c.statusMu.Unlock()

	if c.status != core.StatusStopped {
		return errors.New("cannot start client in non-stopped state")
	}

	if c.startError != nil {
		c.status = core.StatusFailed
		return c.startError
	}

	// Simulate connection delay
	if c.connectDelay > 0 {
		time.Sleep(c.connectDelay)
	}

	oldStatus := c.status
	c.status = core.StatusRunning
	c.isConnected = true

	c.notifyStatusChange(oldStatus, c.status, nil)
	return nil
}

// Stop simulates stopping the client.
func (c *MockClient) Stop() error {
	c.statusMu.Lock()
	defer c.statusMu.Unlock()

	if c.status != core.StatusRunning {
		return errors.New("cannot stop client in non-running state")
	}

	if c.stopError != nil {
		return c.stopError
	}

	oldStatus := c.status
	c.status = core.StatusStopped
	c.isConnected = false

	c.notifyStatusChange(oldStatus, c.status, nil)
	return nil
}

// Status returns the current status of the mock client.
func (c *MockClient) Status() core.Status {
	c.statusMu.RLock()
	defer c.statusMu.RUnlock()
	return c.status
}

// IsConnected returns whether the mock client is connected.
func (c *MockClient) IsConnected() bool {
	c.statusMu.RLock()
	defer c.statusMu.RUnlock()
	return c.isConnected
}

// OnStatusChange registers a callback for status changes.
func (c *MockClient) OnStatusChange(callback func(core.StatusChangeEvent)) {
	c.statusMu.Lock()
	defer c.statusMu.Unlock()
	c.callbacks = append(c.callbacks, callback)
}

// ProcessModel simulates processing a model request.
func (c *MockClient) ProcessModel(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error) {
	c.mu.Lock()
	c.requestsReceived = append(c.requestsReceived, req)
	c.mu.Unlock()

	// Check if context is already canceled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// Continue processing
	}

	// Simulate processing delay
	if c.processDelay > 0 {
		select {
		case <-time.After(c.processDelay):
			// Continue after delay
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	// Return predefined response or error
	if c.processError != nil {
		return nil, c.processError
	}

	if c.processResponse != nil {
		return c.processResponse, nil
	}

	// Default response echoes back the request
	return &core.ModelResponse{
		ID:      req.ID,
		Success: true,
		Results: req.ModelData,
	}, nil
}

// SetStartError configures the mock client to return an error when Start is called.
func (c *MockClient) SetStartError(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.startError = err
}

// SetStopError configures the mock client to return an error when Stop is called.
func (c *MockClient) SetStopError(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stopError = err
}

// SetConnectDelay configures the mock client to delay during connect.
func (c *MockClient) SetConnectDelay(delay time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connectDelay = delay
}

// SetProcessDelay configures the mock client to delay during ProcessModel.
func (c *MockClient) SetProcessDelay(delay time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.processDelay = delay
}

// SetProcessResponse configures the mock client to return a specific response.
func (c *MockClient) SetProcessResponse(resp *core.ModelResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.processResponse = resp
}

// SetProcessError configures the mock client to return an error when ProcessModel is called.
func (c *MockClient) SetProcessError(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.processError = err
}

// GetRequestsReceived returns the requests that have been received by the mock client.
func (c *MockClient) GetRequestsReceived() []*core.ModelRequest {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.requestsReceived
}

// Reset resets the mock client to its initial state.
func (c *MockClient) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.statusMu.Lock()
	defer c.statusMu.Unlock()

	c.status = core.StatusStopped
	c.isConnected = false
	c.processResponse = nil
	c.processError = nil
	c.startError = nil
	c.stopError = nil
	c.connectDelay = 0
	c.processDelay = 0
	c.requestsReceived = make([]*core.ModelRequest, 0)
}

// notifyStatusChange notifies all registered callbacks about a status change.
func (c *MockClient) notifyStatusChange(oldStatus, newStatus core.Status, err error) {
	event := core.StatusChangeEvent{
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Timestamp: time.Now(),
		Error:     err,
	}

	// Copy callbacks to avoid holding the lock during callback execution
	var callbacksCopy []func(core.StatusChangeEvent)
	for _, callback := range c.callbacks {
		callbacksCopy = append(callbacksCopy, callback)
	}

	// Execute callbacks without holding the lock
	for _, callback := range callbacksCopy {
		go callback(event)
	}
}

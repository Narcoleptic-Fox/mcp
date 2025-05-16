package server

import (
	"context"
	"testing"

	"github.com/narcolepticfox/mcp/core"
	"github.com/stretchr/testify/assert"
)

func TestDefaultModelHandler_Methods(t *testing.T) {
	handler := NewDefaultModelHandler()
	methods := handler.Methods()

	assert.NotEmpty(t, methods, "Default handler methods should not be empty")
	assert.Contains(t, methods, "mcp.processModel", "Default handler should implement mcp.processModel method")
}

func TestDefaultModelHandler_ProcessModel(t *testing.T) {
	handler := NewDefaultModelHandler()

	// Create a test request
	req := core.NewModelRequest()
	req.ID = "test-request"
	req.ModelData["name"] = "Test Model"
	req.ModelData["value"] = 42

	// Call the handler
	ctx := context.Background()
	resp, err := handler.ProcessModel(ctx, req)

	// Check results
	assert.NoError(t, err, "ProcessModel should not return an error")
	assert.NotNil(t, resp, "Response should not be nil")

	// Verify basic properties
	assert.Equal(t, req.ID, resp.ID, "Response ID should match request ID")
	assert.True(t, resp.Success, "Response should indicate success")

	// Verify default handler sets expected results
	assert.Equal(t, "processed", resp.Results["status"], "Status should be set to 'processed'")
	assert.Equal(t, "Model processed successfully", resp.Results["message"], "Message should be set correctly")
}

// MockHandler implements the Handler interface for testing
type MockHandler struct {
	methods []string
}

func (m *MockHandler) Methods() []string {
	return m.methods
}

func TestCustomHandler(t *testing.T) {
	// Create a custom handler
	mockHandler := &MockHandler{
		methods: []string{"custom.method1", "custom.method2"},
	}

	// Verify methods
	methods := mockHandler.Methods()
	assert.Len(t, methods, 2, "Should have 2 methods")
	assert.Equal(t, "custom.method1", methods[0], "First method should match")
	assert.Equal(t, "custom.method2", methods[1], "Second method should match")
}

// MockModelHandler implements the ModelHandler interface for testing
type MockModelHandler struct {
	methods         []string
	processResponse *core.ModelResponse
	processError    error
}

func (m *MockModelHandler) Methods() []string {
	return m.methods
}

func (m *MockModelHandler) ProcessModel(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error) {
	if m.processResponse == nil && m.processError == nil {
		// Default behavior
		resp := core.NewModelResponse(req)
		resp.Results["handler"] = "mock"
		return resp, nil
	}
	return m.processResponse, m.processError
}

func TestCustomModelHandler(t *testing.T) {
	// Create a custom model handler
	mockModelHandler := &MockModelHandler{
		methods: []string{"mcp.processModel"},
	}

	// Create a test request
	req := core.NewModelRequest()
	req.ID = "test-custom-handler"

	// Call the handler
	ctx := context.Background()
	resp, err := mockModelHandler.ProcessModel(ctx, req)

	// Check results
	assert.NoError(t, err, "ProcessModel should not return an error")
	assert.NotNil(t, resp, "Response should not be nil")
	assert.Equal(t, req.ID, resp.ID, "Response ID should match request ID")
	assert.Equal(t, "mock", resp.Results["handler"], "Handler should set expected result")
}

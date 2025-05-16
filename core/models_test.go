package core

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewModelRequest(t *testing.T) {
	req := NewModelRequest()

	// Check that the request ID is not empty
	assert.NotEmpty(t, req.ID, "Request ID should not be empty")

	// Verify the ModelData map is initialized
	assert.NotNil(t, req.ModelData, "ModelData should be initialized")

	// Verify the Parameters slice is initialized
	assert.NotNil(t, req.Parameters, "Parameters should be initialized")
	assert.Len(t, req.Parameters, 0, "Parameters should be initialized as empty")
}

func TestNewModelResponse(t *testing.T) {
	// Create a request to link to the response
	req := NewModelRequest()
	req.ID = "test-request-123"

	// Create a response from the request
	resp := NewModelResponse(req)

	// Verify the response attributes
	assert.Equal(t, req.ID, resp.ID, "Response ID should match Request ID")
	assert.True(t, resp.Success, "Response should be marked as successful by default")
	assert.Empty(t, resp.ErrorMessage, "Error message should be empty for successful response")
	assert.NotNil(t, resp.Results, "Results map should be initialized")
	assert.WithinDuration(t, time.Now(), resp.Timestamp, 2*time.Second, "Timestamp should be current time")
}

func TestErrorResponse(t *testing.T) {
	// Create a request to link to the response
	req := NewModelRequest()
	req.ID = "test-request-456"

	err := fmt.Errorf("test error message")

	// Create an error response
	resp := ErrorResponse(req, err)

	// Verify the error response attributes
	assert.Equal(t, req.ID, resp.ID, "Response ID should match Request ID")
	assert.False(t, resp.Success, "Response should be marked as unsuccessful")
	assert.Equal(t, err.Error(), resp.ErrorMessage, "Error message should be set correctly")
	assert.NotNil(t, resp.Results, "Results map should be initialized")
}

func TestParameter(t *testing.T) {
	// Test parameter creation and value handling
	testCases := []struct {
		name  string
		param Parameter
		value interface{}
		typ   string
	}{
		{
			name: "string parameter",
			param: Parameter{
				Name:  "stringParam",
				Value: "test string",
				Type:  "string",
			},
			value: "test string",
			typ:   "string",
		},
		{
			name: "integer parameter",
			param: Parameter{
				Name:  "intParam",
				Value: 42,
				Type:  "int",
			},
			value: 42,
			typ:   "int",
		},
		{
			name: "boolean parameter",
			param: Parameter{
				Name:  "boolParam",
				Value: true,
				Type:  "boolean",
			},
			value: true,
			typ:   "boolean",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.value, tc.param.Value, "Parameter value should match expected value")
			assert.Equal(t, tc.typ, tc.param.Type, "Parameter type should match expected type")
		})
	}
}

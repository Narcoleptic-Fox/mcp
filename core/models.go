// Package core provides the fundamental models and interfaces for the Model Context Protocol (MCP).
package core

import (
	"time"
)

// ModelRequest represents a request to process a model.
// It contains the request identifier, model data, and processing parameters.
type ModelRequest struct {
	ID         string                 `json:"id"`
	ModelData  map[string]interface{} `json:"modelData"`
	Parameters []Parameter            `json:"parameters"`
}

// ModelResponse represents the response from processing a model.
// It includes the request identifier, success status, any error message,
// processing results, and a timestamp.
type ModelResponse struct {
	ID           string                 `json:"id"`
	Success      bool                   `json:"success"`
	ErrorMessage string                 `json:"errorMessage,omitempty"`
	Results      map[string]interface{} `json:"results"`
	Timestamp    time.Time              `json:"timestamp"`
}

// Parameter represents a named parameter with type information for model processing.
type Parameter struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

// NewModelRequest creates a new request with a generated ID.
// The returned request has initialized maps and slices ready to use.
func NewModelRequest() *ModelRequest {
	return &ModelRequest{
		ID:         generateID(),
		ModelData:  make(map[string]interface{}),
		Parameters: make([]Parameter, 0),
	}
}

// NewModelResponse creates a response for a given request.
// The response contains the same ID as the request and is initialized
// with a success status, empty results map, and current timestamp.
func NewModelResponse(req *ModelRequest) *ModelResponse {
	return &ModelResponse{
		ID:        req.ID,
		Success:   true,
		Results:   make(map[string]interface{}),
		Timestamp: time.Now(),
	}
}

// ErrorResponse creates an error response for a given request with the provided error.
// The response is marked as unsuccessful and includes the error message.
func ErrorResponse(req *ModelRequest, err error) *ModelResponse {
	return &ModelResponse{
		ID:           req.ID,
		Success:      false,
		ErrorMessage: err.Error(),
		Results:      make(map[string]interface{}),
		Timestamp:    time.Now(),
	}
}

// generateID creates a new unique ID using a timestamp-based approach.
// Note: In a production environment, consider using UUID or another more robust
// identifier generation method.
func generateID() string {
	return "mcp-" + time.Now().Format("20060102150405")
}

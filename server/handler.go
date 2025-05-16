// Package server provides a server implementation for the Model Context Protocol (MCP).
// It allows applications to receive and process model requests using registered handlers.
package server

import (
	"context"

	"github.com/yourorg/mcp/core"
)

// Handler defines the interface for MCP request handlers.
// Handlers implement specific RPC methods that can be invoked by clients.
type Handler interface {
	// Methods returns the list of method names that this handler implements.
	// These method names are used to register the handler with the server.
	Methods() []string
}

// ModelHandler handles model processing requests.
// It extends the base Handler interface with a method for processing model requests.
type ModelHandler interface {
	Handler
	// ProcessModel processes a model request and returns a response or an error.
	// The context can be used for cancellation and timeout control.
	ProcessModel(context.Context, *core.ModelRequest) (*core.ModelResponse, error)
}

// DefaultModelHandler provides a default implementation of the ModelHandler interface.
// It can be used as a starting point for custom model handlers or for testing.
type DefaultModelHandler struct{}

// NewDefaultModelHandler creates a new instance of DefaultModelHandler.
// This provides a simple handler that can be registered with an MCP server.
func NewDefaultModelHandler() *DefaultModelHandler {
	return &DefaultModelHandler{}
}

// Methods returns the list of method names that this handler implements.
// For DefaultModelHandler, this includes only the model processing method.
func (h *DefaultModelHandler) Methods() []string {
	return []string{"mcp.processModel"}
}

// ProcessModel processes a model request and returns a successful response.
// This default implementation simply acknowledges the request without performing
// any actual model processing. It should be overridden in production handlers.
func (h *DefaultModelHandler) ProcessModel(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error) {
	resp := core.NewModelResponse(req)

	// In a real implementation, this would process the model
	resp.Results["status"] = "processed"
	resp.Results["message"] = "Model processed successfully"

	return resp, nil
}

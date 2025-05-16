# Creating Custom Handlers

This guide explains how to create custom handlers for the Model Context Protocol (MCP) server.

## Overview

Handlers in MCP are responsible for processing specific types of requests. The server routes incoming requests to the appropriate handler based on the method name. By implementing custom handlers, you can add support for new functionality or override the default behavior.

## Types of Handlers

The MCP server supports different types of handlers:

- **Basic Handler**: Implements the `Handler` interface
- **Model Handler**: Implements the `ModelHandler` interface, specifically for model processing

## Implementing a Basic Handler

To implement a basic handler:

1. Create a struct that implements the `Handler` interface:

```go
import "github.com/narcolepticfox/mcp/server"

// MyCustomHandler is a basic handler
type MyCustomHandler struct {
    // Add any required fields here
}

// Methods returns the methods this handler implements
func (h *MyCustomHandler) Methods() []string {
    return []string{"mcp.customMethod"}
}
```

2. Implement your custom method handler:

```go
// HandleCustomMethod handles custom method calls
func (h *MyCustomHandler) HandleCustomMethod(ctx context.Context, req interface{}) (interface{}, error) {
    // Your custom logic goes here
    return map[string]string{"result": "success"}, nil
}
```

## Implementing a Model Handler

To implement a model handler:

1. Create a struct that implements the `ModelHandler` interface:

```go
import (
    "context"
    
    "github.com/narcolepticfox/mcp/core"
    "github.com/narcolepticfox/mcp/server"
)

// MyModelHandler is a custom model handler
type MyModelHandler struct {
    // Add any required fields here
}

// Methods returns the methods this handler implements
func (h *MyModelHandler) Methods() []string {
    return []string{"mcp.processModel"}
}
```

2. Implement the `ProcessModel` method:

```go
// ProcessModel processes a model request
func (h *MyModelHandler) ProcessModel(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error) {
    // Create a response linked to the request
    resp := core.NewModelResponse(req)
    
    // Your model processing logic goes here
    // For example:
    modelName, ok := req.ModelData["name"].(string)
    if !ok {
        resp.Success = false
        resp.ErrorMessage = "model name not provided or not a string"
        return resp, nil
    }
    
    // Process the model...
    resp.Results["processedModel"] = modelName
    resp.Results["status"] = "processed successfully"
    
    return resp, nil
}
```

## Registering a Handler

Once you have implemented your handler, register it with the server:

```go
// Create a new server
srv := server.New()

// Create your custom handlers
modelHandler := &MyModelHandler{}

// Register the handlers
if err := srv.RegisterHandler(modelHandler); err != nil {
    log.Fatalf("Failed to register model handler: %v", err)
}

// Start the server
if err := srv.Start(); err != nil {
    log.Fatalf("Failed to start server: %v", err)
}
```

## Error Handling

When implementing handlers, consider these error handling practices:

1. Use the response's error fields for application-level errors:

```go
func (h *MyModelHandler) ProcessModel(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error) {
    resp := core.NewModelResponse(req)
    
    // Application-level error
    if !isValidRequest(req) {
        resp.Success = false
        resp.ErrorMessage = "invalid request format"
        return resp, nil
    }
    
    // Process the request...
    return resp, nil
}
```

2. Return Go errors for system-level errors:

```go
func (h *MyModelHandler) ProcessModel(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error) {
    // System-level error
    if dbConn == nil {
        return nil, errors.New("database connection unavailable")
    }
    
    resp := core.NewModelResponse(req)
    // Process the request...
    return resp, nil
}
```

## Context Usage

The context passed to handler methods can be used for:

- Request cancellation detection
- Timeouts
- Passing request-scoped values

Example:

```go
func (h *MyModelHandler) ProcessModel(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error) {
    resp := core.NewModelResponse(req)
    
    // Use context for cancellation
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
        // Continue processing...
    }
    
    // Process with timeout
    processingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    result, err := processModelWithContext(processingCtx, req.ModelData)
    if err != nil {
        resp.Success = false
        resp.ErrorMessage = err.Error()
        return resp, nil
    }
    
    resp.Results["result"] = result
    return resp, nil
}
```

## Advanced Handlers

For more advanced scenarios:

### Stateful Handlers

```go
type StatefulModelHandler struct {
    mu       sync.Mutex
    counters map[string]int
}

func NewStatefulModelHandler() *StatefulModelHandler {
    return &StatefulModelHandler{
        counters: make(map[string]int),
    }
}

func (h *StatefulModelHandler) Methods() []string {
    return []string{"mcp.processModel"}
}

func (h *StatefulModelHandler) ProcessModel(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error) {
    resp := core.NewModelResponse(req)
    
    modelID, _ := req.ModelData["id"].(string)
    if modelID != "" {
        h.mu.Lock()
        h.counters[modelID]++
        count := h.counters[modelID]
        h.mu.Unlock()
        
        resp.Results["processCount"] = count
    }
    
    // Process the model...
    
    return resp, nil
}
```

### Chain of Responsibility

```go
type ModelHandlerMiddleware func(core.ModelHandler) core.ModelHandler

func LoggingMiddleware(next core.ModelHandler) core.ModelHandler {
    return &loggingModelHandler{next: next}
}

type loggingModelHandler struct {
    next core.ModelHandler
}

func (h *loggingModelHandler) Methods() []string {
    return h.next.Methods()
}

func (h *loggingModelHandler) ProcessModel(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error) {
    log.Printf("Processing model request: %s", req.ID)
    start := time.Now()
    
    resp, err := h.next.ProcessModel(ctx, req)
    
    log.Printf("Completed model request: %s, took: %v", req.ID, time.Since(start))
    return resp, err
}
```

Usage:

```go
baseHandler := &MyModelHandler{}
handlerWithLogging := LoggingMiddleware(baseHandler)

srv.RegisterHandler(handlerWithLogging)
```

## Best Practices

1. Keep handlers focused on a single responsibility
2. Use appropriate error handling (application vs. system errors)
3. Implement proper validation for incoming requests
4. Handle context cancellation
5. Use thread-safe techniques for stateful handlers
6. Add logging for debugging purposes
7. Consider middleware for cross-cutting concerns

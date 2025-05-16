# Model Context Protocol (MCP) for Go

A Go implementation of the Model Context Protocol, enabling applications to communicate with AI models through a standardized interface.

## Overview

The Model Context Protocol (MCP) provides a standardized way for applications to interact with AI models, regardless of their implementation details. This Go SDK includes both client and server components, allowing applications to:

- Make requests to MCP-compatible services
- Implement MCP-compatible services
- Process model data with customizable handlers
- Handle connection management and error recovery

## Installation

To install the MCP Go SDK, use the following command:

```bash
go get github.com/narcolepticfox/mcp
```

## Requirements

- Go 1.16 or higher

## Quick Start

### Server Example

Create an MCP server that processes model requests:

```go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/narcolepticfox/mcp/core"
	"github.com/narcolepticfox/mcp/server"
)

func main() {
	// Create a server with custom options
	srv := server.New(
		server.WithHost("0.0.0.0"),
		server.WithPort(5000),
	)

	// Register the default model handler
	handler := server.NewDefaultModelHandler()
	if err := srv.RegisterHandler(handler); err != nil {
		log.Fatalf("Failed to register handler: %v", err)
	}

	// Start the server
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Stop the server
	if err := srv.Stop(); err != nil {
		log.Fatalf("Failed to stop server: %v", err)
	}
}
```

### Client Example

Create an MCP client that sends model requests:

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/narcolepticfox/mcp/client"
	"github.com/narcolepticfox/mcp/core"
)

func main() {
	// Create a client with custom options
	c := client.New(
		client.WithServerHost("localhost"),
		client.WithServerPort(5000),
	)

	// Start the client
	if err := c.Start(); err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}
	defer c.Stop()

	// Create a model request
	req := core.NewModelRequest()
	req.ModelData["name"] = "Test Model"
	req.ModelData["value"] = 42
	req.Parameters = append(req.Parameters, core.Parameter{
		Name:  "param1",
		Value: "value1",
		Type:  "string",
	})

	// Send the request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.ProcessModel(ctx, req)
	if err != nil {
		log.Fatalf("Failed to process model: %v", err)
	}

	log.Printf("Response: %+v", resp)
}
```

## Architecture

The MCP Go SDK is organized into three main packages:

### Core Package

The core package contains fundamental types and interfaces used by both clients and servers:

- `ModelRequest` and `ModelResponse` - Data structures for model processing
- `Parameter` - Represents parameters for model requests
- `Status` - Represents component status and lifecycle
- `StatusChangeEvent` - Event emitted when component status changes

### Client Package

The client package provides functionality for applications that need to make requests to MCP services:

- `Client` - Main client interface for connecting to MCP servers
- `Options` - Configuration options for the client
- Functional options for customizing client behavior

### Server Package

The server package enables applications to implement MCP services:

- `Server` - Main server interface for receiving and handling requests
- `Handler` - Interface for implementing custom request handlers
- `ModelHandler` - Specialized handler for model processing requests
- `Options` - Configuration options for the server
- Functional options for customizing server behavior

## Configuration Options

### Server Options

- `WithHost(string)` - Set the host address to bind to
- `WithPort(int)` - Set the port number to listen on
- `WithMaxConcurrentClients(int)` - Set maximum concurrent client connections
- `WithConnectionTimeout(time.Duration)` - Set connection timeout
- `WithTLS(bool)` - Enable/disable TLS
- `WithCertificatePath(string)` - Set path to TLS certificate
- `WithCertificateKeyPath(string)` - Set path to TLS certificate key

### Client Options

- `WithServerHost(string)` - Set the server host to connect to
- `WithServerPort(int)` - Set the server port to connect to
- `WithConnectionTimeout(time.Duration)` - Set connection timeout
- `WithAutoReconnect(bool)` - Enable/disable automatic reconnection
- `WithMaxReconnectAttempts(int)` - Set maximum reconnect attempts
- `WithReconnectDelay(time.Duration)` - Set delay between reconnect attempts
- `WithTLS(bool)` - Enable/disable TLS

## Custom Handlers

You can implement custom handlers to process specific types of requests:

```go
type CustomModelHandler struct{}

func (h *CustomModelHandler) Methods() []string {
	return []string{"mcp.processModel"}
}

func (h *CustomModelHandler) ProcessModel(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error) {
	// Custom processing logic here
	resp := core.NewModelResponse(req)
	// ...
	return resp, nil
}
```

## Error Handling

The MCP SDK includes comprehensive error handling:

- Connection errors
- Timeout errors
- Validation errors
- Protocol errors

All errors are properly typed and include contextual information to aid debugging.

## Status Management

Both client and server components implement the `Component` interface, which provides:

- Status reporting
- Status change notifications
- Lifecycle management (Start/Stop)

## Contributing

Contributions to the MCP Go SDK are welcome! Please feel free to submit pull requests or open issues on the project repository.

## License

[Specify your license here]

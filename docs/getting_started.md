# Getting Started with MCP

This guide will help you get started with the Model Context Protocol (MCP) Go SDK.

## Installation

To install the MCP Go SDK, use the following command:

```bash
go get github.com/narcolepticfox/mcp
```

## Setting Up a Simple Server

Here's how to create a basic MCP server:

1. Create a new Go project:

```bash
mkdir mcp-server-example
cd mcp-server-example
go mod init example.com/mcp-server
```

2. Create a `main.go` file:

```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/narcolepticfox/mcp/core"
	"github.com/narcolepticfox/mcp/server"
)

// CustomModelHandler demonstrates a simple model handler implementation
type CustomModelHandler struct{}

func (h *CustomModelHandler) Methods() []string {
	return []string{"mcp.processModel"}
}

func (h *CustomModelHandler) ProcessModel(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error) {
	log.Printf("Received request: %s", req.ID)
	
	resp := core.NewModelResponse(req)
	resp.Results["message"] = "Hello from the MCP server!"
	
	// Extract model name if provided
	if name, ok := req.ModelData["name"].(string); ok {
		resp.Results["greeting"] = "Hello, " + name + "!"
	}
	
	return resp, nil
}

func main() {
	// Create a server with custom options
	srv := server.New(
		server.WithHost("0.0.0.0"),  // Listen on all interfaces
		server.WithPort(5000),        // Port 5000
	)

	// Register our custom handler
	handler := &CustomModelHandler{}
	if err := srv.RegisterHandler(handler); err != nil {
		log.Fatalf("Failed to register handler: %v", err)
	}

	// Register a status change callback
	srv.OnStatusChange(func(event core.StatusChangeEvent) {
		log.Printf("Server status changed: %s -> %s",
			event.OldStatus, event.NewStatus)
		if event.Error != nil {
			log.Printf("Error: %v", event.Error)
		}
	})

	// Start the server
	log.Println("Starting MCP server on port 5000...")
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Stop the server gracefully
	log.Println("Shutting down MCP server...")
	if err := srv.Stop(); err != nil {
		log.Fatalf("Failed to stop server: %v", err)
	}
	log.Println("Server stopped")
}
```

3. Install dependencies and run the server:

```bash
go mod tidy
go run main.go
```

## Setting Up a Simple Client

Here's how to create a basic MCP client:

1. Create a new Go project:

```bash
mkdir mcp-client-example
cd mcp-client-example
go mod init example.com/mcp-client
```

2. Create a `main.go` file:

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
		client.WithAutoReconnect(true),
	)

	// Register a status change callback
	c.OnStatusChange(func(event core.StatusChangeEvent) {
		log.Printf("Client status changed: %s -> %s",
			event.OldStatus, event.NewStatus)
		if event.Error != nil {
			log.Printf("Error: %v", event.Error)
		}
	})

	// Start the client
	log.Println("Starting MCP client...")
	if err := c.Start(); err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}
	defer func() {
		log.Println("Stopping MCP client...")
		if err := c.Stop(); err != nil {
			log.Printf("Failed to stop client: %v", err)
		}
		log.Println("Client stopped")
	}()

	// Create a model request
	req := core.NewModelRequest()
	req.ModelData["name"] = "MCP User"
	req.ModelData["value"] = 42
	req.Parameters = append(req.Parameters, core.Parameter{
		Name:  "param1",
		Value: "value1",
		Type:  "string",
	})

	// Send the request
	log.Println("Sending model processing request...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := c.ProcessModel(ctx, req)
	if err != nil {
		log.Fatalf("Failed to process model: %v", err)
	}

	log.Println("Request successful!")
	log.Printf("Response ID: %s", resp.ID)
	log.Printf("Success: %t", resp.Success)
	log.Printf("Results: %+v", resp.Results)
	
	if greeting, ok := resp.Results["greeting"].(string); ok {
		log.Printf("Server greeting: %s", greeting)
	}
}
```

3. Install dependencies and run the client:

```bash
go mod tidy
go run main.go
```

## Testing the Connection

1. First, start the server in one terminal window.
2. Then, run the client in another terminal window.
3. You should see:
   - The server logs the incoming request
   - The client displays the response from the server
   - The connection is established successfully

## Next Steps

Now that you have a basic MCP client and server running, you can:

1. Implement more sophisticated handlers for your specific use cases
2. Add validation to your requests and responses
3. Implement error handling and recovery strategies
4. Scale your server to handle more concurrent connections
5. Add security features such as TLS and authentication

For more advanced usage, check out:
- [API Reference](api_reference.md)
- [Creating Custom Handlers](custom_handlers.md)

// Example client application for the Model Context Protocol (MCP).
// This demonstrates how to create a client, connect to a server,
// and send model processing requests.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	if err := c.Start(); err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}

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

	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Stop the client
	if err := c.Stop(); err != nil {
		log.Fatalf("Failed to stop client: %v", err)
	}
}

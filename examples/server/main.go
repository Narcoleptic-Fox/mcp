// Example server application for the Model Context Protocol (MCP).
// This demonstrates how to create a server, register handlers,
// and process incoming model requests.
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

	// Register a status change callback
	srv.OnStatusChange(func(event core.StatusChangeEvent) {
		log.Printf("Server status changed: %s -> %s",
			event.OldStatus, event.NewStatus)
		if event.Error != nil {
			log.Printf("Error: %v", event.Error)
		}
	})

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

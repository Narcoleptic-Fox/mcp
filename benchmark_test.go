package mcp

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/narcolepticfox/mcp/client"
	"github.com/narcolepticfox/mcp/core"
	"github.com/narcolepticfox/mcp/server"
	"github.com/narcolepticfox/mcp/testutil"
)

// BenchmarkLocalRequestResponse measures the round-trip time for local requests.
func BenchmarkLocalRequestResponse(b *testing.B) {
	// Get a free port for testing
	port, err := testutil.GetFreePort()
	if err != nil {
		b.Fatalf("Failed to get free port: %v", err)
	}

	// Create and start server
	srv := server.New(
		server.WithPort(port),
		server.WithMaxConcurrentClients(100),
	)

	// Register default handler
	handler := server.NewDefaultModelHandler()
	err = srv.RegisterHandler(handler)
	if err != nil {
		b.Fatalf("Failed to register handler: %v", err)
	}

	// Start server
	err = srv.Start()
	if err != nil {
		b.Fatalf("Failed to start server: %v", err)
	}
	defer srv.Stop()

	// Create and start client
	c := client.New(
		client.WithServerPort(port),
		client.WithConnectionTimeout(5*time.Second),
	)

	err = c.Start()
	if err != nil {
		b.Fatalf("Failed to start client: %v", err)
	}
	defer c.Stop()

	// Ensure connection is established
	if !testutil.WaitForCondition(5*time.Second, 100*time.Millisecond, func() bool {
		return c.Status() == core.StatusRunning
	}) {
		b.Fatal("Client failed to connect to server")
	}

	// Create a request to reuse
	req := core.NewModelRequest()
	req.ModelData["name"] = "Benchmark Test Model"
	req.ModelData["value"] = 42
	req.Parameters = append(req.Parameters, core.Parameter{
		Name:  "param1",
		Value: "value1",
		Type:  "string",
	})

	// Use a long-lived context for the benchmark
	ctx := context.Background()

	// Reset the benchmark timer to exclude setup time
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		resp, err := c.ProcessModel(ctx, req)
		if err != nil {
			b.Fatalf("ProcessModel failed: %v", err)
		}
		if resp == nil || !resp.Success {
			b.Fatalf("Response unsuccessful: %v", resp)
		}
	}
}

// BenchmarkRequestSizes measures performance with different payload sizes.
func BenchmarkRequestSizes(b *testing.B) {
	// Define different payload sizes to test
	payloadSizes := []int{1, 10, 100, 1000, 10000}

	for _, size := range payloadSizes {
		b.Run(fmt.Sprintf("Payload-%dKB", size), func(b *testing.B) {
			// Get a free port for testing
			port, err := testutil.GetFreePort()
			if err != nil {
				b.Fatalf("Failed to get free port: %v", err)
			}

			// Create and start server
			srv := server.New(server.WithPort(port))

			// Register default handler
			handler := server.NewDefaultModelHandler()
			err = srv.RegisterHandler(handler)
			if err != nil {
				b.Fatalf("Failed to register handler: %v", err)
			}

			// Start server
			err = srv.Start()
			if err != nil {
				b.Fatalf("Failed to start server: %v", err)
			}
			defer srv.Stop()

			// Create and start client
			c := client.New(client.WithServerPort(port))
			err = c.Start()
			if err != nil {
				b.Fatalf("Failed to start client: %v", err)
			}
			defer c.Stop()

			// Create a string payload of the specified size (roughly in KB)
			payload := make([]byte, size*1024)
			for i := range payload {
				payload[i] = byte(i % 256)
			}

			// Create a request with the payload
			req := core.NewModelRequest()
			req.ModelData["payload"] = string(payload)

			// Use a background context
			ctx := context.Background()

			// Reset the benchmark timer to exclude setup time
			b.ResetTimer()

			// Run the benchmark
			for i := 0; i < b.N; i++ {
				_, err := c.ProcessModel(ctx, req)
				if err != nil {
					b.Fatalf("ProcessModel failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkConcurrentRequests measures performance with different levels of concurrency.
func BenchmarkConcurrentRequests(b *testing.B) {
	// Define concurrency levels to test
	concurrencyLevels := []int{1, 5, 10, 25, 50, 100}

	for _, concurrency := range concurrencyLevels {
		b.Run(fmt.Sprintf("Concurrency-%d", concurrency), func(b *testing.B) {
			// Get a free port for testing
			port, err := testutil.GetFreePort()
			if err != nil {
				b.Fatalf("Failed to get free port: %v", err)
			}

			// Create and start server with appropriate max clients setting
			srv := server.New(
				server.WithPort(port),
				server.WithMaxConcurrentClients(concurrency*2), // Extra headroom
			)

			// Register default handler
			handler := server.NewDefaultModelHandler()
			err = srv.RegisterHandler(handler)
			if err != nil {
				b.Fatalf("Failed to register handler: %v", err)
			}

			// Start server
			err = srv.Start()
			if err != nil {
				b.Fatalf("Failed to start server: %v", err)
			}
			defer srv.Stop()

			// Create and start client
			c := client.New(client.WithServerPort(port))
			err = c.Start()
			if err != nil {
				b.Fatalf("Failed to start client: %v", err)
			}
			defer c.Stop()

			// Create a standard request
			req := core.NewModelRequest()
			req.ModelData["name"] = "Concurrent Benchmark"
			req.Parameters = append(req.Parameters, core.Parameter{
				Name:  "benchmark",
				Value: "concurrent",
				Type:  "string",
			})

			// Use a background context
			ctx := context.Background()

			// Set parallelism to our concurrency level
			b.SetParallelism(concurrency)

			// Reset the benchmark timer to exclude setup time
			b.ResetTimer()

			// Run the benchmark
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := c.ProcessModel(ctx, req)
					if err != nil {
						b.Fatalf("ProcessModel failed: %v", err)
					}
				}
			})
		})
	}
}

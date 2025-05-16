// Package client provides a client implementation for the Model Context Protocol (MCP).
// It allows applications to connect to MCP servers and make model processing requests.
package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/narcolepticfox/mcp/core"
	"github.com/sourcegraph/jsonrpc2"
)

// Client implements the MCP client that connects to an MCP server.
// It manages the connection, handles request/response communication, and
// provides methods for model processing operations.
type Client struct {
	opts             Options
	status           core.Status
	statusMu         sync.RWMutex
	conn             *jsonrpc2.Conn
	connMu           sync.RWMutex
	callbacks        []func(core.StatusChangeEvent)
	reconnectAttempt int
	isConnected      bool

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// New creates a new MCP client with the given options.
// It applies all provided option functions to configure the client and
// initializes it with default values for all unspecified options.
func New(options ...Option) *Client {
	opts := DefaultOptions()
	for _, opt := range options {
		opt(&opts)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Client{
		opts:      opts,
		status:    core.StatusStopped,
		callbacks: make([]func(core.StatusChangeEvent), 0),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Start connects to the server and starts the client.
// It establishes a connection to the configured server and initializes
// the JSON-RPC communication channel. Returns an error if the client
// is already running or if connection fails.
func (c *Client) Start() error {
	c.statusMu.Lock()
	if c.status != core.StatusStopped {
		c.statusMu.Unlock()
		return fmt.Errorf("cannot start client in %s state", c.status)
	}
	c.updateStatusLocked(core.StatusStarting, nil)
	c.statusMu.Unlock()

	if err := c.connect(); err != nil {
		c.updateStatus(core.StatusFailed, err)
		return err
	}

	c.updateStatus(core.StatusRunning, nil)
	log.Printf("MCP client connected to %s:%d", c.opts.ServerHost, c.opts.ServerPort)

	return nil
}

// connect establishes a TCP connection to the MCP server and sets up the JSON-RPC communication.
// It creates the necessary streams and handlers, and starts a background goroutine to monitor
// the connection status.
func (c *Client) connect() error {
	// Create TCP connection
	addr := fmt.Sprintf("%s:%d", c.opts.ServerHost, c.opts.ServerPort)

	dialer := &net.Dialer{
		Timeout: c.opts.ConnectionTimeout,
	}

	netConn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", addr, err)
	}

	// Create JSON-RPC stream
	stream := jsonrpc2.NewBufferedStream(netConn, jsonrpc2.VSCodeObjectCodec{})

	// Create JSON-RPC handler
	handler := &rpcHandler{client: c}

	// Create JSON-RPC connection
	c.connMu.Lock()
	c.conn = jsonrpc2.NewConn(c.ctx, stream, handler)
	c.isConnected = true
	c.connMu.Unlock()

	// Monitor connection
	c.wg.Add(1)
	go c.monitorConnection()

	return nil
}

func (c *Client) monitorConnection() {
	defer c.wg.Done()

	c.connMu.RLock()
	conn := c.conn
	c.connMu.RUnlock()

	if conn == nil {
		return
	}

	// Wait for disconnection
	<-conn.DisconnectNotify()

	c.connMu.Lock()
	c.isConnected = false
	c.connMu.Unlock()

	log.Printf("Disconnected from server")

	// Handle reconnection if enabled
	if c.opts.AutoReconnect && c.status == core.StatusRunning {
		c.attemptReconnect()
	}
}

func (c *Client) attemptReconnect() {
	for c.reconnectAttempt < c.opts.MaxReconnectAttempts {
		c.reconnectAttempt++

		log.Printf("Attempting to reconnect (%d/%d)...",
			c.reconnectAttempt, c.opts.MaxReconnectAttempts)

		// Wait before reconnecting
		time.Sleep(c.opts.ReconnectDelay)

		// Check if we're shutting down
		select {
		case <-c.ctx.Done():
			return
		default:
			// Continue with reconnection
		}

		if err := c.connect(); err != nil {
			log.Printf("Reconnection attempt failed: %v", err)
		} else {
			log.Printf("Reconnected to server")
			c.reconnectAttempt = 0
			return
		}
	}

	log.Printf("Max reconnection attempts reached")
	c.updateStatus(core.StatusFailed, errors.New("max reconnection attempts reached"))
}

// Stop disconnects from the server and stops the client.
func (c *Client) Stop() error {
	c.statusMu.Lock()
	if c.status != core.StatusRunning {
		c.statusMu.Unlock()
		return fmt.Errorf("cannot stop client in %s state", c.status)
	}
	c.updateStatusLocked(core.StatusStopping, nil)
	c.statusMu.Unlock()

	// Cancel the context to signal shutdown
	c.cancel()

	// Close the connection
	c.connMu.Lock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	c.isConnected = false
	c.connMu.Unlock()

	// Wait for all goroutines to finish
	c.wg.Wait()

	c.updateStatus(core.StatusStopped, nil)
	log.Printf("MCP client stopped")

	return nil
}

// Status returns the current client status.
func (c *Client) Status() core.Status {
	c.statusMu.RLock()
	defer c.statusMu.RUnlock()
	return c.status
}

// IsConnected returns whether the client is currently connected.
func (c *Client) IsConnected() bool {
	c.connMu.RLock()
	defer c.connMu.RUnlock()
	return c.isConnected
}

// OnStatusChange registers a callback for status changes.
func (c *Client) OnStatusChange(callback func(core.StatusChangeEvent)) {
	c.callbacks = append(c.callbacks, callback)
}

// ProcessModel sends a model processing request to the server.
func (c *Client) ProcessModel(ctx context.Context, req *core.ModelRequest) (*core.ModelResponse, error) {
	c.connMu.RLock()
	conn := c.conn
	c.connMu.RUnlock()

	if conn == nil {
		return nil, errors.New("not connected to server")
	}

	var resp core.ModelResponse
	err := conn.Call(ctx, "mcp.processModel", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("RPC error: %w", err)
	}

	return &resp, nil
}

func (c *Client) updateStatus(newStatus core.Status, err error) {
	c.statusMu.Lock()
	defer c.statusMu.Unlock()
	c.updateStatusLocked(newStatus, err)
}

func (c *Client) updateStatusLocked(newStatus core.Status, err error) {
	oldStatus := c.status
	c.status = newStatus

	event := core.StatusChangeEvent{
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Timestamp: time.Now(),
		Error:     err,
	}

	// Notify callbacks
	for _, callback := range c.callbacks {
		go callback(event)
	}
}

// rpcHandler implements jsonrpc2.Handler for the client.
type rpcHandler struct {
	client *Client
}

// Handle handles JSON-RPC requests from the server.
func (h *rpcHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	// Handle notifications or requests from the server
	// In this simplified example, we just log them
	log.Printf("Received request from server: %s", req.Method)

	// We could dispatch to registered handlers here, similar to the server
}

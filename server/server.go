// Package server provides a server implementation for the Model Context Protocol (MCP).
// It allows applications to receive and process model requests using registered handlers.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/narcolepticfox/mcp/core"
	"github.com/sourcegraph/jsonrpc2"
)

// Server implements the MCP server that listens for client connections and
// routes requests to appropriate handlers. It manages the server lifecycle,
// network listeners, and registered method handlers.
type Server struct {
	opts      Options
	status    core.Status
	statusMu  sync.RWMutex
	listeners []net.Listener
	handlers  map[string]interface{}
	callbacks []func(core.StatusChangeEvent)

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// New creates a new MCP server with the given options.
// It applies all provided option functions to configure the server and
// initializes it with default values for all unspecified options.
func New(options ...Option) *Server {
	opts := DefaultOptions()
	for _, opt := range options {
		opt(&opts)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Server{
		opts:      opts,
		status:    core.StatusStopped,
		handlers:  make(map[string]interface{}),
		callbacks: make([]func(core.StatusChangeEvent), 0),
		ctx:       ctx,
		cancel:    cancel,
	}
}

// RegisterHandler registers a handler with the server for processing model requests.
// It registers the handler for all methods it supports, checking for conflicts
// with already registered handlers. Returns an error if a method is already registered.
func (s *Server) RegisterHandler(handler Handler) error {
	for _, method := range handler.Methods() {
		if _, exists := s.handlers[method]; exists {
			return fmt.Errorf("handler for method %s already registered", method)
		}
		s.handlers[method] = handler
	}
	return nil
}

// Start starts the server and begins listening for client connections.
// It creates network listeners based on the configured options and handles
// incoming client connections. Returns an error if the server is already
// running or if it fails to set up the listeners.
func (s *Server) Start() error {
	s.statusMu.Lock()
	if s.status != core.StatusStopped {
		s.statusMu.Unlock()
		return fmt.Errorf("cannot start server in %s state", s.status)
	}
	s.updateStatusLocked(core.StatusStarting, nil)
	s.statusMu.Unlock()

	// Create TCP listener
	addr := fmt.Sprintf("%s:%d", s.opts.Host, s.opts.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		s.updateStatus(core.StatusFailed, err)
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	s.listeners = append(s.listeners, listener)

	// Start accepting connections
	s.wg.Add(1)
	go s.acceptConnections(listener)

	s.updateStatus(core.StatusRunning, nil)
	log.Printf("MCP server listening on %s", addr)

	return nil
}

func (s *Server) acceptConnections(listener net.Listener) {
	defer s.wg.Done()

	for {
		conn, err := listener.Accept()
		if err != nil {
			// Check if we're shutting down
			select {
			case <-s.ctx.Done():
				return
			default:
				log.Printf("Error accepting connection: %v", err)
				continue
			}
		}

		// Handle each connection in a goroutine
		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	log.Printf("Client connected from %s", conn.RemoteAddr())

	// Create JSON-RPC stream
	stream := jsonrpc2.NewBufferedStream(conn, jsonrpc2.VSCodeObjectCodec{})

	// Create JSON-RPC handler
	handler := &rpcHandler{server: s}

	// Create JSON-RPC connection
	rpcConn := jsonrpc2.NewConn(s.ctx, stream, handler)

	// Wait for connection to close
	<-rpcConn.DisconnectNotify()

	log.Printf("Client disconnected from %s", conn.RemoteAddr())
}

// Stop stops the server.
func (s *Server) Stop() error {
	s.statusMu.Lock()
	if s.status != core.StatusRunning {
		s.statusMu.Unlock()
		return fmt.Errorf("cannot stop server in %s state", s.status)
	}
	s.updateStatusLocked(core.StatusStopping, nil)
	s.statusMu.Unlock()

	// Cancel the context to signal shutdown
	s.cancel()

	// Close all listeners
	for _, listener := range s.listeners {
		listener.Close()
	}

	// Wait for all goroutines to finish
	s.wg.Wait()

	s.updateStatus(core.StatusStopped, nil)
	log.Printf("MCP server stopped")

	return nil
}

// Status returns the current server status.
func (s *Server) Status() core.Status {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	return s.status
}

// OnStatusChange registers a callback for status changes.
func (s *Server) OnStatusChange(callback func(core.StatusChangeEvent)) {
	s.callbacks = append(s.callbacks, callback)
}

func (s *Server) updateStatus(newStatus core.Status, err error) {
	s.statusMu.Lock()
	defer s.statusMu.Unlock()
	s.updateStatusLocked(newStatus, err)
}

func (s *Server) updateStatusLocked(newStatus core.Status, err error) {
	oldStatus := s.status
	s.status = newStatus

	event := core.StatusChangeEvent{
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Timestamp: time.Now(),
		Error:     err,
	}

	// Notify callbacks
	for _, callback := range s.callbacks {
		go callback(event)
	}
}

// rpcHandler implements jsonrpc2.Handler.
type rpcHandler struct {
	server *Server
}

// Handle handles JSON-RPC requests.
func (h *rpcHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	// Find the appropriate handler
	handler, ok := h.server.handlers[req.Method]
	if !ok {
		err := conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeMethodNotFound,
			Message: fmt.Sprintf("method not found: %s", req.Method),
		})
		if err != nil {
			log.Printf("Error replying to client: %v", err)
		}
		return
	}

	// Handle the request based on the method
	switch req.Method {
	case "mcp.processModel":
		h.handleProcessModel(ctx, conn, req, handler)
	default:
		err := conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInvalidRequest,
			Message: fmt.Sprintf("unknown method: %s", req.Method),
		})
		if err != nil {
			log.Printf("Error replying to client: %v", err)
		}
	}
}

func (h *rpcHandler) handleProcessModel(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request, handler interface{}) {
	modelHandler, ok := handler.(ModelHandler)
	if !ok {
		err := conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInternalError,
			Message: "handler is not a ModelHandler",
		})
		if err != nil {
			log.Printf("Error replying to client: %v", err)
		}
		return
	}
	// Parse the request
	var modelReq core.ModelRequest
	if err := json.Unmarshal(*req.Params, &modelReq); err != nil {
		err := conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInvalidParams,
			Message: fmt.Sprintf("invalid params: %v", err),
		})
		if err != nil {
			log.Printf("Error replying to client: %v", err)
		}
		return
	}

	// Process the request
	resp, err := modelHandler.ProcessModel(ctx, &modelReq)
	if err != nil {
		err := conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInternalError,
			Message: fmt.Sprintf("processing error: %v", err),
		})
		if err != nil {
			log.Printf("Error replying to client: %v", err)
		}
		return
	}

	// Send the response
	if err := conn.Reply(ctx, req.ID, resp); err != nil {
		log.Printf("Error replying to client: %v", err)
	}
}

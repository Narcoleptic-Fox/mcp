// Package client provides a client implementation for the Model Context Protocol (MCP).
package client

import "time"

// Options holds configuration parameters for the MCP client.
// It defines connection settings, reconnection behavior, and security options.
type Options struct {
	ServerHost           string        // Hostname or IP address of the MCP server
	ServerPort           int           // TCP port of the MCP server
	ConnectionTimeout    time.Duration // Timeout for establishing a connection
	AutoReconnect        bool          // Whether to automatically attempt reconnection on disconnect
	MaxReconnectAttempts int           // Maximum number of reconnection attempts before giving up
	ReconnectDelay       time.Duration // Time to wait between reconnection attempts
	EnableTLS            bool          // Whether to use TLS for server connections
}

// DefaultOptions returns the default client options.
// These defaults provide sensible starting values that work for most local deployments,
// with automatic reconnection enabled but limited to 3 attempts.
func DefaultOptions() Options {
	return Options{
		ServerHost:           "localhost",
		ServerPort:           5000,
		ConnectionTimeout:    30 * time.Second,
		AutoReconnect:        true,
		MaxReconnectAttempts: 3,
		ReconnectDelay:       time.Second,
		EnableTLS:            false,
	}
}

// Option is a function type that modifies Options.
// It implements the functional options pattern for configuring the client.
type Option func(*Options)

// WithServerHost sets the hostname or IP address of the MCP server.
// This can be a domain name, IPv4 address, or IPv6 address.
func WithServerHost(host string) Option {
	return func(o *Options) {
		o.ServerHost = host
	}
}

// WithServerPort sets the TCP port number of the MCP server.
// This must match the port the server is listening on.
func WithServerPort(port int) Option {
	return func(o *Options) {
		o.ServerPort = port
	}
}

// WithConnectionTimeout sets the maximum time to wait when connecting to the server.
// If the connection isn't established within this time, the attempt is aborted.
func WithConnectionTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.ConnectionTimeout = timeout
	}
}

// WithAutoReconnect enables or disables automatic reconnection.
func WithAutoReconnect(enable bool) Option {
	return func(o *Options) {
		o.AutoReconnect = enable
	}
}

// WithMaxReconnectAttempts sets the maximum number of reconnection attempts.
func WithMaxReconnectAttempts(max int) Option {
	return func(o *Options) {
		o.MaxReconnectAttempts = max
	}
}

// WithReconnectDelay sets the delay between reconnection attempts.
func WithReconnectDelay(delay time.Duration) Option {
	return func(o *Options) {
		o.ReconnectDelay = delay
	}
}

// WithTLS enables TLS.
func WithTLS() Option {
	return func(o *Options) {
		o.EnableTLS = true
	}
}

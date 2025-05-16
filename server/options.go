// Package server provides a server implementation for the Model Context Protocol (MCP).
package server

import "time"

// Options holds configuration parameters for the MCP server.
// It defines network settings, connection limits, timeouts, and TLS configuration.
type Options struct {
	Host                 string        // Network interface to bind to, e.g., "127.0.0.1" for localhost only
	Port                 int           // TCP port to listen on
	MaxConcurrentClients int           // Maximum number of simultaneous client connections
	ConnectionTimeout    time.Duration // Time limit for establishing connections
	EnableTLS            bool          // Whether to use TLS encryption for connections
	CertificatePath      string        // Path to the TLS certificate file when TLS is enabled
	CertificateKeyPath   string        // Path to the TLS certificate key file when TLS is enabled
}

// DefaultOptions returns the default server options.
// These defaults provide a reasonable starting point for most deployments
// with the server listening only on localhost.
func DefaultOptions() Options {
	return Options{
		Host:                 "127.0.0.1",
		Port:                 5000,
		MaxConcurrentClients: 10,
		ConnectionTimeout:    30 * time.Second,
		EnableTLS:            false,
	}
}

// Option is a function type that modifies Options.
// It implements the functional options pattern for configuring the server.
type Option func(*Options)

// WithHost sets the host address for the server to bind to.
// Use "0.0.0.0" to listen on all interfaces, or a specific IP to restrict access.
func WithHost(host string) Option {
	return func(o *Options) {
		o.Host = host
	}
}

// WithPort sets the TCP port number for the server to listen on.
// The port must be available and the process must have permission to bind to it.
func WithPort(port int) Option {
	return func(o *Options) {
		o.Port = port
	}
}

// WithMaxConcurrentClients sets the maximum number of concurrent client connections.
// This helps prevent resource exhaustion by limiting the number of simultaneous
// connections the server will accept.
func WithMaxConcurrentClients(max int) Option {
	return func(o *Options) {
		o.MaxConcurrentClients = max
	}
}

// WithConnectionTimeout sets the connection timeout.
func WithConnectionTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.ConnectionTimeout = timeout
	}
}

// WithTLS enables TLS with the specified certificate and key.
func WithTLS(certPath, keyPath string) Option {
	return func(o *Options) {
		o.EnableTLS = true
		o.CertificatePath = certPath
		o.CertificateKeyPath = keyPath
	}
}

// Add your option functions here, e.g. WithHost, WithPort, etc.

func WithCertificatePath(path string) Option {
	return func(o *Options) {
		o.CertificatePath = path
	}
}

// Add your option functions here, e.g. WithHost, WithPort, etc.

func WithCertificateKeyPath(path string) Option {
	return func(o *Options) {
		o.CertificateKeyPath = path
	}
}

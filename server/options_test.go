package server

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServerDefaultOptions(t *testing.T) {
	options := DefaultOptions()

	// Check that default options are set correctly
	assert.Equal(t, "127.0.0.1", options.Host, "Default Host should be 127.0.0.1")
	assert.Equal(t, 5000, options.Port, "Default Port should be 5000")
	assert.Equal(t, 10, options.MaxConcurrentClients, "Default MaxConcurrentClients should be 10")
	assert.Equal(t, 30*time.Second, options.ConnectionTimeout, "Default ConnectionTimeout should be 30s")
	assert.False(t, options.EnableTLS, "Default EnableTLS should be false")
	assert.Empty(t, options.CertificatePath, "Default CertificatePath should be empty")
	assert.Empty(t, options.CertificateKeyPath, "Default CertificateKeyPath should be empty")
}

func TestWithHost(t *testing.T) {
	options := DefaultOptions()
	option := WithHost("0.0.0.0")
	option(&options)

	assert.Equal(t, "0.0.0.0", options.Host, "Host should be updated")
}

func TestWithPort(t *testing.T) {
	options := DefaultOptions()
	option := WithPort(9999)
	option(&options)

	assert.Equal(t, 9999, options.Port, "Port should be updated")
}

func TestWithMaxConcurrentClients(t *testing.T) {
	options := DefaultOptions()
	option := WithMaxConcurrentClients(50)
	option(&options)

	assert.Equal(t, 50, options.MaxConcurrentClients, "MaxConcurrentClients should be updated")
}

func TestWithConnectionTimeout(t *testing.T) {
	options := DefaultOptions()
	timeout := 10 * time.Second
	option := WithConnectionTimeout(timeout)
	option(&options)

	assert.Equal(t, timeout, options.ConnectionTimeout, "ConnectionTimeout should be updated")
}

func TestWithTLS(t *testing.T) {
	options := DefaultOptions()
	option := WithTLS(true)
	option(&options)

	assert.True(t, options.EnableTLS, "EnableTLS should be updated")
}

func TestWithCertificatePath(t *testing.T) {
	options := DefaultOptions()
	path := "/path/to/cert.pem"
	option := WithCertificatePath(path)
	option(&options)

	assert.Equal(t, path, options.CertificatePath, "CertificatePath should be updated")
}

func TestWithCertificateKeyPath(t *testing.T) {
	options := DefaultOptions()
	path := "/path/to/key.pem"
	option := WithCertificateKeyPath(path)
	option(&options)

	assert.Equal(t, path, options.CertificateKeyPath, "CertificateKeyPath should be updated")
}

func TestServerOptionChaining(t *testing.T) {
	// Test applying multiple options
	server := New(
		WithHost("0.0.0.0"),
		WithPort(8888),
		WithMaxConcurrentClients(20),
		WithConnectionTimeout(15*time.Second),
	)

	// Extract options from server for testing
	options := server.options

	assert.Equal(t, "0.0.0.0", options.Host, "Host should be updated")
	assert.Equal(t, 8888, options.Port, "Port should be updated")
	assert.Equal(t, 20, options.MaxConcurrentClients, "MaxConcurrentClients should be updated")
	assert.Equal(t, 15*time.Second, options.ConnectionTimeout, "ConnectionTimeout should be updated")
}

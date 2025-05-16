package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultOptions(t *testing.T) {
	options := DefaultOptions()

	// Check that default options are set correctly
	assert.Equal(t, "localhost", options.ServerHost, "Default ServerHost should be localhost")
	assert.Equal(t, 5000, options.ServerPort, "Default ServerPort should be 5000")
	assert.Equal(t, 30*time.Second, options.ConnectionTimeout, "Default ConnectionTimeout should be 30s")
	assert.True(t, options.AutoReconnect, "Default AutoReconnect should be true")
	assert.Equal(t, 3, options.MaxReconnectAttempts, "Default MaxReconnectAttempts should be 3")
	assert.Equal(t, time.Second, options.ReconnectDelay, "Default ReconnectDelay should be 1s")
	assert.False(t, options.EnableTLS, "Default EnableTLS should be false")
}

func TestWithServerHost(t *testing.T) {
	options := DefaultOptions()
	option := WithServerHost("test-host")
	option(&options)

	assert.Equal(t, "test-host", options.ServerHost, "ServerHost should be updated")
}

func TestWithServerPort(t *testing.T) {
	options := DefaultOptions()
	option := WithServerPort(9999)
	option(&options)

	assert.Equal(t, 9999, options.ServerPort, "ServerPort should be updated")
}

func TestWithConnectionTimeout(t *testing.T) {
	options := DefaultOptions()
	timeout := 10 * time.Second
	option := WithConnectionTimeout(timeout)
	option(&options)

	assert.Equal(t, timeout, options.ConnectionTimeout, "ConnectionTimeout should be updated")
}

func TestWithAutoReconnect(t *testing.T) {
	options := DefaultOptions()
	option := WithAutoReconnect(false)
	option(&options)

	assert.False(t, options.AutoReconnect, "AutoReconnect should be updated")
}

func TestWithMaxReconnectAttempts(t *testing.T) {
	options := DefaultOptions()
	option := WithMaxReconnectAttempts(10)
	option(&options)

	assert.Equal(t, 10, options.MaxReconnectAttempts, "MaxReconnectAttempts should be updated")
}

func TestWithReconnectDelay(t *testing.T) {
	options := DefaultOptions()
	delay := 5 * time.Second
	option := WithReconnectDelay(delay)
	option(&options)

	assert.Equal(t, delay, options.ReconnectDelay, "ReconnectDelay should be updated")
}

func TestWithTLS(t *testing.T) {
	options := DefaultOptions()
	option := WithTLS(true)
	option(&options)

	assert.True(t, options.EnableTLS, "EnableTLS should be updated")
}

func TestOptionChaining(t *testing.T) {
	// Test applying multiple options
	client := New(
		WithServerHost("custom-host"),
		WithServerPort(8888),
		WithConnectionTimeout(15*time.Second),
		WithAutoReconnect(false),
	)

	// Extract options from client for testing
	options := client.options

	assert.Equal(t, "custom-host", options.ServerHost, "ServerHost should be updated")
	assert.Equal(t, 8888, options.ServerPort, "ServerPort should be updated")
	assert.Equal(t, 15*time.Second, options.ConnectionTimeout, "ConnectionTimeout should be updated")
	assert.False(t, options.AutoReconnect, "AutoReconnect should be updated")
}

// Package core provides the fundamental models and interfaces for the Model Context Protocol (MCP).
package core

import "time"

// Status represents the operational status of an MCP component.
// It uses enumerated values to indicate the component's current state.
type Status int

const (
	// StatusStopped indicates the component is not running.
	StatusStopped Status = iota

	// StatusStarting indicates the component is in the process of starting.
	StatusStarting

	// StatusRunning indicates the component is operational.
	StatusRunning

	// StatusStopping indicates the component is shutting down.
	StatusStopping

	// StatusFailed indicates the component encountered an error.
	StatusFailed
)

// String returns a string representation of the status.
// This implements the Stringer interface for the Status type.
func (s Status) String() string {
	return [...]string{"Stopped", "Starting", "Running", "Stopping", "Failed"}[s]
}

// StatusChangeEvent represents a status change notification.
// It contains the previous and new status, the time of the change, and any associated error.
type StatusChangeEvent struct {
	OldStatus Status    // Status before the change
	NewStatus Status    // Status after the change
	Timestamp time.Time // When the status change occurred
	Error     error     // Error that caused the status change, if any
}

// Component defines the interface for MCP components.
// All components in the MCP system must implement these methods
// to provide consistent lifecycle management and status reporting.
type Component interface {
	// Start initializes the component and begins its operation.
	// Returns an error if the component fails to start.
	Start() error

	// Stop terminates the component's operation in a graceful manner.
	// Returns an error if the component fails to stop properly.
	Stop() error

	// Status returns the current operational status of the component.
	Status() Status

	// OnStatusChange registers a callback function to be called when the component's status changes.
	// The callback receives a StatusChangeEvent containing details about the change.
	OnStatusChange(callback func(StatusChangeEvent))
}

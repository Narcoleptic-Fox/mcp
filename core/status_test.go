package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusString(t *testing.T) {
	// Test that Status values convert to strings correctly
	cases := []struct {
		status Status
		want   string
	}{
		{StatusStopped, "Stopped"},
		{StatusStarting, "Starting"},
		{StatusRunning, "Running"},
		{StatusStopping, "Stopping"},
		{StatusFailed, "Failed"},
	}

	for _, c := range cases {
		t.Run(string(c.status), func(t *testing.T) {
			got := string(c.status)
			assert.Equal(t, c.want, got, "Status string representation should match expected value")
		})
	}
}

// MockComponent implements the Component interface for testing
type MockComponent struct {
	status       Status
	callbacks    []func(StatusChangeEvent)
	startErr     error
	stopErr      error
	startCalled  bool
	stopCalled   bool
	statusCalled bool
}

func NewMockComponent(initialStatus Status) *MockComponent {
	return &MockComponent{
		status:    initialStatus,
		callbacks: make([]func(StatusChangeEvent), 0),
	}
}

func (m *MockComponent) Start() error {
	m.startCalled = true
	if m.startErr != nil {
		return m.startErr
	}
	oldStatus := m.status
	m.status = StatusRunning
	m.notifyStatusChange(oldStatus, m.status, nil)
	return nil
}

func (m *MockComponent) Stop() error {
	m.stopCalled = true
	if m.stopErr != nil {
		return m.stopErr
	}
	oldStatus := m.status
	m.status = StatusStopped
	m.notifyStatusChange(oldStatus, m.status, nil)
	return nil
}

func (m *MockComponent) Status() Status {
	m.statusCalled = true
	return m.status
}

func (m *MockComponent) OnStatusChange(callback func(StatusChangeEvent)) {
	m.callbacks = append(m.callbacks, callback)
}

func (m *MockComponent) notifyStatusChange(oldStatus, newStatus Status, err error) {
	event := StatusChangeEvent{
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Error:     err,
	}
	for _, callback := range m.callbacks {
		callback(event)
	}
}

func TestStatusChangeEvent(t *testing.T) {
	// Test status change notifications
	mockComponent := NewMockComponent(StatusIdle)

	// Track status changes
	var receivedEvents []StatusChangeEvent
	mockComponent.OnStatusChange(func(event StatusChangeEvent) {
		receivedEvents = append(receivedEvents, event)
	})

	// Test Start causing a status change
	err := mockComponent.Start()
	assert.NoError(t, err, "Start should not return an error")
	assert.Equal(t, StatusRunning, mockComponent.Status(), "Status should be running after start")
	assert.True(t, mockComponent.startCalled, "Start method should be called")

	// Test Stop causing a status change
	err = mockComponent.Stop()
	assert.NoError(t, err, "Stop should not return an error")
	assert.Equal(t, StatusIdle, mockComponent.Status(), "Status should be idle after stop")
	assert.True(t, mockComponent.stopCalled, "Stop method should be called")

	// Verify that we received both status change events
	assert.Len(t, receivedEvents, 2, "Should have received 2 status change events")

	// Verify the first event (idle -> running)
	assert.Equal(t, StatusIdle, receivedEvents[0].OldStatus, "First event old status should be idle")
	assert.Equal(t, StatusRunning, receivedEvents[0].NewStatus, "First event new status should be running")
	assert.Nil(t, receivedEvents[0].Error, "First event should not have an error")

	// Verify the second event (running -> idle)
	assert.Equal(t, StatusRunning, receivedEvents[1].OldStatus, "Second event old status should be running")
	assert.Equal(t, StatusIdle, receivedEvents[1].NewStatus, "Second event new status should be idle")
	assert.Nil(t, receivedEvents[1].Error, "Second event should not have an error")
}

func TestComponentErrors(t *testing.T) {
	// Test error handling in Start() and Stop()
	startErr := errors.New("start error")
	stopErr := errors.New("stop error")

	mockComponent := NewMockComponent(StatusIdle)
	mockComponent.startErr = startErr
	mockComponent.stopErr = stopErr

	// Track status changes
	var receivedEvents []StatusChangeEvent
	mockComponent.OnStatusChange(func(event StatusChangeEvent) {
		receivedEvents = append(receivedEvents, event)
	})

	// Test Start with error
	err := mockComponent.Start()
	assert.Error(t, err, "Start should return an error")
	assert.Equal(t, startErr, err, "Start should return the expected error")

	// Test Stop with error
	err = mockComponent.Stop()
	assert.Error(t, err, "Stop should return an error")
	assert.Equal(t, stopErr, err, "Stop should return the expected error")

	// Verify that no status change events were received (since the mock doesn't change status on error)
	assert.Len(t, receivedEvents, 0, "Should not have received any status change events")
}

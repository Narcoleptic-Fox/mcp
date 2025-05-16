package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidationResult(t *testing.T) {
	result := NewValidationResult()

	assert.True(t, result.Valid, "New validation result should be valid")
	assert.NotNil(t, result.Errors, "Errors slice should be initialized")
	assert.Len(t, result.Errors, 0, "Errors slice should be empty")
}

func TestAddError(t *testing.T) {
	result := NewValidationResult()

	// Add a single error
	result.AddError("field1", "error message 1")

	assert.False(t, result.Valid, "Result should be invalid after adding an error")
	assert.Len(t, result.Errors, 1, "Should have one error")
	assert.Equal(t, "field1", result.Errors[0].Field, "Error should have correct field name")
	assert.Equal(t, "error message 1", result.Errors[0].Message, "Error should have correct message")

	// Add another error
	result.AddError("field2", "error message 2")

	assert.False(t, result.Valid, "Result should still be invalid")
	assert.Len(t, result.Errors, 2, "Should have two errors")
	assert.Equal(t, "field2", result.Errors[1].Field, "Second error should have correct field name")
	assert.Equal(t, "error message 2", result.Errors[1].Message, "Second error should have correct message")
}

func TestErrorMethod(t *testing.T) {
	// Test with valid result
	validResult := NewValidationResult()
	err := validResult.Error()
	assert.Nil(t, err, "Valid result should return nil error")

	// Test with a single error
	result := NewValidationResult()
	result.AddError("field1", "message1")

	err = result.Error()
	assert.NotNil(t, err, "Invalid result should return an error")
	assert.Contains(t, err.Error(), "validation failed:", "Error should start with 'validation failed:'")
	assert.Contains(t, err.Error(), "field1: message1", "Error should contain the field and message")

	// Test with multiple errors
	result = NewValidationResult()
	result.AddError("field1", "message1")
	result.AddError("field2", "message2")

	err = result.Error()
	assert.NotNil(t, err, "Invalid result should return an error")
	assert.Contains(t, err.Error(), "field1: message1", "Error should contain the first field and message")
	assert.Contains(t, err.Error(), "field2: message2", "Error should contain the second field and message")
}

func TestValidationError(t *testing.T) {
	// Test creating ValidationError directly
	ve := ValidationError{
		Field:   "testField",
		Message: "test message",
	}

	assert.Equal(t, "testField", ve.Field, "Field should be correctly set")
	assert.Equal(t, "test message", ve.Message, "Message should be correctly set")
}

func TestValidatorImplementation(t *testing.T) {
	// This just tests that the Validator type exists. Actual validation methods would be added here.
	validator := Validator{}

	// Since Validator is currently empty, we just verify it can be instantiated
	assert.NotNil(t, validator, "Should be able to instantiate a Validator")
}

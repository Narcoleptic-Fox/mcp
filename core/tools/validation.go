// Package tools provides utility functions and types for MCP implementations.
package tools

import (
	"errors"
	"fmt"
)

// ValidationResult contains the result of a validation operation.
// It tracks whether validation passed and collects any validation errors.
type ValidationResult struct {
	Valid  bool              // Whether validation passed
	Errors []ValidationError // Collection of validation errors
}

// ValidationError represents a specific error found during validation.
// It identifies both the field that failed validation and the reason.
type ValidationError struct {
	Field   string // Name of the field that failed validation
	Message string // Description of why validation failed
}

// NewValidationResult creates a new, valid validation result with no errors.
// This is typically used at the start of a validation process.
func NewValidationResult() *ValidationResult {
	return &ValidationResult{
		Valid:  true,
		Errors: make([]ValidationError, 0),
	}
}

// AddError adds a validation error to the result and marks it as invalid.
// The field parameter identifies what was being validated, and message
// explains why validation failed.
func (vr *ValidationResult) AddError(field, message string) {
	vr.Valid = false
	vr.Errors = append(vr.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// Error returns all validation errors as a single error object.
// If validation passed, it returns nil. Otherwise, it returns an error
// containing all field-specific error messages.
func (vr *ValidationResult) Error() error {
	if vr.Valid {
		return nil
	}

	errorMsg := "validation failed:"
	for _, err := range vr.Errors {
		errorMsg += fmt.Sprintf(" %s: %s;", err.Field, err.Message)
	}
	return errors.New(errorMsg)
}

// Validator provides methods for validating MCP data structures.
// It contains reusable validation logic that can be applied to various objects.
type Validator struct{}

// NewValidator creates a new validator.
func NewValidator() *Validator {
	return &Validator{}
}

// Validate validates a struct and returns a validation result.
func (v *Validator) Validate(obj interface{}) *ValidationResult {
	result := NewValidationResult()

	// Basic validation: check if nil
	if obj == nil {
		result.AddError("object", "cannot be nil")
		return result
	}

	// In a real implementation, we'd add more validation logic here
	// Example: validate required fields, data types, etc.
	// This would use reflection to examine struct fields and tags

	return result
}

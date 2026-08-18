package server

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the go-playground/validator instance.
type Validator struct {
	v *validator.Validate
}

// ValidationError is a user-facing error produced when request validation
// fails. It is recognized by MapError.
type ValidationError struct {
	Fields map[string]string
}

// Error implements error with a summary message.
func (e *ValidationError) Error() string {
	if len(e.Fields) == 0 {
		return "validation failed"
	}
	parts := make([]string, 0, len(e.Fields))
	for field, msg := range e.Fields {
		parts = append(parts, field+": "+msg)
	}
	return "validation failed: " + strings.Join(parts, "; ")
}

// NewValidator creates a Validator with common custom rules registered.
func NewValidator() *Validator {
	v := validator.New()
	return &Validator{v: v}
}

// ValidateStruct validates a struct and returns a *ValidationError listing
// field-level messages. Returns nil when the input is valid.
func (val *Validator) ValidateStruct(s interface{}) error {
	if val == nil || val.v == nil {
		return nil
	}
	err := val.v.Struct(s)
	if err == nil {
		return nil
	}

	fields := map[string]string{}
	for _, fe := range err.(validator.ValidationErrors) {
		fields[fe.Field()] = messageForTag(fe.Tag(), fe.Param())
	}
	return &ValidationError{Fields: fields}
}

// messageForTag maps a validator tag to a human-readable message.
func messageForTag(tag, param string) string {
	switch tag {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s", param)
	case "max":
		return fmt.Sprintf("must be at most %s", param)
	case "oneof":
		return fmt.Sprintf("must be one of: %s", param)
	case "len":
		return fmt.Sprintf("must be exactly %s characters", param)
	default:
		return "is invalid"
	}
}

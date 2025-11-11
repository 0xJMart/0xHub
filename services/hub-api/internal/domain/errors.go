package domain

import (
	"errors"
	"strings"
)

// ValidationError captures a single field-level validation issue.
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	if e.Field == "" {
		return e.Message
	}
	return e.Field + ": " + e.Message
}

// ValidationErrors aggregates multiple validation issues.
type ValidationErrors []ValidationError

// Error implements the error interface.
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

// HasErrors returns true when the collection has at least one entry.
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// AsError normalises the validation errors to a single error value or nil.
func (ve ValidationErrors) AsError() error {
	if !ve.HasErrors() {
		return nil
	}
	errs := make([]error, len(ve))
	for i, v := range ve {
		errs[i] = v
	}
	return errors.Join(errs...)
}

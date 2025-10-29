package errors

import (
	"github.com/steve-phan/page-insight-tool/internal/models"
)

// ErrorMapper handles conversion from domain errors to HTTP responses
type ErrorMapper struct{}

// NewErrorMapper creates a new error mapper
func NewErrorMapper() *ErrorMapper {
	return &ErrorMapper{}
}

// MapToHTTPError converts domain errors to HTTP error responses
// This function now relies entirely on type-safe domain errors instead of string matching
func (m *ErrorMapper) MapToHTTPError(err error) *models.HTTPError {
	if err == nil {
		return nil
	}

	// All errors should now be domain errors with proper types
	var domainErr *DomainError
	if AsError(err, &domainErr) {
		return &models.HTTPError{
			Message:   domainErr.Message,
			Type:      string(domainErr.Type),
			Code:      domainErr.StatusCode,
			Details:   domainErr.Details,
			RequestID: "", // Will be set by middleware
		}
	}

	// Fallback for any non-domain errors (should be rare in production)
	// This handles edge cases where we might receive unexpected error types
	return &models.HTTPError{
		Message: "Internal server error occurred",
		Type:    string(ErrorTypeInternal),
		Code:    500,
		Details: map[string]interface{}{
			"original_error": err.Error(),
		},
	}
}

// AsError wraps errors.As for easier testing and consistency
func AsError(err error, target interface{}) bool {
	// Simple type assertion for domain errors
	// This can be extended with additional logic if needed
	if domainErr, ok := err.(*DomainError); ok && target != nil {
		if targetPtr, ok := target.(**DomainError); ok {
			*targetPtr = domainErr
			return true
		}
	}
	return false
}
